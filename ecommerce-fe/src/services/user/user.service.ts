import axios from 'axios';
import { User } from '@/types/user.model';
import api from '@/services/user';

// Base URL cho API requests đã được định nghĩa trong api/index.ts
const API_BASE_URL = import.meta.env.VITE_PUBLIC_BROKER_API_URL || 'http://localhost:8080/api/v1';

// Axios instance với cấu hình  
const apiClient = axios.create({
  baseURL: API_BASE_URL,
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

      return this.getUserById(userId);
    } catch (error) {
      console.error('Get current user error:', error);
      throw error;
    }
  }

  /**
   * Lấy thông tin user theo ID
   * @param id ID của user
   * @returns Chi tiết thông tin user
   */
  async getUserById(userId: string): Promise<User> {
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
   * Cập nhật thông tin profile của user
   * @param userId ID của user
   * @param userData Dữ liệu cần cập nhật
   * @returns Thông tin user đã cập nhật
   */
  async updateProfile(userId: string, userData: Partial<User>): Promise<User> {
    try {
      // Deep clone user data to avoid modifying the original
      const userDataToUpdate = JSON.parse(JSON.stringify(userData));

      // Format profile data according to API expectations
      if (userDataToUpdate.profile) {
        // Convert null values to empty strings for the API
        Object.keys(userDataToUpdate.profile).forEach(key => {
          if (userDataToUpdate.profile[key] === null) {
            userDataToUpdate.profile[key] = '';
          }
        });
      }

      console.log('Updating profile with data:', userDataToUpdate);

      const response = await api.put(`/users/me`, userDataToUpdate);

      if (response.status === 200 && response.data.error === false) {
        // Update local storage with the updated user data
        const currentUser = this.getUserFromStorage();
        if (currentUser) {
          const updatedUser = { ...currentUser, ...response.data.data };
          localStorage.setItem('user', JSON.stringify(updatedUser));
        }

        return response.data.data;
      } else {
        throw new Error(response.data.message || 'Failed to update profile');
      }
    } catch (error: any) {
      console.error('Error updating profile:', error);
      throw error;
    }
  }

  /**
   * Thêm địa chỉ mới cho người dùng
   * @param userId ID của user
   * @param address Thông tin địa chỉ
   * @returns Thông tin user đã cập nhật
   */
  async addAddress(userId: string, address: Omit<any, 'id'>) {
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
   * @param userId ID của user
   * @returns Danh sách sản phẩm yêu thích
   */
  async getWishlist(userId: string) {
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