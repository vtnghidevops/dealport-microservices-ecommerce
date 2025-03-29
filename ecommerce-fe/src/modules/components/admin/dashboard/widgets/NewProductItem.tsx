import React from 'react';
import { NewProduct } from '../models/product.model';
import { FiPlusCircle } from "react-icons/fi";
interface NewProductItemProps {
  product: NewProduct;
}

const NewProductItem: React.FC<NewProductItemProps> = ({ product }) => {
  return (
    <div className="flex items-center justify-between py-[8px] border-b border-neutral-200 w-[321px] h-[65px]">
      <div className="flex items-center gap-3 ">
        <div className="w-[46px] h-[46px] bg-white border border-neutral-200 rounded-md flex items-center justify-center overflow-hidden">
          <img 
            src={product.image} 
            alt={product.name} 
            className="w-full h-[35px] object-contain"
          />
        </div>
        <div>
          <div className="font-medium">{product.name}</div>
          <div className="text-green-600">${product.price.toFixed(2)}</div>
        </div>
      </div>
      <button className="bg-green-500 py-6 pl-8 pr-12 text-white rounded-full w-[62px] h-[28px] flex items-center justify-center hover:bg-green-600 transition-colors">
        <FiPlusCircle></FiPlusCircle>
        <span className='ml-4 font-medium text-[12px]'>Add</span>
      </button>
    </div>
  );
};

export default NewProductItem;