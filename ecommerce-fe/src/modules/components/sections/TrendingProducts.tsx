import React from "react";
import { getButtonClass } from "../common/Button";
import ProductCard from "../common/Card";

const TrendingProducts = ({ products_trending }) => {
  const urlImgMenCollection = "images/trending/";
  const products = [
    {
      id: 1,
      image: `${urlImgMenCollection}pants.png`, // Replace with your image path
      price: "25.95",
      discount: "20% off",
    },
    {
      id: 2,
      image: `${urlImgMenCollection}shirt.png`, // Replace with your image path
      price: "",
      discount: "",
    },
    {
      id: 3,
      image: `${urlImgMenCollection}hat.png`, // Replace with your image path
      price: "104.0",
      discount: "",
    },
    {
      id: 4,
      image: `${urlImgMenCollection}shoe.png`, // Replace with your image path
      price: "10.56",
      discount: "",
    },
  ];  
  return (
    <div className="h-full">
      <span className="header-2 font-bold p-[10px] block mb-3">
        Trending Products
      </span>
      <div className="h-full w-full flex justify-start items-start">
        {products_trending.map((product) => (
          <ProductCard
            key={product.id}
            title={product.title}
            description={product.description}
            price={product.price}
            originalPrice={product.originalPrice}
            discount={product.discount}
            reviews={product.review}
            imageUrl={product.imageUrl}
          />
        ))}
        {/* Collection for Men */}
        <div>
          <MenCollection products={products}></MenCollection>
        </div>
      </div>
      <div className="absolute right-[5%] top-0">
        <button
          className={`${getButtonClass(
            ""
          )} w-[8rem] h-[3rem] rounded-3xl border border-black`}
        >
          View All
        </button>
      </div>
    </div>
  );
};
const MenCollection = ({ products }) => {
  return (
    <div className="w-[392px] min-h-[435px] mx-[0.5rem] p-6 rounded-[12px] border border-grep-300 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.1),0_4px_6px_-2px_rgba(0,0,0,0.05)] ">
      <h2 className="title mb-4 p-[10px]">Trend collection for men</h2>
      <div className="flex flex-wrap gap-5 mx-3">
        {products.map((product) => (
          <MenCard key={product.id} {...product}></MenCard>
        ))}
      </div>
    </div>
  );
};
const MenCard = ({ image, price, discount }) => {
  // id: 1,
  // image:  `${urlImgMenCollection}pants.png`, // Replace with your image path
  // price: '25.95',
  // discount: '20% off',
  return (
    <div
      className="bg-white rounded-xl overflow-hidden relative w-[166px] h-[171px] 
  shadow-[0_-2px_4px_rgba(0,0,0,0.05),2px_0_5px_rgba(0,0,0,0.05),-2px_0_5px_rgba(0,0,0,0.05),0_4px_6px_rgba(0,0,0,0.2)]"
    >
      {discount && (
        <div className="absolute top-4 left-4">
          <span className="font-medium text-[12px] p-[10px] z-10">
            {discount}
          </span>
        </div>
      )}
      <a href="#" className="flex justify-center items-center w-full h-full">
        <img
          src={image}
          alt="Black pants"
          className="w-full h-full object-contain "
        />
      </a>
      <div className="absolute top-4 right-0 ">
        {price && (
          <span className="text-ocean-green p-[10px] text-[12px] z-10">
            ${price}
          </span>
        )}
      </div>

      <div className="mt-2 absolute bottom-[7%] w-full">
        <div className="rounded-3xl w-full flex justify-center items-center ">
          <button
            className={`${getButtonClass(
              "secondary"
            )} flex justify-center items-center !text-[10px] w-[90px] h-[20px] `}
          >
            Buy Now
          </button>
        </div>
      </div>
    </div>
  );
};

export { TrendingProducts };
