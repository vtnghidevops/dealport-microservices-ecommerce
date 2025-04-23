import { BannerShowCaseItem } from "../components/homepage/Ads/models/ads.model";
import { DisplayItem } from "../components/homepage/Ads/models/ads.model";
import { GamingItem } from "../components/homepage/Ads/models/ads.model";
import { NewFashionItem } from "../components/homepage/Ads/models/ads.model";
import axios from "axios";

const API_URL = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || "http://localhost:8082/api/v1";

// Helper function to parse UI settings JSON safely
function parseUISettings(data: any): any {
  if (!data || !data.ui_settings) return {};

  try {
    if (typeof data.ui_settings === 'string') {
      return JSON.parse(data.ui_settings);
    } else {
      return data.ui_settings;
    }
  } catch (error) {
    console.error("Error parsing UI settings:", error);
    return {};
  }
}

export const BannerService = {
  getBannerData: async (): Promise<BannerShowCaseItem[]> => {
    try {
      const response = await axios.get(`${API_URL}/ads/placement/banner`);
      console.log("item banner....:", response.data.data)

      // Transform API data to match frontend models
      return response.data.data.map((item: any) => {
        const uiSettings = parseUISettings(item);

        return {
          id: item.id,
          image_url: item.image_url,
          buttonType: uiSettings.button_type,
          buttonText: uiSettings.button_text,
          hasMore: uiSettings.has_more || false,
          type: item.type, // category type
          categorySlug: item.slug // product type
        };
      });
    } catch (error) {
      console.error("Error fetching banner data:", error);
      // Fallback data
      return [];
    }
  }
};

export const DisplayService = {
  getDataDisplay: async (): Promise<DisplayItem[]> => {
    try {
      const response = await axios.get(`${API_URL}/ads/placement/display`);

      // Transform API data to match frontend models
      return response.data.data.map((item: any) => {
        const uiSettings = parseUISettings(item);

        return {
          id: item.id,
          image_url: item.image_url,
          buttonType: uiSettings.button_type,
          price: item.type === 'product' ? item.price : undefined,
          title: item.custom_title || item.name,
          discount_img: uiSettings.discount_img,
          type: item.type,
          productSlug: item.slug,
          categorySlug: item.category_slug // product type
        };
      });
    } catch (error) {
      console.error("Error fetching display data:", error);
      // Fallback data
      return [
        // {
        //   id: 1,
        //   image_url: `/images/ads/display/be-winner.png`,
        // },
        // {
        //   id: 2,
        //   image_url: `/images/ads/display/redmi-y3.png`,
        //   buttonType: "gradient",
        // },
        // {
        //   id: 3,
        //   image_url: `/images/ads/display/ambilighttv.png`,
        //   buttonType: "secondary",
        //   price: "750.99",
        //   title: "Philips 4K Ambilight TV",
        //   discount_img: `/images/ads/display/discount_img.png`,
        // }
      ];
    }
  }
};

export const GamingService = {
  getDataGaming: async (): Promise<GamingItem[]> => {
    try {
      const response = await axios.get(`${API_URL}/ads/placement/gaming`);

      // Transform API data to match frontend models
      return response.data.data.map((item: any) => {
        // Use name as subtitle (backend ensures proper data)
        return {
          id: item.id,
          subtitle: item.name,
          image_url: item.image_url,
          type: item.type,
          categorySlug: item.slug // category type
        };
      });
    } catch (error) {
      console.error("Error fetching gaming data:", error);
      // Fallback data
      return [
      ];
    }
  }
};

export const NewFashionService = {
  getDataNewFashion: async (): Promise<NewFashionItem> => {
    try {
      const response = await axios.get(`${API_URL}/ads/placement/new_fashion`);

      if (response.data && response.data.data && response.data.data.length > 0) {
        const item = response.data.data[0];
        const uiSettings = parseUISettings(item);

        return {
          id: item.id,
          title: item.custom_title || item.name,
          image_url: item.image_url,
          buttonType: uiSettings.button_type || "secondary",
          type: item.type,
          productSlug: item.slug,
          categorySlug: item.category_slug // product type
        };
      }

      // Fallback data if no items returned
      return {
        id: 1,
        title: "New Year! New Fashion",
        image_url: `/images/ads/products/new-fashion.png`,
        buttonType: "secondary",
        type: "category",
        productSlug: "new-fashion",
        categorySlug: "new-fashion"
      };
    } catch (error) {
      console.error("Error fetching new fashion data:", error);
      // Fallback data
      return {
        id: 1,
        title: "New Year! New Fashion",
        image_url: `/images/ads/products/new-fashion.png`,
        buttonType: "secondary",
        type: "category",
        productSlug: "new-fashion",
        categorySlug: "new-fashion"
      };
    }
  }
};



