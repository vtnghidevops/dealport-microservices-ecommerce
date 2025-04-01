import React from 'react';
import NewProductItem from './NewProductItem';
import { NewProduct } from '../models/product.model';


interface NewProductListProps {
  products: NewProduct[];
}

const NewProductList: React.FC<NewProductListProps> = ({ products }) => {
  return (
    <div className="mt-6">
      <h2 className="text-neutral-500 text-[14px] font-medium mb-[10px]">Product</h2>
      
      <div>
        {products.map((product) => (
          <NewProductItem key={product.id} product={product} />
        ))}
      </div>
      
      <div className="mt-[1.5rem] flex justify-center">
        <button className="text-primary font-medium">See more</button>
      </div>
    </div>
  );
};

export default NewProductList;