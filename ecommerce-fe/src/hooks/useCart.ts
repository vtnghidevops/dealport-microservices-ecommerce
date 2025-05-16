import { useContext, useEffect, useState, useRef } from 'react';
import { CartContext, defaultCartTotals } from '@/context/CartContext';
// import { CartItem } from '@/components/cart/models/cart.model';
import { CartItem } from '@/types/cart.model';
import { useToast } from '@/hooks/use-toast';
import cartService from '@/services/user/cart.service';
import { useAuth } from './useAuth';

export interface UseCartReturn {
  cartItems: CartItem[];
  cartTotals: typeof defaultCartTotals;
  error: string | null;
  isLoading: boolean;
  updateQuantity: (id: string, quantity: number) => void;
  removeFromCart: (id: string) => void;
  applyCoupon: (couponCode: string) => Promise<boolean>;
  removeCoupon: () => Promise<boolean>;
  addToCart: (product: Omit<Omit<CartItem, 'quantity'>, 'id'>, quantity?: number) => void;
  clearCart: () => void;
}

export const useCart = (): UseCartReturn => {
  const context = useContext(CartContext);
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const { authState } = useAuth();
  const isAuthenticated = authState.isAuthenticated;
  const cartFetchedRef = useRef(false);
  const debouncedFetchRef = useRef<NodeJS.Timeout | null>(null);

  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }

  const { cartItems, setCartItems, cartTotals, setCartTotals, error, setError } = context;

  // Lấy giỏ hàng từ server khi user đăng nhập
  useEffect(() => {
    // Chỉ gọi fetchCart khi user đăng nhập và cart chưa được tải
    if (isAuthenticated) {
      // Xóa timeout cũ nếu có
      if (debouncedFetchRef.current) {
        clearTimeout(debouncedFetchRef.current);
      }

      // Debounce 100ms để tránh gọi API nhiều lần
      debouncedFetchRef.current = setTimeout(() => {
        // Kiểm tra xem đã tải cart trong session này chưa
        if (!cartFetchedRef.current) {
          fetchCart();
          cartFetchedRef.current = true;
        }
      }, 100);
    } else {
      // Reset flag khi logout
      cartFetchedRef.current = false;
    }

    return () => {
      if (debouncedFetchRef.current) {
        clearTimeout(debouncedFetchRef.current);
      }
    };
  }, [isAuthenticated]);

  // Tính toán tổng tiền dựa trên cart items
  useEffect(() => {
    if (!isAuthenticated) {
      calculateTotals(cartItems);
    }
  }, [cartItems, isAuthenticated]);

  const fetchCart = async () => {
    if (!isAuthenticated) return;

    try {
      setIsLoading(true);
      const cartData = await cartService.getCart();

      // Cập nhật cartItems từ response API
      setCartItems(cartData.items);

      // Cập nhật cartTotals từ response API
      setCartTotals({
        subtotal: cartData.totals.subtotal,
        shipping: cartData.totals.shipping,
        discount: cartData.totals.discount,
        tax: cartData.totals.tax,
        total: cartData.totals.total
      });
    } catch (error: unknown) {
      handleApiError(error);
    } finally {
      setIsLoading(false);
    }
  };

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

  const handleApiError = (error: unknown) => {
    const errorMessage = error instanceof Error ? error.message : 'Error with cart operation';
    setError(errorMessage);
    toast({
      title: 'Error',
      description: errorMessage,
      variant: 'destructive'
    });
  };

  const addToCart = async (product: Omit<Omit<CartItem, 'quantity'>, 'id'>, quantity: number = 1) => {
    try {
      setIsLoading(true);

      if (isAuthenticated) {
        // Gọi API để thêm vào giỏ hàng server
        const response = await cartService.addCartItem({
          productId: product.productId,
          name: product.name,
          price: product.price,
          originalPrice: product.originalPrice,
          quantity: quantity,
          imageUrl: product.imageUrl
        });

        // Cập nhật state từ response
        setCartItems(response.items);
        setCartTotals({
          subtotal: response.totals.subtotal,
          shipping: response.totals.shipping,
          discount: response.totals.discount,
          tax: response.totals.tax,
          total: response.totals.total
        });

        toast({
          title: `Added ${product.name} to cart`,
          description: 'You can now proceed to checkout',
          variant: 'success'
        });
      } else {
        // Xử lý cart local khi chưa đăng nhập
        setCartItems(prevItems => {
          const existingItemIndex = prevItems.findIndex(item => item.productId === product.productId);

          let newItems;
          if (existingItemIndex !== -1) {
            newItems = [...prevItems];
            newItems[existingItemIndex].quantity += quantity;
            toast({
              title: `Updated ${product.name} quantity in cart`,
              variant: 'success'
            });
          } else {
            // Tạo ID tạm thời cho cart item khi lưu local
            const tempId = `temp_${Date.now()}_${product.productId}`;
            newItems = [...prevItems, { ...product, id: tempId, quantity }];
            toast({
              title: `Added ${product.name} to cart`,
              description: 'You can now proceed to checkout',
              variant: 'success'
            });
          }

          return newItems;
        });
      }
    } catch (error: unknown) {
      handleApiError(error);
    } finally {
      setIsLoading(false);
    }
  };

  const updateQuantity = async (id: string, quantity: number) => {
    if (quantity < 1) {
      toast({
        title: 'Warning',
        description: 'Quantity cannot be less than 1',
        variant: 'success'
      });
      return;
    }

    try {
      setIsLoading(true);

      if (isAuthenticated) {
        // Gọi API cập nhật quantity
        const response = await cartService.updateCartItem(id, quantity);

        // Cập nhật state từ response
        setCartItems(response.items);
        setCartTotals({
          subtotal: response.totals.subtotal,
          shipping: response.totals.shipping,
          discount: response.totals.discount,
          tax: response.totals.tax,
          total: response.totals.total
        });
      } else {
        // Xử lý local
        setCartItems(prevItems => {
          const newItems = prevItems.map(item =>
            item.id === id ? { ...item, quantity } : item
          );
          return newItems;
        });
      }

      toast({
        title: 'Cart quantity updated successfully',
        variant: 'success'
      });
    } catch (error: unknown) {
      handleApiError(error);
    } finally {
      setIsLoading(false);
    }
  };

  const removeFromCart = async (id: string) => {
    try {
      setIsLoading(true);

      if (isAuthenticated) {
        // Gọi API xóa sản phẩm
        const response = await cartService.removeCartItem(id);

        // Lấy tên sản phẩm trước khi xóa
        const itemToRemove = cartItems.find(item => item.id === id);

        // Cập nhật state từ response
        setCartItems(response.items);
        setCartTotals({
          subtotal: response.totals.subtotal,
          shipping: response.totals.shipping,
          discount: response.totals.discount,
          tax: response.totals.tax,
          total: response.totals.total
        });

        if (itemToRemove) {
          toast({
            title: `Removed ${itemToRemove.name} from cart`,
            variant: 'success'
          });
        }
      } else {
        // Xử lý local
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
      }
    } catch (error: unknown) {
      handleApiError(error);
    } finally {
      setIsLoading(false);
    }
  };

  const clearCart = async () => {
    try {
      setIsLoading(true);

      if (isAuthenticated) {
        // Gọi API xóa giỏ hàng
        const response = await cartService.clearCart();

        // Kiểm tra response từ API
        if (response.success) {
          // Cập nhật state
          setCartItems([]);
          setCartTotals(defaultCartTotals);
          // Reset flag để cho phép tải lại giỏ hàng
          cartFetchedRef.current = false;

          toast({
            title: response.message || 'Cart cleared successfully',
            variant: 'success'
          });
        } else {
          toast({
            title: 'Failed to clear cart',
            description: response.message || 'An error occurred',
            variant: 'destructive'
          });
        }
      } else {
        // Xử lý local
        setCartItems([]);
        setCartTotals(defaultCartTotals);
        // Reset flag để cho phép tải lại giỏ hàng
        cartFetchedRef.current = false;

        toast({
          title: 'Cart cleared successfully',
          variant: 'success'
        });
      }
    } catch (error: unknown) {
      handleApiError(error);
    } finally {
      setIsLoading(false);
    }
  };

  const applyCoupon = async (couponCode: string): Promise<boolean> => {
    if (!couponCode) {
      toast({
        title: 'Error',
        description: 'Please enter a coupon code',
        variant: 'destructive'
      });
      return false;
    }

    try {
      setIsLoading(true);

      if (isAuthenticated) {
        // Apply coupon via API
        const response = await cartService.applyCoupon(couponCode);

        // Update state from response
        setCartItems(response.items);
        setCartTotals({
          subtotal: response.totals.subtotal,
          shipping: response.totals.shipping,
          discount: response.totals.discount,
          tax: response.totals.tax,
          total: response.totals.total
        });

        // Store coupon code in localStorage for checkout
        localStorage.setItem('appliedCoupon', couponCode);

        toast({
          title: 'Coupon Applied',
          description: `You saved ${response.totals.discount} with this coupon`,
          variant: 'success'
        });
        return true;
      } else {
        // Handle local cart coupon
        // This is just a simulation for local cart
        const discount = 10; // Simulate a $10 discount
        setCartTotals((prev: typeof defaultCartTotals) => ({
          ...prev,
          discount,
          total: prev.total - discount
        }));

        // Store coupon code in localStorage for checkout
        localStorage.setItem('appliedCoupon', couponCode);

        toast({
          title: 'Coupon Applied',
          description: `You saved $${discount} with this coupon`,
          variant: 'success'
        });
        return true;
      }
    } catch (error: unknown) {
      handleApiError(error);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const removeCoupon = async (): Promise<boolean> => {
    try {
      setIsLoading(true);

      if (isAuthenticated) {
        // Remove coupon via API
        const response = await cartService.removeCoupon();

        // Update state from response
        setCartItems(response.items);
        setCartTotals({
          subtotal: response.totals.subtotal,
          shipping: response.totals.shipping,
          discount: response.totals.discount,
          tax: response.totals.tax,
          total: response.totals.total
        });

        // Remove coupon code from localStorage
        localStorage.removeItem('appliedCoupon');

        toast({
          title: 'Coupon Removed',
          variant: 'success'
        });
        return true;
      } else {
        // Handle local cart
        setCartTotals((prev: typeof defaultCartTotals) => ({
          ...prev,
          discount: 0,
          total: prev.subtotal + prev.tax
        }));

        // Remove coupon code from localStorage
        localStorage.removeItem('appliedCoupon');

        toast({
          title: 'Coupon Removed',
          variant: 'success'
        });
        return true;
      }
    } catch (error: unknown) {
      handleApiError(error);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    cartItems,
    cartTotals,
    error,
    isLoading,
    addToCart,
    updateQuantity,
    removeFromCart,
    applyCoupon,
    removeCoupon,
    clearCart
  };
};