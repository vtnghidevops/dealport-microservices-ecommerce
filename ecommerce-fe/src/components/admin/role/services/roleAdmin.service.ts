// components/admin/role/services/roleAdmin.service.ts
import { AdminRole } from '../models/adminRole.model';
import axios from 'axios';


// Create an authenticated API client with the token
const getAuthClient = () => {
  const token = localStorage.getItem('token');
  return axios.create({
    baseURL: import.meta.env.VITE_PUBLIC_PRODUCT_API_URL,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
};

export const fetchAdminRole = async (): Promise<AdminRole> => {
  try {
    // Try to get user data from local storage first
    const userData = localStorage.getItem('user');
    if (!userData) {
      throw new Error('User data not found');
    }

    const user = JSON.parse(userData);

    // Get additional user details from API if needed
    const api = getAuthClient();
    const response = await api.get(`/user/${user.id}`);

    // Map the API response to AdminRole model
    return {
      id: user.id,
      firstName: user.profile?.firstName || response.data?.firstName || '',
      lastName: user.profile?.lastName || response.data?.lastName || '',
      email: user.email,
      phoneNumber: user.profile?.phone || response.data?.phone || '',
      profileImage: user.profile?.avatar || '/images/common/avatars/admin.png',
      dateOfBirth: user.profile?.dateOfBirth || response.data?.dateOfBirth || '',
      location: user.profile?.address || response.data?.address || '',
      creditCard: response.data?.creditCard || '',
      biography: user.profile?.bio || response.data?.bio || '',
      socialMedia: {
        google: true,
        facebook: true,
        twitter: true
      }
    };
  } catch (error) {
    console.error('Error fetching admin data:', error);

    // Fallback to user data in localStorage if API call fails
    const userData = localStorage.getItem('user');
    if (userData) {
      const user = JSON.parse(userData);
      return {
        id: user.id,
        firstName: user.profile?.firstName || '',
        lastName: user.profile?.lastName || '',
        email: user.email,
        phoneNumber: user.profile?.phone || '',
        profileImage: user.profile?.avatar || '/images/common/avatars/admin.png',
        dateOfBirth: user.profile?.dateOfBirth || '',
        location: user.profile?.address || '',
        creditCard: '',
        biography: user.profile?.bio || '',
        socialMedia: {
          google: true,
          facebook: true,
          twitter: true
        }
      };
    }

    throw new Error('Failed to fetch admin data');
  }
};

export const updateAdminRole = async (adminRole: AdminRole): Promise<AdminRole> => {
  try {
    const api = getAuthClient();

    // Create update data in format expected by API
    const updateData = {
      profile: {
        firstName: adminRole.firstName,
        lastName: adminRole.lastName,
        phone: adminRole.phoneNumber,
        dateOfBirth: adminRole.dateOfBirth,
        address: adminRole.location,
        bio: adminRole.biography
      }
    };

    // Call API to update profile
    const response = await api.put(`/user/${adminRole.id}`, updateData);

    // Update local storage with new data
    const userData = localStorage.getItem('user');
    if (userData) {
      const user = JSON.parse(userData);
      user.profile = {
        ...user.profile,
        ...updateData.profile
      };
      localStorage.setItem('user', JSON.stringify(user));
    }

    return adminRole;
  } catch (error) {
    console.error('Error updating admin data:', error);
    throw new Error('Failed to update admin data');
  }
};

export const changePassword = async (
  currentPassword: string,
  newPassword: string
): Promise<boolean> => {
  try {
    const api = getAuthClient();
    await api.post('/auth/password-update', {
      current_password: currentPassword,
      password: newPassword
    });
    return true;
  } catch (error) {
    console.error('Error changing password:', error);
    throw error;
  }
};