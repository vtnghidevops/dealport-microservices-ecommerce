import React from 'react';
import { BestSellingProduct } from '../models/product.model';

interface BestSellingTableProps {
  bestProducts: BestSellingProduct[];
}

const BestSellingTable: React.FC<BestSellingTableProps> = ({ bestProducts }) => {
  return (
    <div className="w-full">
      {/* Table Header */}
      <div className="mt-[1.5rem] flex items-center bg-aqua-spring rounded-lg h-[48px] text-[#62796F] font-medium text-sm ">
        <div className="w-[220px] px-[20px] py-[8px]">PRODUCT</div>
        <div className="w-[152px] pl-[20px] pr-[9px] py-[8px]">TOTAL ORDER</div>
        <div className="w-[152px] px-[22px] py-[8px]">STATUS</div>
        <div className="w-[152px] px-[20px] py-[8px]">PRICE</div>
      </div>

      {/* Table Body */}
      <div className="mt-[0.5rem]">
        {bestProducts && bestProducts.map((item) => (
          <div 
            key={item.product.id} 
            className="flex items-center h-[54px] rounded-md hover:bg-gray-50 transition-colors duration-200 cursor-pointer"
          >
            <div className="flex items-center gap-3 pl-4 w-[240px]">
              <div className="w-10 h-10 relative">
                <img 
                  src={item.product.image} 
                  alt={item.product.name}
                  className="rounded-md object-cover w-full h-full"
                />
              </div>
              <span className="font-medium text-[15px] text-[#333]">{item.product.name}</span>
            </div>
            <div className="w-[152px]">{item.totalOrder}</div>
            <div className="flex items-center w-[152px]">
              <span className={`
                inline-flex items-center px-2 py-1 text-xs rounded-full
                ${item.status === 'Stock' ? 'text-green-600 bg-green-100' : 'text-red-600 bg-red-100'}
              `}>
                <span className={`w-2 h-2 rounded-full mr-1 
                  ${item.status === 'Stock' ? 'bg-green-600' : 'bg-red-600'}
                `}></span>
                {item.status}
              </span>
            </div>
            <div className="pr-4 font-medium">
              ${item.price.toFixed(2)}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default BestSellingTable;