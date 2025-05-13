import { useState } from 'react';
import {
  ProductCardSkeletonGrid,
  ProductDetailSkeleton,
  ProductGridSkeleton,
  SidebarSkeleton,
  FilterSidebarSkeleton,
  SliderSkeleton,
  BannerShowCaseSkeleton,
  DisplayGridSkeleton,
  GamingBannerSkeleton,
  HappyCustomerSkeleton,
  NewFashionSkeleton,
  CategorySkeleton,
  TopProductSkeleton
} from '@/components/ui/skeletons';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

const SkeletonDemo = () => {
  const [loading, setLoading] = useState(true);

  const toggleLoading = () => {
    setLoading(!loading);
  };

  return (
    <div className="container mx-auto py-8 px-4">
      <div className="mb-8 space-y-4">
        <h1 className="text-3xl font-bold">Skeleton Loading Demo</h1>
        <p className="text-gray-500">
          Demonstrating skeleton loading states using Shadcn UI components
        </p>
        <div className="flex flex-wrap gap-3">
          <Button
            onClick={toggleLoading}
            variant={loading ? "destructive" : "default"}
          >
            {loading ? "Hide Skeletons" : "Show Skeletons"}
          </Button>

          <Link to="/skeleton-demo-live">
            <Button variant="outline" className="bg-blue-50 hover:bg-blue-100">
              View Live Simulation Demo
            </Button>
          </Link>
        </div>
      </div>

      <Tabs defaultValue="overall" className="w-full">
        <TabsList className="mb-4 flex flex-wrap">
          <TabsTrigger value="overall">Overall Layout</TabsTrigger>
          <TabsTrigger value="homepage">Homepage Elements</TabsTrigger>
          <TabsTrigger value="slider">Sliders</TabsTrigger>
          <TabsTrigger value="product-grid">Product Grid</TabsTrigger>
          <TabsTrigger value="product-detail">Product Detail</TabsTrigger>
          <TabsTrigger value="sidebar">Sidebars</TabsTrigger>
          <TabsTrigger value="product-card">Product Card</TabsTrigger>
        </TabsList>

        {/* New Overall Layout Tab */}
        <TabsContent value="overall" className="space-y-10">
          <div className="space-y-6">
            <h2 className="text-2xl font-semibold">Slider</h2>
            {loading ? (
              <SliderSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[300px] flex items-center justify-center">
                <p className="text-lg font-medium">Main banner slider content would appear here</p>
              </div>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="space-y-4">
              <h2 className="text-xl font-semibold">New Fashion</h2>
              {loading ? (
                <NewFashionSkeleton />
              ) : (
                <div className="h-[328px] bg-gray-100 p-8 rounded-lg text-center flex items-center justify-center">
                  <p className="text-lg font-medium">New Fashion content</p>
                </div>
              )}
            </div>
            <div className="space-y-4">
              <h2 className="text-xl font-semibold">Gaming Accessories</h2>
              {loading ? (
                <GamingBannerSkeleton />
              ) : (
                <div className="h-[328px] bg-gray-100 p-8 rounded-lg text-center flex items-center justify-center">
                  <p className="text-lg font-medium">Gaming accessories content</p>
                </div>
              )}
            </div>
            <div className="space-y-4">
              <h2 className="text-xl font-semibold">Featured Products</h2>
              {loading ? (
                <DisplayGridSkeleton />
              ) : (
                <div className="h-[328px] bg-gray-100 p-8 rounded-lg text-center flex items-center justify-center">
                  <p className="text-lg font-medium">Featured products content</p>
                </div>
              )}
            </div>
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Banner Showcase</h2>
            {loading ? (
              <BannerShowCaseSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[200px] flex items-center justify-center">
                <p className="text-lg font-medium">Banner showcase content would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Featured Products</h2>
            {loading ? (
              <ProductCardSkeletonGrid count={5} />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center">
                <p className="text-lg font-medium">Featured product cards would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Categories 1</h2>
            {loading ? (
              <CategorySkeleton count={5} />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[150px] flex items-center justify-center">
                <p className="text-lg font-medium">Category listings would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Best Selling Products</h2>
            {loading ? (
              <TopProductSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[500px] flex items-center justify-center">
                <p className="text-lg font-medium">Best selling products grid would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Recommended Products</h2>
            {loading ? (
              <ProductCardSkeletonGrid count={5} />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center">
                <p className="text-lg font-medium">Recommended product cards would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Happy Customers</h2>
            {loading ? (
              <div className="border border-gray-200 rounded-xl overflow-hidden">
                <HappyCustomerSkeleton />
              </div>
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[400px] flex items-center justify-center">
                <p className="text-lg font-medium">Happy customers testimonials would appear here</p>
              </div>
            )}
          </div>
        </TabsContent>

        <TabsContent value="homepage" className="space-y-10">
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Display Grid Skeleton</h2>
            {loading ? (
              <DisplayGridSkeleton />
            ) : (
              <div className="h-[328px] w-[392px] bg-gray-100 p-8 rounded-lg text-center flex items-center justify-center">
                <p className="text-lg font-medium">Actual Display Grid would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Gaming Banner Skeleton</h2>
            {loading ? (
              <GamingBannerSkeleton />
            ) : (
              <div className="h-[328px] w-[464px] bg-gray-100 p-8 rounded-lg text-center flex items-center justify-center">
                <p className="text-lg font-medium">Actual Gaming Banner would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Categories Skeleton</h2>
            {loading ? (
              <CategorySkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[200px] flex items-center justify-center">
                <p className="text-lg font-medium">Actual category items would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Best Selling Products Skeleton</h2>
            {loading ? (
              <TopProductSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[500px] flex items-center justify-center">
                <p className="text-lg font-medium">Actual best selling products would appear here</p>
              </div>
            )}
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Happy Customers Skeleton</h2>
            {loading ? (
              <div className="border border-gray-200 rounded-xl overflow-hidden">
                <HappyCustomerSkeleton />
              </div>
            ) : (
              <div className="h-[500px] bg-gray-100 p-8 rounded-lg text-center flex items-center justify-center">
                <p className="text-lg font-medium">Actual Happy Customers section would appear here</p>
              </div>
            )}
          </div>
        </TabsContent>

        <TabsContent value="slider" className="space-y-10">
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Main Slider Skeleton</h2>
            {loading ? (
              <SliderSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center h-[500px] flex items-center justify-center">
                <p className="text-lg font-medium">Actual slider content would appear here</p>
              </div>
            )}
          </div>
        </TabsContent>

        <TabsContent value="product-grid" className="space-y-4">
          <h2 className="text-2xl font-semibold">Product Grid Skeleton</h2>
          {loading ? (
            <ProductGridSkeleton columns={4} rows={2} />
          ) : (
            <div className="bg-gray-100 p-8 rounded-lg text-center">
              <p className="text-lg font-medium">Actual product grid content would appear here</p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="product-detail" className="space-y-4">
          <h2 className="text-2xl font-semibold">Product Detail Skeleton</h2>
          {loading ? (
            <ProductDetailSkeleton />
          ) : (
            <div className="bg-gray-100 p-8 rounded-lg text-center">
              <p className="text-lg font-medium">Actual product detail content would appear here</p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="sidebar" className="space-y-4">
          <h2 className="text-2xl font-semibold">Navigation Sidebar Skeleton</h2>
          <div className="flex flex-wrap gap-4">
            {loading ? (
              <SidebarSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center w-full md:w-[260px]">
                <p className="text-lg font-medium">Actual sidebar content would appear here</p>
              </div>
            )}
            <div className="flex-1 bg-gray-50 rounded-lg p-8 flex items-center justify-center">
              <p className="text-gray-400">Main content area</p>
            </div>
          </div>

          <h2 className="text-2xl font-semibold mt-8">Filter Sidebar Skeleton</h2>
          <div className="flex flex-wrap gap-4">
            {loading ? (
              <FilterSidebarSkeleton />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center w-full md:w-[300px]">
                <p className="text-lg font-medium">Actual filter sidebar content would appear here</p>
              </div>
            )}
            <div className="flex-1 bg-gray-50 rounded-lg p-8 flex items-center justify-center">
              <p className="text-gray-400">Product grid area</p>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="product-card" className="space-y-4">
          <h2 className="text-2xl font-semibold">Product Card Skeleton</h2>
          <div className="max-w-sm mx-auto">
            {loading ? (
              <ProductCardSkeletonGrid count={1} />
            ) : (
              <div className="bg-gray-100 p-8 rounded-lg text-center">
                <p className="text-lg font-medium">Actual product card content would appear here</p>
              </div>
            )}
          </div>

          <h2 className="text-2xl font-semibold mt-8">Product Card Grid</h2>
          {loading ? (
            <ProductCardSkeletonGrid count={8} />
          ) : (
            <div className="bg-gray-100 p-8 rounded-lg text-center">
              <p className="text-lg font-medium">Actual product cards would appear here</p>
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default SkeletonDemo; 