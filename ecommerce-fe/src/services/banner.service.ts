// Global banner service
import { Banner, BannerType, SliderBannerItem } from '@/types/banner.model';

const urlImgSlider = "/images/slider/";

// Mock slider banner data
const sliderBannerData: SliderBannerItem[] = [
  {
    id: "1",
    title: "Discover the Latest Deals",
    subtitle: "Limited Time Offer",
    description: "Explore our collection of premium products with exclusive discounts",
    discount: "Up to 50% Off!",
    highlight_text: "Shop Now",
    image_url: `${urlImgSlider}sliderbar-1.png`,
    link_url: "/products/deals",
    action_text: "Shop Now",
    isActive: true,
    priority: 1,
    background_color: "#f8f9fa",
    text_color: "#212529",
    animation_type: "fade"
  },
  {
    id: "2",
    title: "New Season Collection",
    subtitle: "Spring/Summer 2024",
    description: "Refresh your wardrobe with our latest arrivals",
    discount: "20% Off for Members",
    highlight_text: "Exclusive",
    image_url: `${urlImgSlider}sliderbar-2.png`,
    link_url: "/collections/summer",
    action_text: "Explore",
    isActive: true,
    priority: 2,
    background_color: "#e9ecef",
    text_color: "#343a40",
    animation_type: "slide"
  },
  {
    id: "3",
    title: "Tech Innovations",
    subtitle: "Next Generation Gadgets",
    description: "Discover cutting-edge technology for your lifestyle",
    discount: "Free Shipping on Orders $100+",
    highlight_text: "New Arrivals",
    image_url: `${urlImgSlider}sliderbar-3.png`,
    link_url: "/collections/tech",
    action_text: "Discover",
    isActive: true,
    priority: 3,
    background_color: "#d8e2dc",
    text_color: "#2b2d42",
    animation_type: "zoom"
  }
];

// Mock general banners data
const bannersData: Banner[] = [
  // {
  //   id: "1",
  //   title: "Summer Sale",
  //   subtitle: "Hot Deals",
  //   description: "Get up to 40% off on all summer essentials",
  //   image_url: "/images/banners/summer-sale.jpg",
  //   link_url: "/promotions/summer",
  //   action_text: "Shop Now",
  //   position: "homepage_top",
  //   isActive: true,
  //   startDate: "2024-06-01",
  //   endDate: "2024-08-31",
  //   priority: 1,
  //   type: BannerType.SEASONAL,
  //   createdAt: "2024-05-15T00:00:00Z",
  //   updatedAt: "2024-05-15T00:00:00Z"
  // },
  // {
  //   id: "2",
  //   title: "New Electronics",
  //   subtitle: "Tech Innovations",
  //   description: "Discover the latest gadgets and electronics",
  //   image_url: "/images/banners/electronics.jpg",
  //   link_url: "/category/electronics",
  //   action_text: "Explore",
  //   position: "homepage_middle",
  //   isActive: true,
  //   priority: 2,
  //   type: BannerType.CATEGORY,
  //   createdAt: "2024-05-10T00:00:00Z",
  //   updatedAt: "2024-05-10T00:00:00Z"
  // }
];

export const bannerService = {
  // Slider banner methods
  getSliderBanners: (): Promise<SliderBannerItem[]> => {
    return Promise.resolve(sliderBannerData);
  },

  getSliderBannerById: (id: string): Promise<SliderBannerItem | undefined> => {
    return Promise.resolve(sliderBannerData.find(banner => banner.id === id));
  },

  // General banner methods
  getBanners: (type?: BannerType, limit?: number): Promise<Banner[]> => {
    let filteredBanners = bannersData;

    if (type) {
      filteredBanners = filteredBanners.filter(banner => banner.type === type);
    }

    // Only return active banners by default
    filteredBanners = filteredBanners.filter(banner => banner.isActive);

    // Sort by priority (lower number = higher priority)
    filteredBanners = filteredBanners.sort((a, b) =>
      (a.priority || 999) - (b.priority || 999)
    );

    if (limit) {
      filteredBanners = filteredBanners.slice(0, limit);
    }

    return Promise.resolve(filteredBanners);
  },

  getBannerById: (id: string): Promise<Banner | undefined> => {
    return Promise.resolve(bannersData.find(banner => banner.id === id));
  }
}; 