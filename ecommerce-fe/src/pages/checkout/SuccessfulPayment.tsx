import React, { useEffect, useState, useRef } from "react";
import { useNavigate, useSearchParams, useLocation } from "react-router-dom";
import { BsCheckCircleFill } from "react-icons/bs";
import { useCart } from "../../hooks/useCart";
import { useToast } from "@/hooks/use-toast";

const SuccessfulPayment: React.FC = () => {
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { cartItems, cartTotals, clearCart } = useCart();
  const { toast } = useToast();
  const [animateCheck, setAnimateCheck] = useState(false);
  const [orderTotal, setOrderTotal] = useState<number>(0);
  // Add ref to track if cart has been cleared
  const hasCartBeenCleared = useRef(false);

  // Get values either from URL parameters or from location state
  const locationState = location.state || {};

  // Get order ID from URL parameters or state (if provided by payment processor)
  const orderId = searchParams.get("orderId") || locationState.orderId || "ORD" + Math.floor(Math.random() * 100000);
  const paymentMethod = searchParams.get("method") || locationState.method || "online";
  const transactionId = searchParams.get("transactionId") || locationState.transactionId || "";
  // Use correct status according to schema
  const status = 'paid';

  useEffect(() => {
    // Get the stored total from sessionStorage if it exists
    const storedTotal = sessionStorage.getItem('orderTotal');
    if (storedTotal) {
      setOrderTotal(parseFloat(storedTotal));
    } else if (cartTotals.total > 0) {
      // If cartTotals.total is valid, use it and store for future reference
      setOrderTotal(cartTotals.total);
      sessionStorage.setItem('orderTotal', cartTotals.total.toString());
    }

    // Trigger animation after a short delay for better visual effect
    setTimeout(() => {
      setAnimateCheck(true);
    }, 300);

    // Clear cart after successful payment - but only once
    const handlePaymentSuccess = async () => {
      // Check if cart has already been cleared in this session
      if (hasCartBeenCleared.current) return;

      try {
        hasCartBeenCleared.current = true; // Mark as cleared before the API call
        await clearCart();
        toast({
          title: "Payment Successful",
          description: "Your order has been placed successfully",
          variant: "success"
        });

        // Store payment status in sessionStorage
        sessionStorage.setItem('paymentStatus', status);
      } catch (error) {
        console.error("Error clearing cart:", error);
        // If there's an error, we might want to reset the flag to try again
        hasCartBeenCleared.current = false;
      }
    };

    handlePaymentSuccess();
    // Remove cartTotals.total from dependencies to avoid infinite loop
  }, [clearCart, toast, status]);

  const handleViewOrders = () => {
    navigate("/user/orders-history");
  };

  // Use the orderTotal state instead of directly using cartTotals.total
  return (
    <div className="max-w-[1440px] min-w-[1024px] px-[5rem] py-[3rem] mx-auto">
      <div className="max-w-2xl mx-auto bg-white p-8 rounded-lg shadow-sm">
        <div className="text-center mb-5">
          <div className="mb-[2rem] relative flex justify-center items-center h-24">
            {/* Ripple effect circle */}
            <div className={`absolute rounded-full bg-green-100 transition-all duration-1000 ${animateCheck ? 'w-24 h-24 opacity-0' : 'w-0 h-0 opacity-100'}`}></div>

            {/* Second ripple for multi-stage effect */}
            <div className={`absolute rounded-full bg-green-200 transition-all duration-700 ${animateCheck ? 'w-20 h-20 opacity-0' : 'w-0 h-0 opacity-100'}`}></div>

            {/* Green circle that grows */}
            <div className={`absolute rounded-full bg-green-50 transition-all duration-500 ${animateCheck ? 'w-16 h-16 scale-100' : 'w-0 h-0 scale-0'}`}></div>

            {/* The check icon with scale and rotation animation */}
            <BsCheckCircleFill
              className={`text-success text-6xl relative z-10 transition-all duration-500 ${animateCheck
                ? 'opacity-100 scale-100 rotate-0'
                : 'opacity-0 scale-0 -rotate-90'
                }`}
            />
          </div>
          <h1 className="text-2xl font-bold text-gray-800 mb-2">
            Payment Successful!
          </h1>
          <p className="text-gray-600">
            Your order has been confirmed and is being processed.
          </p>
        </div>

        <div className="py-5 mb-5 border-t border-b border-gray-200">
          <div className="flex justify-between mb-2">
            <span className="text-gray-600">Order ID:</span>
            <span className="font-medium">{orderId}</span>
          </div>
          <div className="flex justify-between mb-2">
            <span className="text-gray-600">Payment Method:</span>
            <span className="font-medium capitalize">{paymentMethod}</span>
          </div>
          <div className="flex justify-between mb-2">
            <span className="text-gray-600">Status:</span>
            <span className="font-medium text-green-600 capitalize">{status}</span>
          </div>
          {transactionId && (
            <div className="flex justify-between mb-2">
              <span className="text-gray-600">Transaction ID:</span>
              <span className="font-medium">{transactionId}</span>
            </div>
          )}
          <div className="flex justify-between">
            <span className="text-gray-600">Total Amount:</span>
            <span className="font-bold">${orderTotal.toFixed(0)} USD</span>
          </div>
        </div>

        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-3">Order Details</h2>
          <div className="space-y-3">
            {cartItems.map((item) => (
              <div key={item.id} className="flex items-center">
                {item.imageUrl && (
                  <div className="mr-3">
                    <img
                      src={item.imageUrl}
                      alt={item.name}
                      className="w-[40px] h-[40px] object-cover rounded-md"
                    />
                  </div>
                )}
                <div className="flex-1">
                  <p className="font-medium">{item.name}</p>
                  <div className="text-sm text-gray-500">
                    {item.quantity} x ${item.price}
                  </div>
                </div>
                <div className="font-medium">
                  ${(item.price * item.quantity).toFixed(2)}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="text-center mt-5">
          <p className="text-gray-800 mb-4">
            Please check your email for detailed information.
          </p>
          <div className="flex gap-16 justify-center mt-5">
            <button
              onClick={() => navigate("/")}
              className="w-[170px] h-[50px] px-6 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
            >
              Continue Shopping
            </button>
            <button
              onClick={handleViewOrders}
              className="w-[170px] h-[50px] px-6 py-2 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
            >
              View Orders
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SuccessfulPayment;