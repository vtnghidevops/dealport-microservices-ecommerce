// components/admin/product/add/components/ProductPricing.tsx
import React, { ChangeEvent, useEffect } from "react";
import { Product } from "@/types/product.model";

interface ProductPricingProps {
  product: Product;
  onChange: (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  discountAmount: string;
}

export const ProductPricing: React.FC<ProductPricingProps> = ({
  product,
  onChange,
  discountAmount,
}) => {
  // Handle numeric input to ensure only numbers and decimal points
  const handleNumericInput = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    // Only allow numbers and a single decimal point
    const regex = /^[0-9]*\.?[0-9]*$/;

    if (value === '' || regex.test(value)) {
      onChange({
        ...e,
        target: {
          ...e.target,
          name,
          value
        }
      });
    }
  };

  // Calculate discounted price based on original price and discount percentage
  const calculateDiscountedPrice = () => {
    if (product.originalPrice && product.discount) {
      const discountedPrice = product.originalPrice - (product.originalPrice * product.discount / 100);
      return discountedPrice.toFixed(2);
    }
    return product.price ? product.price.toString() : '';
  };

  // Handle date changes for expiration dates
  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e);
  };

  // Update price whenever original_price or discount changes
  useEffect(() => {
    if (product.originalPrice && product.discount) {
      const discountedPrice = product.originalPrice - (product.originalPrice * product.discount / 100);

      // Create synthetic event to update price
      const priceEvent = {
        target: {
          name: "price",
          value: discountedPrice.toFixed(2)
        }
      } as ChangeEvent<HTMLInputElement>;

      onChange(priceEvent);
    }
  }, [product.originalPrice, product.discount, onChange]);

  // Handle price change directly
  const handlePriceChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { value } = e.target;
    // Only allow numbers and a single decimal point
    const regex = /^[0-9]*\.?[0-9]*$/;

    if (value === '' || regex.test(value)) {
      onChange({
        ...e,
        target: {
          ...e.target,
          name: "price",
          value
        }
      });
    }
  };

  // Format discount amount to display with 2 decimal places
  const formatSavingsAmount = (amount: string) => {
    if (!amount) return "0.00";
    return parseFloat(amount).toFixed(2);
  };

  return (
    <div className="bg-white rounded-lg p-[1.25rem] mt-5 mb-6">
      <h2 className="block text-cyprus font-bold text-[22px] mb-6">
        Pricing
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Original Price
        </label>
        <div className="flex">
          <div className="flex-grow relative">
            <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
            <input
              type="text"
              name="originalPrice"
              value={product.originalPrice || ""}
              onChange={handleNumericInput}
              placeholder="999.89"
              className="w-full border pl-[25px] text-cyprus bg-neutral-50 border-gray-200 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
            />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-4">
        <div>
          <label className="block text-cyprus font-bold text-[15px] mb-2">
            Discount Percentage <span className="text-neutral-500">(Optional)</span>
          </label>
          <div className="relative">
            <input
              type="text"
              name="discount"
              value={product.discount || ""}
              onChange={handleNumericInput}
              placeholder="10"
              className="w-full border text-cyprus bg-neutral-50 border-gray-200 rounded-lg p-2 pr-8 focus:outline-none focus:border-ocean-green"
            />
            <span className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500">%</span>
          </div>
        </div>

        <div>
          <label className="block text-cyprus font-bold text-[15px] mb-2">
            Sale Price <span className="text-neutral-500">{product.discount > 0 ? "(Calculated)" : ""}</span>
          </label>
          <div className="relative">
            <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
            <input
              type="text"
              name="price"
              value={product.price || ""}
              onChange={handlePriceChange}
              placeholder="899.90"
              className="w-full border pl-[25px] text-cyprus bg-neutral-50 border-gray-200 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
              readOnly={product.discount > 0}
            />
          </div>
          {parseFloat(discountAmount) > 0 && (
            <p className="text-sm text-green-600 mt-1">
              Save: <span className="font-semibold">${formatSavingsAmount(discountAmount)}</span>
            </p>
          )}
        </div>
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Sale Period
        </label>
        <div className="grid grid-cols-2 gap-4">
          <div className="relative">
            <input
              type="date"
              name="expirationStart"
              value={product.uiMetadata?.expirationStart || ""}
              onChange={handleDateChange}
              placeholder="Start Date"
              className="w-full border border-gray-200 rounded-lg p-2 pr-10 text-cyprus bg-neutral-50 focus:outline-none focus:border-ocean-green"
            />
            <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            </div>
          </div>
          <div className="relative">
            <input
              type="date"
              name="expirationEnd"
              value={product.uiMetadata?.expirationEnd || ""}
              onChange={handleDateChange}
              placeholder="End Date"
              className="w-full border border-gray-200 rounded-lg p-2 pr-10 bg-neutral-50 focus:outline-none focus:border-ocean-green"
            />
            <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};