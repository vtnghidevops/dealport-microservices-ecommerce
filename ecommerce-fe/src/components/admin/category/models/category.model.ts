// Category model definition for the application
import { Category as GlobalCategory } from '@/types/category.model';

// Extend the global category model to make it compatible with the admin category model
export interface Category extends Omit<GlobalCategory, 'id'> {
  id: string; // ID as string in admin module
  createdAt: Date;
  updatedAt: Date;
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

