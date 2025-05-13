// components/user/OrderCard.tsx
import React, { useState } from 'react';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { FiChevronDown, FiChevronUp } from 'react-icons/fi';
import { Order } from '@/services/user/order.service';
import { useNavigate } from 'react-router-dom';
import { differenceInMinutes } from 'date-fns';

interface OrderCardProps {
  order: Order;
  onRetryPayment?: (orderId: string) => void;
}

const OrderCard: React.FC<OrderCardProps> = ({ order, onRetryPayment }) => {
  const [expanded, setExpanded] = useState(false);
  const navigate = useNavigate();

  // Log the order on component mount
  // useEffect(() => {
  //   console.log("OrderCard: Received order to render:", order);
  //   console.log("OrderCard: Order items:", order.items);
  //   console.log("OrderCard: Order status:", order.status);
  //   console.log("OrderCard: Order dates - created:", order.createdAt, "updated:", order.updatedAt);
  // }, [order]);

  // Check if the pending online payment is still valid (within 1 hour)
  const isPaymentStillValid = () => {
    if (order.status !== 'pending' || !isOnlinePayment(order.paymentMethod)) {
      return false;
    }

    try {
      const createdAt = new Date(order.createdAt);
      const now = new Date();
      const minutesDifference = differenceInMinutes(now, createdAt);

      console.log(`OrderCard: Payment created ${minutesDifference} minutes ago`);
      // Return true if created less than 60 minutes ago
      return minutesDifference < 60;
    } catch (error) {
      console.error('Error calculating time difference:', error);
      return false;
    }
  };

  // Check if payment method is an online payment
  const isOnlinePayment = (method: string | undefined) => {
    if (!method) return false;

    const onlinePaymentMethods = ['momo', 'vnpay'];
    return onlinePaymentMethods.includes(method.toLowerCase());
  };

  // Get a friendly payment method name for display
  const getPaymentMethodDisplay = (method: string | undefined) => {
    if (!method || method === '') return 'COD';

    // Map payment method to user-friendly names
    const methodMap: { [key: string]: string } = {
      'momo': 'MoMo',
      'vnpay': 'VNPay',
      'cod': 'Cash On Delivery'
    };

    // Convert to lowercase for case-insensitive comparison
    const methodLower = method.toLowerCase();
    return methodMap[methodLower] || method;
  };

  const handleRetryPayment = () => {
    if (onRetryPayment) {
      onRetryPayment(order.id);
    } else {
      // Fallback to direct navigation
      navigate(`/checkout/payment/${order.id}`);
    }
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) {
      // console.log("OrderCard: Missing date string");
      return 'Unknown date';
    }
    try {
      // console.log(`OrderCard: Formatting date string: "${dateString}"`);
      const date = new Date(dateString);
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    } catch (error) {
      // console.error('OrderCard: Error formatting date:', error);
      return 'Invalid date';
    }
  };

  //   const getStatusColor = (status: string | undefined) => {
  //     if (!status) {
  //       console.log("OrderCard: Missing status");
  //       return 'bg-gray-100 text-gray-800';
  //     }

  //   switch (status.toLowerCase()) {
  //     case 'pending':
  //       return 'bg-purple-100 text-purple-800';
  //     case 'processing':
  //       return 'bg-yellow-100 text-yellow-800';
  //     case 'paid':
  //       return 'bg-green-100 text-green-800';
  //     case 'shipped':
  //       return 'bg-blue-100 text-blue-800';
  //     case 'delivered':
  //       return 'bg-green-100 text-green-800';
  //     case 'cancelled':
  //       return 'bg-red-100 text-red-800';
  //     case 'refunded':
  //       return 'bg-orange-100 text-orange-800';
  //     default:
  //       console.log(`OrderCard: Unknown status "${status}"`);
  //       return 'bg-gray-100 text-gray-800';
  //   }
  // };

  const getPaymentStatusColor = (status: string | undefined) => {
    if (!status) {
      console.log("OrderCard: Missing payment status");
      return 'bg-gray-100 text-gray-800';
    }

    switch (status.toLowerCase()) {
      case 'paid':
        return 'bg-green-100 text-green-800';
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        console.log(`OrderCard: Unknown payment status "${status}"`);
        return 'bg-gray-100 text-gray-800';
    }
  };

  // Helper function to safely access nested properties
  const safelyGetItems = () => {
    // console.log("OrderCard: Checking items array:", order.items);
    return Array.isArray(order.items) ? order.items : [];
  }

  // Safely get total with fallback
  const getTotal = () => {
    if (typeof order.total === 'number') {
      return order.total;
    }
    // Try to calculate from items if total is missing
    const items = safelyGetItems();
    if (items.length > 0) {
      return items.reduce((sum, item) => {
        const price = typeof item.price === 'number' ? item.price : 0;
        const quantity = typeof item.quantity === 'number' ? item.quantity : 1;
        return sum + (price * quantity);
      }, 0);
    }
    return 0;
  };

  // Safely format a number with fallback
  const formatPrice = (price: number | undefined) => {
    if (typeof price !== 'number') return '$0.00';
    return `$${price.toFixed(2)}`;
  };

  // Check if payment can still be attempted
  const canRetryPayment = isPaymentStillValid();

  return (
    <Card className='p-[10px] mb-5'>
      <CardHeader className="pb-2">
        <div className="flex flex-col md:flex-row justify-between">
          <div>
            <CardTitle className="text-lg font-bold">Order #{order.id || 'Unknown'}</CardTitle>
            <p className="text-[15px] text-gray-500">Placed on {formatDate(order.createdAt)}</p>
          </div>
          <div className="flex items-center mt-2 md:mt-0 gap-2">
            {/* <Badge className={`${getStatusColor(order.status)} capitalize`}>
              {order.status || 'Unknown'} oke
            </Badge> */}
            <Badge className={`${getPaymentStatusColor(order.status)} capitalize`}>
              {order.status || 'Unknown'}
            </Badge>
            {isOnlinePayment(order.paymentMethod) && (
              <Badge className="bg-blue-100 text-blue-800 capitalize">
                {order.paymentMethod}
              </Badge>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="pb-2">
        <div className="flex justify-between mb-2">
          <span className="text-[15px] text-gray-500">Items: {safelyGetItems().length}</span>
          <span className="font-bold">{formatPrice(getTotal())}</span>
        </div>

        {expanded && (
          <div className="mt-4 space-y-3 border-t pt-3">
            {safelyGetItems().map((item, index) => (
              <div key={item.id || index} className="flex items-center">
                {item.imageUrl && (
                  <div className="mr-3">
                    <img
                      src={item.imageUrl}
                      alt={item.name || 'Product'}
                      className="w-[30px] h-[30px] object-cover rounded"
                    />
                  </div>
                )}
                <div className="flex-1">
                  <p className="font-medium">{item.name || 'Unknown Product'}</p>
                  <div className="text-sm text-gray-500">
                    {item.quantity || 1} x {formatPrice(item.price)}
                  </div>
                </div>
                <div className="font-bold">
                  {formatPrice((item.price || 0) * (item.quantity || 1))}
                </div>
              </div>
            ))}

            {order.transactionId && (
              <div className="pt-2 border-t mt-2">
                <p className="text-sm text-gray-600">
                  <span className="font-medium">Payment Method:</span> {getPaymentMethodDisplay(order.paymentMethod)}
                </p>
                <p className="text-sm text-gray-600">
                  <span className="font-medium">Transaction ID:</span> {order.transactionId}
                </p>
              </div>
            )}

            {!order.transactionId && order.paymentMethod && (
              <div className="pt-2 border-t mt-2">
                <p className="text-sm text-gray-600">
                  <span className="font-medium">Payment Method:</span> {getPaymentMethodDisplay(order.paymentMethod)}
                </p>
                {canRetryPayment && (
                  <p className="text-xs text-orange-600 mt-1">
                    This payment can be retried within {60 - differenceInMinutes(new Date(), new Date(order.createdAt))} minutes
                  </p>
                )}
              </div>
            )}
          </div>
        )}
      </CardContent>

      <CardFooter className="pt-2 flex justify-between">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setExpanded(!expanded)}
          className="flex items-center text-[15px]"
        >
          {expanded ? (
            <>
              <FiChevronUp className="mr-1 !w-[15px] !h-[15px]" /> Hide Details
            </>
          ) : (
            <>
              <FiChevronDown className="mr-1 !w-[15px] !h-[15px]" /> View Details
            </>
          )}
        </Button>
        <div className="flex gap-2">
          {canRetryPayment && (
            <Button
              variant="destructive"
              className='bg-[#ff6b6b] text-white hover:bg-red-400 w-[120px] h-[40px] text-[14px] p-2'
              onClick={handleRetryPayment}
            >
              Pay Now
            </Button>
          )}
          {order.status !== 'delivered' && order.status !== 'cancelled' && order.status !== 'refunded' && (
            <Button
              variant="outline"
              className='bg-[#0496FF] text-white hover:bg-blue-600 w-[120px] h-[40px] text-[14px] p-2'
            >
              Track Order
            </Button>
          )}
        </div>
      </CardFooter>
    </Card>
  );
};

export default OrderCard;