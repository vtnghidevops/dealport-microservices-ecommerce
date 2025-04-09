// src/components/product/services/category.service.ts
import { Category } from '../models/category.model';

const categories: Category[] = [
  {
    id: '1',
    name: 'Grocery',
    slug: 'grocery',
    description: 'Fresh food products and everyday essentials',
    image_url: '/images/categories/grocery.png',
    productCount: 5
  },
  {
    id: '2',
    name: 'Home',
    slug: 'home',
    description: 'Furniture, decor, and home improvement products',
    image_url: '/images/categories/home.png',
    productCount: 5
  },
  {
    id: '3',
    name: 'Fashion',
    slug: 'fashion',
    description: 'Clothing, shoes, and accessories for all ages',
    image_url: '/images/categories/fashion.png',
    productCount: 5
  },
  {
    id: '4',
    name: 'Electronic',
    slug: 'electronic',
    description: 'Latest gadgets and electronic devices',
    image_url: '/images/categories/electronic.png',
    productCount: 5
  },
  {
    id: '5',
    name: 'Toys',
    slug: 'toys',
    description: 'Fun and educational toys for children',
    image_url: '/images/categories/toys.png',
    productCount: 5
  },
  {
    id: '6',
    name: 'Grocery',
    slug: 'grocery-2',
    description: 'More grocery items and specialty foods',
    image_url: '/images/categories/grocery.png',
    productCount: 3
  },
];

export const categoryService = {
  getCategories: () => {
    return Promise.resolve(categories);
  },
  
  getCategoryBySlug: (slug: string) => {
    return Promise.resolve(categories.find(category => category.slug === slug));
  }
};