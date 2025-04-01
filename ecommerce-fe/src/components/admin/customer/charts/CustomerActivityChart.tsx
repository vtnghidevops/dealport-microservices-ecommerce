import React from 'react';
import { CustomerChartData } from '../models/customer.model';
import { CustomerOverview } from '../models/customer.model';

interface CustomerActivityChartProps {
  chartData: CustomerChartData[];
  overview?: CustomerOverview;
}

const CustomerActivityChart: React.FC<CustomerActivityChartProps> = ({ 
  chartData,
  overview
}) => {
  // Find the peak day (highest value)
  const peakDay = chartData.reduce((peak, current) => 
    current.count > peak.count ? current : peak, chartData[0]);
  
  // Find the index of Wednesday for the vertical line
  const wednesdayIndex = chartData.findIndex(data => data.day === 'Wed');
  const wednesdayPosition = ((wednesdayIndex) / (chartData.length - 1)) * 100;

  // Find max value for scaling (round up to nearest 10k for nice y-axis labels)
  const maxValue = Math.ceil(Math.max(...chartData.map(item => item.count)) / 10000) * 10000;
  
  return (
    <div className="bg-white rounded-lg shadow p-5 w-[826px] h-[445px]">
      {/* Header */}
      <div className="mb-4 flex justify-between items-center h-[38px]">
        <h3 className="text-[18px] font-bold text-cyprus">Customer Overview</h3>
        <div className="flex items-center gap-2 w-[200px] h-full">
          <div className='rounded-lg p-4 gap-4 w-[178px] h-full bg-aqua-spring flex justify-center items-center'>
            <button className="h-[30px] w-[83px] text-[12px] text-ocean-green bg-white rounded-lg">
              This week
            </button>
            <button className="h-[30px] w-[83px] text-[12px] text-neutral-500 rounded-lg">
              Last week
            </button>
          </div>
          <button className="p-2 text-gray-400">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
            </svg>
          </button>
        </div>
      </div>

      {/* Customer Metrics */}
      <div className="grid grid-cols-4 gap-4 mt-5 shadow-sm">
        <div className="border-b-2 border-success pb-2 w-[182px] h-[78px]">
          <div className="text-[24px] font-bold text-cyprus">
            {overview ? formatNumber(overview.activeCustomers.count) : "25k"}
          </div>
          <div className="text-sm text-gray-500">
            {overview
              ? overview.activeCustomers.chartLabel
              : "Active Customers"}
          </div>
        </div>
        <div className="border-b-2 border-gray-200 pb-2 w-[182px] h-[78px]">
          <div className="text-[24px] font-bold text-cyprus">
            {overview ? formatNumber(overview.repeatCustomers.count) : "5.6k"}
          </div>
          <div className="text-sm text-gray-500">
            {overview
              ? overview.repeatCustomers.chartLabel
              : "Repeat Customers"}
          </div>
        </div>
        <div className="border-b-2 border-gray-200 pb-2 w-[182px] h-[78px]">
          <div className="text-[24px] font-bold text-cyprus">
            {overview ? formatNumber(overview.shopVisitor.count) : "250k"}
          </div>
          <div className="text-sm text-gray-500">
            {overview ? overview.shopVisitor.chartLabel : "Shop Visitor"}
          </div>
        </div>
        <div className="border-b-2 border-gray-200 pb-2 w-[182px] h-[78px]">
          <div className="text-3xl font-medium text-gray-900">
            {overview ? `${overview.conversionRate.rate}%` : "5.5%"}
          </div>
          <div className="text-sm text-gray-500">
            {overview ? overview.conversionRate.chartLabel : "Conversion Rate"}
          </div>
        </div>
      </div>

      {/* Chart */}
      <div className="mt-8 relative" style={{ height: "220px" }}>
        {/* Y-axis labels */}
        <div className="absolute left-0 inset-y-0 flex flex-col justify-between text-xs text-gray-500 pr-2">
          {generateYAxisLabels(maxValue).map((label, index) => (
            <div key={index}>{label}</div>
          ))}
        </div>

        {/* Chart area */}
        <div className="absolute left-10 right-0 bottom-6 top-0">
          {/* Wednesday vertical line */}
          {wednesdayIndex !== -1 && (
            <div
              className="absolute top-0 bottom-0 border-l border-green-300 border-dashed z-10"
              style={{ left: `${wednesdayPosition}%` }}
            ></div>
          )}

          {/* Area chart */}
          <svg
            className="absolute inset-0 w-full h-full"
            preserveAspectRatio="none"
          >
            {/* Area fill with gradient */}
            <defs>
              <linearGradient
                id="chartGradient"
                x1="0%"
                y1="0%"
                x2="0%"
                y2="100%"
              >
                <stop offset="0%" stopColor="#10B981" stopOpacity="0.2" />
                <stop offset="100%" stopColor="#10B981" stopOpacity="0.05" />
              </linearGradient>
            </defs>

            {/* The chart line path */}
            <path
              d={generatePathD(chartData, maxValue)}
              stroke="#10B981"
              strokeWidth="2"
              fill="none"
              strokeLinejoin="round"
              strokeLinecap="round"
            />

            {/* The area fill path */}
            <path
              d={generateAreaD(chartData, maxValue)}
              fill="url(#chartGradient)"
            />
          </svg>

          {/* Peak day callout */}
          <div
            className="absolute bg-green-100 rounded-lg p-2 text-xs text-green-800 z-20"
            style={{
              left: `${
                (chartData.findIndex((d) => d.day === peakDay.day) /
                  (chartData.length - 1)) *
                100
              }%`,
              top: "10px",
              transform: "translateX(-50%)",
            }}
          >
            {peakDay.day}
            <br />
            {peakDay.count.toLocaleString()}
          </div>
        </div>

        {/* X-axis labels */}
        <div className="absolute left-10 right-0 bottom-0 flex justify-between text-xs text-gray-500">
          {chartData.map((data, index) => (
            <div
              key={index}
              className={data.day === "Wed" ? "font-medium text-green-700" : ""}
            >
              {data.day}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

// Helper functions
function formatNumber(value: number): string {
  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(1)}M`;
  } else if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}k`;
  }
  return value.toString();
}

function generateYAxisLabels(max: number): string[] {
  const labels: string[] = [];
  const step = max / 5;
  
  for (let i = 5; i >= 0; i--) {
    const value = i * step;
    labels.push(formatNumber(value));
  }
  
  return labels;
}

function generatePathD(data: CustomerChartData[], maxValue: number): string {
  const height = 220 - 30; // Chart height minus bottom padding for x-axis
  
  return data.map((point, i) => {
    const x = `${(i / (data.length - 1)) * 100}%`;
    const y = height - (point.count / maxValue) * height;
    return `${i === 0 ? 'M' : 'L'} ${x},${y}`;
  }).join(' ');
}

function generateAreaD(data: CustomerChartData[], maxValue: number): string {
  const height = 220 - 30; // Chart height minus bottom padding for x-axis
  
  const linePath = data.map((point, i) => {
    const x = `${(i / (data.length - 1)) * 100}%`;
    const y = height - (point.count / maxValue) * height;
    return `${i === 0 ? 'M' : 'L'} ${x},${y}`;
  }).join(' ');
  
  // Complete the path by drawing to bottom right, bottom left, and back to start
  return `${linePath} L 100%,${height} L 0,${height} Z`;
}

export default CustomerActivityChart;