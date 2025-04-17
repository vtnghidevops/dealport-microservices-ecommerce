export type UserRole = 'customer' | 'admin' | 'manager' | 'support';
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
  type: 'credit_card' | 'paypal' | 'bank_transfer' | 'other';
  provider?: string;
  accountNumber?: string;
  expiryDate?: string;
  isDefault?: boolean;
}

export interface UserPreferences {
  newsletter: boolean;
  marketingEmails: boolean;
  orderNotifications: boolean;
  twoFactorAuth: boolean;
  language: string;
  currency: string;
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
  username?: string;
  profile: UserProfile;
  role: UserRole;
  status: UserStatus;
  addresses: UserAddress[];
  paymentMethods?: UserPaymentMethod[];
  preferences?: UserPreferences;
  createdAt: string;
  updatedAt: string;
  lastLogin?: string;
  wishlist?: string[]; // Product IDs
  recentlyViewed?: string[]; // Product IDs
  cartId?: string;
  orderCount?: number;
  totalSpent?: number;
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
  phone?: string;
  acceptTerms: boolean;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}
