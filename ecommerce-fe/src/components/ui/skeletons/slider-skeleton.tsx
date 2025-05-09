import { Skeleton } from "@/components/ui/skeleton";

export const SliderSkeleton = () => {
  return (
    <div className="relative w-full overflow-hidden">
      {/* Main slider image skeleton */}
      <div className="relative">
        <Skeleton className="w-full h-[400px] md:h-[500px] rounded-md" />
      </div>
    </div>
  );
};
