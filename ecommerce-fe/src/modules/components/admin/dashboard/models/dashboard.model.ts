export interface DashboardSummary {
  totalSales: {
    amount: number;
    currency: string;
    percentChange: number;
    previousAmount: number;
    lastDays: number;
  };
  totalOrders: {
    count: number;
    percentChange: number;
    previousCount: number;
    lastDays: number;
  };
  pendingCanceled: {
    pending: {
      count: number;
      userCount: number;
    };
    canceled: {
      count: number;
      percentChange: number;
    };
    lastDays: number;
  };
  weeklyReport: {
    customers: number;
    totalProducts: number;
    stockProducts: number;
    outOfStock: number;
    revenue: number;
    chartData: WeeklyChartData;
  };
  userStats: {
    totalInLastMinutes: number;
    minutesDuration: number;
    userPerMinute: UserPerMinuteData[];
  };
  salesByCountry: CountrySales[];
}

export interface WeeklyChartData {
  days: string[];
  values: number[];
  highlights?: {
    day: string;
    value: number;
    label: string;
  }[];
}

export interface UserPerMinuteData {
  minute: string;
  count: number;
}

export interface CountrySales {
  country: string;
  code: string;
  amount: number;
  percentChange: number;
}
