import { Category, CategoryFilter, CategoryResponse, CategoryFilterCounts } from '../models/category.model';
import { CategorySummaryData } from '../cards/models/card.model';
import { CategoryService as GlobalCategoryService } from "@/services/product/product.service";
import { Category as GlobalCategory } from '@/types/category.model';
import axios from 'axios';

// Base URL for product service API from environment variables
const API_BASE_URL = import.meta.env.API_URL || "http://localhost:8082/api/v1";

/**
 * Helper function to convert a global category to admin category format
 */
const convertToAdminCategory = (globalCategory: GlobalCategory): Category => {
  return {
    ...globalCategory,
    id: String(globalCategory.id), // Convert ID to string
    createdAt: new Date(), // Add createdAt if missing
    updatedAt: new Date(), // Add updatedAt if missing
  };
};

/**
 * Helper function to convert an admin category to global category format
 */
const convertToGlobalCategory = (adminCategory: Category): GlobalCategory => {
  return {
    ...adminCategory,
    id: Number(adminCategory.id), // Convert ID to number
  };
};

/**
 * CategoryService class containing all category-related API calls
 * Uses the global CategoryService for real API integration
 */
export class CategoryService {
  /**
   * Get all categories with optional filtering and pagination
   */
  static getCategories = async (filter: CategoryFilter): Promise<CategoryResponse> => {
    try {
      // Create API filter parameters
      const apiFilters: Record<string, string> = {};

      if (filter.search) {
        apiFilters.name = filter.search;
      }

      if (filter.status !== 'all') {
        apiFilters.is_active = filter.status === 'active' ? 'true' : 'false';
      }

      // Get categories from API
      const globalCategories = await GlobalCategoryService.getAllCategories(apiFilters);

      // Convert global categories to admin categories
      const categories: Category[] = globalCategories.map(convertToAdminCategory);

      // Apply additional product filters not handled by the API
      let filteredCategories = [...categories];

      if (filter.productFilter && filter.productFilter !== 'all') {
        switch (filter.productFilter) {
          case 'featured':
            // Filter for categories with more products
            filteredCategories = filteredCategories.filter(
              category => (category.productCount || 0) > 20
            );
            break;
          case 'onSale':
            // Filter for categories that might have products on sale
            filteredCategories = filteredCategories.filter(
              category => {
                const count = category.productCount || 0;
                return count >= 10 && count <= 30;
              }
            );
            break;
          case 'outOfStock':
            // Filter for categories with few products
            filteredCategories = filteredCategories.filter(
              category => (category.productCount || 0) < 10
            );
            break;
        }
      }

      // Get total count before pagination
      const total = filteredCategories.length;

      // Manual pagination
      const page = filter.page || 1;
      const limit = filter.limit || 10;
      const start = (page - 1) * limit;
      const paginatedCategories = filteredCategories.slice(start, start + limit);

      return {
        categories: paginatedCategories,
        total
      };
    } catch (error) {
      console.error("Error fetching categories:", error);
      return {
        categories: [],
        total: 0
      };
    }
  };

  /**
   * Get category summary data for dashboard cards
   */
  static getCategorySummary = async (): Promise<CategorySummaryData> => {
    try {
      const globalCategories = await GlobalCategoryService.getAllCategories();

      // Convert global categories to admin categories
      const categories: Category[] = globalCategories.map(convertToAdminCategory);

      const totalCategories = categories.length;
      const activeCategories = categories.filter(cat => cat.isActive).length;
      const featuredCategories = Math.floor(totalCategories * 0.4); // Estimate 40% as featured

      // Find most popular category based on product count
      let popularCategory: Category = {
        id: '0',
        name: 'Unknown',
        slug: '',
        description: '',
        imageUrl: '',
        productCount: 0,
        isActive: false,
        isVisible: false,
        createdAt: new Date(),
        updatedAt: new Date()
      };

      if (categories.length > 0) {
        popularCategory = categories.reduce((prev, current) =>
          (prev.productCount || 0) > (current.productCount || 0) ? prev : current
        );
      }

      return {
        totalCategories,
        activeCategories,
        featuredCategories,
        popularCategory: {
          name: popularCategory.name,
          productCount: popularCategory.productCount || 0,
        },
      };
    } catch (error) {
      console.error("Error fetching category summary:", error);
      return {
        totalCategories: 0,
        activeCategories: 0,
        featuredCategories: 0,
        popularCategory: {
          name: 'Unknown',
          productCount: 0,
        },
      };
    }
  };

  /**
   * Get category by ID
   */
  static getCategoryById = async (id: string): Promise<Category | undefined> => {
    try {
      const globalCategory = await GlobalCategoryService.getCategoryById(id);
      if (globalCategory) {
        return convertToAdminCategory(globalCategory);
      }
      return undefined;
    } catch (error) {
      console.error(`Error fetching category with ID ${id}:`, error);
      return undefined;
    }
  };

  /**
   * Create a new category
   */
  static createCategory = async (category: Omit<Category, "id">): Promise<Category> => {
    try {
      // Use the API to create a category
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('Authentication token is missing');
      }

      const response = await axios.post(
        `${API_BASE_URL}/categories`,
        category,
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          }
        }
      );

      const createdCategory = response.data.data;
      return convertToAdminCategory(createdCategory);
    } catch (error) {
      console.error("Error creating category:", error);
      throw error;
    }
  };

  /**
   * Update an existing category
   */
  static updateCategory = async (category: Category): Promise<Category> => {
    try {
      // Use the API to update a category
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('Authentication token is missing');
      }

      // Convert to global category format for API
      const globalCategory = convertToGlobalCategory(category);

      const response = await axios.put(
        `${API_BASE_URL}/categories/${globalCategory.id}`,
        globalCategory,
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          }
        }
      );

      const updatedCategory = response.data.data;
      return convertToAdminCategory(updatedCategory);
    } catch (error) {
      console.error("Error updating category:", error);
      throw error;
    }
  };

  /**
   * Delete a category by ID
   */
  static deleteCategory = async (id: string): Promise<boolean> => {
    try {
      // Use the API to delete a category
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('Authentication token is missing');
      }

      await axios.delete(
        `${API_BASE_URL}/categories/${id}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      );

      return true;
    } catch (error) {
      console.error(`Error deleting category with ID ${id}:`, error);
      return false;
    }
  };

  /**
   * Get counts for each filter type
   */
  static getCategoryFilterCounts = async (): Promise<CategoryFilterCounts> => {
    try {
      const globalCategories = await GlobalCategoryService.getAllCategories();

      // Convert global categories to admin categories
      const categories: Category[] = globalCategories.map(convertToAdminCategory);

      // Count for all categories
      const all = categories.length;

      // Count for categories with more products (featured)
      const featured = categories.filter(cat => (cat.productCount || 0) > 20).length;

      // Count for categories with medium products (on sale)
      const onSale = categories.filter(cat => {
        const count = cat.productCount || 0;
        return count >= 10 && count <= 30;
      }).length;

      // Count for categories with low products (out of stock)
      const outOfStock = categories.filter(cat => (cat.productCount || 0) < 10).length;

      return {
        all,
        featured,
        onSale,
        outOfStock
      };
    } catch (error) {
      console.error("Error fetching filter counts:", error);
      return {
        all: 0,
        featured: 0,
        onSale: 0,
        outOfStock: 0
      };
    }
  };
}