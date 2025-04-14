import { useContext, useEffect } from 'react';
import { CartContext, defaultCartTotals } from '@/context/CartContext';
import { CartItem } from '@/components/cart/models/cart.model';
import { showNotification } from '@/utils/notifications';
export interface UseCartReturn {
  cartItems: CartItem[];
  cartTotals: typeof defaultCartTotals;
  error: string | null;
  updateQuantity: (id: string, quantity: number) => void;
  removeFromCart: (id: string) => void;
  applyCoupon: (couponCode: string) => Promise<boolean>;
  addToCart: (product: Omit<CartItem, 'quantity'>, quantity?: number) => void;
}

export const useCart = (): UseCartReturn => {
  const context = useContext(CartContext);
  
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }

  const { cartItems, setCartItems, cartTotals, setCartTotals, error, setError } = context;
  useEffect(() => {
    calculateTotals(cartItems);
  }, [cartItems]);

  const calculateTotals = (items: CartItem[]) => {
    if (items.length === 0) {
      setCartTotals(defaultCartTotals);
      return;
    }

    const subtotal = items.reduce((total, item) => total + (item.price * item.quantity), 0);
    const tax = Number((subtotal * 0.18).toFixed(2));
    const discount = items.length > 0 ? 24 : 0;
    
    setCartTotals({
      subtotal,
      shipping: 'Free',
      discount,
      tax,
      total: Number((subtotal - discount + tax).toFixed(2))
    });
  };

  const addToCart = (product: Omit<CartItem, 'quantity'>, quantity: number = 1) => {
    try {
      setCartItems(prevItems => {
        const existingItemIndex = prevItems.findIndex(item => item.id === product.id);
        
        let newItems;
        if (existingItemIndex !== -1) {
          newItems = [...prevItems];
          newItems[existingItemIndex].quantity += quantity;
          showNotification(`Updated ${product.name} quantity in cart`);
        } else {
          newItems = [...prevItems, { ...product, quantity }];
          showNotification(`Added ${product.name} to cart`);
        }
    
        return newItems;
      });
      
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to add item to cart';
      setError(errorMessage);
      showNotification(errorMessage, 'error');
    }
  };


  const updateQuantity = (id: string, quantity: number) => {
    if (quantity < 1) {
      showNotification('Quantity cannot be less than 1', 'warning');
      return;
    }
    
    try {
      setCartItems(prevItems => {
        const newItems = prevItems.map(item =>
          item.id === id ? { ...item, quantity } : item
        );
        showNotification('Cart quantity updated successfully');
        return newItems;
      });
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update quantity';
      setError(errorMessage);
      showNotification(errorMessage, 'error');
    }
  };

  const removeFromCart = (id: string) => {
    try {
      setCartItems(prevItems => {
        const itemToRemove = prevItems.find(item => item.id === id);
        const newItems = prevItems.filter(item => item.id !== id);
        if (itemToRemove) {
          showNotification(`Removed ${itemToRemove.name} from cart`);
        }
        return newItems;
      });
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to remove item';
      setError(errorMessage);
      showNotification(errorMessage, 'error');
    }
  };

  const applyCoupon = async (couponCode: string): Promise<boolean> => {
    try {
      if (couponCode === 'DISCOUNT10') {
        setCartTotals(prev => ({
          ...prev,
          discount: 10,
          total: prev.total - 10
        }));
        showNotification('Coupon applied successfully!');
        return true;
      } else {
        showNotification('Invalid coupon code', 'error');
        return false;
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to apply coupon';
      setError(errorMessage);
      showNotification(errorMessage, 'error');
      return false;
    }
  };


  return {
    cartItems,
    cartTotals,
    error,
    addToCart,
    updateQuantity,
    removeFromCart,
    applyCoupon
  };
};