import React, { useState, useEffect } from 'react';
import { cn } from '@/lib/utils';

interface TransitionContainerProps {
  isLoading: boolean;
  children: React.ReactNode;
  skeleton: React.ReactNode;
  className?: string;
  duration?: number; // Thời gian chuyển tiếp (transition duration) tính bằng ms
  delay?: number; // Độ trễ trước khi bắt đầu chuyển tiếp (ms)
}

export const TransitionContainer: React.FC<TransitionContainerProps> = ({
  isLoading,
  children,
  skeleton,
  className,
  duration = 400,
  delay = 100
}) => {
  const [showSkeleton, setShowSkeleton] = useState(isLoading);
  const [showContent, setShowContent] = useState(!isLoading);
  const [skeletonOpacity, setSkeletonOpacity] = useState(1);
  const [contentOpacity, setContentOpacity] = useState(isLoading ? 0 : 1);

  useEffect(() => {
    let skeletonTimer: NodeJS.Timeout;
    let contentTimer: NodeJS.Timeout;

    if (!isLoading && showSkeleton) {
      // Khi loading hoàn tất, bắt đầu chuyển đổi
      setSkeletonOpacity(0); // Bắt đầu fade out skeleton

      skeletonTimer = setTimeout(() => {
        setShowSkeleton(false); // Ẩn skeleton hoàn toàn sau khi đã fade out
        setShowContent(true); // Hiển thị content

        contentTimer = setTimeout(() => {
          setContentOpacity(1); // Fade in content
        }, 50); // Độ trễ nhỏ để đảm bảo content đã được render
      }, duration);
    } else if (isLoading && !showSkeleton) {
      // Khi quay lại trạng thái loading
      setContentOpacity(0); // Fade out content

      contentTimer = setTimeout(() => {
        setShowContent(false); // Ẩn content
        setShowSkeleton(true); // Hiển thị skeleton

        skeletonTimer = setTimeout(() => {
          setSkeletonOpacity(1); // Fade in skeleton
        }, 50);
      }, duration);
    }

    return () => {
      clearTimeout(skeletonTimer);
      clearTimeout(contentTimer);
    };
  }, [isLoading, duration]);

  return (
    <div className={cn("relative", className)}>
      {showSkeleton && (
        <div
          className="transition-opacity w-full"
          style={{
            opacity: skeletonOpacity,
            transition: `opacity ${duration}ms ease-in-out ${delay}ms`
          }}
        >
          {skeleton}
        </div>
      )}

      {showContent && (
        <div
          className={cn("transition-opacity w-full", showSkeleton ? "absolute top-0 left-0" : "")}
          style={{
            opacity: contentOpacity,
            transition: `opacity ${duration}ms ease-in-out ${delay}ms`
          }}
        >
          {children}
        </div>
      )}
    </div>
  );
};

export default TransitionContainer; 