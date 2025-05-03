import {
  Customer,
  CustomerOverview,
  CustomerStatus,
  CustomerChartData,
  CustomerFilterParams
} from '../models/customer.model';
import axios from 'axios';

// Base URL for API requests
const API_BASE_URL = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8080/api/v1';
const BROKER_ENDPOINT = `${API_BASE_URL}/broker`;

// Create an axios instance with authorization configuration
const getAuthClient = () => {
  const token = localStorage.getItem('token');
  return axios.create({
    baseURL: API_BASE_URL,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
};

// Fallback mock data for cases where API is unavailable
const fallbackMockData = {
  customers: [
    { id: 'CUST001', name: 'John Doe', phone: '+1234567890', orderCount: 25, totalSpend: 3450.00, status: CustomerStatus.ACTIVE },
    { id: 'CUST002', name: 'Jane Smith', phone: '+1987654321', orderCount: 5, totalSpend: 250.00, status: CustomerStatus.INACTIVE },
    { id: 'CUST003', name: 'Emily Davis', phone: '+1122334455', orderCount: 30, totalSpend: 4600.00, status: CustomerStatus.VIP },
  ],
  overview: {
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
  },
  chartData: [
    { day: 'Sun', count: 12000 },
    { day: 'Mon', count: 19000 },
    { day: 'Tue', count: 17000 },
    { day: 'Wed', count: 22000 },
    { day: 'Thu', count: 25000 },
    { day: 'Fri', count: 23000 },
    { day: 'Sat', count: 18000 },
  ],
};

// Helper function to convert API data to Customer object
const mapApiDataToCustomer = (userData: any): Customer => {
  return {
    id: userData.id?.toString() || '',
    name: `${userData.first_name || ''} ${userData.last_name || ''}`.trim() || userData.email || 'Unknown User',
    phone: userData.phone || 'N/A',
    orderCount: userData.order_count || 0,
    totalSpend: userData.total_spend || 0,
    status: userData.vip ? CustomerStatus.VIP :
      userData.active ? CustomerStatus.ACTIVE : CustomerStatus.INACTIVE
  };
};

export const CustomerService = {
  getCustomers: async (params: CustomerFilterParams = { page: 1, limit: 10 }): Promise<{ customers: Customer[], total: number }> => {
    try {
      const api = getAuthClient();

      // Building query parameters
      const queryParams = new URLSearchParams();
      queryParams.append('page', params.page.toString());
      queryParams.append('limit', params.limit.toString());

      if (params.searchTerm) {
        queryParams.append('search', params.searchTerm);
      }

      if (params.status) {
        switch (params.status) {
          case CustomerStatus.ACTIVE:
            queryParams.append('active', 'true');
            break;
          case CustomerStatus.INACTIVE:
            queryParams.append('active', 'false');
            break;
          case CustomerStatus.VIP:
            queryParams.append('vip', 'true');
            break;
        }
      }

      // Call broker service to get user data from user-service
      const response = await api.post(`${BROKER_ENDPOINT}/user-service/users/list`, {
        params: queryParams.toString()
      });

      if (response.data && response.data.data) {
        const users = response.data.data;
        const total = response.data.meta?.total || users.length;

        // Map API data to Customer model
        const customers = users.map(mapApiDataToCustomer);

        return {
          customers,
          total
        };
      } else {
        throw new Error('Invalid API response format');
      }
    } catch (error) {
      console.error('Error fetching customers:', error);

      // Fallback to mock data if API fails
      const filteredCustomers = fallbackMockData.customers.filter(customer => {
        if (params.status && customer.status !== params.status) {
          return false;
        }

        if (params.searchTerm) {
          const searchLower = params.searchTerm.toLowerCase();
          return customer.name.toLowerCase().includes(searchLower) ||
            customer.id.toLowerCase().includes(searchLower);
        }

        return true;
      });

      const total = filteredCustomers.length;
      const start = (params.page - 1) * params.limit;
      const end = start + params.limit;
      const paginatedCustomers = filteredCustomers.slice(start, end);

      return {
        customers: paginatedCustomers,
        total
      };
    }
  },

  getCustomerOverview: async (): Promise<CustomerOverview> => {
    try {
      const api = getAuthClient();

      // Call broker service to get customer statistics from user-service
      const response = await api.post(`${BROKER_ENDPOINT}/user-service/users/statistics`);

      if (response.data && response.data.data) {
        const stats = response.data.data;

        return {
          totalCustomers: {
            count: stats.total_users || 0,
            growth: stats.user_growth || 0,
            period: 'Last 7 days'
          },
          newCustomers: {
            count: stats.new_users || 0,
            growth: stats.new_user_growth || 0,
            period: 'Last 7 days'
          },
          visitors: {
            count: stats.visitors || 0,
            growth: stats.visitor_growth || 0,
            period: 'Last 7 days'
          },
          activeCustomers: {
            count: stats.active_users || 0,
            chartLabel: 'Active Customers'
          },
          repeatCustomers: {
            count: stats.repeat_customers || 0,
            chartLabel: 'Repeat Customers'
          },
          shopVisitor: {
            count: stats.shop_visitors || 0,
            chartLabel: 'Shop Visitor'
          },
          conversionRate: {
            rate: stats.conversion_rate || 0,
            chartLabel: 'Conversion Rate'
          }
        };
      } else {
        throw new Error('Invalid API response format');
      }
    } catch (error) {
      console.error('Error fetching customer overview:', error);

      // Fallback to mock data if API fails
      return fallbackMockData.overview;
    }
  },

  getCustomerChartData: async (): Promise<CustomerChartData[]> => {
    try {
      const api = getAuthClient();

      // Call broker service to get chart data from user-service
      const response = await api.post(`${BROKER_ENDPOINT}/user-service/users/activity-chart`);

      if (response.data && response.data.data) {
        // Map API response to CustomerChartData format
        return response.data.data.map((item: any) => ({
          day: item.day || item.date || 'Unknown',
          count: item.count || item.value || 0
        }));
      } else {
        throw new Error('Invalid API response format');
      }
    } catch (error) {
      console.error('Error fetching customer chart data:', error);

      // Fallback to mock data if API fails
      return fallbackMockData.chartData;
    }
  },

  updateCustomerStatus: async (customerId: string, status: CustomerStatus): Promise<boolean> => {
    try {
      const api = getAuthClient();

      // Prepare update data based on status
      const updateData: any = {};

      switch (status) {
        case CustomerStatus.ACTIVE:
          updateData.active = true;
          updateData.vip = false;
          break;
        case CustomerStatus.INACTIVE:
          updateData.active = false;
          updateData.vip = false;
          break;
        case CustomerStatus.VIP:
          updateData.active = true;
          updateData.vip = true;
          break;
      }

      // Call broker service to update user status in user-service
      const response = await api.post(`${BROKER_ENDPOINT}/user-service/users/${customerId}/update-status`, updateData);

      return response.status === 200 || response.status === 204;
    } catch (error) {
      console.error(`Error updating customer ${customerId} status:`, error);
      return false;
    }
  },

  deleteCustomer: async (customerId: string): Promise<boolean> => {
    try {
      const api = getAuthClient();

      // Call broker service to delete user in user-service
      const response = await api.post(`${BROKER_ENDPOINT}/user-service/users/${customerId}/delete`);

      return response.status === 200 || response.status === 204;
    } catch (error) {
      console.error(`Error deleting customer ${customerId}:`, error);
      return false;
    }
  }
};