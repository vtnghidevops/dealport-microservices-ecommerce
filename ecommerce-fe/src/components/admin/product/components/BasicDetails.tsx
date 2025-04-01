// components/admin/product/add/components/BasicDetails.tsx
import React, { ChangeEvent } from "react";
import { Product } from "../models/product.model";

interface BasicDetailsProps {
  product: Product;
  onChange: (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
}

export const BasicDetails: React.FC<BasicDetailsProps> = ({
  product,
  onChange,
}) => {
  return (
    <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        Basic Details
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Product Name
        </label>
        <input
          type="text"
          name="name"
          value={product.name}
          onChange={onChange}
          placeholder="iPhone 15"
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        />
      </div>

      <div className="mb-4 mt-12">
        <label className="block text-cyprus font-bold text-[15px] mb-12">
          Product Description
        </label>
        <textarea
          name="description"
          value={product.description}
          onChange={onChange}
          placeholder="The iPhone 15 delivers cutting-edge performance with the A16 Bionic chip, an immersive Super Retina XDR display, advanced dual-camera system, and exceptional battery life, all encased in stunning aerospace-grade aluminum."
          rows={5}
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        ></textarea>
        <div className="flex justify-end mt-2 space-x-2 border-b border-neutral-300">
          <button className="p-2 hover:bg-gray-100 rounded">
            <svg
              className="w-5 h-5 text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
              />
            </svg>
          </button>
          <button className="p-2 hover:bg-gray-100 rounded">
            <svg
              className="w-5 h-5 text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
};