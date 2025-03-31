// components/admin/product/add/components/ProductInventory.tsx
import React, { ChangeEvent } from "react";
import { Product } from "../../models/product.model";

interface ProductInventoryProps {
  product: Product;
  onChange: (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => void;
  onRadioChange: (e: ChangeEvent<HTMLInputElement>) => void;
  handlePublish: () => void;
  handleSaveDraft: () => void;
}

export const ProductInventory: React.FC<ProductInventoryProps> = ({
  product,
  onChange,
  onRadioChange,
  handlePublish,
  handleSaveDraft
}) => {
  return (
    <div className="bg-white p-[1.25rem] rounded-lg  mb-6">
      <h2 className="block text-cyprus font-bold text-[22px] mb-12 mt-[1rem]">
        Inventory
      </h2>
      <div className="flex items-center justify-between">
        <div className="mb-4 w-[272px]">
          <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
            Stock Quantity
          </label>
          <input
            type="text"
            name="stockQuantity"
            value={product.stockQuantity}
            onChange={onChange}
            placeholder="Unlimited"
            className="mb-4 w-full border border-gray-200 rounded-lg p-2 bg-neutral-50 focus:outline-none focus:border-ocean-green"
          />
          <div className="flex items-center mt-2">
            <label className="flex items-center cursor-pointer">
              <div
                className={`mr-2 relative inline-block w-[48px] h-[24px] transition-all duration-200 ease-in-out rounded-full ${
                  product.stockQuantity === "Unlimited"
                    ? "bg-green-500"
                    : "bg-gray-300"
                }`}
              >
                <input
                  type="checkbox"
                  className="opacity-0 absolute w-full h-full"
                  checked={product.stockQuantity === "Unlimited"}
                  onChange={(e) =>
                    onChange({
                      ...e,
                      target: {
                        ...e.target,
                        name: "stockQuantity",
                        value: e.target.checked ? "Unlimited" : "0",
                      },
                    } as ChangeEvent<HTMLInputElement>)
                  }
                />
                <span
                  className={`absolute left-4 top-1/2 -translate-y-1/2 bg-white w-[15px] h-[15px] rounded-full transition-transform duration-200 ease-in-out transform ${
                    product.stockQuantity === "Unlimited"
                      ? "translate-x-5"
                      : "translate-x-0"
                  }`}
                ></span>
              </div>
              <span className="text-sm text-gray-700">Unlimited</span>
            </label>
          </div>
        </div>

        <div className="w-[272px] mb-[2.1rem]">
          <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
            Stock Status
          </label>
          <div className="relative">
            <select
              name="stockStatus"
              value={product.stockStatus}
              onChange={(e) => onChange(e)}
              className="text-cyprus w-full border border-gray-200 bg-neutral-50 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:border-ocean-green"
            >
              <option value="In Stock">In Stock</option>
              <option value="Out of Stock">Out of Stock</option>
              <option value="Pre-order">Pre-order</option>
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
      </div>

      <div className="mb-4 mt-[0.5rem]">
        <label className="flex items-center cursor-pointer">
          <input
            type="checkbox"
            name="highlighted"
            checked={product.highlighted}
            onChange={(e) =>
              onChange({
                ...e,
                target: { ...e.target, name: "highlighted", type: "checkbox" },
              } as ChangeEvent<HTMLInputElement>)
            }
            className="h-[20px] w-[20px] text-ocean-green focus:ring-green-500 border-gray-300 rounded"
          />
          <span className="ml-2 text-[15px] text-neutral-500">
            Highlight this product in a featured section.
          </span>
        </label>
      </div>

      <div className="flex justify-end mt-6 space-x-3">
            <button
              onClick={handleSaveDraft}
              className="px-5 py-2 border border-gray-200 rounded-lg flex items-center hover:bg-gray-50"
            >
              <svg
                className="w-5 h-5 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
                />
              </svg>
              Save to draft
            </button>
            <button
              onClick={handlePublish}
              className="px-5 py-2 bg-ocean-green hover:bg-green-600 text-white rounded-lg"
            >
              Publish Product
            </button>
          </div>
    </div>
  );
};