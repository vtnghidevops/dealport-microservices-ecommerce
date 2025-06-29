import React from 'react';
import { OrderSummaryCardProps } from './models/card.model';
import { IoIosArrowRoundDown } from "react-icons/io";
import { IoIosArrowRoundUp } from "react-icons/io";

export const OrderSummaryCard: React.FC<OrderSummaryCardProps> = ({
  title = 'Orders',
  value = 0,
  growthRate = 0,
  period = 'Last 7 days',
}) => {
  // Ensure we have valid numbers
  const safeValue = typeof value === 'number' && !isNaN(value) ? value : 0;
  const safeGrowthRate = typeof growthRate === 'number' && !isNaN(growthRate) ? growthRate : 0;

  const getGrowthRateColor = () => {
    if (title === 'Canceled Orders') {
      return safeGrowthRate > 0 ? 'text-error' : 'text-success';
    }
    return safeGrowthRate > 0 ? 'text-success' : 'text-error';
  };

  const getGrowthRateIcon = () => {
    if (title === 'Canceled Orders') {
      return safeGrowthRate > 0 ? <IoIosArrowRoundUp className='h-[16px] w-[16px]' /> : <IoIosArrowRoundDown className='h-[16px] w-[16px]' />;
    }
    return safeGrowthRate > 0 ? <IoIosArrowRoundUp className='h-[16px] w-[16px]' /> : <IoIosArrowRoundDown className='h-[16px] w-[16px]' />;
  };

  return (
    <div className="rounded-lg p-5 w-[270px] h-[135px] bg-white filter drop-shadow-lg">
      <div className="flex items-center justify-between mb-2">
        <h3 className="button-text text-cyprus text-[18px] ">{title}</h3>
        <button className="text-gray-700 text-[18px]">
          <svg
            stroke="currentColor"
            fill="none"
            strokeWidth="2"
            viewBox="0 0 24 24"
            aria-hidden="true"
            height="1em"
            width="1em"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
            ></path>
          </svg>
        </button>
      </div>
      <div className="flex justify-between items-end">
        <div>
          <div className="flex items-center gap-8 w-[160px] h-[34px]">
            <span className="header-2 font-bold text-cyprus">
              {safeValue.toLocaleString()}
            </span>
            <div className={`mt-[1rem] text-sm font-medium flex items-center ${getGrowthRateColor()}`}>
              <div className='mr-1 text-[12px]'>{getGrowthRateIcon()} </div>
              <div >{Math.abs(safeGrowthRate)}%</div>
            </div>
          </div>

          <p className="text-xs pt-[0.5rem] text-neutral-500">{period || 'Last 7 days'}</p>
        </div>
      </div>
    </div>
  );
};