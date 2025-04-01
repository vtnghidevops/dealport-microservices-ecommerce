// components/admin/role/services/adminRole.service.ts
import { AdminRole } from '../models/adminRole.model';

// Mock data for development
const mockAdminRole: AdminRole = {
  id: '1',
  firstName: 'Wade',
  lastName: 'Warren',
  email: 'wade.warren@example.com',
  phoneNumber: '(406) 555-0120',
  profileImage: '/images/common/avatars/admin.png',
  dateOfBirth: '12-January-1999',
  location: '2972 Westheimer Rd. Santa Ana, Illinois 85486',
  creditCard: '843-4359-4444',
  biography: '',
  socialMedia: {
    google: true,
    facebook: true,
    twitter: true
  }
};

export const fetchAdminRole = async (): Promise<AdminRole> => {
  // In a real app, you'd make an API call here
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(mockAdminRole);
    }, 500);
  });
};

export const updateAdminRole = async (adminRole: AdminRole): Promise<AdminRole> => {
  // In a real app, you'd make an API call here
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(adminRole);
    }, 500);
  });
};

export const changePassword = async (
  currentPassword: string,
  newPassword: string
): Promise<boolean> => {
  // In a real app, you'd make an API call here
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(true);
    }, 500);
  });
};