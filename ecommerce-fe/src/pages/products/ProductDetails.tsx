import React, { useRef, useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Product } from '@/types/product.model';
import ProductService from '@/services/product/product.service';
import NotFound from '../system/NotFound';
import ProductInformation from '@/components/product/ProductInfomation';
import RelatedProduct from '@/components/product/RelatedProduct';
import { CiHeart } from "react-icons/ci";
import { IoShareSocialOutline } from "react-icons/io5";
import { IoIosArrowForward, IoIosArrowBack } from "react-icons/io";
import { TbScale } from "react-icons/tb";
import ProductComments from '@/components/product/ProductComments';
import { useCart } from '@/hooks/useCart';
import { ProductDetailSkeleton } from '@/components/ui/skeletons';

const ProductDetail: React.FC = () => {
  const { categorySlug, productSlug } = useParams<{
    categorySlug: string;
    productSlug: string;
  }>();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [quantity, setQuantity] = useState<number>(1);
  const [selectedImage, setSelectedImage] = useState<string>('');
  const [currentImageIndex, setCurrentImageIndex] = useState<number>(0);
  const [thumbnailStartIndex, setThumbnailStartIndex] = useState<number>(0);
  const maxVisibleThumbnails = 5; // Maximum number of visible thumbnails
  const commentsRef = useRef<HTMLDivElement>(null);
  const { addToCart } = useCart();

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        if (categorySlug && productSlug) {
          const data = await ProductService.getProductBySlug(productSlug);
          setProduct(data || null);
          if (data) {
            setSelectedImage(data.imageUrl);
          }
        }
      } catch (error) {
        console.error('Error fetching product details:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [categorySlug, productSlug]);
  console.log("product to receive: ", product)

  if (loading) {
    return (
      <div className="container mx-auto px-4 md:px-[5rem] py-[1rem]">
        <ProductDetailSkeleton />
      </div>
    );
  }

  if (!product) {
    return <NotFound />;
  }

  // Create an array of all images for the slider
  const allImages = [product.imageUrl, ...(product.imgSlider || [])];

  // Calculate the end index for visible thumbnails
  const thumbnailEndIndex = Math.min(thumbnailStartIndex + maxVisibleThumbnails, allImages.length);

  // Get only the thumbnails that should be visible in the current view
  const visibleThumbnails = allImages.slice(thumbnailStartIndex, thumbnailEndIndex);

  // Check if navigation buttons should be shown
  const showPrevThumbnailsButton = thumbnailStartIndex > 0;
  const showNextThumbnailsButton = thumbnailEndIndex < allImages.length;

  const handleQuantityChange = (newQty: number) => {
    if (newQty >= 1 && newQty <= product.stockQuantity) {
      setQuantity(newQty);
    }
  };

  const decreaseQuantity = () => {
    if (quantity > 1) {
      setQuantity(quantity - 1);
    }
  };

  const increaseQuantity = () => {
    if (quantity < product.stockQuantity) {
      setQuantity(quantity + 1);
    }
  };

  // Navigation functions for main image slider
  const goToPreviousImage = () => {
    const newIndex = currentImageIndex > 0 ? currentImageIndex - 1 : allImages.length - 1;
    setCurrentImageIndex(newIndex);
    setSelectedImage(allImages[newIndex]);

    // Make sure the thumbnail for the current image is visible
    if (newIndex < thumbnailStartIndex) {
      setThumbnailStartIndex(Math.max(0, newIndex));
    } else if (newIndex >= thumbnailStartIndex + maxVisibleThumbnails) {
      setThumbnailStartIndex(Math.max(0, newIndex - maxVisibleThumbnails + 1));
    }
  };

  const goToNextImage = () => {
    const newIndex = currentImageIndex < allImages.length - 1 ? currentImageIndex + 1 : 0;
    setCurrentImageIndex(newIndex);
    setSelectedImage(allImages[newIndex]);

    // Make sure the thumbnail for the current image is visible
    if (newIndex < thumbnailStartIndex) {
      setThumbnailStartIndex(Math.max(0, newIndex));
    } else if (newIndex >= thumbnailStartIndex + maxVisibleThumbnails) {
      setThumbnailStartIndex(Math.max(0, newIndex - maxVisibleThumbnails + 1));
    }
  };

  // Select a specific image
  const selectImage = (image: string, index: number) => {
    setSelectedImage(image);
    setCurrentImageIndex(index);
  };

  // Thumbnail navigation functions
  const goToPreviousThumbnails = () => {
    setThumbnailStartIndex(Math.max(0, thumbnailStartIndex - 1));
  };

  const goToNextThumbnails = () => {
    setThumbnailStartIndex(Math.min(allImages.length - maxVisibleThumbnails, thumbnailStartIndex + 1));
  };

  // Format discount percentage
  const discountPercentage = product.discount ? `${product.discount}% OFF` : null;


  // Hàm scroll đến phần comments
  const scrollToComments = () => {
    commentsRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const productForCart = {
    productId: Number(product.id),
    name: product.name,
    price: product.price,
    originalPrice: product.originalPrice || undefined,
    imageUrl: product.imageUrl,
    stockQuantity: product.stockQuantity
  };
  return (
    <div className="container mx-auto px-4 md:px-[5rem] py-[1rem]">
      {/* Breadcrumb */}
      <div className="flex items-center text-sm text-gray-600 mb-6">
        <Link to="/" className="hover:text-primary text-base">
          Home
        </Link>
        <span className="mx-2">/</span>
        <Link to="/shop" className="hover:text-primary text-base">
          Shop
        </Link>
        <span className="mx-2">/</span>
        <Link to={`/category/${categorySlug}`} className="hover:text-primary text-base">
          {(categorySlug || "")
            .split("-")
            .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
            .join(" ")}
        </Link>
        <span className="mx-2">/</span>
        <span className="text-primary text-base">{product.name}</span>
      </div>

      {/* Product detail */}
      <div className="flex flex-between py-6 gap-3">
        {/* Product Image Section */}
        <div className="relative space-y-4 w-[50%] flex flex-col justify-center items-center">
          {/* Main Image with Navigation Controls */}
          <div className="relative w-full h-[350px] rounded-lg overflow-hidden border border-gray-200">
            {/* Previous Button */}
            <button
              onClick={goToPreviousImage}
              className="absolute left-2 top-1/2 transform -translate-y-1/2 bg-white hover:bg-gray-100 text-gray-800 rounded-full w-10 h-10 flex items-center justify-center z-10 shadow-md transition-all duration-300"
              aria-label="Previous image"
            >
              <IoIosArrowBack className="text-xl"></IoIosArrowBack>
            </button>

            {/* Main Image */}
            <div className="w-full h-full flex items-center justify-center">
              <img
                src={selectedImage || product.imageUrl}
                alt={product.name}
                className="w-full h-full object-contain max-h-[15rem] max-w-[30rem]"
              />
            </div>

            {/* Next Button */}
            <button
              onClick={goToNextImage}
              className="absolute right-2 top-1/2 transform -translate-y-1/2 bg-white hover:bg-gray-100 text-gray-800 rounded-full w-10 h-10 flex items-center justify-center z-10 shadow-md transition-all duration-300"
              aria-label="Next image"
            >
              <IoIosArrowForward className="text-xl"></IoIosArrowForward>
            </button>
          </div>

          {/* Thumbnail Navigation with Sliding */}
          <div className="relative w-full mt-5">
            {/* Thumbnails Container */}
            <div className="flex items-center justify-center gap-3 w-full">
              {/* Previous Thumbnails Button */}
              {showPrevThumbnailsButton && (
                <button
                  onClick={goToPreviousThumbnails}
                  className="h-[35px] w-[35px] absolute left-0 top-1/2 transform -translate-y-1/2 bg-white hover:bg-gray-100 text-gray-800 rounded-full  flex items-center justify-center z-10 shadow-md transition-all duration-300"
                  aria-label="Previous thumbnails"
                >
                  <IoIosArrowBack className="text-sm"></IoIosArrowBack>
                </button>
              )}

              {/* Visible Thumbnails */}
              <div className="flex items-center justify-center gap-3 w-full px-10">
                {visibleThumbnails.map((img, relativeIndex) => {
                  const absoluteIndex = thumbnailStartIndex + relativeIndex;
                  return (
                    <button
                      key={absoluteIndex}
                      onClick={() => selectImage(img, absoluteIndex)}
                      className={`justify-center items-center rounded-lg flex w-[70px] h-[70px] border-2 overflow-hidden transition-all duration-300
                      ${currentImageIndex === absoluteIndex
                          ? "border-blue-500"
                          : "border-gray-200"
                        }`}
                    >
                      <img
                        src={img}
                        alt={`${product.name} view ${absoluteIndex + 1}`}
                        className="w-[55px] h-[55px] object-cover"
                      />
                    </button>
                  );
                })}
              </div>

              {/* Next Thumbnails Button */}
              {showNextThumbnailsButton && (
                <button
                  onClick={goToNextThumbnails}
                  className="h-[35px] w-[35px] absolute right-0 top-1/2 transform -translate-y-1/2 bg-white hover:bg-gray-100 text-gray-800 rounded-full flex items-center justify-center z-10 shadow-md transition-all duration-300"
                  aria-label="Next thumbnails"
                >
                  <IoIosArrowForward className="text-sm"></IoIosArrowForward>
                </button>
              )}
            </div>

            {/* Pagination Indicator (Optional) */}
            {allImages.length > maxVisibleThumbnails && (
              <div className="flex justify-center mt-2 gap-1">
                {Array.from({
                  length: Math.ceil(allImages.length / maxVisibleThumbnails),
                }).map((_, index) => {
                  const isActive =
                    index ===
                    Math.floor(thumbnailStartIndex / maxVisibleThumbnails);
                  return (
                    <span
                      key={index}
                      className={`block w-2 h-2 rounded-full ${isActive ? "bg-blue-500" : "bg-gray-300"
                        }`}
                    />
                  );
                })}
              </div>
            )}
          </div>
        </div>

        {/* Product Info */}
        <div className="space-y-4 px-2 w-[75%] flex flex-col justify-center h-[80%]">
          {/* Rating */}
          <div className="mb-4 flex items-center">
            {[...Array(5)].map((_, i) => (
              <span
                key={i}
                className={`text-lg ${i < product.reviewsAvg.rating
                  ? "text-[#FF9017]"
                  : "text-gray-300"}`}
              >
                ★
              </span>
            ))}
            <div className="flex items-center">
              <span className="flex items-center h-[28px] ml-2 text-cyprus text-sm font-bold">
                {product.reviewsAvg.rating} Star Rating
              </span>
              <span className="text-gray-500 text-sm ml-2">
                ({product.reviewsAvg.count} Users feedback )
              </span>
            </div>
          </div>
          <h2 className="text-xl font-medium !mb-2">{product.name}</h2>

          {/* SKU and Category */}
          <div className="text-sm text-gray-600 !mb-3">
            <div className="flex items-center justify-between">
              <p>
                Category:{" "}
                <span className="font-bold text-sm">
                  {(categorySlug || "")
                    .split("-")
                    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
                    .join(" ")}
                </span>
              </p>
              <p className="mt-1">
                Availability:{" "}
                <span
                  className={`text-sm font-bold ${product.stockQuantity > 0 ? "text-success" : "text-error"
                    }`}
                >
                  {product.stockQuantity > 0 ? "In Stock" : "Out of Stock"}
                </span>
              </p>
            </div>
          </div>

          {/* Price */}
          <div className="flex items-center !mb-[1rem] gap-2 justify-between">
            <div className="flex items-center gap-2">
              <div className="text-2xl font-semibold text-[#2DA5F3]">
                ${product.price.toLocaleString()}
              </div>
              {product.originalPrice && product.originalPrice > 0 && product.originalPrice > product.price && (
                <div className="text-gray-500 line-through self-center">
                  ${product.originalPrice.toLocaleString()}
                </div>
              )}
              {discountPercentage && (
                <div className="text-sm bg-pending text-gray-900 font-bold px-2 py-1 rounded self-center">
                  {discountPercentage}
                </div>
              )}
            </div>
          </div>

          {/* Short Desc */}
          <div className="text-sm text-gray-600 !mb-[1rem]">
            <div className="font-bold mb-2">Description: </div>
            <div className="whitespace-pre-line text-neutral-500">
              {product.description || "No short description available"}
            </div>
          </div>

          {/* Quantity Selector */}
          <div className="flex items-center !mb-[1rem]">
            <div className="px-[10px] w-[140px] h-[50px] flex items-center justify-between border border-gray-300 rounded-full mr-3">
              <button
                onClick={decreaseQuantity}
                className="px-3 py-1 text-lg cursor-pointer text-gray-600 hover:bg-gray-100 hover:text-blue-600 transition-colors duration-200 rounded-md"
                disabled={quantity <= 1}
              >
                -
              </button>
              <input
                type="text"
                value={quantity.toString().padStart(2, "0")}
                onChange={(e) => {
                  const val = parseInt(e.target.value.replace(/^0+/, ""));
                  if (!isNaN(val) && val >= 1 && val <= product.stockQuantity) {
                    handleQuantityChange(val);
                  }
                }}
                className="w-[20px] text-center focus:outline-none"
                readOnly
              />
              <button
                onClick={increaseQuantity}
                className="px-3 py-1 text-lg text-gray-600 hover:bg-gray-100 hover:text-blue-600 transition-colors duration-200 rounded-md"
                disabled={quantity >= product.stockQuantity}
              >
                +
              </button>
            </div>
            <button
              onClick={() => addToCart(productForCart)}
              disabled={product.stockQuantity <= 0}
              className="h-[50px] w-[150px] text-[14px] flex justify-center items-center font-medium px-2 rounded-full bg-[#0496FF] text-white py-3  hover:bg-blue-500 disabled:bg-gray-400 mr-3 transition-all duration-300 transform"
            >
              ADD TO CART
            </button>
            <button
              disabled={product.stockQuantity <= 0}
              className="h-[50px] w-[150px] text-[14px] flex justify-center items-center font-medium px-2 rounded-full text-[#0496FF] border border-[#0496FF] py-3 hover:bg-[#0496FF] hover:text-white transition-all duration-300 transform hover:scale-105 disabled:bg-gray-400 disabled:border-gray-400 disabled:text-gray-600 mr-3"
            >
              BUY NOW
            </button>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-wrap gap-16 mb-6 items-center justify-between">
            <div className="flex items-center gap-16">
              <button className="flex items-center px-4 py-3 rounded-lg text-gray-600 hover:text-blue-500 transition-all duration-300 ">
                <CiHeart className="mr-2 text-xl transition-transform duration-300 group-hover:scale-110" />
                Add to Wishlist
              </button>

              <button className="flex items-center px-4 py-3 rounded-lg text-gray-600 hover:text-blue-500 transition-all duration-300">
                <TbScale className="mr-2 text-xl transition-transform duration-300 group-hover:scale-110" />
                Add to Compare
              </button>
            </div>

            <button className="flex items-center px-4 py-3 rounded-lg text-gray-600 hover:text-blue-500 transition-all duration-300 ">
              <IoShareSocialOutline className="mr-2 text-xl transition-transform duration-300 group-hover:scale-110" />
              Share
            </button>
          </div>
        </div>
      </div>

      {/* Product Information (Tabs) */}
      <ProductInformation product={product} onWriteReview={scrollToComments} />

      {/* Related Products */}
      <div className="mt-12">
        <h2 className='text-[25px] font-bold'>You may also like</h2>
        <RelatedProduct
          categorySlug={categorySlug || ""}
          currentProductId={product.id.toString()}
          name={product.name}
          description={product.description}
          price={product.price}
          imageUrl={product.imageUrl}
          categoryId={product.categoryId}
          type="normal"
          id={product.id}
          slug={product.slug}
          stockQuantity={product.stockQuantity}
          originalPrice={product.originalPrice}
          discount={product.discount}
          brand={product.brand}
          tags={product.tags}
          features={product.features}
          shippingInfo={product.shippingInfo}
          reviewsAvg={product.reviewsAvg}
          orders={product.orders}
          imgSlider={product.imgSlider}
        />
      </div>

      {/* Product Comments */}
      <div ref={commentsRef}>
        <ProductComments
          productId={product.id.toString()}
          productName={product.name}
          productImage={product.imageUrl}
        />
      </div>
    </div>
  );
};

export default ProductDetail;