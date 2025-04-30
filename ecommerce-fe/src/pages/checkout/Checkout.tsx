// pages/checkout/Checkout.tsx

import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import BillingInformation from "../../components/checkout/BillingInformation";
import PaymentOption from "../../components/checkout/PaymentOption";
import OrderSummary from "../../components/checkout/OrderSummary";
import { useCart } from "@/hooks/useCart";
import { useCheckout } from "@/hooks/useCheckout";
import { BillingInfo, ShippingInfo } from '@/types/checkout.model';
import { useToast } from "@/hooks/use-toast";

const Checkout: React.FC = () => {
  const { cartItems, cartTotals } = useCart();
  const [paymentMethod, setPaymentMethod] = useState<string>("cod");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();
  const { validateCheckout, createOrder, processPayment } = useCheckout();
  const { toast } = useToast();

  // State để lưu thông tin thanh toán từ form
  const [billingInfo, setBillingInfo] = useState<BillingInfo>({
    firstName: "",
    lastName: "",
    address: "",
    country: "VN",
    state: "",
    city: "",
    zipCode: "",
    email: "",
    phone: ""
  });

  const [shippingInfo, setShippingInfo] = useState<ShippingInfo>({
    shipToDifferentAddress: false,
    shippingMethod: "standard"
  });

  // Hàm cập nhật billing info từ form BillingInformation
  const handleBillingInfoChange = (info: Partial<BillingInfo>) => {
    setBillingInfo(prev => ({ ...prev, ...info }));
  };

  // Hàm cập nhật shipping info từ form BillingInformation
  const handleShippingInfoChange = (info: Partial<ShippingInfo>) => {
    setShippingInfo(prev => ({ ...prev, ...info }));
  };

  const handlePlaceOrder = async (): Promise<void> => {
    if (isSubmitting) return;

    try {
      setIsSubmitting(true);

      // Validate thông tin thanh toán
      const validationResult = await validateCheckout(
        billingInfo,
        shippingInfo,
        paymentMethod
      );

      if (!validationResult.valid) {
        // Hiển thị thông báo lỗi
        return;
      }

      // Tạo đơn hàng
      const orderResult = await createOrder(
        billingInfo,
        shippingInfo,
        paymentMethod
      );

      // Xử lý thanh toán
      if (paymentMethod === "momo") {
        // Xử lý thanh toán MoMo
        await processPayment(
          orderResult.id,
          paymentMethod,
          orderResult.totals.total,
          `${window.location.origin}/user/checkout/success?orderId=${orderResult.id}`
        );
      } else if (paymentMethod === "cod") {
        // Thanh toán Cash on Delivery
        toast({
          title: "Order Placed Successfully",
          description: "Your order has been placed. You will pay on delivery.",
          variant: "success"
        });
        navigate(`/user/checkout/success?orderId=${orderResult.id}&method=cash`);
      } else if (paymentMethod === "vnpay") {
        toast({
          title: "Payment Method Unavailable",
          description: "VNPAY payment is not yet available. Please choose another payment method.",
          variant: "destructive"
        });
        setIsSubmitting(false);
        return;
      } else {
        toast({
          title: "Payment Method Error",
          description: "Selected payment method is not supported",
          variant: "destructive"
        });
      }
    } catch (error) {
      toast({
        title: "Error Placing Order",
        description: error instanceof Error ? error.message : "Failed to place order",
        variant: "destructive"
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  if (cartItems.length === 0) {
    return (
      <div className="max-w-[1440px] min-w-[1024px] px-[5rem] py-[1.5rem] mx-auto text-center">
        <h1 className="text-2xl font-bold mb-5">Checkout</h1>
        <div className="max-w-[312px] h-[280px] relative flex justify-center items-center mx-auto">
          <img
            src="/images/system/empty_cart.png"
            alt="empty_cart"
            width={312}
            height={160}
            className="object-contain w-full h-full"
          />
        </div>
        <p className="mb-5 mt-5">Your cart is empty. Please add some products before checkout.</p>
        <button
          onClick={() => navigate('/products')}
          className="bg-[#8AC732] text-white font-sans w-[146px] h-[48px] rounded-[100px] font-medium text-sm hover:bg-[#7AB52B] transition-colors"
        >
          Continue Shopping
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-[1440px] min-w-[1024px] px-[5rem] py-[1.5rem] mx-auto">
      <h1 className="text-2xl font-bold mb-5">Checkout</h1>
      <div className="flex item-center">
        <div className="w-[65%]">
          <BillingInformation
            onBillingInfoChange={handleBillingInfoChange}
            onShippingInfoChange={handleShippingInfoChange}
          />
          <PaymentOption
            paymentMethod={paymentMethod}
            setPaymentMethod={setPaymentMethod}
          />
        </div>
        <div className="w-[35%]">
          <OrderSummary
            onPlaceOrder={handlePlaceOrder}
            cartItems={cartItems}
            cartTotals={cartTotals}
            isSubmitting={isSubmitting}
          />
        </div>
      </div>
    </div>
  );
};

export default Checkout;
