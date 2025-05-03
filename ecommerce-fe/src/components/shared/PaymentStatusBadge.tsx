import React from 'react';
import { PaymentInfo } from '@/types/checkout.model';
import { getPaymentStatusText, getPaymentStatusColor } from '@/utils/payment-helpers';

interface PaymentStatusBadgeProps {
  paymentInfo: PaymentInfo;
  className?: string;
}

/**
 * A component to display payment status with appropriate styling
 */
const PaymentStatusBadge: React.FC<PaymentStatusBadgeProps> = ({ paymentInfo, className = '' }) => {
  const statusText = getPaymentStatusText(paymentInfo);
  const colorClass = getColorClass(getPaymentStatusColor(paymentInfo));

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClass} ${className}`}>
      {statusText}
    </span>
  );
};

// Helper function to convert color name to Tailwind CSS class
function getColorClass(color: string): string {
  switch (color) {
    case 'green':
      return 'bg-green-100 text-green-800';
    case 'red':
      return 'bg-red-100 text-red-800';
    case 'blue':
      return 'bg-blue-100 text-blue-800';
    case 'orange':
      return 'bg-orange-100 text-orange-800';
    case 'purple':
      return 'bg-purple-100 text-purple-800';
    case 'gray':
    default:
      return 'bg-gray-100 text-gray-800';
  }
}

export default PaymentStatusBadge; 