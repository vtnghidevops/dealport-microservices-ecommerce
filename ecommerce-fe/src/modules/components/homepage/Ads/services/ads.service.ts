import React from "react";
import { BannerShowCaseItem } from "../models/ads.model";
import { DisplayItem } from "../models/ads.model";
import { GamingItem } from "../models/ads.model";
import { NewFashionItem } from "../models/ads.model";

const urlImgProduct = "images/ads/products/";
const urlImgGaming = "images/ads/gaming/";
const urlImgDisplay = "images/ads/display/";
const urlImgBanner = "images/ads/banner/";

const bannerData: BannerShowCaseItem[] = [
  {
    id: 1,
    image: `${urlImgBanner}trousers_fashion.png`,
    buttonType: "",
    buttonText: "",
    hasMore: false,
  },
  {
    id: 2,
    image: `${urlImgBanner}watchmen_fashion.png`,
    buttonType: "gray",
    buttonText: "Shop Now",
    hasMore: false,
  },
  {
    id: 3,
    image: `${urlImgBanner}denim_fashion.png`,
    buttonType: "",
    buttonText: "",
    hasMore: true,
  },
  {
    id: 4,
    image: `${urlImgBanner}dometic.png`,
    buttonType: "black",
    buttonText: "Shop Now",
    hasMore: false,
  },
]

const displayData: DisplayItem[] = [
  {
    id: 1,
    image: `${urlImgDisplay}be-winner.png`,
  },
  {
    id: 2,
    image: `${urlImgDisplay}redmi-y3.png`,
    buttonType: "gradient",
  },
  {
    id: 3,
    image: `${urlImgDisplay}ambilighttv.png`,
    buttonType: "secondary",
    price: "750.99",
    title: "Philips 4K Ambilight TV",
    discount_img: `${urlImgDisplay}discount_img.png`,
  },
]

const gamingData: GamingItem[] = [
  {
    id: 1,
    subtitle: "Headsets",
    image: `${urlImgGaming}headsets.png`,
    href: "#",
  },
  {
    id: 2,
    subtitle: "Mouse",
    image: `${urlImgGaming}mouse.png`,
    href: "#",
  },
  {
    id: 3,
    subtitle: "Controller",
    image: `${urlImgGaming}controller.png`,
    href: "#",
  },
  {
    id: 4,
    subtitle: "Chair",
    image: `${urlImgGaming}chair.png`,
    href: "#",
  },
]

const newFashionData: NewFashionItem = {
  id: 1,
  title: "New Year! New Fashion",
  image: `${urlImgProduct}new-fashion.png`,
  buttonType: "secondary"
}

export const BannerService = {
  getBannerData: (): Promise<BannerShowCaseItem[]> => {
    return Promise.resolve(bannerData);
  }
};

export const DisplayService = {
  getDataDisplay: (): Promise<DisplayItem[]> => {
    return Promise.resolve(displayData);
  }
};

export const GamingService = {
  getDataGaming: (): Promise<GamingItem[]> => {
    return Promise.resolve(gamingData);
  }
};

export const NewFashionService = {
  getDataNewFashion: (): Promise<NewFashionItem> => {
    return Promise.resolve(newFashionData);
  }
};



