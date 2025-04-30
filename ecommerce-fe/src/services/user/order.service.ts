import axios from 'axios';

// Base API URL already includes the /api/v1 prefix
const API_URL = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8080/api/v1';

// Get auth header for authenticated requests
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  if (!token) return {};
  return { Authorization: `Bearer ${token}` };
};

// Order status types
export enum OrderStatus {
  Pending = 'pending',
  Processing = 'processing',
  Paid = 'paid',
  Shipped = 'shipped',
  Delivered = 'delivered',
  Cancelled = 'cancelled',
  Refunded = 'refunded'
}

// Order item interface
export interface OrderItem {
  id: string;
  productId: string;
  name: string;
  price: number;
  quantity: number;
  imageUrl?: string;
  subtotal?: number;
}

// Shipping info interface
export interface ShippingInfo {
  address: string;
  city: string;
  companyName: string;
  country: string;
  firstName: string;
  lastName: string;
  region: string;
  shipToDifferentAddress: boolean;
  shippingCost: number;
  shippingMethod: string;
  zipCode: string;
}

// Billing info interface
export interface BillingInfo {
  address: string;
  city: string;
  companyName: string;
  country: string;
  email: string;
  firstName: string;
  lastName: string;
  phone: string;
  region: string;
  zipCode: string;
}

// Payment info interface
export interface PaymentInfo {
  amount: number;
  currency: string;
  paymentDate: string;
  paymentMethod: string;
  status: string;
  transactionId: string;
}

// Order interface
export interface Order {
  id: string;
  userId: string;
  status: OrderStatus;
  items: OrderItem[];
  total: number;
  subtotal: number;
  tax: number;
  shipping: number | string;
  discount: number;
  createdAt: string;
  updatedAt: string;
  paymentMethod: string;
  paymentStatus: string;
  transactionId?: string;
  orderNumber?: string;
  billingInfo?: BillingInfo;
  shippingInfo?: ShippingInfo;
  paymentInfo?: PaymentInfo;
  notes?: string;
  totals?: {
    subtotal: number;
    shipping: number;
    discount: number;
    tax: number;
    total: number;
  };
}

// API response interface
interface ApiResponse<T> {
  error?: boolean;
  message?: string;
  data?: T | DataWrapper<T>; // Data can be direct array or wrapped
  orders?: T; // Some APIs might use 'orders' instead of 'data'
  id?: string; // Allow for the case when response data is a single order
  [key: string]: any; // Allow additional properties
}

// Data wrapper (for nested data structure)
interface DataWrapper<T> {
  limit?: number;
  page?: number;
  total?: number;
  orders?: T;
  [key: string]: any;
}

// Backend order interface based on what's coming from the API
interface BackendOrder {
  id: string;
  userId?: string;
  status: string;
  items: any[];
  totals?: {
    subtotal: number;
    shipping: number;
    discount: number;
    tax: number;
    total: number;
  };
  billingInfo?: any; // camelCase alternative
  shippingInfo?: any; // camelCase alternative
  paymentInfo?: any; // camelCase alternative
  createdAt?: string;
  updatedAt?: string;
  created_at?: string; // Backend might be using snake_case
  updated_at?: string; // Backend might be using snake_case
  orderNumber?: string;
  order_number?: string; // Backend might be using snake_case
}

/**
 * Adapt backend order format to frontend order format
 */
const adaptOrder = (backendOrder: BackendOrder): Order => {
  // console.log('Adapting backend order to frontend format:', backendOrder);

  // Extract payment info - this is the key issue
  const paymentInfo = backendOrder.paymentInfo || {};

  return {
    id: backendOrder.id,
    userId: backendOrder.userId || '',
    status: backendOrder.status as OrderStatus,
    items: Array.isArray(backendOrder.items)
      ? backendOrder.items.map(item => ({
        id: item.id || '',
        productId: item.product_id || item.productId || '',
        name: item.name || 'Unknown Product',
        price: typeof item.price === 'number' ? item.price : 0,
        quantity: typeof item.quantity === 'number' ? item.quantity : 1,
        imageUrl: item.image_url || item.imageUrl
      }))
      : [],
    subtotal: backendOrder.totals?.subtotal || 0,
    shipping: backendOrder.totals?.shipping || 0,
    discount: backendOrder.totals?.discount || 0,
    tax: backendOrder.totals?.tax || 0,
    total: backendOrder.totals?.total || 0,
    createdAt: backendOrder.createdAt || backendOrder.created_at || new Date().toISOString(),
    updatedAt: backendOrder.updatedAt || backendOrder.updated_at || new Date().toISOString(),
    // Sửa lại cách lấy paymentMethod
    paymentMethod: paymentInfo.paymentMethod || '',
    paymentStatus: paymentInfo.status || 'pending',
    transactionId: paymentInfo.transactionId || '',
    orderNumber: backendOrder.orderNumber || '',
    billingInfo: backendOrder.billingInfo as BillingInfo,
    shippingInfo: backendOrder.shippingInfo as ShippingInfo,
    paymentInfo: backendOrder.paymentInfo as PaymentInfo,
    totals: backendOrder.totals,
  };
};

class OrderService {
  private baseUrl = `${API_URL}/checkout/orders`;

  /**
   * Fetch all orders for the authenticated user
   */
  async getOrders(): Promise<Order[]> {
    try {
      const headers = getAuthHeader();

      console.log('----------- ORDER SERVICE DEBUG -----------');
      console.log('Fetching orders from:', this.baseUrl);
      console.log('Authorization token present:', !!headers.Authorization);

      const response = await axios.get<ApiResponse<BackendOrder[]>>(this.baseUrl, { headers });

      console.log('API Response received, status:', response.status);
      console.log('Response headers:', response.headers);
      console.log('Raw response data type:', typeof response.data);
      console.log('Raw response data:', JSON.stringify(response.data, null, 2));

      // Process different response formats and convert to frontend order format
      let backendOrders: BackendOrder[] = [];

      // Check different possible formats
      if (Array.isArray(response.data)) {
        console.log('Response is an array with', response.data.length, 'orders');
        backendOrders = response.data;
      } else if (typeof response.data === 'object' && response.data !== null) {
        console.log('Response is an object with keys:', Object.keys(response.data));

        // Check for data property
        if (response.data.data !== undefined) {
          // console.log('Data property exists, type:', typeof response.data.data);

          if (Array.isArray(response.data.data)) {
            // console.log('Data property is an array with', response.data.data.length, 'items');
            backendOrders = response.data.data;
          } else if (typeof response.data.data === 'object' && response.data.data !== null) {
            // console.log('Data property is an object with keys:', Object.keys(response.data.data));

            // Check for orders array inside data object (handle nested structure)
            const dataWrapper = response.data.data as DataWrapper<BackendOrder[]>;
            if (Array.isArray(dataWrapper.orders)) {
              // console.log('Found orders array inside data object with', dataWrapper.orders.length, 'orders');
              backendOrders = dataWrapper.orders;
            }
          } else {
            console.log('Data property is NOT an array:', response.data.data);
          }
        }

        // Check for orders property directly on response.data
        if (response.data.orders !== undefined) {
          console.log('Orders property exists, type:', typeof response.data.orders);

          if (Array.isArray(response.data.orders)) {
            console.log('Orders property is an array with', response.data.orders.length, 'items');
            backendOrders = response.data.orders;
          } else {
            console.log('Orders property is NOT an array:', response.data.orders);
          }
        }

        // If response.data itself contains order properties, it might be a single order
        if (response.data.id && typeof response.data.id === 'string') {
          console.log('Found a single order with ID:', response.data.id);
          console.log('Single order details:', response.data);
          backendOrders = [response.data as BackendOrder];
        }
      }

      // Try to adapt to unexpected response format as a fallback
      if (backendOrders.length === 0 && typeof response.data === 'object' && !Array.isArray(response.data)) {
        const possibleOrdersArray = Object.values(response.data);
        console.log('Trying to extract orders from object values:', possibleOrdersArray);

        for (const value of possibleOrdersArray) {
          if (Array.isArray(value)) {
            console.log('Found array in object values with length:', value.length);
            backendOrders = value as BackendOrder[];
            break;
          } else if (typeof value === 'object' && value !== null && 'orders' in value && Array.isArray(value.orders)) {
            console.log('Found nested orders array in object values with length:', value.orders.length);
            backendOrders = value.orders as BackendOrder[];
            break;
          }
        }
      }

      // Convert backend orders to frontend format
      if (backendOrders.length > 0) {
        // console.log(`Converting ${backendOrders.length} backend orders to frontend format`);
        const frontendOrders = backendOrders.map(adaptOrder);
        //  console.log('Converted orders:', frontendOrders);
        return frontendOrders;
      }

      console.warn('Could not extract orders from response - returning empty array');
      return [];
    } catch (error) {
      console.error('----------- ORDER SERVICE ERROR -----------');
      console.error('Error fetching orders:', error);

      if (axios.isAxiosError(error)) {
        console.error('Axios error details:');
        console.error('Status:', error.response?.status);
        console.error('Status text:', error.response?.statusText);
        console.error('Response data:', error.response?.data);
        console.error('Request URL:', error.config?.url);
        console.error('Request method:', error.config?.method);
        console.error('Request headers:', error.config?.headers);
      }

      this.handleError('Error fetching orders', error);
      return [];
    }
  }

  /**
   * Fetch a specific order by ID
   */
  async getOrderById(orderId: string): Promise<Order | null> {
    try {
      const headers = getAuthHeader();

      const response = await axios.get<ApiResponse<BackendOrder>>(`${this.baseUrl}/${orderId}`, { headers });

      // Check if the response has data property
      if (response.data && response.data.data) {
        return adaptOrder(response.data.data as BackendOrder);
      }

      // If response.data itself contains order properties
      if (response.data && response.data.id) {
        return adaptOrder(response.data as unknown as BackendOrder);
      }

      return null;
    } catch (error) {
      console.error(`Error fetching order ${orderId}:`, error);
      this.handleError('Error fetching order details', error);
      return null;
    }
  }

  // Handle API errors
  private handleError(message: string, error: any): void {
    console.error(`${message}:`, error);

    if (axios.isAxiosError(error)) {
      if (error.response?.status === 401) {
        throw new Error('Vui lòng đăng nhập để xem đơn hàng');
      } else if (error.response?.status === 404) {
        throw new Error('Không tìm thấy thông tin đơn hàng');
      } else if (error.response?.data?.message) {
        throw new Error(error.response.data.message);
      } else if (error.response) {
        throw new Error(`Lỗi từ máy chủ: ${error.response.status} - ${error.response.statusText}`);
      } else if (error.request) {
        throw new Error('Không thể kết nối đến máy chủ. Vui lòng kiểm tra kết nối mạng.');
      }
    }

    throw new Error('Đã xảy ra lỗi khi truy xuất thông tin đơn hàng');
  }
}

// Export singleton instance
const orderService = new OrderService();
export default orderService; 