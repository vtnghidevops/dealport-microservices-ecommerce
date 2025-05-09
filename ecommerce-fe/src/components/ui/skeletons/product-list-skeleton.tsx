import { Skeleton } from "@/components/ui/skeleton";

export const ProductListSkeleton = () => {
  return (
    <div className="container mx-auto px-4 md:px-8 py-4">
      {/* Breadcrumb skeleton */}
      <div className="flex items-center mb-6">
        <div className="flex gap-2 items-center">
          <Skeleton className="h-4 w-14 rounded" />
          <span className="mx-2">/</span>
          <Skeleton className="h-4 w-14 rounded" />
          <span className="mx-2">/</span>
          <Skeleton className="h-4 w-24 rounded" />
        </div>
      </div>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Sidebar skeleton */}
        <div className="w-full lg:w-1/4 mb-5">
          <div className="bg-white rounded-lg shadow-md p-5 mb-6 min-h-[600px] md:min-h-[800px] sticky top-4">
            <Skeleton className="h-7 w-36 mb-6" />

            {/* Category filter */}
            <div className="space-y-3 mb-8">
              <Skeleton className="h-5 w-full" />
              <Skeleton className="h-5 w-full" />
              <Skeleton className="h-5 w-full" />
              <Skeleton className="h-5 w-4/5" />
              <Skeleton className="h-5 w-3/5" />
            </div>

            {/* Price filter */}
            <Skeleton className="h-7 w-28 mb-4" />
            <div className="space-y-4 mb-8">
              <div className="flex justify-between">
                <Skeleton className="h-6 w-20" />
                <Skeleton className="h-6 w-20" />
              </div>
              <Skeleton className="h-8 w-full rounded-lg" />
            </div>

            {/* Rating filter */}
            <Skeleton className="h-7 w-28 mb-4" />
            <div className="space-y-3 mb-8">
              {Array(4).fill(0).map((_, i) => (
                <div key={i} className="flex items-center gap-3">
                  <Skeleton className="h-5 w-5 rounded-sm" />
                  <div className="flex gap-1">
                    {Array(5 - i).fill(0).map((_, j) => (
                      <Skeleton key={j} className="h-5 w-5 rounded-full" />
                    ))}
                  </div>
                </div>
              ))}
            </div>

            {/* Brand filter */}
            <Skeleton className="h-7 w-28 mb-4" />
            <div className="space-y-3 mb-8">
              {Array(5).fill(0).map((_, i) => (
                <div key={i} className="flex items-center gap-3">
                  <Skeleton className="h-5 w-5 rounded-sm" />
                  <Skeleton className="h-5 w-36" />
                </div>
              ))}
            </div>

            {/* Tags filter */}
            <Skeleton className="h-7 w-28 mb-4" />
            <div className="flex flex-wrap gap-2">
              {Array(6).fill(0).map((_, i) => (
                <Skeleton key={i} className="h-8 w-20 rounded-full mb-2" />
              ))}
            </div>
          </div>
        </div>

        {/* Main content skeleton */}
        <div className="w-full lg:w-3/4">
          {/* Header with search and sort */}
          <div className="bg-white p-5 rounded-lg shadow-md mb-6">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-5">
              <div className="w-full md:w-[400px] h-[44px] mb-4 md:mb-0">
                <Skeleton className="h-full w-full rounded-lg" />
              </div>

              <div className="flex items-center space-x-3">
                <Skeleton className="h-5 w-24" />
                <Skeleton className="h-[44px] w-[160px] rounded-lg" />
              </div>
            </div>

            {/* Active filters */}
            <div className="flex flex-wrap gap-2 py-3">
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-8 w-28 rounded-full" />
              <Skeleton className="h-8 w-36 rounded-full" />
              <Skeleton className="h-8 w-32 rounded-full" />
            </div>
          </div>

          {/* Product grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 mb-5">
            {Array(9).fill(0).map((_, index) => (
              <div key={index} className="bg-white rounded-lg shadow-md overflow-hidden mb-5">
                <Skeleton className="h-[200px] w-full" />
                <div className="p-4 space-y-3">
                  <Skeleton className="h-6 w-4/5" />
                  <Skeleton className="h-5 w-3/5" />
                  <div className="flex justify-between items-center pt-2">
                    <Skeleton className="h-7 w-24" />
                    <Skeleton className="h-9 w-9 rounded-full" />
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Pagination */}
          <div className="flex justify-center mt-8 mb-10">
            <div className="flex gap-2">
              <Skeleton className="h-10 w-10 rounded-md" />
              {Array(3).fill(0).map((_, i) => (
                <Skeleton key={i} className="h-10 w-10 rounded-md mx-1" />
              ))}
              <Skeleton className="h-10 w-10 rounded-md" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}; 