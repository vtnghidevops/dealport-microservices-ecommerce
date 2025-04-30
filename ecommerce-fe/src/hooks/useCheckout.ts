import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useToast } from './use-toast';
import { useAuth } from './useAuth';
import { useCart } from './useCart';
import checkoutService from '@/services/user/checkout.service';
import paymentService from '@/services/user/payment.service'

import {
  BillingInfo,
  ShippingInfo,
  CheckoutOrder,
  ValidationError,
  CheckoutValidationResponse,
  PaymentResponse
} from '@/types/checkout.model';

export interface UseCheckoutReturn {
  isLoading: boolean;
  error: string | null;
  validateCheckout: (
    billingInfo: BillingInfo,
    shippingInfo: ShippingInfo,
    paymentMethod: string
  ) => Promise<CheckoutValidationResponse>;
  createOrder: (
    billingInfo: BillingInfo,
    shippingInfo: ShippingInfo,
    paymentMethod: string,
    notes?: string
  ) => Promise<CheckoutOrder>;
  processPayment: (
    orderId: string,
    paymentMethod: string,
    amount: number,
    returnUrl?: string
  ) => Promise<PaymentResponse>;
  getOrder: (orderId: string) => Promise<CheckoutOrder>;
  listOrders: (page?: number, pageSize?: number) => Promise<{
    orders: CheckoutOrder[];
    total: number;
    page: number;
    size: number;
  }>;
}

export const useCheckout = (): UseCheckoutReturn => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { authState } = useAuth();
  const { cartItems, cartTotals, clearCart } = useCart();
  const { toast } = useToast();
  const navigate = useNavigate();

  const handleError = (error: unknown) => {
    const errorMessage = error instanceof Error
      ? error.message
      : 'An error occurred during checkout';

    setError(errorMessage);
    toast({
      title: 'Checkout Error',
      description: errorMessage,
      variant: 'destructive'
    });

    return errorMessage;
  };

  const validateCheckout = async (
    billingInfo: BillingInfo,
    shippingInfo: ShippingInfo,
    paymentMethod: string
  ): Promise<CheckoutValidationResponse> => {
    try {
      setIsLoading(true);
      setError(null);

      // Gọi API validate checkout
      const result = await checkoutService.validateCheckout(
        cartItems,
        billingInfo,
        shippingInfo,
        paymentMethod
      );

      if (!result.valid && result.errors.length > 0) {
        // Hiển thị lỗi từ validation
        const firstError = result.errors[0];
        toast({
          title: 'Validation Error',
          description: `${firstError.field}: ${firstError.message}`,
          variant: 'destructive'
        });
      }

      return result;
    } catch (error) {
      handleError(error);
      return { valid: false, errors: [{ field: 'system', message: handleError(error) }] };
    } finally {
      setIsLoading(false);
    }
  };

  const createOrder = async (
    billingInfo: BillingInfo,
    shippingInfo: ShippingInfo,
    paymentMethod: string,
    notes?: string
  ): Promise<CheckoutOrder> => {
    if (!authState.isAuthenticated) {
      const errorMessage = 'You must be logged in to create an order';
      setError(errorMessage);
      toast({
        title: 'Authentication Required',
        description: errorMessage,
        variant: 'destructive'
      });

      // Redirect to login
      navigate('/auth/login?redirect=checkout');

      throw new Error(errorMessage);
    }

    try {
      setIsLoading(true);
      setError(null);

      // Get coupon code from local storage if available
      // This is a workaround since we don't have direct access to the coupon code from the cart
      const couponCode = localStorage.getItem('appliedCoupon') || undefined;

      // Gọi API create order
      const orderResult = await checkoutService.createOrder(
        cartItems,
        billingInfo,
        shippingInfo,
        paymentMethod,
        couponCode,
        notes
      );

      // Nếu tạo đơn hàng thành công, xóa giỏ hàng và coupon code khỏi localStorage
      await clearCart();
      localStorage.removeItem('appliedCoupon');

      toast({
        title: 'Order Created',
        description: `Order #${orderResult.orderNumber} has been created successfully`,
        variant: 'success'
      });

      return orderResult;
    } catch (error) {
      handleError(error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const processPayment = async (
    orderId: string,
    paymentMethod: string,
    amount: number,
    returnUrl?: string
  ): Promise<PaymentResponse> => {
    try {
      setIsLoading(true);
      setError(null);

      // Default return URL if not specified
      const defaultReturnUrl = window.location.origin + '/user/checkout/success';
      const effectiveReturnUrl = returnUrl || defaultReturnUrl;

      // Process payment based on payment method
      if (paymentMethod === 'momo') {
        try {
          // Use MoMo QuickPay payment service
          const momoPayment = await paymentService.createMomoPayment(
            orderId,
            Math.round(amount * 23000), // Convert USD to VND (approximate rate)
            effectiveReturnUrl
          );

          // Redirect to MoMo payment page
          if (momoPayment && momoPayment.paymentUrl) {
            // Log payment info before redirect
            console.log(`Redirecting to MoMo payment URL: ${momoPayment.paymentUrl}`);
            console.log(`MoMo payment details:`, {
              orderId: momoPayment.orderId,
              amount: momoPayment.amount,
              requestId: momoPayment.requestId
            });

            // Redirect to payment page
            window.location.href = momoPayment.paymentUrl;
            return {
              success: true,
              status: 'redirecting',
              redirectUrl: momoPayment.paymentUrl,
              orderId: orderId
            };
          } else {
            throw new Error('MoMo payment URL not provided in response');
          }
        } catch (error) {
          console.error('Error creating MoMo payment:', error);
          toast({
            title: 'Payment Error',
            description: error instanceof Error ? error.message : 'Failed to initiate MoMo payment',
            variant: 'destructive'
          });
          throw error;
        }
      } else if (paymentMethod === 'cod') {
        // For Cash on Delivery (COD)
        const paymentResult = await checkoutService.processPayment(
          orderId,
          paymentMethod,
          amount,
          'USD',
          effectiveReturnUrl
        );

        if (paymentResult.success) {
          toast({
            title: 'Order Confirmed',
            description: 'Your order has been confirmed and will be processed',
            variant: 'success'
          });

          // For COD, go directly to success page
          navigate(`/user/checkout/success?orderId=${orderId}&method=${paymentMethod}`);
        } else {
          toast({
            title: 'Order Processing Failed',
            description: paymentResult.status || 'There was an issue processing your order',
            variant: 'destructive'
          });
        }

        return paymentResult;
      } else {
        throw new Error('Unsupported payment method');
      }

      // If execution reaches here, something went wrong with the payment setup
      throw new Error('Failed to initialize payment. Please try again.');

    } catch (error) {
      handleError(error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const getOrder = async (orderId: string): Promise<CheckoutOrder> => {
    try {
      setIsLoading(true);
      setError(null);

      // Gọi API get order
      const order = await checkoutService.getOrder(orderId);
      return order;
    } catch (error) {
      handleError(error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const listOrders = async (page: number = 1, pageSize: number = 10) => {
    try {
      setIsLoading(true);
      setError(null);

      // Gọi API list orders
      const orders = await checkoutService.listOrders(page, pageSize);
      return orders;
    } catch (error) {
      handleError(error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isLoading,
    error,
    validateCheckout,
    createOrder,
    processPayment,
    getOrder,
    listOrders
  };
}; 