// src/components/product/services/comment.service.ts

import { Comment } from '../models/comment.model';
import { Reply } from '../models/comment.model';
import { PaginatedResponse } from '../models/comment.model';
const comments: Comment[] = [
  {
    id: '1',
    productId: '1',
    userId: '101',
    userName: 'Phan Khánh Linh',
    userAvatar: '/images/avatars/p.png', 
    content: 'Cho em hỏi sản phẩm này còn hàng ở gò đâu tây ninh không ạ ?',
    rating: 0, // Không đánh giá
    createdAt: '2024-03-11T00:00:00Z',
    likes: 0,
    replies: [
      {
        id: '101',
        commentId: '1',
        userId: 'admin1',
        userName: 'Thành Nhân',
        userAvatar: '/images/avatars/admin_logo.png',
        content: 'Chào chị Linh, Dạ, Laptop MSI Gaming Thin A15 B7UC-261VN R5 7535HS/16GB/512GB/15.6" FHD/RTX3050_4GB/Win11_Balo với thiết kế ấn tượng cùng cấu hình vượt trội, sản phẩm đang có giá ưu đãi chỉ còn 17.490.000 đ áp dụng đến 13/03. Mẫu này chưa có sẵn hàng tại Tây Ninh, chị tham khảo chờ hàng từ 3-5 ngày làm việc ạ. Nếu cần thêm thông tin khác chị gọi tổng đài miễn phí 18006601 hoặc có thể chat qua Zalo tại đây. Thân mến!',
        createdAt: '2024-03-11T01:00:00Z',
        isAdmin: true,
        likes: 0
      }
    ]
  },
  {
    id: '2',
    productId: '1',
    userId: '102',
    userName: 'Mạc Quang Huy',
    userAvatar: '/images/avatars/m.png',
    content: 'Con này ở hải dương còn hàng không',
    rating: 0,
    createdAt: '2024-02-11T00:00:00Z',
    likes: 0,
    replies: []
  }
];

export const commentService = {
  getCommentsByProductId: async (
    productId: string,
    page: number,
    limit: number,
    rating?: number,
  ): Promise<PaginatedResponse> => {
    let filteredComments = comments
      .filter((c) => c.productId === productId)
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      );

    if (rating !== undefined) {
      filteredComments = filteredComments.filter((c) => c.rating === rating);
    }
    const total = filteredComments.length;
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;
    
    return {
      data: filteredComments.slice(startIndex, endIndex),
      pagination: {
        total,
        currentPage: page,
        totalPages: Math.ceil(total / limit),
        limit,
      },
    };
  },
  
  addComment: (comment: Omit<Comment, 'id' | 'createdAt' | 'likes' | 'replies'>) => {
    const newComment: Comment = {
      ...comment,
      id: Date.now().toString(),
      createdAt: new Date().toISOString(),
      likes: 0,
      replies: []
    };
    // comments.push(newComment);
    comments.unshift(newComment);
    return Promise.resolve(newComment);
  },
  
  addReply: (commentId: string, reply: Omit<Reply, 'id' | 'commentId' | 'createdAt' | 'likes'>) => {
    const comment = comments.find(c => c.id === commentId);
    if (!comment) return Promise.reject('Comment not found');
    
    const newReply: Reply = {
      ...reply,
      id: Date.now().toString(),
      commentId,
      createdAt: new Date().toISOString(),
      likes: 0
    };
    
    if (!comment.replies) {
      comment.replies = [];
    }
    comment.replies.push(newReply);
    console.log('New reply:', comment);
    return Promise.resolve(newReply);
  },
  
  likeComment: (commentId: string) => {
    const comment = comments.find(c => c.id === commentId);
    if (!comment) return Promise.reject('Comment not found');
    
    // Toggle like status instead of always incrementing
    if (comment.isLiked) {
        comment.likes -= 1;
        comment.isLiked = false;
    } else {
        comment.likes += 1;
        comment.isLiked = true;
    }
    
    return Promise.resolve({
        likes: comment.likes,
        isLiked: comment.isLiked
    });
  },
  
  likeReply: (commentId: string, replyId: string) => {
    const comment = comments.find(c => c.id === commentId);
    if (!comment) return Promise.reject('Comment not found');
    
    if (!comment.replies) {
      comment.replies = [];
    }
    const reply = comment.replies.find(r => r.id === replyId);
    if (!reply) return Promise.reject('Reply not found');
    
    reply.likes += 1;
    return Promise.resolve(reply);
  }
};