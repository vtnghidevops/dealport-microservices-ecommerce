import { Skeleton } from "@/components/ui/skeleton";

export const SliderSkeleton = () => {
  return (
    <div className="w-full mb-8 relative z-0">
      {/* Categories bar skeleton */}
      <div className="w-full bg-white border-b border-gray-200 shadow-sm relative z-10">
        <div className="max-w-7xl mx-auto">
          <div className="flex items-center gap-4 h-[48px] overflow-x-auto px-4 md:px-6">
            {Array(8).fill(0).map((_, index) => (
              <Skeleton key={index} className="h-6 w-20 flex-shrink-0" />
            ))}
            <Skeleton className="h-6 w-16 flex-shrink-0 ml-auto" />
          </div>
        </div>
      </div>

      {/* Banner skeleton */}
      <div className="w-full h-[400px] sm:h-[450px] md:h-[500px] relative border border-gray-200">
        <Skeleton className="w-full h-full rounded-none" />
      </div>
    </div>
  );
};
