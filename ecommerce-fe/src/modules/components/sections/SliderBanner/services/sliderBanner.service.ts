import { SliderBannerItem } from "../models/sliderBanner.model";

const urlImgSlider = "images/slidebar/";
// Mock data - trong thực tế sẽ fetch từ API
// title: string;
// discount: number | string;
// image: string;
const sliderBannerData: SliderBannerItem[] = [
  {
    title: "Discover the Latest Deals –",
    discount: "Up to 50% Off!",
    image: `${urlImgSlider}slidebar-1.png`,
  },
];



// resolve data trendingProduct
export const SliderBannerService = {
  getDataSlider: (): Promise<SliderBannerItem[]> => {
    return Promise.resolve(sliderBannerData);
  }
};