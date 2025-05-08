import axios from 'axios';
import { CartItem, CartTotalsData } from '@/types/cart.model';

const VITE_PUBLIC_BROKER_API_URL = import.meta.env.VITE_VITE_PUBLIC_BROKER_API_URL || 'http://localhost:8080';

export interface CartItemRequest {
  productId: number;
  name: string;
  price: number;
  originalPrice?: number;
  quantity: number;
  imageUrl: string;
}

export interface CartResponse {
  id: string;
  userId: string;
  items: CartItem[];
  totals: CartTotalsData;
  couponCode: string;
  discountAmount: number;
  itemCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface StatusResponse {
  success: boolean;
  message: string;
}

// Lấy token xác thực từ localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('You must be logged in to manage your cart');
  }
  return { Authorization: `Bearer ${token}` };
};

// Lấy user ID từ localStorage
const getCurrentUserId = (): string => {
  const userId = localStorage.getItem('user_id');
  if (!userId) {
    throw new Error('User ID not found. Please log in again.');
  }
  return userId;
};

// Singleton instance of cart service
class CartService {
  private baseUrl = `${VITE_PUBLIC_BROKER_API_URL}/api/v1/cart`;

  // Get the user's cart
  async getCart(): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.get(this.baseUrl, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error getting cart', error);
      throw error;
    }
  }

  // Add an item to the cart
  async addCartItem(item: CartItemRequest): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(`${this.baseUrl}/items`, {
        productId: item.productId,
        name: item.name,
        price: item.price,
        originalPrice: item.originalPrice,
        quantity: item.quantity,
        imageUrl: item.imageUrl
      }, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error adding item to cart', error);
      throw error;
    }
  }

  // Update an item in the cart
  async updateCartItem(itemId: string, quantity: number): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.put(`${this.baseUrl}/items/${itemId}`, { quantity }, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error updating cart item', error);
      throw error;
    }
  }

  // Remove an item from the cart
  async removeCartItem(itemId: string): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.delete(`${this.baseUrl}/items/${itemId}`, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error removing cart item', error);
      throw error;
    }
  }

  // Clear the cart
  async clearCart(): Promise<StatusResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.delete(this.baseUrl, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error clearing cart', error);
      throw error;
    }
  }

  // Apply a coupon to the cart
  async applyCoupon(couponCode: string): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(`${this.baseUrl}/coupon`, { couponCode }, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error applying coupon', error);
      throw error;
    }
  }

  // Remove a coupon from the cart
  async removeCoupon(): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.delete(`${this.baseUrl}/coupon`, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error removing coupon', error);
      throw error;
    }
  }

  // Sync local cart with server (used when a user logs in)
  async syncCart(items: CartItem[]): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(`${this.baseUrl}/sync`, { items }, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError('Error syncing cart', error);
      throw error;
    }
  }

  // Handle API errors
  private handleError(message: string, error: any): void {
    console.error(`${message}:`, error);

    if (axios.isAxiosError(error)) {
      // Handle specific error cases
      if (error.response?.status === 401) {
        // Unauthorized error
        throw new Error('Please login to manage your cart');
      } else if (error.response?.status === 400) {
        // Bad request
        throw new Error(error.response.data.message || 'Invalid cart operation');
      } else if (error.response?.status === 404) {
        // Not found
        throw new Error('Cart item not found');
      } else if (error.response?.data?.message) {
        // Server provided an error message
        throw new Error(error.response.data.message);
      }
    }

    // Default error message
    throw new Error('Failed to process cart operation. Please try again.');
  }
}

// Export a singleton instance
const cartService = new CartService();
export default cartService;
