// src/components/product/services/category.service.ts
import { Category } from '../models/category.model';

// Mock data
const categories: Category[] = [
  {
    id: '1',
    name: 'Máy tính xách tay',
    slug: 'may-tinh-xach-tay',
    description: 'Các loại máy tính xách tay chất lượng cao',
    imageUrl: '/images/categories/laptop.jpg',
    productCount: 24
  },
  // Thêm danh mục khác...
];

export const categoryService = {
  getCategories: () => {
    return Promise.resolve(categories);
  },
  
  getCategoryBySlug: (slug: string) => {
    return Promise.resolve(categories.find(category => category.slug === slug));
  }
};