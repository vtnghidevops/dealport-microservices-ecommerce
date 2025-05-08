// src/components/product/components/RelatedProduct.tsx
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Product } from '@/types/product.model';
import ProductService from '@/services/product/product.service';
import ProductCard from './ProductCard';
import Loading from '@/components/shared/Loading';

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

  useEffect(() => {
    const fetchRelatedProducts = async () => {
      try {
        const products = await ProductService.getProductsByCategorySlug(categorySlug);
        // Filter out current product and limit to 4 related products
        const filtered = products.products
          ? products.products
            .filter((product: Product) => String(product.id) === currentProductId)
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
    return <Loading />;
  }

  if (relatedProducts.length === 0) {
    return null;
  }

  return (
    <div className="py-10">

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