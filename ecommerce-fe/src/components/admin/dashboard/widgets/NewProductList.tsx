import React from "react";
import { Product } from '@/types/product.model';

interface NewProductListProps {
  products: Product[];
}

const NewProductList: React.FC<NewProductListProps> = ({ products = [] }) => {
  return (
    <div className="space-y-3">
      {products.slice(0, 5).map((product) => (
        <div
          key={product.id}
          className="flex items-center border-b border-gray-100 pb-2"
        >
          <div className="flex-shrink-0 w-10 h-10 bg-gray-200 rounded-md overflow-hidden">
            <img
              src={product.imageUrl}
              alt={product.name}
              className="w-full h-full object-cover"
            />
          </div>
          <div className="ml-3 flex-grow">
            <h5 className="text-sm font-medium text-gray-800">{product.name}</h5>
            <p className="text-xs text-gray-500">${product.price.toFixed(2)}</p>
          </div>
          <button className="text-xs text-primary">Details</button>
        </div>
      ))}
    </div>
  );
};

export default NewProductList;