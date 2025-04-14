import React, { useState } from 'react';
import { useCart } from '@/hooks/useCart'

interface MessageState {
  text: string;
  isError: boolean;
}
const CouponCode: React.FC = () => {
  const [couponCode, setCouponCode] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<MessageState>({ text: '', isError: false });
  const { applyCoupon } = useCart();

  const validateCouponCode = (code: string): boolean => {
    return code.length >= 3 && code.length <= 20;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedCode = couponCode.trim();
    
    if (!trimmedCode) {
      setMessage({ text: 'Please enter a coupon code', isError: true });
      return;
    }

    if (!validateCouponCode(trimmedCode)) {
      setMessage({ text: 'Invalid coupon code format', isError: true });
      return;
    }
    
    setIsSubmitting(true);
    setMessage({ text: '', isError: false });
    
    try {
      const success = await applyCoupon(trimmedCode);
      if (success) {
        setMessage({ text: 'Coupon applied successfully!', isError: false });
        setCouponCode('');
      } else {
        setMessage({ text: 'Invalid coupon code', isError: true });
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';
      setMessage({ 
        text: `Failed to apply coupon: ${errorMessage}`, 
        isError: true 
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-5">
      <h2 className="text-lg font-bold mb-4">Coupon Code</h2>
      
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Enter coupon code"
          className={`mb-8 w-full p-3 border rounded transition-colors
            ${message.isError ? 'border-red-300' : 'border-gray-300'}
            ${isSubmitting ? 'bg-gray-50' : 'bg-white'}`}
          value={couponCode}
          onChange={(e) => {
            setCouponCode(e.target.value);
            if (message.text) setMessage({ text: '', isError: false });
          }}
          disabled={isSubmitting}
        />
        
        {message.text && (
          <div className={`text-sm mb-3 ${message.isError ? 'text-error' : 'text-success'}`}>
            {message.text}
          </div>
        )}
        
        <button
          type="submit"
          disabled={isSubmitting || !couponCode.trim()}
          className={`w-1/2 p-3 rounded font-medium transition-all
            ${isSubmitting || !couponCode.trim() 
              ? 'bg-gray-400 cursor-not-allowed' 
              : 'bg-[#0496FF] hover:bg-blue-600'} 
            text-white`}
        >
          {isSubmitting ? (
            <span className="flex items-center justify-center">
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Applying...
            </span>
          ) : 'APPLY COUPON'}
        </button>
      </form>
    </div>
  );
};

export default CouponCode;