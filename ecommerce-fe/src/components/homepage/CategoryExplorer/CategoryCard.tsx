import React from "react";
import { CategoryExplorerCard } from "./models/category.model";
import { Link } from "react-router-dom";
import { generateSlug } from "../../../utils/helpers"; 

const CategoryCard: React.FC<CategoryExplorerCard> = ({
  category,
  itemWidth,
  itemHeight,
  onItemClick,
}) => {
  // Tạo slug từ tên danh mục
  const categorySlug = category.slug || generateSlug(category.name);
  // console.log(category.name,":", categorySlug);
  return (
    <Link 
      to={`/${categorySlug}`} 
      className={`flex-shrink-0 ${itemWidth} cursor-pointer h-[220px] w-[180px] mr-[3rem] rounded-xl`}
      onClick={(e) => {
        // Nếu muốn giữ lại onItemClick callback đã có
        if (category.id !== undefined) {
          // Ngăn chuyển hướng trang nếu cần xử lý sự kiện click đặc biệt
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
            src={category.image_url}
            alt={category.name}
            className="w-[148px] h-[140px] object-cover transition-transform duration-300 ease-out hover:scale-110"
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