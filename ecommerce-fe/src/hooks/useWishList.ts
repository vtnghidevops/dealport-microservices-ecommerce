import { useContext } from 'react';
import { WishlistContext } from '../context/WishlistContext';
import { Product } from '@/types/product.model';
import { toast } from '@/hooks/use-toast';

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
        toast({
          title: `${product.name} has been added to your wishlist`,
          variant: 'success'
        });
        return [...prev, product];
      }
      return prev;
    });
  };

  const removeFromWishlist = (productId: number) => {
    setWishlistItems(prev => {
      const itemToRemove = prev.find(item => item.id === productId);
      const filtered = prev.filter(item => item.id !== productId);

      if (itemToRemove) {
        toast({
          title: `${itemToRemove.name} has been removed from your wishlist`,
          variant: 'success'
        });
      }
      return filtered;
    });
  };

  const isInWishlist = (productId: number) => {
    return wishlistItems.some(item => item.id === productId);
  };

  const clearWishlist = () => {
    setWishlistItems([]);
    toast({
      title: 'Your wishlist has been cleared',
      variant: 'success'
    });
  };

  return {
    wishlistItems,
    addToWishlist,
    removeFromWishlist: (productId: string) => removeFromWishlist(Number(productId)),
    isInWishlist: (productId: string) => isInWishlist(Number(productId)), 
    clearWishlist
  };
};