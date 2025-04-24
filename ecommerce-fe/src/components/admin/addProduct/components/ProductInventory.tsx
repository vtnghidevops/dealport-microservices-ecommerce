// components/admin/product/add/components/ProductInventory.tsx
import React, { ChangeEvent } from "react";
import { Product } from "../models/product.model";

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
  // Handle quantity input to ensure only numbers
  const handleQuantityChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { value } = e.target;
    // Only allow numbers for stock quantity
    if (value === '' || /^[0-9]+$/.test(value)) {
      onChange({
        ...e,
        target: {
          ...e.target,
          name: "stock_quantity",
          value
        }
      } as ChangeEvent<HTMLInputElement>);
    }
  };

  const isUnlimited = product.stock_quantity === 999999;

  return (
    <div className="bg-white p-[1.25rem] rounded-lg mb-5">
      <h2 className="block text-cyprus font-bold text-[22px] mb-6">
        Inventory
      </h2>
      <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-6">
        <div className="mb-4 w-full md:w-[272px]">
          <label className="block text-cyprus font-bold text-[15px] mb-2">
            Stock Quantity
          </label>
          <input
            type="text"
            name="stock_quantity"
            value={isUnlimited ? "" : product.stock_quantity}
            onChange={handleQuantityChange}
            placeholder="Enter quantity"
            disabled={isUnlimited}
            className={`mb-4 w-full border border-gray-200 rounded-lg p-2 bg-neutral-50 focus:outline-none focus:border-ocean-green ${isUnlimited ? 'opacity-50' : ''}`}
          />
          <div className="flex items-center mt-2">
            <label className="flex items-center cursor-pointer">
              <div
                className={`mr-2 relative inline-block w-[48px] h-[24px] transition-all duration-200 ease-in-out rounded-full ${isUnlimited ? "bg-green-500" : "bg-gray-300"
                  }`}
              >
                <input
                  type="checkbox"
                  className="opacity-0 absolute w-full h-full"
                  checked={isUnlimited}
                  onChange={(e) => {
                    onChange({
                      ...e,
                      target: {
                        ...e.target,
                        name: "stock_quantity",
                        value: e.target.checked ? "999999" : "0",
                      },
                    } as ChangeEvent<HTMLInputElement>);
                  }}
                />
                <span
                  className={`absolute left-1 top-1/2 -translate-y-1/2 bg-white w-[18px] h-[18px] rounded-full transition-transform duration-200 ease-in-out transform ${isUnlimited ? "translate-x-full" : "translate-x-0"
                    }`}
                ></span>
              </div>
              <span className="text-sm text-gray-700">Unlimited</span>
            </label>
          </div>
        </div>

        <div className="w-full md:w-[272px] mb-4">
          <label className="block text-cyprus font-bold text-[15px] mb-2">
            Stock Status
          </label>
          <div className="relative">
            <select
              name="stock_status"
              value={product.stock_status}
              onChange={onChange}
              className="text-cyprus w-full border border-gray-200 bg-neutral-50 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:border-ocean-green"
            >
              <option value="In Stock">In Stock</option>
              <option value="Out of Stock">Out of Stock</option>
              <option value="Pre-order">Pre-order</option>
              <option value="Back Order">Back Order</option>
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

      <div className="mb-4 mt-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          SKU (Stock Keeping Unit)
        </label>
        <input
          type="text"
          name="sku"
          value={product.sku || ''}
          onChange={onChange}
          placeholder="e.g. PRD-12345"
          className="w-full border border-gray-200 rounded-lg p-2 bg-neutral-50 focus:outline-none focus:border-ocean-green"
        />
      </div>

      <div className="mb-4 mt-4">
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


    </div>
  );
};