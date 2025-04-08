import { MenCollectionItem } from "../models/trendingProducts.model";

const urlImgMenCollection = "images/trending/";
// Mock data - trong thực tế sẽ fetch từ API
const menCollectionData: MenCollectionItem[] = [
  {
    id: '1',
    image_url: `${urlImgMenCollection}pants.png`, // Replace with your image path
    price: 25.95,
    discount: 20,
  },
  {
    id: '2',
    image_url: `${urlImgMenCollection}shirt.png`, // Replace with your image path
    price: 0,
    discount: 0,
  },
  {
    id: '3',
    image_url: `${urlImgMenCollection}hat.png`, // Replace with your image path
    price: 104.0,
    discount: 0,
  },
  {
    id: '4',
    image_url: `${urlImgMenCollection}shoe.png`, // Replace with your image path
    price: 10.56,
    discount: 0,
  },
];



// resolve data trendingProduct
export const MenCollectionService = {
  getDataMen: (): Promise<MenCollectionItem[]> => {
    return Promise.resolve(menCollectionData);
  }
};