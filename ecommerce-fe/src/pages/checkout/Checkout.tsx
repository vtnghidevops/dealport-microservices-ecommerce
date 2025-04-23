// pages/checkout/Checkout.tsx

import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import BillingInformation from "../../components/checkout/BillingInformation";
import PaymentOption from "../../components/checkout/PaymentOption";
import OrderSummary from "../../components/checkout/OrderSummary";
import { useCart } from "../../hooks/useCart";

const Checkout: React.FC = () => {
  const { cartItems, cartTotals } = useCart();
  const [paymentMethod, setPaymentMethod] = useState<string>("cash");
  const navigate = useNavigate();
  
  const handlePlaceOrder = async (): Promise<void> => {
    try {
      // Here would typically be the API call to process the order
      if (paymentMethod === "momo" || paymentMethod === "banking") {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000));
        // Redirect to success page with payment method
        navigate(`/user/checkout/success?method=${paymentMethod}`);
      } else {
        // For cash on delivery
        await new Promise(resolve => setTimeout(resolve, 500));
        navigate(`/user/checkout/success?method=cash`);
      }
    } catch (error) {
      console.error("Error placing order:", error);
    }
  }

    return (
      <div className="max-w-[1440px] min-w-[1024px] px-[5rem] py-[1.5rem] mx-auto ">
        <h1 className="text-2xl font-bold mb-5">Checkout</h1>
        <div className="flex item-center">
          <div className="w-[65%]">
            <BillingInformation />
            <PaymentOption
              paymentMethod={paymentMethod}
              setPaymentMethod={setPaymentMethod}
            />
          </div>
          <div className="w-[35%]">
            <OrderSummary onPlaceOrder={handlePlaceOrder} cartItems={cartItems} cartTotals={cartTotals} />
          </div>
        </div>
      </div>
    );
  };


export default Checkout;
