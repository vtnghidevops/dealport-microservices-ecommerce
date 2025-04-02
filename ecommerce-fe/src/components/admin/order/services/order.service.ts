// src/services/order.service.ts
import { Order, OrderStatus, OrderFilterParams, OrderSummary } from '../models/order.model';

// Mock data for frontend testing
const mockOrders: Order[] = [
  {
    id: '1',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD001',
        productName: 'Wireless Bluetooth Headphones',
        productImage: '/images/products/headphones.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 49.99,
    status: 'Delivered',
    paymentStatus: 'Paid',
  },
  {
    id: '2',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD002',
        productName: "Men's T-Shirt",
        productImage: '/images/products/tshirt.png',
        quantity: 1,
        price: 14.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 14.99,
    status: 'Pending',
    paymentStatus: 'Unpaid',
  },
  {
    id: '3',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD003',
        productName: "Men's Leather Wallet",
        productImage: '/images/products/wallet.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 49.99,
    status: 'Delivered',
    paymentStatus: 'Paid',
  },
  {
    id: '4',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD004',
        productName: "Memory Foam Pillow",
        productImage: '/images/products/pillow.png',
        quantity: 1,
        price: 39.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 39.99,
    status: 'Shipped',
    paymentStatus: 'Paid',
  },
  {
    id: '5',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD005',
        productName: "Adjustable Dumbbells",
        productImage: '/images/products/dumbbells.png',
        quantity: 1,
        price: 14.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 14.99,
    status: 'Pending',
    paymentStatus: 'Unpaid',
  },
  {
    id: '6',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD006',
        productName: "Coffee Maker",
        productImage: '/images/products/coffee-maker.png',
        quantity: 1,
        price: 79.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 79.99,
    status: 'Cancelled',
    paymentStatus: 'Unpaid',
  },
  {
    id: '7',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD007',
        productName: "Casual Baseball Cap",
        productImage: '/images/products/cap.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 49.99,
    status: 'Delivered',
    paymentStatus: 'Paid',
  },
  {
    id: '8',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD007',
        productName: "Casual Baseball Cap",
        productImage: '/images/products/cap.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 49.99,
    status: 'Delivered',
    paymentStatus: 'Paid',
  },
  {
    id: '9',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD007',
        productName: "Casual Baseball Cap",
        productImage: '/images/products/cap.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 49.99,
    status: 'Delivered',
    paymentStatus: 'Paid',
  },
  {
    id: '10',
    orderId: 'ORD0001',
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD007',
        productName: "Casual Baseball Cap",
        productImage: '/images/products/cap.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    date: '01-01-2025',
    totalAmount: 49.99,
    status: 'Delivered',
    paymentStatus: 'Paid',
  },
];

// Mock order summary data
const mockOrderSummary: OrderSummary = {
  totalOrders: 324,
  newOrders: 42,
  completedOrders: 265,
  cancelledOrders: 17,
  lastUpdated: 'Last 7 days',
  growthRate: {
    total: 12.5,
    new: 8.3,
    completed: 15.2,
    cancelled: -2.1
  }
};

class OrderService {
  /**
   * Fetch orders with filtering and pagination
   */
  async fetchOrders(params: OrderFilterParams): Promise<{ orders: Order[], total: number }> {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 500));
    
    let filteredOrders = [...mockOrders];
    
    // Apply filters if provided
    if (params) {
      if (params.status) {
        filteredOrders = filteredOrders.filter(order => order.status === params.status);
      }
      
      if (params.customerId) {
        filteredOrders = filteredOrders.filter(order => order.customerId === params.customerId);
      }
      
      if (params.productId) {
        filteredOrders = filteredOrders.filter(order => 
          order.products.some(product => product.productId === params.productId)
        );
      }
      
      if (params.searchTerm) {
        const searchLower = params.searchTerm.toLowerCase();
        filteredOrders = filteredOrders.filter(order => 
          order.orderId.toLowerCase().includes(searchLower) ||
          order.customerName.toLowerCase().includes(searchLower) ||
          order.products.some(product => product.productName.toLowerCase().includes(searchLower))
        );
      }
      
      if (params.dateRange) {
        const [startDate, endDate] = params.dateRange;
        filteredOrders = filteredOrders.filter(order => {
          const orderDate = new Date(order.date);
          return orderDate >= startDate && orderDate <= endDate;
        });
      }
    }
    
    // Get total count before pagination
    const total = filteredOrders.length;
    
    // Apply pagination
    const startIndex = (params.page - 1) * params.limit;
    const endIndex = startIndex + params.limit;
    filteredOrders = filteredOrders.slice(startIndex, endIndex);
    
    return { orders: filteredOrders, total };
  }
  
  /**
   * Fetch order summary statistics
   */
  async fetchOrderSummary(): Promise<OrderSummary> {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 300));
    
    return mockOrderSummary;
  }
  
  /**
   * Update an order's status
   */
  async updateOrderStatus(orderId: string, status: OrderStatus): Promise<boolean> {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 400));
    
    const orderIndex = mockOrders.findIndex(order => order.id === orderId);
    if (orderIndex !== -1) {
      mockOrders[orderIndex].status = status;
      return true;
    }
    return false;
  }
  
  /**
   * Update payment status
   */
  async updatePaymentStatus(orderId: string, paymentStatus: 'Paid' | 'Unpaid'): Promise<boolean> {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 400));
    
    const orderIndex = mockOrders.findIndex(order => order.id === orderId);
    if (orderIndex !== -1) {
      mockOrders[orderIndex].paymentStatus = paymentStatus;
      return true;
    }
    return false;
  }

  async createOrder (orderData: Omit<Order, 'id' | 'orderId' | 'status' | 'date'>): Promise<Order> {
    try {
      // Here you would call your API endpoint to create an order
      // For now, we'll simulate an API call with a timeout
      return new Promise(resolve => {
        setTimeout(() => {
          resolve({
            id: Date.now().toString(),
            orderId: `ORD${Math.floor(Math.random() * 10000).toString().padStart(4, '0')}`,
            ...orderData,
            status: 'Pending',
            date: new Date().toLocaleDateString('en-US', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric'
            }).replace(/\//g, '-'),
          });
        }, 500);
      });
    } catch (error) {
      console.error('Error creating order:', error);
      throw error;
    }
  };
}

// Export a singleton instance
export const orderService = new OrderService();