// components/admin/product/ProductManagement.tsx
import React, { useState } from 'react';
import { ProductHeader } from "./components/ProductHeader";
import { BasicDetails } from "./components/BasicDetails";
import { ProductPricing } from "./components/ProductPricing";
import { ProductInventory } from "./components/ProductInventory";
import { ProductImages } from "./components/ProductImages";
import { ProductCategories } from "./components/ProductCategories";
import { useProductForm } from "./hooks/useProductForm";
import { useImageUpload } from "./hooks/useImageUpload";
import AdminHeader from '../layout/AdminHeader';

export const ProductManagement: React.FC = () => {
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
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Add Product" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            {/* <AddProduct /> */}
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
          </div>
        </main>
      </div>
    </div>
  );
};