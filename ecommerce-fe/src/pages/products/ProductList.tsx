import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Product } from '../../components/product/models/product.model';
import { Category } from '../../components/product/models/category.model';
import { productService } from '../../components/product/services/product.service';
import { categoryService } from '../../components/product/services/category.service';
import ProductGrid from '../../components/product/components/ProductGrid';
import NotFound from '../system/NotFound';

const ProductList: React.FC = () => {
  const { categorySlug } = useParams<{ categorySlug: string }>();
  const [products, setProducts] = useState<Product[]>([]);
  const [category, setCategory] = useState<Category | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [notFound, setNotFound] = useState<boolean>(false);
  
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        if (categorySlug) {
          const categoryData = await categoryService.getCategoryBySlug(categorySlug);
          if (!categoryData) {
            setNotFound(true);
            return;
          }
          setCategory(categoryData);
          
          const productsData = await productService.getProductsByCategorySlug(categorySlug);
          setProducts(productsData || []);
        }
      } catch (error) {
        console.error('Error fetching data:', error);
        setNotFound(true);
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, [categorySlug]);
  
  if (loading) {
    return <div className="flex justify-center items-center h-64">Loading...</div>;
  }

  if (notFound) {
    return <NotFound />;
  }
  
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-2">{category?.name}</h1>
      {category?.description && (
        <p className="mb-6 text-gray-600">{category.description}</p>
      )}
      
      <ProductGrid products={products} />
    </div>
  );
};

export default ProductList;