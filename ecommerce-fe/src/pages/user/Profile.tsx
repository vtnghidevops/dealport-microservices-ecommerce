// pages/user/Profile.tsx
import React, { useState } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { useToast } from '@/hooks/use-toast';
import UserLayout from '../../components/layouts/UserLayout';
import UserAvatar from '../../components/user/UserAvatar';
import { FiEdit, FiCheck } from 'react-icons/fi';
import { User, UserProfile } from '@/types/user.model';

const Profile: React.FC = () => {
  const { authState, updateProfile } = useAuth();
  const { toast } = useToast();
  const [editingField, setEditingField] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    firstName: authState.user?.profile?.firstName || '',
    lastName: authState.user?.profile?.lastName || '',
    email: authState.user?.email || '',
    phone: authState.user?.profile?.phone || '',
    dateOfBirth: authState.user?.profile?.dateOfBirth || '',
    displayName: authState.user?.profile?.firstName || '',
    username: authState.user?.email?.split('@')[0] || '',
    secondaryEmail: '',
    country: 'Bangladesh',
    region: 'Dhaka',
    city: 'Dhaka',
    zipCode: '1207',
    // address: authState.user?.profile?.address || '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const startEditing = (fieldName: string) => {
    setEditingField(fieldName);
  };

  const cancelEditing = () => {
    setEditingField(null);
    // Reset form data to original values
    setFormData({
      ...formData,
      firstName: authState.user?.profile?.firstName || '',
      lastName: authState.user?.profile?.lastName || '',
      email: authState.user?.email || '',
      phone: authState.user?.profile?.phone || '',
      dateOfBirth: authState.user?.profile?.dateOfBirth || '',
    });
  };

  const updateField = async (fieldName: string) => {
    try {
      let updatedProfile: Partial<UserProfile> = {};

      if (fieldName === 'name') {
        updatedProfile = {
          firstName: formData.firstName,
          lastName: formData.lastName
        };
      } else if (fieldName === 'phone') {
        updatedProfile = {
          phone: formData.phone
        };
      } else if (fieldName === 'dateOfBirth') {
        updatedProfile = {
          dateOfBirth: formData.dateOfBirth
        };
      }

      // Create user update object
      const userUpdate: Partial<User> = {
        profile: {
          ...authState.user?.profile,
          ...updatedProfile
        } as UserProfile
      };

      await updateProfile(userUpdate);
      setEditingField(null);
      toast({
        title: 'Profile Updated',
        description: 'Your profile has been updated successfully',
        variant: 'success',
      });
    } catch (error) {
      toast({
        title: 'Update Failed',
        description: 'Failed to update profile',
        variant: 'error',
      });
    }
  };

  const handleAvatarChange = (file: File) => {
    // Here you would implement the logic to upload the avatar file
    console.log('Avatar file changed:', file);
    toast({
      title: 'Avatar Updated',
      description: 'Your profile picture has been updated successfully',
      variant: 'success',
    });
  };

  return (
    <UserLayout>
      <div className="max-w-4xl mx-auto p-5 border border-neutral-100 rounded-lg">
        <h1 className="text-[20px] font-medium font-sans mb-5">ACCOUNT SETTING</h1>

        <div className="bg-white rounded-lg shadow-sm p-6 mb-5">
          <div className="flex flex-col md:flex-row gap-8">
            <div className="w-full md:w-1/5 flex flex-col items-center">
              <UserAvatar
                user={authState.user}
                editable={true}
                size="xl"
                onImageChange={handleAvatarChange}
              />
            </div>

            <div className="w-full md:w-4/5">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5 mb-5">
                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Display name
                  </label>
                  <Input
                    value={formData.displayName}
                    readOnly
                    className="bg-gray-50"
                  />
                </div>

                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Username
                  </label>
                  <Input
                    value={formData.username}
                    readOnly
                    className="bg-gray-50"
                  />
                </div>
              </div>

              <div className="mb-5">
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Email
                </label>
                <Input
                  name="email"
                  value={formData.email}
                  readOnly
                  className="bg-gray-50"
                />
              </div>

              <div className="mb-5">
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Secondary Email
                </label>
                <Input
                  name="secondaryEmail"
                  value={formData.secondaryEmail}
                  placeholder="Add secondary email"
                  readOnly
                  className="bg-gray-50"
                />
              </div>

              <div className="mb-5">
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Phone Number
                </label>
                <div className="flex items-center">
                  <Input
                    name="phone"
                    value={formData.phone}
                    onChange={handleChange}
                    readOnly={editingField !== 'phone'}
                    className={editingField !== 'phone' ? "bg-gray-50 flex-grow" : "flex-grow"}
                  />
                  {editingField === 'phone' ? (
                    <div className="flex ml-2 gap-3">
                      <Button variant="ghost" className="hover:bg-neutral-300 bg-neutral-200 text-[15px] text-gray-700 !h-[40px] !w-[80px]" onClick={cancelEditing}>
                        Cancel
                      </Button>
                      <Button className="hover:bg-blue-500 bg-[#0496FF] text-[15px] text-white !h-[40px] !w-[80px]" onClick={() => updateField('phone')}>
                        Save
                        <FiCheck className="ml-1 !w-[15px] !h-[15px]" />
                      </Button>
                    </div>
                  ) : (
                    <Button
                      variant="outline"
                      className="ml-2 w-[40px] h-[40px"
                      onClick={() => startEditing('phone')}
                    >
                      <FiEdit className="!w-[20px] !h-[20px]" />
                    </Button>
                  )}
                </div>
              </div>

              <div className="mb-5">
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Full Name
                </label>
                <div className="flex items-center">
                  <div className="flex-grow grid grid-cols-2 gap-2">
                    <Input
                      name="firstName"
                      value={formData.firstName}
                      onChange={handleChange}
                      readOnly={editingField !== 'name'}
                      className={editingField !== 'name' ? "bg-gray-50" : ""}
                      placeholder="First Name"
                    />
                    <Input
                      name="lastName"
                      value={formData.lastName}
                      onChange={handleChange}
                      readOnly={editingField !== 'name'}
                      className={editingField !== 'name' ? "bg-gray-50" : ""}
                      placeholder="Last Name"
                    />
                  </div>
                  {editingField === 'name' ? (
                    <div className="flex ml-2 gap-3">
                      <Button variant="ghost" className="hover:bg-neutral-300 bg-neutral-200 text-[15px] text-gray-700 !h-[40px] !w-[80px]" onClick={cancelEditing}>
                        Cancel
                      </Button>
                      <Button className="hover:bg-blue-500 bg-[#0496FF] text-[15px] text-white !h-[40px] !w-[80px]" onClick={() => updateField('name')}>
                        Save
                        <FiCheck className="ml-1 !w-[15px] !h-[15px]" />
                      </Button>
                    </div>
                  ) : (
                    <Button
                      variant="outline"
                      className="ml-2 w-[40px] h-[40px]"
                      onClick={() => startEditing('name')}
                    >
                      <FiEdit className="!w-[20px] !h-[20px]" />
                    </Button>
                  )}
                </div>
              </div>

              <div className="mb-5">
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Date of Birth
                </label>
                <div className="flex items-center">
                  <Input
                    name="dateOfBirth"
                    type="date"
                    value={formData.dateOfBirth}
                    onChange={handleChange}
                    readOnly={editingField !== 'dateOfBirth'}
                    className={editingField !== 'dateOfBirth' ? "bg-gray-50 flex-grow text-neutral-500" : "flex-grow text-[15px] font-sans"}
                  />
                  {editingField === 'dateOfBirth' ? (
                    <div className="flex ml-2 gap-2">
                      <Button variant="ghost" className="hover:bg-neutral-300 bg-neutral-200 text-[15px] text-gray-700 !h-[40px] !w-[80px]" onClick={cancelEditing}>
                        Cancel
                      </Button>
                      <Button className="hover:bg-blue-500 bg-[#0496FF] text-[15px] text-white !h-[40px] !w-[80px]" onClick={() => updateField('dateOfBirth')}>
                        Save
                        <FiCheck className="ml-1 !w-[15px] !h-[15px]" />
                      </Button>
                    </div>
                  ) : (
                    <Button
                      variant="outline"
                      className="ml-2 w-[40px] h-[40px]"
                      onClick={() => startEditing('dateOfBirth')}
                    >
                      <FiEdit className="!w-[20px] !h-[20px]" />
                    </Button>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Country/Region
                  </label>
                  <Input
                    value={formData.country}
                    readOnly
                    className="bg-gray-50 text-neutral-500"
                  />
                </div>
                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Status
                  </label>
                  <Input
                    value="Active"
                    readOnly
                    className="bg-gray-50 text-neutral-500  "
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm p-6 mb-5">
          <h2 className="text-[18px] font-sans font-medium mb-4">BILLING ADDRESS</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                First Name
              </label>
              <Input
                value={formData.firstName}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>
            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Last Name
              </label>
              <Input
                value={formData.lastName}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Address
              </label>
              <Input
                value={formData.city + ', ' + formData.country}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Country
              </label>
              <Input
                value={formData.country}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Region/State
              </label>
              <Input
                value={formData.region}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                City
              </label>
              <Input
                value={formData.city}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Zip Code
              </label>
              <Input
                value={formData.zipCode}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Email
              </label>
              <Input
                value={formData.email}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>

            <div>
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Phone Number
              </label>
              <Input
                value={formData.phone}
                readOnly
                className="bg-gray-50 text-neutral-500"
              />
            </div>
          </div>
          <div className="mt-5">
            <Button
              variant="outline" 
              className="hover:bg-orange-600 hover:text-white text-[15px] bg-orange-500 font-sans font-medium text-white w-[110px] h-[40px] p-4"
            >
              Edit Address
            </Button>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm p-6">
          <h2 className="text-[18px] font-sans font-medium mb-4">CHANGE PASSWORD</h2>
          <div className="space-y-4">
            <div className="mb-5">
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Current Password
              </label>
              <Input
                type="password"
                placeholder="Enter current password"
                className="bg-gray-50 text-neutral-500"
              />
            </div>
            <div className="!mb-5">
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                New Password
              </label>
              <Input
                type="password"
                placeholder="Min. 8 characters"
                className="bg-gray-50 text-neutral-500"
              />
            </div>
            <div className="!mb-5">
              <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                Confirm Password
              </label>
              <Input
                type="password"
                placeholder="Confirm new password"
                className="bg-gray-50 text-neutral-500"
              />
            </div>
            <Button className="bg-orange-500 hover:bg-orange-600 text-[15px] font-sans font-medium text-white w-[150px] h-[50px] !p-5">
              Change Password
            </Button>
          </div>
        </div>
      </div>
    </UserLayout>
  );
};

export default Profile;