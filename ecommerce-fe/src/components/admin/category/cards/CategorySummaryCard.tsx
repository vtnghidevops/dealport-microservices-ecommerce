import React, { useEffect, useState } from 'react';
import { CategorySummaryData } from './models/card.model';
import { CategoryService } from '../services/category.service';

/**
 * Component that displays category summary information in cards
 * Shows key metrics about categories
 */
const CategorySummaryCard: React.FC = () => {
  const [summaryData, setSummaryData] = useState<CategorySummaryData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    // Fetch summary data when component mounts
    const fetchSummaryData = async () => {
      try {
        const data = await CategoryService.getCategorySummary();
        setSummaryData(data);
      } catch (error) {
        console.error('Error fetching category summary:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchSummaryData();
  }, []);

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="bg-white p-6 rounded-lg shadow-sm animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
            <div className="h-6 bg-gray-200 rounded w-2/4 mb-6"></div>
          </div>
        ))}
      </div>
    );
  }

  if (!summaryData) {
    return (
      <div className="bg-white p-6 rounded-lg shadow-sm">
        <p className="text-gray-500">No summary data available</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      {/* Total Categories Card */}
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <h3 className="text-sm font-medium text-gray-500 mb-1">Total Categories</h3>
        <p className="text-2xl font-semibold text-gray-800">{summaryData.totalCategories}</p>
      </div>

      {/* Active Categories Card */}
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <h3 className="text-sm font-medium text-gray-500 mb-1">Active Categories</h3>
        <p className="text-2xl font-semibold text-green-600">{summaryData.activeCategories}</p>
      </div>

      {/* Featured Categories Card */}
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <h3 className="text-sm font-medium text-gray-500 mb-1">Featured Categories</h3>
        <p className="text-2xl font-semibold text-blue-600">{summaryData.featuredCategories}</p>
      </div>

      {/* Popular Category Card */}
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <h3 className="text-sm font-medium text-gray-500 mb-1">Popular Category</h3>
        <p className="text-lg font-semibold text-gray-800">{summaryData.popularCategory.name}</p>
        <p className="text-sm text-gray-500">{summaryData.popularCategory.productCount} Products</p>
      </div>
    </div>
  );
};

export default CategorySummaryCard;