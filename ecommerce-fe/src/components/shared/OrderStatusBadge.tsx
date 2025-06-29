import React from 'react';

interface OrderStatusBadgeProps {
  status: string;
  className?: string;
}

/**
 * Get the display text for an order status
 */
export function getOrderStatusText(status: string): string {
  switch (status?.toLowerCase()) {
    case 'pending':
      return 'Pending';
    case 'processing':
      return 'Processing';
    case 'paid':
      return 'Paid';
    case 'shipped':
      return 'Shipped';
    case 'delivered':
      return 'Delivered';
    case 'cancelled':
    case 'canceled':
      return 'Cancelled';
    case 'refunded':
      return 'Refunded';
    default:
      return status || 'Unknown';
  }
}

/**
 * Get the color for an order status badge
 */
export function getOrderStatusColor(status: string): string {
  switch (status?.toLowerCase()) {
    case 'pending':
      return 'yellow';
    case 'processing':
      return 'blue';
    case 'paid':
      return 'teal';
    case 'shipped':
      return 'indigo';
    case 'delivered':
      return 'green';
    case 'cancelled':
    case 'canceled':
      return 'red';
    case 'refunded':
      return 'purple';
    default:
      return 'gray';
  }
}

/**
 * A component to display order status with appropriate styling
 */
const OrderStatusBadge: React.FC<OrderStatusBadgeProps> = ({ status, className = '' }) => {
  const statusText = getOrderStatusText(status);
  const colorClass = getColorClass(getOrderStatusColor(status));

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
    case 'yellow':
      return 'bg-yellow-100 text-yellow-800';
    case 'red':
      return 'bg-red-100 text-red-800';
    case 'blue':
      return 'bg-blue-100 text-blue-800';
    case 'teal':
      return 'bg-teal-100 text-teal-800';
    case 'indigo':
      return 'bg-indigo-100 text-indigo-800';
    case 'purple':
      return 'bg-purple-100 text-purple-800';
    case 'gray':
    default:
      return 'bg-gray-100 text-gray-800';
  }
}

export default OrderStatusBadge; 