// src/components/product/components/RelatedProduct.tsx
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Product } from '@/types/product.model';
import { productService } from '../services/product.service';
import ProductCard from './ProductCard';

interface RelatedProductProps extends Product {
  currentProductId: string;
}

const RelatedProduct: React.FC<RelatedProductProps> = ({ 
  categorySlug, 
  currentProductId,
  name = "You may also like" 
}) => {
  const [relatedProducts, setRelatedProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchRelatedProducts = async () => {
      try {
        const products = await productService.getProductsByCategorySlug(categorySlug);
        // Filter out current product and limit to 4 related products
        const filtered = products
          ? products
              .filter(product => product.id !== currentProductId)
              .slice(0, 4)
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

  if (loading) {
    return (
      <div className="flex justify-center items-center h-40 w-full">
        <div className="animate-pulse flex space-x-4">
          <div className="flex-1 space-y-4 py-1">
            <div className="h-4 bg-gray-200 rounded w-3/4"></div>
            <div className="space-y-2">
              <div className="h-4 bg-gray-200 rounded"></div>
              <div className="h-4 bg-gray-200 rounded w-5/6"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (relatedProducts.length === 0) {
    return null;
  }

  return (
    <div className="py-10">
      <h2 className="text-2xl font-bold mb-6 text-gray-800  pb-2">{name}</h2>
      <div className="flex gap-16">
        {relatedProducts.map((product) => (
          <Link
        key={product.id}
        to={`/${product.categorySlug}/${product.slug}`}
        className="w-[15rem] min-h-[320px]"
          >
        <ProductCard product={product} />
          </Link>
        ))}
      </div>
    </div>
  );
};

export default RelatedProduct;