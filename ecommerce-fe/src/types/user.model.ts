export type UserRole = 'user' | 'admin';
export type UserStatus = 'active' | 'inactive' | 'suspended' | 'pending';

export interface UserAddress {
  id: string;
  isDefault: boolean;
  name: string;
  phone: string;
  street: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  type?: 'billing' | 'shipping';
}


export interface UserPaymentMethod {
  id: string;
  type: 'credit_card' | 'mono' | 'other';
  provider?: string;
  accountNumber?: string;
  expiryDate?: string;
  isDefault?: boolean;
}


export interface UserProfile {
  firstName: string;
  lastName: string;
  avatar?: string;
  phone?: string;
  dateOfBirth?: string;
  gender?: 'male' | 'female' | 'other' | 'prefer_not_to_say';
}

export interface User {
  id: string;
  email: string;
  username: string;
  role: UserRole;
  status: UserStatus;
  isActive: boolean;
  profile?: UserProfile;
  addresses?: UserAddress[];
  paymentMethods?: UserPaymentMethod[];
  wishlist?: string[]; // Product IDs
  createdAt?: string;
  updatedAt?: string;
  lastLogin?: string;
  cartId?: string;
  orderCount?: number;
}

export interface UserLoginCredentials {
  email: string;
  password: string;
}

export interface UserRegistrationData {
  email: string;
  password: string;
  confirmPassword: string;
  firstName: string;
  lastName: string;
  username: string;
  phone?: string;
  acceptTerms: boolean;
}

export interface AuthState {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}
