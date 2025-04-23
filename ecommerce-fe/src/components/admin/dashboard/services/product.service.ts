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
//   image_url: string;
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
