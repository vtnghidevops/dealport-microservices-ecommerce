import React from "react";
import { getButtonClass } from "../common/Button";

const BannerShowCase = ({categories}) => {
  return (
    <div className="px-4 py-8">
      <div className="flex gap-[20px] w-full overflow-x-auto ml-[2.5rem]">
        {categories.map((category) => (
          <ItemBannerShowCase
            key={category.id}
            image={category.image}
            buttonType={category.buttonType}
            buttonText={category.buttonText}
            hasMore={category.hasMore}
          ></ItemBannerShowCase>
        ))}
      </div>
    </div>
  );
}
const ItemBannerShowCase = ({ image, buttonType, buttonText, hasMore }) => {
  if (buttonText && buttonType) {
    if (buttonType == "gray") {
      return (
        <div className="h-[190px] w-[327px] rounded-[12px] bg-white relative">
          <img
            src={image}
            className="w-full h-full rounded-[12px] object-cover"
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
        <div className="h-[190px] w-[327px] rounded-[12px] bg-white relative">
          <img
            src={image}
            className="w-full h-full rounded-[12px] object-cover"
          ></img>
          <div className="absolute top-[75%] left-[30%]">
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
    }
  } else if (hasMore) {
    return (
      <div className="h-[190px] w-[327px] rounded-[12px] bg-white relative">
        <img
          src={image}
          className="w-full h-full rounded-[12px] object-cover"
        ></img>
        <a
          href="#"
          className="absolute flex items-center justify-center bottom-0 w-full bg-white bg-opacity-50"
        >
          <span className="text-cyprus ">See more</span>
        </a>
      </div>
    );
  } else {
    return (
      <div className="h-[190px] w-[327px] rounded-[12px] bg-white">
        <a href="#" className="block h-full w-full">
          <img
            src={image}
            className="w-full h-full rounded-[12px] object-cover"
            alt="Product image"
          />
        </a>
      </div>
    );
  }
};

export default BannerShowCase;