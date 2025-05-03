// context/AuthContext.tsx
import React, { createContext, useState, ReactNode, useEffect } from 'react';
import { User, AuthState, UserLoginCredentials, UserRegistrationData } from '@/types/user.model';
import authService, { login as loginApi, register as registerApi } from '@/services/auth/auth.service';
import userService from '@/services/user/user.service';

// Initial state
const initialState: AuthState = {
  user: null,
  accessToken: null,
  isAuthenticated: false,
  isLoading: false,
  error: null
};

// Create the context
export const AuthContext = createContext<{
  authState: AuthState;
  login: (email: string, password: string) => Promise<{ success: boolean, error?: string }>;
  register: (userData: UserRegistrationData) => Promise<boolean>;
  logout: () => void;
  updateProfile: (userData: Partial<User>) => Promise<void>;
}>({
  authState: initialState,
  login: async () => ({ success: false, error: 'Not implemented' }),
  register: async () => false,
  logout: () => { },
  updateProfile: async () => { }
});

// Create provider
export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [authState, setAuthState] = useState<AuthState>(() => {
    // Check local storage for existing session
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');

    if (token && userStr) {
      try {
        const user = JSON.parse(userStr);
        console.log("Found existing session in localStorage:", { user, hasToken: !!token });
        return {
          user: user,
          accessToken: token,
          isAuthenticated: true,
          isLoading: false,
          error: null
        };
      } catch (error) {
        console.error('Error parsing stored user data:', error);
        return initialState;
      }
    }

    return initialState;
  });

  // Kiểm tra token khi component mount
  useEffect(() => {
    // Start the automatic token refresh
    authService.setupAutoRefresh();

    // Cleanup on unmount
    return () => {
      authService.stopAutoRefresh();
    };
  }, []);

  // Thêm biến để theo dõi xem đã kiểm tra token hay chưa
  const [sessionValidated, setSessionValidated] = useState(false);

  // Run token validation when component mounts
  useEffect(() => {
    const validateSession = async () => {
      // Tránh kiểm tra lặp lại nếu đã kiểm tra trong session này
      if (sessionValidated) {
        console.log("Session already validated, skipping");
        return;
      }

      // Kiểm tra ratelimit
      const lastValidationTime = parseInt(localStorage.getItem('lastTokenValidation') || '0');
      const currentTime = Date.now();
      const minValidationInterval = 2 * 60 * 1000; // 2 phút

      if (currentTime - lastValidationTime < minValidationInterval) {
        console.log("Skipping validation due to rate limiting");
        setSessionValidated(true);
        return;
      }

      console.log("Validating session on app load...");
      localStorage.setItem('lastTokenValidation', currentTime.toString());

      // Nếu đã có token, kiểm tra tính hợp lệ
      if (authState.accessToken) {
        try {
          setAuthState(prev => ({ ...prev, isLoading: true }));

          // Kiểm tra token
          const validateResult = await authService.validateToken();

          // Kiểm tra các claims để debug
          if (validateResult.claims) {
            // Nếu có role trong claims, lưu tạm vào localStorage
            if (validateResult.claims.role) {
              const userStr = localStorage.getItem('user');
              if (userStr) {
                try {
                  const userData = JSON.parse(userStr);
                  userData.role = validateResult.claims.role;
                  localStorage.setItem('user', JSON.stringify(userData));
                } catch (e) {
                  console.error("Failed to update role from claims", e);
                }
              }
            }
          }

          if (validateResult.valid) {
            // Nếu token hợp lệ, lấy thông tin user đầy đủ
            try {
              const userId = validateResult.user_id;
              const userData = await userService.getUserById(userId) as User;

              // Đảm bảo role từ API /user/me được sử dụng
              const updatedUserData = {
                ...userData,
                role: userData.role || (authState.user?.role || 'user') // Ưu tiên role từ API
              };

              setAuthState({
                user: updatedUserData,
                accessToken: authState.accessToken,
                isAuthenticated: true,
                isLoading: false,
                error: null
              });

              // Make sure localStorage has the latest user data
              localStorage.setItem('user', JSON.stringify(updatedUserData));
              console.log("Auth state updated with fresh user data, role:", updatedUserData.role);
            } catch (userError) {
              console.error('Error fetching user data:', userError);
              // Nếu không lấy được thông tin chi tiết, vẫn giữ thông tin cơ bản
              setAuthState(prev => ({
                ...prev,
                isLoading: false,
                isAuthenticated: true // Still authenticated even if we couldn't get full user details
              }));
              console.log("Using basic user data from localStorage due to fetch error");
            }
          } else {
            console.log("Token invalid, attempting refresh...");
            // Token không hợp lệ, thử làm mới token
            try {
              const refreshResult = await authService.refreshToken();
              console.log("Token refresh successful:", refreshResult);

              // Get user again with the new token
              try {
                // Use a safe fallback for userId since refreshResult.user_id might not exist
                const userId = authState.user?.id;
                if (!userId) {
                  throw new Error('Cannot get user data after refresh: missing user ID');
                }

                console.log("Fetching user data after token refresh for ID:", userId);
                const userData = await userService.getUserById(userId) as User;

                setAuthState({
                  user: userData,
                  accessToken: refreshResult.access_token,
                  isAuthenticated: true,
                  isLoading: false,
                  error: null
                });

                // Update user in localStorage
                localStorage.setItem('user', JSON.stringify(userData));
                console.log("Auth state updated after token refresh");
              } catch (userError) {
                console.error('Error fetching user data after token refresh:', userError);
                // Keep using existing user data
                setAuthState({
                  user: authState.user,
                  accessToken: refreshResult.access_token,
                  isAuthenticated: true,
                  isLoading: false,
                  error: null
                });
              }
            } catch (refreshError) {
              console.error('Error refreshing token:', refreshError);
              // Nếu không làm mới được, đăng xuất
              console.log("Token refresh failed, logging out");
              authService.logout();
              setAuthState(initialState);
            }
          }
        } catch (error) {
          console.error('Session validation error:', error);
          // Nếu kiểm tra token thất bại, thử làm mới token trước khi đăng xuất
          try {
            console.log("Validation failed, attempting token refresh as fallback");
            const refreshResult = await authService.refreshToken();
            console.log("Emergency token refresh successful");

            setAuthState(prev => ({
              ...prev,
              accessToken: refreshResult.access_token,
              isAuthenticated: true,
              isLoading: false
            }));
          } catch (refreshError) {
            console.error('Emergency token refresh also failed:', refreshError);
            // Now we finally log out
            console.log("All authentication recovery attempts failed, logging out");
            authService.logout();
            setAuthState(initialState);
          }
        }
      } else {
        console.log("No token found in state, checking localStorage directly");
        // Double-check localStorage directly in case state wasn't initialized properly
        const token = localStorage.getItem('token');
        const userStr = localStorage.getItem('user');

        if (token && userStr) {
          try {
            console.log("Found token in localStorage but not in state, reinitializing");
            const user = JSON.parse(userStr);
            setAuthState({
              user: user,
              accessToken: token,
              isAuthenticated: true,
              isLoading: false,
              error: null
            });

            // Now we have a token in state, so we'll re-run this effect
            // No need to validate here as the effect will run again
          } catch (error) {
            console.error('Error parsing user from localStorage in double-check:', error);
          }
        }
      }

      // Đánh dấu đã validate để tránh kiểm tra liên tục
      setSessionValidated(true);
    };

    validateSession();
  }, [authState.accessToken, sessionValidated]); // Thêm sessionValidated vào dependencies

  // Login function
  const login = async (email: string, password: string): Promise<{ success: boolean, error?: string }> => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));

      // Gọi API đăng nhập
      const loginResult = await loginApi({ email, password });

      if (loginResult.success) {
        try {
          const userData = await userService.getUserById(loginResult.user_id);
          const userWithProfile = {
            ...userData,
            role: userData.role || loginResult.user_info.role || 'user',
            profile: userData.profile || {
              firstName: loginResult.user_info.first_name,
              lastName: loginResult.user_info.last_name
            }
          };
          setAuthState({
            user: userWithProfile as User,
            accessToken: loginResult.access_token,
            isAuthenticated: true,
            isLoading: false,
            error: null
          });
          localStorage.setItem('user', JSON.stringify(userWithProfile));
          window.dispatchEvent(new Event('storage'));
        } catch (userError) {
          const basicUserInfo = {
            id: loginResult.user_id,
            email: loginResult.user_info.email,
            role: loginResult.user_info.role || 'user',
            username: loginResult.user_info.username,
            profile: {
              firstName: loginResult.user_info.first_name,
              lastName: loginResult.user_info.last_name
            }
          };
          setAuthState({
            user: basicUserInfo as unknown as User,
            accessToken: loginResult.access_token,
            isAuthenticated: true,
            isLoading: false,
            error: null
          });
          window.dispatchEvent(new Event('storage'));
        }
        return { success: true };
      } else {
        setAuthState(prev => ({
          ...prev,
          isLoading: false,
          error: loginResult.message || 'Authentication failed'
        }));
        return { success: false, error: loginResult.message || 'Authentication failed' };
      }
    } catch (error: any) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: error.message || 'Authentication failed. Please check your credentials.'
      }));
      return { success: false, error: error.message || 'Authentication failed. Please check your credentials.' };
    }
  };

  // Register function
  const register = async (userData: UserRegistrationData) => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));

      // Gọi API đăng ký
      const registerResult = await registerApi(userData);

      if (registerResult.success) {
        // Không tự động đăng nhập sau khi đăng ký, chỉ cập nhật trạng thái
        setAuthState(prev => ({
          ...prev,
          isLoading: false,
          // Đặt thông báo thành công để hiển thị cho người dùng
          error: null
        }));

        // Trả về true để component gọi hàm này biết đăng ký thành công
        return true;
      } else {
        setAuthState(prev => ({
          ...prev,
          isLoading: false,
          error: registerResult.message || 'Registration failed'
        }));
        return false;
      }
    } catch (error: any) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: error.message || 'Registration failed. Please try again.'
      }));
      return false;
    }
  };

  // Logout function
  const logout = () => {
    authService.logout();
    setAuthState(initialState);
  };

  // Update profile function
  const updateProfile = async (userData: Partial<User>) => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));

      if (!authState.user?.id) {
        throw new Error('User ID not found');
      }

      // Gọi API cập nhật thông tin
      const updatedUser = await userService.updateProfile(authState.user.id, userData);

      setAuthState(prev => ({
        ...prev,
        user: updatedUser as User,
        isLoading: false
      }));
    } catch (error: any) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: error.message || 'Profile update failed.'
      }));
    }
  };

  return (
    <AuthContext.Provider value={{ authState, login, register, logout, updateProfile }}>
      {children}
    </AuthContext.Provider>
  );
};