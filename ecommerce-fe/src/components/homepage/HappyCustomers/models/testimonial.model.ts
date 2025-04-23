export interface TestimonialItem {
  id: string;
  name: string;
  avatarUrl: string;
  review: string;
  rating: number;
}

export interface TestimonialItemProps {
  testimonial: TestimonialItem;
}
