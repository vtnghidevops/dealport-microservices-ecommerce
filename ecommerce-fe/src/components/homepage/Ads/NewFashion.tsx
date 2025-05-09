import React, { useState, useEffect } from "react";
import { getButtonClass, ButtonType } from "../../../utils/buttonUtils";
import { NewFashionItem } from "./models/ads.model";
import { NewFashionService } from "../../../services/product/ads.service";
import { Link } from "react-router-dom";
import { NewFashionSkeleton } from "@/components/ui/skeletons";
const NewFashion: React.FC = () => {
  // fetch data
  const [newFashionData, setNewFashionData] = useState<NewFashionItem>();
  const [loading, setLoading] = useState<boolean>(true);

  // fetch data trending Products
  useEffect(() => {
    const fetchBannerData = async () => {
      try {
        const data = await NewFashionService.getDataNewFashion();
        setNewFashionData(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching testimonials:", error);
        setLoading(false);
      }
    };

    fetchBannerData();
  }, []);
  // console.log("item newfashion:", newFashionData)


  if (loading || !newFashionData) {
    return <NewFashionSkeleton count={1} />;
  }
  const { title, imageUrl, buttonType, type, productSlug, categorySlug } = newFashionData;

  // Determine the correct navigation path based on item type and available data
  let navigationPath = "#";
  if (type === "product" && productSlug) {
    if (categorySlug) {
      navigationPath = `/category/${categorySlug}/${productSlug}`;
    } else {
      navigationPath = `/products/${productSlug}`;
    }
  } else if (type === "category" && productSlug) {
    navigationPath = `/category/${categorySlug}`;
  }

  return (
    <div className="bg-white text-black h-[328px] w-[464px] rounded-2xl relative border border-gray-300 mt-[-5%] shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),0px_2px_4px_-2px_rgba(0,0,0,0.05)]">
      <span className="header-2 !text-[22px] absolute block top-[5%] left-[5%]">
        {title}
      </span>
      <Link to={navigationPath}>
        <div className="rounded-xl relative w-full h-[80%] p-3">
          <img
            src={imageUrl}
            className="max-w-[27rem] absolute top-[20%] rounded-xl"
            alt={title}
          />
        </div>
        <div className="absolute top-[83%] left-[30%]">
          <div className="w-[12rem] h-[4rem] rounded-3xl ">
            <button className={getButtonClass(buttonType as ButtonType)}>
              Shop Now
            </button>
          </div>
        </div>
      </Link>
    </div>
  );
};

export default NewFashion;
