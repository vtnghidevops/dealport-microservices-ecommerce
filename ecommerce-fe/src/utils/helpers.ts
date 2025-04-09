import { CategoryItem } from './../components/homepage/CategoryExplorer/models/category.model';
import { NavigateFunction } from 'react-router-dom';
export const handleViewAll = (navigate: NavigateFunction) => {
  console.log("View all categories clicked");
  navigate('/categories');
  // Navigate to all categories page
};

export const handleCategoryClick = (navigate: NavigateFunction, category: CategoryItem) => {
  console.log("Category clicked:", category);
  navigate(`/${category.slug}`);
};

export const handleProductClick = (product: string, index: number) => {
  console.log("Product clicked:", product, index);
  // Điều hướng đến trang chi tiết sản phẩm
};

export function generateSlug(text: string): string {
  return text
    .toLowerCase()
    .normalize('NFD') // tách dấu ra khỏi kí tự
    .replace(/[\u0300-\u036f]/g, '') // loại bỏ dấu
    .replace(/[đĐ]/g, 'd') // đổi đ thành d
    .replace(/[^a-z0-9\s]/g, '') // loại bỏ ký tự đặt biệt
    .replace(/\s+/g, '-') // thay khoảng trắng bằng dấu gạch ngang
    .replace(/^-+|-+$/g, ''); // loại bỏ dấu gạch ngang ở đầu và cuối
}