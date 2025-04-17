import React, { createContext, useState, useEffect } from 'react';
import { Product } from '@/types/product.model';

export interface WishlistContextType {
  wishlistItems: Product[];
  setWishlistItems: React.Dispatch<React.SetStateAction<Product[]>>;
}

const WishlistContext = createContext<WishlistContextType | undefined>(undefined);

export const WishlistProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [wishlistItems, setWishlistItems] = useState<Product[]>(() => {
    try {
      const savedWishlist = localStorage.getItem('wishlist');
      return savedWishlist ? JSON.parse(savedWishlist) : [];
    } catch {
      return [];
    }
  });

  useEffect(() => {
    localStorage.setItem('wishlist', JSON.stringify(wishlistItems));
  }, [wishlistItems]);

  return (
    <WishlistContext.Provider value={{ wishlistItems, setWishlistItems }}>
      {children}
    </WishlistContext.Provider>
  );
};

export { WishlistContext };