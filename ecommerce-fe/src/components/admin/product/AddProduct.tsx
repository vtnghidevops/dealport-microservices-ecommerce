// components/admin/product/add/AddProduct.tsx
import React from "react";
import { ProductHeader } from "./components/ProductHeader";
import { BasicDetails } from "./components/BasicDetails";
import { ProductPricing } from "./components/ProductPricing";
import { ProductInventory } from "./components/ProductInventory";
import { ProductImages } from "./components/ProductImages";
import { ProductCategories } from "./components/ProductCategories";
import { useProductForm } from "./hooks/useProductForm";
import { useImageUpload } from "./hooks/useImageUpload";

export const AddProduct: React.FC = () => {
  const {
    product,
    handleInputChange,
    handleRadioChange,
    handleColorSelect,
    handlePublish,
    handleSaveDraft,
    calculateSavings
  } = useProductForm();

  const {
    selectedImages,
    mainImage,
    handleImageUpload,
    removeImage,
    setMainImage
  } = useImageUpload();

  return (
    <div className="container mx-auto">
      <ProductHeader 
        onPublish={handlePublish} 
        onSaveDraft={handleSaveDraft} 
      />

      <div className="flex flex-wrap -mx-3 mt-[1.5rem] drop-shadow-sm filter">
        <div className="w-[611px] h-auto md:w-7/12 px-3 mb-6">
          <BasicDetails 
            product={product} 
            onChange={handleInputChange} 
          />
          
          <ProductPricing 
            product={product}
            onChange={handleInputChange}
            discountAmount={calculateSavings()}
          />
          
          <ProductInventory 
            product={product}
            onChange={handleInputChange}
            onRadioChange={handleRadioChange}
            handlePublish={handlePublish}
            handleSaveDraft={handleSaveDraft}
          />

        </div>

        <div className="w-full md:w-5/12 px-3">
          <ProductImages 
            selectedImages={selectedImages}
            mainImage={mainImage}
            onImageUpload={handleImageUpload}
            onRemoveImage={removeImage}
            onSetMainImage={setMainImage}
          />
          
          <ProductCategories 
            product={product}
            onColorSelect={handleColorSelect}
          />
        </div>
      </div>
    </div>
  );
};