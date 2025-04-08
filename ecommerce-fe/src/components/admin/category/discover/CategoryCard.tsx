import React from "react";

interface CategoryCardProps {
  name: string;
  image_url: string;
}

// Card component for displaying category in the discover section
const CategoryCard: React.FC<CategoryCardProps> = ({ name, image_url }) => {
  return (
    <div className="w-[242px] h-[88px] flex items-center gap-12 p-12 bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow cursor-pointer">
      <div className="flex items-center justify-center border border-neutral-300 rounded-lg">
        <img src={image_url} alt={name} className="w-[64px] h-[64px]" />
      </div>
      <h4 className="text-[18px] font-medium text-black">{name}</h4>
    </div>
  );
};

export default CategoryCard;
