// components/admin/product/add/hooks/useImageUpload.ts
import { useState, ChangeEvent } from "react";

export function useImageUpload() {
  const [selectedImages, setSelectedImages] = useState<string[]>([]);
  const [mainImage, setMainImage] = useState<string | null>(null);

  const handleImageUpload = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files) return;

    const newImages = Array.from(files).map((file) =>
      URL.createObjectURL(file)
    );
    
    setSelectedImages((prev) => [...prev, ...newImages]);

    if (!mainImage && newImages.length > 0) {
      setMainImage(newImages[0]);
    }
  };

  const removeImage = (index: number) => {
    const newImages = [...selectedImages];
    const removedImage = newImages.splice(index, 1)[0];
    setSelectedImages(newImages);
    
    if (mainImage === removedImage) {
      setMainImage(newImages.length > 0 ? newImages[0] : null);
    }
  };

  return {
    selectedImages,
    mainImage,
    handleImageUpload,
    removeImage,
    setMainImage,
  };
}