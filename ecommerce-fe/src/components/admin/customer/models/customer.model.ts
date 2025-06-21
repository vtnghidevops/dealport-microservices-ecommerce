export enum CustomerStatus {
  ACTIVE = "Active",
  INACTIVE = "Inactive",
  VIP = "VIP",
}

export interface Customer {
  id: string;
  name: string;
  phone: string;
  orderCount: number;
  totalSpend: number;
  status: CustomerStatus;
}

export interface CustomerFilterParams {
  page: number;
  limit: number;
  searchTerm?: string;
  status?: CustomerStatus;
}

export interface CustomerOverview {
  totalCustomers: {
    count: number;
    growth: number;
    period: string;
  };
  newCustomers: {
    count: number;
    growth: number;
    period: string;
  };
  visitors: {
    count: number;
    growth: number;
    period: string;
  };
  activeCustomers: {
    count: number;
    chartLabel: string;
  };
  repeatCustomers: {
    count: number;
    chartLabel: string;
  };
  conversionRate: {
    rate: number;
    chartLabel: string;
  };
}

export interface CustomerChartData {
  day: string;
  count: number;
}
