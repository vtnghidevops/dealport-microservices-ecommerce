// pages/checkout/Checkout.tsx

import React, { useState } from "react";
import BillingInformation from "../../components/checkout/BillingInformation";
import PaymentOption from "../../components/checkout/PaymentOption";
import OrderSummary from "../../components/checkout/OrderSummary";
import { useCart } from "../../hooks/useCart";

const Checkout: React.FC = () => {
  const { cartItems, cartTotals } = useCart();
  const [paymentMethod, setPaymentMethod] = useState<string>("cash");

  const handlePlaceOrder = async () => {
    try {
      if (paymentMethod === "momo" || paymentMethod === "banking") {
        // Giả lập API call
        const response = await fetch("/api/payment", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ method: paymentMethod, amount: cartTotals }),
        });
        const data = await response.json();
        alert(`Payment successful: ${data.message}`);
      } else {
        alert("Order placed successfully with Cash on Delivery!");
      }
    } catch (error) {
      console.error("Error placing order:", error);
    }
  };

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
          <OrderSummary cartItems={cartItems} cartTotals={cartTotals} />
        </div>
      </div>
    </div>
  );
};

export default Checkout;
