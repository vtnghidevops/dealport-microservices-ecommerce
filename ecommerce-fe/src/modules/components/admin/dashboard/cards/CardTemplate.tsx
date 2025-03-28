import React from 'react';
import { HiOutlineDotsVertical } from "react-icons/hi";
import { CardBase } from './models/card.model';

// Define possible card types with their specific data structures
interface SalesCardData extends CardBase {
  type: 'sales';
  amount: number;
  currency: string;
  percentChange: number;
  previousAmount: number;
}

interface OrdersCardData extends CardBase {
  type: 'orders';
  count: number;
  percentChange: number;
  previousCount: number;
}

interface PendingCanceledCardData extends CardBase {
  type: 'pendingCanceled';
  pending: {
    count: number;
    userCount: number;
  };
  canceled: {
    count: number;
    percentChange: number;
  };
}

// Union type for all possible card data types
type DashboardCard = SalesCardData | OrdersCardData | PendingCanceledCardData;

const DashboardCard: React.FC<DashboardCard> = (props) => {
  const { title, lastDays } = props;
  
  // Determine if change is positive (for displaying arrow up/down)
  const getChangeDirection = (change: number) => {
    return change >= 0;
  };

  // Format number to currency or compact form
  const formatNumber = (value: number, style: 'currency' | 'decimal' = 'decimal', currencySymbol = '$') => {
    // Convert currency symbol to ISO code
    const currencyCodeMap: Record<string, string> = {
      '$': 'USD',
      '€': 'EUR',
      '£': 'GBP',
      '¥': 'JPY',
      '₹': 'INR',
    };
    
    // Use USD as default if the symbol is not recognized
    const currencyCode = currencyCodeMap[currencySymbol] || 'USD';
    
    return new Intl.NumberFormat('en-US', {
      style: style === 'currency' ? 'currency' : 'decimal',
      currency: style === 'currency' ? currencyCode : undefined,
      notation: 'compact',
      maximumFractionDigits: 0
    }).format(value);
  };

  // Render card content based on type
  const renderCardContent = () => {
    switch(props.type) {
      case 'sales':
        const isPositiveSales = getChangeDirection(props.percentChange);
        return (
          <>
            <div className="flex items-center gap-2 mb-2 mt-[1.25rem] h-[38px]">
              <h2 className="text-2xl font-bold">{formatNumber(props.amount, 'currency', props.currency)}</h2>
              <span className={`text-xs px-2 py-1 ${isPositiveSales ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'} rounded-full flex items-center`}>
                {renderChangeIcon(isPositiveSales)}
                Sales {Math.abs(props.percentChange)}%
              </span>
            </div>
            <p className="text-xs text-gray-500">
              Previous {lastDays} days: 
              <span className="text-blue-500 ml-[0.3rem]">
                {`${formatNumber(props.previousAmount, 'currency', props.currency)}`}
              </span>
            </p>
          </>
        );
      
      case 'orders':
        const isPositiveOrders = getChangeDirection(props.percentChange);
        return (
          <>
            <div className="flex items-center gap-2 mb-2 mt-[1.25rem] h-[38px]">
              <h2 className="text-2xl font-bold">
                {formatNumber(props.count)}
              </h2>
              <span
                className={`text-xs px-2 py-1 ${
                  isPositiveOrders
                    ? "bg-green-100 text-green-600"
                    : "bg-red-100 text-red-600"
                } rounded-full flex items-center`}
              >
                {renderChangeIcon(isPositiveOrders)}
                Order {Math.abs(props.percentChange)}%
              </span>
            </div>
            <p className="text-xs text-gray-500">
              Previous {lastDays} days:
              <span className="text-blue-500 ml-[0.3rem]">
                {`${formatNumber(props.previousCount)}`}
              </span>
            </p>
          </>
        );
      
      case 'pendingCanceled':
        const isPositiveCanceled = getChangeDirection(props.canceled.percentChange);
        return (
          <>
            <div className="mb-2 mt-[1.25rem] h-[38px] flex items-center justify-start">
              <div className="flex justify-between mb-1 flex-col w-[139px] h-[50px] mt-[1.25rem]">
                <span className="text-[14px]">Pending</span>
                <div className="text-[22px] font-bold flex items-center">
                  {props.pending.count}
                  <span className='text-[16px] font-medium text-ocean-green ml-[0.5rem]'>{`User ${props.pending.userCount}`}</span>
                </div>

              </div>
              <span className='h-[50px] flex items-center text-gray-400'>|</span>
              <div className="flex justify-start items-center flex-col w-[139px] h-[50px] mt-[1.25rem]">
                <span className="pl-[2rem] text-[14px] block w-full pd-[2rem]">Canceled</span>
                <div className="flex items-center w-full pl-[2rem]">
                  <span className="text-[22px] font-bold mr-1 text-red-500">{props.canceled.count}</span>
                  <span className={`text-xs px-2 py-1 ${isPositiveCanceled ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'} rounded-full flex items-center`}>
                    {renderChangeIcon(!isPositiveCanceled)}
                    {Math.abs(props.canceled.percentChange)}%
                  </span>
                </div>
              </div>
             
            </div>
          </>
        );
    }
  };
  
  // Helper to render change icon
  const renderChangeIcon = (isPositive: boolean) => {
    if (isPositive) {
      return (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 mr-1" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clipRule="evenodd" />
        </svg>
      );
    } else {
      return (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 mr-1" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
        </svg>
      );
    }
  };

  return (
    <div className="bg-white rounded-lg p-4 filter drop-shadow-lg h-[222px] w-[362px]">
      <div className='p-5 h-full relative'>
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-[18px] font-bold">{title}</h3>
          <button className="text-gray-700">
            <HiOutlineDotsVertical></HiOutlineDotsVertical>
          </button>
        </div>
        <p className="text-[14px] text-gray-500 mb-1">Last {lastDays} days</p>
        {renderCardContent()}
        <button className="flex justify-center items-center absolute right-[5%] bottom-[5%] mt-4 px-4 py-2 text-sm h-[32px] w-[92px] rounded-[25px] border border-primary text-primary">
          Details
        </button>
      </div>
    </div>
  );
};

export default DashboardCard;