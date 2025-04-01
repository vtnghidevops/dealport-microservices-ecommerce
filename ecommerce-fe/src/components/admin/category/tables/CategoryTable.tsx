// CategoryTable component được cập nhật theo mẫu OrderTable
import React from "react";
import { Category } from "../models/category.model";

interface CategoryTableProps {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (id: number) => void;
}

const CategoryTable: React.FC<CategoryTableProps> = ({
  categories,
  onEdit,
  onDelete,
}) => {
  // Icons for product categories
  const getCategoryIcon = (icon: string) => {
    // This would map to actual icons in a real implementation
    const iconMap: Record<string, React.ReactNode> = {
      headphones: <div className="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center">🎧</div>,
      tshirt: <div className="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center">👕</div>,
      wallet: <div className="w-6 h-6 bg-brown-100 rounded-full flex items-center justify-center">👛</div>,
      pillow: <div className="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center">🛏️</div>,
      coffee: <div className="w-6 h-6 bg-brown-100 rounded-full flex items-center justify-center">☕</div>,
      cap: <div className="w-6 h-6 bg-green-100 rounded-full flex items-center justify-center">🧢</div>,
      webcam: <div className="w-6 h-6 bg-black-100 rounded-full flex items-center justify-center">🎥</div>,
      bulb: <div className="w-6 h-6 bg-yellow-100 rounded-full flex items-center justify-center">💡</div>,
    };
    
    return iconMap[icon] || <div className="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center">📦</div>;
  };

  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse bg-white rounded-lg">
        <thead className="bg-aqua-spring">
          <tr className="h-[56px]">
            <th className="py-4 px-6 text-center text-[15px] font-medium text-cyprus">No.</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Product</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Created Date</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Order</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Actions</th>
          </tr>
        </thead>
        <tbody>
          {categories.map((category, index) => (
            <tr key={category.id} className="border-b h-[68px]">
              <td className="py-4 px-6">
                <div className="flex items-center justify-center">
                  <input
                    type="checkbox"
                    className="mr-2 h-[16px] w-[16px] rounded border border-ocean-green checked:bg-success checked:border-success appearance-none relative checked:after:content-['✓'] checked:after:text-white checked:after:absolute checked:after:text-xs checked:after:left-[3px]"
                  />
                  <span>{index + 1}</span>
                </div>
              </td>
              <td className="py-4 px-6">
                <div className="flex items-center">
                  <div className="border border-neutral-200 w-[40px] h-[40px] mr-3 rounded flex items-center justify-center overflow-hidden">
                    {getCategoryIcon(category.icon)}
                  </div>
                  <span className="text-[15px] max-w-[140px]">{category.name}</span>
                </div>
              </td>
              <td className="py-4 px-6 text-[15px]">
                {category.createdAt}
              </td>
              <td className="py-4 px-6 text-[15px]">
                {category.productCount || 0}
              </td>
              <td className="py-4 px-6">
                <div className="flex space-x-2">
                  <button 
                    onClick={() => onEdit(category)}
                    className="text-gray-500 hover:text-blue-600"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 0L19 7.586l-5.586 5.586-2.762.552.552-2.762L19 4.414z" />
                    </svg>
                  </button>
                  <button 
                    onClick={() => onDelete(category.id)}
                    className="text-gray-500 hover:text-red-600"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default CategoryTable;