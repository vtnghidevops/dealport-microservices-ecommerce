import { SliderBannerItem } from "../models/sliderBanner.model";

const urlImgSlider = "images/slidebar/";
// Mock data - trong thực tế sẽ fetch từ API
// title: string;
// discount: number | string;
// image: string;
const sliderBannerData: SliderBannerItem[] = [
  {
    id: "1",
    title: "Discover the Latest Deals –",
    discount: "Up to 50% Off!",
    image_url: `${urlImgSlider}slidebar-1.png`,
  },
];



// resolve data trendingProduct
export const SliderBannerService = {
  getDataSlider: (): Promise<SliderBannerItem[]> => {
    return Promise.resolve(sliderBannerData);
  }
};