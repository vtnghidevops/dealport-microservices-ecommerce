

// src/components/product/services/product.service.ts
import { Product } from '../models/product.model';

// Mock data
const products: Product[] = [
  {
    id: '1',
    name: 'Dell XPS 13',
    slug: 'dell-xps-13',
    price: 25990000,
    description: 'Laptop cao cấp với màn hình InfinityEdge',
    imageUrl: '/images/products/dell-xps-13.jpg',
    categoryId: '1',
    categorySlug: 'may-tinh-xach-tay',
    stock: 15,
    rating: 4.8
  },
  {
    id: '2',
    name: 'Lenovo ThinkPad X1 Carbon',
    slug: 'lenovo-thinkpad-x1-carbon',
    price: 32990000,
    description: 'Laptop doanh nhân siêu bền',
    imageUrl: '/images/products/lenovo-thinkpad-x1.jpg',
    categoryId: '1',
    categorySlug: 'may-tinh-xach-tay',
    stock: 8,
    rating: 4.7
  },
  // Thêm các sản phẩm khác...
];

export const productService = {
  getProducts: () => {
    return Promise.resolve(products);
  },
  
  getProductsByCategory: (categoryId: string) => {
    return Promise.resolve(products.filter(product => product.categoryId === categoryId));
  },
  
  getProductsByCategorySlug: (categorySlug: string) => {
    return Promise.resolve(products.filter(product => product.categorySlug === categorySlug));
  },
  
  getProductById: (productId: string) => {
    return Promise.resolve(products.find(product => product.id === productId));
  },
  
  getProductBySlug: (productSlug: string) => {
    return Promise.resolve(products.find(product => product.slug === productSlug));
  },
  
  getProductByCategoryAndSlug: (categorySlug: string, productSlug: string) => {
    return Promise.resolve(
      products.find(product => 
        product.categorySlug === categorySlug && product.slug === productSlug
      )
    );
  }
};