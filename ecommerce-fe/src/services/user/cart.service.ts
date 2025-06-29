import axios from "axios";
import { CartItem, CartTotalsData } from "@/types/cart.model";
import { getApiUrl, getAuthHeader } from "@/utils/api-config";

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

// Helper function to get current user ID from localStorage
// const getCurrentUserId = (): string | null => {
//   const userStr = localStorage.getItem('user');
//   if (!userStr) return null;

//   try {
//     const user = JSON.parse(userStr);
//     return user?.id || null;
//   } catch (e) {
//     console.error('Error parsing user from localStorage', e);
//     return null;
//   }
// };

// Singleton instance of cart service
class CartService {
  private baseUrl = getApiUrl("cart");

  // Get the user's cart
  async getCart(): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.get(this.baseUrl, { headers });
      return response.data.data;
    } catch (error) {
      this.handleError("Error getting cart", error);
      throw error;
    }
  }

  // Add an item to the cart
  async addCartItem(item: CartItemRequest): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/items`,
        {
          productId: item.productId,
          name: item.name,
          price: item.price,
          originalPrice: item.originalPrice,
          quantity: item.quantity,
          imageUrl: item.imageUrl,
        },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError("Error adding item to cart", error);
      throw error;
    }
  }

  // Update an item in the cart
  async updateCartItem(
    itemId: string,
    quantity: number
  ): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.put(
        `${this.baseUrl}/items/${itemId}`,
        { quantity },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError("Error updating cart item", error);
      throw error;
    }
  }

  // Remove an item from the cart
  async removeCartItem(itemId: string): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.delete(`${this.baseUrl}/items/${itemId}`, {
        headers,
      });
      return response.data.data;
    } catch (error) {
      this.handleError("Error removing cart item", error);
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
      this.handleError("Error clearing cart", error);
      throw error;
    }
  }

  // Apply a coupon to the cart
  async applyCoupon(couponCode: string): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/coupon`,
        { couponCode },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError("Error applying coupon", error);
      throw error;
    }
  }

  // Remove a coupon from the cart
  async removeCoupon(): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.delete(`${this.baseUrl}/coupon`, {
        headers,
      });
      return response.data.data;
    } catch (error) {
      this.handleError("Error removing coupon", error);
      throw error;
    }
  }

  // Refresh cart TTL (time-to-live)
  async refreshCartTTL(): Promise<StatusResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/refresh-ttl`,
        {},
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError("Error refreshing cart TTL", error);
      throw error;
    }
  }

  // Sync local cart with server (used when a user logs in)
  async syncCart(items: CartItem[]): Promise<CartResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/sync`,
        { items },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError("Error syncing cart", error);
      throw error;
    }
  }

  // Handle API errors
  private handleError(message: string, error: unknown): void {
    console.error(`${message}:`, error);

    if (axios.isAxiosError(error)) {
      // Handle specific error cases
      if (error.response?.status === 401) {
        // Unauthorized error
        throw new Error("Please login to manage your cart");
      } else if (error.response?.status === 400) {
        // Bad request
        throw new Error(
          error.response.data.message || "Invalid cart operation"
        );
      } else if (error.response?.status === 404) {
        // Not found
        throw new Error("Cart item not found");
      } else if (error.response?.data?.message) {
        // Server provided an error message
        throw new Error(error.response.data.message);
      }
    }

    // Default error message
    throw new Error("Failed to process cart operation. Please try again.");
  }
}

// Export a singleton instance
const cartService = new CartService();
export default cartService;
