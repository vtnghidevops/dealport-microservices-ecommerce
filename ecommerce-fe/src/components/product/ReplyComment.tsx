import React, { useEffect, useState, useRef } from 'react';

interface ReplyCommentProps {
  commentId: string;
  userName: string;
  onSubmit: (commentId: string, content: string) => void;
  onCancel: () => void;
}

const ReplyComment: React.FC<ReplyCommentProps> = ({
  commentId,
  userName,
  onSubmit,
  onCancel
}) => {
  const [replyContent, setReplyContent] = useState(`@${userName} `); // Initialize with @mention
  const [isSubmitting, setIsSubmitting] = useState(false); // Add submitting state
  const inputRef = useRef<HTMLInputElement>(null);

  // Focus input when component mounts and place cursor at end
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
      inputRef.current.setSelectionRange(replyContent.length, replyContent.length);
    }
  }, [replyContent]);

  const handleSubmit = async () => {
    if (!replyContent.trim() || isSubmitting) return;

    try {
      setIsSubmitting(true);
      // Remove @mention from the start of content before submitting
      const content = replyContent.replace(`@${userName} `, '').trim();
      await onSubmit(commentId, content);
      setReplyContent("");
    } catch (error) {
      console.error("Error submitting reply:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey && !isSubmitting) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleChangeMention = (e: React.ChangeEvent<HTMLInputElement>) => {
    const mention = `@${userName} `;
    const newValue = e.target.value;

    // Always ensure @mention stays at the beginning
    if (!newValue.startsWith(mention)) {
      setReplyContent(mention + newValue.replace(mention, ''));
    } else {
      setReplyContent(newValue);
    }
  };

  return (
    <div className="bg-gray-50 p-6 rounded-lg mt-5 w-[80%]">
      <div className="flex items-center ml-[3rem] mb-3">
        <button className="mr-2" onClick={onCancel}>
          <svg
            className="w-5 h-5 mr-1 text-cyprus font-medium"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
          >
            <path d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
          </svg>
        </button>
        <span className="text-base font-medium">
          Repling: <strong>{userName}</strong>
        </span>
        <button className="ml-auto" onClick={onCancel}>
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M18 6L6 18M6 6L18 18"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </button>
      </div>

      <div className="flex items-start">
        <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center text-white font-semibold mr-4">
          <span>A</span>
        </div>

        <div className="flex-1">
          <div className="rounded-lgoverflow-hidden">
            <div className="relative w-full mr-2 flex items-center">
              <input
                ref={inputRef}
                className="ml-4 h-[50px] flex items-center pl-5 border border-gray-300 rounded-lg overflow-hidden w-full p-4 focus:outline-none"
                placeholder={`Reply to ${userName}...`}
                value={replyContent}
                onChange={handleChangeMention}
                onKeyDown={handleKeyDown}
                disabled={isSubmitting}
                maxLength={3000}
              />
              <span className="absolute right-1 bottom-1 text-gray-500 mr-2">
                {replyContent.length}/3000
              </span>
            </div>
          </div>
          <div className="flex justify-between items-center">
            <div className="ml-[5.5rem] w-[10rem] mt-5 flex items-center">
              <button
                type="submit"
                onClick={handleSubmit}
                className="font-[18px] ml-[-5rem] cursor-pointer h-[50px] w-full px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                disabled={!replyContent.trim() || isSubmitting}
              >
                {isSubmitting ? "Sending..." : "Send Comment"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ReplyComment;