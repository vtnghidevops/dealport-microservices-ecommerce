import { Skeleton } from "@/components/ui/skeleton";

interface TopProductSkeletonProps {
  count?: number;
}

export const TopProductSkeleton = ({ count = 6 }: TopProductSkeletonProps) => {
  // Ensure count does not exceed 6
  const itemCount = Math.min(count, 6);

  // Define grid positions that match the real component's layout
  const defaultGridPositions = [
    "col-span-1 row-span-1", // Item 1
    "col-span-1 row-span-1", // Item 2
    "col-span-1 row-span-1", // Item 3
    "col-span-1 row-span-2", // Item 4 (tall column)
    "col-span-2 row-span-1 row-start-2", // Item 5 (wide row)
    "col-span-1 row-span-1 row-start-2", // Item 6
  ];

  const createSkeletonItems = () => {
    const items = [];

    for (let i = 0; i < itemCount; i++) {
      const gridPosition = defaultGridPositions[i] || "col-span-1 row-span-1";

      // Determine height based on grid position
      let height = "h-[220px]";
      let width = "w-full";

      // Tall column item (position 4)
      if (i === 3) {
        height = "h-[450px]";
      }

      // Wide row item (position 5)
      if (i === 4) {
        width = "w-[700px]";
      }

      items.push(
        <div key={i} className={`${gridPosition} transform transition-all duration-300 hover:scale-[1.02]`}>
          <div className="rounded-xl relative mr-[0.5rem] mt-[0.5rem] border border-aqua-spring/60 overflow-hidden bg-white p-[10px] transition-all duration-300 hover:shadow-md hover:border-aqua-spring/80">
            {/* Main image skeleton */}
            <Skeleton className={`${height} ${width} rounded-xl bg-aqua-spring/80`} />

            {/* Product info for trending product cards */}
            {(i < 3 || i > 3) && (
              <div className="mt-4 flex flex-col gap-2.5">
                <Skeleton className="h-6 w-3/4" /> {/* Title */}
                <Skeleton className="h-4 w-5/6" /> {/* Description */}

                <div className="flex items-center mt-1.5">
                  <div className="flex space-x-1">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <Skeleton key={star} className="h-3.5 w-3.5 rounded-sm" />
                    ))}
                  </div>
                  <Skeleton className="h-3 w-16 ml-2" /> {/* Review count */}
                </div>

                <div className="flex items-center mt-1">
                  <Skeleton className="h-5 w-16 bg-aqua-spring" /> {/* Price */}
                  <Skeleton className="h-4 w-12 ml-2 bg-aqua-spring/70" /> {/* Original price */}
                </div>

                <div className="flex justify-between items-center mt-2">
                  <Skeleton className="h-4 w-20" /> {/* View details */}
                  <Skeleton className="h-9 w-28 rounded-md bg-aqua-spring hover:bg-aqua-spring/90" /> {/* Add to cart */}
                </div>
              </div>
            )}

            {/* Custom tall column for Men's Collection (position 3) */}
            {i === 3 && (
              <div className="absolute bottom-4 left-0 right-0 px-4">
                <Skeleton className="h-7 w-48 mb-3 mx-auto rounded-lg bg-aqua-spring" /> {/* Collection title */}
                <div className="grid grid-cols-2 gap-3">
                  <Skeleton className="h-16 w-full rounded-lg hover:bg-aqua-spring/90" /> {/* Product 1 */}
                  <Skeleton className="h-16 w-full rounded-lg hover:bg-aqua-spring/90" /> {/* Product 2 */}
                  <Skeleton className="h-16 w-full rounded-lg hover:bg-aqua-spring/90" /> {/* Product 3 */}
                  <Skeleton className="h-16 w-full rounded-lg hover:bg-aqua-spring/90" /> {/* Product 4 */}
                </div>
              </div>
            )}
          </div>
        </div>
      );
    }

    return items;
  };

  return (
    <div className="w-full py-6 md:px-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-[2rem]">
        <Skeleton className="h-8 w-48 bg-aqua-spring/90" /> {/* Title */}
        <Skeleton className="w-[8rem] h-[3rem] rounded-3xl bg-aqua-spring" /> {/* View All button */}
      </div>

      {/* Grid layout with animation */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 transition-all duration-300">
        {createSkeletonItems()}
      </div>
    </div>
  );
}; 