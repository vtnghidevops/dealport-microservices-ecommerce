import React from 'react';

interface BadgeProps {
  label: string;
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'info';
  size?: 'sm' | 'md';
}

// Badge component for displaying status or other labels
const Badge: React.FC<BadgeProps> = ({ 
  label, 
  variant = 'primary', 
  size = 'md' 
}) => {
  // Style definitions based on variant
  const variantStyles = {
    primary: 'bg-blue-100 text-blue-800',
    success: 'bg-green-100 text-green-800',
    warning: 'bg-yellow-100 text-yellow-800',
    danger: 'bg-red-100 text-red-800',
    info: 'bg-gray-100 text-gray-800',
  };
  
  // Size definitions
  const sizeStyles = {
    sm: 'text-xs px-2 py-0.5',
    md: 'text-sm px-2.5 py-1',
  };

  return (
    <span 
      className={`
        inline-block rounded-full font-medium 
        ${variantStyles[variant]} 
        ${sizeStyles[size]}
      `}
    >
      {label}
    </span>
  );
};

export default Badge;