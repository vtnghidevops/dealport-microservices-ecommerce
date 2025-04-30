// pages/user/Profile.tsx
import React, { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { useToast } from '@/hooks/use-toast';
import UserLayout from '../../components/layouts/UserLayout';
import UserAvatar from '../../components/user/UserAvatar';
import { FiEdit, FiCheck } from 'react-icons/fi';
import { User, UserProfile, UserAddress } from '@/types/user.model';
import userService from '@/services/user/user.service';

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
    displayName: authState.user?.profile?.firstName ?
      `${authState.user?.profile?.firstName} ${authState.user?.profile?.lastName || ''}`.trim() :
      authState.user?.username || '',
    username: authState.user?.profile?.lastName && authState.user?.profile?.firstName ?
      `${authState.user?.profile?.lastName}${authState.user?.profile?.firstName}` :
      authState.user?.username || '',
    address: '',
    country: 'Vietnam',
    region: '',
    city: '',
    zipCode: '',
    addressId: '',
  });

  // Update formData when authState.user changes
  useEffect(() => {
    if (authState.user) {
      setFormData(prev => ({
        ...prev,
        firstName: authState.user?.profile?.firstName || '',
        lastName: authState.user?.profile?.lastName || '',
        email: authState.user?.email || '',
        phone: authState.user?.profile?.phone || '',
        dateOfBirth: authState.user?.profile?.dateOfBirth || '',
        displayName: authState.user?.profile?.firstName ?
          `${authState.user?.profile?.firstName} ${authState.user?.profile?.lastName || ''}`.trim() :
          authState.user?.username || '',
        username: authState.user?.profile?.lastName && authState.user?.profile?.firstName ?
          `${authState.user?.profile?.lastName}${authState.user?.profile?.firstName}` :
          authState.user?.username || '',
      }));
    }
  }, [authState.user]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const startEditing = (fieldName: string) => {
    setEditingField(fieldName);

    // If starting to edit address, populate with current address data if available
    if (fieldName === 'address' && authState.user?.addresses && authState.user.addresses.length > 0) {
      const defaultAddress = authState.user.addresses.find(addr => addr.isDefault) || authState.user.addresses[0];
      setFormData(prev => ({
        ...prev,
        address: defaultAddress.street || '',
        city: defaultAddress.city || '',
        region: defaultAddress.state || '',
        zipCode: defaultAddress.zipCode || '',
        country: defaultAddress.country || 'Vietnam',
        addressId: defaultAddress.id || '',
      }));
    }
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
      if (fieldName === 'name') {
        const updatedProfile: Partial<UserProfile> = {
          firstName: formData.firstName,
          lastName: formData.lastName
        };

        const userUpdate: Partial<User> = {
          profile: {
            ...(authState.user?.profile || {}),
            ...updatedProfile
          } as UserProfile
        };

        await updateProfile(userUpdate);

        // Also update username to match lastName + firstName pattern
        setFormData(prev => ({
          ...prev,
          username: `${formData.lastName}${formData.firstName}`,
          displayName: `${formData.firstName} ${formData.lastName}`.trim()
        }));
      } else if (fieldName === 'phone') {
        const updatedProfile: Partial<UserProfile> = {
          phone: formData.phone
        };

        const userUpdate: Partial<User> = {
          profile: {
            ...(authState.user?.profile || {}),
            ...updatedProfile
          } as UserProfile
        };

        await updateProfile(userUpdate);
      } else if (fieldName === 'dateOfBirth') {
        const updatedProfile: Partial<UserProfile> = {
          dateOfBirth: formData.dateOfBirth
        };

        const userUpdate: Partial<User> = {
          profile: {
            ...(authState.user?.profile || {}),
            ...updatedProfile
          } as UserProfile
        };

        await updateProfile(userUpdate);
      } else if (fieldName === 'address') {
        // Basic validation
        if (!formData.address || !formData.city || !formData.country) {
          toast({
            title: 'Validation Error',
            description: 'Address, city and country are required fields',
            variant: 'error',
          });
          return;
        }

        // Format address data for API call
        const addressData = {
          name: `${formData.firstName} ${formData.lastName}`.trim(),
          phone: formData.phone || '',
          line1: formData.address,
          city: formData.city,
          state: formData.region,
          postal_code: formData.zipCode,
          country: formData.country,
          is_default: true,
          address_type: 'billing'
        };

        try {
          let updatedAddress;

          // Check if we're updating an existing address or creating a new one
          if (formData.addressId) {
            // Update existing address
            updatedAddress = await userService.updateAddress(formData.addressId, addressData);
            console.log('Address updated:', updatedAddress);
          } else {
            // Create new address
            updatedAddress = await userService.createAddress(addressData);
            console.log('New address created:', updatedAddress);
          }

          // Refresh user data to get updated addresses
          const userId = userService.getUserIdFromStorage();
          if (userId) {
            try {
              const updatedUser = await userService.getUserById(userId);
              console.log('User data refreshed:', updatedUser);

              // If the function call didn't update the global state, we could manually update it here
              // This depends on how your auth context is set up
            } catch (refreshError) {
              console.error('Error refreshing user data:', refreshError);
            }
          }

          toast({
            title: 'Address Updated',
            description: 'Your address has been updated successfully',
            variant: 'success',
          });
        } catch (error) {
          console.error('Error updating address:', error);
          toast({
            title: 'Address Update Failed',
            description: 'Failed to update your address. Please try again.',
            variant: 'error',
          });
          return; // Don't clear editing field on error
        }
      }

      setEditingField(null);

      if (fieldName !== 'address') {
        toast({
          title: 'Profile Updated',
          description: 'Your profile has been updated successfully',
          variant: 'success',
        });
      }
    } catch (error) {
      console.error('Update error:', error);
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
                    placeholder="Your display name"
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
                    placeholder="Your username"
                  />
                  {/* <p className="text-sm text-gray-500 mt-1">
                    Username is automatically generated as lastName + firstName
                  </p> */}
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
                  placeholder="Your email address"
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
                    placeholder="Add your phone number"
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
                    type="date"
                    name="dateOfBirth"
                    value={formData.dateOfBirth}
                    onChange={handleChange}
                    readOnly={editingField !== 'dateOfBirth'}
                    className={editingField !== 'dateOfBirth' ? "bg-gray-50 flex-grow" : "flex-grow"}
                  />
                  {editingField === 'dateOfBirth' ? (
                    <div className="flex ml-2 gap-3">
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
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm p-6 mb-5">
          <h2 className="text-[18px] font-sans font-medium mb-4">BILLING ADDRESS</h2>

          {editingField === 'address' ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div>
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Address
                </label>
                <Input
                  name="address"
                  value={formData.address}
                  onChange={handleChange}
                  placeholder="Enter your street address"
                />
              </div>

              <div>
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Country
                </label>
                <Input
                  name="country"
                  value={formData.country}
                  onChange={handleChange}
                  placeholder="Country"
                />
              </div>

              <div>
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Region/State
                </label>
                <Input
                  name="region"
                  value={formData.region}
                  onChange={handleChange}
                  placeholder="Region/State"
                />
              </div>

              <div>
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  City
                </label>
                <Input
                  name="city"
                  value={formData.city}
                  onChange={handleChange}
                  placeholder="City"
                />
              </div>

              <div>
                <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                  Zip Code
                </label>
                <Input
                  name="zipCode"
                  value={formData.zipCode}
                  onChange={handleChange}
                  placeholder="Zip Code"
                />
              </div>

              <div className="md:col-span-2 mt-4 flex gap-3">
                <Button
                  variant="ghost"
                  className="p-3 hover:bg-neutral-300 bg-neutral-200 text-[15px] text-gray-700 !h-[40px]"
                  onClick={cancelEditing}
                >
                  Cancel
                </Button>
                <Button
                  className="p-3 hover:bg-blue-500 bg-[#0496FF] text-[15px] text-white !h-[40px]"
                  onClick={() => updateField('address')}
                >
                  Save Address
                  <FiCheck className="ml-1 !w-[15px] !h-[15px]" />
                </Button>
              </div>
            </div>
          ) : (
            <>
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
                    value={authState.user?.addresses && authState.user.addresses.length > 0
                      ? authState.user.addresses[0].street
                      : 'No address provided'}
                    readOnly
                    className="bg-gray-50 text-neutral-500"
                  />
                </div>

                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Country
                  </label>
                  <Input
                    value={authState.user?.addresses && authState.user.addresses.length > 0
                      ? authState.user.addresses[0].country
                      : 'Vietnam'}
                    readOnly
                    className="bg-gray-50 text-neutral-500"
                  />
                </div>

                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Region/State
                  </label>
                  <Input
                    value={authState.user?.addresses && authState.user.addresses.length > 0
                      ? authState.user.addresses[0].state
                      : ''}
                    readOnly
                    className="bg-gray-50 text-neutral-500"
                  />
                </div>

                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    City
                  </label>
                  <Input
                    value={authState.user?.addresses && authState.user.addresses.length > 0
                      ? authState.user.addresses[0].city
                      : ''}
                    readOnly
                    className="bg-gray-50 text-neutral-500"
                  />
                </div>

                <div>
                  <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">
                    Zip Code
                  </label>
                  <Input
                    value={authState.user?.addresses && authState.user.addresses.length > 0
                      ? authState.user.addresses[0].zipCode
                      : ''}
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
                  onClick={() => startEditing('address')}
                  className="hover:bg-orange-600 bg-orange-500 text-[15px] font-sans font-medium text-white w-[110px] h-[40px] p-4"
                >
                  Edit Address
                </Button>
              </div>
            </>
          )}
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