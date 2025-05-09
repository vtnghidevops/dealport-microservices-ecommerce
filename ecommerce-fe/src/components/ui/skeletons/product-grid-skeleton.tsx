import { Skeleton } from "@/components/ui/skeleton";
import { ProductCardSkeleton } from './product-card-skeleton';

interface ProductGridSkeletonProps {
  columns?: number;
  rows?: number;
  withFilters?: boolean;
  withPagination?: boolean;
}

export const ProductGridSkeleton = ({
  columns = 4,
  rows = 2,
  withFilters = true,
  withPagination = true
}: ProductGridSkeletonProps) => {
  const count = columns * rows;

  return (
    <div className="w-full space-y-6">
      {/* Header with filters */}
      {withFilters && (
        <div className="flex flex-col md:flex-row justify-between space-y-4 md:space-y-0 mb-6">
          {/* Search bar */}
          <Skeleton className="h-10 w-full md:w-[300px] rounded-lg" />

          {/* Sort options */}
          <div className="flex items-center space-x-2">
            <Skeleton className="h-4 w-16" /> {/* Label */}
            <Skeleton className="h-10 w-[200px] rounded-md" /> {/* Select dropdown */}
          </div>
        </div>
      )}

      {/* Active filters */}
      {withFilters && (
        <div className="flex flex-wrap gap-2 items-center py-3 px-4 bg-gray-50 rounded-lg mb-6">
          <Skeleton className="h-4 w-24" /> {/* Label */}

          {/* Filter tags */}
          {Array(3).fill(0).map((_, index) => (
            <div key={index} className="flex items-center space-x-1 bg-white rounded-full px-3 py-1">
              <Skeleton className="h-4 w-20" />
              <Skeleton className="h-4 w-4 rounded-full" />
            </div>
          ))}
        </div>
      )}

      {/* Product grid */}
      <div className={`grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-${columns} gap-4`}>
        {Array(count).fill(0).map((_, index) => (
          <ProductCardSkeleton key={index} />
        ))}
      </div>

      {/* Pagination */}
      {withPagination && (
        <div className="flex justify-center items-center mt-8 space-x-2">
          <Skeleton className="h-9 w-9 rounded-md" /> {/* Previous */}
          {Array(5).fill(0).map((_, index) => (
            <Skeleton key={index} className="h-9 w-9 rounded-md" />
          ))}
          <Skeleton className="h-9 w-9 rounded-md" /> {/* Next */}
        </div>
      )}
    </div>
  );
}; 