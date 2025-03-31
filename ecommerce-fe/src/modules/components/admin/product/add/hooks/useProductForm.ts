// components/admin/product/add/hooks/useProductForm.ts
import { useState, ChangeEvent } from "react";
import { Product } from "../../models/product.model";
import { productService } from "../../services/product.service";

export function useProductForm() {
  const [product, setProduct] = useState<Product>({
    id: "",
    name: "",
    description: "",
    price: 0,
    discountedPrice: 0,
    saleAmount: 0,
    taxIncluded: true,
    expirationStart: "",
    expirationEnd: "",
    stockQuantity: "Unlimited",
    stockStatus: "In Stock",
    highlighted: false,
    images: [],
    categories: [],
    tags: [],
    color: "",
  });

  const handleInputChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) =>  {
    const { name, value, type } = e.target as HTMLInputElement;

    if (type === "checkbox") {
      const checked = (e.target as HTMLInputElement).checked;
      setProduct((prev) => ({ ...prev, [name]: checked }));
    } else {
      setProduct((prev) => ({ ...prev, [name]: value }));
    }
  };

  const handleRadioChange = (e: ChangeEvent<HTMLInputElement>) => {
    setProduct((prev) => ({
      ...prev,
      [e.target.name]: e.target.value === "Yes",
    }));
  };

  const handleColorSelect = (color: string) => {
    setProduct((prev) => ({ ...prev, color }));
  };

  const calculateSavings = () => {
    const originalPrice = parseFloat(String(product.price)) || 0;
    const discountPrice = parseFloat(String(product.discountedPrice)) || 0;

    if (discountPrice && discountPrice < originalPrice) {
      const savings = originalPrice - discountPrice;
      return savings.toFixed(2);
    }

    return "0.00";
  };

  const handlePublish = () => {
    console.log("Publishing product:", product);
    // Integrate with productService here
    productService.createProduct(product);
    alert("Product published successfully!");
  };

  const handleSaveDraft = () => {
    console.log("Saving draft:", product);
    // Save draft logic here
    alert("Draft saved successfully!");
  };

  return {
    product,
    setProduct,
    handleInputChange,
    handleRadioChange,
    handleColorSelect,
    handlePublish,
    handleSaveDraft,
    calculateSavings,
  };
}