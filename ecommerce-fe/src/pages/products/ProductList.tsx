import React, { useEffect, useState } from 'react';
import { useParams, Link, } from 'react-router-dom';
import { Product } from '@/types/product.model';
import { Category } from '@/types/category.model';
import ProductService from '@/services/product/product.service';
import { CategoryService } from '@/services/product/product.service';
import ProductFilter from '@/components/product/ProductFilter';
import ProductGrid from '@/components/product/ProductGrid';
import NotFound from '../system/NotFound';
import { CiSearch } from "react-icons/ci";
import Pagination from '@/components/common/Pagination';
import { normalizeText } from '@/utils/helpers';
import { ProductListSkeleton } from '@/components/ui/skeletons';

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

  const [priceFilter, setPriceFilter] = useState<{ min: number | null, max: number | null }>({
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
          const categoryData = await CategoryService.getCategoryBySlug(categorySlug);
          if (!categoryData) {
            setNotFound(true);
            return;
          }
          setCategory(categoryData);

          const productsData = await ProductService.getProductsByCategorySlug(categorySlug);
          setProducts(productsData.products || []);
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
  //console.log("products in category list", products)
  //console.log("category in category list", category)

  const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
    const query = event.target.value;
    setSearchQuery(query);

    if (!categorySlug) return;

    // Get original products and filter
    ProductService.getProductsByCategorySlug(categorySlug).then(products => {
      if (!products) return;

      const filtered = products.products.filter(product => {
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
  };

  const handleBrandFilter = (brands: string[]) => {
    setSelectedBrands(brands);
  };

  const handleTagFilter = (tag: string) => {
    setSelectedTag(currentTag => currentTag === tag ? null : tag);
  };

  // Combined filter function
  const applyFilters = async () => {
    if (!categorySlug) return;

    try {
      const response = await ProductService.getProductsByCategorySlug(categorySlug);

      if (!response) return;

      let filteredProducts = [...response.products];

      // Apply rating filter (áp dụng bộ lọc đánh giá)
      if (filterRating !== null) {
        filteredProducts = filteredProducts.filter(
          (product: Product) => product.reviewsAvg?.rating >= filterRating
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
          (product: Product) => product.brand && selectedBrands.includes(product.brand)
        );
      }

      // Apply tag filter
      if (selectedTag) {
        const normalizedSelectedTag = normalizeText(selectedTag);
        filteredProducts = filteredProducts.filter(
          (product: Product) => product.tags?.some(tag =>
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
          filteredProducts.sort((a, b) => (b.reviewsAvg?.rating || 0) - (a.reviewsAvg?.rating || 0));
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
  };

  const handleFilterByRating = (rating: number | null) => {
    setFilterRating(rating);
  };

  if (loading) {
    return <ProductListSkeleton />;
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
        <Link to="/products" className="hover:text-primary text-base">
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