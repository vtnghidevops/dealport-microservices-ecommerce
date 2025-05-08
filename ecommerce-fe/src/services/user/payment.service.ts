import axios from 'axios';

// Base API URL already includes the /api/v1 prefix
const API_URL = import.meta.env.API_URL || 'http://localhost:8080/api/v1';

// Get auth header for authenticated requests
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  if (!token) return {};
  return { Authorization: `Bearer ${token}` };
};

// Interface for generic payment requests
interface PaymentRequest {
  orderId: string;
  amount: number;
  paymentMethod: string;
  returnUrl: string;
  currency?: string;
}

// Interface for generic payment results
interface PaymentResult {
  success: boolean;
  transactionId?: string;
  status?: string;
  redirectUrl?: string;
  message?: string;
}

class PaymentService {
  // The base URL should just add the payments path without double adding /api/v1
  private baseUrl = `${API_URL}/payments`;
  private checkoutBaseUrl = `${API_URL}/checkout/payments`;

  /**
   * Generic method to process a payment for an order with any payment method
   * @param paymentRequest - Payment request details
   * @returns Payment result including redirect URL if needed
   */
  async processPayment(paymentRequest: PaymentRequest): Promise<PaymentResult> {
    try {
      const headers = getAuthHeader();

      console.log('Processing payment for order:', paymentRequest.orderId);
      console.log('Payment details:', paymentRequest);

      // Default currency to USD if not provided
      if (!paymentRequest.currency) {
        paymentRequest.currency = 'USD';
      }

      const response = await axios.post(
        `${this.checkoutBaseUrl}/process`,
        paymentRequest,
        { headers }
      );

      console.log('Payment response:', response.data);

      if (response.data && response.data.success) {
        return {
          success: true,
          transactionId: response.data.transactionId,
          status: response.data.status,
          redirectUrl: response.data.redirectUrl,
          message: response.data.message || 'Payment initiated successfully'
        };
      }

      return {
        success: false,
        message: response.data?.message || 'Payment failed'
      };
    } catch (error) {
      console.error('Error processing payment:', error);

      // Handle Axios errors
      if (axios.isAxiosError(error)) {
        const errorMessage = error.response?.data?.message || error.message;
        return {
          success: false,
          message: `Payment failed: ${errorMessage}`
        };
      }

      return {
        success: false,
        message: 'An unexpected error occurred during payment processing'
      };
    }
  }

  /**
   * Check the status of a payment
   * @param transactionId - Transaction ID to check
   * @returns Payment status result
   */
  async checkPaymentStatus(transactionId: string): Promise<PaymentResult> {
    try {
      const headers = getAuthHeader();

      const response = await axios.get(
        `${this.checkoutBaseUrl}/status/${transactionId}`,
        { headers }
      );

      console.log('Payment status response:', response.data);

      if (response.data && response.data.success) {
        return {
          success: true,
          status: response.data.status,
          message: response.data.message || 'Payment status retrieved successfully'
        };
      }

      return {
        success: false,
        message: response.data?.message || 'Failed to get payment status'
      };
    } catch (error) {
      console.error('Error checking payment status:', error);

      if (axios.isAxiosError(error)) {
        const errorMessage = error.response?.data?.message || error.message;
        return {
          success: false,
          message: `Error checking payment status: ${errorMessage}`
        };
      }

      return {
        success: false,
        message: 'An unexpected error occurred while checking payment status'
      };
    }
  }

  /**
   * Create MoMo payment request using ATM banking method (direct bank transfer)
   * @param orderId - Order ID to be paid
   * @param amount - Amount to pay in VND
   * @param returnUrl - URL to redirect after payment completion
   * @returns Payment session data including payment URL
   */
  async createMomoPayment(
    orderId: string,
    amount: number,
    returnUrl: string
  ): Promise<{
    paymentUrl: string;
    orderId: string;
    requestId: string;
    transactionId?: string;
    amount: number;
  }> {
    try {
      const headers = getAuthHeader();
      const payload = {
        orderId,
        amount,
        orderInfo: `Payment for order ${orderId}`,
        returnUrl
      };

      // Log payment request details
      console.log('Creating MoMo ATM payment with details:', {
        url: `${this.baseUrl}/momo/create`,
        payload
      });

      const response = await axios.post(
        `${this.baseUrl}/momo/create`,
        payload,
        { headers }
      );

      // Log successful response
      console.log('MoMo ATM payment created successfully:', response.data);

      if (!response.data.data) {
        throw new Error('Invalid response format from payment service');
      }

      return response.data.data;
    } catch (error) {
      // Log error before handling
      console.error('Error in createMomoPayment:', error);

      // Check if the error response contains detailed information
      if (axios.isAxiosError(error) && error.response?.data) {
        console.error('Error details from server:', error.response.data);

        // If server sent specific error message, use that
        if (error.response.data.message) {
          throw new Error(`MoMo payment error: ${error.response.data.message}`);
        }
      }

      this.handleError('Error creating MoMo payment', error);
      throw error;
    }
  }

  /**
   * Verify MoMo payment status after redirect
   * @param params - URL parameters from MoMo redirect
   * @returns Payment verification result
   */
  async verifyMomoPayment(params: Record<string, string>): Promise<{
    success: boolean;
    orderId: string;
    message: string;
    transactionId?: string;
    amount?: number;
  }> {
    try {
      const headers = getAuthHeader();

      // Log verification request
      console.log('Verifying MoMo payment with params:', params);

      const response = await axios.post(
        `${this.baseUrl}/momo/verify`,
        { params },
        { headers }
      );

      // Log successful verification
      console.log('MoMo payment verification result:', response.data);

      if (!response.data.data) {
        throw new Error('Invalid response format from verification service');
      }

      return response.data.data;
    } catch (error) {
      console.error('Error in verifyMomoPayment:', error);

      // Check if the error response contains detailed information
      if (axios.isAxiosError(error) && error.response?.data) {
        console.error('Error details from server:', error.response.data);
      }

      this.handleError('Error verifying MoMo payment', error);
      throw error;
    }
  }

  // Handle API errors
  private handleError(message: string, error: any): void {
    console.error(`${message}:`, error);

    if (axios.isAxiosError(error)) {
      if (error.response?.status === 401) {
        throw new Error('Vui lòng đăng nhập để thực hiện thanh toán');
      } else if (error.response?.status === 404) {
        throw new Error('Không tìm thấy dịch vụ thanh toán. Vui lòng thử lại sau.');
      } else if (error.response?.data?.message) {
        throw new Error(error.response.data.message);
      } else if (error.response) {
        throw new Error(`Lỗi từ máy chủ: ${error.response.status} - ${error.response.statusText}`);
      } else if (error.request) {
        throw new Error('Không thể kết nối đến máy chủ. Vui lòng kiểm tra kết nối mạng.');
      }
    }

    throw new Error('Không thể xử lý thanh toán. Vui lòng thử lại sau.');
  }
}

// Export singleton instance
const paymentService = new PaymentService();
export default paymentService; 