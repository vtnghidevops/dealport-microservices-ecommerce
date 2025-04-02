import React from 'react';
import { Line } from 'react-chartjs-2';
import { 
  Chart as ChartJS, 
  CategoryScale, 
  LinearScale, 
  PointElement, 
  LineElement, 
  Title, 
  Tooltip, 
  Filler, 
  Legend 
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend
);

interface WeeklyReportChartProps {
  data: {
    days: string[];
    values: number[];
    highlights?: {
      day: string;
      value: number;
      label: string;
    }[];
  };
}

const WeeklyReportChart: React.FC<WeeklyReportChartProps> = ({ data }) => {
  // Aqua-spring color theme
  const aquaSpringColor = 'rgba(126, 192, 151, 1)';
  const aquaSpringBgColor = 'rgba(126, 192, 151, 0.2)';
  
  const chartData = {
    labels: data.days,
    datasets: [
      {
        fill: true,
        label: 'Weekly Report',
        data: data.values,
        borderColor: aquaSpringColor,
        backgroundColor: aquaSpringBgColor,
        borderWidth: 2,
        tension: 0.4,
        pointBackgroundColor: aquaSpringColor
      }
    ]
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: 'rgba(240, 253, 244, 1)',
        titleColor: '#065f46',
        bodyColor: '#065f46',
        padding: 10,
        cornerRadius: 8,
        bodyFont: {
          size: 12
        },
        titleFont: {
          size: 12,
          weight: 'bold' as const
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
        border: {
          display: false
        },
        ticks: {
          font: {
            size: 12,
          },
          color: '#94a3b8',
          padding: 8,
          callback: function(tickValue: number | string) {
            return tickValue + 'k';
          }
        }
      },
      x: {
        grid: {
          display: false,
        },
        border: {
          display: false
        },
        ticks: {
          font: {
            size: 12,
            weight: 500
          },
          color: '#94a3b8',
          padding: 8
        }
      }
    },
    elements: {
      point: {
        radius: 0, // Hide points by default
        hitRadius: 10
      }
    }
  };

  // Find active day (e.g., "Wed" in your image)
  const activeDay = "Wed";
  const activeDayIndex = data.days.indexOf(activeDay);

  return (
    <div className="relative h-[250px]">
      <Line data={chartData} options={options} height={256} />
      
      {/* Highlights */}
      {data.highlights?.map((highlight, index) => {
        const dayIndex = data.days.indexOf(highlight.day);
        const xPercent = (dayIndex / (data.days.length - 1)) * 100;
        
        return (
          <div 
            key={index}
            className="absolute bg-green-50 px-4 py-2 rounded-lg text-sm text-green-800 shadow-sm border border-green-100"
            style={{
              left: `${xPercent}%`,
              top: '20%',
              transform: 'translate(-50%, 0)'
            }}
          >
            <div className="font-medium">{highlight.day}</div>
            <div className="font-bold">{highlight.value}k</div>
          </div>
        );
      })}

      {/* Active day indicator (vertical dashed line) */}
      {activeDayIndex !== -1 && (
        <div 
          className="absolute bottom-0 border-l-2 border-dashed border-green-500 h-[80%]"
          style={{
            left: `${(activeDayIndex / (data.days.length - 1)) * 100}%`,
            bottom: '10%'
          }}
        ></div>
      )}
    </div>
  );
};

export default WeeklyReportChart;