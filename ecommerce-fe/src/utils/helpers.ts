import { CategoryItem } from './../components/homepage/CategoryExplorer/models/category.model';
export const handleViewAll = () => {
  console.log("View all categories clicked");
  // Navigate to all categories page
};

export const handleCategoryClick = (category: CategoryItem) => {
  console.log("Category clicked:", category);
  // Navigate or perform actions
};

export const handleProductClick = (product, index) => {
  console.log("Product clicked:", product);
  // Điều hướng đến trang chi tiết sản phẩm
};