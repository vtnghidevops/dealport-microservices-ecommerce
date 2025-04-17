import { TopProductItem } from "../models/topProducts.model";

// Mock data - trong thực tế sẽ fetch từ API
const urlImgTop = "/images/bestselling/";
// mongodb
const topProductsData: TopProductItem[] = [
    {
      id: '1', 
      type: "top-sale",
      name: "Computer Accessories",
      image_url: `${urlImgTop}computer_accessories.png`,
      setUpDesign: "row",
      isCommingSoon: true,
      description: "category",
      slug: "computer-accessories",
      price: 0,
      categoryId: "",
      categorySlug: "category",
      stockQuantity: 0,
      gridSpan: {
        col: 1,
        row: 1
      }
    },
    {
      id: '2',
      type: "top-sale",
      name: "Men's Casual Outfit",
      price: 200,
      image_url: `${urlImgTop}football.png`,
      actionLabel: "Visit store",
      setUpDesign: "row",
      gridSpan: {
        col: 1,
        row: 1
      },
      description: "category",
      slug: "category",
      categoryId: "",
      categorySlug: "category",
      stockQuantity: 0,
    },
    {
      id: '3',
      type: "top-sale",
      name: "Pome Granate Juice",
      price: 49,
      image_url: `${urlImgTop}juice.png`,
      actionLabel: "Buy now",
      setUpDesign: "row",
      gridSpan: {
        col: 1,
        row: 1
      },
      description: "category",
      slug: "category",
      categoryId: "",
      categorySlug: "category",
      stockQuantity: 0,
    },
    {
      id: '4',
      type: "top-sale",
      name: "Dog Food Made With Love",
      image_url: `${urlImgTop}dog_food.png`,
      actionLabel: "Shop Now",
      setUpDesign: "col",
      gridSpan: {
        col: 1,
        row: 2
      },
      description: "category",
      slug: "category",
      categoryId: "",
      categorySlug: "category",
      price: 0,
      stockQuantity: 0,
    },
    {
      id: '5',
      type: "top-sale", 
      name: "Security Camera System",
      price: 420,
      image_url: `${urlImgTop}security_camera_left.png`,
      actionLabel: "Visit store",
      setUpDesign: "double",
      gridSpan: {
        col: 2,
        row: 1
      },
      isImageDouble: true,
      image_double_url: `${urlImgTop}security_camera_right.png`,
      description: "category",
      slug: "category",
      categoryId: "",
      categorySlug: "category",
      stockQuantity: 0,
    },
    {
      id: '6',
      type: "top-sale",
      name: "Premium Cosmetic Set",
      price: 149,
      image_url: `${urlImgTop}skincare.png`,
      actionLabel: "Buy Now",
      setUpDesign: "row",
      gridSpan: {
        col: 1,
        row: 1
      },
      description: "category",
      slug: "category",
      categoryId: "",
      categorySlug: "category",
      stockQuantity: 0,
    }
];



// resolve data trendingProduct
export const TopProductsService = {
  getTopProductsData: (): Promise<TopProductItem[]> => {
    return Promise.resolve(topProductsData);
  }
};