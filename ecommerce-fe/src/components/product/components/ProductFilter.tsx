
import React, { useState, useEffect } from 'react';
import { Category } from '../models/category.model';
import { categoryService } from '../services/category.service';
import { FaStar } from 'react-icons/fa';
import { GoChevronDown, GoChevronUp } from "react-icons/go";
interface ProductFilterProps {
  currentCategory: Category | null;
  selectedRating: number | null;
  onRatingFilter: (rating: number | null) => void;
  onPriceRangeFilter?: (minPrice: number | null, maxPrice: number | null) => void;
  onBrandFilter?: (brands: string[]) => void;
  onTagFilter?: (tag: string) => void;
}

const ProductFilter: React.FC<ProductFilterProps> = ({
  currentCategory,
  selectedRating,
  onRatingFilter,
  onPriceRangeFilter,
  onBrandFilter,
  onTagFilter
}) => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedPriceRange, setSelectedPriceRange] = useState<string>('all');
  const [selectedBrands, setSelectedBrands] = useState<string[]>([]);
  const [minPrice, setMinPrice] = useState<string>('');
  const [maxPrice, setMaxPrice] = useState<string>('');
  const [expandedSections, setExpandedSections] = useState({
    categories: true,
    priceRange: true,
    brands: true,
    rating: true,
    tags: true
  });

  // Popular brands data
  const brands = [
    { name: 'Apple', checked: false },
    { name: 'Google', checked: false },
    { name: 'Microsoft', checked: false },
    { name: 'Samsung', checked: false },
    { name: 'Dell', checked: false },
    { name: 'HP', checked: false },
    { name: 'Symphony', checked: false },
    { name: 'Xiaomi', checked: false },
    { name: 'Sony', checked: false },
    { name: 'Panasonic', checked: false },
    { name: 'LG', checked: false },
    { name: 'Intel', checked: false },
    { name: 'One Plus', checked: false },
  ];

  // Popular tags data
  const popularTags = [
    'Game', 'iPhone', 'TV', 'Asus Laptops', 
    'Macbook', 'SSD', 'Graphics Card', 
    'Power Bank', 'Smart TV', 'Speaker',
    'Tablet', 'Microwave', 'Samsung'
  ];

  // Price ranges
  const priceRanges = [
    { id: 'all', label: 'All Price', min: null, max: null },
    { id: 'under20', label: 'Under $20', min: 0, max: 20 },
    { id: '25to100', label: '$25 to $100', min: 25, max: 100 },
    { id: '100to300', label: '$100 to $300', min: 100, max: 300 },
    { id: '300to500', label: '$300 to $500', min: 300, max: 500 },
    { id: '500to1000', label: '$500 to $1,000', min: 500, max: 1000 },
    { id: '1000to10000', label: '$1,000 to $10,000', min: 1000, max: 10000 },
  ];

  useEffect(() => {
    const fetchCategories = async () => {
      setIsLoading(true);
      try {
        const categoriesData = await categoryService.getCategories();
        if (Array.isArray(categoriesData)) {
          setCategories(categoriesData);
        } else {
          console.error('Invalid categories data received:', categoriesData);
          setCategories([]);
        }
      } catch (error) {
        console.error('Error fetching categories:', error);
        setCategories([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchCategories();
  }, []);


  const toggleSection = (section: keyof typeof expandedSections) => {
    setExpandedSections({
      ...expandedSections,
      [section]: !expandedSections[section]
    });
  };

  const handlePriceRangeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id } = e.target;
    setSelectedPriceRange(id);
    
    const selectedRange = priceRanges.find(range => range.id === id);
    if (selectedRange && onPriceRangeFilter) {
      onPriceRangeFilter(selectedRange.min, selectedRange.max);
      // Reset custom price inputs when selecting a predefined range
      setMinPrice('');
      setMaxPrice('');
    }
  };

  const handleCustomPriceFilter = () => {
    const min = minPrice ? parseFloat(minPrice) : null;
    const max = maxPrice ? parseFloat(maxPrice) : null;
    
    if (onPriceRangeFilter && (min !== null || max !== null)) {
      onPriceRangeFilter(min, max);
      // Reset predefined price range selection
      setSelectedPriceRange('');
    }
  };

  const handleBrandChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value, checked } = e.target;
    let updatedBrands: string[];
    
    if (checked) {
      updatedBrands = [...selectedBrands, value];
    } else {
      updatedBrands = selectedBrands.filter(brand => brand !== value);
    }
    
    setSelectedBrands(updatedBrands);
    if (onBrandFilter) {
      onBrandFilter(updatedBrands);
    }
  };

  const handleRatingClick = (rating: number) => {
    if (rating === selectedRating) {
      // If clicking the same rating, clear the filter
      onRatingFilter(null);
    } else {
      onRatingFilter(rating);
    }
  };

  const handleTagClick = (tag: string) => {
    if (onTagFilter) {
      onTagFilter(tag);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-sm">
      {/* Categories Section */}
      <div className="mb-[1.5rem]">
        <div
          className="flex justify-between items-center cursor-pointer"
          onClick={() => toggleSection("categories")}
        >
          <h3 className="text-[17px] font-medium font-sans text-gray-800">
            CATEGORY
          </h3>
          <span>
            {expandedSections.categories ? <GoChevronDown /> : <GoChevronUp />}
          </span>
        </div>

        {expandedSections.categories && (
          <div className="mt-3 space-y-2 transform duration-300 ease-linear">
            {isLoading ? (
              <div className="text-sm text-gray-500">Loading categories...</div>
            ) : categories.length > 0 ? (
              categories.map((category) => (
                <div key={category.id} className="flex items-center">
                  <input
                    type="radio"
                    id={`category-${category.id}`}
                    name="category"
                    className="h-[20px] w-[20px] text-orange-500 focus:ring-orange-500 border-gray-300 rounded-full"
                    checked={currentCategory?.id === category.id}
                    readOnly
                  />
                  <label
                    htmlFor={`category-${category.id}`}
                    className="ml-2 text-[16px] text-gray-700"
                  >
                    {category.name}
                  </label>
                </div>
              ))
            ) : (
              <div className="text-sm text-gray-500">
                No categories available
              </div>
            )}
          </div>
        )}
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200 my-6"></div>

      {/* Price Range Section */}
      <div className="mb-[1.5rem] mt-5">
        <div
          className="flex justify-between items-center cursor-pointer"
          onClick={() => toggleSection("priceRange")}
        >
          <h3 className="text-[17px] font-medium font-sans text-gray-800">
            PRICE RANGE
          </h3>
          <span>
            {expandedSections.priceRange ? <GoChevronDown /> : <GoChevronUp />}
          </span>
        </div>

        {expandedSections.priceRange && (
          <div className="mt-3">
            {/* Price slider */}
            <div className="mb-[1rem]">
              <div className="relative mb-[1rem]">
                <div className="h-1 bg-gray-200 rounded-full">
                  <div
                    className="absolute h-1 bg-orange-500 rounded-full"
                    style={{ left: "0%", width: "70%" }}
                  ></div>
                  <div
                    className="absolute w-[15px] h-[15px] bg-white border-2 border-orange-500 rounded-full -mt-1.5 -ml-2"
                    style={{ left: "0%" }}
                  ></div>
                  <div
                    className="absolute w-[15px] h-[15px] bg-white border-2 border-orange-500 rounded-full -mt-1.5 -ml-2"
                    style={{ left: "70%" }}
                  ></div>
                </div>
              </div>

              {/* Custom price range inputs */}
              <div className="flex gap-3 mb-[1rem]">
                <div className="flex-1">
                  <input
                    type="text"
                    placeholder="Min price"
                    className="w-full p-2 text-sm border border-gray-300 rounded-md"
                    value={minPrice}
                    onChange={(e) => setMinPrice(e.target.value)}
                  />
                </div>
                <div className="flex-1">
                  <input
                    type="text"
                    placeholder="Max price"
                    className="w-full p-2 text-sm border border-gray-300 rounded-md"
                    value={maxPrice}
                    onChange={(e) => setMaxPrice(e.target.value)}
                  />
                </div>
              </div>

              <button
                className="w-full py-2 bg-[#0496FF] text-white rounded-md text-[16px] font-medium hover:bg-blue-500 transition duration-150"
                onClick={handleCustomPriceFilter}
              >
                Apply Filter
              </button>
            </div>

            {/* Predefined price ranges */}
            <div className="space-y-2">
              {priceRanges.map((range) => (
                <div key={range.id} className="flex items-center">
                  <input
                    type="radio"
                    id={range.id}
                    name="priceRange"
                    className="h-[20px] w-[20px] text-orange-500 focus:ring-orange-500 border-gray-300 rounded-full"
                    checked={selectedPriceRange === range.id}
                    onChange={handlePriceRangeChange}
                  />
                  <label
                    htmlFor={range.id}
                    className="ml-2 text-[16px] text-gray-700"
                  >
                    {range.label}
                  </label>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200 my-6"></div>

      {/* Brands Section */}
      <div className="mb-[1.5rem] mt-5">
        <div
          className="flex justify-between items-center cursor-pointer"
          onClick={() => toggleSection("brands")}
        >
          <h3 className="text-[17px] font-medium font-sans text-gray-800">
            POPURLAR BRANDS
          </h3>
          <span>
            {expandedSections.brands ? <GoChevronDown /> : <GoChevronUp />}
          </span>
        </div>

        {expandedSections.brands && (
          <div className="mt-3 grid grid-cols-2 gap-x-4 gap-y-2">
            {brands.map((brand) => (
              <div key={brand.name} className="flex items-center gap-y-3">
                <input
                  type="checkbox"
                  id={`brand-${brand.name}`}
                  value={brand.name}
                  className="h-[20px] w-[20px] text-orange-500 focus:ring-orange-500 border-gray-200 rounded"
                  checked={selectedBrands.includes(brand.name)}
                  onChange={handleBrandChange}
                />
                <label
                  htmlFor={`brand-${brand.name}`}
                  className="ml-2 text-sm text-gray-700"
                >
                  {brand.name}
                </label>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200 my-6"></div>

      {/* Rating Section */}
      <div className="mb-[1.5rem] mt-5">
        <div
          className="flex justify-between items-center cursor-pointer"
          onClick={() => toggleSection("rating")}
        >
          
          <h3 className="text-[17px] font-medium font-sans text-gray-800">
            RATINGS
          </h3>
          <span>
            {expandedSections.rating ? <GoChevronDown /> : <GoChevronUp />}
          </span>
        </div>

        {expandedSections.rating && (
          <div className="mt-3 space-y-2">
            {[5, 4, 3, 2, 1].map((rating) => (
              <div
                key={rating}
                className="flex items-center cursor-pointer"
                onClick={() => handleRatingClick(rating)}
              >
                <div className="flex items-center">
                  {Array.from({ length: 5 }).map((_, index) => (
                    <FaStar
                      key={index}
                      className={`${
                        index < rating ? "text-yellow-400" : "text-gray-300"
                      } ${
                        selectedRating === rating ? "scale-110" : ""
                      } w-[14px] h-[14px]`}
                    />
                  ))}
                  <span className="ml-2 text-sm text-gray-700">
                    {rating === 5 ? "& Above" : "& Up"}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200 my-6"></div>

      {/* Popular Tags Section */}
      <div className="mb-[1.5rem] mt-5">
        <div
          className="flex justify-between items-center cursor-pointer"
          onClick={() => toggleSection("tags")}
        >
          <h3 className="text-[17px] font-medium font-sans text-gray-800">
            POPURLAR TAGS
          </h3>
          <span>
            {expandedSections.tags ? <GoChevronDown /> : <GoChevronUp />}
          </span>
        </div>

        {expandedSections.tags && (
          <div className="mt-3 flex flex-wrap gap-2">
            {popularTags.map((tag) => (
              <button
                key={tag}
                onClick={() => handleTagClick(tag)}
                className="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md transition-colors"
              >
                {tag}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductFilter;