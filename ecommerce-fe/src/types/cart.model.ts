export interface CartItem {
  id?: string;
  productId: number;
  name: string;
  price: number;
  originalPrice?: number;
  quantity: number;
  imageUrl: string;
}

export interface CartTotalsData {
  subtotal: number;
  shipping: string;
  discount: number;
  tax: number;
  total: number;
}