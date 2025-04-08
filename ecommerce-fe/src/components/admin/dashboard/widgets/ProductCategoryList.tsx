import React from "react";
import { ProductCategory } from "../models/product.model";
import { IoChevronForward } from "react-icons/io5";

interface ProductCategoryListProps {
  categories: ProductCategory[];
}

const ProductCategoryList: React.FC<ProductCategoryListProps> = ({
  categories,
}) => {
  return (
    <div className="w-full">
      <h2 className="text-neutral-500 text-[14px] font-medium mb-[10px]">
        Categories
      </h2>

      <div className="space-y-3">
        {categories.map((category) => (
          <div
            key={category.id}
            className="w-[321px] h-[58px] flex items-center justify-between p-6 bg-white rounded-lg filter drop-shadow-md cursor-pointer hover:shadow-md transition-shadow "
          >
            <div className="flex items-center gap-3">
              <div className="w-[46px] h-[46px] bg-gray-50 rounded-md flex items-center justify-center overflow-hidden">
                <img
                  src={category.image_url}
                  alt={category.name}
                  className="w-full h-full object-contain"
                />
              </div>
              <span className="font-medium">{category.name}</span>
            </div>
            <IoChevronForward className="text-gray-400" />
          </div>
        ))}
      </div>

      <div className="mt-[15px] flex justify-center">
        <button className="text-primary font-medium text-[14px]">
          See more
        </button>
      </div>
    </div>
  );
};

export default ProductCategoryList;
