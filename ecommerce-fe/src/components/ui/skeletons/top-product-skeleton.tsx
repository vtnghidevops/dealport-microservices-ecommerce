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
      let height = "h-[200px]";
      let width = "w-full";

      // Tall column item (position 4)
      if (i === 3) {
        height = "h-[420px]";
      }

      // Wide row item (position 5)
      if (i === 4) {
        width = "w-[700px]";
      }

      items.push(
        <div key={i} className={`${gridPosition}`}>
          <div className="rounded-xl relative mr-[0.5rem] mt-[0.5rem] border border-gray-200 overflow-hidden">
            <Skeleton className={`${height} ${width} rounded-xl`} />

            {/* Price tag for regular items (not for tall column) */}
            {i !== 3 && (
              <div className="absolute bottom-2 right-5">
                <Skeleton className="h-5 w-16 rounded-md" />
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
        <Skeleton className="h-8 w-48" /> {/* Title */}
        <Skeleton className="w-[8rem] h-[3rem] rounded-3xl" /> {/* View All button */}
      </div>

      {/* Grid layout */}
      <div className="grid grid-cols-4 gap-4">
        {createSkeletonItems()}
      </div>
    </div>
  );
}; 