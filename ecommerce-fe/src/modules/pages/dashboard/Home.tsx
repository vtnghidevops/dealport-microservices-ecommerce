import React from "react";
import HeroBanner from "../../components/sections/HeroBanner";
import NewFashion from "../../components/sections/NewFashion";
import GamingBanner from "../../components/sections/GamingBanner";
import DisplayGrid from "../../components/sections/DisplayGrid";
import BannerShowCase from "../../components/sections/BannerShowcase";
import { LimitedDeal } from "../../components/sections/LimitedDeal";
import { HappyCustomers } from "../../components/sections/HappyCustomers";
import { TrendingProducts } from "../../components/sections/TrendingProducts";
import { CategoryExplorer } from "../../components/sections/CategoryExplorer";
import { TopProducts } from "../../components/sections/TopProducts";
const urlImgProduct = "images/products/";
const urlImgGaming = "images/gaming/";
const urlImgDisplay = "images/display/";
const urlImgBanner = "images/banner/";
const urlImgTop = "images/bestselling/";
const urlLimited = "images/limited/" ;


// Limited Time Deal
const limitedProducts= [
  {
    id: 1,
    title: "Samsung Galaxy S24",
    description: "Gentle yet effective, our hydrating serum infuses skin with essential moisture while brightening and plumping for a radiant complexion.",
    price: 29.99,
    originalPrice: 39.99,
    discount: 25,
    review: {
      rating: 4.8,
      count: 345,
    },
    imageUrl: `${urlLimited}samsung_s24.png`,
  },
  {
    id: 2,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  },
  {
    id: 3,
    title: "Winter fashion jacket",
    description: "Delivers extreme hydration and helps strengthen the skin barrier. Perfect for both daytime wear and overnight rejuvenation. [[2]]",
    price: 39.99,
    originalPrice: 49.99,
    discount: 20,
    review: {
      rating: 4.9,
      count: 412,
    },
    imageUrl: `${urlLimited}winter_jacket.png`,
  },
  {
    id: 4,
    title: "New Balance 574 Senekers",
    description: "Illuminating serum infused with glow-boosting intelligent botanicals that leave the skin visibly radiant while reducing redness. [[4]]",
    price: 44.99,
    originalPrice: 59.99,
    discount: 25,
    review: {
      rating: 4.6,
      count: 278,
    },
    imageUrl: `${urlLimited}senekers.png`,
  },
  {
    id: 5,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  },
  {
    id: 6,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  },
  {
    id: 7,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  }
];

export default function Home() {
  return (
    <div>
      <HeroBanner></HeroBanner>
      <div>
        {/* Ads */}
        <div className="relative flex items-center justify-center w-full h-[550px]">
          <div className="absolute left-[3%] top-0">
            <NewFashion
              title={"New Year! New Fashion"}
              imgURL={`${urlImgProduct}new-fashion.png`}
              ButtonType={"secondary"}
            ></NewFashion>
          </div>
          <div className="top-0 absolute left-[37%]">
            <GamingBanner
              title="Gaming accessories"
              items={[
                {
                  subtitle: "Headsets",
                  imgURL: `${urlImgGaming}headsets.png`,
                  href: "#",
                },
                {
                  subtitle: "Mouse",
                  imgURL: `${urlImgGaming}mouse.png`,
                  href: "#",
                },
                {
                  subtitle: "Controller",
                  imgURL: `${urlImgGaming}controller.png`,
                  href: "#",
                },
                {
                  subtitle: "Chair",
                  imgURL: `${urlImgGaming}chair.png`,
                  href: "#",
                },
              ]}
            ></GamingBanner>
          </div>
          <div className="absolute left-[71%] top-0 ">
            {/* image, title, price, discount_img, ctaText, buttonType */}
            <DisplayGrid
              products={[
                {
                  image: `${urlImgDisplay}be-winner.png`,
                },
                {
                  image: `${urlImgDisplay}redmi-y3.png`,
                  buttonType: "gradient",
                },
                {
                  image: `${urlImgDisplay}ambilighttv.png`,
                  buttonType: "secondary",
                  price: "750.99",
                  title: "Philips 4K Ambilight TV",
                  discount_img: `${urlImgDisplay}discount_img.png`,
                },
              ]}
            ></DisplayGrid>
          </div>
          <div className="absolute w-full top-[58%]">
            {/* img, buttonType, buttonText, hasMore */}
            <BannerShowCase
              categories={[
                {
                  image: `${urlImgBanner}trousers_fashion.png`,
                  buttonType: "",
                  buttonText: "",
                  hasMore: false,
                },
                {
                  image: `${urlImgBanner}watchmen_fashion.png`,
                  buttonType: "gray",
                  buttonText: "Shop Now",
                  hasMore: false,
                },
                {
                  image: `${urlImgBanner}denim_fashion.png`,
                  buttonType: "",
                  buttonText: "",
                  hasMore: true,
                },
                {
                  image: `${urlImgBanner}dometic.png`,
                  buttonType: "black",
                  buttonText: "Shop Now",
                  hasMore: false,
                },
              ]}
            />
          </div>
        </div>

        {/* Trending Products */}
        <div className="relative w-full h-max pl-[6rem]"> 
          <TrendingProducts></TrendingProducts>
        </div>

        {/* Category */}
        <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
          <CategoryExplorer
          />
        </div>

        {/* Best selling */}
        <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
          <TopProducts></TopProducts>
        </div>

        {/* Limited time deal */}
        <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
            <LimitedDeal products={limitedProducts}></LimitedDeal>
        </div>

        {/* Happy Customers */}
        <div>
            <HappyCustomers></HappyCustomers>
        </div>
      </div>
    </div>
  );
}
