import React from "react";
import { CartItem } from "../cart/models/cart.model";
import { defaultCartTotals } from "@/context/CartContext";

interface OrderSummaryProps {
  cartItems: CartItem[];
  cartTotals: typeof defaultCartTotals;
  onPlaceOrder: () => void;
}

const OrderSummary: React.FC<OrderSummaryProps> = ({
  cartItems,
  cartTotals,
  onPlaceOrder
}) => {
  return (
    <div className="p-4 ml-5">
      <h2 className="text-lg font-bold mb-3">Order Summary</h2>

      <div className="space-y-3 mb-4">
        {cartItems.map((item) => (
          <div key={item.id} className="flex items-center rounded-sm">
            {item.image && (
              <div className="max-w-[50px] max-h-[50px] rounded-sm flex-shrink-0 mr-3">
                <img
                  src={item.image}
                  alt={item.name}
                  className="rounded-md w-[40px] max-w-[50px] max-h-[50px] object-cover"
                />
              </div>
            )}
            <div className="flex flex-col justify-center flex-1 min-w-0">
              {" "}
              {/* Added min-w-0 to enable truncation */}
              <p className="text-[15px] text-neutral-500 font-medium truncate">
                {item.name}
              </p>
              <div className="flex items-end">
                <span className="text-[13px] text-neutral-500">
                  {item.quantity} x{" "}
                </span>
                <span className="text-[#2DA5F3] font-medium ml-1">
                  ${item.price}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="space-y-2 text-sm pt-5">
        <div className="flex justify-between mb-3">
          <span className="text-[15px] font-medium text-neutral-500">
            Sub-total
          </span>
          <span className="text-[15px] font-bold text-neutral-800">
            ${cartTotals.subtotal}
          </span>
        </div>

        <div className="flex justify-between mb-3 ">
          <span className="text-[15px] font-medium text-neutral-500">
            Shipping
          </span>
          <span className="text-[15px] font-bold text-neutral-800">
            {cartTotals.shipping === "Free"
              ? "Free"
              : `$${cartTotals.shipping}`}
          </span>
        </div>

        {cartTotals.discount > 0 && (
          <div className="flex justify-between mb-3">
            <span className="text-[15px] font-medium text-neutral-500">
              Discount
            </span>
            <span className="text-[15px] font-bold text-neutral-800">
              ${cartTotals.discount}
            </span>
          </div>
        )}

        <div className="flex justify-between mb-3">
          <span className="text-[15px] font-medium text-neutral-500">Tax</span>
          <span className="text-[15px] font-bold text-neutral-800">
            ${cartTotals.tax.toFixed(2)}
          </span>
        </div>
      </div>

      <div className="flex justify-between font-medium text-base mt-4 pt-3 border-t">
        <span className="text-[15px] font-medium text-neutral-500">Total</span>
        <span className="text-[15px] font-bold text-neutral-800">
          ${cartTotals.total.toFixed(2)} USD
        </span>
      </div>

      <button onClick={onPlaceOrder} className="mt-5 w-full bg-orange-500 hover:bg-orange-400 text-white py-3 rounded font-sans font-medium flex items-center justify-center transition-colors duration-200">
        PLACE ORDER <span className="ml-2">→</span>
      </button>
    </div>
  );
};

export default OrderSummary;