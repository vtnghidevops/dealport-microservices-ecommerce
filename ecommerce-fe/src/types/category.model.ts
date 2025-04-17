// Model Category chung cho toàn hệ thống
export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  image_url?: string;
  icon?: string;
  banner_url?: string;
  productCount?: number;
  createdAt?: string;
  updatedAt?: string;
  isActive?: boolean;
  isVisible?: boolean;
  displayOrder?: number;
  metaTitle?: string;
  metaDescription?: string;
  parentId?: string | null;
  ancestors?: string[]; // Array of parent category IDs for breadcrumb navigation
  level?: number; // Depth in category tree (0 for root)
  attributes?: CategoryAttribute[];
}

export interface CategoryAttribute {
  id: string;
  name: string;
  type: 'text' | 'number' | 'boolean' | 'select';
  required?: boolean;
  options?: string[]; // For select type attributes
}

// Category item cho category explorer
export interface CategoryItem {
  id: string;
  name: string;
  image_url: string;
  link?: string;
  slug: string;
  productCount?: number;
  children?: CategoryItem[];
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

export interface CategoryTree {
  id: string;
  name: string;
  slug: string;
  image_url?: string;
  productCount?: number;
  children: CategoryTree[];
}

export interface CategoryFilter {
  search?: string;
  parentId?: string | null;
  isActive?: boolean;
  sortBy?: 'name' | 'productCount' | 'createdAt';
  sortOrder?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}