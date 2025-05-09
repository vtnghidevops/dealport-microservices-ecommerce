import { Skeleton } from "@/components/ui/skeleton";

export const TestimonialCardSkeleton = ({ isGreen = false }: { isGreen?: boolean }) => {
  return (
    <div
      className={`testimonial-card h-[180px] w-[440px] !mx-4 my-3 rounded-xl p-[22px] shadow-sm border border-gray-300 ${isGreen ? 'bg-aqua-spring/30' : 'bg-white'
        }`}
    >
      <div className="flex items-center mb-4 p-2">
        {/* Avatar */}
        <Skeleton className="w-[3rem] h-[3rem] rounded-lg mr-12" />

        <div className='items-center flex w-[70%] gap-8'>
          {/* Username */}
          <Skeleton className="h-5 w-24" />

          {/* Rating stars */}
          <div className="flex space-x-1">
            {Array(5).fill(0).map((_, index) => (
              <Skeleton key={index} className="h-5 w-5" />
            ))}
          </div>
        </div>
      </div>

      {/* Review text */}
      <div className="space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-[90%]" />
        <Skeleton className="h-4 w-[80%]" />
      </div>
    </div>
  );
};

export const HappyCustomerSkeleton = () => {
  return (
    <section className="py-16 overflow-hidden bg-white mt-[4rem]">
      <div className="max-w-7xl mx-auto">
        {/* Title */}
        <div className="flex justify-center mb-4">
          <Skeleton className="h-10 w-[300px]" />
        </div>

        

        {/* Container for testimonials */}
        <div className="testimonials-container relative w-full">
          {/* Top row */}
          <div className="flex space-x-6 mb-8 overflow-hidden">
            {Array(3).fill(0).map((_, index) => (
              <TestimonialCardSkeleton
                key={`top-${index}`}
                isGreen={index % 2 === 1}
              />
            ))}
          </div>

          {/* Bottom row */}
          <div className="flex space-x-6 overflow-hidden">
            {Array(3).fill(0).map((_, index) => (
              <TestimonialCardSkeleton
                key={`bottom-${index}`}
                isGreen={index % 2 === 0}
              />
            ))}
          </div>
        </div>

        {/* Button */}
        <div className="text-center mt-[4rem] flex items-center justify-center">
          <Skeleton className="h-[4rem] w-[12rem] rounded-full" />
        </div>
      </div>
    </section>
  );
}; 