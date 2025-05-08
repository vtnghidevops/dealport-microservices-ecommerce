import { DashboardSummary } from '../models/dashboard.model';
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_PUBLIC_BROKER_VITE_PUBLIC_BROKER_API_URL || 'http://localhost:8080/api/v1';

interface ChartDataPoint {
  day: string;
  date?: string;
  count: number;
  value?: number;
}

export const DashboardService = {
  getDashboardSummary: async (): Promise<DashboardSummary> => {
    try {
      // Parallel requests to backend services for statistics data
      const [userStats, orderStats, productStats] = await Promise.all([
        axios.get(`${API_BASE_URL}/users/statistics`).then(res => res.data.data),
        axios.get(`${API_BASE_URL}/orders/statistics`).then(res => res.data.data).catch(() => null),
        axios.get(`${API_BASE_URL}/products/statistics`).then(res => res.data.data).catch(() => null)
      ]);

      console.log('Dashboard data loaded from backend:', { userStats, orderStats, productStats });

      // User activity data
      const userActivity = await axios.get(`${API_BASE_URL}/users/activity-chart?days=7&chart_type=activity`)
        .then(res => res.data.data)
        .catch(() => ({ chart_data: [] }));

      // Order data by date for weekly report
      const orderChart = await axios.get(`${API_BASE_URL}/orders/activity-chart?days=7&chart_type=orders`)
        .then(res => res.data.data)
        .catch(() => ({
          chart_data: Array.from({ length: 7 }, (_, i) => ({
            day: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'][i],
            count: Math.floor(Math.random() * 50) + 10,
          }))
        }));

      // Map backend data to dashboard model
      return {
        totalSales: {
          amount: orderStats?.totalSales || 350000,
          currency: '$',
          percentChange: orderStats?.salesGrowth || 10.4,
          previousAmount: orderStats?.previousPeriodSales || 235000,
          lastDays: 7
        },
        totalOrders: {
          count: orderStats?.totalOrders || 10700,
          percentChange: orderStats?.orderGrowth || 14.4,
          previousCount: orderStats?.previousPeriodOrders || 7600,
          lastDays: 7
        },
        pendingCanceled: {
          pending: {
            count: orderStats?.pendingOrders || 509,
            userCount: orderStats?.uniqueCustomersWithPendingOrders || 204
          },
          canceled: {
            count: orderStats?.canceledOrders || 94,
            percentChange: orderStats?.cancelationRateChange || 14.4
          },
          lastDays: 7
        },
        weeklyReport: {
          customers: userStats?.totalUsers || 52000,
          totalProducts: productStats?.totalProducts || 3500,
          stockProducts: productStats?.inStockProducts || 2500,
          outOfStock: productStats?.outOfStockProducts || 500,
          revenue: orderStats?.totalSales || 250000,
          chartData: {
            days: orderChart.chart_data.map((point: ChartDataPoint) => point.day),
            values: orderChart.chart_data.map((point: ChartDataPoint) => point.count),
            highlights: [
              {
                day: orderChart.chart_data.reduce((max: number, point: ChartDataPoint, index: number) =>
                  point.count > orderChart.chart_data[max].count ? index : max, 0) >= 0
                  ? orderChart.chart_data[orderChart.chart_data.reduce((max: number, point: ChartDataPoint, index: number) =>
                    point.count > orderChart.chart_data[max].count ? index : max, 0)].day
                  : 'Thu',
                value: Math.max(...orderChart.chart_data.map((p: ChartDataPoint) => p.count)),
                label: (Math.max(...orderChart.chart_data.map((p: ChartDataPoint) => p.count)) / 1000).toFixed(1) + 'k'
              }
            ]
          }
        },
        userStats: {
          totalInLastMinutes: userStats?.activeUsers || 21500,
          minutesDuration: 30,
          userPerMinute: userActivity.chart_data?.length
            ? userActivity.chart_data.map((point: ChartDataPoint, i: number) => ({
              minute: `${i}`,
              count: point.count
            }))
            : Array.from({ length: 24 }, (_, i) => ({
              minute: `${i}`,
              count: Math.floor(Math.random() * 50) + 10
            }))
        },
        salesByCountry: [
          {
            country: 'US',
            code: 'us',
            amount: 30000,
            percentChange: 25.8
          },
          {
            country: 'Brazil',
            code: 'br',
            amount: 30000,
            percentChange: -18.8
          },
          {
            country: 'Australia',
            code: 'au',
            amount: 25000,
            percentChange: 35.8
          }
        ]
      };
    } catch (error) {
      console.error('Error fetching dashboard data:', error);

      // Return fallback data if API requests fail
      return {
        totalSales: {
          amount: 350000,
          currency: '$',
          percentChange: 10.4,
          previousAmount: 235000,
          lastDays: 7
        },
        totalOrders: {
          count: 10700,
          percentChange: 14.4,
          previousCount: 7600,
          lastDays: 7
        },
        pendingCanceled: {
          pending: {
            count: 509,
            userCount: 204
          },
          canceled: {
            count: 94,
            percentChange: 14.4
          },
          lastDays: 7
        },
        weeklyReport: {
          customers: 52000,
          totalProducts: 3500,
          stockProducts: 2500,
          outOfStock: 500,
          revenue: 250000,
          chartData: {
            days: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
            values: [10, 20, 30, 45, 50, 40, 30],
            highlights: [
              {
                day: 'Thu',
                value: 50,
                label: '14k'
              }
            ]
          }
        },
        userStats: {
          totalInLastMinutes: 21500,
          minutesDuration: 30,
          userPerMinute: Array.from({ length: 24 }, (_, i) => ({
            minute: `${i}`,
            count: Math.floor(Math.random() * 50) + 10
          }))
        },
        salesByCountry: [
          {
            country: 'US',
            code: 'us',
            amount: 30000,
            percentChange: 25.8
          },
          {
            country: 'Brazil',
            code: 'br',
            amount: 30000,
            percentChange: -18.8
          },
          {
            country: 'Australia',
            code: 'au',
            amount: 25000,
            percentChange: 35.8
          }
        ]
      };
    }
  }
};
