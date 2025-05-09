import { Skeleton } from "@/components/ui/skeleton";

export const ProductDetailSkeleton = () => {
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="flex items-center text-sm text-gray-600 mb-6">
        <Skeleton className="h-4 w-14 rounded mr-1" />
        <span className="mx-2">/</span>
        <Skeleton className="h-4 w-12 rounded mr-1" />
        <span className="mx-2">/</span>
        <Skeleton className="h-4 w-24 rounded" />
      </div>

      {/* Product detail - Two columns layout */}
      <div className="flex flex-col md:flex-row py-6 gap-8 mb-5">
        {/* Product Image Section - Left Column */}
        <div className="relative w-full md:w-1/2 flex flex-col justify-center items-center mb-5">
          {/* Main Image with Navigation Controls */}
          <div className="relative w-full h-[300px] md:h-[400px] rounded-lg overflow-hidden border border-gray-200 shadow-md">
            <Skeleton className="h-full w-full" />
          </div>

          {/* Thumbnail Navigation with Sliding */}
          <div className="w-full mt-5 mb-5">
            <div className="flex items-center justify-center gap-3 w-full overflow-x-auto pb-2">
              {Array(4).fill(0).map((_, index) => (
                <Skeleton key={index} className="w-[70px] h-[70px] md:w-[80px] md:h-[80px] rounded-lg shadow-md flex-shrink-0" />
              ))}
            </div>
          </div>
        </div>

        {/* Product Details Section - Right Column */}
        <div className="w-full md:w-1/2 space-y-5 bg-white p-4 md:p-6 rounded-lg shadow-md">
          {/* Rating stars */}
          <div className="flex items-center gap-2">
            <div className="flex">
              {Array(5).fill(0).map((_, i) => (
                <Skeleton key={i} className="h-5 w-5 rounded-full mr-1" />
              ))}
            </div>
            <Skeleton className="h-5 w-32 rounded ml-2" />
          </div>

          {/* Title */}
          <Skeleton className="h-8 w-full md:w-4/5 rounded mb-4" />

          {/* Price section with discount */}
          <div className="flex items-center gap-4 mb-5">
            <Skeleton className="h-9 w-32 rounded" />
            <Skeleton className="h-7 w-24 rounded" />
            <Skeleton className="h-7 w-20 rounded-full bg-yellow-100" />
          </div>

          {/* Stock and category info */}
          <div className="flex flex-col sm:flex-row sm:justify-between gap-3 mb-4">
            <div className="flex items-center gap-2">
              <Skeleton className="h-6 w-24 rounded" />
              <Skeleton className="h-6 w-28 rounded font-bold" />
            </div>
            <div className="flex items-center gap-2">
              <Skeleton className="h-6 w-24 rounded" />
              <Skeleton className="h-6 w-28 rounded-full bg-green-100" />
            </div>
          </div>

          {/* Description */}
          <div className="space-y-3 mb-5">
            <Skeleton className="h-6 w-32 rounded mb-2" />
            <Skeleton className="h-5 w-full rounded" />
            <Skeleton className="h-5 w-full rounded" />
            <Skeleton className="h-5 w-3/4 rounded" />
          </div>

          {/* Divider */}
          <div className="my-5 border-t border-gray-200"></div>

          {/* Quantity control */}
          <div className="mb-5">
            <div className="flex flex-wrap items-center gap-4">
              <div className="flex items-center border border-gray-300 rounded-full p-1 w-[140px]">
                <Skeleton className="h-[42px] w-[140px] rounded-full" />
              </div>
              <Skeleton className="h-[46px] w-[140px] rounded-full" />
              <Skeleton className="h-[46px] w-[140px] rounded-full" />
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex flex-wrap items-center gap-4 pt-4 mb-6">
            <Skeleton className="h-[50px] w-full sm:w-[180px] rounded-lg" />
            <Skeleton className="h-[50px] w-full sm:w-[180px] rounded-lg" />
          </div>

          {/* Product meta info */}
          <div className="space-y-4 pt-4 border-t border-gray-200">
            <div className="flex gap-3">
              <Skeleton className="h-6 w-24 rounded" />
              <Skeleton className="h-6 w-36 rounded" />
            </div>
            <div className="flex gap-3">
              <Skeleton className="h-6 w-24 rounded" />
              <div className="flex flex-wrap gap-2">
                <Skeleton className="h-6 w-20 rounded-full" />
                <Skeleton className="h-6 w-20 rounded-full" />
                <Skeleton className="h-6 w-20 rounded-full" />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Product tabs (Description, Reviews, etc.) */}
      <div className="mt-8 mb-10">
        {/* Tab navigation */}
        <div className="border-b border-gray-200 mb-6">
          <div className="flex space-x-4 md:space-x-8 overflow-x-auto">
            <Skeleton className="h-10 w-32 rounded-t-lg" />
            <Skeleton className="h-10 w-32 rounded" />
            <Skeleton className="h-10 w-32 rounded" />
          </div>
        </div>

        {/* Tab content */}
        <div className="space-y-5 p-5 bg-white rounded-lg shadow-md mb-10">
          <Skeleton className="h-7 w-64 rounded mb-4" />
          <div className="space-y-4">
            <Skeleton className="h-5 w-full rounded" />
            <Skeleton className="h-5 w-full rounded" />
            <Skeleton className="h-5 w-full rounded" />
            <Skeleton className="h-5 w-3/4 rounded" />
          </div>
        </div>
      </div>

      {/* Related products section */}
      <div className="mt-12 space-y-6 mb-10">
        <div className="flex justify-between items-center mb-6">
          <Skeleton className="h-8 w-56 rounded" />
          <Skeleton className="h-10 w-32 rounded-full" />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
          {Array(4).fill(0).map((_, index) => (
            <div key={index} className="bg-white rounded-lg shadow-md overflow-hidden mb-5">
              <Skeleton className="h-[200px] w-full" />
              <div className="p-4 space-y-3">
                <Skeleton className="h-6 w-4/5" />
                <Skeleton className="h-5 w-3/5" />
                <div className="flex justify-between items-center pt-2">
                  <Skeleton className="h-7 w-24" />
                  <Skeleton className="h-8 w-8 rounded-full" />
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}; 