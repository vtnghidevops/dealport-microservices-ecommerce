import { Category } from '@/types/category.model';
import { NavigateFunction } from 'react-router-dom';
import { Product } from '@/types/product.model'
export const handleViewAll = (navigate: NavigateFunction) => {
  console.log("View all products clicked");
  navigate('/products');
  // Navigate to all products page
};

export const handleCategoryClick = (navigate: NavigateFunction, category: Category) => {
  console.log("Category clicked:", category);
  navigate(`/category/${category.slug}`);
};

export const handleProductItemClick = (navigate: NavigateFunction, product: Product, index: number) => {
  console.log("Product clicked:", product, index);
  // Make sure we're navigating to a product with its category
  if (product.categorySlug) {
    navigate(`/category/${product.categorySlug}/${product.slug}`);
  } else {
    navigate(`/products/${product.slug}`);
  }
};

export function generateSlug(text: string): string {
  return text
    .toLowerCase()
    .normalize('NFD') // tách dấu ra khỏi kí tự (separate accents from characters)
    .replace(/[\u0300-\u036f]/g, '') // loại bỏ dấu (remove accents)
    .replace(/[đĐ]/g, 'd') // đổi đ thành d (replace đ with d)
    .replace(/[^a-z0-9\s]/g, '') // loại bỏ ký tự đặt biệt (remove special characters)
    .replace(/\s+/g, '-') // thay khoảng trắng bằng dấu gạch ngang (replace spaces with hyphens)
    .replace(/^-+|-+$/g, ''); // loại bỏ dấu gạch ngang ở đầu và cuối (remove hyphens at the beginning and end)
}

export const normalizeText = (text: string): string => {
  return text
    .toLowerCase()
    .trim()
    .replace(/\s+/g, ' ');
};