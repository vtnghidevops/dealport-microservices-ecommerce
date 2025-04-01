import { 
  Customer, 
  CustomerOverview, 
  CustomerStatus, 
  CustomerChartData,
  CustomerFilterParams 
} from '../models/customer.model';

// Mock data for development
const mockCustomers: Customer[] = [
  { id: 'CUST001', name: 'John Doe', phone: '+1234567890', orderCount: 25, totalSpend: 3450.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST002', name: 'Jane Smith', phone: '+1987654321', orderCount: 5, totalSpend: 250.00, status: CustomerStatus.INACTIVE },
  { id: 'CUST003', name: 'Emily Davis', phone: '+1122334455', orderCount: 30, totalSpend: 4600.00, status: CustomerStatus.VIP },
  { id: 'CUST004', name: 'Michael Brown', phone: '+1555666777', orderCount: 12, totalSpend: 1200.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST005', name: 'Sarah Wilson', phone: '+1888999000', orderCount: 8, totalSpend: 780.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST006', name: 'Robert Taylor', phone: '+1222333444', orderCount: 3, totalSpend: 150.00, status: CustomerStatus.INACTIVE },
  { id: 'CUST007', name: 'Jessica Martinez', phone: '+1777888999', orderCount: 18, totalSpend: 2100.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST008', name: 'David Anderson', phone: '+1444555666', orderCount: 7, totalSpend: 520.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST009', name: 'Christopher Lee', phone: '+1333222111', orderCount: 22, totalSpend: 3200.00, status: CustomerStatus.VIP },
  { id: 'CUST010', name: 'Amanda Thomas', phone: '+1666777888', orderCount: 4, totalSpend: 320.00, status: CustomerStatus.INACTIVE },
  { id: 'CUST011', name: 'Ryan Garcia', phone: '+1999888777', orderCount: 9, totalSpend: 870.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST012', name: 'Jennifer White', phone: '+1111222333', orderCount: 15, totalSpend: 1800.00, status: CustomerStatus.ACTIVE },
  { id: 'CUST013', name: 'Kevin Harris', phone: '+1888777666', orderCount: 6, totalSpend: 450.00, status: CustomerStatus.INACTIVE },
  { id: 'CUST014', name: 'Lisa Robinson', phone: '+1444333222', orderCount: 28, totalSpend: 3900.00, status: CustomerStatus.VIP },
  { id: 'CUST015', name: 'Daniel King', phone: '+1777666555', orderCount: 11, totalSpend: 950.00, status: CustomerStatus.ACTIVE },
];

const mockOverview: CustomerOverview = {
  totalCustomers: {
    count: 11040,
    growth: 14.4,
    period: 'Last 7 days'
  },
  newCustomers: {
    count: 2370,
    growth: 20,
    period: 'Last 7 days'
  },
  visitors: {
    count: 250000,
    growth: 20,
    period: 'Last 7 days'
  },
  activeCustomers: {
    count: 25000,
    chartLabel: 'Active Customers'
  },
  repeatCustomers: {
    count: 5600,
    chartLabel: 'Repeat Customers'
  },
  shopVisitor: {
    count: 250000,
    chartLabel: 'Shop Visitor'
  },
  conversionRate: {
    rate: 5.5,
    chartLabel: 'Conversion Rate'
  }
};

const mockChartData: CustomerChartData[] = [
  { day: 'Sun', count: 12000 },
  { day: 'Mon', count: 19000 },
  { day: 'Tue', count: 17000 },
  { day: 'Wed', count: 22000 },
  { day: 'Thu', count: 25000 },
  { day: 'Fri', count: 23000 },
  { day: 'Sat', count: 18000 },
];

export const CustomerService = {
  getCustomers: async (params: CustomerFilterParams = { page: 1, limit: 10 }): Promise<{ customers: Customer[], total: number }> => {
    // In a real app, this would be an API call with filtering, pagination, etc.
    let filteredCustomers = [...mockCustomers];
    
    // Apply status filter
    if (params.status) {
      filteredCustomers = filteredCustomers.filter(customer => customer.status === params.status);
    }
    
    // Apply search filter
    if (params.searchTerm) {
      const searchLower = params.searchTerm.toLowerCase();
      filteredCustomers = filteredCustomers.filter(customer => 
        customer.name.toLowerCase().includes(searchLower) ||
        customer.id.toLowerCase().includes(searchLower)
      );
    }
    
    // Calculate pagination
    const total = filteredCustomers.length;
    const start = (params.page - 1) * params.limit;
    const end = start + params.limit;
    const paginatedCustomers = filteredCustomers.slice(start, end);
    
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 300));
    
    return {
      customers: paginatedCustomers,
      total
    };
  },

  getCustomerOverview: async (): Promise<CustomerOverview> => {
    // In a real app, this would be an API call
    await new Promise(resolve => setTimeout(resolve, 300));
    return mockOverview;
  },

  getCustomerChartData: async (): Promise<CustomerChartData[]> => {
    // In a real app, this would be an API call
    await new Promise(resolve => setTimeout(resolve, 300));
    return mockChartData;
  },

  updateCustomerStatus: async (customerId: string, status: CustomerStatus): Promise<boolean> => {
    // In a real app, this would be an API call
    console.log(`Updating customer ${customerId} status to ${status}`);
    await new Promise(resolve => setTimeout(resolve, 300));
    return true;
  }
};