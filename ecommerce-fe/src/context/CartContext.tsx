import React, { createContext, useState, useEffect, ReactNode } from 'react';
import { CartItem, CartTotalsData } from '@/components/cart/models/cart.model';

export const defaultCartTotals: CartTotalsData = {
  subtotal: 0,
  shipping: 'Free',
  discount: 0,
  tax: 0,
  total: 0
};

export interface CartContextType {
  cartItems: CartItem[];
  cartTotals: CartTotalsData;
  error: string | null;
  setCartItems: React.Dispatch<React.SetStateAction<CartItem[]>>;
  setCartTotals: React.Dispatch<React.SetStateAction<CartTotalsData>>;
  setError: React.Dispatch<React.SetStateAction<string | null>>;
}

export const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [cartItems, setCartItems] = useState<CartItem[]>(() => {
    try {
      const savedCart = localStorage.getItem("cart");
      return savedCart ? JSON.parse(savedCart) : [];
    } catch {
      return [];
    }
  });
  const [cartTotals, setCartTotals] = useState<CartTotalsData>(defaultCartTotals);
  const [error, setError] = useState<string | null>(null);
  
  useEffect(() => {
    localStorage.setItem("cart", JSON.stringify(cartItems));
  }, [cartItems]);


  
  const value = {
    cartItems,
    cartTotals,
    error,
    setCartItems,
    setCartTotals,
    setError
  };


  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};