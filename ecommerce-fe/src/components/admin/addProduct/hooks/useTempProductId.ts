import { useState, useEffect } from 'react';

/**
 * Hook to manage a temporary product ID for image uploads 
 * before a product is created
 */
export function useTempProductId() {
  const [tempProductId, setTempProductId] = useState<string>('0');

  useEffect(() => {
    // Check for existing temp ID in localStorage
    const storedId = localStorage.getItem('tempProductId');

    if (!storedId) {
      // Generate a new temporary negative ID to avoid conflicts with real products
      // We use negative numbers to clearly indicate temporary status
      const newTempId = `-${Date.now()}`;
      localStorage.setItem('tempProductId', newTempId);
      setTempProductId(newTempId);
    } else {
      setTempProductId(storedId);
    }

    // Cleanup on component unmount
    return () => {
      // Only clear if it was a temporary ID (negative number)
      const currentId = localStorage.getItem('tempProductId');
      if (currentId && currentId.startsWith('-')) {
        localStorage.removeItem('tempProductId');
      }
    };
  }, []);

  // Function to update the temp ID with a real product ID when it's created
  const updateWithRealProductId = (realId: number) => {
    const realIdStr = realId.toString();
    localStorage.setItem('tempProductId', realIdStr);
    setTempProductId(realIdStr);
  };

  // Function to reset the temp ID
  const resetTempProductId = () => {
    const newTempId = `-${Date.now()}`;
    localStorage.setItem('tempProductId', newTempId);
    setTempProductId(newTempId);
  };

  return {
    tempProductId,
    updateWithRealProductId,
    resetTempProductId
  };
} 