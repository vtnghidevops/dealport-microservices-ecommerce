export interface Comment {
  id: string;
  productId: string;
  userId: string;
  userName: string;
  userAvatar?: string;
  content: string;
  rating: number;
  createdAt: string;
  isAdmin?: boolean;
  likes: number;
  replies?: Reply[];
  isLiked?: boolean; // Indicates if the current user has liked this comment
}

export interface Reply {
  id: string;
  commentId: string;
  userId: string;
  userName: string;
  userAvatar?: string;
  content: string;
  createdAt: string;
  isAdmin?: boolean;
  likes: number;
  isLiked?: boolean;
}

export interface PaginatedResponse {
  data: Comment[];
  pagination: {
    total: number;
    currentPage: number;
    totalPages: number;
    limit: number;
  };
}