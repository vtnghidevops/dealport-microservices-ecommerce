import React, { useState, useRef, useEffect } from 'react';
import { commentService } from "@/services/product/comment.service";
import { useAuth } from "@/hooks/useAuth";

interface RatingModalProps {
  isOpen: boolean;
  onClose: () => void;
  productId: string;
  productName: string;
  productImage?: string;
  onRatingSubmit?: () => void;
}

const RatingModal: React.FC<RatingModalProps> = ({
  isOpen,
  onClose,
  productId,
  productName,
  productImage,
  onRatingSubmit
}) => {
  const [rating, setRating] = useState<number>(5);
  const [hoverRating, setHoverRating] = useState<number | null>(null);
  const [comment, setComment] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { authState } = useAuth();
  const modalRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Predefined rating labels
  const ratingLabels = ['Very Bad', 'Bad', 'Normal', 'Good', 'Excellent'];

  useEffect(() => {
    // Focus the textarea when modal opens
    if (isOpen && textareaRef.current) {
      setTimeout(() => {
        textareaRef.current?.focus();
      }, 100);
    }

    // Close modal on Escape key
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    // Close modal if clicking outside
    const handleClickOutside = (e: MouseEvent) => {
      if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown);
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, onClose]);

  // Prevent body scrolling when modal is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'auto';
    }
    return () => {
      document.body.style.overflow = 'auto';
    };
  }, [isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitting) return;

    setIsSubmitting(true);
    try {
      if (!authState.isAuthenticated) {
        alert('Please log in to submit a rating');
        return;
      }

      await commentService.addComment({
        productId,
        userId: authState.user?.id || '',
        userName: authState.user?.profile ?
          `${authState.user.profile.firstName} ${authState.user.profile.lastName}` :
          authState.user?.username || '',
        content: comment,
        rating,
      });

      setComment('');
      setRating(5);
      setHoverRating(null);
      onRatingSubmit?.();
      onClose();
    } catch (error) {
      console.error("Error submitting rating:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div
        ref={modalRef}
        className="bg-white rounded-lg shadow-xl w-[600px] max-w-[95vw] max-h-[95vh] overflow-y-auto"
      >
        <div className="flex justify-between items-center p-4 border-b">
          <h2 className="text-xl font-semibold">Rating & Comment</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700"
            aria-label="Close"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div className="p-6">
          <div className="flex mb-6">
            <div className="w-16 h-16 border border-gray-200 rounded-md overflow-hidden flex-shrink-0">
              {productImage ? (
                <img
                  src={productImage}
                  alt={productName}
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="w-full h-full bg-gray-100 flex items-center justify-center text-gray-400">
                  <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </div>
              )}
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium">{productName}</h3>
            </div>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="mb-6">
              <h3 className="text-lg font-medium mb-4">Overall Rating</h3>
              <div className="flex justify-between items-center max-w-[500px] mb-4">
                {[1, 2, 3, 4, 5].map((starValue) => (
                  <div
                    key={starValue}
                    className="flex flex-col items-center space-y-2"
                  >
                    <button
                      type="button"
                      onClick={() => setRating(starValue)}
                      onMouseEnter={() => setHoverRating(starValue)}
                      onMouseLeave={() => setHoverRating(null)}
                      className="focus:outline-none transition duration-150"
                      aria-label={`Rate ${starValue} stars`}
                    >
                      <svg
                        className={`w-10 h-10 ${(hoverRating !== null ? starValue <= hoverRating : starValue <= rating)
                          ? "text-yellow-400"
                          : "text-gray-300"
                          }`}
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                      </svg>
                    </button>
                    <span className="text-sm text-gray-600">{ratingLabels[starValue - 1]}</span>
                  </div>
                ))}
              </div>
            </div>

            

            <div className="mb-6">
              <div className="relative">
                <textarea
                  ref={textareaRef}
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                  placeholder="Xin mời chia sẻ một số cảm nhận về sản phẩm (nhập tối thiểu 15 kí tự)"
                  className="w-full min-h-[150px] p-4 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  maxLength={3000}
                ></textarea>
                <span className="absolute bottom-2 right-2 text-sm text-gray-500">
                  {comment.length}/3000
                </span>
              </div>
            </div>

            <div className="flex justify-center pb-4">
              <button
                type="submit"
                disabled={comment.length < 15 || isSubmitting}
                className="bg-red-600 text-white font-semibold py-3 px-6 rounded-lg text-lg w-full max-w-xs hover:bg-red-700 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? "Đang gửi..." : "GỬI ĐÁNH GIÁ"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default RatingModal; 