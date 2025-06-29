// Model Category chung cho toàn hệ thống
export interface Category {
  id: number;
  name: string;
  slug: string;
  description: string;
  imageUrl: string;
  productCount: number;
  isActive: boolean;
  isVisible: boolean;
}

export interface CategoryFilterCounts {
  all: number;
  featured: number;
  onSale: number;
  outOfStock: number;
}

export interface CategorySummaryData {
  totalCategories: number;
  activeCategories: number;
  featuredCategories: number;
  popularCategory: {
    name: string;
    productCount: number;
  };
  conversionRate?: number;
  avgRevenuePerCategory?: number;
}

export interface CategoryResponse {
  categories: Category[];
  total: number;
  page?: number;
  limit?: number;
}

export interface CategoryFilter {
  search?: string;
  isActive?: boolean;
  sortBy?: 'name' | 'product_count' | 'created_at';
  sortOrder?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}