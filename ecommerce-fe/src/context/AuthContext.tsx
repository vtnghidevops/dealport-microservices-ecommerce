// context/AuthContext.tsx
import React, { createContext, useState, ReactNode, useEffect } from 'react';
import { User, AuthState, UserRegistrationData, UserRole } from '@/types/user.model';
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

// Define AuthContextType interface
type AuthContextType = {
  authState: AuthState;
  login: (email: string, password: string) => Promise<{ success: boolean, error?: string }>;
  register: (userData: UserRegistrationData) => Promise<boolean>;
  logout: (logoutFromAllDevices?: boolean) => void;
  logoutFromAllDevices: () => void;
  updateProfile: (userData: Partial<User>) => Promise<void>;
};

// Create the context with a default empty value
export const AuthContext = createContext<AuthContextType>({
  authState: initialState,
  login: async () => ({ success: false }),
  register: async () => false,
  logout: () => { },
  logoutFromAllDevices: () => { },
  updateProfile: async () => { },
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

  // Thêm helper function to parse JWT token
  const parseJwt = (token: string) => {
    try {
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));

      return JSON.parse(jsonPayload);
    } catch (e) {
      console.error("Failed to parse JWT token:", e);
      return null;
    }
  };

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

          // Extract role directly from JWT token
          let tokenRole: UserRole = 'user';
          try {
            const tokenData = parseJwt(authState.accessToken);
            if (tokenData && tokenData.role) {
              // Make sure tokenRole is a valid UserRole value
              tokenRole = tokenData.role === 'admin' ? 'admin' : 'user';
              console.log("Role extracted from JWT token during validation:", tokenRole);
            }
          } catch (jwtError) {
            console.error("Error parsing JWT during validation:", jwtError);
          }

          // Kiểm tra token
          const validateResult = await authService.validateToken();

          // Kiểm tra các claims để debug
          if (validateResult.claims) {
            console.log("Claims from validation:", validateResult.claims);
            // Nếu có role trong claims, lưu tạm vào localStorage
            if (validateResult.claims.role) {
              const userStr = localStorage.getItem('user');
              if (userStr) {
                try {
                  const userData = JSON.parse(userStr);
                  // Use role from JWT token
                  userData.role = tokenRole;
                  localStorage.setItem('user', JSON.stringify(userData));
                  console.log("Updated user role in localStorage to:", tokenRole);
                } catch (e) {
                  console.error("Failed to update role from claims", e);
                }
              }
            }
          }

          if (validateResult.valid) {
            // Nếu token hợp lệ, lấy thông tin user đầy đủ
            try {
              // const userId = validateResult.user_id;
              const userData = await userService.getUserById() as User;

              // Use role from JWT token
              const updatedUserData = {
                ...userData,
                role: tokenRole
              };

              console.log("Setting auth state with role from token:", tokenRole);

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
              // console.log("Token refresh successful:", refreshResult);

              // Get user again with the new token
              try {
                // Use a safe fallback for userId since refreshResult.user_id might not exist
                const userId = authState.user?.id;
                if (!userId) {
                  throw new Error('Cannot get user data after refresh: missing user ID');
                }

                // console.log("Fetching user data after token refresh for ID:", userId);
                const userData = await userService.getUserById() as User;

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
          // Extract role from JWT token directly
          const token = loginResult.access_token;
          let tokenRole: UserRole = 'user'; // Default role

          // Parse JWT to extract role claim
          try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
              return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            const tokenData = JSON.parse(jsonPayload);
            // Ensure it's a valid UserRole
            tokenRole = tokenData.role === 'admin' ? 'admin' : 'user';
            console.log("Role extracted from JWT token:", tokenRole);
          } catch (jwtError) {
            console.error("Error parsing JWT token:", jwtError);
          }

          const userData = await userService.getUserById();
          const userWithProfile = {
            ...userData,
            // Prioritize JWT token role over other sources
            role: tokenRole,
            profile: userData.profile || {
              firstName: loginResult.user_info.first_name,
              lastName: loginResult.user_info.last_name
            }
          };

          console.log("Setting user role from token:", tokenRole);

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
          console.error("Error getting user details:", userError);

          // Extract role from JWT token directly
          const token = loginResult.access_token;
          let tokenRole: UserRole = 'user'; // Default role

          // Parse JWT to extract role claim
          try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
              return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            const tokenData = JSON.parse(jsonPayload);
            // Ensure it's a valid UserRole
            tokenRole = tokenData.role === 'admin' ? 'admin' : 'user';
            console.log("Fallback: Role extracted from JWT token:", tokenRole);
          } catch (jwtError) {
            console.error("Error parsing JWT token:", jwtError);
          }

          const basicUserInfo = {
            id: loginResult.user_id,
            email: loginResult.user_info.email,
            // Prioritize JWT token role
            role: tokenRole,
            username: loginResult.user_info.username,
            profile: {
              firstName: loginResult.user_info.first_name,
              lastName: loginResult.user_info.last_name
            }
          };

          console.log("Setting basic user info with role from token:", tokenRole);

          setAuthState({
            user: basicUserInfo as unknown as User,
            accessToken: loginResult.access_token,
            isAuthenticated: true,
            isLoading: false,
            error: null
          });
          localStorage.setItem('user', JSON.stringify(basicUserInfo));
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
  const logout = (logoutFromAllDevices = false) => {
    console.log("AuthContext: Calling logout with logoutFromAllDevices =", logoutFromAllDevices);
    authService.logout(logoutFromAllDevices);

    // Update auth state with initialState
    setAuthState({
      user: null,
      accessToken: null,
      isAuthenticated: false,
      isLoading: false,
      error: null
    });

    console.log("AuthContext: Auth state updated after logout");
  };

  // Specialized function to logout from all devices
  const logoutFromAllDevices = () => {
    console.log("AuthContext: Calling logoutFromAllDevices");
    logout(true);
  };

  // Update profile function
  const updateProfile = async (userData: Partial<User>) => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));

      if (!authState.user?.id) {
        throw new Error('User ID not found');
      }

      // Gọi API cập nhật thông tin
      const updatedUser = await userService.updateProfile(userData);

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
    <AuthContext.Provider value={{ authState, login, register, logout, logoutFromAllDevices, updateProfile }}>
      {children}
    </AuthContext.Provider>
  );
};