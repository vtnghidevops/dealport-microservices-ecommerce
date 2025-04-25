// // components/admin/product/add/components/ProductCategories.tsx
// import React from "react";
// import { Product } from "../models/product.model";

// interface ProductCategoriesProps {
//   product: Product;
//   onColorSelect: (color: string) => void;
// }

// export const ProductCategories: React.FC<ProductCategoriesProps> = ({
//   product,
//   onColorSelect,
// }) => {
//   return (
//     <div className="bg-white rounded-lg p-[1rem] shadow-sm mt-5">
//       <h2 className="block text-cyprus font-bold text-[22px] mb-12 mt-[1rem]">Categories</h2>

//       <div className="mb-4">
//         <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//           Product Categories
//         </label>
//         <div className="relative">
//           <select className="w-full border border-gray-200 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500">
//             <option>Select your product</option>
//           </select>
//           <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
//             <svg
//               className="w-5 h-5"
//               fill="none"
//               stroke="currentColor"
//               viewBox="0 0 24 24"
//             >
//               <path
//                 strokeLinecap="round"
//                 strokeLinejoin="round"
//                 strokeWidth={2}
//                 d="M19 9l-7 7-7-7"
//               />
//             </svg>
//           </div>
//         </div>
//       </div>

//       <div className="mb-4">
//         <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//           Product Tag
//         </label>
//         <div className="relative">
//           <select className="w-full border border-gray-200 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500">
//             <option>Select your product</option>
//           </select>
//           <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
//             <svg
//               className="w-5 h-5"
//               fill="none"
//               stroke="currentColor"
//               viewBox="0 0 24 24"
//             >
//               <path
//                 strokeLinecap="round"
//                 strokeLinejoin="round"
//                 strokeWidth={2}
//                 d="M19 9l-7 7-7-7"
//               />
//             </svg>
//           </div>
//         </div>
//       </div>

//       <div>
//         <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//           Select your color
//         </label>
//         <div className="flex space-x-2">
//           {["#D9F3D9", "#FFD8D8", "#D8E7F3", "#FFF8D8", "#2F3133"].map(
//             (color) => (
//               <button
//                 key={color}
//                 className={`w-[48px] h-[48px] rounded-md ${
//                   product.color === color
//                     ? "drop-shadow-xl filter"
//                     : ""
//                 }`}
//                 style={{ backgroundColor: color }}
//                 onClick={() => onColorSelect(color)}
//               ></button>
//             )
//           )}
//         </div>
//       </div>
//     </div>
//   );
// };
// components/admin/product/add/components/ProductCategories.tsx
import React, { useEffect } from "react";
import { Product } from "@/types/product.model";

interface ProductCategoriesProps {
  product: Product;
  categories: any[];
  onChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  onColorSelect: (color: string) => void;
}

export const ProductCategories: React.FC<ProductCategoriesProps> = ({
  product,
  categories,
  onChange,
  onColorSelect,
}) => {
  // Log available categories on mount
  useEffect(() => {
   //  console.log("ProductCategories - Available categories:", categories);
    // console.log("ProductCategories - Current product category:", product.categoryId, product.categorySlug);
  }, [categories, product.categoryId, product.categorySlug]);

  // Handle category change with additional logging
  const handleCategoryChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const categoryId = parseInt(e.target.value);
    const selectedCategory = categories.find(cat => cat.id === categoryId);

    // console.log("ProductCategories - Selected category:", categoryId, selectedCategory);
    onChange(e);
  };

  return (
    <div className="bg-white rounded-lg p-[1.25rem] shadow-sm mb-6 mt-5">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        Categories
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Product Category
        </label>
        <div className="relative">
          <select
            name="categoryId"
            value={product.categoryId || ""}
            onChange={handleCategoryChange}
            className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:border-ocean-green"
          >
            <option value="">Select a category</option>
            {categories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
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
        {product.categoryId ? (
          <p className="text-sm text-green-600 mt-1">
            Selected category: {categories.find(c => c.id === parseInt(product.categoryId))?.name || "Unknown"} (ID: {product.categoryId})
          </p>
        ) : (
          <p className="text-sm text-red-600 mt-1">
            Please select a category
          </p>
        )}
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Product Type
        </label>
        <div className="relative">
          <select
            name="type"
            value={product.type}
            onChange={onChange}
            className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:border-ocean-green"
          >
            <option value="normal">Normal</option>
            <option value="trending">Trending</option>
            <option value="top-sale">Top Sale</option>
            <option value="new">New</option>
            <option value="limited">Limited</option>
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
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Select your color
        </label>
        <div className="flex space-x-2">
          {["#D9F3D9", "#FFD8D8", "#D8E7F3", "#FFF8D8", "#2F3133"].map(
            (color) => (
              <button
                key={color}
                className={`w-[48px] h-[48px] rounded-md ${product.uiMetadata?.color === color
                  ? "ring-2 ring-offset-2 ring-ocean-green"
                  : ""
                  }`}
                style={{ backgroundColor: color }}
                onClick={() => onColorSelect(color)}
                type="button"
              ></button>
            )
          )}
        </div>
      </div>
    </div>
  );
};