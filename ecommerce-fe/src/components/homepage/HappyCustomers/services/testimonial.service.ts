import { TestimonialItem } from "../models/testimonial.model";

// Mock data - trong thực tế sẽ fetch từ API
const testimonialData: TestimonialItem[] = [
  {
    id: '1',
    name: 'Emily R.',
    avatarUrl: '/images/common/avatars/emily.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  },
  {
    id: '2',
    name: 'John D',
    avatarUrl: '/images/common/avatars/johnd.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  },
  {
    id: '3',
    name: 'Ahmed M.',
    avatarUrl: '/images/common/avatars/ahmedM.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  },
  {
    id: '4',
    name: 'Alex T',
    avatarUrl: '/images/common/avatars/alexT.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  },
  {
    id: '5',
    name: 'Priya R',
    avatarUrl: '/images/common/avatars/priyaR.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  },
  {
    id: '6',
    name: 'David H',
    avatarUrl: '/images/common/avatars/davidH.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  },
  {
    id: '7',
    name: 'Ahmen M',
    avatarUrl: '/images/common/avatars/ahmedN.png',
    review: '"Fast delivery and fantastic quality! The customer support team was quick to resolve my query. Dealport has earned a loyal customer."',
    rating: 5,
  }
];

export const TestimonialService = {
  getTestimonials: (): Promise<TestimonialItem[]> => {
    return Promise.resolve(testimonialData);
  }
};