import { BannerShowCaseItem } from "../../components/homepage/Ads/models/ads.model";
import { DisplayItem } from "../../components/homepage/Ads/models/ads.model";
import { GamingItem } from "../../components/homepage/Ads/models/ads.model";
import { NewFashionItem } from "../../components/homepage/Ads/models/ads.model";
import axios from "axios";
import { getApiUrl } from '@/utils/api-config';

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
      const response = await axios.get(getApiUrl('ads/placement/banner'));
      // console.log("item banner....:", response.data.data)

      // Check if data exists before mapping
      if (!response.data || !response.data.data) {
        return [];
      }

      // Transform API data to match frontend models
      return response.data.data.map((item: any) => {
        const uiSettings = parseUISettings(item);

        return {
          id: item.id,
          imageUrl: item.imageUrl,
          buttonType: uiSettings.buttonType,
          buttonText: uiSettings.buttonText,
          hasMore: uiSettings.hasMore || false,
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
      const response = await axios.get(getApiUrl('ads/placement/display'));

      // Check if data exists before mapping
      if (!response.data || !response.data.data) {
        return [];
      }

      // Transform API data to match frontend models
      return response.data.data.map((item: any) => {
        // Parse uiSettings if it's a string
        let uiSettings = item.uiSettings || item.ui_settings;
        if (typeof uiSettings === 'string') {
          try {
            uiSettings = JSON.parse(uiSettings);
          } catch (err) {
            console.error("Error parsing uiSettings:", err);
            uiSettings = {};
          }
        }
        // console.log("-----------------------")
        // console.log("Button type:", uiSettings.button_type);
        // console.log("Discount img:", uiSettings.discount_img);
        // console.log("Price:", item.price);
        // console.log("Custom title:", item.customTitle);
        // console.log("Name:", item.name);
        // console.log("Type:", item.type);
        // console.log("Product slug:", item.slug);
        // console.log("-----------------------")

        return {
          id: item.id,
          imageUrl: item.imageUrl,
          buttonType: uiSettings.button_type,
          price: item.type === 'product' ? item.price : undefined,
          title: item.customTitle || item.name,
          discountImg: uiSettings.discount_img,
          type: item.type,
          productSlug: item.slug,
          categorySlug: item.categorySlug // product type
        };
      });
    } catch (error) {
      console.error("Error fetching display data:", error);
      // Fallback data
      return [];
    }
  }
};

export const GamingService = {
  getDataGaming: async (): Promise<GamingItem[]> => {
    try {
      const response = await axios.get(getApiUrl('ads/placement/gaming'));

      // Check if data exists before mapping
      if (!response.data || !response.data.data) {
        return [];
      }

      // Transform API data to match frontend models
      return response.data.data.map((item: any) => {
        // Use name as subtitle (backend ensures proper data)
        return {
          id: item.id,
          subtitle: item.name,
          imageUrl: item.imageUrl,
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
      const response = await axios.get(getApiUrl('ads/placement/new_fashion'));

      if (response.data && response.data.data && response.data.data.length > 0) {
        const item = response.data.data[0];
        const uiSettings = parseUISettings(item);

        return {
          id: item.id,
          title: "New Year! New Fashion", //item.custom_title || item.name,
          imageUrl: item.imageUrl,
          buttonType: uiSettings.buttonType || "secondary",
          type: item.type,
          productSlug: item.slug,
          categorySlug: item.categorySlug // product type
        };
      }

      // Fallback data if no items returned
      return {
        id: 1,
        title: "New Year! New Fashion",
        imageUrl: `/images/ads/products/new-fashion.png`,
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
        imageUrl: `/images/ads/products/new-fashion.png`,
        buttonType: "secondary",
        type: "category",
        productSlug: "new-fashion",
        categorySlug: "new-fashion"
      };
    }
  }
};



