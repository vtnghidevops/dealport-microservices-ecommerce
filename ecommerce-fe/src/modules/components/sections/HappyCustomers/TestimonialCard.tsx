import React from 'react';
import { TestimonialItemProps } from './models/testimonial.model';

interface ExtendedTestimonialCardProps extends TestimonialItemProps {
  isGreen?: boolean;
}

const TestimonialCard: React.FC<ExtendedTestimonialCardProps> = ({ testimonial, isGreen = false }) => {
  const { name, avatarUrl, review, rating } = testimonial;
  const initial = name.charAt(0);
  
  return (
    <div 
      className={`testimonial-card h-[180px] w-[440px] mx-3 my-3 rounded-xl p-[22px] shadow-sm border border-gray-300 hover:shadow-md transition-shadow ${
        isGreen ? 'bg-aqua-spring' : 'bg-white'
      }`}
    >
      <div className="flex items-center mb-4 p-2">
        {avatarUrl ? (
          <img 
            src={avatarUrl} 
            alt={name} 
            className="w-[3rem] h-[3rem] rounded-lg mr-12 object-cover"
          />
        ) : (
          <div className="w-12 h-12 rounded-lg bg-gray-200 flex items-center justify-center mr-4">
            <span className="text-gray-600 font-bold">{initial}</span>
          </div>
        )}
        <div className='items-center flex w-[70%] gap-8'>
          <h3 className="title font-bold text-cyprus">{name}</h3>
          <div className="flex text-yellow-400 text-[24px] mb-1">
            {Array.from({ length: rating }).map((_, i) => (
              <span key={i}>★</span>
            ))}
          </div>
        </div>
      </div>
      <p className="text-gray-700">{review}</p>
    </div>
  );
};

export default TestimonialCard;
