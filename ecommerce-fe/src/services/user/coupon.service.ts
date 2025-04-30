import axios from 'axios';

const API_URL = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8080';

export interface Coupon {
  id: string;
  code: string;
  discount: number;
  discountType: 'percentage' | 'fixed';
  minOrderAmount: number;
  maxUsage: number;
  usageCount: number;
  validFrom: string;
  validTo: string;
  isActive: boolean;
  description?: string;
}

export interface CouponResponseData {
  success: boolean;
  message: string;
  data?: Coupon | Coupon[];
}

// Get authentication header from localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  if (!token) return {};
  return { Authorization: `Bearer ${token}` };
};

class CouponService {
  private baseUrl = `${API_URL}/api/v1/coupons`;

  /**
   * Get all coupons with optional pagination
   */
  async getCoupons(page: number = 1, limit: number = 10): Promise<{ coupons: Coupon[], total: number }> {
    try {
      const response = await axios.get(`${this.baseUrl}`, {
        headers: getAuthHeader(),
        params: { page, limit }
      });

      return {
        coupons: response.data.data as Coupon[],
        total: response.data.meta?.total || response.data.data.length
      };
    } catch (error) {
      console.error('Error fetching coupons:', error);
      throw this.handleError(error);
    }
  }

  /**
   * Get a single coupon by ID
   */
  async getCoupon(id: string): Promise<Coupon> {
    try {
      const response = await axios.get(`${this.baseUrl}/${id}`, {
        headers: getAuthHeader()
      });

      return response.data.data as Coupon;
    } catch (error) {
      console.error(`Error fetching coupon ${id}:`, error);
      throw this.handleError(error);
    }
  }

  /**
   * Create a new coupon
   */
  async createCoupon(coupon: Omit<Coupon, 'id' | 'usageCount'>): Promise<Coupon> {
    try {
      const response = await axios.post(`${this.baseUrl}`, coupon, {
        headers: getAuthHeader()
      });

      return response.data.data as Coupon;
    } catch (error) {
      console.error('Error creating coupon:', error);
      throw this.handleError(error);
    }
  }

  /**
   * Update an existing coupon
   */
  async updateCoupon(id: string, coupon: Partial<Coupon>): Promise<Coupon> {
    try {
      const response = await axios.put(`${this.baseUrl}/${id}`, coupon, {
        headers: getAuthHeader()
      });

      return response.data.data as Coupon;
    } catch (error) {
      console.error(`Error updating coupon ${id}:`, error);
      throw this.handleError(error);
    }
  }

  /**
   * Delete a coupon
   */
  async deleteCoupon(id: string): Promise<{ success: boolean, message: string }> {
    try {
      const response = await axios.delete(`${this.baseUrl}/${id}`, {
        headers: getAuthHeader()
      });

      return {
        success: true,
        message: response.data.message
      };
    } catch (error) {
      console.error(`Error deleting coupon ${id}:`, error);
      throw this.handleError(error);
    }
  }

  /**
   * Apply a coupon to the cart
   */
  async applyCoupon(couponCode: string): Promise<any> {
    try {
      const response = await axios.post(`${API_URL}/api/v1/cart/coupon`,
        { couponCode },
        { headers: getAuthHeader() }
      );

      return response.data;
    } catch (error) {
      console.error(`Error applying coupon ${couponCode}:`, error);
      throw this.handleError(error);
    }
  }

  /**
   * Remove coupon from cart
   */
  async removeCoupon(): Promise<any> {
    try {
      const response = await axios.delete(`${API_URL}/api/v1/cart/coupon`, {
        headers: getAuthHeader()
      });

      return response.data;
    } catch (error) {
      console.error('Error removing coupon:', error);
      throw this.handleError(error);
    }
  }

  /**
   * Handle API errors
   */
  private handleError(error: any): Error {
    if (axios.isAxiosError(error)) {
      // Handle specific error cases
      if (error.response?.status === 401) {
        return new Error('Unauthorized: Please login to manage coupons');
      } else if (error.response?.status === 403) {
        return new Error('Forbidden: You do not have permission to access this resource');
      } else if (error.response?.status === 404) {
        return new Error('Coupon not found');
      } else if (error.response?.data?.message) {
        return new Error(error.response.data.message);
      }
    }

    return new Error('An error occurred while processing your request');
  }
}

// Export a singleton instance
const couponService = new CouponService();
export default couponService; 