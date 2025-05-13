import { Skeleton } from "@/components/ui/skeleton";
import { ProductCardSkeleton } from "./product-card-skeleton";

interface ProductListSkeletonProps {
  itemCount?: number;
  withFilters?: boolean;
  withPagination?: boolean;
}

export const ProductListSkeleton = ({
  itemCount = 12,
  withFilters = true,
  withPagination = true,
}: ProductListSkeletonProps) => {
  return (
    <div className="container mx-auto px-4 my-8 max-w-7xl">
      {/* Breadcrumb */}
      <div className="flex items-center mb-6">
        <Skeleton className="h-4 w-16 rounded-md" />
        <div className="mx-2">/</div>
        <Skeleton className="h-4 w-20 rounded-md" />
        <div className="mx-2">/</div>
        <Skeleton className="h-4 w-24 rounded-md" />
      </div>

      <div className="flex flex-col lg:flex-row gap-6">
        {/* Sidebar */}
        {withFilters && (
          <div className="w-full lg:w-64 flex-shrink-0">
            <div className="sticky top-24 space-y-6">
              {/* Category filter */}
              <div className="border border-aqua-spring/60 rounded-lg p-4 shadow-sm hover:shadow-md transition-all duration-300 hover:border-aqua-spring/80">
                <Skeleton className="h-6 w-32 mb-4 bg-aqua-spring" />
                {[1, 2, 3, 4, 5].map((i) => (
                  <div key={i} className="flex items-center mb-3">
                    <Skeleton className="h-4 w-4 rounded-sm mr-3" />
                    <Skeleton className="h-4 flex-1" />
                  </div>
                ))}
                <Skeleton className="h-8 w-full mt-4 rounded-md bg-aqua-spring/80 hover:bg-aqua-spring/90" />
              </div>

              {/* Price filter */}
              <div className="border border-aqua-spring/60 rounded-lg p-4 shadow-sm hover:shadow-md transition-all duration-300 hover:border-aqua-spring/80">
                <Skeleton className="h-6 w-24 mb-4 bg-aqua-spring" />
                <div className="space-y-4">
                  <div className="flex justify-between">
                    <Skeleton className="h-5 w-20 bg-aqua-spring/80" />
                    <Skeleton className="h-5 w-20 bg-aqua-spring/80" />
                  </div>
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-8 w-full rounded-md bg-aqua-spring/80 hover:bg-aqua-spring/90" />
                </div>
              </div>

              {/* Rating filter */}
              <div className="border border-aqua-spring/60 rounded-lg p-4 shadow-sm hover:shadow-md transition-all duration-300 hover:border-aqua-spring/80">
                <Skeleton className="h-6 w-28 mb-4 bg-aqua-spring" />
                {[5, 4, 3, 2, 1].map((rating) => (
                  <div key={rating} className="flex items-center mb-3">
                    <Skeleton className="h-4 w-4 rounded-sm mr-3" />
                    <div className="flex space-x-1 mr-2">
                      {Array.from({ length: 5 }).map((_, i) => (
                        <Skeleton key={i} className="h-3.5 w-3.5 rounded-sm" />
                      ))}
                    </div>
                    <Skeleton className="h-4 w-12" />
                  </div>
                ))}
                <Skeleton className="h-8 w-full mt-4 rounded-md bg-aqua-spring/80 hover:bg-aqua-spring/90" />
              </div>
            </div>
          </div>
        )}

        {/* Product Grid */}
        <div className="flex-1">
          {/* Sorting and view options */}
          <div className="flex justify-between items-center mb-6">
            <Skeleton className="h-5 w-48 bg-aqua-spring/80" />
            <div className="flex items-center space-x-3">
              <Skeleton className="h-9 w-36 rounded-md bg-aqua-spring/80" />
              <div className="flex space-x-2">
                <Skeleton className="h-9 w-9 rounded-md bg-aqua-spring/70" />
                <Skeleton className="h-9 w-9 rounded-md bg-aqua-spring/70" />
              </div>
            </div>
          </div>

          {/* Products grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {Array.from({ length: itemCount }).map((_, index) => (
              <div key={index} className="transform transition-all duration-300 hover:scale-[1.02]">
                <ProductCardSkeleton />
              </div>
            ))}
          </div>

          {/* Pagination */}
          {withPagination && (
            <div className="flex justify-center items-center mt-10">
              <div className="flex space-x-2">
                <Skeleton className="h-10 w-10 rounded-md bg-aqua-spring/80" />
                {[1, 2, 3, 4, 5].map((i) => (
                  <Skeleton key={i} className="h-10 w-10 rounded-md bg-aqua-spring/80" />
                ))}
                <Skeleton className="h-10 w-10 rounded-md bg-aqua-spring/80" />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}; 