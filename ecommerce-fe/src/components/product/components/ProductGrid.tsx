// export default ProductGrid;
import React from 'react';
import ProductCard from './ProductCard';
import { useNavigate } from 'react-router-dom';
import { Product } from '@/types/product.model';

interface ProductGridProps {
  products: Product[];
}

const ProductGrid: React.FC<ProductGridProps> = ({ products }) => {
  const navigate = useNavigate();
  
  const handleProductClick = (product: Product, e: React.MouseEvent) => {
    e.preventDefault();
    
    // Construct the product URL
    const productUrl = `/${product.categorySlug}/${product.slug}`;
    
    // Navigate programmatically
    navigate(productUrl);
    
    // Use setTimeout to ensure scrolling happens after navigation starts
    setTimeout(() => {
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    }, 100);
  };
  
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-16">
      {products.map(product => (
        <ProductCard 
          product={product}
          onClick={(e) => handleProductClick(product, e)}
        />
      ))}
    </div>
  );
};

export default ProductGrid;