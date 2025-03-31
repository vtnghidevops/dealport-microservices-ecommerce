// components/admin/product/add/components/ProductImages.tsx
import React, { ChangeEvent } from "react";
import { FaImage } from "react-icons/fa6";
import { PiRepeatFill } from "react-icons/pi";
import { CiCirclePlus } from "react-icons/ci";
interface ProductImagesProps {
  selectedImages: string[];
  mainImage: string | null;
  onImageUpload: (e: ChangeEvent<HTMLInputElement>) => void;
  onRemoveImage: (index: number) => void;
  onSetMainImage: (image: string) => void;
}

export const ProductImages: React.FC<ProductImagesProps> = ({
  selectedImages,
  mainImage,
  onImageUpload,
  onRemoveImage,
  onSetMainImage,
}) => {
  return (
    <div className="bg-white rounded-lg p-[1.25rem] shadow-sm mb-6 h-[445px]">
      <h2 className="font-bold text-[22px] text-cyprus">
        Upload Product Image
      </h2>
      <p className="text-[15px] text-cyprus font-bold mt-[1rem]">
        Product Image
      </p>

      <div className="relative mb-4 p-4 flex flex-col items-center">
        {mainImage ? (
          <div className="border border-neutral-300 rounded-lg relative w-full mb-3 h-[266px]">
            <img
              src={mainImage}
              className="w-full h-full object-contain"
              alt="Product"
            />
          </div>
        ) : (
          <div className="border border-neutral-300 rounded-lg h-[266px] w-full  mb-3 flex items-center justify-center bg-white">
            <svg
              className="w-16 h-16 text-gray-300"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        )}

        <div className="left-[5%] bottom-[16%] absolute flex justify-between w-[94px] h-[36px] mb-4">
          <button
            className="text-neutral-500 px-3 py-1 border border-gray-200 rounded-lg bg-gray-50 text-sm flex items-center"
            onClick={() => document.getElementById("file-upload")?.click()}
          >
            <FaImage className="mr-2"></FaImage>
            Browse
            <input
              id="file-upload"
              type="file"
              className="hidden"
              onChange={onImageUpload}
              multiple
              accept="image/*"
            />
          </button>

          {mainImage && (
            <button
              className="ml-[12rem] w-[88px] h-[36px] px-3 py-1 border border-gray-200 rounded-lg bg-white text-sm flex items-center drop-shadow-sm filter"
              onClick={() => {
                onSetMainImage("");
                onRemoveImage(selectedImages.indexOf(mainImage));
              }}
            >
              <PiRepeatFill className="min-w-[14px] h-[14px] mr-1"></PiRepeatFill>
              Replace
            </button>
          )}
        </div>
      </div>
      <div className="flex overflow-x-auto">
        {selectedImages.length > 0 &&
          selectedImages.map((img, idx) => (
            <div
              key={idx}
              className="flex-shrink-0 min-w-[98px] h-[105px] relative rounded border border-neutral-200"
            >
              <img
                src={img}
                className="w-[98px] h-[98px] object-cover rounded"
                alt={`Product ${idx + 1}`}
              />
              <button
                className="absolute w-[16px] h-[16px] top-[10px] right-[10px] border border-neutral-500 bg-white rounded-full p-0.5 transform translate-x-1/3 -translate-y-1/3 shadow"
                onClick={() => onRemoveImage(idx)}
              >
                <svg
                  className=" text-black"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          ))}

        <div
          className="ml-3 flex-shrink-0 min-w-[150px] h-[105px] border border-dashed border-gray-300 rounded flex items-center justify-center cursor-pointer hover:bg-gray-50"
          onClick={() => document.getElementById("file-upload")?.click()}
        >
          <div className="flex items-center flex-col justify-center gap-8">
            <CiCirclePlus className="w-24 h-24 rounded-full bg-ocean-green text-white"></CiCirclePlus>
            <span className="text-[15px] font-medium text-ocean-green">Add Image</span>
          </div>
        </div>
      </div>
    </div>
  );
};