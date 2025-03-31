// components/admin/product/add/components/ProductPricing.tsx
import React, { ChangeEvent } from "react";
import { Product } from "../../models/product.model";

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
  return (
    <div className="bg-white rounded-lg p-[1.25rem] mt-5 mb-6">
      <h2 className="block text-cyprus font-bold text-[22px] mb-12 mt-[1rem]">
        Pricing
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Product Price
        </label>
        <div className="flex">
          <div className="flex-grow">
            <input
              type="text"
              name="price"
              value={product.price || ""}
              onChange={onChange}
              placeholder="$999.89"
              className="w-full border text-cyprus bg-neutral-50 border-gray-200 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
            />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-4">
        <div>
          <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
            Discounted Price <span className="text-neutral-500">(Optional)</span>
          </label>
          <div className="border border-gray-200 bg-neutral-50 rounded-lg flex items-center w-[272px] gap-4">
            <span className="ml-2 flex items-center justify-center bg-aqua-spring w-[32px] h-[32px] text-black rounded-lg mr-2">
              $
            </span>
            <input
              type="text"
              name="discountedPrice"
              value={product.discountedPrice || ""}
              onChange={onChange}
              placeholder="$99"
              className="w-[108px] flex-grow py-[10px] text-cyprus bg-neutral-50 px-12 focus:outline-none focus:border-ocean-green"
            />
            <label className="block text-[15px] font-bold text-cyprus mb-1 w-[103px] h-[18px]">
              Sale = <span className="text-cyprus">${discountAmount}</span>
            </label>
          </div>
        </div>
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Expiration
        </label>
        <div className="grid grid-cols-2 gap-4">
          <div className="relative">
            <input
              type="text"
              onFocus={(e) => (e.target.type = "date")}
              onBlur={(e) => {
                if (!e.target.value) e.target.type = "text";
              }}
              name="expirationStart"
              value={product.expirationStart}
              onChange={onChange}
              placeholder="Start"
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
              type="text"
              onFocus={(e) => (e.target.type = "date")}
              onBlur={(e) => {
                if (!e.target.value) e.target.type = "text";
              }}
              name="expirationEnd"
              value={product.expirationEnd}
              onChange={onChange}
              placeholder="End"
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