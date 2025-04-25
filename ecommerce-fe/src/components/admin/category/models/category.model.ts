// Category model definition for the application
export interface Category {
  id: string;
  name: string;
  slug: string;
  imageUrl: string; // Path or name of icon
  productCount: number; // Optional count of products in this category
  createdAt: Date;
  updatedAt: Date;
  isActive: boolean;
}

// Category filter options
export interface CategoryFilter {
  search: string;
  status: 'all' | 'active' | 'inactive';
  page?: number;
  limit?: number;
  productFilter?: 'all' | 'featured' | 'onSale' | 'outOfStock'; // Thêm trường này
}
/**
 * Add filter counts type to existing category model
 */
export interface CategoryFilterCounts {
  all: number;
  featured: number;
  onSale: number;
  outOfStock: number;
}
// Available category types for the discover section
export type CategoryType =
  | 'Electronics'
  | 'Fashion'
  | 'Accessories'
  | 'Home & Kitchen'
  | 'Sports & Outdoors'
  | 'Toys & Games'
  | 'Health & Fitness'
  | 'Books';

/**
 * Response structure for paginated category API calls
 */
export interface CategoryResponse {
  categories: Category[];
  total: number;
}

