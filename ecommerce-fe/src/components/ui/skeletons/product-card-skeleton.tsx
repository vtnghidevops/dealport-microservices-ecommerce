import { Skeleton } from "@/components/ui/skeleton";

export const ProductCardSkeleton = () => {
  return (
    <div className="rounded-lg border shadow-sm p-3 space-y-3 min-w-[16rem] max-w-[16rem] mr-5">
      {/* Image skeleton */}
      <Skeleton className="h-[200px] w-full rounded-md" />

      {/* Title skeleton */}
      <Skeleton className="h-5 w-full" />

      {/* Price skeleton */}
      <Skeleton className="h-4 w-1/4" />

      {/* Description skeleton */}
      <div className="space-y-2">
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-5/6" />
      </div>

      {/* Button skeleton */}
      <Skeleton className="h-9 w-full mt-2" />
    </div>
  );
};

export const ProductCardSkeletonGrid = ({ count = 5 }: { count?: number }) => {
  return (
    <div className="flex items-center overflow-hidden">
      {Array(count).fill(0).map((_, index) => (
        <ProductCardSkeleton key={index} />
      ))}
    </div>
  );
}; 