import React, {useState, useEffect} from "react";
import { ButtonType, getButtonClass } from "@/utils/buttonUtils";
import { BannerShowCaseItem } from "./models/ads.model";
import { BannerService } from "./services/ads.service";

const BannerShowCase: React.FC = () => {
    // fetch data
    const [bannerData, setBannerData] = useState<BannerShowCaseItem[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    
      // fetch data trending Products
      useEffect(() => {
        const fetchBannerData = async () => {
          try {
            const data = await BannerService.getBannerData();
            setBannerData(data);
            setLoading(false);
          } catch (error) {
            console.error("Error fetching testimonials:", error);
            setLoading(false);
          }
        };
    
        fetchBannerData();
      }, []);

  return (
    <div className="px-4 py-8">
      <div className="flex gap-[20px] w-full overflow-x-auto ml-[2.5rem]">
        {bannerData.map((category, index) => (
          <ItemBannerShowCase
            key={index}
            id={category.id}
            image_url={category.image_url}
            buttonType={category.buttonType}
            buttonText={category.buttonText}
            hasMore={category.hasMore}
          ></ItemBannerShowCase>
        ))}
      </div>
    </div>
  );
}

const ItemBannerShowCase: React.FC<BannerShowCaseItem> = ({ image_url, buttonType, buttonText, hasMore }) => {
  // Check if button text is "Shop Now" to disable image hover scale
  const isShopNow = buttonText === "Shop Now";
  
  if (buttonText && buttonType) {
    if (buttonType == "gray") {
      return (
        <div className="h-[190px] w-[327px] rounded-[12px] bg-white relative overflow-hidden cursor-pointer">
          <img
            src={image_url}
            className={`w-full h-full rounded-[12px] object-cover ${!isShopNow ? "transition-transform duration-300 ease-out hover:scale-110" : ""}`}
            alt="Banner image"
          ></img>
          <div className="absolute top-[70%] left-[8%]">
            <div className="w-[8rem] h-[2rem] rounded-3xl ">
              <button
                className={`${getButtonClass(
                  buttonType
                )} w-full h-full text-[12px]`}
              >
                {buttonText}
              </button>
            </div>
          </div>
        </div>
      );
    } else {
      return (
        <div className="h-[190px] w-[327px] rounded-[12px] bg-white relative overflow-hidden cursor-pointer">
          <img
            src={image_url}
            className={`w-full h-full rounded-[12px] object-cover ${!isShopNow ? "transition-transform duration-300 ease-out hover:scale-110" : ""}`}
            alt="Banner image"
          ></img>
          <div className="absolute top-[75%] left-[30%]">
            <div className="w-[8rem] h-[2rem] rounded-3xl ">
              <button
                className={`${getButtonClass(
                  buttonType as ButtonType
                )} w-full h-full text-[12px]`}
              >
                {buttonText}
              </button>
            </div>
          </div>
        </div>
      );
    }
  } else if (hasMore) {
    return (
      <div className="h-[190px] w-[327px] rounded-[12px] bg-white relative overflow-hidden cursor-pointer">
        <img
          src={image_url}
          className="w-full h-full rounded-[12px] object-cover transition-transform duration-300 ease-out hover:scale-110"
          alt="Banner image"
        ></img>
        <a
          href="#"
          className="absolute flex items-center justify-center bottom-0 w-full bg-white bg-opacity-50 py-2 transition-all duration-300 hover:bg-opacity-70"
        >
          <span className="text-cyprus transition-all duration-300 hover:font-semibold">See more</span>
        </a>
      </div>
    );
  } else {
    return (
      <div className="h-[190px] w-[327px] rounded-[12px] bg-white overflow-hidden cursor-pointer">
        <a href="#" className="block h-full w-full">
          <img
            src={image_url}
            className="w-full h-full rounded-[12px] object-cover transition-transform duration-300 ease-out hover:scale-110"
            alt="Banner image"
          />
        </a>
      </div>
    );
  }
};

export default BannerShowCase;