// export default GamingBanner;
import React, { useState, useEffect } from "react";
import { GamingItem } from "./models/ads.model";
import { GamingService } from "./services/ads.service";

const GamingBanner: React.FC = () => {
  // fetch data
  const [gamingData, setGamingData] = useState<GamingItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // fetch data trending Products
  useEffect(() => {
    const fetchBannerData = async () => {
      try {
        const data = await GamingService.getDataGaming();
        setGamingData(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching testimonials:", error);
        setLoading(false);
      }
    };

    fetchBannerData();
  }, []);
  
  return (
    <div className="bg-white mt-[-5%] text-black h-[328px] w-[464px] rounded-2xl relative border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]">
      <span className="header-2 !text-[22px] absolute block top-[5%] left-[5%]">
        Gaming Accessories
      </span>
      <div className="grid grid-cols-2 gap-12 mt-4 absolute top-[15%] left-[5%]">
        {gamingData.map((item) => (
          <a key={item.id} href={item.href}>
            <SubGamingBanner subtitle={item.subtitle} image={item.image} />
          </a>
        ))}
      </div>
    </div>
  );
};

interface SubGamming {
  subtitle: string;
  image: string;
}

const SubGamingBanner: React.FC<SubGamming> = ({ subtitle, image }) => {
  return (
    <div className="w-[206px] h-[117px] rounded-xl bg-white relative border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]">
      <span className="text-[12px] block absolute z-10 pl-3 pb-1 bottom-0">
        {subtitle}
      </span>
      <div className="h-[100%] flex justify-center items-center overflow-hidden">
        <img
          src={image}
          className="h-[100%] max-w-[8rem] object-contain py-0.5 transition-transform duration-300 ease-out hover:scale-110"
          alt={subtitle}
        ></img>
      </div>
    </div>
  );
};

export default GamingBanner;