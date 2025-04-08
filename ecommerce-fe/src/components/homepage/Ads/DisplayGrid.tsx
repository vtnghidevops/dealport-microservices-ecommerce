import React, { useState, useEffect } from "react";
import { ButtonType, getButtonClass } from "@/utils/buttonUtils";
import { DisplayItem } from "./models/ads.model";
import { DisplayService } from "./services/ads.service";

interface CardMain {
  id?: number;
  image_url: string;
  buttonType: string;
  price?: string | number;
  title?: string;
  discount_img?: string;
}

interface CardFirst {
  image_url: string;
}

interface CardSecond {
  image_url: string;
  buttonType: string;
}

const DisplayGrid: React.FC = () => {
  // fetch data
  const [displayData, setDisplayData] = useState<DisplayItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // fetch data trending Products
  useEffect(() => {
    const fetchDisplayData = async () => {
      try {
        const data = await DisplayService.getDataDisplay();
        setDisplayData(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching: ", error);
        setLoading(false);
      }
    };
    fetchDisplayData();
  }, []);

  return (
    <div className="h-[328px] w-[392px] mt-[-5%] text-black rounded-2xl relative">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 ">
        {displayData.map((product, index) => (
          <DisplayCard key={index} {...product} />
        ))}
      </div>
    </div>
  );
};
const DisplayCard: React.FC<DisplayItem> = ({
  image_url,
  title,
  price,
  discount_img,
  buttonType,
}) => {
  // Kiểm tra các props để quyết định render card nào
  if (title && price) {
    // Nếu có title và price, render CardMain (card đầy đủ thông tin)
    return (
      <div className="w-[392px] h-[165px] rounded-lg overflow-hidden shadow-md transition-all hover:shadow-xl">
        <CardMain
          image_url={image_url}
          title={title}
          price={price}
          discount_img={discount_img}
          buttonType={buttonType || ""}
        />
      </div>
    );
  } else if (buttonType) {
    // Nếu có buttonType nhưng không có title/price, render CardSecond
    return (
      <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-md transition-all hover:shadow-xl">
        <CardSecond image_url={image_url} buttonType={buttonType} />
      </div>
    );
  } else {
    // Mặc định render CardFirst (chỉ có hình ảnh và link)
    return (
      <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-md transition-all">
        <CardFirst image_url={image_url} />
      </div>
    );
  }
};

const CardFirst: React.FC<CardFirst> = ({ image_url }) => {
  return (
    <a href="#" className="relative h-[100%] block">
      <img src={image_url} className="absolute rounded-xl "></img>
      <a
        href="#"
        className="absolute text-white text-[10px] bottom-0 ml-[1rem] mb-[0.5rem] "
      >
        More Detail
      </a>
    </a>
  );
};
const CardSecond: React.FC<CardSecond> = ({ image_url, buttonType }) => {
  return (
    <a href="#" className="block relative">
      <img src={image_url}></img>
      <div className="absolute bottom-2 left-1">
        <div className="rounded-3xl">
          <button
            className={`${getButtonClass(
              buttonType as ButtonType || ""
            )} text-white w-[5rem] h-[1.5rem] text-[8px]`}
          >
            Shop Now
          </button>
        </div>
      </div>
    </a>
  );
};
const CardMain: React.FC<CardMain> = ({
  image_url,
  title,
  price,
  discount_img,
  buttonType,
}) => {
  return (
    <div className="relative h-full w-full border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]">
      <div className="absolute left-[-1rem] max-w-[14rem]">
        <img src={image_url}></img>
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
              buttonType as ButtonType
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
