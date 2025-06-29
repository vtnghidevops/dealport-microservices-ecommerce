// src/services/comment.service.ts

import { Comment, Reply, CommentPaginatedResponse } from '@/types/comment.model';
import { ProductReview } from '@/types/product.model';
import ProductService from './product.service';

// Convert ProductReview from API to Comment type for frontend
const convertProductReviewToComment = (review: ProductReview): Comment => {
  return {
    id: review.id?.toString() || '',
    productId: review.productId.toString(),
    userId: review.userId?.toString() || '',
    userName: review.userName || '',
    content: review.comment,
    rating: review.rating,
    createdAt: review.createdAt || new Date().toISOString(),
    likes: 0,
    replies: [],
  };
};

export const commentService = {
  getCommentsByProductId: async (
    productId: string,
    page: number,
    limit: number,
    rating?: number,
  ): Promise<CommentPaginatedResponse> => {
    try {
      // Use the ProductService to fetch reviews from the API
      const response = await ProductService.getProductReviews(productId, page, limit);

      // Convert API responses to frontend Comment format
      const comments: Comment[] = response.reviews.map(convertProductReviewToComment);

      // If there are ratings filter, filter them on client-side
      // In a real app, you might want to add this filtering to the API call
      let filteredComments = comments;
      if (rating !== undefined) {
        filteredComments = comments.filter((c) => c.rating === rating);
      }

      return {
        data: filteredComments,
        pagination: {
          total: response.pagination.total_items,
          currentPage: response.pagination.current_page,
          totalPages: response.pagination.total_pages,
          limit: response.pagination.page_size,
        },
      };
    } catch (error) {
      console.error("Error fetching comments from API:", error);
      return {
        data: [],
        pagination: {
          total: 0,
          currentPage: page,
          totalPages: 0,
          limit,
        },
      };
    }
  },

  addComment: async (comment: Omit<Comment, 'id' | 'createdAt' | 'likes' | 'replies'>) => {
    try {
      // Convert comment to ProductReview format for API
      const reviewData: Omit<ProductReview, 'id'> = {
        productId: parseInt(comment.productId),
        userId: comment.userId, // Keep userId as string (UUID format)
        userName: comment.userName,
        rating: comment.rating,
        comment: comment.content,
      };

      // Call the API to add the review
      const response = await ProductService.addProductReview(reviewData);

      // Return a properly formatted Comment object
      return {
        id: response.id?.toString() || Date.now().toString(),
        productId: comment.productId,
        userId: comment.userId,
        userName: comment.userName,
        content: comment.content,
        rating: comment.rating,
        createdAt: response.createdAt || new Date().toISOString(),
        likes: 0,
        replies: []
      };
    } catch (error) {
      console.error("Error adding comment via API:", error);

      // Fallback to creating a mock comment if API fails
      const newComment: Comment = {
        ...comment,
        id: Date.now().toString(),
        createdAt: new Date().toISOString(),
        likes: 0,
        replies: []
      };

      return newComment;
    }
  },

  addReply: async (commentId: string, reply: Omit<Reply, 'id' | 'commentId' | 'createdAt' | 'likes'>) => {
    try {
      // Note: If there's no API endpoint for replies, we can implement this later
      // For now, handle this client-side
      const newReply: Reply = {
        ...reply,
        id: Date.now().toString(),
        commentId,
        createdAt: new Date().toISOString(),
        likes: 0
      };

      // We'd need to fetch the comment first, add the reply, then update it
      // This is a placeholder until API support is added

      return newReply;
    } catch (error) {
      console.error("Error adding reply:", error);

      // Fallback
      const newReply: Reply = {
        ...reply,
        id: Date.now().toString(),
        commentId,
        createdAt: new Date().toISOString(),
        likes: 0
      };

      return newReply;
    }
  },

  likeComment: async (commentId: string) => {
    try {
      // Note: If there's no API endpoint for liking comments, implement client-side
      // This is a placeholder until API support is added
      return {
        id: commentId,
        likes: 1,
        isLiked: true
      };
    } catch (error) {
      console.error("Error liking comment:", error);
      return {
        id: commentId,
        likes: 1,
        isLiked: true
      };
    }
  },

  likeReply: async (replyId: string) => {
    try {
      // Note: If there's no API endpoint for liking replies, implement client-side
      // This is a placeholder until API support is added
      return {
        id: replyId,
        likes: 1,
        isLiked: true
      };
    } catch (error) {
      console.error("Error liking reply:", error);
      return {
        id: replyId,
        likes: 1,
        isLiked: true
      };
    }
  }
}; 