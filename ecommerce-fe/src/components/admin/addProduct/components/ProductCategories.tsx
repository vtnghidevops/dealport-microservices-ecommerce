// components/admin/product/add/components/ProductCategories.tsx
import React from "react";
import { Product } from "../models/product.model";

interface ProductCategoriesProps {
  product: Product;
  onColorSelect: (color: string) => void;
}

export const ProductCategories: React.FC<ProductCategoriesProps> = ({
  product,
  onColorSelect,
}) => {
  return (
    <div className="bg-white rounded-lg p-[1rem] shadow-sm mt-5">
      <h2 className="block text-cyprus font-bold text-[22px] mb-12 mt-[1rem]">Categories</h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Product Categories
        </label>
        <div className="relative">
          <select className="w-full border border-gray-200 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option>Select your product</option>
          </select>
          <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
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
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </div>
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Product Tag
        </label>
        <div className="relative">
          <select className="w-full border border-gray-200 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option>Select your product</option>
          </select>
          <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
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
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </div>
      </div>

      <div>
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Select your color
        </label>
        <div className="flex space-x-2">
          {["#D9F3D9", "#FFD8D8", "#D8E7F3", "#FFF8D8", "#2F3133"].map(
            (color) => (
              <button
                key={color}
                className={`w-[48px] h-[48px] rounded-md ${
                  product.color === color
                    ? "drop-shadow-xl filter"
                    : ""
                }`}
                style={{ backgroundColor: color }}
                onClick={() => onColorSelect(color)}
              ></button>
            )
          )}
        </div>
      </div>
    </div>
  );
};