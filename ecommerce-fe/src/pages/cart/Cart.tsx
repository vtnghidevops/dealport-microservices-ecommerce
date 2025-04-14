import React from 'react';
import ShoppingCartTable from '../../components/cart/ShoppingCartTable';
import CartTotals from '../../components/cart/CartTotals';
import CouponCode from '../../components/cart/CouponCode';
import { useCart } from '@/hooks/useCart';
import EmptyCart from '../system/EmptyCart';
// import Loading from '@/components/shared/Loading';

const Cart: React.FC = () => {
  const { cartItems, error, updateQuantity, removeFromCart } = useCart();

  // if (isLoading) return <Loading size="large" fullscreen />;
  if (error) return <div className="text-red-500 p-8">{error}</div>;
  if (cartItems.length === 0) return <EmptyCart />;
  return (
    <div className="max-w-[1440px] min-w-[1024px] mx-auto px-[5rem] py-[1.5rem] ">
      <div className="flex flex-col lg:flex-row gap-8 min-h-[calc(100vh-200px)]">
        <div className="mr-5 lg:w-2/3 rounded-lg shadow p-5 h-fit">
          <h1 className="text-[18px] font-bold mb-5">Shopping Cart</h1>
          <ShoppingCartTable
            key={cartItems.length}
            cartItems={cartItems} 
            updateQuantity={updateQuantity} 
            removeFromCart={removeFromCart} 
          />
        </div>

        {cartItems.length > 0 && (
          <div className="lg:w-[35%] lg:min-w-[320px] ">
            <CartTotals />
            <div className="mt-5">
              <CouponCode />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Cart;