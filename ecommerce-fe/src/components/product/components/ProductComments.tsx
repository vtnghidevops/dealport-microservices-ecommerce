
import React, { useEffect, useState, useRef } from "react";
import { Comment } from "@/types/comment.model";
import { commentService } from "@/services/comment.service";
import ReplyComponent from "./ReplyComment";
import Pagination from "@/components/common/Pagination";

interface ProductCommentsProps {
  productId: string;
}

const ProductComments: React.FC<ProductCommentsProps> = ({ productId }) => {
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState("");
  const [rating, setRating] = useState<number>(5);
  const [loading, setLoading] = useState(false);
  const [replyTo, setReplyTo] = useState<string | null>(null);
  const [activeFilter, setActiveFilter] = useState<number | "all">("all");
  const [currentUser] = useState({ avatar: "", name: "Test User" }); // Add mock user data

  // Thêm state cho pagination
  const [currentPage, setCurrentPage] = useState(1);
  // const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const commentSectionRef = useRef<HTMLDivElement>(null);
  const limit = 5; 
  const [expandedReplies, setExpandedReplies] = useState<string[]>([]);

  const fetchComments = async () => {
    try {
      setLoading(true);
      const response = await commentService.getCommentsByProductId(
        productId,
        currentPage,
        limit,
        activeFilter === "all" ? undefined : (activeFilter as number)
      );
    
    // Store old replies state with full state information
    const oldRepliesState = new Map(
      comments.flatMap(comment => 
        (comment.replies || []).map(reply => 
          [`${comment.id}-${reply.id}`, {
            isLiked: reply.isLiked,
            likes: reply.likes,
            id: reply.id,
            content: reply.content,
            userName: reply.userName,
            userAvatar: reply.userAvatar,
            isAdmin: reply.isAdmin,
            createdAt: reply.createdAt
          }]
        )
      )
    );
    // If filtering by rating, don't show replies
    // If filter is "all", keep replies for expanded comments 
    // Apply old state to new replies
    const commentsWithPreservedReplies = response.data.map(comment => ({
      ...comment,
      replies: activeFilter === "all" 
        ? (comment.replies || []).map(reply => {
            const oldReplyState = oldRepliesState.get(`${comment.id}-${reply.id}`);
            if (oldReplyState) {
              // Preserve the entire state of the reply if it existed before
              return {
                ...reply,
                isLiked: oldReplyState.isLiked,
                likes: oldReplyState.likes
              };
            }
            return reply;
          })
        : []
    }));
      setComments(commentsWithPreservedReplies);
      setTotal(response.pagination.total);
    } catch (error) {
      console.error("Error fetching comments:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchComments();
    if (currentPage !== 1) {
      setCurrentPage(1);
    }
  }, [productId, currentPage, activeFilter]);


  const handlePageChange = (page: number) => {
    if (page === currentPage) return; // Prevent unnecessary updates
    setCurrentPage(page);
    setLoading(true); // Show loading when changing pages
    commentSectionRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;

    try {
      // add new comment
      const comment = await commentService.addComment({
        productId,
        userId: "user123",
        userName: currentUser.name,
        content: newComment,
        rating,
      });

      // Update UI right now
      setComments(prev => [comment, ...prev.slice(0, limit - 1)]);
      setTotal(prev => prev + 1);
      setNewComment("");
      setRating(5);

      // reset pagination to first page
      if (currentPage !== 1) {
        setCurrentPage(1);
      }
    } catch (error) {
      console.error("Error adding comment:", error);
    }
  };

  const handleReply = async (commentId: string, content: string) => {
    if (!content.trim()) return;
  
    try {
      const reply = await commentService.addReply(commentId, {
        userId: "user123",
        userName: "Admin",
        content: content,
      });
  
      // // Check if reply was already added to prevent duplicates
      // setComments((prev) =>
      //   prev.map((comment) => {
      //     if (comment.id !== commentId) return comment;
          
      //     // Check if reply already exists
      //     const replyExists = comment.replies?.some(r => r.id === reply.id);
      //     if (replyExists) return comment;
  
      //     return {
      //       ...comment,
      //       replies: [...(comment.replies || []), reply]
      //     };
      //   })
      // );
      setComments((prev) =>
        prev.map((comment) => {
          if (comment.id !== commentId) return comment;
          
          if (!expandedReplies.includes(comment.id)) {
            setExpandedReplies(prev => [...prev, comment.id]);
          }
  
          const replyExists = comment.replies?.some(r => r.id === reply.id);
          if (replyExists) return comment;
  
          return {
            ...comment,
            replies: [...(comment.replies || []), reply]
          };
        })
      );
  
      setReplyTo(null);
    } catch (error) {
      console.error("Error adding reply:", error);
    }
  };


  const handleLikeComment = async (commentId: string) => {
    const targetComment = comments.find(comment => comment.id === commentId);
    if (!targetComment) return;

    try {
        // Update UI first (optimistic update)
        setComments((prev) =>
            prev.map((comment) =>
                comment.id === commentId
                    ? {
                        ...comment,
                        likes: comment.isLiked ? comment.likes - 1 : comment.likes + 1,
                        isLiked: !comment.isLiked,
                    }
                    : comment
            )
        );

        // Call API
        const response = await commentService.likeComment(commentId);
        
        // Update with actual server response
        if (response) {
            setComments((prev) =>
                prev.map((comment) =>
                    comment.id === commentId
                        ? {
                            ...comment,
                            likes: response.likes,
                            isLiked: response.isLiked,
                        }
                        : comment
                )
            );
        }

    } catch (error) {
        console.error("Error toggling like:", error);
        // Revert to original state if error occurs
        setComments((prev) =>
            prev.map((comment) =>
                comment.id === commentId
                    ? {
                        ...comment,
                        likes: targetComment.likes,
                        isLiked: targetComment.isLiked,
                    }
                    : comment
            )
        );
    }
  };

  const handleLikeReply = async (commentId: string, replyId: string) => {
  const targetComment = comments.find((comment) => comment.id === commentId);
  const targetReply = targetComment?.replies?.find(
    (reply) => reply.id === replyId
  );
  if (!targetComment || !targetReply) return;

  try {
    // Optimistic update
    setComments((prev) =>
      prev.map((comment) => {
        if (comment.id !== commentId) return comment;
        return {
          ...comment,
          replies: comment.replies?.map((reply) => {
            if (reply.id !== replyId) return reply;
            return {
              ...reply,
              likes: reply.isLiked ? reply.likes - 1 : reply.likes + 1,
              isLiked: !reply.isLiked,
            };
          }),
        };
      })
    );

    // Call API
    await commentService.likeReply(commentId, replyId);
  } catch (error) {
    console.error("Error liking reply:", error);
    // Revert on error
    setComments((prev) =>
      prev.map((comment) => {
        if (comment.id !== commentId) return comment;
        return {
          ...comment,
          replies: comment.replies?.map((reply) => {
            if (reply.id !== replyId) return reply;
            return targetReply;
          }),
        };
      })
    );
  }
  };

  const handleToggleReply = (commentId: string) => {
    if (replyTo === commentId) {
      setReplyTo(null);
    } else {
      setReplyTo(commentId);
      if (!expandedReplies.includes(commentId)) {
        setExpandedReplies(prev => [...prev, commentId]);
      }
    }
  };


  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays <= 30) {
      return `${diffDays} ngày trước`;
    } else {
      const diffMonths = Math.floor(diffDays / 30);
      return `${diffMonths} tháng trước`;
    }
  };


  return (
    <div className="mt-12" ref={commentSectionRef}>
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-semibold mb-6">{total} Comment</h2>

        {/* Rating Filter */}
        <div className="flex items-center space-x-4 mb-8">
          {["all", 5, 4, 3, 2, 1].map((value) => (
            <button
              key={value}
              className={`flex items-center justify-center h-[40px] w-[75px] px-4 py-2 rounded-full transition-colors
          ${value === "all" ? "w-[75px]" : "w-[70px]"}
          ${
            activeFilter === value
              ? "border border-blue-500 text-blue-500"
              : "bg-white border border-gray-300 hover:bg-gray-100"
          }`}
              onClick={() => setActiveFilter(value as number | "all")}
            >
              {value === "all" ? (
                "Tất cả"
              ) : (
                <>
                  {value} <span className="ml-1">★</span>
                </>
              )}
            </button>
          ))}
        </div>
      </div>

      {/* Comment Form */}
      <div className="mb-8">
        <form onSubmit={handleSubmitComment}>
          <div className="flex justify-between items-center mt-5 gap-16 ">
            <div className="relative w-[90%] flex items-center">
              <input
                className="h-[50px] flex items-center pl-5 border border-gray-300 rounded-lg overflow-hidden w-full p-4 focus:outline-none"
                placeholder="Nhập nội dung bình luận"
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
              />
              <span className="absolute right-1 bottom-1 text-gray-500 mr-2">
                {newComment.length}/3000
              </span>
            </div>

            <div className="!w-[15%] bg-gray-50 p-4 flex justify-between items-center">
              <div className="flex items-center w-full">
                <button
                  type="submit"
                  className="font-[18px]  cursor-pointer h-[50px] w-full px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                  disabled={!newComment.trim()}
                >
                  Send Comment
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>

      {/* Comments List */}
      <div className="space-y-8">
        {loading ? (
          <div className="text-center">
            <div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="mt-2 text-gray-600">Loading Comment...</p>
          </div>
        ) : comments.length > 0 ? (
          <>
            {comments.map((comment) => (
              <div key={comment.id} className=" pb-8 mb-5">
                <div className="flex items-start">
                  <div className="w-[40px] h-[40px] mr-5 rounded-full bg-gray-300 flex items-center justify-center text-xl font-bold text-white">
                    {comment.userAvatar ? (
                      <img
                        src={comment.userAvatar}
                        alt={comment.userName}
                        className="w-full h-full rounded-full object-cover"
                      />
                    ) : (
                      comment.userName.charAt(0).toUpperCase()
                    )}
                  </div>

                  <div className="flex-1">
                    <div className="flex items-center">
                      <h3 className="font-bold text-lg">{comment.userName}</h3>
                      {comment.isAdmin && (
                        <span className="ml-2 px-2 py-1 bg-gray-200 text-xs rounded-md">
                          Admin
                        </span>
                      )}
                      <span className="ml-2 text-gray-500 text-sm">
                        • {formatDate(comment.createdAt)}
                      </span>
                    </div>

                    {/* Add Rating Display */}
                    {comment.rating > 0 && (
                      <div className="flex items-center mt-1 mb-2">
                        {[...Array(5)].map((_, index) => (
                          <svg
                            key={index}
                            className={`w-[15px] h-[15px] ${
                              index < comment.rating
                                ? "text-yellow-400"
                                : "text-gray-300"
                            }`}
                            fill="currentColor"
                            viewBox="0 0 20 20"
                          >
                            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                          </svg>
                        ))}
                      </div>
                    )}

                    <div className="mt-2 text-cyprus">{comment.content}</div>

                    <div className="mt-3 flex items-center space-x-4">
                      <button
                        className={`justify-center h-[30px] w-[60px] border rounded-full flex items-center transition-all duration-200
                  ${
                    comment.isLiked
                      ? "border-blue-500 text-blue-500 bg-blue-50"
                      : "border-gray-300 text-gray-500 hover:text-blue-500 hover:border-blue-500 hover:bg-blue-50"
                  }`}
                        onClick={() => handleLikeComment(comment.id)}
                      >
                        <svg
                          className={`w-5 h-5 mr-1 font-medium ${
                            comment.isLiked ? "text-blue-500" : "text-cyprus"
                          }`}
                          viewBox="0 0 24 24"
                          fill={comment.isLiked ? "currentColor" : "none"}
                          stroke="currentColor"
                          strokeWidth="2"
                        >
                          <path d="M14 9V5a3 3 0 00-3-3l-4 9v11h11.28a2 2 0 002-1.7l1.38-9a2 2 0 00-2-2.3H14z" />
                          <path d="M7 22H4a2 2 0 01-2-2v-7a2 2 0 012-2h3" />
                        </svg>
                        <span
                          className={`font-medium ${
                            comment.isLiked ? "text-blue-500" : "text-gray-800"
                          }`}
                        >
                          {comment.likes}
                        </span>
                      </button>

                      <button
                        className="!ml-[10px] justify-center h-[30px] w-[90px] border border-gray-300 rounded-full flex items-center text-gray-500 hover:text-blue-500 hover:border-blue-500 hover:bg-blue-50 transition-all duration-200"
                        onClick={() => handleToggleReply(comment.id)}
                      >
                        <svg
                          className="w-5 h-5 mr-1 text-cyprus font-medium"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          strokeWidth="2"
                        >
                          <path d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                        </svg>
                        <span className="text-gray-800 font-medium">
                          Reply
                        </span>
                      </button>
                        {/* Hide or Show Comment reply */}
                      {comment.replies && comment.replies.length > 0 && (
                        <button
                          className="min-w-[5rem] !ml-[10px] justify-center h-[30px] px-8 border border-gray-300 rounded-full flex items-center text-gray-500 hover:text-blue-500 hover:border-blue-500 hover:bg-blue-50 transition-all duration-200"
                          onClick={() => {
                            if (expandedReplies.includes(comment.id)) {
                              setExpandedReplies((prev) =>
                                prev.filter((id) => id !== comment.id)
                              );
                            } else {
                              setExpandedReplies((prev) => [
                                ...prev,
                                comment.id,
                              ]);
                            }
                          }}
                        >
                          <span className="text-gray-800 font-medium">
                            {expandedReplies.includes(comment.id)
                              ? "Hiden"
                              : `Show ${comment.replies.length} reply`}
                          </span>
                        </button>
                      )}
                    </div>

                    {/* Show Reply Form */}
                    {replyTo === comment.id && (
                      <ReplyComponent
                        replyToUser={comment.userName}
                        onSubmit={(content) => handleReply(comment.id, content)}
                        onClose={() => setReplyTo(null)}
                        currentUserAvatar={currentUser?.avatar} // Nếu có thông tin user hiện tại
                        currentUserInitial={currentUser?.name?.charAt(0) || "U"} // Lấy chữ cái đầu của tên user
                      />
                    )}

                    {/* Response Replies from Form */}
                    {comment.replies &&
                      comment.replies.length > 0 &&
                      expandedReplies.includes(comment.id) && (
                        <div className="mt-4 pl-4">
                          {comment.replies.map((reply) => (
                            <div key={reply.id} className="mt-5">
                              <div className="flex items-start bg-neutral-100 py-5 px-3 rounded-lg">
                                <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center text-md font-bold text-white mr-3">
                                  {reply.userAvatar ? (
                                    <img
                                      src={reply.userAvatar}
                                      alt={reply.userName}
                                      className="w-full h-full rounded-full object-cover"
                                    />
                                  ) : (
                                    reply.userName.charAt(0).toUpperCase()
                                  )}
                                </div>

                                <div className="flex-1">
                                  <div className="flex items-center">
                                    <h4 className="font-bold">
                                      {reply.userName}
                                    </h4>
                                    {reply.isAdmin && (
                                      <span className="ml-2 px-2 py-1 bg-gray-200 text-xs rounded-md">
                                        Admin
                                      </span>
                                    )}
                                    <span className="ml-2 text-gray-500 text-sm">
                                      • {formatDate(reply.createdAt)}
                                    </span>
                                  </div>

                                  <div className="mt-1 text-gray-800">
                                    <span className="text-blue-600 font-medium">
                                      @{comment.userName}
                                    </span>{" "}
                                    {reply.content}
                                  </div>

                                  <div className="mt-2">
                                    <button
                                      onClick={() =>
                                        handleLikeReply(comment.id, reply.id)
                                      }
                                      className={`bg-white justify-center h-[30px] w-[60px] border rounded-full flex items-center transition-all duration-200 ${
                                        reply.isLiked
                                          ? "text-blue-600 border-blue-600"
                                          : "text-gray-500 border-gray-500"
                                      } hover:text-blue-600 hover:border-blue-600 transition-colors
                                    `}
                                    >
                                      <svg
                                        className={`w-5 h-5 mr-1 ${
                                          reply.isLiked
                                            ? "text-blue-500"
                                            : "text-cyprus"
                                        }`}
                                        viewBox="0 0 24 24"
                                        fill={
                                          reply.isLiked
                                            ? "currentColor"
                                            : "none"
                                        }
                                        stroke="currentColor"
                                        strokeWidth="2"
                                      >
                                        <path d="M14 9V5a3 3 0 00-3-3l-4 9v11h11.28a2 2 0 002-1.7l1.38-9a2 2 0 00-2-2.3H14z" />
                                        <path d="M7 22H4a2 2 0 01-2-2v-7a2 2 0 012-2h3" />
                                      </svg>
                                      {reply.likes}
                                    </button>
                                  </div>
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                  </div>
                </div>
              </div>
            ))}

            {/* Pagination */}
            {comments.length > 0 && total > limit && (
              <div className="mt-8 border-t border-gray-200 pt-4">
                <Pagination
                  currentPage={currentPage}
                  totalItems={total} // Changed from totalPages to total
                  onPageChange={handlePageChange}
                  pageSize={limit}
                />
              </div>
            )}
          </>
        ) : (
          <div className="text-center py-8">
            <p className="text-gray-500">
              Chưa có bình luận nào. Hãy là người đầu tiên bình luận!
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductComments;
