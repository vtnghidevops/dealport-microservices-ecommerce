import axios from "axios";
import { Category } from "@/types/category.model";
import { TopProductItem } from "@/components/homepage/BestSelling/models/topProducts.model";
import { Product, ProductReview } from "@/types/product.model";
import { Banner, BannerType, SliderBannerItem, BannerResponse } from '@/types/banner.model';
import { TestimonialItem } from "@/components/homepage/HappyCustomers/models/testimonial.model";
import { getApiUrl } from '@/utils/api-config';

// Additional interface definitions to fix type errors
interface BestSellingProduct extends Omit<Product, 'originalPrice' | 'orders'> {
  sold: number;
  onSale: boolean;
  originalPrice: number | null;
  orders?: number;
  discount_percent: number;
}

interface NewProduct extends Omit<Product, 'originalPrice'> {
  onSale: boolean;
  originalPrice: number | null;
  discount_percent: number;
}

interface ProductCategory {
  id: string;
  name: string;
  slug: string;
  description?: string;
  imageUrl?: string;
  parentId?: string;
  children?: ProductCategory[];
  level?: number;
  isActive?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

interface ProductBrand {
  id: string;
  name: string;
  slug: string;
  description?: string;
  logoUrl?: string;
  websiteUrl?: string;
  isActive?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

/**
 * Product Service
 * 
 * Use methods to get products by type:
 * - getAllProducts: Get all products with pagination and optional filters (lấy tất cả sản phẩm với phân trang và tùy chọn filter)
 * - getProductsByType: Generic method to get products by type and filters (phương thức chung để lấy sản phẩm theo type và filter)
 * - getTrendingProducts: Get trending products (type="trending") (lấy sản phẩm trending)
 * - getMenCollection: Get trending products for men (type="trending", category_slug="men") (lấy sản phẩm trending cho nam)
 * - getTopSaleProducts: Get top-sale products with UI metadata (type="top-sale") (lấy sản phẩm top-sale với metadata UI)
 * - getLimitedProducts: Get limited edition products (type="limited") (lấy sản phẩm limited edition)
 */

// API Response interfaces
export interface PaginationMeta {
  current_page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  status: number;
  message?: string;
  data: T;
  meta?: PaginationMeta;
}

// Error handling helper
const handleApiError = (error: any): never => {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.message || error.message;
    console.error(`API Error: ${message}`);
    throw new Error(message);
  }
  console.error("Unexpected error:", error);
  throw error;
};

// Service implementation
const ProductService = {
  // Products API
  getAllProducts: async (
    page = 1,
    pageSize = 10,
    filters?: Record<string, string>
  ): Promise<{ products: Product[], pagination: PaginationMeta }> => {
    try {
      // Build query parameters
      const params = new URLSearchParams({
        page: page.toString(),
        limit: pageSize.toString(),
        ...filters
      });

      const response = await axios.get<ApiResponse<Product[]>>(
        `${getApiUrl('products')}?${params.toString()}`
      );

      return {
        products: response.data.data,
        pagination: response.data.meta as PaginationMeta
      };
    } catch (error) {
      return handleApiError(error);
    }
  },

  getProductById: async (id: string): Promise<Product> => {
    try {
      const response = await axios.get<ApiResponse<Product>>(
        `${getApiUrl(`products/${id}`)}`
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  getProductBySlug: async (slug: string): Promise<Product> => {
    try {
      const response = await axios.get<ApiResponse<Product>>(
        `${getApiUrl(`products/slug/${slug}`)}`
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  createProduct: async (product: Omit<Product, "id">): Promise<Product> => {
    try {
      const response = await axios.post<ApiResponse<Product>>(
        getApiUrl('products'),
        product
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  updateProduct: async (id: string, product: Partial<Product>): Promise<Product> => {
    try {
      const response = await axios.put<ApiResponse<Product>>(
        getApiUrl(`products/${id}`),
        product
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  // Patch update for partial updates (only changed fields)
  patchProduct: async (id: string, partialProduct: Partial<Product>): Promise<Product> => {
    try {
      const response = await axios.patch<ApiResponse<Product>>(
        getApiUrl(`products/${id}`),
        partialProduct
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  deleteProduct: async (id: string): Promise<boolean> => {
    try {
      await axios.delete(getApiUrl(`products/${id}`));
      return true;
    } catch (error) {
      return handleApiError(error);
    }
  },

  // Product Reviews
  getProductReviews: async (
    productId: string,
    page = 1,
    pageSize = 10
  ): Promise<{ reviews: ProductReview[], pagination: PaginationMeta }> => {
    try {
      const response = await axios.get<ApiResponse<ProductReview[]>>(
        `${getApiUrl(`products/${productId}/reviews`)}?page=${page}&limit=${pageSize}`
      );

      return {
        reviews: response.data.data,
        pagination: response.data.meta as PaginationMeta
      };
    } catch (error) {
      return handleApiError(error);
    }
  },

  addProductReview: async (review: Omit<ProductReview, "id">): Promise<ProductReview> => {
    try {
      const response = await axios.post<ApiResponse<ProductReview>>(
        getApiUrl(`products/${review.productId}/reviews`),
        review
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  // Generic method to get products by type
  getProductsByType: async (type: string, page = 1, pageSize = 10, additionalFilters?: Record<string, string>): Promise<Product[]> => {
    try {
      const filters = {
        type,
        ...additionalFilters
      };

      const { products } = await ProductService.getAllProducts(page, pageSize, filters);
      return products;
    } catch (error) {
      console.error(`Error fetching products with type ${type}:`, error);
      return [];
    }
  },

  // Get trending products
  getTrendingProducts: async (page = 1, pageSize = 10): Promise<Product[]> => {
    return ProductService.getProductsByType("trending", page, pageSize);
  },

  // Get Men's Collection (trending products in men category)
  getMenCollection: async (page = 1, pageSize = 10): Promise<Product[]> => {
    return ProductService.getProductsByType("trending", page, pageSize, { category_slug: "men" });
  },

  // Get Top-Sale Products with UI metadata from backend
  getTopSaleProducts: async (page = 1, pageSize = 10): Promise<TopProductItem[]> => {
    try {
      // Get top-sale products
      const response = await axios.get<ApiResponse<TopProductItem[]>>(
        `${getApiUrl('products')}?page=${page}&limit=${pageSize}&type=top-sale`
      );

      // Process and return the data
      if (response.data && response.data.data) {
        return response.data.data.map(product => {
          // Parse the uiMetadata if it's a string
          let metadata = product.uiMetadata || {};
          if (typeof product.uiMetadata === 'string') {
            try {
              metadata = JSON.parse(product.uiMetadata);
            } catch (err) {
              console.error("Error parsing uiMetadata:", err);
            }
          }

          return {
            ...product,
            uiMetadata: {
              setUpDesign: metadata.setUpDesign || 'row',
              isCommingSoon: metadata.isCommingSoon || false
            }
          };
        });
      }

      return [];
    } catch (error) {
      console.error('Error fetching top sale products:', error);
      return [];
    }
  },

  // Get Limited Edition Products
  getLimitedProducts: async (page = 1, pageSize = 10): Promise<Product[]> => {
    return ProductService.getProductsByType("limited", page, pageSize);
  },

  // Get Best Selling Products
  getBestSellingProducts: async (page = 1, pageSize = 10): Promise<BestSellingProduct[]> => {
    try {
      const response = await axios.get<ApiResponse<BestSellingProduct[]>>(
        `${getApiUrl('products')}?page=${page}&limit=${pageSize}&sort_by=orders&sort_dir=desc`
      );

      // Type cast to BestSellingProduct and add some additional UI fields
      return (response.data.data || []).map(product => ({
        ...product,
        sold: product.orders || 0,
        onSale: product.discount_percent > 0,
        originalPrice: product.discount_percent > 0 ? product.price / (1 - product.discount_percent / 100) : null
      }));
    } catch (error) {
      console.error('Error fetching best selling products:', error);
      return [];
    }
  },

  // Get new products (most recently added)
  getNewProducts: async (page = 1, pageSize = 10): Promise<NewProduct[]> => {
    try {
      const response = await axios.get<ApiResponse<NewProduct[]>>(
        `${getApiUrl('products')}?page=${page}&limit=${pageSize}&sort_by=created_at&sort_dir=desc`
      );

      // Return the data with some additional UI fields
      return (response.data.data || []).map(product => ({
        ...product,
        onSale: product.discount_percent > 0,
        originalPrice: product.discount_percent > 0 ? product.price / (1 - product.discount_percent / 100) : null
      }));
    } catch (error) {
      console.error('Error fetching new products:', error);
      return [];
    }
  },

  // Get product categories with optional parent ID filter
  getProductCategories: async (parentId?: string): Promise<ProductCategory[]> => {
    try {
      let url = getApiUrl('categories');
      if (parentId) {
        url += `?parent_id=${parentId}`;
      }

      const response = await axios.get<ApiResponse<ProductCategory[]>>(url);
      return response.data.data || [];
    } catch (error) {
      console.error('Error fetching product categories:', error);
      return [];
    }
  },

  // Get category tree (hierarchical structure)
  getCategoryTree: async (): Promise<ProductCategory[]> => {
    try {
      const response = await axios.get<ApiResponse<ProductCategory[]>>(
        getApiUrl('categories/tree')
      );
      return response.data.data || [];
    } catch (error) {
      console.error('Error fetching category tree:', error);
      return [];
    }
  },

  // Get all product brands
  getProductBrands: async (): Promise<ProductBrand[]> => {
    try {
      const response = await axios.get<ApiResponse<ProductBrand[]>>(
        getApiUrl('brands')
      );
      return response.data.data || [];
    } catch (error) {
      console.error('Error fetching product brands:', error);
      return [];
    }
  },

  // Product Health Check
  getProductHealth: async (): Promise<boolean> => {
    try {
      const response = await axios.get(
        getApiUrl('products/health')
      );
      return response.status === 200;
    } catch (error) {
      console.error('Product service health check failed:', error);
      return false;
    }
  },

  // Get products by category slug
  getProductsByCategorySlug: async (
    categorySlug: string,
    page = 1,
    pageSize = 10
  ): Promise<{ products: Product[], pagination: PaginationMeta }> => {
    return ProductService.getAllProducts(page, pageSize, {
      category_slug: categorySlug
    });
  },

  // Get products - this is a legacy method for compatibility
  getProducts: async (): Promise<Product[]> => {
    const { products } = await ProductService.getAllProducts(1, 100);
    return products;
  },

  getProductsByCategory: async (categoryId: string): Promise<Product[]> => {
    const { products } = await ProductService.getAllProducts(1, 100, {
      category_id: categoryId
    });
    return products;
  },

};

const CategoryService = {
  // Categories API
  getAllCategories: async (filters?: Record<string, string>): Promise<Category[]> => {
    try {
      const params = new URLSearchParams(filters);
      const response = await axios.get<ApiResponse<Category[]>>(
        `${getApiUrl('categories')}?${params.toString()}`
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  getCategoryById: async (id: string): Promise<Category> => {
    try {
      const response = await axios.get<ApiResponse<Category>>(
        getApiUrl(`categories/${id}`)
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },

  getCategoryBySlug: async (slug: string): Promise<Category> => {
    try {
      const response = await axios.get<ApiResponse<Category>>(
        getApiUrl(`categories/slug/${slug}`)
      );
      return response.data.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

const BannerService = {
  // Slider banner methods
  getSliderBanners: async (): Promise<SliderBannerItem[]> => {
    try {
      const response = await axios.get(getApiUrl('banners'), {
        params: {
          type: 'hero',
          is_active: true,
          order_by: 'priority',
          order_dir: 'ASC'
        }
      });

      const { data } = response.data as BannerResponse;
      return data.filter(banner => banner.type === 'hero').map(banner => ({
        id: banner.id,
        title: banner.title,
        subtitle: banner.subtitle,
        description: banner.description,
        imageUrl: banner.imageUrl,
        linkUrl: banner.linkUrl,
        actionText: banner.actionText,
        isActive: banner.isActive,
        priority: banner.priority,
        discount: banner.discount,
        highlightText: banner.highlightText,
        backgroundColor: banner.backgroundColor,
        textColor: banner.textColor,
        animationType: banner.animationType
      }));
    } catch (error) {
      console.error('Error fetching slider banners:', error);
      return [];
    }
  },

  getSliderBannerById: async (id: number): Promise<SliderBannerItem | undefined> => {
    try {
      const response = await axios.get(getApiUrl(`banners/${id}`));
      const banner = response.data.data;
      return {
        id: banner.id,
        title: banner.title,
        subtitle: banner.subtitle,
        description: banner.description,
        imageUrl: banner.imageUrl,
        linkUrl: banner.linkUrl,
        actionText: banner.actionText,
        isActive: banner.isActive,
        priority: banner.priority,
        discount: banner.discount,
        highlightText: banner.highlightText,
        backgroundColor: banner.backgroundColor,
        textColor: banner.textColor,
        animationType: banner.animationType
      };
    } catch (error) {
      console.error(`Error fetching slider banner with ID ${id}:`, error);
      return undefined;
    }
  },

  // General banner methods
  getBanners: async (type?: BannerType, limit?: number): Promise<Banner[]> => {
    try {
      const params: Record<string, any> = {
        is_active: true,
        order_by: 'priority',
        order_dir: 'ASC'
      };

      if (type) {
        params.type = type;
      }

      if (limit) {
        params.page_size = limit;
      }

      const response = await axios.get(getApiUrl('banners'), { params });
      return response.data.data;
    } catch (error) {
      console.error('Error fetching banners:', error);
      return [];
    }
  },

  getBannerById: async (id: number): Promise<Banner | undefined> => {
    try {
      const response = await axios.get(getApiUrl(`banners/${id}`));
      return response.data.data;
    } catch (error) {
      console.error(`Error fetching banner with ID ${id}:`, error);
      return undefined;
    }
  }
};

const TestimonialService = {
  getTestimonials: async (limit: number = 7): Promise<TestimonialItem[]> => {
    try {
      // Try to fetch from API (cố gắng lấy dữ liệu từ API)
      const response = await axios.get(getApiUrl('testimonials'), {
        params: { limit }
      });

      if (response.data && response.data.data) {
        return response.data.data;
      }

      throw new Error('Invalid API response format');
    } catch (error) {
      console.error('Error fetching testimonials from API, using fallback data:', error);
      return []
    }
  }
};
export default ProductService;
export { BannerService, CategoryService, TestimonialService };
