// components/admin/product/add/components/ProductShipping.tsx
import React from "react";
import { ShippingInfo } from "../models/product.model";

interface ProductShippingProps {
  shipping: ShippingInfo;
  onChange: (field: string, value: string) => void;
}

export const ProductShipping: React.FC<ProductShippingProps> = ({
  shipping,
  onChange,
}) => {
  return (
    <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        Shipping Information
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Courier Shipping
        </label>
        <input
          type="text"
          value={shipping.courier || ""}
          onChange={(e) => onChange("courier", e.target.value)}
          placeholder="e.g. 2-3 Business Days"
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        />
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Local Shipping
        </label>
        <input
          type="text"
          value={shipping.local || ""}
          onChange={(e) => onChange("local", e.target.value)}
          placeholder="e.g. 1-2 Business Days"
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        />
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          UPS Shipping
        </label>
        <input
          type="text"
          value={shipping.ups || ""}
          onChange={(e) => onChange("ups", e.target.value)}
          placeholder="e.g. 3-5 Business Days"
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        />
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Global Shipping
        </label>
        <input
          type="text"
          value={shipping.global || ""}
          onChange={(e) => onChange("global", e.target.value)}
          placeholder="e.g. 7-14 Business Days"
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        />
      </div>
    </div>
  );
};