import React from 'react';
import { CartItem } from '@/types/cart.model';
import CartItemRow from './CartItemRow';

interface ShoppingCartTableProps {
  cartItems: CartItem[];
  updateQuantity: (id: string, quantity: number) => void;
  removeFromCart: (id: string) => void;
}

const ShoppingCartTable: React.FC<ShoppingCartTableProps> = ({ 
  cartItems, 
  updateQuantity, 
  removeFromCart 
}) => {

  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-100">
            <tr>
              <th scope="col" className="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Products
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Price
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Quantity
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Sub-total
              </th>
              <th scope="col" className="relative px-6 py-3">
                <span className="sr-only">Remove</span>
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {cartItems.map((item) => (
              <CartItemRow 
                key={item.id} 
                item={item} 
                updateQuantity={updateQuantity} 
                removeFromCart={removeFromCart} 
              />
            ))}
          </tbody>
        </table>
      </div>
      
      <div className="ml-5 mr-5 mb-8 px-6 py-4 flex justify-between items-center bg-white">
        <a 
          href="/" 
          className="text-sans inline-flex items-center px-5 py-2 border-2 border-[#0496FF] text-[#0496FF] rounded hover:bg-blue-50"
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          RETURN TO SHOP
        </a>
        
        <button 
          className="text-sans inline-flex items-center px-5 py-2 border-2 border-[#0496FF] text-[#0496FF] rounded hover:bg-blue-50"
        >
          Update Cart
        </button>
      </div>
    </div>
  );
};

export default ShoppingCartTable;
