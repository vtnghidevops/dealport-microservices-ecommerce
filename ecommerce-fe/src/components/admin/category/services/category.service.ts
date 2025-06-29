import {
  Category,
  CategoryFilter,
  CategoryResponse,
  CategoryFilterCounts,
} from "../models/category.model";
import { CategorySummaryData } from "../cards/models/card.model";
import { CategoryService as GlobalCategoryService } from "@/services/product/product.service";
import { Category as GlobalCategory } from "@/types/category.model";
import axios from "axios";
import { getApiUrl, getAuthHeader } from "@/utils/api-config";

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
 * CategoryService class containing all category-related API calls
 * Uses the global CategoryService for real API integration
 */
export class CategoryService {
  /**
   * Get all categories with optional filtering and pagination
   */
  static getCategories = async (
    filter: CategoryFilter
  ): Promise<CategoryResponse> => {
    try {
      // Create API filter parameters
      const apiFilters: Record<string, string> = {};

      if (filter.search) {
        apiFilters.name = filter.search;
      }

      if (filter.status !== "all") {
        apiFilters.is_active = filter.status === "active" ? "true" : "false";
      }

      // Get categories from API
      const globalCategories = await GlobalCategoryService.getAllCategories(
        apiFilters
      );

      // Convert global categories to admin categories
      const categories: Category[] = globalCategories.map(
        convertToAdminCategory
      );

      // Apply additional product filters not handled by the API
      let filteredCategories = [...categories];

      if (filter.productFilter && filter.productFilter !== "all") {
        switch (filter.productFilter) {
          case "featured":
            // Filter for categories with more products
            filteredCategories = filteredCategories.filter(
              (category) => (category.productCount || 0) > 20
            );
            break;
          case "onSale":
            // Filter for categories that might have products on sale
            filteredCategories = filteredCategories.filter((category) => {
              const count = category.productCount || 0;
              return count >= 10 && count <= 30;
            });
            break;
          case "outOfStock":
            // Filter for categories with few products
            filteredCategories = filteredCategories.filter(
              (category) => (category.productCount || 0) < 10
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
      const paginatedCategories = filteredCategories.slice(
        start,
        start + limit
      );

      return {
        categories: paginatedCategories,
        total,
      };
    } catch (error) {
      console.error("Error fetching categories:", error);
      return {
        categories: [],
        total: 0,
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
      const categories: Category[] = globalCategories.map(
        convertToAdminCategory
      );

      const totalCategories = categories.length;
      const activeCategories = categories.filter((cat) => cat.isActive).length;
      const featuredCategories = Math.floor(totalCategories * 0.4); // Estimate 40% as featured

      // Find most popular category based on product count
      let popularCategory: Category = {
        id: "0",
        name: "Unknown",
        slug: "",
        description: "",
        imageUrl: "",
        productCount: 0,
        isActive: false,
        isVisible: false,
        createdAt: new Date(),
        updatedAt: new Date(),
      };

      if (categories.length > 0) {
        popularCategory = categories.reduce((prev, current) =>
          (prev.productCount || 0) > (current.productCount || 0)
            ? prev
            : current
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
          name: "Unknown",
          productCount: 0,
        },
      };
    }
  };

  /**
   * Get category by ID
   */
  static getCategoryById = async (
    id: string
  ): Promise<Category | undefined> => {
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
  static createCategory = async (
    category: Omit<Category, "id">
  ): Promise<Category> => {
    try {
      // Use the API to create a category
      const headers = getAuthHeader();

      const response = await axios.post(getApiUrl("categories"), category, {
        headers,
      });

      const newCategory = response.data.data || response.data;
      return convertToAdminCategory({
        ...newCategory,
        id: Number(newCategory.id),
      });
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
      const headers = getAuthHeader();

      const response = await axios.put(
        getApiUrl(`categories/${category.id}`),
        {
          ...category,
          id: Number(category.id),
        },
        { headers }
      );

      const updatedCategory = response.data.data || response.data;
      return convertToAdminCategory(updatedCategory);
    } catch (error) {
      console.error(`Error updating category ${category.id}:`, error);
      throw error;
    }
  };

  /**
   * Delete a category
   */
  static deleteCategory = async (id: string): Promise<boolean> => {
    try {
      // Use the API to delete a category
      const headers = getAuthHeader();

      await axios.delete(getApiUrl(`categories/${id}`), { headers });

      return true;
    } catch (error) {
      console.error(`Error deleting category ${id}:`, error);
      return false;
    }
  };

  /**
   * Get counts for each category filter
   */
  static getCategoryFilterCounts = async (): Promise<CategoryFilterCounts> => {
    try {
      const globalCategories = await GlobalCategoryService.getAllCategories();
      const categories: Category[] = globalCategories.map(
        convertToAdminCategory
      );

      // Count categories by status
      const activeCount = categories.filter((cat) => cat.isActive).length;
      const inactiveCount = categories.length - activeCount;

      // Count categories by product count
      const featuredCount = categories.filter(
        (cat) => (cat.productCount || 0) > 20
      ).length;
      const onSaleCount = categories.filter((cat) => {
        const count = cat.productCount || 0;
        return count >= 10 && count <= 30;
      }).length;
      const outOfStockCount = categories.filter(
        (cat) => (cat.productCount || 0) < 10
      ).length;

      return {
        all: categories.length,
        active: activeCount,
        inactive: inactiveCount,
        featured: featuredCount,
        onSale: onSaleCount,
        outOfStock: outOfStockCount,
      };
    } catch (error) {
      console.error("Error getting category filter counts:", error);
      return {
        all: 0,
        active: 0,
        inactive: 0,
        featured: 0,
        onSale: 0,
        outOfStock: 0,
      };
    }
  };
}
