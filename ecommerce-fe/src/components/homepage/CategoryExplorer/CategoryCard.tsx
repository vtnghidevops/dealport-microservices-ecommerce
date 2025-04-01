import React from "react";
import { CategoryExplorerCard } from "./models/category.model";

const CategoryCard: React.FC<CategoryExplorerCard> = ({
  category,
  itemWidth,
  itemHeight,
  onItemClick,
}) => {
  return (
    <div
      className={`flex-shrink-0 ${itemWidth} cursor-pointer h-[220px] w-[180px] mr-[3rem] rounded-xl`}
      onClick={() =>
        category.id !== undefined
          ? onItemClick?.(category, category.id as number)
          : undefined
      }
    >
      <div className="relative rounded-lg overflow-hidden shadow-sm border border-gray-300 h-full">
        <div
          className={`${itemHeight} w-full flex items-center justify-center overflow-hidden`}
        >
          <img
            src={category.image}
            alt={category.name}
            className="w-[148px] h-[140px] object-cover transition-transform duration-300 ease-out hover:scale-110"
          />
        </div>
        <div className="p-2 text-center absolute bottom-0 flex justify-center w-full bg-white">
          <h3 className="text-sm font-medium text-gray-800">{category.name}</h3>
        </div>
      </div>
    </div>
  );
};

export default CategoryCard;