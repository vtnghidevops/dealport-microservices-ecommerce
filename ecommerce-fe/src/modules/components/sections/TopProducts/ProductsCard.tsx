import React from "react";
import { ProductCard } from "./models/topProducts.model";
import { RxArrowTopRight } from "react-icons/rx";


const buttonClasses =
  "absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-white flex items-center justify-center";
// Product Card Component
const ProductCardItem: React.FC<ProductCard> = ({ product, onClick }) => {
  const {
    name,
    price,
    image,
    actionLabel,
    setUpDesign,
    isCommingSoon,
    image_double,
  } = product;
  if (setUpDesign == "row") {
    if (isCommingSoon) {
      return (
        <div className="rounded-xl mr-[0.5rem] mt-[0.5rem] border border-gray-200">
          <a href="#">
            <img src={image} alt={name} className="rounded-xl h-[200px] w-[310px]"></img>
          </a>
        </div>
      );
    } else {
      return (
        <div className="rounded-xl relative mr-[0.5rem] mt-[0.5rem]  border border-gray-200">
          <a href="#">
            <img src={image} alt={name} className="rounded-xl h-[200px] w-[310px]"></img>
          </a>
          <span className="absolute bottom-0 p-2 text-white title font-bold right-5">
            ${price}
          </span>
          {actionLabel === "Visit store" ? (
            <button className="absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-white  flex items-center justify-center ">
              <span className="flex justify-center items-center ml-[-0.5rem]">
                {actionLabel}
              </span>
              <div className="bg-white border border-black w-[1.2rem] h-[1.2rem] rounded-[50%] absolute flex items-center justify-center right-2">
                <RxArrowTopRight></RxArrowTopRight>
              </div>
            </button>
          ) : (
            <button className="absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-white  flex items-center justify-center ">
              <span className="flex justify-center items-center">
                {actionLabel}
              </span>
            </button>
          )}
        </div>
      );
    }
  } else if (setUpDesign == "col") {
    return (
      <div className="rounded-xl mr-[0.5rem] relative mt-[0.5rem] ">
        <a href="#" className="block h-full">
          <img
            src={image}
            alt={name}
            className="rounded-xl w-[310px] h-[420px]"
          ></img>
        </a>
        <button className="text-white absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-ocean-blue flex items-center justify-center ">
          <span className="flex justify-center items-center ml-[-0.5rem]">
            {actionLabel}
          </span>
        </button>
      </div>
    );
  } else {
    return (
      <div className="flex rounded-xl mr-[0.5rem] items-center relative mt-[0.5rem] ">
        <a href="#" className="flex max-w-[38rem]">
          <img
            src={image}
            alt={name}
            className="max-w-[19.2rem] aspect-[16/9] h-[200px] w-[310px] rounded-l-xl"
          />
          <img
            src={image_double}
            alt={name}
            className="max-w-[19.2rem] aspect-[16/9] h-[200px] w-[310px] rounded-r-xl"
          />
        </a>
        {actionLabel === "Visit store" ? (
          <button className={buttonClasses}>
            <span className="flex justify-center items-center ml-[-0.5rem]">
              {actionLabel}
            </span>
            <div className="bg-white border border-black w-[1.2rem] h-[1.2rem] rounded-[50%] absolute flex items-center justify-center right-2">
              <RxArrowTopRight />
            </div>
          </button>
        ) : (
          <button className={buttonClasses}>
            <span className="flex justify-center items-center">
              {actionLabel}
            </span>
          </button>
        )}
      </div>
    );
  }
};

export default ProductCardItem;