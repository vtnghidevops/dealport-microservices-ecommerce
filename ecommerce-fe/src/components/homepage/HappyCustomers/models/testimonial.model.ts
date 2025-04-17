// Backend
//     id: 1,
//     name: 'Emily R.',
//     avatarUrl: '/images/common/avatars/emily.png',
//     review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
//     rating: 5,
export interface TestimonialItem {
  id: string;
  userId?: string;
  name: string;
  avatarUrl?: string;
  review: string;
  rating: number;
  createdAt?: string;
  productId?: string;
  productName?: string;
  status?: 'published' | 'pending' | 'rejected';
  helpfulCount?: number;
}

// stand for 5 item with only 1 var
export interface TestimonialItemProps {
  testimonial: TestimonialItem;
}

export interface TestimonialFilter {
  rating?: number;
  sortBy?: 'newest' | 'highest-rating' | 'lowest-rating' | 'most-helpful';
  productId?: string;
  limit?: number;
}