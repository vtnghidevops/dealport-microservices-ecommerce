import React from "react";

export default function GamingBanner ({title, items}){
  return (
    <div className="bg-white mt-[-5%] text-black h-[328px] w-[464px] rounded-2xl relative border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]" >
      <span className="header-2 !text-[22px] absolute block top-[5%] left-[5%]">{title}</span>
      <div className="grid grid-cols-2 gap-12 mt-4 absolute top-[15%] left-[5%]">
        {items.map((item, index) => (
          <a key={index} href={item.href}>
            <SubGamingBanner
              subtitle={item.subtitle}
              imgURL={item.imgURL}
            />
          </a>
        ))}
      </div>
    </div>
  );
}
const SubGamingBanner = ({subtitle, imgURL}) => {
  return (
    <div className="w-[206px] h-[117px] rounded-xl  bg-white relative border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]">
      <span className="text-[12px] block absolute z-10 pl-3 pb-1 bottom-0">{subtitle}</span>
      <div className="h-[100%] flex justify-center items-center">
        <img src={imgURL} className="h-[100%] max-w-[8rem] object-contain py-0.5"></img>
      </div>
    </div>
  );
}