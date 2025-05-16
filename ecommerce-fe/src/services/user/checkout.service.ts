import axios from 'axios';
import { CartItem } from '@/types/cart.model';
import {
  BillingInfo,
  ShippingInfo,
  //PaymentInfo,
  CheckoutOrder,
  //ValidationError,
  CheckoutValidationResponse,
  PaymentResponse
} from '@/types/checkout.model';
import { getApiUrl, getAuthHeader } from '@/utils/api-config';

// Singleton instance of checkout service
class CheckoutService {
  private baseUrl = getApiUrl('checkout');

  // Validate checkout data before placing order
  async validateCheckout(
    items: CartItem[],
    billingInfo: BillingInfo,
    shippingInfo: ShippingInfo,
    paymentMethod: string
  ): Promise<CheckoutValidationResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/validate`,
        {
          items,
          billingInfo,
          shippingInfo,
          paymentMethod
        },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError('Error validating checkout', error);
      throw error;
    }
  }

  // Create a new order
  async createOrder(
    items: CartItem[],
    billingInfo: BillingInfo,
    shippingInfo: ShippingInfo,
    paymentMethod: string,
    couponCode?: string,
    notes?: string
  ): Promise<CheckoutOrder> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/orders`,
        {
          items,
          billingInfo,
          shippingInfo,
          paymentMethod,
          couponCode,
          notes
        },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError('Error creating order', error);
      throw error;
    }
  }

  // Process payment for an order
  async processPayment(
    orderId: string,
    paymentMethod: string,
    amount: number,
    currency: string = 'USD',
    returnUrl?: string
  ): Promise<PaymentResponse> {
    try {
      const headers = getAuthHeader();
      const response = await axios.post(
        `${this.baseUrl}/orders/${orderId}/payment`,
        {
          paymentMethod,
          amount,
          currency,
          returnUrl
        },
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError('Error processing payment', error);
      throw error;
    }
  }

  // Get an order by ID
  async getOrder(orderId: string): Promise<CheckoutOrder> {
    try {
      const headers = getAuthHeader();
      const response = await axios.get(
        `${this.baseUrl}/orders/${orderId}`,
        { headers }
      );
      return response.data.data;
    } catch (error) {
      this.handleError('Error fetching order', error);
      throw error;
    }
  }

  // List user's orders with pagination
  async listOrders(page: number = 1, pageSize: number = 10): Promise<{
    orders: CheckoutOrder[];
    total: number;
    page: number;
    size: number;
  }> {
    try {
      const headers = getAuthHeader();
      const response = await axios.get(
        `${this.baseUrl}/orders`,
        {
          headers,
          params: { page, pageSize }
        }
      );
      return response.data.data;
    } catch (error) {
      this.handleError('Error fetching orders', error);
      throw error;
    }
  }

  // Handle API errors
  private handleError(message: string, error: any): void {
    console.error(`${message}:`, error);

    if (axios.isAxiosError(error)) {
      // Handle specific error cases
      if (error.response?.status === 401) {
        throw new Error('Please login to complete checkout');
      } else if (error.response?.status === 400) {
        throw new Error(error.response.data.message || 'Invalid checkout data');
      } else if (error.response?.status === 404) {
        throw new Error('Order not found');
      } else if (error.response?.data?.message) {
        throw new Error(error.response.data.message);
      }
    }

    // Default error message
    throw new Error('Failed to process checkout. Please try again.');
  }
}

// Export a singleton instance
const checkoutService = new CheckoutService();
export default checkoutService; 