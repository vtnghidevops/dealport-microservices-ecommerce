import {
  Customer,
  CustomerOverview,
  CustomerStatus,
  CustomerChartData,
  CustomerFilterParams
} from '../models/customer.model';
import axios from 'axios';
import { getApiUrl, getAuthHeader } from '@/utils/api-config';

// Helper function to convert API data to Customer object
const mapApiDataToCustomer = (userData: any): Customer => {
  return {
    id: userData.id?.toString() || '',
    name: userData.name || 'Unknown User',
    phone: userData.phone || 'N/A',
    orderCount: userData.orderCount || 0,
    totalSpend: userData.totalSpend || 0,
    status: userData.status || CustomerStatus.ACTIVE
  };
};

export const CustomerService = {
  getCustomers: async (params: CustomerFilterParams = { page: 1, limit: 10 }): Promise<{ customers: Customer[], total: number }> => {
    try {
      const headers = getAuthHeader();

      // Build filter parameters based on CustomerFilterParams
      const requestData: any = {
        page: params.page,
        limit: params.limit
      };

      // Add search term if provided
      if (params.searchTerm) {
        requestData.search = params.searchTerm;
      }

      // Add status filter if provided
      if (params.status) {
        requestData.status = params.status;
      }

      // Call the customer list API (sử dụng endpoint admin mới)
      const response = await axios.post(
        getApiUrl('users/admin/list'),
        requestData,
        { headers }
      );

      if (response.data && !response.data.error && response.data.data) {
        const { customers, total } = response.data.data;

        // Map API data to Customer model
        const mappedCustomers = customers.map(mapApiDataToCustomer);

        return {
          customers: mappedCustomers,
          total: total || mappedCustomers.length
        };
      } else {
        throw new Error(response.data?.message || 'Failed to fetch customers');
      }
    } catch (error) {
      console.error('Error fetching customers:', error);
      // Throw the error to be handled by the component
      throw error;
    }
  },

  getCustomerOverview: async (): Promise<CustomerOverview> => {
    try {
      const headers = getAuthHeader();

      // Call the customer statistics API (sử dụng endpoint admin mới)
      const response = await axios.post(
        getApiUrl('users/admin/statistics'),
        {},
        { headers }
      );

      if (response.data && !response.data.error && response.data.data) {
        // Trích xuất dữ liệu thống kê từ API response
        const stats = response.data.data;

        // Cấu trúc object đúng theo định dạng CustomerOverview
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
        throw new Error(response.data?.message || 'Failed to fetch customer overview');
      }
    } catch (error) {
      console.error('Error fetching customer overview:', error);
      // Throw the error to be handled by the component
      throw error;
    }
  },

  getCustomerChartData: async (): Promise<CustomerChartData[]> => {
    try {
      const headers = getAuthHeader();

      // Call the customer activity chart API (sử dụng endpoint admin mới)
      const response = await axios.post(
        getApiUrl('users/admin/activity-chart'),
        {},
        { headers }
      );

      if (response.data && !response.data.error && response.data.data) {
        const chartData = response.data.data.chart_data || [];

        // Map API data to CustomerChartData format
        return chartData.map((item: any) => ({
          day: item.day || 'Unknown',
          count: item.count || 0
        }));
      } else {
        throw new Error(response.data?.message || 'Failed to fetch customer chart data');
      }
    } catch (error) {
      console.error('Error fetching customer chart data:', error);
      // Throw the error to be handled by the component
      throw error;
    }
  },

  updateCustomerStatus: async (customerId: string, status: CustomerStatus): Promise<boolean> => {
    try {
      const headers = getAuthHeader();

      // Prepare update data based on status
      const updateData: any = {
        status: status
      };

      // Gọi API cập nhật trạng thái người dùng
      const response = await axios.put(
        getApiUrl(`users/${customerId}/status`),
        updateData,
        { headers }
      );

      return response.data && !response.data.error;
    } catch (error) {
      console.error('Error updating customer status:', error);
      return false;
    }
  },

  deleteCustomer: async (customerId: string): Promise<boolean> => {
    try {
      const headers = getAuthHeader();

      // Gọi API xóa người dùng
      const response = await axios.delete(
        getApiUrl(`users/${customerId}`),
        { headers }
      );

      return response.data && !response.data.error;
    } catch (error) {
      console.error('Error deleting customer:', error);
      return false;
    }
  }
};