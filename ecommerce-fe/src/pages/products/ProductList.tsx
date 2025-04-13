import React, { useEffect, useState } from 'react';
import { useParams, Link,  } from 'react-router-dom';
import { Product } from '../../components/product/models/product.model';
import { Category } from '../../components/product/models/category.model';
import { productService } from '../../components/product/services/product.service';
import { categoryService } from '../../components/product/services/category.service';
import ProductFilter from '../../components/product/components/ProductFilter';
import ProductGrid from '../../components/product/components/ProductGrid';
import NotFound from '../system/NotFound';
import { CiSearch } from "react-icons/ci";
import Pagination from '@/components/common/Pagination';
import { normalizeText } from '@/utils/helpers';
const ProductListPage: React.FC = () => {
  const { categorySlug } = useParams<{ categorySlug: string }>();
  const [products, setProducts] = useState<Product[]>([]);
  const [category, setCategory] = useState<Category | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [notFound, setNotFound] = useState<boolean>(false);
  const [sortOption, setSortOption] = useState<string>('popular');
  const [filterRating, setFilterRating] = useState<number | null>(null);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const pageSize: number = 8; // Number of products per page
  // Calculate pagination
  const indexOfLastProduct = currentPage * pageSize;
  const indexOfFirstProduct = indexOfLastProduct - pageSize;
  const currentProducts = products.slice(indexOfFirstProduct, indexOfLastProduct);
  
  const [priceFilter, setPriceFilter] = useState<{min: number | null, max: number | null}>({
    min: null, 
    max: null
  });
  const [selectedBrands, setSelectedBrands] = useState<string[]>([]);
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>("");
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
  
  const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
    const query = event.target.value;
    setSearchQuery(query);
    
    if (!categorySlug) return;
    
    // Get original products and filter
    productService.getProductsByCategorySlug(categorySlug).then(products => {
      if (!products) return;
      
      const filtered = products.filter(product => {
        const searchIn = [
          product.name,
          product.description,
          product.brand,
          ...(product.tags || [])
        ].map(text => normalizeText(text || ""));
        
        return searchIn.some(text => 
          text.includes(normalizeText(query))
        );
      });
      
      setProducts(filtered);
      setCurrentPage(1); // Reset to first page when searching
    });
  };

    // Handler functions for filters
    const handlePriceRangeFilter = (min: number | null, max: number | null) => {
      setPriceFilter({ min, max });
      applyFilters();
    };
  
    const handleBrandFilter = (brands: string[]) => {
      setSelectedBrands(brands);
      applyFilters();
    };
  
    const handleTagFilter = (tag: string) => {
      setSelectedTag(currentTag => currentTag === tag ? null : tag);
      applyFilters();
    };
  
    // Combined filter function
    const applyFilters = async () => {
      if (!categorySlug) return;
      
      try {
        let filteredProducts = await productService.getProductsByCategorySlug(categorySlug);
        
        if (!filteredProducts) return;
  
        // Apply rating filter
        if (filterRating !== null) {
          filteredProducts = filteredProducts.filter(
            product => product.rating && product.rating >= filterRating
          );
        }
  
        // Apply price filter
        if (priceFilter.min !== null || priceFilter.max !== null) {
          filteredProducts = filteredProducts.filter(product => {
            const price = product.price;
            if (priceFilter.min !== null && priceFilter.max !== null) {
              return price >= priceFilter.min && price <= priceFilter.max;
            }
            if (priceFilter.min !== null) {
              return price >= priceFilter.min;
            }
            if (priceFilter.max !== null) {
              return price <= priceFilter.max;
            }
            return true;
          });
        }
  
        // Apply brand filter
        if (selectedBrands.length > 0) {
          filteredProducts = filteredProducts.filter(
            product => product.brand && selectedBrands.includes(product.brand)
          );
        }
  
        // Apply tag filter
        if (selectedTag) {
          const normalizedSelectedTag = normalizeText(selectedTag);
          filteredProducts = filteredProducts.filter(
            product => product.tags?.some(tag => 
              normalizeText(tag) === normalizedSelectedTag
            )
          );
        }
        // Apply search filter
        if (searchQuery.trim()) {
          filteredProducts = filteredProducts.filter(product => {
            const searchIn = [
              product.name,
              product.description,
              product.brand,
              ...(product.tags || [])
            ].map(text => normalizeText(text || ""));
            
            return searchIn.some(text => 
              text.includes(normalizeText(searchQuery))
            );
          });
        }
        // Apply sorting
        switch (sortOption) {
          case 'price-low':
            filteredProducts.sort((a, b) => a.price - b.price);
            break;
          case 'price-high':
            filteredProducts.sort((a, b) => b.price - a.price);
            break;
          case 'rating':
            filteredProducts.sort((a, b) => (b.rating || 0) - (a.rating || 0));
            break;
          default:
            // Default sorting (popular)
            break;
        }
  
        setProducts(filteredProducts);
        
      } catch (error) {
        console.error('Error applying filters:', error);
      }
    };
  
    useEffect(() => {
      applyFilters();
    }, [categorySlug, filterRating, priceFilter, selectedBrands, selectedTag, sortOption]);

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setSortOption(e.target.value);
    
    // Sort products based on the selected option
    const sortedProducts = [...products];
    switch (e.target.value) {
      case 'price-low':
        sortedProducts.sort((a, b) => a.price - b.price);
        break;
      case 'price-high':
        sortedProducts.sort((a, b) => b.price - a.price);
        break;
      case 'rating':
        sortedProducts.sort((a, b) => (b.rating || 0) - (a.rating || 0));
        break;
      default:
        // Default sorting (popular)
        break;
    }
    setProducts(sortedProducts);
  };

  const handleFilterByRating = (rating: number | null) => {
    setFilterRating(rating);
    
    // Reset to original products list if no rating filter
    if (rating === null) {
      productService.getProductsByCategorySlug(categorySlug || '').then(data => {
        setProducts(data || []);
      });
      return;
    }
    
    // Filter products by rating
    productService.getProductsByCategorySlug(categorySlug || '').then(data => {
      if (data) {
        const filteredProducts = data.filter((product) => 
          (product.rating && product.rating >= rating)
        );
        setProducts(filteredProducts);
      }
    });
  };
  
  if (loading) {
    return (
      <div className="container mx-auto p-4 min-h-screen flex justify-center items-center">
        <div className="flex flex-col items-center">
          <div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Loading products...</p>
        </div>
      </div>
    );
  }

  if (notFound) {
    return <NotFound />;
  }
  
  return (
    <div className="container mx-auto px-[5rem] py-[1rem]">
      {/* Breadcrumb */}
      <div className="flex items-center text-sm text-neutral-500 mb-6">
        <Link to="/" className="hover:text-primary text-base">
          Home
        </Link>
        <span className="mx-2">/</span>
        <Link to="/shop" className="hover:text-primary text-base">
          Shop
        </Link>
        <span className="mx-2">/</span>
        <span className="text-primary text-base">{category?.name}</span>
      </div>

      <div className="flex flex-col md:flex-row gap-8">
        {/* Sidebar */}
        <div className="w-full md:w-1/4">
          <ProductFilter
            currentCategory={category}
            selectedRating={filterRating}
            onRatingFilter={handleFilterByRating}
            onPriceRangeFilter={handlePriceRangeFilter}
            onBrandFilter={handleBrandFilter}
            onTagFilter={handleTagFilter}
            selectedTag={selectedTag}
          />
        </div>

        {/* Main content */}
        <div className="w-full md:w-3/4">
          <div className="bg-white p-6 rounded-lg shadow-sm mb-6">
            {/* Category header */}
            <div className="flex md:flex-row justify-between items-start md:items-center mb-6">
              <div className="border border-neutral-100-100 gap-8 flex items-center w-[424px] px-[16px] py-[12px] h-[44px] bg-white rounded-lg">
                <input
                  type="text"
                  value={searchQuery}
                  onChange={handleSearch}
                  className="text-left w-[364px] focus:outline-none bg-transparent text-cyprus placeholder:text-neutral-500"
                  placeholder="Search for anything..."
                />
                <CiSearch className="text-neutral-500 text-lg cursor-pointer" />
              </div>

              {/* Sort options */}
              <div className="flex items-center mt-4 md:mt-0">
                <label htmlFor="sort" className="text-sm text-gray-600 mr-2">
                  Sort by:
                </label>
                <select
                  id="sort"
                  value={sortOption}
                  onChange={handleSortChange}
                  className="h-[44px] border rounded-md py-1 px-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="popular">Most Popular</option>
                  <option value="price-low">Price: Low to High</option>
                  <option value="price-high">Price: High to Low</option>
                  <option value="rating">Highest Rating</option>
                </select>
              </div>
            </div>

            {/* Active Filters Section */}
            <div className="mb-[1rem] relative mt-5 px-[24px] flex flex-wrap items-center h-[44px] bg-aqua-spring rounded-lg">
              <div className="mr-2 text-sm text-neutral-500">
                Active Filters:
              </div>

              {/* Category filter tag */}
              {category && (
                <div className="flex items-center bg-white rounded-full px-3 py-1 text-sm mr-2">
                  <span>{category.name}</span>
                </div>
              )}
                
                
                {/* Search filter tag */}
              {searchQuery && (
                <div className="flex items-center bg-white rounded-full px-3 py-1 text-sm mr-2">
                  <span>Search: {searchQuery}</span>
                  <button
                    className="ml-2 text-gray-500 hover:text-gray-700"
                    onClick={() => {
                      setSearchQuery("");
                      applyFilters(); // Reset to filtered products without search
                    }}
                  >
                    ×
                  </button>
                </div>
              )}

              {/* Rating filter tag */}
              {filterRating && (
                <div className="flex items-center bg-white rounded-full px-3 py-1 text-sm mr-2">
                  <span>{filterRating} Star & Up</span>
                  <button
                    className="ml-2 text-gray-500 hover:text-gray-700"
                    onClick={() => handleFilterByRating(null)}
                  >
                    ×
                  </button>
                </div>
              )}

              {/* Price filter tag */}
              {(priceFilter.min !== null || priceFilter.max !== null) && (
                <div className="flex items-center bg-white rounded-full px-3 py-1 text-sm mr-2">
                  <span>
                    Price: ${priceFilter.min || 0} - ${priceFilter.max || "∞"}
                  </span>
                  <button
                    className="ml-2 text-gray-500 hover:text-gray-700"
                    onClick={() => handlePriceRangeFilter(null, null)}
                  >
                    ×
                  </button>
                </div>
              )}

              {/* Brand filter tags */}
              {selectedBrands.map((brand) => (
                <div
                  key={brand}
                  className="flex items-center bg-white rounded-full px-3 py-1 text-sm mr-2"
                >
                  <span>{brand}</span>
                  <button
                    className="ml-2 text-gray-500 hover:text-gray-700"
                    onClick={() =>
                      handleBrandFilter(
                        selectedBrands.filter((b) => b !== brand)
                      )
                    }
                  >
                    ×
                  </button>
                </div>
              ))}

              {/* Tag filter */}
              {selectedTag && (
                <div className="flex items-center bg-white rounded-full px-3 py-1 text-sm mr-2">
                  <span>{selectedTag}</span>
                  <button
                    className="ml-2 text-gray-500 hover:text-gray-700"
                    onClick={() => handleTagFilter(selectedTag)}
                  >
                    ×
                  </button>
                </div>
              )}
            </div>

            {/* Products */}
            {products.length > 0 ? (
              <>
                <ProductGrid products={currentProducts} />
                <Pagination
                  currentPage={currentPage}
                  totalItems={products.length}
                  pageSize={pageSize}
                  onPageChange={(page) => setCurrentPage(page)}
                  showNavigation={false}
                />
              </>
            ) : (
              <div className="flex justify-center items-center h-64 bg-gray-50 rounded-lg">
                <p className="text-gray-500">
                  No products found matching your criteria.
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductListPage;