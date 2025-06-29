import React, { useEffect, useState } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useToast } from '../../hooks/use-toast';
import { Button } from '../../components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../../components/ui/card';
import { Skeleton } from '../../components/ui/skeleton';
import { BiError } from "react-icons/bi";
import { FiArrowLeft } from "react-icons/fi";
import orderService from '@/services/user/order.service';
import paymentService from '@/services/user/payment.service';

const RetryPayment: React.FC = () => {
  const { orderId } = useParams<{ orderId: string }>();
  const location = useLocation();
  const navigate = useNavigate();
  const { authState } = useAuth();
  const { toast } = useToast();

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [order, setOrder] = useState<any>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  // Payment details can be either from location state or fetched from the order
  const paymentDetails = location.state || {};

  useEffect(() => {
    // Fetch order details if not provided in location state
    const fetchOrderDetails = async () => {
      if (!authState.isAuthenticated) {
        setError('Please login to retry payment');
        setIsLoading(false);
        return;
      }

      if (!orderId) {
        setError('Order ID is required');
        setIsLoading(false);
        return;
      }

      try {
        setIsLoading(true);
        const orderData = await orderService.getOrderById(orderId);

        if (!orderData) {
          setError('Order not found');
          return;
        }

        // Check if order is eligible for retry payment
        if (orderData.status !== 'pending') {
          setError('This order is not eligible for payment retry');
          return;
        }

        // Check if order is still within time window (1 hour)
        const createdAt = new Date(orderData.createdAt);
        const now = new Date();
        const timeDiff = now.getTime() - createdAt.getTime();
        const hoursDiff = timeDiff / (1000 * 60 * 60);

        if (hoursDiff > 1) {
          setError('Payment window has expired. Please create a new order.');
          return;
        }

        setOrder(orderData);
      } catch (err) {
        console.error('Error fetching order:', err);
        setError(err instanceof Error ? err.message : 'Failed to fetch order details');
      } finally {
        setIsLoading(false);
      }
    };

    fetchOrderDetails();
  }, [orderId, authState.isAuthenticated]);

  const handleRetryPayment = async () => {
    if (!order || !orderId) {
      toast({
        title: 'Error',
        description: 'Order information is missing',
        variant: 'destructive'
      });
      return;
    }

    try {
      setIsProcessing(true);

      // Determine which payment method to use
      // Default to momo if unknown or not provided
      const paymentMethod = paymentDetails.paymentMethod || order.paymentMethod || 'momo';

      // Create payment request
      const paymentRequest = {
        orderId: orderId,
        amount: order.total,
        paymentMethod: paymentMethod,
        returnUrl: `${window.location.origin}/user/checkout/success`,
      };

      // Process payment
      const result = await paymentService.processPayment(paymentRequest);

      if (result.success) {
        // If redirect URL is provided, redirect to payment gateway
        if (result.redirectUrl) {
          window.location.href = result.redirectUrl;
        } else {
          // Otherwise navigate to success page with status "paid" instead of "success" to match schema
          navigate('/user/checkout/success', {
            state: {
              orderId: orderId,
              transactionId: result.transactionId,
              method: paymentMethod,
              status: 'paid' // Set status to paid according to schema
            }
          });
        }
      } else {
        toast({
          title: 'Payment Failed',
          description: result.message || 'Unable to process payment',
          variant: 'destructive'
        });
      }
    } catch (err) {
      console.error('Error processing payment:', err);
      toast({
        title: 'Payment Error',
        description: err instanceof Error ? err.message : 'Failed to process payment',
        variant: 'destructive'
      });
    } finally {
      setIsProcessing(false);
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-md mx-auto p-6 mt-10">
        <Card>
          <CardHeader>
            <Skeleton className="h-8 w-3/4 mb-2" />
            <Skeleton className="h-4 w-full" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-20 w-full mb-4" />
            <Skeleton className="h-10 w-full" />
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-md mx-auto p-6 mt-10">
        <Card className="border-red-200">
          <CardHeader>
            <div className="flex items-center gap-2">
              <BiError className="h-5 w-5 text-red-500" />
              <CardTitle className="text-red-600">Payment Error</CardTitle>
            </div>
            <CardDescription>There was a problem with your payment retry</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-center text-red-500 my-4">{error}</p>
          </CardContent>
          <CardFooter className="flex justify-center">
            <Button
              className="mr-2"
              variant="outline"
              onClick={() => navigate('/user/orders-history')}
            >
              <FiArrowLeft className="h-4 w-4 mr-2" />
              Back to Orders
            </Button>
            <Button onClick={() => navigate('/')}>
              Continue Shopping
            </Button>
          </CardFooter>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto p-6 mt-10">
      <Card>
        <CardHeader>
          <CardTitle>Retry Payment</CardTitle>
          <CardDescription>
            Complete your payment for order #{orderId}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="font-medium">Order Amount:</span>
              <span>${order?.total.toFixed(2)}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Payment Method:</span>
              <span className="capitalize">
                {order?.paymentMethod === 'Unknown' || !order?.paymentMethod ?
                  (paymentDetails.paymentMethod || 'MoMo') : order.paymentMethod}
              </span>
            </div>
          </div>
        </CardContent>
        <CardFooter className="flex justify-between">
          <Button
            variant="outline"
            disabled={isProcessing}
            onClick={() => navigate('/user/orders-history')}
          >
            Cancel
          </Button>
          <Button
            onClick={handleRetryPayment}
            disabled={isProcessing}
            className="bg-[#0496FF] hover:bg-blue-600"
          >
            {isProcessing ? 'Processing...' : 'Pay Now'}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default RetryPayment; 