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
  isUploading?: boolean;
}

// Đơn giản hóa việc xử lý URL hình ảnh
const getDisplayImageUrl = (url: string): string => {
  // Nếu URL đã là URL đầy đủ, giữ nguyên
  if (url && url.startsWith('http')) {
    return url;
  }

  // Nếu URL là đường dẫn tương đối của product image
  if (url && url.startsWith('/api/products/images/')) {
    return `http://localhost:8080${url}`;
  }

  // Trường hợp khác - giữ nguyên URL
  return url;
};

export const ProductImages: React.FC<ProductImagesProps> = ({
  selectedImages,
  mainImage,
  onImageUpload,
  onRemoveImage,
  onSetMainImage,
  isUploading = false,
}) => {
  return (
    <div className="bg-white rounded-lg p-[1.25rem] shadow-sm mb-6 h-auto">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px] flex items-center">
        Upload Product Image
        {isUploading && (
          <span className="ml-2 inline-block w-5 h-5 border-2 border-gray-300 border-t-ocean-green rounded-full animate-spin"></span>
        )}
      </h2>

      <div className="relative mb-4 p-4 flex flex-col items-center">
        {mainImage ? (
          <div className="border border-neutral-300 rounded-lg relative w-full mb-3 aspect-[4/3] overflow-hidden">
            <div className="w-full h-full flex items-center justify-center bg-neutral-50">
              <img
                src={getDisplayImageUrl(mainImage)}
                className="max-w-full max-h-full object-contain"
                alt="Product main image"
                style={{ maxHeight: "100%", maxWidth: "100%" }}
              />
            </div>
          </div>
        ) : (
          <div className="border border-neutral-300 rounded-lg aspect-[4/3] w-full mb-3 flex items-center justify-center bg-neutral-50">
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

        <div className="mt-2 flex justify-start w-full">
          <button
            className="text-neutral-500 px-3 py-1 border border-gray-200 rounded-lg bg-gray-50 text-sm flex items-center mr-2"
            onClick={() => document.getElementById("file-upload")?.click()}
            disabled={isUploading}
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
              disabled={isUploading}
            />
          </button>

          {mainImage && (
            <button
              className="px-3 py-1 border border-gray-200 rounded-lg bg-white text-sm flex items-center drop-shadow-sm filter"
              onClick={() => {
                // Find index of main image in selectedImages array
                const index = selectedImages.indexOf(mainImage);
                if (index !== -1) {
                  onRemoveImage(index);
                }
              }}
              disabled={isUploading}
            >
              <PiRepeatFill className="min-w-[14px] h-[14px] mr-1"></PiRepeatFill>
              Replace
            </button>
          )}
        </div>
      </div>

      <div className="mt-4">
        <p className="text-sm text-gray-500 mb-2">All Images</p>
        <div className="flex flex-wrap gap-2">
          {selectedImages.length > 0 &&
            selectedImages.map((img, idx) => (
              <div
                key={idx}
                className={`relative rounded border ${img === mainImage ? 'border-ocean-green' : 'border-neutral-200'} cursor-pointer`}
                onClick={() => onSetMainImage(img)}
              >
                <div className="w-[80px] h-[80px] flex items-center justify-center bg-neutral-50 overflow-hidden">
                  <img
                    src={getDisplayImageUrl(img)}
                    className="max-w-full max-h-full object-contain"
                    alt={`Product ${idx + 1}`}
                  />
                </div>
                <button
                  className="absolute w-[16px] h-[16px] top-[2px] right-[2px] border border-neutral-500 bg-white rounded-full p-0.5 transform shadow"
                  onClick={(e) => {
                    e.stopPropagation();
                    onRemoveImage(idx);
                  }}
                  disabled={isUploading}
                >
                  <svg
                    className="text-black"
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
                {img === mainImage && (
                  <div className="absolute bottom-0 w-full bg-ocean-green text-white text-[9px] text-center py-0.5">
                    Main
                  </div>
                )}
              </div>
            ))}

          <div
            className={`w-[80px] h-[80px] border border-dashed border-gray-300 rounded flex items-center justify-center ${isUploading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer hover:bg-gray-50'}`}
            onClick={() => !isUploading && document.getElementById("file-upload")?.click()}
          >
            <CiCirclePlus className="w-8 h-8 text-ocean-green"></CiCirclePlus>
          </div>
        </div>
      </div>
    </div>
  );
};