import { Skeleton } from "@/components/ui/skeleton";

export const SidebarSkeleton = () => {
  return (
    <div className="w-full md:w-[260px] flex-shrink-0 border border-neutral-50 rounded-lg">
      <div className="bg-white p-4 rounded-lg shadow-sm">
        <div className="space-y-1">
          {/* Header */}
          <Skeleton className="h-6 w-3/4 mb-4" />

          {/* Navigation items */}
          {Array(8).fill(0).map((_, index) => (
            <div
              key={index}
              className="flex items-center px-3 py-3"
            >
              <Skeleton className="h-5 w-5 mr-3" /> {/* Icon */}
              <Skeleton className="h-4 w-24" /> {/* Label */}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

// Filter Sidebar Skeleton
export const FilterSidebarSkeleton = () => {
  return (
    <div className="w-full md:w-[300px] flex-shrink-0 border border-neutral-50 rounded-lg">
      <div className="bg-white p-4 rounded-lg shadow-sm space-y-6">
        {/* Categories section */}
        <div className="space-y-3">
          <Skeleton className="h-6 w-1/2" /> {/* Section title */}
          <div className="space-y-2 pl-2">
            {Array(5).fill(0).map((_, index) => (
              <div key={index} className="flex items-center space-x-2">
                <Skeleton className="h-4 w-4 rounded-sm" /> {/* Checkbox */}
                <Skeleton className="h-4 w-20" /> {/* Category name */}
                <Skeleton className="h-3 w-6 ml-auto" /> {/* Count */}
              </div>
            ))}
          </div>
        </div>

        {/* Price range section */}
        <div className="space-y-3">
          <Skeleton className="h-6 w-1/2" /> {/* Section title */}
          <Skeleton className="h-12 w-full" /> {/* Price slider */}
          <div className="flex justify-between">
            <Skeleton className="h-8 w-[45%]" /> {/* Min price input */}
            <Skeleton className="h-8 w-[45%]" /> {/* Max price input */}
          </div>
        </div>

        {/* Rating filter */}
        <div className="space-y-3">
          <Skeleton className="h-6 w-1/2" /> {/* Section title */}
          <div className="space-y-2">
            {Array(5).fill(0).map((_, index) => (
              <div key={index} className="flex items-center">
                <div className="flex">
                  {Array(5).fill(0).map((_, starIndex) => (
                    <Skeleton key={starIndex} className="h-4 w-4 mr-1" />
                  ))}
                </div>
                <Skeleton className="h-3 w-10 ml-2" /> {/* Count */}
              </div>
            ))}
          </div>
        </div>

        {/* Brand filter */}
        <div className="space-y-3">
          <Skeleton className="h-6 w-1/2" /> {/* Section title */}
          <div className="space-y-2 pl-2">
            {Array(4).fill(0).map((_, index) => (
              <div key={index} className="flex items-center space-x-2">
                <Skeleton className="h-4 w-4 rounded-sm" /> {/* Checkbox */}
                <Skeleton className="h-4 w-24" /> {/* Brand name */}
              </div>
            ))}
          </div>
        </div>

        {/* Apply filters button */}
        <Skeleton className="h-10 w-full" />
      </div>
    </div>
  );
}; 