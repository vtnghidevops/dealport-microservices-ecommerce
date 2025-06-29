import { Product } from './product.model';
import { CartItem } from './cart.model';

export type OrderStatus =
  | 'pending'
  | 'processing'
  | 'paid'
  | 'shipped'
  | 'delivered'
  | 'canceled'
  | 'refunded';

export interface Customer {
  name: string;
  email: string;
  phone: string;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  country: string;
  zipCode: string;
}

export interface ShippingMethod {
  id: string;
  name: string;
  price: number;
  estimatedDelivery: string;
}

export interface PaymentMethod {
  id: string;
  name: string;
  cardNumber?: string;
  expiryDate?: string;
}

export interface Order {
  id: string;
  userId: string;
  orderNumber: string;
  items: {
    product: Product;
    quantity: number;
    price: number;
  }[];
  billingInfo: BillingInfo;
  shippingInfo: ShippingInfo;
  status: OrderStatus;
  subtotal: number;
  tax: number;
  shippingCost: number;
  discount: number;
  total: number;
  couponCode?: string;
  paymentMethod: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CheckoutState {
  customer: Customer;
  shippingAddress: Address;
  billingAddress: Address;
  shippingMethod: ShippingMethod | null;
  paymentMethod: PaymentMethod | null;
  order: Order | null;
  loading: boolean;
  error: string | null;
}

export interface CheckoutFeature {
  placeOrder: (order: Order) => Promise<Order>;
  getShippingMethods: () => Promise<ShippingMethod[]>;
  getPaymentMethods: () => Promise<PaymentMethod[]>;
  validateCustomerInfo: (customer: Customer) => {
    valid: boolean;
    message?: string;
  }
}

export interface OrderProduct {
  productId: number;
  name: string;
  price: number;
  quantity: number;
  imageUrl?: string;
  subtotal?: number;
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

// Interfaces aligned with server responses
export interface BillingInfo {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
}

export interface ShippingInfo {
  shipToDifferentAddress: boolean;
  firstName?: string;
  lastName?: string;
  companyName?: string;
  address?: string;
  country?: string;
  region?: string;
  city?: string;
  zipCode?: string;
  shippingMethod: string;
  shippingCost?: number;
}

export interface PaymentInfo {
  paymentMethod: string;
  transactionId?: string;
  status?: string;
  amount?: number;
  currency?: string;
  paymentDate?: string;
}

export interface OrderTotals {
  subtotal: number;
  shipping: number | string;
  discount: number;
  tax: number;
  total: number;
}

export interface CheckoutOrder {
  id: string;
  userId: string;
  orderNumber: string;
  status: string;
  items: CartItem[];
  billingInfo: BillingInfo;
  shippingInfo: ShippingInfo;
  paymentInfo: PaymentInfo;
  totals: OrderTotals;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ValidationError {
  field: string;
  message: string;
}

export interface CheckoutValidationResponse {
  valid: boolean;
  errors: ValidationError[];
}

export interface PaymentResponse {
  success: boolean;
  transactionId?: string;
  status: string;
  redirectUrl?: string;
  orderId?: string;
  data?: any;
}