// components/admin/role/models/adminRole.model.ts
export interface SocialMedia {
  google?: boolean;
  facebook?: boolean;
  twitter?: boolean;
  linkedin?: boolean;
}

export interface AdminRole {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  profileImage: string;
  dateOfBirth: string;
  location: string;
  creditCard: string;
  biography?: string;
  socialMedia: SocialMedia;
}

export interface PasswordChangeData {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}