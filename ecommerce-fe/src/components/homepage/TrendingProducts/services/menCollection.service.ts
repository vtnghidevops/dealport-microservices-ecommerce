import { Product } from "@/types/product.model";

const urlImgMenCollection = "/images/trending/";
// Mock data - trong thực tế sẽ fetch từ API
const menCollectionData: Product[] = [
  {
    type: "trending",
    id: '1',
    image_url: `${urlImgMenCollection}pants.png`, // Replace with your image path
    price: 25.95,
    discount: 20,
    name: "Pants",
    description: "Pants",
    slug: "pants",
    categoryId: "1",
    categorySlug: "men",
    stockQuantity: 100,
  },
  {
    type: "trending",
    id: '2',
    image_url: `${urlImgMenCollection}shirt.png`, // Replace with your image path
    price: 0,
    discount: 0,
    name: "Shirt",
    description: "Shirt",
    slug: "shirt",
    categoryId: "1",
    categorySlug: "men",
    stockQuantity: 100,
  },
  {
    type: "trending",
    id: '3',
    image_url: `${urlImgMenCollection}hat.png`, // Replace with your image path
    price: 104.0,
    discount: 0,
    name: "Hat",
    description: "Hat",
    slug: "hat",
    categoryId: "1",
    categorySlug: "men",
    stockQuantity: 100,
  },
  {
    type: "trending",
    id: '4',
    image_url: `${urlImgMenCollection}shoe.png`, // Replace with your image path
    price: 10.56,
    discount: 0,
    name: "Shoe",
    description: "Shoe",
    slug: "shoe",
    categoryId: "1",
    categorySlug: "men",
    stockQuantity: 100,
  },
];



// resolve data trendingProduct
export const MenCollectionService = {
  getDataMen: (): Promise<Product[]> => {
    return Promise.resolve(menCollectionData);
  }
};