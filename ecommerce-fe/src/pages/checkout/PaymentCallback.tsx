import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import paymentService from '@/services/user/payment.service';

/**
 * This component handles payment callbacks from external payment providers like MoMo
 * It verifies the payment status and redirects the user accordingly
 */
const PaymentCallback: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [verifying, setVerifying] = useState(true);
  const [message, setMessage] = useState('Verifying payment status...');

  useEffect(() => {
    const verifyPayment = async () => {
      try {
        setVerifying(true);

        // Get all URL parameters
        const searchParams = new URLSearchParams(location.search);
        const paramMap: Record<string, string> = {};

        // Convert URLSearchParams to a plain object
        searchParams.forEach((value, key) => {
          paramMap[key] = value;
        });

        // Check payment provider
        if (paramMap.partnerCode || paramMap.orderId) {
          // This appears to be a MoMo callback
          // console.log('Detected MoMo payment callback');

          // Get orderId from the params
          const { orderId } = paramMap;
          if (!orderId) {
            throw new Error('Order ID not found in callback parameters');
          }

          // Verify the payment with our backend
          setMessage('Verifying payment with MoMo...');
          const result = await paymentService.verifyMomoPayment(paramMap);

          if (result.success) {
            // Payment was successful
            toast({
              title: 'Payment successful',
              description: 'Your payment has been processed successfully.',
              variant: 'success',
            });

            // Redirect to success page with order details
            navigate(`/user/checkout/success?orderId=${result.orderId}&method=momo`);
          } else {
            // Payment failed
            toast({
              title: 'Payment verification failed',
              description: result.message || 'Unable to verify payment.',
              variant: 'destructive',
            });

            // Redirect to payment failure page
            navigate(`/user/checkout/failed?orderId=${orderId}&message=${encodeURIComponent(result.message || 'Verification failed')}`);
          }
        } else {
          // Unknown payment callback format
          throw new Error('Unrecognized payment callback format');
        }
      } catch (error) {
        console.error('Error verifying payment:', error);

        setMessage('Payment verification failed');
        toast({
          title: 'Error verifying payment',
          description: error instanceof Error ? error.message : 'Unknown error occurred',
          variant: 'destructive',
        });

        // Redirect to payment failure page with generic error message
        navigate('/user/checkout/failed?message=verification_failed');
      } finally {
        setVerifying(false);
      }
    };

    verifyPayment();
  }, [location.search, navigate, toast]);

  return (
    <div className="flex flex-col items-center justify-center min-h-[50vh] p-6">
      <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-8 text-center">
        <h1 className="text-2xl font-bold mb-4">Payment Verification</h1>
        {verifying ? (
          <>
            <Loader2 className="animate-spin h-8 w-8 mb-4 mx-auto text-blue-500" />
            <p className="text-gray-600">{message}</p>
            <p className="mt-2 text-sm text-gray-500">Please do not close this page...</p>
          </>
        ) : (
          <p className="text-gray-600">{message}</p>
        )}
      </div>
    </div>
  );
};

export default PaymentCallback; 