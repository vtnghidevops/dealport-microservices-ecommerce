import React from "react";
import { CartItem } from "@/types/cart.model";
import { defaultCartTotals } from "@/context/CartContext";

interface OrderSummaryProps {
  cartItems: CartItem[];
  cartTotals: typeof defaultCartTotals;
  onPlaceOrder: () => void;
  isSubmitting?: boolean;
}

const OrderSummary: React.FC<OrderSummaryProps> = ({
  cartItems,
  cartTotals,
  onPlaceOrder,
  isSubmitting = false
}) => {
  // Store the total in sessionStorage before placing order
  const handlePlaceOrder = () => {
    // Store the current cart total in sessionStorage
    sessionStorage.setItem('orderTotal', cartTotals.total.toString());
    // Call the original onPlaceOrder function
    onPlaceOrder();
  };

  return (
    <div className="p-4 ml-5">
      <h2 className="text-lg font-bold mb-3">Order Summary</h2>

      <div className="space-y-3 mb-4">
        {cartItems.map((item) => (
          <div key={item.id} className="flex items-center rounded-sm">
            {item.imageUrl && (
              <div className="max-w-[50px] max-h-[50px] rounded-sm flex-shrink-0 mr-3">
                <img
                  src={item.imageUrl}
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
            ${cartTotals.subtotal.toFixed(0)}
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
              ${cartTotals.discount.toFixed(0)}
            </span>
          </div>
        )}

        <div className="flex justify-between mb-3">
          <span className="text-[15px] font-medium text-neutral-500">Tax</span>
          <span className="text-[15px] font-bold text-neutral-800">
            ${cartTotals.tax.toFixed(0)}
          </span>
        </div>
      </div>

      <div className="flex justify-between font-medium text-base mt-4 pt-3 border-t">
        <span className="text-[15px] font-medium text-neutral-500">Total</span>
        <span className="text-[15px] font-bold text-neutral-800">
          ${cartTotals.total.toFixed(0)} USD
        </span>
      </div>

      <button
        onClick={handlePlaceOrder}
        disabled={isSubmitting}
        className={`mt-5 w-full py-3 rounded font-sans font-medium flex items-center justify-center transition-colors duration-200
          ${isSubmitting
            ? 'bg-gray-400 text-white cursor-not-allowed'
            : 'bg-orange-500 hover:bg-orange-400 text-white'
          }`}
      >
        {isSubmitting ? (
          <>
            <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            PROCESSING...
          </>
        ) : (
          <>
            PLACE ORDER <span className="ml-2">→</span>
          </>
        )}
      </button>
    </div>
  );
};

export default OrderSummary;