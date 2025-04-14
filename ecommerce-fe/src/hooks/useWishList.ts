import { useContext } from 'react';
import { WishlistContext } from '../context/WishlistContext';
import { Product } from '../components/product/models/product.model';
import { showNotification } from '@/utils/notifications';
export interface UseWishlistReturn {
  wishlistItems: Product[];
  addToWishlist: (product: Product) => void;
  removeFromWishlist: (productId: string) => void;
  isInWishlist: (productId: string) => boolean;
  clearWishlist: () => void;
}

export const useWishlist = (): UseWishlistReturn => {
  const context = useContext(WishlistContext);
  
  if (!context) {
    throw new Error('useWishlist must be used within a WishlistProvider');
  }

  const { wishlistItems, setWishlistItems } = context;

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
    showNotification('Your wishlist has been cleared');
  };

  return {
    wishlistItems,
    addToWishlist,
    removeFromWishlist,
    isInWishlist,
    clearWishlist
  };
};