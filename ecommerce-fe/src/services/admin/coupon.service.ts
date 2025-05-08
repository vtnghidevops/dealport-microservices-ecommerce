import axios from 'axios';

const VITE_PUBLIC_BROKER_API_URL = import.meta.env.VITE_PUBLIC_BROKER_API_URL || 'http://localhost:8080';

// Frontend interface using camelCase (dữ liệu đến từ API đã ở dạng camelCase)
export interface Coupon {
  id: string;
  code: string;
  discount: number;
  discountType: 'percentage' | 'fixed';
  minOrderAmount: number;
  maxDiscount?: number;
  maxUsage: number;
  usageCount: number;
  validFrom: string;
  validTo: string;
  isActive: boolean;
  description?: string;
  createdAt?: string;
  updatedAt?: string;
}

// API response interface
export interface CouponResponseData {
  data?: Coupon | Coupon[];
  meta?: {
    total: number;
    page: number;
    limit: number;
  };
  error?: boolean;
  message?: string;
}

// Format date to RFC3339
const formatDateToRFC3339 = (dateString: string) => {
  // If already in RFC3339 format, return as is
  if (dateString.includes('T')) {
    return dateString;
  }
  // Otherwise, add time component
  return `${dateString}T00:00:00Z`;
};

// Prepare coupon data for API (maintain camelCase but format dates)
const prepareCouponForAPI = (coupon: Partial<Coupon>): Partial<Coupon> => {
  const result = { ...coupon };

  // Format dates to RFC3339 if present
  if (coupon.validFrom) {
    result.validFrom = formatDateToRFC3339(coupon.validFrom);
  }

  if (coupon.validTo) {
    result.validTo = formatDateToRFC3339(coupon.validTo);
  }

  return result;
};

// Get authentication header from localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  console.log('Current auth token:', token ? `${token.substring(0, 15)}...` : 'No token found');
  if (!token) return {};
  return { Authorization: `Bearer ${token}` };
};

class CouponService {
  private baseUrl = `${VITE_PUBLIC_BROKER_API_URL}/coupons`;

  /**
   * Get all coupons with optional pagination
   */
  async getCoupons(page: number = 1, limit: number = 10): Promise<{ coupons: Coupon[], total: number }> {
    try {
      const response = await axios.get<CouponResponseData>(`${this.baseUrl}`, {
        headers: getAuthHeader(),
        params: { page, limit }
      });

      console.log('API Response:', response.data);

      // Get coupons array
      let coupons: Coupon[] = [];
      if (response.data.data && Array.isArray(response.data.data)) {
        coupons = response.data.data;
      }

      return {
        coupons,
        total: response.data.meta?.total || coupons.length
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
      const response = await axios.get<CouponResponseData>(`${this.baseUrl}/${id}`, {
        headers: getAuthHeader()
      });

      if (!response.data.data) {
        throw new Error('No coupon data returned from API');
      }

      return response.data.data as Coupon;
    } catch (error) {
      console.error(`Error fetching coupon ${id}:`, error);
      throw this.handleError(error);
    }
  }

  /**
   * Get a coupon by code
   */
  async getCouponByCode(code: string): Promise<Coupon> {
    try {
      const response = await axios.get<CouponResponseData>(`${this.baseUrl}/code/${code}`, {
        headers: getAuthHeader()
      });

      if (!response.data.data) {
        throw new Error('No coupon data returned from API');
      }

      return response.data.data as Coupon;
    } catch (error) {
      console.error(`Error fetching coupon by code ${code}:`, error);
      throw this.handleError(error);
    }
  }

  /**
   * Create a new coupon
   */
  async createCoupon(coupon: Omit<Coupon, 'id' | 'usageCount'>): Promise<Coupon> {
    try {
      console.log('Creating coupon with data:', coupon);
      console.log('API URL:', `${this.baseUrl}`);

      // Prepare coupon data for API (keep camelCase, format dates)
      const payload = prepareCouponForAPI(coupon);
      console.log('Sending payload:', payload);

      const response = await axios.post<CouponResponseData>(`${this.baseUrl}`, payload, {
        headers: {
          ...getAuthHeader(),
          'Content-Type': 'application/json'
        }
      });

      console.log('API response:', response.data);

      if (!response.data.data) {
        throw new Error('No coupon data returned from API');
      }

      // Handle both array and single object response
      if (Array.isArray(response.data.data)) {
        return response.data.data[0];
      }

      return response.data.data as Coupon;
    } catch (error) {
      console.error('Error creating coupon:', error);
      if (axios.isAxiosError(error)) {
        console.error('API Error details:', {
          status: error.response?.status,
          statusText: error.response?.statusText,
          data: error.response?.data,
          headers: error.response?.headers
        });
      }
      throw this.handleError(error);
    }
  }

  /**
   * Update an existing coupon
   */
  async updateCoupon(id: string, coupon: Partial<Coupon>): Promise<Coupon> {
    try {
      // Prepare coupon data for API (keep camelCase, format dates)
      const payload = prepareCouponForAPI(coupon);

      const response = await axios.put<CouponResponseData>(`${this.baseUrl}/${id}`, payload, {
        headers: {
          ...getAuthHeader(),
          'Content-Type': 'application/json'
        }
      });

      if (!response.data.data) {
        throw new Error('No coupon data returned from API');
      }

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
        success: !response.data.error,
        message: response.data.message || 'Coupon deleted successfully'
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
      const response = await axios.post(`${VITE_PUBLIC_BROKER_API_URL}/cart/coupon`,
        { coupon_code: couponCode },
        {
          headers: {
            ...getAuthHeader(),
            'Content-Type': 'application/json'
          }
        }
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
      const response = await axios.delete(`${VITE_PUBLIC_BROKER_API_URL}/cart/coupon`, {
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
      } else if (error.response?.data?.error) {
        return new Error(error.response.data.error);
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