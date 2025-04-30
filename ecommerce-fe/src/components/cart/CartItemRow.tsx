import React from 'react';
import { CartItem } from '@/types/cart.model';
import { CiCircleRemove } from "react-icons/ci";

interface CartItemRowProps {
  item: CartItem;
  updateQuantity: (id: string, quantity: number) => void;
  removeFromCart: (id: string) => void;
}

const CartItemRow: React.FC<CartItemRowProps> = ({
  item,
  updateQuantity,
  removeFromCart
}) => {
  const handleQuantityChange = (newQuantity: number) => {
    if (newQuantity >= 1) {
      updateQuantity(item.id, newQuantity);
    }
  };

  const formatPrice = (price: number) => {
    return `$${price}`;
  };

  return (
    <tr className='h-[80px] hover:bg-neutral-50 border-b border-gray-100 transition-colors duration-150'>
      <td className="px-6 py-4 w-[40%]">
        <div className="ml-5 flex items-center">
          <button
            onClick={() => removeFromCart(item.id)}
            className="mr-8 text-gray-400 hover:text-red-500 transform hover:scale-110 transition-all duration-200 ease-in-out"
            aria-label="Remove item"
          >
            <CiCircleRemove className='h-5 w-5' />
          </button>
          <div className="ml-1 h-[45px] w-[45px] flex-shrink-0 border border-gray-100 rounded overflow-hidden">
            <img className="h-full w-full object-contain" src={item.imageUrl} alt={item.name} />
          </div>
          <div className="ml-5">
            <div className="text-sm font-medium text-gray-900 line-clamp-2">{item.name}</div>
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap w-1/8">
        <div className="flex items-center">
          {item.originalPrice && (
            <span className="mr-2 text-xs text-gray-400 line-through">${item.originalPrice}</span>
          )}
          <span className="text-sm font-bold text-gray-900">{formatPrice(item.price)}</span>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap w-1/10">
        <div className="flex border border-gray-300 rounded w-[68%] shadow-sm">
          <button
            onClick={() => handleQuantityChange(item.quantity - 1)}
            disabled={item.quantity <= 1}
            className="px-3 py-1 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            aria-label="Decrease quantity"
          >
            —
          </button>
          <div className="px-3 py-1 min-w-[40px] text-center border-x border-gray-200">
            {item.quantity.toString().padStart(2, '0')}
          </div>
          <button
            onClick={() => handleQuantityChange(item.quantity + 1)}
            className="px-3 py-1 hover:bg-gray-100"
            aria-label="Increase quantity"
          >
            +
          </button>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap w-1/8 text-left">
        <span className="text-sm font-bold text-gray-900">
          {formatPrice(Number((item.price * item.quantity).toFixed(2)))}
        </span>
      </td>
    </tr>
  );
};

export default CartItemRow;