import { TopProductItem } from "../models/topProducts.model";

// Mock data - trong thực tế sẽ fetch từ API
const urlImgTop = "images/bestselling/";
const topProductsData: TopProductItem[] = [
    {
      id: 1,
      name: "Computer Accessories",
      image: `${urlImgTop}computer_accessories.png`,
      setUpDesign: "row",
      isCommingSoon: true,
      gridSpan: {
        col: 1,
        row: 1
      }
    },
    {
      id: 2,
      name: "Men's Casual Outfit",
      price: "200",
      image: `${urlImgTop}football.png`,
      actionLabel: "Visit store",
      setUpDesign: "row",
      gridSpan: {
        col: 1,
        row: 1
      }
    },
    {
      id: 3,
      name: "Pome Granate Juice",
      price: "49",
      image: `${urlImgTop}juice.png`,
      actionLabel: "Buy now",
      setUpDesign: "row",
      gridSpan: {
        col: 1,
        row: 1
      }
      
    },
    {
      id: 4,
      name: "Dog Food Made With Love",
      image: `${urlImgTop}dog_food.png`,
      actionLabel: "Shop Now",
      setUpDesign: "col",
      gridSpan: {
        col: 1,
        row: 2
      },
    },
    {
      id: 5,
      name: "Security Camera System",
      price: "420",
      image: `${urlImgTop}security_camera_left.png`,
      actionLabel: "Visit store",
      setUpDesign: "double",
      gridSpan: {
        col: 2,
        row: 1
      },
      isImageDouble: true,
      image_double: `${urlImgTop}security_camera_right.png`,
    },
    {
      id: 6,
      name: "Premium Cosmetic Set",
      price: "149",
      image: `${urlImgTop}skincare.png`,
      actionLabel: "Buy Now",
      setUpDesign: "row",
      gridSpan: {
        col: 1,
        row: 1
      }
  
    }
];



// resolve data trendingProduct
export const TopProductsService = {
  getTopProductsData: (): Promise<TopProductItem[]> => {
    return Promise.resolve(topProductsData);
  }
};