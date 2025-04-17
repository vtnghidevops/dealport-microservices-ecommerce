import { ReactNode } from 'react';
import { WishlistProvider } from '@/context/WishlistContext';
import { CartProvider } from '@/context/CartContext';
import { Toaster } from '@/components/ui/toaster';
interface AppProviderProps {
  children: ReactNode;
}

export const AppProvider = ({ children }: AppProviderProps) => {
  return (
    <CartProvider>
      <WishlistProvider>
        {children}
        <Toaster />
      </WishlistProvider>
    </CartProvider>
  );
};