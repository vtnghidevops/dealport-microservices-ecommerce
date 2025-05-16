// components/admin/product/ProductManagement.tsx
import React, { useEffect, useState } from 'react';
import { ProductHeader } from "./components/ProductHeader";
import { BasicDetails } from "./components/BasicDetails";
import { ProductPricing } from "./components/ProductPricing";
import { ProductInventory } from "./components/ProductInventory";
import { ProductImages } from "./components/ProductImages";
import { ProductCategories } from "./components/ProductCategories";
import { ProductShipping } from "./components/ProductShipping";
import { ProductTags } from "./components/ProductTags";
import { ProductBrands } from "./components/ProductBrands";
import { UIMetadata } from "./components/UIMetadata";
import { useProductForm } from "./hooks/useProductForm";
import { useImageUpload } from "./hooks/useImageUpload";
import AdminHeader from '../layout/AdminHeader';
import { useToast } from "@/hooks/use-toast";
import { productService } from "./services/product.service";

export const AddProduct: React.FC = () => {
  const [showDebugData, setShowDebugData] = useState<boolean>(false);
  const [debugData, setDebugData] = useState<string>('');
  const [isUploading, setIsUploading] = useState<boolean>(false);
  const { toast } = useToast();

  const {
    product,
    setProduct,
    categories,
    loading,
    error,
    handleInputChange,
    handleRadioChange,
    handleColorSelect,
    // handleFeatureChange,
    handleTagsChange,
    handleShippingInfoChange,
    handleUIMetadataChange,
    handlePublish: originalHandlePublish,
    handleSaveDraft,
    calculateSavings,
    // validateProduct
  } = useProductForm();

  const {
    images,
    mainImageIndex,
    uploadError,
    handleImageUpload,
    removeImage,
    setMainImage,
    prepareImagesForSubmit,
    getImageUrlsForProduct
  } = useImageUpload();

  // Simulate initial loading for better UX
  useEffect(() => {
    const timer = setTimeout(() => {
      // Do something when timer ends
    }, 1000);

    return () => clearTimeout(timer);
  }, []);

  // Sync images between hooks and update product.imgSlider
  useEffect(() => {
    if (images.length > 0) {
      // Get image URLs as a string array for product.imgSlider
      const imageUrls = getImageUrlsForProduct();

      // Update product with just the image URLs (as strings)
      setProduct(prev => ({
        ...prev,
        imgSlider: imageUrls
      }));

      // console.log("Updated product images:", imageUrls);
    }
  }, [images, setProduct, getImageUrlsForProduct]);

  // Validate product data before publishing
  const validateProductData = (): string | null => {
    if (!product.name || product.name.trim() === '') {
      return 'Product name is required';
    }
    if (!product.price || product.price <= 0) {
      return 'Valid product price is required';
    }
    if (!product.categoryId || product.categoryId === '0') {
      return 'Please select a category';
    }
    return null;
  };

  // Override handlePublish to upload images first
  const handlePublish = async () => {
    try {
      // 1. Validate product data
      const validationError = validateProductData();
      if (validationError) {
        const errorMsg = `Validation error: ${validationError}`;
        console.error("❌ " + errorMsg);
        toast({
          title: "Validation Error",
          description: errorMsg,
          variant: "destructive",
        });
        return;
      }

      // 2. Set loading state
      setIsUploading(true);

      // 3. First, create a basic product without images to get a valid product ID
      // console.log("🔄 Creating initial product without images to get an ID...");

      // Make a copy of the product and set imgSlider to empty for initial creation
      // const initialProduct = {
      //   ...product,
      //   imgSlider: []
      // };

      try {
        // Call the original publish method to create product
        const createdProduct = await originalHandlePublish();

        if (!createdProduct) {
          throw new Error("Failed to create product: No response from server");
        }

        // Get the product ID directly from the created product
        const productId = createdProduct.id;

        // Double check ID was returned
        if (!productId) {
          // console.error("❌ Created product missing ID:", createdProduct);
          throw new Error("Created product is missing ID");
        }

        // console.log("✅ Product created with ID:", productId);

        // Update product state explicitly to ensure ID is captured
        setProduct(prevProduct => ({
          ...prevProduct,
          id: productId
        }));

        // 4. Now check if we have local images to upload
        let hasLocalImages = false;
        if (images.length > 0) {
          hasLocalImages = images.some(img => img.url.startsWith('blob:') || img.url.startsWith('data:'));
        }

        // console.log("🖼️ Images check: ", images.length, "images,", hasLocalImages ? "has local images" : "all images are server URLs");

        // 5. If we have local images, upload them with the product ID
        if (hasLocalImages) {
          // console.log("🔄 Uploading local images to server with product ID:", productId);

          try {
            // Upload images with the valid product ID
            const serverImageUrls = await prepareImagesForSubmit(productId);
            // console.log("✅ Local images uploaded successfully, received URLs:", serverImageUrls);

            // Update the product with server URLs
            if (serverImageUrls.length > 0) {
              // console.log("🔄 Updating product with image URLs...");

              // Important: Update local product state with images
              setProduct(prev => ({
                ...prev,
                id: productId, // Ensure ID is set
                imgSlider: serverImageUrls
              }));

              // Update the product in the backend
              const result = await productService.patchProduct(productId, {
                imgSlider: serverImageUrls
              });

              if (result) {
                // console.log("✅ Product updated with images successfully!");
                toast({
                  variant: "success",
                  title: "Success",
                  description: "Product published and images uploaded successfully!"
                });
              } else {
                console.error("❌ Failed to update product with image URLs");
                toast({
                  title: "Warning",
                  description: "Product created but failed to update with images",
                  variant: "destructive",
                });
              }
            }
          } catch (uploadError: any) {
            //  console.error("❌ Error during image upload:", uploadError);
            toast({
              title: "Image Upload Failed",
              description: "Product was created but image upload failed: " + (uploadError.message || "Unknown error"),
              variant: "destructive",
            });
          }
        } else {
          // No images to upload
          toast({
            variant: "success",
            title: "Success",
            description: "Product published successfully!"
          });
        }
      } catch (createError: any) {
        // console.error("❌ Error creating product:", createError);
        toast({
          title: "Create Product Failed",
          description: createError.message || "Unknown error creating product",
          variant: "destructive",
        });
      }
    } catch (error: any) {
      // console.error("❌ Uncaught error in handlePublish:", error);
      toast({
        title: "Error",
        description: error.message || "An unexpected error occurred",
        variant: "destructive",
      });
    } finally {
      setIsUploading(false);
    }
  };

  // Handle category selection and update categorySlug
  const handleCategoryChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const categoryId = e.target.value;  // Lấy giá trị từ select (Get value from select)
    const selectedCategory = categories.find(cat => cat.id === parseInt(categoryId));

    if (selectedCategory) {
      //console.log("Selected category in component:", selectedCategory);

      setProduct(prev => ({
        ...prev,
        categoryId: categoryId, // Giữ là string vì model định nghĩa là string (Keep as string as per model definition)
        categorySlug: selectedCategory.slug || ""
      }));
    } else {
      console.warn("Selected category not found:", categoryId);

      setProduct(prev => ({
        ...prev,
        categoryId: categoryId, // Giữ là string vì model định nghĩa là string (Keep as string as per model definition)
        categorySlug: ""
      }));
    }
  };

  // Show debug data in a formatted modal/alert
  const showDebugInfo = () => {
    // Create a processed copy of the product data like it would be sent to the backend
    const debugProduct = {
      ...product,
      categoryId: product.categoryId,
      // Ensure slug is set
      slug: product.slug || product.name.toLowerCase().replace(/\s+/g, '-').replace(/[^\w\-]+/g, '')
    };

    // Format the data for display
    const debugDataStr = JSON.stringify(debugProduct, null, 2);
    setDebugData(debugDataStr);
    setShowDebugData(true);

    // Additional logging focused on slug
    // console.log("Debug - Product name:", product.name);
    // console.log("Debug - Product slug:", product.slug);
    // console.log("Debug - Category:", product.categoryId, product.categorySlug);
    // console.log("Debug - imgSlider:", debugProduct.imgSlider?.length || 0, "images");
  };

  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Add Product" />
        <main className="p-[1rem]">
          <div className="w-full md:w-[1116px] mx-auto">
            <div className="container mx-auto">
              {error && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                  {error}
                </div>
              )}

              {uploadError && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                  Image upload error: {uploadError}
                </div>
              )}

              <ProductHeader
                onPublish={handlePublish}
                onSaveDraft={handleSaveDraft}
                loading={loading || isUploading}
              />

              <div className="flex flex-wrap -mx-3 mt-[1.5rem] drop-shadow-sm filter">
                <div className="w-full md:w-7/12 px-3 mb-6">
                  <BasicDetails
                    product={product}
                    onChange={handleInputChange}
                  />

                  <ProductPricing
                    product={product}
                    onChange={handleInputChange}
                    discountAmount={calculateSavings().toString()}
                  />

                  <ProductInventory
                    product={product}
                    onChange={handleInputChange}
                    onRadioChange={handleRadioChange}
                    handlePublish={handlePublish}
                    handleSaveDraft={handleSaveDraft}
                  />

                  <ProductShipping
                    shipping={product.shippingInfo}
                    onChange={handleShippingInfoChange}
                  />

                  <UIMetadata
                    product={product}
                    onChange={handleUIMetadataChange}
                  />
                </div>

                <div className="w-full md:w-5/12 px-3">
                  <ProductImages
                    selectedImages={images.map(img => img.url)}
                    mainImage={images[mainImageIndex]?.url || null}
                    onImageUpload={(e: React.ChangeEvent<HTMLInputElement>) => {
                      const files = Array.from(e.target.files || []);
                      if (files.length > 0) {
                        handleImageUpload(files);
                      }
                    }}
                    onRemoveImage={removeImage}
                    onSetMainImage={(image: string) => {
                      const index = images.findIndex(img => img.url === image);
                      if (index !== -1) {
                        setMainImage(index);
                      }
                    }}
                    isUploading={isUploading}
                  />

                  <ProductBrands
                    brand={product.brand || ""}
                    onChange={handleInputChange}
                  />

                  <ProductCategories
                    product={product}
                    categories={categories}
                    onChange={handleCategoryChange}
                    onColorSelect={handleColorSelect}
                  />

                  <ProductTags
                    tags={product.tags || []}
                    onChange={handleTagsChange}
                  />
                </div>
              </div>

              <div className="flex justify-between mt-6 space-x-3 px-3 mb-6">
                <button
                  onClick={showDebugInfo}
                  className="px-5 py-2 border border-blue-300 rounded-lg flex items-center hover:bg-blue-50 text-blue-700"
                >
                  <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Debug Data
                </button>

                <div className="flex space-x-3">
                  <button
                    onClick={handleSaveDraft}
                    className="px-5 py-2 border border-gray-200 rounded-lg flex items-center hover:bg-gray-50"
                  >
                    <svg
                      className="w-5 h-5 mr-2"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
                      />
                    </svg>
                    Save to draft
                  </button>
                  <button
                    onClick={handlePublish}
                    className="px-5 py-2 bg-ocean-green hover:bg-green-600 text-white rounded-lg"
                  >
                    Publish Product
                  </button>
                </div>
              </div>

              {/* Debug Data Modal */}
              {showDebugData && (
                <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
                  <div className="bg-white rounded-lg shadow-xl max-w-4xl max-h-[80vh] w-full overflow-hidden">
                    <div className="p-4 border-b flex justify-between items-center">
                      <h3 className="text-lg font-bold">Debug Product Data</h3>
                      <button
                        onClick={() => setShowDebugData(false)}
                        className="text-gray-500 hover:text-gray-700"
                      >
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                    <div className="p-4 overflow-auto max-h-[calc(80vh-120px)]">
                      <div className="border rounded p-4 bg-gray-50 mb-4">
                        <h4 className="font-bold mb-2">Required Field Status:</h4>
                        <ul className="grid grid-cols-2 gap-2">
                          <li className={`${product.name ? 'text-green-600' : 'text-red-600'}`}>
                            Name: {product.name ? '✓' : '✗'}
                          </li>
                          <li className={`${product.slug ? 'text-green-600' : 'text-red-600'}`}>
                            Slug: {product.slug ? '✓' : '✗'} {product.slug}
                          </li>
                          <li className={`${product.price > 0 ? 'text-green-600' : 'text-red-600'}`}>
                            Price: {product.price > 0 ? '✓' : '✗'} ({product.price})
                          </li>
                          <li className={`${parseInt(product.categoryId) > 0 ? 'text-green-600' : 'text-red-600'}`}>
                            Category ID: {parseInt(product.categoryId) > 0 ? '✓' : '✗'} ({product.categoryId})
                          </li>
                          <li className={`${product.categorySlug ? 'text-green-600' : 'text-red-600'}`}>
                            Category Slug: {product.categorySlug ? '✓' : '✗'} ({product.categorySlug})
                          </li>
                          <li className={`${product.brand ? 'text-green-600' : 'text-red-600'}`}>
                            Brand: {product.brand ? '✓' : '✗'}
                          </li>
                          <li className={`${product.imgSlider?.length > 0 ? 'text-green-600' : 'text-red-600'}`}>
                            Images: {product.imgSlider?.length > 0 ? '✓' : '✗'} ({product.imgSlider?.length})
                          </li>
                          <li className={`text-green-600`}>
                            UI Metadata: ✓ (Optional)
                          </li>
                        </ul>
                      </div>
                      {/* Add explicit slug display */}
                      <div className="border rounded p-4 bg-blue-50 mb-4">
                        <h4 className="font-bold mb-2">Generated Slug Information:</h4>
                        <div>
                          <p><span className="font-semibold">Product Name:</span> {product.name}</p>
                          <p><span className="font-semibold">Generated Slug:</span> {product.slug}</p>
                          <p className="text-xs text-gray-600 mt-2">
                            Note: The slug is automatically generated from the product name, but can be overridden if needed.
                          </p>
                        </div>
                      </div>
                      <div className="border rounded p-4 bg-green-50 mb-4">
                        <h4 className="font-bold mb-2">UI Metadata:</h4>
                        <div>
                          <pre className="text-xs bg-white p-3 rounded overflow-x-auto">
                            {JSON.stringify(product.uiMetadata, null, 2) || "{}"}
                          </pre>
                          <p className="text-xs text-gray-600 mt-2">
                            Custom UI settings for frontend display and behavior. This is used for special rendering features.
                          </p>
                        </div>
                      </div>
                      <pre className="text-xs bg-gray-800 text-white p-4 rounded overflow-x-auto">
                        {debugData}
                      </pre>
                    </div>
                    <div className="p-4 border-t flex justify-end">
                      <button
                        onClick={() => {
                          navigator.clipboard.writeText(debugData);
                          alert('Debug data copied to clipboard!');
                        }}
                        className="mr-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                      >
                        Copy to Clipboard
                      </button>
                      <button
                        onClick={() => setShowDebugData(false)}
                        className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                      >
                        Close
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};