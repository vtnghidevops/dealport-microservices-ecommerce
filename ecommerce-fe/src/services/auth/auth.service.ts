import axios from 'axios';
import { UserLoginCredentials, UserRegistrationData, UserRole } from '@/types/user.model';

// Base URL for API requests - thay bằng URL thực tế của broker-service
const API_BASE_URL = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8080';

// Thêm log để kiểm tra API_BASE_URL
console.log("Auth Service API URL:", API_BASE_URL);

// Axios instance với cấu hình chung
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Thêm interceptor cho token authorization và refresh token khi cần
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Enhance the interceptor with a flag to prevent infinite loops
let isRefreshing = false;
let refreshSubscribers: Array<(token: string) => void> = [];

// Subscribe to token refresh
function subscribeTokenRefresh(cb: (token: string) => void) {
  refreshSubscribers.push(cb);
}

// Execute subscribers after token refresh
function onRefreshed(token: string) {
  refreshSubscribers.forEach(cb => cb(token));
  refreshSubscribers = [];
}

// Add response interceptor to handle token expiration
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Avoid infinite loops with the refresh token endpoint
    if (error.response?.config?.url?.includes('/auth/refresh')) {
      // Reset the refreshing state if the refresh token request itself fails
      isRefreshing = false;
      refreshSubscribers = [];
      // Clear tokens from localStorage
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('user');
      window.dispatchEvent(new Event('storage'));
      return Promise.reject(error);
    }

    // If error is 401 Unauthorized and we haven't tried to refresh the token yet
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      // Mark the request as retried
      originalRequest._retry = true;

      // If we're already refreshing the token, wait for the new token
      if (isRefreshing) {
        console.log("Another request is already refreshing the token, waiting...");
        try {
          // Wait for the new token
          const newToken = await new Promise<string>((resolve, reject) => {
            subscribeTokenRefresh((token: string) => {
              resolve(token);
            });

            // Add a timeout to prevent waiting forever
            setTimeout(() => {
              reject(new Error('Token refresh timeout'));
            }, 10000);
          });

          // Retry the original request with the new token
          originalRequest.headers['Authorization'] = `Bearer ${newToken}`;
          return axios(originalRequest);
        } catch (subscribeError) {
          console.error("Error waiting for token refresh:", subscribeError);
          return Promise.reject(error);
        }
      }

      // Start refreshing the token
      isRefreshing = true;

      try {
        console.log("Token expired. Attempting to refresh...");

        const refreshToken = localStorage.getItem('refreshToken');
        if (!refreshToken) {
          throw new Error('No refresh token available');
        }

        // Trực tiếp gọi API refresh thay vì tạo instance AuthService mới
        const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
          refresh_token: refreshToken
        });

        const responseData = response.data.data || response.data;

        if (responseData.access_token) {
          console.log("Token refreshed successfully, updating all pending requests");

          // Update token in localStorage
          localStorage.setItem('token', responseData.access_token);
          if (responseData.refresh_token) {
            localStorage.setItem('refreshToken', responseData.refresh_token);
          }

          // Update all pending requests with the new token
          originalRequest.headers['Authorization'] = `Bearer ${responseData.access_token}`;
          onRefreshed(responseData.access_token);

          // Reset refreshing state
          isRefreshing = false;

          // Retry the original request with the new token
          return axios(originalRequest);
        } else {
          throw new Error('Failed to refresh token');
        }
      } catch (refreshError) {
        console.error("Token refresh failed:", refreshError);

        // Reset refreshing state
        isRefreshing = false;
        refreshSubscribers = [];

        // Clear auth data from localStorage
        localStorage.removeItem('token');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('user');

        // Notify all components about auth state change
        window.dispatchEvent(new Event('storage'));

        // Only redirect if it's a navigation-capable environment (browser)
        if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
          console.log("Redirecting to login page");
          // Use a small delay to allow the current code to complete
          setTimeout(() => {
            window.location.href = '/login';
          }, 100);
        }

        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

/**
 * Service cho authentication-service
 * Tương tác với auth service thông qua broker
 */
class AuthService {
  /**
   * Đăng nhập với email và password
   * @param credentials thông tin đăng nhập
   * @returns thông tin user và token
   */
  async login(credentials: UserLoginCredentials) {
    try {
      console.log("Login attempt with:", credentials.email);
      const response = await apiClient.post('/auth/login', credentials);
      console.log("Login response:", response.data);

      // Kiểm tra cấu trúc response từ server
      const responseData = response.data.data || response.data;
      console.log("Login response data parsed:", responseData);
      console.log("User info from response:", responseData.user_info);

      if (responseData.access_token) {
        // Lưu token vào localStorage
        localStorage.setItem('token', responseData.access_token);
        localStorage.setItem('refreshToken', responseData.refresh_token);

        // Giải mã JWT token để lấy role
        const tokenData = this.parseJwt(responseData.access_token);
        console.log("Token data extracted:", tokenData);

        // Explicitly check and log the role from token
        const tokenRole: UserRole = tokenData?.role === 'admin' ? 'admin' : 'user';
        console.log("Role from JWT token:", tokenRole);

        // Nếu broker trả về thông tin user cơ bản, lưu tạm để hiển thị ngay
        if (responseData.user_info) {
          const basicUserInfo = {
            id: responseData.user_id,
            email: responseData.user_info.email,
            // IMPORTANT: Always prioritize role from JWT token
            role: tokenRole,
            username: responseData.user_info.username || responseData.user_info.email.split('@')[0],
            profile: {
              firstName: responseData.user_info.first_name,
              lastName: responseData.user_info.last_name
            }
          };
          console.log("Saving user info to localStorage with role from token:", basicUserInfo);
          localStorage.setItem('user', JSON.stringify(basicUserInfo));
        }
      }

      return {
        access_token: responseData.access_token,
        refresh_token: responseData.refresh_token,
        user_id: responseData.user_id,
        user_info: responseData.user_info || {},
        success: responseData.success || !!responseData.access_token,
        message: responseData.message || response.data.message || "Login successful"
      };
    } catch (error: any) {
      console.error('Login error:', error);
      if (error.response && error.response.data) {
        throw new Error(error.response.data.message || error.response.data.error || "Login failed");
      }
      throw error;
    }
  }

  /**
   * Đăng ký người dùng mới
   * @param userData thông tin đăng ký
   * @returns kết quả đăng ký
   */
  async register(userData: UserRegistrationData) {
    try {
      const registerData = {
        email: userData.email,
        password: userData.password,
        first_name: userData.firstName,
        last_name: userData.lastName,
        username: userData.username, // Use the provided username directly
        phone: userData.phone ? userData.phone : null
      };

      // Log registration data for debugging (without password)
      const debugData = {
        email: userData.email,
        first_name: userData.firstName,
        last_name: userData.lastName,
        username: userData.username
      };
      console.log("Sending registration data:", debugData);

      const response = await apiClient.post('/auth/register', registerData);
      console.log("Register response:", response.data);

      // Handle both data.data and direct data formats
      const responseData = response.data.data || response.data;

      return {
        success: responseData.success || !response.data.error || true,
        message: responseData.message || response.data.message || "Registration successful",
        user_id: responseData.user_id
      };
    } catch (error: any) {
      console.error('Registration error:', error);
      if (error.response && error.response.data) {
        throw new Error(error.response.data.message || error.response.data.error || "Registration failed");
      }
      throw error;
    }
  }

  /**
   * Đăng xuất - xóa token và thông tin user khỏi localStorage
   * @param logoutFromAllDevices Nếu true, sẽ đăng xuất khỏi tất cả các thiết bị
   */
  logout(logoutFromAllDevices = false) {
    console.log("Logging out user...", logoutFromAllDevices ? "from all devices" : "from current session");

    // Lấy thông tin user trước khi xóa localStorage
    const token = localStorage.getItem('token');
    const userId = token ? this.getUserIdFromToken(token) : null;
    const userEmail = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user') || '{}').email : '';

    console.log(`DEBUG logout: userId=${userId}, email=${userEmail}, logoutFromAllDevices=${logoutFromAllDevices}`);

    // Xóa thông tin trong localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');

    // Trigger a storage event to notify all components
    window.dispatchEvent(new Event('storage'));

    // Gửi request đến server để vô hiệu hóa token
    // Đảm bảo request này được gửi dù có user hay không
    if (userId && userEmail) {
      console.log(`DEBUG logout: Sending logout request to server for user ${userId}`);

      // Sử dụng axios trực tiếp thay vì apiClient để tránh interceptor
      // và đảm bảo request được gửi dù đã xóa token
      const logoutUrl = `${API_BASE_URL}/auth/logout`;
      const logoutData = {
        user_id: userId,
        email: userEmail,
        logout_all_devices: logoutFromAllDevices
      };

      //console.log("DEBUG logout: Request URL:", logoutUrl);
      //console.log("DEBUG logout: Request data:", logoutData);

      // Thêm timeout để đảm bảo request hoàn thành
      axios.post(logoutUrl, logoutData, {
        timeout: 5000 // 5 seconds timeout
      })
        .then(response => {
          console.log("DEBUG logout: Logout API response:", response.data);
        })
        .catch(err => {
          console.error("DEBUG logout: Logout API call failed:", err);
          // Continue with client-side logout regardless of server response
        });
    } else {
      console.warn("DEBUG logout: No user info available, skipping server logout request");
    }
  }

  /**
   * Kiểm tra tính hợp lệ của token hiện tại
   * @returns thông tin token hợp lệ
   */
  async validateToken() {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No token found');
      }

      // Đổi từ POST sang GET và truyền token qua header Authorization
      const response = await apiClient.get('/auth/validate');
      const responseData = response.data.data || response.data;

      return {
        valid: responseData.valid || !response.data.error,
        user_id: responseData.user_id,
        claims: responseData.claims || {}
      };
    } catch (error) {
      console.error('Token validation error:', error);
      throw error;
    }
  }

  /**
   * Làm mới token khi hết hạn
   * @returns token mới và các thông tin liên quan
   */
  async refreshToken() {
    try {
      const refreshToken = localStorage.getItem('refreshToken');
      if (!refreshToken) {
        throw new Error('No refresh token found');
      }

      // Kiểm tra ratelimit
      const lastRefreshTime = parseInt(localStorage.getItem('lastTokenRefresh') || '0');
      const currentTime = Date.now();
      const minRefreshInterval = 60 * 1000; // Ít nhất 1 phút giữa các lần refresh

      if (currentTime - lastRefreshTime < minRefreshInterval) {
        console.log("Skipping refresh due to rate limiting");
        throw new Error('Rate limit exceeded for token refresh');
      }

      // Lưu thời gian refresh
      localStorage.setItem('lastTokenRefresh', currentTime.toString());

      console.log("Attempting to refresh token with refresh token:", refreshToken.substring(0, 15) + "...");
      const response = await apiClient.post('/auth/refresh', { refresh_token: refreshToken });
      console.log("Refresh token response:", response.data);

      const responseData = response.data.data || response.data;

      if (!responseData.access_token) {
        throw new Error('Invalid response: No access token in refresh response');
      }

      console.log("Successfully refreshed token");

      // Always update localStorage with new tokens
      localStorage.setItem('token', responseData.access_token);
      if (responseData.refresh_token) {
        localStorage.setItem('refreshToken', responseData.refresh_token);
      }

      // Trigger a storage event so other components know tokens changed
      window.dispatchEvent(new Event('storage'));

      return {
        access_token: responseData.access_token,
        refresh_token: responseData.refresh_token || refreshToken, // Use old refresh if new one not provided
        success: responseData.success || !response.data.error,
        user_id: responseData.user_id || this.getUserIdFromToken(responseData.access_token)
      };
    } catch (error) {
      console.error('Token refresh error:', error);
      // Don't logout here - let the caller decide what to do
      throw error;
    }
  }

  /**
   * Get user ID from a JWT token
   * @param token The JWT token to parse
   * @returns The user ID from the token, or null if not found
   */
  getUserIdFromToken(token: string): string | null {
    try {
      const decoded = this.parseJwt(token);
      // Look for common user ID field names in JWT claims
      return decoded?.user_id || decoded?.sub || decoded?.userId || null;
    } catch (e) {
      console.error("Failed to extract user ID from token:", e);
      return null;
    }
  }

  /**
   * Kiểm tra xem người dùng đã đăng nhập chưa
   * @returns true nếu đã đăng nhập
   */
  isLoggedIn(): boolean {
    return !!localStorage.getItem('token');
  }

  /**
   * Lấy token hiện tại
   * @returns token đang lưu trong localStorage
   */
  getToken(): string | null {
    return localStorage.getItem('token');
  }

  /**
   * Hiển thị thông tin debug về trạng thái authentication hiện tại
   * Hữu ích khi gỡ lỗi
   */
  debugAuthState() {
    const token = localStorage.getItem('token');
    const refreshToken = localStorage.getItem('refreshToken');
    const userStr = localStorage.getItem('user');
    let user = null;

    try {
      if (userStr) {
        user = JSON.parse(userStr);
      }
    } catch (e) {
      console.error("Failed to parse user JSON from localStorage");
    }

    const tokenData = token ? this.parseJwt(token) : null;

    console.log("=== AUTH DEBUG ===");
    console.log("API Base URL:", API_BASE_URL);
    console.log("Has token:", !!token);
    console.log("Has refresh token:", !!refreshToken);
    console.log("Has user data:", !!user);
    console.log("Token data:", tokenData);
    console.log("User data:", user);

    if (token) {
      const expiryTime = tokenData?.exp ? new Date(tokenData.exp * 1000) : 'unknown';
      const isExpired = tokenData?.exp ? (tokenData.exp * 1000 < Date.now()) : 'unknown';
      console.log("Token expiration:", expiryTime);
      console.log("Token expired:", isExpired);
    }

    return {
      hasToken: !!token,
      hasRefreshToken: !!refreshToken,
      hasUserData: !!user,
      tokenData,
      userData: user,
      isExpired: tokenData?.exp ? (tokenData.exp * 1000 < Date.now()) : null
    };
  }

  /**
   * Parse JWT token to get payload data
   * @param token JWT token
   * @returns Decoded token payload
   */
  parseJwt(token: string) {
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
  }

  /**
   * Kiểm tra token có sắp hết hạn không và tự động làm mới nếu cần
   * @param expiryThresholdMinutes Số phút trước khi token hết hạn để làm mới (mặc định 10 phút)
   * @returns true nếu token còn hiệu lực, false nếu không có token hoặc token đã hết hạn
   */
  async checkAndRefreshTokenIfNeeded(expiryThresholdMinutes = 10): Promise<boolean> {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.log("No token found to check expiration");
        return false;
      }

      // Decode token to get expiration time
      const tokenData = this.parseJwt(token);
      if (!tokenData || !tokenData.exp) {
        console.log("Token doesn't have expiration claim");
        return false;
      }

      // Calculate time until expiration
      const expirationTime = tokenData.exp * 1000; // Convert to milliseconds
      const currentTime = Date.now();
      const timeUntilExpiry = expirationTime - currentTime;
      const expiryThresholdMs = expiryThresholdMinutes * 60 * 1000;

      // Chỉ log khi cần thiết để tránh spam console
      if (timeUntilExpiry <= expiryThresholdMs) {
        console.log(`Token expiry check: ${Math.floor(timeUntilExpiry / 60000)} minutes remaining`);
      }

      // Lưu thời gian kiểm tra cuối cùng để tránh refresh quá thường xuyên
      const lastRefreshCheck = parseInt(localStorage.getItem('lastTokenRefreshCheck') || '0');
      const minRefreshInterval = 5 * 60 * 1000; // Tối thiểu 5 phút giữa các lần refresh

      if (timeUntilExpiry <= 0) {
        // Token đã hết hạn
        if (currentTime - lastRefreshCheck > minRefreshInterval) {
          console.log("Token expired, attempting refresh");
          localStorage.setItem('lastTokenRefreshCheck', currentTime.toString());
          await this.refreshToken();
          return true;
        } else {
          console.log("Token expired but skipping refresh due to rate limiting");
          return false;
        }
      } else if (timeUntilExpiry <= expiryThresholdMs) {
        // Token sắp hết hạn
        if (currentTime - lastRefreshCheck > minRefreshInterval) {
          console.log(`Token will expire in ${Math.floor(timeUntilExpiry / 60000)} minutes, refreshing proactively`);
          localStorage.setItem('lastTokenRefreshCheck', currentTime.toString());
          await this.refreshToken();
          return true;
        } else {
          console.log("Token expiring soon but skipping refresh due to rate limiting");
          return true; // Vẫn trả về true vì token vẫn còn hiệu lực
        }
      }

      // Token còn hiệu lực và không gần hết hạn
      return true;
    } catch (error) {
      console.error("Error checking token expiration:", error);
      return false;
    }
  }

  /**
   * Set up automatic token refresh in the background
   * @param checkIntervalMinutes How often to check if token needs refresh (default 30 minutes)
   */
  setupAutoRefresh(checkIntervalMinutes = 30) {
    if (typeof window === 'undefined') return; // Only run in browser

    // Kiểm tra xem interval đã được cài đặt chưa, tránh đặt nhiều interval
    if ((window as any).__authRefreshInterval) {
      console.log("Auto refresh already set up, skipping");
      return;
    }

    // Check immediately on setup if user is logged in, but only if no recent check
    if (this.isLoggedIn()) {
      const lastRefreshCheck = parseInt(localStorage.getItem('lastTokenRefreshCheck') || '0');
      const currentTime = Date.now();
      // Chỉ kiểm tra nếu lần cuối kiểm tra cách đây hơn 5 phút
      if (currentTime - lastRefreshCheck > 5 * 60 * 1000) {
        this.checkAndRefreshTokenIfNeeded().catch(err =>
          console.log("Initial token check failed:", err)
        );
      }
    }

    // Kiểm tra định kỳ
    const intervalId = setInterval(() => {
      // Chỉ kiểm tra nếu người dùng đã đăng nhập
      if (this.isLoggedIn()) {
        this.checkAndRefreshTokenIfNeeded().catch(err =>
          console.log("Periodic token check failed:", err)
        );
      }
    }, checkIntervalMinutes * 60 * 1000); // Kiểm tra mỗi 30 phút

    // Lưu ID interval để có thể hủy khi cần
    (window as any).__authRefreshInterval = intervalId;
    console.log(`Auto refresh set up with interval of ${checkIntervalMinutes} minutes`);

    // Thêm event listener cho việc thay đổi visibility tab
    document.addEventListener('visibilitychange', () => {
      // Chỉ kiểm tra khi tab hiển thị và người dùng đã đăng nhập
      // Thêm kiểm tra tỷ lệ giới hạn để không refresh quá thường xuyên
      if (document.visibilityState === 'visible' && this.isLoggedIn()) {
        const lastRefreshCheck = parseInt(localStorage.getItem('lastTokenRefreshCheck') || '0');
        const currentTime = Date.now();
        if (currentTime - lastRefreshCheck > 5 * 60 * 1000) {
          this.checkAndRefreshTokenIfNeeded().catch(err =>
            console.log("Visibility change token check failed:", err)
          );
        }
      }
    });

    return intervalId;
  }

  /**
   * Stop the automatic token refresh
   */
  stopAutoRefresh() {
    if (typeof window === 'undefined') return;

    if ((window as any).__authRefreshInterval) {
      clearInterval((window as any).__authRefreshInterval);
      (window as any).__authRefreshInterval = null;
    }
  }
}

// Create and export instance
const authService = new AuthService();

// Export service instance and direct API functions
export default authService;

// Legacy functions to support direct import
export const login = async (data: UserLoginCredentials) => {
  return authService.login(data);
};

export const register = async (data: UserRegistrationData) => {
  return authService.register(data);
};

// Define a UserProfileData interface for the updateProfile function
interface UserProfileData {
  firstName?: string;
  lastName?: string;
  email?: string;
  [key: string]: any;
}

export const getProfile = async () => {
  try {
    const response = await apiClient.get("/user/profile");
    return response.data;
  } catch (error: any) {
    console.error('Error getting profile:', error);
    throw error;
  }
};

export const updateProfile = async (data: UserProfileData) => {
  try {
    const response = await apiClient.put("/user/profile", data);
    return response.data;
  } catch (error: any) {
    console.error('Error updating profile:', error);
    throw error;
  }
};

export const logout = async (logoutFromAllDevices = false) => {
  return authService.logout(logoutFromAllDevices);
};

/**
 * Đăng xuất khỏi tất cả các thiết bị
 * Tiện ích để gọi logout với tham số logoutFromAllDevices = true
 */
export const logoutFromAllDevices = async () => {
  return authService.logout(true);
};

export const forgotPassword = async (email: string) => {
  try {
    // Kiểm tra tài khoản tồn tại trước khi yêu cầu đặt lại mật khẩu
    await checkAccountExists(email);

    // Yêu cầu password reset thông qua API password-reset của broker
    const response = await apiClient.post("auth/password-reset", { email });
    console.log("Password reset response:", response.data);
    return response.data;
  } catch (error: any) {
    console.error('Error initiating password reset:', error);
    throw error;
  }
};

export const resetPassword = async (email: string, password: string, token: string) => {
  try {
    // Normalize email to lowercase to ensure consistency with the backend
    const normalizedEmail = email.toLowerCase().trim();

    // Kiểm tra tài khoản tồn tại trước khi cập nhật mật khẩu
    await checkAccountExists(normalizedEmail);

    // Log debug info about token without revealing sensitive details
    console.log(`Token provided (length: ${token.length})`);

    // Prepare the request payload
    const payload = {
      token,  // Send the token as-is, without any modifications
      password,
      email: normalizedEmail // Use normalized email
    };

    console.log("Sending password update request with payload:",
      { ...payload, password: "[MASKED]", token: `${token.substring(0, 10)}...` });

    // Sử dụng API password-update để cập nhật mật khẩu
    const response = await apiClient.post("/auth/password-update", payload);
    console.log("Password update response:", response.data);

    // Ensure we return a consistent response format
    return {
      success: response.data.success || response.data.data?.success || !response.data.error,
      message: response.data.message || response.data.data?.message || "Password updated successfully",
      data: response.data.data || response.data
    };
  } catch (error: any) {
    console.error('Error resetting password:', error);
    throw new Error(error.response?.data?.message || error.message || "Failed to update password");
  }
};

// OTP related functions
export const requestOTP = async (email: string, purpose: string) => {
  try {
    // Kiểm tra tài khoản tồn tại trước khi gửi OTP
    await checkAccountExists(email);

    // Sử dụng API request-otp để yêu cầu OTP mới
    const response = await apiClient.post("/auth/request-otp", {
      email,
      purpose
    });
    console.log("Request OTP response:", response.data);
    return response.data;
  } catch (error: any) {
    console.error('Error requesting OTP:', error);
    throw error;
  }
};

export const verifyOTP = async (email: string, otp: string, purpose: string) => {
  try {
    // Kiểm tra tài khoản tồn tại trước khi xác thực OTP
    await checkAccountExists(email);

    if (purpose === 'password_reset') {
      // Sử dụng API verify-password-reset để xác minh OTP cho đặt lại mật khẩu
      const response = await apiClient.post("/auth/verify-password-reset", {
        email,
        otp,
        purpose
      });
      console.log("Verify password reset OTP response:", response.data);
      return response.data;
    } else {
      // Sử dụng API verify-registration để xác minh OTP cho đăng ký
      const response = await apiClient.post("/auth/verify-registration", {
        email,
        otp,
        purpose
      });
      console.log("Verify registration OTP response:", response.data);
      return response.data;
    }
  } catch (error: any) {
    console.error('Error verifying OTP:', error);
    throw error;
  }
};

export const resendOTP = async (email: string, purpose: string) => {
  try {
    // Kiểm tra tài khoản tồn tại trước khi gửi lại OTP
    await checkAccountExists(email);

    // Gửi lại OTP bằng cách dùng lại endpoint request-otp
    const response = await apiClient.post("/auth/request-otp", {
      email,
      purpose
    });
    console.log("Resend OTP response:", response.data);
    return response.data;
  } catch (error: any) {
    console.error('Error resending OTP:', error);
    throw error;
  }
};

// Thêm phương thức kiểm tra tài khoản tồn tại
export const checkAccountExists = async (email: string): Promise<boolean> => {
  try {
    // Gọi API kiểm tra tài khoản
    const response = await apiClient.post("/auth/check-account", { email });
    console.log("Check account response:", response.data);

    const exists = response.data.exists || response.data.data?.exists || false;

    if (!exists) {
      throw new Error("Account not found. Please check your email or create a new account.");
    }

    return exists;
  } catch (error: any) {
    // Nếu lỗi do API (như 404 Not Found), có nghĩa là tài khoản không tồn tại
    if (error.response && error.response.status === 404) {
      throw new Error("Account not found. Please check your email or create a new account.");
    }

    // Nếu lỗi từ phản hồi API
    if (error.response && error.response.data) {
      throw new Error(error.response.data.message || error.response.data.error || "Failed to verify account");
    }

    // Nếu là lỗi tự tạo, ném lại
    if (error instanceof Error) {
      throw error;
    }

    // Lỗi không xác định
    throw new Error("Failed to check if account exists");
  }
};

export const changePassword = async (currentPassword: string, newPassword: string) => {
  try {
    // Get token from localStorage
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error("You are not logged in. Please log in to change your password.");
    }

    // Sử dụng endpoint password-update với current_password để cập nhật mật khẩu
    const response = await apiClient.post("/auth/password-update", {
      current_password: currentPassword,
      password: newPassword,
      token: token  // Send the token explicitly in the request body
    });

    console.log("Change password response:", response.data);
    return response.data;
  } catch (error: any) {
    console.error('Error changing password:', error);

    // Handle specific error messages from backend
    if (error.response && error.response.data) {
      throw new Error(error.response.data.message || error.response.data.error || "Failed to change password");
    }

    throw error;
  }
}; 