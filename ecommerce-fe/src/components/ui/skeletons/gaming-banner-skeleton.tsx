import { Skeleton } from "@/components/ui/skeleton";

export const SubGamingBannerSkeleton = () => {
  return (
    <div className="w-[206px] h-[117px] rounded-xl bg-white relative border border-gray-100 shadow-sm">
      {/* Main image area */}
      <div className="h-full flex justify-center items-center overflow-hidden">
        <Skeleton className="h-[90px] w-[140px] object-contain" />
      </div>

      {/* Subtitle */}
      <div className="absolute bottom-0 left-1/2">
        <Skeleton className="h-2 w-[100px]" />
      </div>
    </div>
  );
};

export const GamingBannerSkeleton = () => {
  return (
    <div className="bg-white mt-[-5%] text-black h-[328px] w-[464px] rounded-2xl relative border border-gray-100 shadow-sm p-4">
      {/* Title */}
      <div className="absolute top-[5%] left-[5%]">
        <Skeleton className="h-7 w-[180px]" /> {/* Gaming Accessories */}
      </div>

      {/* Grid with 4 gaming items */}
      <div className="grid grid-cols-2 gap-12 mt-4 absolute top-[15%] left-[5%]">
        {Array(4).fill(0).map((_, index) => (
          <SubGamingBannerSkeleton key={index} />
        ))}
      </div>
    </div>
  );
}; 