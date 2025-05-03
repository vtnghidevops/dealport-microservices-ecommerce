import { PaymentInfo } from '@/types/checkout.model';

/**
 * Get a user-friendly payment status text based on payment method and status
 */
export function getPaymentStatusText(paymentInfo: PaymentInfo): string {
  const { paymentMethod, status } = paymentInfo;

  // Handle COD (Cash on Delivery) special cases
  if (paymentMethod === 'cod') {
    if (status === 'pending' || status === 'pending_delivery' || status === 'cod_pending') {
      return 'Payment will be collected upon delivery';
    }
    if (status === 'completed') {
      return 'Payment collected on delivery';
    }
    return 'Cash on delivery';
  }

  // Handle Online payment methods (MoMo, VNPay, etc.)
  switch (status) {
    case 'pending':
      return 'Payment pending';
    case 'processing':
      return 'Payment being processed';
    case 'completed':
      return 'Payment successful';
    case 'failed':
      return 'Payment failed';
    case 'refunded':
      return 'Payment refunded';
    default:
      return status || 'Unknown payment status';
  }
}

/**
 * Get a color for the payment status badge
 */
export function getPaymentStatusColor(paymentInfo: PaymentInfo): string {
  const { paymentMethod, status } = paymentInfo;

  // For COD payments
  if (paymentMethod === 'cod') {
    if (status === 'pending' || status === 'pending_delivery' || status === 'cod_pending') {
      return 'orange';
    }
    if (status === 'completed') {
      return 'green';
    }
    return 'blue';
  }

  // For other payment methods
  switch (status) {
    case 'pending':
      return 'orange';
    case 'processing':
      return 'blue';
    case 'completed':
      return 'green';
    case 'failed':
      return 'red';
    case 'refunded':
      return 'purple';
    default:
      return 'gray';
  }
}

/**
 * Returns a standardized payment method display name
 */
export function getPaymentMethodDisplayName(method: string): string {
  switch (method?.toLowerCase()) {
    case 'cod':
      return 'Cash on Delivery';
    case 'momo':
      return 'MoMo Wallet';
    case 'vnpay':
      return 'VNPay';
    case 'card':
      return 'Credit/Debit Card';
    case 'bank_transfer':
      return 'Bank Transfer';
    default:
      return method || 'Unknown';
  }
} 