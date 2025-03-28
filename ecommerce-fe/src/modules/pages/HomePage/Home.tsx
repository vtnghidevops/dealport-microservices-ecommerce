import React from "react";
import Header from "../../components/layouts/Header";
import Footer from "../../components/layouts/Footer";
import { HappyCustomers } from "../../components/sections/HappyCustomers";
import { TrendingProducts } from "../../components/sections/TrendingProducts";
import { CategoryExplorer } from "../../components/sections/CategoryExplorer";
import { TopProducts } from "../../components/sections/TopProducts";
import { LimitedDeal } from "../../components/sections/LimitedDeal";
import { HeroBanner } from "../../components/sections/SliderBanner";
import { BannerShowCase } from "../../components/sections/Ads";
import { DisplayGrid } from "../../components/sections/Ads";
import { NewFashion } from "../../components/sections/Ads";
import { GamingBanner } from "../../components/sections/Ads";


export default function Home() {
  return (
    <>
    <Header></Header>
    <div>
      <HeroBanner></HeroBanner>
      <div>
        {/* Ads */}
        <div className="relative flex items-center justify-center w-full h-[550px]">
          <div className="absolute left-[3%] top-0">
            <NewFashion></NewFashion>
          </div>
          <div className="top-0 absolute left-[37%]">
            <GamingBanner></GamingBanner>
          </div>
          <div className="absolute left-[71%] top-0 ">
            {/* image, title, price, discount_img, ctaText, buttonType */}
            <DisplayGrid></DisplayGrid>
          </div>
          <div className="absolute w-full top-[58%]">
            {/* img, buttonType, buttonText, hasMore */}
            <BannerShowCase/>
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
            <LimitedDeal></LimitedDeal>
        </div>

        {/* Happy Customers */}
        <div>
            <HappyCustomers></HappyCustomers>
        </div>
      </div>
    </div>
    <Footer></Footer>
    </>
  );
}
