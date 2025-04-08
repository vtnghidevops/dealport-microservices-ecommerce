// Backend
//     id: 1,
//     name: 'Emily R.',
//     avatarUrl: '/images/common/avatars/emily.png',
//     review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
//     rating: 5,
export interface TestimonialItem {
  id: string;
  name: string;
  avatarUrl?: string;
  review: string;
  rating: number;
}

// stand for 5 item with only 1 var
export interface TestimonialItemProps {
  testimonial: TestimonialItem;
}