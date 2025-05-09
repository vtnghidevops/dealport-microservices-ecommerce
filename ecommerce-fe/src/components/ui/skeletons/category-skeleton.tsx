import { Skeleton } from "@/components/ui/skeleton";

interface CategorySkeletonProps {
  count?: number;
  variant?: "homepage" | "admin";
}

export const CategorySkeleton = ({ count = 5, variant = "homepage" }: CategorySkeletonProps) => {
  // Homepage category skeleton (tall card with image at top and name at bottom)
  if (variant === "homepage") {
    return (
      <div className="flex overflow-x-auto gap-4 pb-4 -mx-4 px-4 scrollbar-hide">
        {Array(count).fill(0).map((_, index) => (
          <div key={index} className="flex-shrink-0 h-[220px] w-[200px] rounded-xl mr-5">
            <div className="relative rounded-lg overflow-hidden shadow-sm border border-gray-200 h-full">
              {/* Category image */}
              <div className="w-full h-[180px] flex items-center justify-center overflow-hidden">
                <Skeleton className="w-[148px] h-[140px] rounded-md" />
              </div>
              {/* Category name */}
              <div className="p-2 text-center absolute bottom-0 flex justify-center w-full bg-white">
                <Skeleton className="h-4 w-24" />
              </div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  // Admin category skeleton (horizontal card with image on left and name on right)
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-4 my-8 mb-5">
      {Array(count).fill(0).map((_, index) => (
        <div key={index} className="w-full h-[88px] flex items-center gap-3 p-3 bg-white rounded-lg shadow-sm border">
          {/* Category image */}
          <div className="flex items-center justify-center border border-neutral-200 rounded-lg">
            <Skeleton className="w-[64px] h-[64px] rounded-lg" />
          </div>
          {/* Category name */}
          <Skeleton className="h-5 w-24" />
        </div>
      ))}
    </div>
  );
}; 