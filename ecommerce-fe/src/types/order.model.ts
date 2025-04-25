export type OrderStatus = 'Pending' | 'Processing' | 'Shipped' | 'Delivered' | 'Cancelled';

export interface Order {
  id: number;
  orderId: string;
  customerId: string;
  customerName: string;
  products: OrderProduct[];
  date: string;
  totalAmount: number;
  status: OrderStatus;
  paymentStatus: 'Paid' | 'Unpaid';
}

export interface OrderProduct {
  productId: number;
  productName: string;
  productImage: string;
  quantity: number;
  price: number;
  id?: number;
  name?: string;
  imageUrl?: string;
  category?: string;
}

export interface OrderSummary {
  totalOrders: number;
  newOrders: number;
  completedOrders: number;
  cancelledOrders: number;
  lastUpdated: string;
  growthRate: {
    total: number;
    new: number;
    completed: number;
    cancelled: number;
  };
}

export interface OrderFilterParams {
  status?: OrderStatus;
  dateRange?: [Date, Date];
  customerId?: string;
  productId?: string;
  searchTerm?: string;
  page: number;
  limit: number;
}

export interface NewOrderData {
  customerId: string;
  customerName: string;
  products: Product[];
  totalAmount: number;
  paymentStatus: string;
}

// Cần import Product từ product.model.ts
import { Product } from './product.model';