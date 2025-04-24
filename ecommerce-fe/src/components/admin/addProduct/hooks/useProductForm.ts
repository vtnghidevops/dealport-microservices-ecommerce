// components/admin/product/add/hooks/useProductForm.ts
import { useState, useEffect, ChangeEvent } from "react";
import { Product, ShippingInfo } from "@/types/product.model";
import { productService } from "../services/product.service";

// Interface to handle image objects from backend
interface ProductImageResponse {
  url: string;
  is_primary?: boolean;
  display_order?: number;
  product_id?: number;
}

// Initial empty product to use as a starting state
const initialProduct: Product = { // TypeScript hack: undefined for new products, will be assigned by backend
  id: undefined as unknown as number,
  name: "",
  description: "",
  slug: "",
  categorySlug: "",
  brand: "",
  type: "normal",
  price: 0,
  originalPrice: 0,
  discount: 0,
  categoryId: "0",
  stockQuantity: 10,
  image_url: "",
  imgSlider: [], // This will always be a string[] in the frontend
  tags: [],
  features: [],
  shippingInfo: {
    courier: "",
    local: "",
    ups: "",
    global: ""
  },
  uiMetadata: {
    color: "#FFFFFF",
    backgroundColor: "#F8F8F8",
    showBadge: false,
    badgeText: "",
    badgeColor: "#FF0000",
    rating: 0,
    displayMode: "default",
    showPromotionalLabel: false,
    promotionalText: ""
  },
  reviewsAvg: {
    rating: 0,
    count: 0
  },
  orders: 0
};

export function useProductForm() {
  const [product, setProduct] = useState<Product>(initialProduct);
  const [categories, setCategories] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch categories on component mount
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const categoryData = await productService.getCategories();
        setCategories(categoryData || []);
        console.log("Categories loaded:", categoryData);
      } catch (error) {
        console.error("Error fetching categories:", error);
        setError("Failed to load categories");
      }
    };

    fetchCategories();
  }, []);

  // Generate slug from product name
  const generateSlug = (name: string): string => {
    return name
      .toString()
      .toLowerCase()
      .trim()
      .replace(/\s+/g, '-')           // Replace spaces with -
      .replace(/[^\w\-]+/g, '')       // Remove all non-word chars
      .replace(/\-\-+/g, '-')         // Replace multiple - with single -
      .replace(/^-+/, '')             // Trim - from start of text
      .replace(/-+$/, '');            // Trim - from end of text
  };

  // Auto-generate slug when name changes
  useEffect(() => {
    if (product.name && !product.slug) {
      setProduct((prev: Product) => ({
        ...prev,
        slug: generateSlug(prev.name)
      }));
    }
  }, [product.name]);

  const handleInputChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;

    // Handle checkbox fields
    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setProduct((prev: Product) => ({
        ...prev,
        [name]: checked
      }));
      return;
    }

    // Handle numeric fields
    if (
      ['price', 'originalPrice', 'discount', 'stockQuantity'].includes(name) &&
      value !== ''
    ) {
      let parsedValue: number;

      // Special handling for price fields to allow decimal values
      if (['price', 'originalPrice', 'discount'].includes(name)) {
        parsedValue = parseFloat(value);
      } else {
        parsedValue = parseInt(value);
      }

      // Only update if the parsed value is a valid number
      if (!isNaN(parsedValue)) {
        setProduct((prev: Product) => ({
          ...prev,
          [name]: parsedValue
        }));
      }
      return;
    }

    // Special handling for categoryId
    if (name === 'categoryId') {
      const categoryIdNum = parseInt(value, 10); // Parse to number for comparison
      // Find the matching category to get the slug
      const selectedCategory = categories.find(cat => cat.id === categoryIdNum);

      setProduct((prev: Product) => ({
        ...prev,
        categoryId: value, // Keep as string to match Product interface
        categorySlug: selectedCategory?.slug || ""
      }));

      return;
    }

    // Handle slug field with special formatting
    if (name === 'slug') {
      setProduct((prev: Product) => ({
        ...prev,
        slug: generateSlug(value)
      }));
      return;
    }

    // Default handling for all other fields
    setProduct((prev: Product) => ({
      ...prev,
      [name]: value
    }));
  };

  const handleRadioChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProduct((prev: Product) => ({
      ...prev,
      [name]: value
    }));
  };

  const handleFeatureChange = (features: string[]) => {
    setProduct((prev: Product) => ({ ...prev, features }));
  };

  const handleTagsChange = (tags: string[]) => {
    setProduct((prev: Product) => ({ ...prev, tags }));
  };

  const handleShippingInfoChange = (field: string, value: string) => {
    setProduct((prev: Product) => {
      const shippingInfo: ShippingInfo = {
        courier: prev.shippingInfo?.courier || "",
        local: prev.shippingInfo?.local || "",
        ups: prev.shippingInfo?.ups || "",
        global: prev.shippingInfo?.global || "",
        [field]: value
      };

      return {
        ...prev,
        shippingInfo
      };
    });
  };

  // Updated to ensure imgSlider is always a string array in the frontend
  const handleImagesChange = (images: any[]) => {
    // Images can be either ProductImage objects or simple URLs

    // Create a proper ProductImage array for the backend
    const productImages = images.map((img, index) => {
      // If img is already a ProductImage object, use it
      if (typeof img === 'object' && img !== null && 'url' in img) {
        return {
          ...img,
          is_primary: index === 0, // First image is primary
          display_order: index
        };
      }
      // Otherwise create a new ProductImage object from the URL
      return {
        url: img,
        is_primary: index === 0,
        display_order: index
      };
    });

    // Make sure we always store strings in imgSlider, not objects
    const imgSlider = images.map(img => {
      if (typeof img === 'object' && img !== null && 'url' in img) {
        return img.url;
      }
      return img;
    });

    // Also update image_url to use the first image
    const image_url = imgSlider.length > 0 ? imgSlider[0] : '';

    setProduct((prev: Product) => ({
      ...prev,
      imgSlider,
      image_url,
      images: productImages // Add the images array with ProductImage objects
    }));
  };

  const handleUIMetadataChange = (metadata: any) => {
    setProduct((prev: Product) => ({
      ...prev,
      uiMetadata: metadata
    }));
  };

  // Tính phần trăm giảm giá dựa trên giá gốc và số tiền giảm
  const calculateDiscountPercentage = () => {
    if (product.originalPrice > 0) {
      if (product.discount > 0) {
        // Nếu có discount (số tiền giảm trực tiếp), tính % = (số tiền giảm / giá gốc) * 100
        const percentage = (product.discount / product.originalPrice) * 100;
        return Math.round(percentage);
      } else if (product.price > 0) {
        // Nếu có giá sau giảm, tính % = ((giá gốc - giá sau giảm) / giá gốc) * 100
        const discountAmount = product.originalPrice - product.price;
        const percentage = (discountAmount / product.originalPrice) * 100;
        return Math.round(percentage);
      }
    }
    // Trả về 0 nếu không có thông tin đủ để tính
    return 0;
  };

  // Tính số tiền tiết kiệm (giảm) dựa trên giá gốc và discount
  const calculateSavings = () => {
    if (product.originalPrice > 0 && product.price > 0) {
      // Nếu có cả giá gốc và giá đã giảm, thì tiền giảm = giá gốc - giá đã giảm
      return product.originalPrice - product.price;
    } else if (product.originalPrice > 0 && product.discount > 0) {
      // Discount đã là số tiền giảm trực tiếp
      return product.discount;
    }
    return 0;
  };

  const handleColorSelect = (color: string) => {
    setProduct((prev: Product) => ({
      ...prev,
      uiMetadata: {
        ...prev.uiMetadata,
        color
      }
    }));
  };

  // Validate product data before submitting to backend
  const validateProduct = (data: Partial<Product> = product): string | null => {
    // Validate required fields
    if (!data.name || data.name.trim() === '') {
      return 'Product name is required';
    }

    if (!data.price || data.price <= 0) {
      return 'Valid product price is required';
    }

    if (!data.categoryId || data.categoryId === '0') {
      return 'Please select a category';
    }

    // More validations can be added here
    return null;
  };

  const handlePublish = async () => {
    setLoading(true);
    setError(null);

    try {
      // Create a processed copy of the product to submit
      const updatedProduct = { ...product };

      // Generate slug from name if empty
      if (!updatedProduct.slug) {
        updatedProduct.slug = generateSlug(updatedProduct.name);
        console.log("Generated slug for publishing:", updatedProduct.slug);
      }

      // Double-check that categorySlug is set
      if (!updatedProduct.categorySlug && updatedProduct.categoryId) {
        const selectedCategory = categories.find(cat => cat.id === parseInt(updatedProduct.categoryId));
        if (selectedCategory) {
          updatedProduct.categorySlug = selectedCategory.slug;
          console.log("Fixed missing categorySlug:", updatedProduct.categorySlug);
        }
      }

      // ADDITIONAL CRITICAL FIX: Ensure imgSlider contains only string URLs
      if (updatedProduct.imgSlider && Array.isArray(updatedProduct.imgSlider)) {
        // Check if any object exists in imgSlider and convert to strings
        const normalizedImgUrls = updatedProduct.imgSlider.map(img => {
          if (typeof img === 'object' && img !== null && 'url' in img) {
            console.log('⚠️ Found complex object in imgSlider, converting to URL string:', img);
            return (img as any).url;
          }
          return img;
        });

        // Update imgSlider with normalized data
        if (JSON.stringify(normalizedImgUrls) !== JSON.stringify(updatedProduct.imgSlider)) {
          console.log('📝 Normalized imgSlider before sending to API:', normalizedImgUrls);
          updatedProduct.imgSlider = normalizedImgUrls;

          // Also update image_url
          updatedProduct.image_url = normalizedImgUrls.length > 0 ? normalizedImgUrls[0] : "";
        }
      }

      // We don't need to update state here as we're just sending data to the API
      // Don't call setProduct(updatedProduct) here to prevent infinite update loop

      // Cập nhật cách tính discount - Discount là số tiền giảm trực tiếp, không phải phần trăm
      if (updatedProduct.originalPrice > 0 && updatedProduct.discount > 0) {
        // Tính giá sau khi giảm = giá gốc - số tiền giảm
        const discountedPrice = updatedProduct.originalPrice - updatedProduct.discount;

        // Đảm bảo giá không âm
        updatedProduct.price = Math.max(0, parseFloat(discountedPrice.toFixed(2)));

        console.log(`Giá gốc: ${updatedProduct.originalPrice}, Số tiền giảm: ${updatedProduct.discount}, Giá sau giảm: ${updatedProduct.price}`);
      }

      console.log("Final product data before validation:", updatedProduct);

      // Display key data in console for debugging
      console.group("📋 KEY PRODUCT DATA");
      console.log(`Name: ${updatedProduct.name}`);
      console.log(`Slug: ${updatedProduct.slug}`);
      console.log(`Category ID: ${updatedProduct.categoryId} (${typeof updatedProduct.categoryId})`);
      console.log(`Category Slug: ${updatedProduct.categorySlug}`);
      console.log(`Price: ${updatedProduct.price}`);
      console.log(`Original Price: ${updatedProduct.originalPrice}`);
      console.log(`Brand: ${updatedProduct.brand}`);
      console.log(`Images: ${updatedProduct.imgSlider.length} images`);
      if (updatedProduct.imgSlider.length > 0) {
        console.log(`Primary Image: ${updatedProduct.imgSlider[0] || 'None'}`);
      }
      console.log(`UI Metadata: ${JSON.stringify(updatedProduct.uiMetadata)}`);
      console.groupEnd();

      // Validate product data before submitting
      const validationError = validateProduct(updatedProduct);
      if (validationError) {
        setError(validationError);
        setLoading(false);
        return;
      }

      // Remove id property before sending to API (backend will assign a new ID)
      const { id, ...productDataForBackend } = updatedProduct;

      console.log("Sending product data to API:", productDataForBackend);
      const result = await productService.createProduct(productDataForBackend);

      console.log("API response:", result);

      if (result) {
        // IMPORTANT! Update the product state with the new ID from the API response
        console.log("Updating product state with new ID:", result.id);
        if (result.id) {
          setProduct(prev => ({
            ...prev,
            id: result.id
          }));
        } else {
          console.error("❌ API response missing product ID!", result);
        }

        // Don't show alert in component, we're using Toast in AddProduct component
        // alert("Product published successfully!");

        return result; // Return the result for the parent component to use
      } else {
        setError("Failed to publish product");
        return null;
      }
    } catch (err) {
      console.error("Error creating product:", err);
      setError(err instanceof Error ? err.message : "An unknown error occurred");
      throw err; // Re-throw to let the parent component handle
    } finally {
      setLoading(false);
    }
  };

  const handleSaveDraft = () => {
    // In a real implementation, this might save to localStorage or a draft API endpoint
    localStorage.setItem('productDraft', JSON.stringify(product));
    alert("Draft saved successfully!");
  };

  return {
    product,
    setProduct,
    categories,
    loading,
    error,
    handleInputChange,
    handleRadioChange,
    handleFeatureChange,
    handleTagsChange,
    handleShippingInfoChange,
    handleImagesChange,
    handleColorSelect,
    handleUIMetadataChange,
    handlePublish,
    handleSaveDraft,
    calculateSavings,
    calculateDiscountPercentage,
    validateProduct
  };
}