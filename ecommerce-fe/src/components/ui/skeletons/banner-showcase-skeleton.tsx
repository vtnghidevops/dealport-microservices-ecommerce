import { Skeleton } from "@/components/ui/skeleton";

export const BannerShowCaseSkeleton = ({ count = 4 }: { count?: number }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5 mt-8">
      {Array(count).fill(0).map((_, index) => (
        <div key={index} className="relative rounded-lg overflow-hidden w-[22rem] h-[200px]">
          <Skeleton className="w-full h-full absolute" />
        </div>
      ))}
    </div>
  );
};


