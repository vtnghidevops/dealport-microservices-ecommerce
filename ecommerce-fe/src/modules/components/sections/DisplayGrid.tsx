import React from "react";
import { Button, getButtonClass } from "../common/Button";
const DisplayGrid = ({ products }) => {
  return (
    <div className="h-[328px] w-[392px] mt-[-5%] text-black rounded-2xl relative">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 ">
        {products.map((product) => (
          <DisplayCard key={product.id} {...product} />
        ))}
      </div>
    </div>
  );
};
const DisplayCard = ({ image, title, price, discount_img, buttonType }) => {
  // Kiểm tra các props để quyết định render card nào
  if (title && price) {
    // Nếu có title và price, render CardMain (card đầy đủ thông tin)
    return (
      <div className="w-[392px] h-[165px] rounded-lg overflow-hidden shadow-md transition-all hover:shadow-xl">
        <CardMain
          img={image}
          title={title}
          price={price}
          discount_img={discount_img}
          buttonType={buttonType}
        />
      </div>
    );
  } else if (buttonType) {
    // Nếu có buttonType nhưng không có title/price, render CardSecond
    return (
      <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-md transition-all hover:shadow-xl">
        <CardSecond img={image} buttonType={buttonType} />
      </div>
    );
  } else {
    // Mặc định render CardFirst (chỉ có hình ảnh và link)
    return (
      <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-md transition-all">
        <CardFirst img={image} />
      </div>
    );
  }
};
const CardFirst = ({ img }) => {
  return (
    <a href="#" className="relative h-[100%] block">
      <img src={img} className="absolute rounded-xl "></img>
      <a
        href="#"
        className="absolute text-white text-[10px] bottom-0 ml-[1rem] mb-[0.5rem] "
      >
        More Detail
      </a>
    </a>
  );
};
const CardSecond = ({ img, buttonType }) => {
  return (
    <a href="#" className="block relative">
      <img src={img}></img>
      <div className="absolute bottom-2 left-1">
        <div className="rounded-3xl">
          <button
            className={`${getButtonClass(
              buttonType
            )} text-white w-[5rem] h-[1.5rem] text-[8px]`}
          >
            Shop Now
          </button>
        </div>
      </div>
    </a>
  );
};
const CardMain = ({ img, title, price, discount_img, buttonType }) => {
  return (
    <div className="relative h-full w-full border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]">
      <div className="absolute left-[-1rem] max-w-[14rem]">
        <img src={img}></img>
      </div>
      <div className="absolute max-w-[60px] top-[10%] left-[42%]">
        <img src={discount_img}></img>
      </div>
      <div className="absolute top-[27%] right-[8%]">
        <h3>{title}</h3>
        <span>${price}</span>
      </div>
      <a className="absolute right-[7%] block bottom-5">
        <div className="rounded-3xl">
          <button
            className={`${getButtonClass(
              buttonType
            )} w-[7rem] h-[2rem] text-[10px]`}
          >
            Shop Now
          </button>
        </div>
      </a>
    </div>
  );
};
export default DisplayGrid;
