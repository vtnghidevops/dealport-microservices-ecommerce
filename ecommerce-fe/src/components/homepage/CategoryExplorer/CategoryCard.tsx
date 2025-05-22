import React from "react";
import { Category } from "@/types/category.model";
import { Link } from "react-router-dom";
import { generateSlug } from "../../../utils/helpers";

export interface CategoryExplorerCard {
  title?: string;
  category: Category;
  viewAllLabel?: string;
  itemWidth?: string;
  itemHeight?: string;
  className?: string;
  showNavigationArrow?: boolean;
  onViewAllClick?: () => void;
  onItemClick?: (category: Category, index: number) => void;
}

const CategoryCard: React.FC<CategoryExplorerCard> = ({
  category,
  itemWidth,
  itemHeight,
  onItemClick,
}) => {
  // Tạo slug từ tên danh mục (Create slug from category name)
  const categorySlug = category.slug || generateSlug(category.name);
  // console.log(category.name,":", categorySlug);
  return (
    <Link
      to={`/category/${categorySlug}`}
      className={`flex-shrink-0 ${itemWidth} cursor-pointer h-[220px] w-[180px] mr-[3rem] rounded-xl`}
      onClick={(e) => {
        // Nếu muốn giữ lại onItemClick callback đã có (If you want to keep the existing onItemClick callback)
        if (category.id !== undefined) {
          // Ngăn chuyển hướng trang nếu cần xử lý sự kiện click đặc biệt (Prevent page navigation if special click handling is needed)
          e.preventDefault();
          onItemClick?.(category, Number(category.id));
        }
      }}
    >
      <div className="relative rounded-lg overflow-hidden shadow-sm border border-gray-300 h-full">
        <div
          className={`${itemHeight} w-full flex items-center justify-center overflow-hidden`}
        >
          <img
            src={category.imageUrl}
            alt={category.name}
            className="rounded-2xl w-[148px] h-[140px] object-cover transition-transform duration-300 ease-out hover:scale-110"
          />
        </div>
        <div className="p-2 text-center absolute bottom-0 flex justify-center w-full bg-white">
          <h3 className="text-sm font-medium text-gray-800">{category.name}</h3>
        </div>
      </div>
    </Link>
  );
};

export default CategoryCard;