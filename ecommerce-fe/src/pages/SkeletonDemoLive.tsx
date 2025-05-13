import React, { useState, useEffect } from 'react';
import {
  SliderSkeleton,
  NewFashionSkeleton,
  GamingBannerSkeleton,
  DisplayGridSkeleton,
  BannerShowCaseSkeleton,
  CategorySkeleton,
  TopProductSkeleton,
  HappyCustomerSkeleton,
  ProductListSkeleton,
  ProductDetailSkeleton
} from '@/components/ui/skeletons';
import { Button } from '@/components/ui/button';

// Các component thật (giả mạo cho demo)
const RealHeroBanner = () => <div className="bg-blue-500 text-white h-[500px] flex items-center justify-center"><span className="text-3xl font-bold">Hero Banner Content</span></div>;
const RealNewFashion = () => <div className="bg-pink-300 h-[328px] flex items-center justify-center"><span className="text-xl font-bold">New Fashion Content</span></div>;
const RealGamingBanner = () => <div className="bg-purple-300 h-[328px] flex items-center justify-center"><span className="text-xl font-bold">Gaming Banner Content</span></div>;
const RealDisplayGrid = () => <div className="bg-green-300 h-[328px] flex items-center justify-center"><span className="text-xl font-bold">Display Grid Content</span></div>;
const RealBannerShowcase = () => <div className="bg-yellow-300 h-[200px] flex items-center justify-center"><span className="text-xl font-bold">Banner Showcase Content</span></div>;
const RealCategoryExplorer = () => <div className="bg-gray-300 h-[220px] flex items-center justify-center"><span className="text-xl font-bold">Category Explorer Content</span></div>;
const RealTopProducts = () => <div className="bg-red-300 h-[500px] flex items-center justify-center"><span className="text-xl font-bold">Top Products Content</span></div>;
const RealHappyCustomers = () => <div className="bg-blue-300 h-[400px] flex items-center justify-center"><span className="text-xl font-bold">Happy Customers Content</span></div>;
const RealProductList = () => <div className="bg-indigo-200 rounded-lg p-6 min-h-[600px] flex items-center justify-center"><span className="text-3xl font-bold">Product List Content</span></div>;
const RealProductDetail = () => <div className="bg-teal-200 rounded-lg p-6 min-h-[800px] flex items-center justify-center"><span className="text-3xl font-bold">Product Detail Content</span></div>;

const SkeletonDemoLive: React.FC = () => {
  // State quản lý trạng thái loading cho từng phần
  const [loading, setLoading] = useState({
    heroBanner: true,
    newFashion: true,
    gamingBanner: true,
    displayGrid: true,
    bannerShowcase: true,
    categoryExplorer: true,
    topProducts: true,
    happyCustomers: true,
    productList: true,
    productDetail: true
  });

  // Các hàm để điều khiển loading của từng phần
  const resetAllLoading = () => {
    setLoading({
      heroBanner: true,
      newFashion: true,
      gamingBanner: true,
      displayGrid: true,
      bannerShowcase: true,
      categoryExplorer: true,
      topProducts: true,
      happyCustomers: true,
      productList: true,
      productDetail: true
    });

    // Giả lập thời gian tải dữ liệu khác nhau cho từng thành phần
    setTimeout(() => setLoading(prev => ({ ...prev, heroBanner: false })), 2000);
    setTimeout(() => setLoading(prev => ({ ...prev, newFashion: false })), 3000);
    setTimeout(() => setLoading(prev => ({ ...prev, gamingBanner: false })), 3500);
    setTimeout(() => setLoading(prev => ({ ...prev, displayGrid: false })), 4000);
    setTimeout(() => setLoading(prev => ({ ...prev, bannerShowcase: false })), 4500);
    setTimeout(() => setLoading(prev => ({ ...prev, categoryExplorer: false })), 5000);
    setTimeout(() => setLoading(prev => ({ ...prev, topProducts: false })), 5500);
    setTimeout(() => setLoading(prev => ({ ...prev, happyCustomers: false })), 6000);
    setTimeout(() => setLoading(prev => ({ ...prev, productList: false })), 6500);
    setTimeout(() => setLoading(prev => ({ ...prev, productDetail: false })), 7000);
  };

  // Tự động bắt đầu loading khi component mount
  useEffect(() => {
    resetAllLoading();
  }, []);

  return (
    <div className="container mx-auto py-8 px-4 max-w-[1600px]">
      <div className="mb-8 space-y-4">
        <h1 className="text-3xl font-bold">Skeleton Demo (Live Simulation)</h1>
        <p className="text-gray-500">
          Xem skeleton loading trước khi nội dung thực được hiển thị.
        </p>
        <Button
          onClick={resetAllLoading}
          className="bg-blue-600 hover:bg-blue-700"
        >
          Reload All Skeletons
        </Button>
      </div>

      <div className="space-y-16">
        {/* Hero Banner */}
        <div>
          <h2 className="text-2xl font-semibold mb-4">Hero Banner</h2>
          {loading.heroBanner ? <SliderSkeleton /> : <RealHeroBanner />}
        </div>

        {/* Ads Section */}
        <div>
          <h2 className="text-2xl font-semibold mb-4">Ads Section</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              {loading.newFashion ? <NewFashionSkeleton /> : <RealNewFashion />}
              <p className="text-center mt-2 text-gray-500">New Fashion</p>
            </div>
            <div>
              {loading.gamingBanner ? <GamingBannerSkeleton /> : <RealGamingBanner />}
              <p className="text-center mt-2 text-gray-500">Gaming Banner</p>
            </div>
            <div>
              {loading.displayGrid ? <DisplayGridSkeleton /> : <RealDisplayGrid />}
              <p className="text-center mt-2 text-gray-500">Display Grid</p>
            </div>
          </div>
        </div>

        {/* Banner Showcase */}
        <div>
          <h2 className="text-2xl font-semibold mb-4">Banner Showcase</h2>
          {loading.bannerShowcase ? <BannerShowCaseSkeleton /> : <RealBannerShowcase />}
        </div>

        {/* Category Explorer */}
        <div>
          <h2 className="text-2xl font-semibold mb-4">Category Explorer</h2>
          {loading.categoryExplorer ? <CategorySkeleton variant="homepage" count={6} /> : <RealCategoryExplorer />}
        </div>

        {/* Top Products */}
        <div>
          <h2 className="text-2xl font-semibold mb-4">Top Products</h2>
          {loading.topProducts ? <TopProductSkeleton /> : <RealTopProducts />}
        </div>

        {/* Happy Customers */}
        <div>
          <h2 className="text-2xl font-semibold mb-4">Happy Customers</h2>
          {loading.happyCustomers ? (
            <div className="border border-gray-200 rounded-xl overflow-hidden">
              <HappyCustomerSkeleton />
            </div>
          ) : (
            <RealHappyCustomers />
          )}
        </div>

        {/* Product List */}
        <div className="mt-16">
          <h2 className="text-2xl font-semibold mb-4">Product List Page</h2>
          <p className="text-gray-500 mb-6 text-lg">Trang hiển thị danh sách sản phẩm với bộ lọc và phân trang</p>
          <div className="border border-gray-200 rounded-xl overflow-hidden shadow-lg">
            {loading.productList ? <ProductListSkeleton /> : <RealProductList />}
          </div>
        </div>

        {/* Product Detail */}
        <div className="mt-16 mb-16">
          <h2 className="text-2xl font-semibold mb-4">Product Detail Page</h2>
          <p className="text-gray-500 mb-6 text-lg">Trang hiển thị chi tiết sản phẩm với hình ảnh, thông tin và tabs</p>
          <div className="border border-gray-200 rounded-xl overflow-hidden shadow-lg">
            {loading.productDetail ? <ProductDetailSkeleton /> : <RealProductDetail />}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SkeletonDemoLive; 