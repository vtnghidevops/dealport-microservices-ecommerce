export interface TestimonialItem {
  id: string;
  userName: string;
  avatar: string;
  reviewText: string;
  rating: number;
}

export interface TestimonialItemProps {
  testimonial: TestimonialItem;
}
