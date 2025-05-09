import { useState, useEffect } from "react";
import { HappyCustomers } from "@/components/homepage/HappyCustomers";
import { TrendingProducts } from "@/components/homepage/TrendingProducts";
import { CategoryExplorer } from "@/components/homepage/CategoryExplorer";
import TopProducts from "@/components/homepage/BestSelling/TopProducts";
import { LimitedDeal } from "@/components/homepage/LimitedDeal";
import { HeroBanner } from "@/components/homepage/SliderBanner";
import { BannerShowCase } from "@/components/homepage/Ads";
import { DisplayGrid } from "@/components/homepage/Ads";
import { NewFashion } from "@/components/homepage/Ads";
import { GamingBanner } from "@/components/homepage/Ads";
import {
  SliderSkeleton,
  NewFashionSkeleton,
  GamingBannerSkeleton,
  DisplayGridSkeleton,
  BannerShowCaseSkeleton,
  CategorySkeleton,
  TopProductSkeleton,
  HappyCustomerSkeleton,
  ProductCardSkeletonGrid
} from "@/components/ui/skeletons";

export default function Home() {
  const [loading, setLoading] = useState(true);

  // Simulate loading data
  useEffect(() => {
    // Force the page to show skeletons first
    setLoading(true);

    // Simulate data loading with timeout
    const timer = setTimeout(() => {
      setLoading(false);
    }, 3000);

    return () => clearTimeout(timer);
  }, []);

  if (loading) {
    return (
      <>
        <div>
          <SliderSkeleton />
          <div>
            {/* Ads Skeletons */}
            <div className="relative flex items-center justify-center w-full h-[550px]">
              <div className="absolute left-[3%] top-[-5%]">
                <NewFashionSkeleton />
              </div>
              <div className="top-0 absolute left-[37%]">
                <GamingBannerSkeleton />
              </div>
              <div className="absolute left-[71%] top-0">
                <DisplayGridSkeleton />
              </div>
              <div className="absolute w-full top-[58%]">
                <BannerShowCaseSkeleton />
              </div>
            </div>

            {/* Trending Products Skeleton */}
            <div className="relative w-full h-max pl-[6rem]">
              <TopProductSkeleton />
            </div>

            {/* Category Skeleton */}
            <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
              <CategorySkeleton count={6} variant="homepage" />
            </div>

            {/* Best selling Skeleton */}
            <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
              <TopProductSkeleton />
            </div>

            {/* Limited time deal Skeleton */}
            <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
              <ProductCardSkeletonGrid count={5} />
            </div>

            {/* Happy Customers Skeleton */}
            <div>
              <div className="border border-gray-200 rounded-xl overflow-hidden">
                <HappyCustomerSkeleton />
              </div>
            </div>
          </div>
        </div>
      </>
    );
  }

  return (
    <>
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
              <BannerShowCase />
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
            <TopProducts products={[]} />
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
    </>
  );
}
