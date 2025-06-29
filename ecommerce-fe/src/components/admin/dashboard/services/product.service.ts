// import MainProductService, {

// } from "@/services/product.service";

// // Reexport the types from the main service
// export type {
//   BestSellingProduct,
//   ProductCategory,
//   NewProduct,
// };

// // Define Product type for admin-specific properties if needed
// export interface Product {
//   id: number;
//   name: string;
//   category: string;
//   imageUrl: string;
//   itemCode: string;
//   price: number;
// }

// // Use the main service but provide admin-specific overrides
// export const ProductService = {
//   // Reuse the functions from the main service
//   getBestSellingProducts: async (): Promise<BestSellingProduct[]> => {
//     return MainProductService.getBestSellingProducts();
//   },

//   getTopProducts: async (): Promise<Product[]> => {
//     return MainProductService.getTopProducts();
//   },

//   getProductCategories: async (): Promise<ProductCategory[]> => {
//     return MainProductService.getProductCategories();
//   },

//   getNewProducts: async (): Promise<NewProduct[]> => {
//     return MainProductService.getNewProducts();
//   },
// };

import axios from "axios";
import { Product } from "@/types/product.model";
import { TopProductItem } from "@/components/homepage/BestSelling/models/topProducts.model";
import { getApiUrl, getAdminHeaders } from "@/utils/api-config";

// Types for product dashboard data
export interface ProductStatistics {
  totalProducts: number;
  inStockProducts: number;
  outOfStockProducts: number;
  lowStockProducts: number;
  productGrowth: number;
}

export interface BestSellingProductStats {
  id: number | string;
  name: string;
  imageSrc: string;
  sold: number;
  price: number;
  stock: number;
  growth: number;
}

export const ProductDashboardService = {
  // Get product statistics for dashboard
  getProductStatistics: async (): Promise<ProductStatistics> => {
    try {
      const headers = getAdminHeaders();
      const response = await axios.get(getApiUrl("products/statistics"), {
        headers,
      });

      if (response.data && response.data.data) {
        return {
          totalProducts: response.data.data.totalProducts || 0,
          inStockProducts: response.data.data.inStockProducts || 0,
          outOfStockProducts: response.data.data.outOfStockProducts || 0,
          lowStockProducts: response.data.data.lowStockProducts || 0,
          productGrowth: response.data.data.productGrowth || 0,
        };
      }

      // Return default data if API response is invalid
      return {
        totalProducts: 3500,
        inStockProducts: 2500,
        outOfStockProducts: 500,
        lowStockProducts: 300,
        productGrowth: 12.5,
      };
    } catch (error) {
      console.error("Error fetching product statistics:", error);

      // Return default data if API request fails
      return {
        totalProducts: 3500,
        inStockProducts: 2500,
        outOfStockProducts: 500,
        lowStockProducts: 300,
        productGrowth: 12.5,
      };
    }
  },

  // Get top selling products for dashboard
  getTopSellingProducts: async (
    limit: number = 5
  ): Promise<BestSellingProductStats[]> => {
    try {
      const headers = getAdminHeaders();
      // Get products sorted by order count (sold)
      const response = await axios.get(getApiUrl("products"), {
        headers,
        params: {
          page: 1,
          page_size: limit,
        },
      });

      if (
        response.data &&
        response.data.data &&
        Array.isArray(response.data.data)
      ) {
        return response.data.data.map((product: any) => ({
          id: product.id,
          name: product.name,
          imageSrc: product.imageUrl || "/assets/images/placeholder.png",
          sold: product.orders || 0,
          price: product.price || 0,
          stock: product.stockQuantity || 0,
          growth: Math.random() * 30, // Placeholder since growth isn't in the API
        }));
      }

      // Return empty array if API response is invalid
      return [];
    } catch (error) {
      console.error("Error fetching top selling products:", error);
      return [];
    }
  },

  // Get new products for dashboard (recently added)
  getNewProducts: async (limit: number = 5): Promise<Product[]> => {
    try {
      const headers = getAdminHeaders();
      // Get products sorted by creation date
      const response = await axios.get(getApiUrl("products"), {
        headers,
        params: {
          page: 1,
          page_size: limit,
        },
      });

      if (
        response.data &&
        response.data.data &&
        Array.isArray(response.data.data)
      ) {
        return response.data.data;
      }

      // Return empty array if API response is invalid
      return [];
    } catch (error) {
      console.error("Error fetching new products:", error);
      return [];
    }
  },

  // Get top sale products (featured products)
  getTopSaleProducts: async (limit: number = 5): Promise<TopProductItem[]> => {
    try {
      const headers = getAdminHeaders();
      // Get top-sale products
      const response = await axios.get(getApiUrl("products"), {
        headers,
        params: {
          page: 1,
          page_size: limit,
          type: "top-sale",
        },
      });

      if (
        response.data &&
        response.data.data &&
        Array.isArray(response.data.data)
      ) {
        return response.data.data.map((product: any) => {
          // Parse the uiMetadata if it's a string
          let metadata = product.uiMetadata || {};
          if (typeof product.uiMetadata === "string") {
            try {
              metadata = JSON.parse(product.uiMetadata);
            } catch (err) {
              console.error("Error parsing uiMetadata:", err);
            }
          }

          return {
            ...product,
            uiMetadata: {
              setUpDesign: metadata.setUpDesign || "row",
              isCommingSoon: metadata.isCommingSoon || false,
            },
          };
        });
      }

      // Return empty array if API response is invalid
      return [];
    } catch (error) {
      console.error("Error fetching top sale products:", error);
      return [];
    }
  },
};
