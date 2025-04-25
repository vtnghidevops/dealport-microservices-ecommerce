// // components/admin/product/add/hooks/useImageUpload.ts
// import { useState } from 'react';
// import { ProductImage } from '@/types/product.model';
// import { productService } from '../services/product.service';

// export function useImageUpload() {
//   const [images, setImages] = useState<ProductImage[]>([]);
//   const [mainImageIndex, setMainImageIndex] = useState<number>(0);
//   const [isUploading, setIsUploading] = useState<boolean>(false);
//   const [uploadError, setUploadError] = useState<string | null>(null);

//   const handleImageUpload = async (files: File[]) => {
//     if (files.length === 0) return;

//     setIsUploading(true);
//     setUploadError(null);

//     try {
//       // Upload images to server and get URLs
//       const uploadedUrls = await productService.uploadImages(files);

//       console.log("Uploaded image URLs:", uploadedUrls);

//       // Create ProductImage objects from URLs
//       const newImages = uploadedUrls.map((url, index) => ({
//         url,
//         is_primary: images.length === 0 && index === 0, // First image is primary only if no existing images
//         display_order: images.length + index
//       }));

//       const updatedImages = [...images, ...newImages];

//       // If this is the first upload, set the first image as main
//       if (images.length === 0 && newImages.length > 0) {
//         setMainImageIndex(0);
//         // Ensure at least one image is marked as primary
//         updatedImages[0].is_primary = true;
//       }

//       setImages(updatedImages);

//       console.log("Updated images array:", updatedImages);
//     } catch (error) {
//       setUploadError(error instanceof Error ? error.message : 'Failed to upload images');
//       console.error('Error uploading images:', error);
//     } finally {
//       setIsUploading(false);
//     }
//   };

//   const removeImage = (index: number) => {
//     if (index < 0 || index >= images.length) return;

//     const newImages = [...images];
//     newImages.splice(index, 1);

//     // Update main image index if needed
//     if (index === mainImageIndex) {
//       // If we removed the main image, set the first image as main
//       setMainImageIndex(newImages.length > 0 ? 0 : -1);
//     } else if (index < mainImageIndex) {
//       // If we removed an image before the main image, adjust the index
//       setMainImageIndex(mainImageIndex - 1);
//     }

//     // Ensure there is still one primary image
//     if (newImages.length > 0) {
//       // Reset all images to non-primary
//       newImages.forEach(img => img.is_primary = false);

//       // Set the new main image as primary
//       const newMainIndex = mainImageIndex >= newImages.length ? 0 : (mainImageIndex < 0 ? 0 : mainImageIndex);
//       newImages[newMainIndex].is_primary = true;
//       setMainImageIndex(newMainIndex);
//     }

//     setImages(newImages);
//     console.log("Images after removal:", newImages);
//   };

//   const setMainImage = (index: number) => {
//     if (index >= 0 && index < images.length) {
//       const newImages = [...images];

//       // Remove primary status from all images
//       newImages.forEach(img => img.is_primary = false);

//       // Set new image as primary
//       newImages[index].is_primary = true;
//       setMainImageIndex(index);
//       setImages(newImages);

//       console.log("Main image set to index:", index, newImages[index]);
//     }
//   };

//   return {
//     images,
//     mainImageIndex,
//     isUploading,
//     uploadError,
//     handleImageUpload,
//     removeImage,
//     setMainImage
//   };
// }
// components/admin/product/add/hooks/useImageUpload.ts
import { useState } from 'react';
import { ProductImage } from '@/types/product.model';
import { productService } from '../services/product.service';

export function useImageUpload() {
  const [images, setImages] = useState<ProductImage[]>([]);
  const [mainImageIndex, setMainImageIndex] = useState<number>(0);
  const [isUploading, setIsUploading] = useState<boolean>(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  // Track which images are local vs server
  const [localImages, setLocalImages] = useState<{ [key: string]: File }>({});

  const handleImageUpload = async (files: File[]) => {
    if (files.length === 0) return;

    setIsUploading(true);
    setUploadError(null);

    try {
      // Create immediate local preview URLs
      const newImages: ProductImage[] = [];
      const newLocalImages = { ...localImages };

      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        // Create blob URL for immediate preview
        const localUrl = URL.createObjectURL(file);

        // Save the local URL and corresponding file
        newLocalImages[localUrl] = file;

        // Create a ProductImage object with the local URL
        newImages.push({
          url: localUrl,
          isPrimary: images.length === 0 && i === 0,
          displayOrder: images.length + i
        });
      }

      // Update local images map
      setLocalImages(newLocalImages);

      // Update UI with local previews
      const updatedImages = [...images, ...newImages];

      // Set first image as main if this is the first upload
      if (images.length === 0 && newImages.length > 0) {
        setMainImageIndex(0);
        updatedImages[0].isPrimary = true;
      }

      setImages(updatedImages);

    } catch (error) {
      setUploadError(error instanceof Error ? error.message : 'Failed to process images');
      console.error('Error processing images:', error);
    } finally {
      setIsUploading(false);
    }
  };

  const removeImage = (index: number) => {
    if (index < 0 || index >= images.length) return;

    const imageToRemove = images[index];
    const newImages = [...images];

    // Remove from local images if it's a blob URL
    if (imageToRemove.url.startsWith('blob:')) {
      URL.revokeObjectURL(imageToRemove.url);
      const newLocalImages = { ...localImages };
      delete newLocalImages[imageToRemove.url];
      setLocalImages(newLocalImages);
    }

    newImages.splice(index, 1);

    // Update main image index if needed
    if (index === mainImageIndex) {
      // If we removed the main image, set the first image as main
      setMainImageIndex(newImages.length > 0 ? 0 : -1);
    } else if (index < mainImageIndex) {
      // If we removed an image before the main image, adjust the index
      setMainImageIndex(mainImageIndex - 1);
    }

    // Ensure there is still one primary image
    if (newImages.length > 0) {
      // Reset all images to non-primary
      newImages.forEach(img => img.isPrimary = false);

      // Set the new main image as primary
      const newMainIndex = mainImageIndex >= newImages.length ? 0 : (mainImageIndex < 0 ? 0 : mainImageIndex);
      newImages[newMainIndex].isPrimary = true;
      setMainImageIndex(newMainIndex);
    }

    setImages(newImages);
  };

  const setMainImage = (index: number) => {
    if (index >= 0 && index < images.length) {
      const newImages = [...images];

      // Remove primary status from all images
      newImages.forEach(img => img.isPrimary = false);

      // Set new image as primary
      newImages[index].isPrimary = true;
      setMainImageIndex(index);
      setImages(newImages);
    }
  };

  // This function uploads all local images to the server before submitting the product
  const prepareImagesForSubmit = async (productId: number): Promise<string[]> => {
    if (!productId || isNaN(Number(productId))) {
      console.error('❌ ERROR - prepareImagesForSubmit: Missing or invalid product ID', productId);
      setUploadError('Cannot upload images without a valid product ID');
      throw new Error('Product ID is required for image upload');
    }

   // console.log('🧪 DEBUG - prepareImagesForSubmit: Using product ID:', productId);

    // Step 1: Find all local images that need to be uploaded
    const localImagesUrls = Object.keys(localImages);
    const dataUrls = images.filter(img => img.url.startsWith('data:'));
    const blobUrls = images.filter(img => img.url.startsWith('blob:'));

    // Make sure we capture ALL blob URLs, not just those in localImages
    blobUrls.forEach(img => {
      if (!localImagesUrls.includes(img.url)) {
        console.warn('⚠️ WARNING - Found blob URL not in localImages:', img.url);
        // Add to localImages if we can find the corresponding file
        // Note: This is a fallback and may not work in all cases
      }
    });

    const needsUpload = localImagesUrls.length > 0 || dataUrls.length > 0 || blobUrls.length > 0;

    // console.log('🧪 DEBUG - prepareImagesForSubmit: Found', localImagesUrls.length, 'blob URLs in localImages,',
      // blobUrls.length, 'total blob URLs, and', dataUrls.length, 'data URLs that need uploading');

    // If no uploads needed, just return current URLs that are not blob: or data: URLs
    if (!needsUpload) {
      const validServerUrls = images
        .map(img => img.url)
        .filter(url => !url.startsWith('blob:') && !url.startsWith('data:'));

      // console.log('🧪 DEBUG - No local images to upload, returning valid server URLs:', validServerUrls);
      return validServerUrls;
    }

    setIsUploading(true);

    try {
      // Step 2: Collect all files that need to be uploaded
      const filesToUpload: File[] = [];

      // Add files from localImages (blob URLs)
      localImagesUrls.forEach(url => {
        if (localImages[url]) {
          filesToUpload.push(localImages[url]);
        }
      });

      // Convert data URLs to files if any
      if (dataUrls.length > 0) {
        // console.log('🧪 DEBUG - Converting', dataUrls.length, 'data URLs to files');

        // Convert data URLs to files
        const dataUrlFiles = await Promise.all(
          dataUrls.map(async (img, index) => {
            try {
              const response = await fetch(img.url);
              const blob = await response.blob();
              return new File([blob], `image-${index}.jpg`, { type: blob.type || 'image/jpeg' });
            } catch (error) {
              console.error('Failed to convert data URL to file:', error);
              return null;
            }
          })
        );

        // Add valid data URL files to upload list
        dataUrlFiles.filter(Boolean).forEach((file) => {
          if (file) filesToUpload.push(file);
        });
      }

      // Step 3: Upload all files to server
      if (filesToUpload.length === 0) {
        console.warn('⚠️ WARNING - No files to upload after preparation');
        // Return only valid server URLs (not blob: or data:)
        return images
          .map(img => img.url)
          .filter(url => !url.startsWith('blob:') && !url.startsWith('data:'));
      }

     // console.log('🧪 DEBUG - Uploading', filesToUpload.length, 'files to server with product ID:', productId);

      // Upload files to server with product ID
      const serverUrls = await productService.uploadImages(filesToUpload, productId);
      // console.log('🧪 DEBUG - Server responded with', serverUrls.length, 'URLs:', serverUrls);

      // Step 4: Replace local URLs with server URLs
      if (serverUrls.length > 0) {
        // Create a map of local->server URLs by index
        const serverUrlMap = new Map<string, string>();

        // First map all blob URLs
        localImagesUrls.forEach((localUrl, index) => {
          if (index < serverUrls.length) {
            serverUrlMap.set(localUrl, serverUrls[index]);
          }
        });

        // Then map any data URLs - offset by number of blob URLs
        const blobCount = localImagesUrls.length;
        dataUrls.forEach((img, index) => {
          const serverIndex = blobCount + index;
          if (serverIndex < serverUrls.length) {
            serverUrlMap.set(img.url, serverUrls[serverIndex]);
          }
        });

        // Update image objects with server URLs
        const updatedImages = images.map(img => {
          const serverUrl = serverUrlMap.get(img.url);
          if (serverUrl) {
            return {
              ...img,
              url: serverUrl
            };
          }
          // Don't add hostname to relative paths - server expects relative URLs
          return img;
        });

        // Step 5: Clean up and update state
        // Revoke blob URLs
        localImagesUrls.forEach(url => {
          URL.revokeObjectURL(url);
        });

        setImages(updatedImages);
        setLocalImages({});

        // Return just the valid URLs for the product (no blob: or data: URLs)
        const finalUrls = updatedImages
          .map(img => img.url)
          .filter(url => !url.startsWith('blob:') && !url.startsWith('data:'));

       // console.log('✅ Final image URLs after upload:', finalUrls);
        return finalUrls;
      } else {
        // If no URLs returned but upload didn't throw, return existing valid images
        console.warn('⚠️ Warning: Upload succeeded but no URLs returned');
        return images
          .map(img => img.url)
          .filter(url => !url.startsWith('blob:') && !url.startsWith('data:'));
      }
    } catch (error) {
      console.error('❌ ERROR - Image upload failed:', error);
      setUploadError('Failed to upload images to server: ' +
        (error instanceof Error ? error.message : 'Unknown error'));

      // Return current valid URLs even on failure (filter out blob: and data:)
      return images
        .map(img => img.url)
        .filter(url => !url.startsWith('blob:') && !url.startsWith('data:'));
    } finally {
      setIsUploading(false);
    }
  };

  // Get image URLs in the right format for the product object
  const getImageUrlsForProduct = (): string[] => {
    return images.map(img => {
      // Return URLs as is - don't add hostname
      return img.url;
    });
  };

  return {
    images,
    mainImageIndex,
    isUploading,
    uploadError,
    handleImageUpload,
    removeImage,
    setMainImage,
    prepareImagesForSubmit,
    getImageUrlsForProduct
  };
}