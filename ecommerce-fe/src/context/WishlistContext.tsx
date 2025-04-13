import React, { createContext, useContext, useState, useEffect } from 'react';
import { Product } from '../components/product/models/product.model';
import { enqueueSnackbar } from 'notistack';

interface WishlistContextType {
  wishlistItems: Product[];
  addToWishlist: (product: Product) => void;
  removeFromWishlist: (productId: string) => void;
  isInWishlist: (productId: string) => boolean;
  clearWishlist: () => void;
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

  const showNotification = (message: string) => {
    enqueueSnackbar(message, {
      variant: 'success',
      autoHideDuration: 2000,
      anchorOrigin: {
        vertical: 'bottom',
        horizontal: 'right'
      },
      preventDuplicate: true // Prevent duplicate notifications
    });
  };

  const addToWishlist = (product: Product) => {
    setWishlistItems(prev => {
      const exists = prev.some(item => item.id === product.id);
      if (!exists) {
        showNotification(`${product.name} has been added to your wishlist`);
        return [...prev, product];
      }
      return prev;
    });
  };
  const removeFromWishlist = (productId: string) => {
    setWishlistItems(prev => {
      const itemToRemove = prev.find(item => item.id === productId);
      const filtered = prev.filter(item => item.id !== productId);
      
      if (itemToRemove) {
        showNotification(`${itemToRemove.name} has been removed from your wishlist`);
      }
      return filtered;
    });
  };
  
  const isInWishlist = (productId: string) => {
    return wishlistItems.some(item => item.id === productId);
  };

  const clearWishlist = () => {
    setWishlistItems([]);
    enqueueSnackbar('Your wishlist has been cleared', {
      variant: 'success',
      autoHideDuration: 2000,
      anchorOrigin: {
        vertical: 'bottom',
        horizontal: 'right'
      }
    });
  };

  const value = {
    wishlistItems,
    addToWishlist,
    removeFromWishlist,
    isInWishlist,
    clearWishlist
  };
  return (
    <WishlistContext.Provider value={value}>
      {children}
    </WishlistContext.Provider>
  );
};

export const useWishlist = (): WishlistContextType => {
  const context = useContext(WishlistContext);
  if (!context) {
    throw new Error('useWishlist must be used within a WishlistProvider');
  }
  return context;
};
