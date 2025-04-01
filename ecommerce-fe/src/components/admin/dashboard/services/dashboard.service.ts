import { DashboardSummary } from '../models/dashboard.model';

export const DashboardService = {
  getDashboardSummary: async (): Promise<DashboardSummary> => {
    // Giả lập API call
    await new Promise(resolve => setTimeout(resolve, 500));
    
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
};
