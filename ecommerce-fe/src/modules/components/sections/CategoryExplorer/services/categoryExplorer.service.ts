import { Category } from './../../../../models/category.models';
import { CategoryItem } from "../models/category.model";

const urlImgExploring = "images/exploring/";
// Mock data - trong thực tế sẽ fetch từ API
const categoryData: CategoryItem[] = [
    {
      id: 1,
      name: "Grocery",
      image: `${urlImgExploring}grocery.png`,
    },
    {
      id: 2,
      name: "Home",
      image: `${urlImgExploring}/home.png`,
    },
    {
      id: 3,
      name: "Fashion",
      image: `${urlImgExploring}/fashion.png`,
    },
    {
      id: 4,
      name: "Electronic",
      image: `${urlImgExploring}electronic.png`,
    },
    {
      id: 5,
      name: "Toys",
      image: `${urlImgExploring}toys.png`,
    },
    {
      id: 6,
      name: "Grocery",
      image: `${urlImgExploring}grocery.png`,
    },
];



// resolve data trendingProduct
export const CategoryExplorerService = {
  getCategoryData: (): Promise<CategoryItem[]> => {
    return Promise.resolve(categoryData);
  }
};