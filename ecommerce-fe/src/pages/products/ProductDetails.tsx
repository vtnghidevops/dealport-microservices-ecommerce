// src/pages/product/ProductDetail.tsx
import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Product } from '../../components/product/models/product.model';
import { productService } from '../../components/product/services/product.service';
import NotFound from '../system/NotFound';
const ProductDetail: React.FC = () => {
  const { categorySlug, productSlug } = useParams<{
    categorySlug: string;
    productSlug: string;
  }>();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  
  useEffect(() => {
    const fetchProduct = async () => {
      try {
        if (categorySlug && productSlug) {
          const data = await productService.getProductByCategoryAndSlug(categorySlug, productSlug);
          setProduct(data || null);
        }
      } catch (error) {
        console.error('Error fetching product details:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchProduct();
  }, [categorySlug, productSlug]);
  
  if (loading) {
    return <div className="flex justify-center items-center h-64">Loading...</div>;
  }
  
  if (!product) {
    return <NotFound />;
  }
  
  return (
    <div className="container mx-auto p-4">
      {/* Breadcrumb */}
      <div className="flex items-center text-sm text-gray-600 mb-6">
        <Link to="/" className="hover:text-blue-600">Trang chủ</Link>
        <span className="mx-2">/</span>
        <Link to={`/${categorySlug}`} className="hover:text-blue-600">
          {(categorySlug || '').split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')}
        </Link>
        <span className="mx-2">/</span>
        <span className="text-gray-900">{product.name}</span>
      </div>
      
      {/* Product detail */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="rounded-lg overflow-hidden">
          <img
            src={product.imageUrl}
            alt={product.name}
            className="w-full h-auto object-cover"
          />
        </div>
        
        <div>
          <h1 className="text-3xl font-bold mb-2">{product.name}</h1>
          <div className="text-2xl font-semibold text-red-600 mb-4">
            {product.price.toLocaleString()} ₫
          </div>
          
          {product.rating && (
            <div className="mb-4 flex items-center">
              <span className="text-yellow-500">★</span>
              <span className="ml-1">{product.rating}</span>
            </div>
          )}
          
          <div className="mb-4">
            <h2 className="text-xl font-semibold mb-2">Mô tả sản phẩm</h2>
            <p className="text-gray-700">{product.description}</p>
          </div>
          
          <div className="mb-6">
            <span className={`${product.stock > 0 ? 'text-green-600' : 'text-red-600'}`}>
              {product.stock > 0 ? `Còn hàng (${product.stock})` : 'Hết hàng'}
            </span>
          </div>
          
          <button
            disabled={product.stock <= 0}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:bg-gray-400"
          >
            Thêm vào giỏ hàng
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProductDetail;