import React, { useEffect, useState } from "react";
import TestimonialCard from "./TestimonialCard";
import { TestimonialService } from "./services/testimonial.service";
import { TestimonialItem } from "./models/testimonial.model";
import { Button } from "../../common/Button";
import "../../../app.css";

const HappyCustomers: React.FC = () => {
  const [testimonials, setTestimonials] = useState<TestimonialItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchTestimonials = async () => {
      try {
        const data = await TestimonialService.getTestimonials();
        setTestimonials(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching testimonials:", error);
        setLoading(false);
      }
    };

    fetchTestimonials();
  }, []);

  // Chia testimonials thành 2 hàng
  const topRow = testimonials.slice(0, 3);
  const bottomRow = testimonials.slice(3);

  if (loading) {
    return <div className="text-center py-16">Loading testimonials...</div>;
  }

  return (
    <section className="py-16 overflow-hidden bg-white mt-[4rem]">
      <div className="max-w-7xl mx-auto">
        <h2 className="text-3xl font-semibold text-center text-ocean-green mb-4 px-4">
          Our Happy Customers
        </h2>
        <p className="text-center max-w-3xl mx-auto text-gray-700 mb-12 px-4">
          Don't just take our word for it – see how our products and services
          have delighted customers across the globe, one experience at a time.
        </p>

        {/* Container cho toàn bộ phần testimonials - full width để hiện thị overflowing */}
        <div className="testimonials-container relative w-full">
          {/* Dòng trên di chuyển từ phải sang trái - mở rộng full width */}
          <div className="scroll-container scroll-left overflow-visible">
            <div className="scroll-content">
              {[...topRow, ...topRow, ...topRow].map((testimonial, index) => (
                <div
                  key={`top-${testimonial.id}-${index}`}
                  className="testimonial-wrapper"
                >
                  <TestimonialCard
                    testimonial={testimonial}
                    isGreen={index % 2 === 1}
                  />
                </div>
              ))}
            </div>
          </div>

          {/* Khoảng cách giữa 2 hàng */}
          <div className="h-8"></div>

          {/* Dòng dưới di chuyển từ trái sang phải - mở rộng full width */}
          <div className="scroll-container scroll-right overflow-visible">
            <div className="scroll-content">
              {[...bottomRow, ...bottomRow, ...bottomRow].map(
                (testimonial, index) => (
                  <div
                    key={`bottom-${testimonial.id}-${index}`}
                    className="testimonial-wrapper"
                  >
                    <TestimonialCard
                      testimonial={testimonial}
                      isGreen={index % 2 === 0}
                    />
                  </div>
                )
              )}
            </div>
          </div>
        </div>

        <div className="text-center mt-[4rem] flex items-center justify-center ">
          <Button type="primary-cy" text="GET STARTED"></Button>
        </div>
      </div>
    </section>
  );
};

export default HappyCustomers;
