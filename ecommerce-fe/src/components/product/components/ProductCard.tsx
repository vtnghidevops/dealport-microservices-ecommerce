// src/components/product/components/ProductCard.tsx
import React from 'react';
import { Link } from 'react-router-dom';
import { Product } from '../models/product.model';

interface ProductCardProps {
  product: Product;
}

const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  return (
    <Link to={`/${product.categorySlug}/${product.slug}`} className="block">
      <div className="border rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow">
        <img src={product.imageUrl} alt={product.name} className="w-full h-48 object-cover" />
        <div className="p-4">
          <h3 className="text-lg font-medium">{product.name}</h3>
          <p className="text-gray-600 mt-1">{product.price.toLocaleString()} $</p>
          {product.rating && (
            <div className="mt-2 flex items-center">
              <span className="text-yellow-500">★</span>
              <span className="ml-1">{product.rating}</span>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
};

export default ProductCard;