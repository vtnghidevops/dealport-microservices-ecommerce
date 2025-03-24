import React, {useState, useEffect} from "react";
import { Button } from "../../common/Button";
import { NewFashionItem } from "./models/ads.model";
import { NewFashionService } from "./services/ads.service";

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

  if (loading || !newFashionData) {
    return <div>Loading...</div>;
  }
  const { title, image, buttonType } = newFashionData;
  
  return (
    <div className="bg-white text-black h-[328px] w-[464px] rounded-2xl relative border border-gray-300 mt-[-5%] shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),0px_2px_4px_-2px_rgba(0,0,0,0.05)]">
      <span className="header-2 !text-[22px] absolute block top-[5%] left-[5%]">
        {title}
      </span>
      <div className="rounded-xl relative w-full h-[80%] p-3">
        <img
          src={image}
          className="max-w-[27rem] absolute top-[20%] rounded-xl"
        ></img>
      </div>
      <div className="absolute top-[83%] left-[30%]">
        <Button type={buttonType} text="Shop now"></Button>
      </div>
    </div>
  );
};

export default NewFashion;
