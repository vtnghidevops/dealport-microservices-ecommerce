import { CategoryItem } from "../models/category.model";

const urlImgExploring = "images/categories/";
// Mock data - trong thực tế sẽ fetch từ API
const categoryData: CategoryItem[] = [
    {
      id: '1',
      name: "Grocery",
      image_url: `${urlImgExploring}grocery.png`,
      slug: "grocery"
    },
    {
      id: '2',
      name: "Home",
      image_url: `${urlImgExploring}home.png`,
      slug: "home"
    },
    {
      id: '3',
      name: "Fashion",
      image_url: `${urlImgExploring}fashion.png`,
      slug: "fashion"
    },
    {
      id: '4',
      name: "Electronic",
      image_url: `${urlImgExploring}electronic.png`,
      slug: "electronic"
    },
    {
      id: '5',
      name: "Toys",
      image_url: `${urlImgExploring}toys.png`,
      slug: "toys"
    },
    {
      id: '6',
      name: "Grocery",
      image_url: `${urlImgExploring}grocery.png`,
      slug: "grocery-2"
    },
];



// resolve data trendingProduct
export const CategoryExplorerService = {
  getCategoryData: (): Promise<CategoryItem[]> => {
    return Promise.resolve(categoryData);
  }
};