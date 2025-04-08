import { Category, CategoryFilter, CategoryResponse, CategoryFilterCounts } from '../models/category.model';
import { CategorySummaryData } from '../cards/models/card.model';

/**
 * Mock data for categories
 * In a real application, this would be fetched from an API
 */
const mockCategories: Category[] = [
  {
    id: '1',
    name: 'Wireless Bluetooth Headphones',
    slug: 'wireless-bluetooth-headphones',
    image_url: 'headphones',
    productCount: 25,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '2',
    name: 'Men\'s T-Shirt',
    slug: 'mens-t-shirt',
    image_url: 'tshirt',
    productCount: 20,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '3',
    name: 'Men\'s Leather Wallet',
    slug: 'mens-leather-wallet',
    image_url: 'wallet',
    productCount: 35,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '4',
    name: 'Memory Foam Pillow',
    slug: 'memory-foam-pillow',
    image_url: 'pillow',
    productCount: 40,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '5',
    name: 'Coffee Maker',
    slug: 'coffee-maker',
    image_url: 'coffee',
    productCount: 45,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '6',
    name: 'Casual Baseball Cap',
    slug: 'casual-baseball-cap',
    image_url: 'cap',
    productCount: 55,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '7',
    name: 'Full HD Webcam',
    slug: 'full-hd-webcam',
    image_url: 'webcam',
    productCount: 20,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '8',
    name: 'Smart LED Color Bulb',
    slug: 'smart-led-color-bulb',
    image_url: 'bulb',
    productCount: 16,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '9',
    name: 'Men\'s T-Shirt',
    slug: 'mens-t-shirt-2',
    image_url: 'tshirt',
    productCount: 10,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
  {
    id: '10',
    name: 'Men\'s Leather Wallet',
    slug: 'mens-leather-wallet-2',
    image_url: 'wallet',
    productCount: 35,
    createdAt: new Date('2025-01-01'),
    updatedAt: new Date('2025-01-01'),
    isActive: true,
  },
];

/**
 * Mock discover categories for the UI
 */
export const discoverCategories = [
  { name: 'Electronics', image_url: '/images/exploring/electronic.png' },
  { name: 'Fashion', image_url: '/images/exploring/fashion.png' },
  { name: 'Accessories', image_url: '/images/exploring/fashion.png' },
  { name: 'Home & Kitchen', image_url: '/images/exploring/home.png' },
  { name: 'Sports & Outdoors', image_url: '/images/exploring/grocery.png' },
  { name: 'Toys & Games', image_url: '/images/exploring/toys.png' },
  { name: 'Health & Fitness', image_url: '/images/exploring/toys.png' },
  { name: 'Books', image_url: '/images/exploring/toys.png' },
];

/**
 * CategoryService class containing all category-related API calls
 * This mock service simulates API calls with Promises
 */
export class CategoryService {
  /**
   * Get all categories with optional filtering and pagination
   */
  static getCategories = async (filter: CategoryFilter): Promise<CategoryResponse> => {
    let filteredCategories = [...mockCategories];
    
    if (filter.search) {
      const searchTerm = filter.search.toLowerCase();
      filteredCategories = filteredCategories.filter(category => 
        category.name.toLowerCase().includes(searchTerm)
      );
    }
    
    if (filter.status !== 'all') {
      filteredCategories = filteredCategories.filter(category => 
        filter.status === 'active' ? category.isActive : !category.isActive
      );
    }
    
    // Apply product filter if specified
    if (filter.productFilter) {
      switch (filter.productFilter) {
        case 'featured':
          // Filter for categories with more than 30 products (mock featured)
          filteredCategories = filteredCategories.filter(
            category => (category.productCount || 0) > 30
          );
          break;
        case 'onSale':
          // Filter for categories with products between 20 and 40 (mock on sale)
          filteredCategories = filteredCategories.filter(
            category => {
              const count = category.productCount || 0;
              return count >= 20 && count <= 40;
            }
          );
          break;
        case 'outOfStock':
          // Filter for categories with less than 15 products (mock out of stock)
          filteredCategories = filteredCategories.filter(
            category => (category.productCount || 0) < 15
          );
          break;
        // 'all' case doesn't need filtering
      }
    }
    
    // Get total count before pagination
    const total = filteredCategories.length;
    
    // Pagination
    const page = filter.page || 1;
    const limit = filter.limit || 10;
    const start = (page - 1) * limit;
    const paginatedCategories = filteredCategories.slice(start, start + limit);
    
    return Promise.resolve({
      categories: paginatedCategories,
      total
    });
  };

  /**
   * Get category summary data for dashboard cards
   */
  static getCategorySummary = async (): Promise<CategorySummaryData> => {
    const totalCategories = mockCategories.length;
    const activeCategories = mockCategories.filter(
      (cat) => cat.isActive
    ).length;
    const featuredCategories = Math.floor(totalCategories * 0.4); // 40% of categories are featured (mock)

    // Find most popular category based on product count
    const popularCategory = mockCategories.reduce((prev, current) =>
      (prev.productCount || 0) > (current.productCount || 0) ? prev : current
    );

    return Promise.resolve({
      totalCategories,
      activeCategories,
      featuredCategories,
      popularCategory: {
        name: popularCategory.name,
        productCount: popularCategory.productCount || 0,
      },
    });
  };

  /**
   * Get category by ID
   */
  static getCategoryById = async (
    id: string
  ): Promise<Category | undefined> => {
    const category = mockCategories.find((cat) => cat.id === id);
    return Promise.resolve(category);
  };

  /**
   * Create a new category
   */
  static createCategory = async (
    category: Omit<Category, "id">
  ): Promise<Category> => {
    const newCategory = {
      ...category,
      id: String(Math.max(...mockCategories.map((c) => parseInt(c.id))) + 1),
    };
    mockCategories.push(newCategory as Category);
    return Promise.resolve(newCategory as Category);
  };

  /**
   * Update an existing category
   */
  static updateCategory = async (category: Category): Promise<Category> => {
    const index = mockCategories.findIndex((c) => c.id === category.id);
    if (index !== -1) {
      mockCategories[index] = category;
    }
    return Promise.resolve(category);
  };

  /**
   * Delete a category
   */
  static deleteCategory = async (id: string ): Promise<boolean> => {
    const index = mockCategories.findIndex((c) => c.id === id);
    if (index !== -1) {
      mockCategories.splice(index, 1);
      return Promise.resolve(true);
    }
    return Promise.resolve(false);
  };

  /**
   * Get counts for different product filters
   * In a real implementation, this would fetch from an API
   */
  static getCategoryFilterCounts = async (): Promise<CategoryFilterCounts> => {
    // Calculate total products
    const totalProducts = mockCategories.reduce(
      (sum, cat) => sum + (cat.productCount || 0),
      0
    );

    // Mock featured products (about 30% of total)
    const featured = Math.floor(totalProducts * 0.3);

    // Mock on sale products (about 20% of total)
    const onSale = Math.floor(totalProducts * 0.2);

    // Mock out of stock products (about 10% of total)
    const outOfStock = Math.floor(totalProducts * 0.1);

    return Promise.resolve({
      all: totalProducts,
      featured,
      onSale,
      outOfStock,
    });
  };
}

// Export a singleton instance
export const categoryService = new CategoryService();