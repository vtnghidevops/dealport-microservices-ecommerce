import React from 'react';
import { OrderSummaryCardProps } from './models/card.model';
import { IoIosArrowRoundDown } from "react-icons/io";
import { IoIosArrowRoundUp } from "react-icons/io";
export const OrderSummaryCard: React.FC<OrderSummaryCardProps> = ({ 
  title, 
  value, 
  growthRate, 
  period, 
}) => {
  const getGrowthRateColor = () => {
    if (title === 'Canceled Orders') {
      return growthRate > 0 ? 'text-error' : 'text-success';
    }
    return growthRate > 0 ? 'text-success' : 'text-error';
  };

  const getGrowthRateIcon = () => {
    if (title === 'Canceled Orders') {
      return growthRate > 0 ? <IoIosArrowRoundUp className='h-[16px] w-[16px]'/> : <IoIosArrowRoundDown className='h-[16px] w-[16px]'/>;
    }
    return growthRate > 0 ? <IoIosArrowRoundUp className='h-[16px] w-[16px]'/> : <IoIosArrowRoundDown className='h-[16px] w-[16px]'/>;
  };

  

  return (
    <div className="rounded-lg p-5 w-[270px] h-[135px] bg-white filter drop-shadow-lg">
      <div className="flex items-center justify-between mb-2">
        <h3 className="button-text text-cyprus text-[18px] ">{title}</h3>
        <button className="text-gray-700 text-[18px]">
          <svg
            stroke="currentColor"
            fill="none"
            stroke-width="2"
            viewBox="0 0 24 24"
            aria-hidden="true"
            height="1em"
            width="1em"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
            ></path>
          </svg>
        </button>
      </div>
      <div className="flex justify-between items-end">
        <div>
          <div className="flex items-center gap-8 w-[160px] h-[34px]">
            <span className="header-2 font-bold text-cyprus">
              {value.toLocaleString()}
            </span>
            <div className={`mt-[1rem] text-sm font-medium flex items-center ${getGrowthRateColor()}`}>
              <div className='mr-1 text-[12px]'>{getGrowthRateIcon()} </div>
              <div >{Math.abs(growthRate)}%</div>
            </div>
          </div>

          <p className="text-xs pt-[0.5rem] text-neutral-500">{period}</p>
        </div>
      </div>
    </div>
  );
};