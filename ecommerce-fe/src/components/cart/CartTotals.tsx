import React from 'react';
import { useCart } from '@/hooks/useCart';
import { useNavigate } from 'react-router-dom';

const CartTotals: React.FC = () => {
  const { cartTotals } = useCart();
  const navigate = useNavigate();
  return (
    <div className="bg-white rounded-lg shadow p-5 ">
      <h2 className="text-lg font-bold mb-4">Cart Totals</h2>
      
      <div className="border-b border-gray-200 py-3 flex justify-between">
        <span className="text-gray-600">Sub-total</span>
        <span className="font-bold">${cartTotals.subtotal.toFixed(0)}</span>
      </div>
      
      <div className="border-b border-gray-200 py-3 flex justify-between">
        <span className="text-gray-600">Shipping</span>
        <span className="font-bold">{cartTotals.shipping}</span>
      </div>
      
      <div className="border-b border-gray-200 py-3 flex justify-between">
        <span className="text-gray-600">Discount</span>
        <span className="font-bold">${cartTotals.discount.toFixed(0)}</span>
      </div>
      
      <div className="border-b border-gray-200 py-3 flex justify-between">
        <span className="text-gray-600">Tax</span>
        <span className="font-bold">${cartTotals.tax.toFixed(0)}</span>
      </div>
      
      <div className="py-3 flex justify-between items-center">
        <span className="text-lg font-medium">Total</span>
        <span className="text-sans text-lg font-bold">${cartTotals.total.toFixed(0)} USD</span>
      </div>
      
      <button onClick={() => navigate('/user/checkout')} className="text-sans mt-4 w-full bg-[#FA8232] text-white py-3 px-4 rounded flex items-center justify-center font-medium hover:bg-orange-600 transition-colors">
        PROCEED TO CHECKOUT
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </button>
    </div>
  );
};

export default CartTotals;
