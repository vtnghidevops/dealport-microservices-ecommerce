// components/admin/product/add/components/ProductHeader.tsx
import React from "react";
import { CiCirclePlus } from "react-icons/ci";

interface ProductHeaderProps {
  onPublish: () => void;
  onSaveDraft: () => void;
}

export const ProductHeader: React.FC<ProductHeaderProps> = ({
  onPublish,
  onSaveDraft,
}) => {
  return (
    <div className="flex items-center justify-between mb-6">
      <h1 className="text-[18px] font-bold text-cyprus">Add New Product</h1>
      <div className="flex items-center space-x-2">
        <div className="mr-[0.5rem] bg-white border rounded-lg border-gray-200 relative w-[315p] h-[48px] px-12 py-6">
          <input
            type="text"
            placeholder="Search product for add"
            className="pl-10 pr-4 py-2 text-neutral-600 w-72 placeholder:text-neutral-600 focus:outline-none"
          />
          <svg
            className="w-5 h-5 absolute top-[54%] left-3 -translate-y-1/2 text-neutral-600"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fillRule="evenodd"
              d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
              clipRule="evenodd"
            />
          </svg>
        </div>
        <button
          onClick={onPublish}
          className="py-6 px-12 bg-ocean-green border rounded-lg border-gray-200 text-white relative w-[315p] h-[48px] hover:bg-green-600 flex items-center"
        >
          Publish Product
        </button>
        <button
          onClick={onSaveDraft}
          className="py-6 px-12 bg-white border-neutral-300 rounded-lg border text-cyprus relative w-[315p] h-[48px] flex items-center hover:bg-gray-50"
        >
          <svg
            className="w-5 h-5 mr-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
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
        <button className="flex items-center justify-center text-neutral-500 border border-neutral-300 bg-white rounded-lg w-[48px] h-[48px] hover:bg-gray-100 p-2">
          <CiCirclePlus className="p-4 w-full h-full" />
        </button>
      </div>
    </div>
  );
};