import React from "react";
import { Category } from "@/types/category.model";

interface ProductCategoryListProps {
  categories: Category[];
}

const ProductCategoryList: React.FC<ProductCategoryListProps> = ({ categories = [] }) => {
  return (
    <div className="grid grid-cols-3 gap-3">
      {categories.slice(0, 6).map((category) => (
        <div
          key={category.id}
          className="flex flex-col items-center bg-gray-50 rounded-md p-2 cursor-pointer hover:bg-gray-100 transition-colors duration-200"
        >
          <div className="w-10 h-10 flex items-center justify-center">
            <img
              src={category.imageUrl || '/assets/images/placeholder.png'}
              alt={category.name}
              className="max-w-full max-h-full object-contain"
            />
          </div>
          <span className="mt-2 text-xs text-gray-700 font-medium text-center">
            {category.name}
          </span>
        </div>
      ))}
    </div>
  );
};

export default ProductCategoryList;
