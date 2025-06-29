import axios from 'axios';
import { User } from '@/types/user.model';
import api from '@/services/user';
import { BASE_API_URL } from '@/utils/api-config';

// Axios instance với cấu hình  
const apiClient = axios.create({
  baseURL: BASE_API_URL,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Thêm interceptor cho token authorization
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

/**
 * Service cho user-service
 * Tương tác với user service thông qua broker
 */
class UserService {
  /**
   * Lấy thông tin user hiện tại từ token
   * @returns Chi tiết thông tin user
   */
  async getCurrentUser() {
    try {
      const userId = this.getUserIdFromStorage();
      if (!userId) {
        throw new Error('User ID not found');
      }

      return this.getUserById();
    } catch (error) {
      console.error('Get current user error:', error);
      throw error;
    }
  }

  /**
   * Lấy thông tin user theo ID
   * @returns Chi tiết thông tin user
   */
  async getUserById(): Promise<User> {
    try {
      const response = await api.get(`/users/me`);

      if (response.status === 200 && response.data.error === false) {
        return response.data.data;
      } else {
        throw new Error(response.data.message || 'Failed to fetch user data');
      }
    } catch (error: any) {
      console.error('Error fetching user data:', error);
      throw error;
    }
  }

  /**
   * Cập nhật thông tin user
   * @param userData Thông tin user cần cập nhật
   * @returns Chi tiết thông tin user sau khi cập nhật
   */
  async updateProfile(userData: Partial<User>): Promise<User> {
    try {
      const response = await api.put(`/users/profile`, userData);

      if (response.status === 200 && response.data.error === false) {
        return response.data.data;
      } else {
        throw new Error(response.data.message || 'Failed to update user profile');
      }
    } catch (error: any) {
      console.error('Error updating user profile:', error);
      throw error;
    }
  }

  /**
   * Cập nhật mật khẩu user
   * @param currentPassword Mật khẩu hiện tại
   * @param newPassword Mật khẩu mới
   * @returns Kết quả cập nhật mật khẩu
   */
  async updatePassword(currentPassword: string, newPassword: string): Promise<any> {
    try {
      const response = await api.put(`/users/password`, {
        current_password: currentPassword,
        new_password: newPassword
      });

      if (response.status === 200 && response.data.error === false) {
        return response.data;
      } else {
        throw new Error(response.data.message || 'Failed to update password');
      }
    } catch (error: any) {
      console.error('Error updating password:', error);
      throw error;
    }
  }

  /**
   * Lấy danh sách địa chỉ của user
   * @returns Danh sách địa chỉ
   */
  async getUserAddresses(): Promise<any[]> {
    try {
      const response = await api.get(`/users/addresses`);

      if (response.status === 200 && response.data.error === false) {
        return response.data.data || [];
      } else {
        throw new Error(response.data.message || 'Failed to fetch user addresses');
      }
    } catch (error: any) {
      console.error('Error fetching user addresses:', error);
      return [];
    }
  }

  /**
   * Thêm địa chỉ mới cho người dùng
   * @param address Thông tin địa chỉ
   * @returns Thông tin user đã cập nhật
   */
  async addAddress(address: Omit<any, 'id'>) {
    try {
      // Format địa chỉ theo yêu cầu của API
      const addressData = {
        line1: address.street,
        line2: address.additionalInfo || '',
        city: address.city,
        state: address.state,
        postal_code: address.zipCode,
        country: address.country,
        is_default: address.isDefault || false,
        address_type: address.type || 'shipping'
      };

      const response = await api.post(`/users/me/addresses`, addressData);
      return response.data;
    } catch (error) {
      console.error('Add address error:', error);
      throw error;
    }
  }

  /**
   * Xóa địa chỉ của người dùng
   * @param addressId ID của địa chỉ
   * @returns Kết quả xóa
   */
  async deleteAddress(addressId: string): Promise<boolean> {
    try {
      const response = await api.delete(`/users/me/addresses/${addressId}`);

      return response.status === 200 && response.data.error === false;
    } catch (error: any) {
      console.error('Error deleting address:', error);
      throw error;
    }
  }

  /**
   * Lấy danh sách sản phẩm yêu thích
   * @returns Danh sách sản phẩm yêu thích
   */
  async getWishlist() {
    try {
      const response = await api.get(`/users/me/wishlist`);
      return response.data;
    } catch (error) {
      console.error('Get wishlist error:', error);
      throw error;
    }
  }

  /**
   * Lấy ID của user từ localStorage
   * @returns User ID hoặc null
   */
  getUserIdFromStorage(): string | null {
    const userStr = localStorage.getItem('user');
    if (!userStr) return null;

    try {
      const user = JSON.parse(userStr);
      return user.id || null;
    } catch (error) {
      console.error('Error parsing user from localStorage:', error);
      return null;
    }
  }

  /**
   * Lấy thông tin user cơ bản từ localStorage
   * @returns Thông tin user hoặc null
   */
  getUserFromStorage(): Partial<User> | null {
    const userStr = localStorage.getItem('user');
    if (!userStr) return null;

    try {
      return JSON.parse(userStr);
    } catch (error) {
      console.error('Error parsing user from localStorage:', error);
      return null;
    }
  }

  /**
   * Cập nhật thông tin địa chỉ
   * @param addressId ID của địa chỉ
   * @param addressData Dữ liệu cập nhật
   * @returns Địa chỉ đã cập nhật
   */
  async updateAddress(addressId: string, addressData: any): Promise<any> {
    try {
      const response = await api.put(`/users/me/addresses/${addressId}`, addressData);

      if (response.status === 200 && response.data.error === false) {
        return response.data.data;
      } else {
        throw new Error(response.data.message || 'Failed to update address');
      }
    } catch (error: any) {
      console.error('Error updating address:', error);
      throw error;
    }
  }

  /**
   * Tạo địa chỉ mới
   * @param addressData Thông tin địa chỉ
   * @returns Địa chỉ đã tạo
   */
  async createAddress(addressData: any): Promise<any> {
    try {
      const response = await api.post(`/users/me/addresses`, addressData);

      if (response.status === 201 && response.data.error === false) {
        return response.data.data;
      } else {
        throw new Error(response.data.message || 'Failed to create address');
      }
    } catch (error: any) {
      console.error('Error creating address:', error);
      throw error;
    }
  }
}

export default new UserService(); 