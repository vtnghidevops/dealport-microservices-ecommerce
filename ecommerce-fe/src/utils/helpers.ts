import { CategoryItem } from './../components/homepage/CategoryExplorer/models/category.model';
export const handleViewAll = () => {
  console.log("View all categories clicked");
  // Navigate to all categories page
};

export const handleCategoryClick = (category: CategoryItem) => {
  console.log("Category clicked:", category);
  // Navigate or perform actions
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