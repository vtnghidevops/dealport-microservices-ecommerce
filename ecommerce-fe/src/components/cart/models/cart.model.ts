export interface CartItem {
  id: string;
  name: string;
  price: number;
  originalPrice?: number;
  quantity: number;
  image: string;
}

export interface CartTotalsData {
  subtotal: number;
  shipping: string;
  discount: number;
  tax: number;
  total: number;
}