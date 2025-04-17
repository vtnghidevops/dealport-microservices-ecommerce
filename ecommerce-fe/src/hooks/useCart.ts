import { useContext, useEffect } from 'react';
import { CartContext, defaultCartTotals } from '@/context/CartContext';
import { CartItem } from '@/components/cart/models/cart.model';
import { useToast } from '@/hooks/use-toast';
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
  const { toast } = useToast();
  
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
          toast({
            title: `Updated ${product.name} quantity in cart`,
            variant: 'success'
          });
        } else {
          newItems = [...prevItems, { ...product, quantity }];
          toast({
            title: `Added ${product.name} to cart`,
            description: 'You can now proceed to checkout',
            variant: 'success'
          });
        }
    
        return newItems;
      });
      
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to add item to cart';
      setError(errorMessage);
      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'error'
      });
    }
  };


  const updateQuantity = (id: string, quantity: number) => {
    if (quantity < 1) {
      toast({
        title: 'Warning',
        description: 'Quantity cannot be less than 1',
        variant: 'success'
      });
      return;
    }
    
    try {
      setCartItems(prevItems => {
        const newItems = prevItems.map(item =>
          item.id === id ? { ...item, quantity } : item
        );
        toast({
          title: 'Cart quantity updated successfully',
          variant: 'success'
        });
        return newItems;
      });
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update quantity';
      setError(errorMessage);
      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'destructive'
      });
    }
  };

  const removeFromCart = (id: string) => {
    try {
      setCartItems(prevItems => {
        const itemToRemove = prevItems.find(item => item.id === id);
        const newItems = prevItems.filter(item => item.id !== id);
        if (itemToRemove) {
          toast({
            title: `Removed ${itemToRemove.name} from cart`,
            variant: 'success'
          });
        }
        return newItems;
      });
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to remove item';
      setError(errorMessage);
      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'destructive'
      });
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
        toast({
          title: 'Coupon applied successfully!',
          variant: 'success'
        });
        return true;
      } else {
        toast({
          title: 'Invalid coupon code',
          variant: 'destructive'
        });
        return false;
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to apply coupon';
      setError(errorMessage);
      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'destructive'
      });
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