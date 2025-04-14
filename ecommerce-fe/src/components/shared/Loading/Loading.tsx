import React from 'react';

interface LoadingProps {
  size?: 'small' | 'medium' | 'large';
  color?: string;
  fullscreen?: boolean;
}

const Loading: React.FC<LoadingProps> = ({ 
  size = 'medium', 
  color = '#FA8232',
  fullscreen = false
}) => {
  const sizeMap = {
    small: 'h-8 w-8 border-2',
    medium: 'h-12 w-12 border-4',
    large: 'h-16 w-16 border-4'
  };

  const containerClasses = fullscreen 
    ? "fixed inset-0 bg-black bg-opacity-50 z-50 flex justify-center items-center" 
    : "flex justify-center items-center p-8";

  return (
    <div className={containerClasses}>
      <div 
        className={`animate-spin rounded-full ${sizeMap[size]} border-t-transparent`}
        style={{ borderColor: `${color} transparent transparent ${color}` }}
      />
    </div>
  );
};

export default Loading;