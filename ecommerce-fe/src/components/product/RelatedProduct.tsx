// src/components/product/components/RelatedProduct.tsx
import React, { useEffect, useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import { Product } from '@/types/product.model';
import ProductService from '@/services/product/product.service';
import ProductCard from './ProductCard';
import Loading from '@/components/shared/Loading';
import { IoIosArrowBack, IoIosArrowForward } from 'react-icons/io';

interface RelatedProductProps extends Product {
  currentProductId: string;
}

const RelatedProduct: React.FC<RelatedProductProps> = ({
  categorySlug,
  currentProductId,
  // name = "You may also like"
}) => {
  const [relatedProducts, setRelatedProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [showLeftArrow, setShowLeftArrow] = useState<boolean>(false);
  const [showRightArrow, setShowRightArrow] = useState<boolean>(true);

  useEffect(() => {
    const fetchRelatedProducts = async () => {
      try {
        const products = await ProductService.getProductsByCategorySlug(categorySlug);
        // Filter out current product and limit to 8 related products
        const filtered = products.products
          ? products.products
            .filter((product: Product) => String(product.id) !== currentProductId)
            .slice(0, 8)
          : [];
        setRelatedProducts(filtered);
      } catch (error) {
        console.error('Error fetching related products:', error);
      } finally {
        setLoading(false);
      }
    };

    if (categorySlug) {
      fetchRelatedProducts();
    }
  }, [categorySlug, currentProductId]);

  // Handle scroll arrows visibility
  useEffect(() => {
    const checkScrollPosition = () => {
      if (!scrollContainerRef.current) return;

      const { scrollLeft, scrollWidth, clientWidth } = scrollContainerRef.current;
      setShowLeftArrow(scrollLeft > 0);
      setShowRightArrow(scrollLeft < scrollWidth - clientWidth - 10); // 10px buffer
    };

    const scrollContainer = scrollContainerRef.current;
    if (scrollContainer) {
      scrollContainer.addEventListener('scroll', checkScrollPosition);
      // Initial check
      checkScrollPosition();
    }

    return () => {
      if (scrollContainer) {
        scrollContainer.removeEventListener('scroll', checkScrollPosition);
      }
    };
  }, [relatedProducts]);

  const scrollLeft = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: -500,
        behavior: 'smooth'
      });
    }
  };

  const scrollRight = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: 500,
        behavior: 'smooth'
      });
    }
  };

  if (loading) {
    return <Loading />;
  }

  if (relatedProducts.length === 0) {
    return null;
  }

  return (
    <div className="py-10">
      <div className="relative">
        {/* Left Arrow */}
        {showLeftArrow && (
          <button
            onClick={scrollLeft}
            className="absolute left-0 top-1/2 transform -translate-y-1/2 z-10 bg-white rounded-full w-10 h-10 shadow-md flex items-center justify-center -ml-5"
            aria-label="Scroll left"
          >
            <IoIosArrowBack className="text-gray-600 text-xl" />
          </button>
        )}

        {/* Product Slider */}
        <div
          ref={scrollContainerRef}
          className="flex overflow-hidden pb-4 snap-x snap-mandatory px-2"
        >
          {relatedProducts.map((product) => (
            <div
              key={product.id}
              className="mr-5 w-[250px] min-w-[250px] max-w-[250px] h-[380px] min-h-[380px] max-h-[380px] flex-shrink-0 snap-start"
            >
              <Link
                to={`/category/${product.categorySlug}/${product.slug}`}
                className="w-full h-full block"
              >
                <ProductCard product={product} />
              </Link>
            </div>
          ))}
        </div>

        {/* Right Arrow */}
        {showRightArrow && (
          <button
            onClick={scrollRight}
            className="absolute right-0 top-1/2 transform -translate-y-1/2 z-10 bg-white rounded-full w-10 h-10 shadow-md flex items-center justify-center -mr-5"
            aria-label="Scroll right"
          >
            <IoIosArrowForward className="text-gray-600 text-xl" />
          </button>
        )}
      </div>
    </div>
  );
};

export default RelatedProduct;