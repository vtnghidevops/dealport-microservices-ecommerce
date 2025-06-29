import { Skeleton } from "@/components/ui/skeleton";

export const DisplayCardMainSkeleton = () => {
  return (
    <div className="w-[392px] h-[165px] rounded-lg overflow-hidden shadow-sm relative">
      <Skeleton className="w-full h-full" />

     
    </div>
  );
};

export const DisplayCardSkeleton = () => {
  return (
    <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-sm relative">
      <Skeleton className="w-full h-full" />

      {/* More details text */}
      <div className="absolute bottom-2 left-5">
        <Skeleton className="h-3 w-[4rem]" />
      </div>
    </div>
  );
};

export const DisplayGridSkeleton = () => {
  return (
    <div className="h-[328px] w-[392px] mt-[-5%] text-black rounded-2xl relative">
      <div className="flex items-center flex-wrap gap-8">
        <div className="space-y-4 flex items-center gap-5">
          <DisplayCardSkeleton />
          <DisplayCardSkeleton />
        </div>
        <div>
          <DisplayCardMainSkeleton />
        </div>
        
      </div>
    </div>
  );
}; 