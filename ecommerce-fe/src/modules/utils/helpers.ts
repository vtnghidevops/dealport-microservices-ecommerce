import { CategoryItem } from "../components/sections/CategoryExplorer/models/category.model";
export const handleViewAll = () => {
  console.log("View all categories clicked");
  // Navigate to all categories page
};

export const handleCategoryClick = (category: CategoryItem) => {
  console.log("Category clicked:", category);
  // Navigate or perform actions
};