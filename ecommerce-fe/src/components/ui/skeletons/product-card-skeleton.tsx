import { Skeleton } from "@/components/ui/skeleton";

interface ProductCardSkeletonProps {
  count?: number;
}

export const ProductCardSkeleton = ({ count = 1 }: ProductCardSkeletonProps) => {
  const skeletons = Array.from({ length: count }).map((_, index) => (
    <div key={index} className="mx-2 p-[10px] rounded-xl max-w-[272px] min-w-[272px] border border-aqua-spring/70 shadow-sm hover:shadow-md transition-all duration-300 transform hover:scale-[1.02]">
      {/* Image container with hover effect */}
      <div className="relative overflow-hidden rounded-xl">
        <Skeleton className="min-w-[248px] min-h-[180px] max-w-[248px] max-h-[180px] rounded-xl bg-aqua-spring/90" />
        {/* Wishlist button */}
        <div className="absolute top-2 right-[5%]">
          <Skeleton className="h-[1.5rem] w-[1.5rem] rounded-full" />
        </div>
        {/* Sale badge */}
        <div className="absolute top-2 left-[5%]">
          <Skeleton className="h-[1.2rem] w-[2.5rem] rounded-md" />
        </div>
      </div>

      {/* Product details */}
      <div className="mt-4 space-y-2.5">
        <Skeleton className="h-6 w-3/4" /> {/* Title */}
        <Skeleton className="h-4 w-full" /> {/* Description */}

        <div className="flex items-center mt-1.5">
          <div className="flex space-x-0.5">
            {[1, 2, 3, 4, 5].map((star) => (
              <Skeleton key={star} className="h-3.5 w-3.5 rounded-sm" />
            ))}
          </div>
          <Skeleton className="h-3 w-16 ml-2" /> {/* Review count */}
        </div>

        <div className="flex items-center">
          <Skeleton className="h-5 w-16 bg-aqua-spring" /> {/* Price */}
          <Skeleton className="h-4 w-16 ml-2 bg-aqua-spring/70" /> {/* Original price */}
          <Skeleton className="h-4 w-12 ml-2 rounded-sm bg-aqua-spring/90" /> {/* Discount */}
        </div>

        <div className="flex justify-between items-center mt-[15px]">
          <Skeleton className="h-4 w-24" /> {/* View Details */}
          <Skeleton className="h-[39px] w-[120px] rounded-md bg-aqua-spring hover:bg-aqua-spring/90" /> {/* Add to Cart button */}
        </div>
      </div>
    </div>
  ));

  return <>{skeletons}</>;
};

export const ProductCardSkeletonGrid = ({ count = 5 }: { count?: number }) => {
  return (
    <div className="flex items-center gap-4 overflow-hidden">
      {Array(count).fill(0).map((_, index) => (
        <ProductCardSkeleton key={index} />
      ))}
    </div>
  );
}; 