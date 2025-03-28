import React from 'react';

interface UserActivityChartProps {
  data: {
    minute: string;
    count: number;
  }[];
}

const UserActivityChart: React.FC<UserActivityChartProps> = ({ data }) => {
  // Find the maximum value to calculate relative heights
  const maxValue = Math.max(...data.map(item => item.count));
  
  return (
    <div className="flex items-end space-x-0.5 w-[319px] h-[50px]">
      {data.map((item, index) => {
        // Calculate height as percentage of max value
        const heightPercentage = (item.count / maxValue) * 90;
        
        return (
          <div 
            key={index}
            className="flex-grow bg-ocean-green rounded-sm" 
            style={{ 
              width: "4px",
              height: `${heightPercentage}%`,
              minWidth: "3px"
            }}
            title={`${item.count} users in minute ${item.minute}`}
          />
        );
      })}
    </div>
  );
};

export default UserActivityChart;