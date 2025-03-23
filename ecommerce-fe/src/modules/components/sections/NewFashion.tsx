import React from "react";
import { Button } from "../common/Button";

export default function NewFashion({ title, imgURL, ButtonType }) {
  return (
    <div className="bg-white text-black h-[328px] w-[464px] rounded-2xl relative border border-gray-300 mt-[-5%] shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),0px_2px_4px_-2px_rgba(0,0,0,0.05)]">
      <span className="header-2 !text-[22px] absolute block top-[5%] left-[5%]">
        {title}
      </span>
      <div className="rounded-xl relative w-full h-[80%] p-3">
        <img
          src={imgURL}
          className="max-w-[27rem] absolute top-[20%] rounded-xl"
        ></img>
      </div>
      <div className="absolute top-[83%] left-[30%]">
        <Button type={ButtonType} text="Shop now"></Button>
      </div>
    </div>
  );
}
