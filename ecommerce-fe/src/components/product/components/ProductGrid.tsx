// import React from 'react';
// import ProductCard from './ProductCard';
// import { Product } from '../models/product.model';
// import { useNavigate } from 'react-router-dom';
// interface ProductGridProps {
//   products: Product[];
// }

// const ProductGrid: React.FC<ProductGridProps> = ({ products }) => {
//   // Handler for product click
// const navigate = useNavigate();
// const handleProductClick = (productUrl: string, e: React.MouseEvent) => {
//   e.preventDefault();
  
//   const productUrl = `/${product.categorySlug}/${product.slug}`;
//   // Navigate programmatically
//   navigate(productUrl);
  
//   // Scroll to top
//   window.scrollTo({
//     top: 0,
//     behavior: 'smooth'
//   });
// };
//   return (
//     <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-16">
//       {products.map(product => (
//         <ProductCard key={product.id} product={product} onClick={(e) => handleProductClick(product.slug, e)}/>
//       ))}
//     </div>
//   );
// };

// export default ProductGrid;
import React from 'react';
import ProductCard from './ProductCard';
import { Product } from '../models/product.model';
import { useNavigate } from 'react-router-dom';

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
          key={product.id} 
          product={product} 
          onClick={(e) => handleProductClick(product, e)}
        />
      ))}
    </div>
  );
};

export default ProductGrid;