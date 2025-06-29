// pages/user/Addresses.tsx
import React, { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import UserLayout from '@/components/layouts/UserLayout';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { FiEdit, FiTrash2 } from 'react-icons/fi';
import AddressForm from '@/components/user/AddressForm';
import { UserAddress } from '@/types/user.model';

// Mock address data
const mockAddresses: UserAddress[] = [
  {
    id: 'addr-1',
    isDefault: true,
    name: 'John Doe',
    phone: '0123456789',
    street: '123 Main St',
    city: 'New York',
    state: 'NY',
    zipCode: '10001',
    country: 'United States'
  },
  {
    id: 'addr-2',
    isDefault: false,
    name: 'John Doe',
    phone: '9876543210',
    street: '456 Park Ave',
    city: 'Los Angeles',
    state: 'CA',
    zipCode: '90001',
    country: 'United States'
  }
];

const Addresses: React.FC = () => {
  const { authState } = useAuth();
  const [addresses, setAddresses] = useState<UserAddress[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showAddForm, setShowAddForm] = useState(false);
  const [editingAddress, setEditingAddress] = useState<UserAddress | null>(null);

  useEffect(() => {
    // In a real app, fetch addresses from API
    const fetchAddresses = async () => {
      try {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 800));
        setAddresses(mockAddresses);
      } catch (error) {
        console.error('Error fetching addresses:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchAddresses();
  }, [authState.user?.id]);

  const handleAddAddress = (address: Omit<UserAddress, 'id' | 'isDefault'>) => {
    // In a real app, would save to API
    const newAddress: UserAddress = {
      ...address,
      id: `addr-${Date.now()}`,
      isDefault: addresses.length === 0 // Make default if it's the first address
    };

    setAddresses([...addresses, newAddress]);
    setShowAddForm(false);
  };

  const handleUpdateAddress = (updatedAddress: UserAddress) => {
    // In a real app, would update via API
    setAddresses(addresses.map(addr =>
      addr.id === updatedAddress.id ? updatedAddress : addr
    ));
    setEditingAddress(null);
  };

  const handleRemoveAddress = (id: string) => {
    // In a real app, would delete via API
    setAddresses(addresses.filter(addr => addr.id !== id));
  };

  const handleSetDefault = (id: string) => {
    // In a real app, would update via API
    setAddresses(addresses.map(addr => ({
      ...addr,
      isDefault: addr.id === id
    })));
  };

  if (isLoading) {
    return (
      <UserLayout>
        <div className="max-w-4xl mx-auto p-6 flex justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      </UserLayout>
    );
  }

  return (
    <UserLayout>
      <div className="max-w-4xl mx-auto px-6">
        <div className="flex justify-between items-center mb-5">
          <h1 className="text-[20px] font-medium font-sans">My Addresses</h1>
          <Button
            onClick={() => setShowAddForm(true)}
            className="flex items-center text-[15px] font-bold font-sans gap-2 !p-4 !w-[170px] !h-[50px] bg-[#0496FF] text-white hover:bg-blue-600"
          >
            Add New Address
          </Button>
        </div>

        {showAddForm && (
          <Card className="mb-5 bg-white rounded-lg shadow-sm !p-5 ">
            <CardHeader className="px-0 pt-0">
              <CardTitle className="text-[18px] font-sans font-medium ">Add New Address</CardTitle>
            </CardHeader>
            <CardContent className="px-0 pb-0">
              <AddressForm
                onSubmit={handleAddAddress}
                onCancel={() => setShowAddForm(false)}
              />
            </CardContent>
          </Card>
        )}

        {editingAddress && (
          <Card className="mb-5 bg-white rounded-lg shadow-sm p-5">
            <CardHeader className="px-0 pt-0">
              <CardTitle className="text-[18px] font-sans font-medium">Edit Address</CardTitle>
            </CardHeader>
            <CardContent className="px-4 pb-0">
              <AddressForm
                address={editingAddress}
                onSubmit={(address) => handleUpdateAddress({ ...address, id: editingAddress.id })}
                onCancel={() => setEditingAddress(null)}
              />
            </CardContent>
          </Card>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {addresses.map(address => (
            <Card key={address.id} className={address.isDefault ? "border-blue-500 p-[10px] bg-white rounded-lg shadow-sm" : "p-[10px] bg-white rounded-lg shadow-sm"}>
              <CardHeader className="pb-2 px-0 pt-0">
                <div className="flex justify-between items-center">
                  <CardTitle className="text-[18px] font-sans font-medium">{address.name}</CardTitle>
                  {address.isDefault && (
                    <span className="bg-blue-100 text-blue-800 text-xs font-sans px-2 py-1 rounded">
                      Default
                    </span>
                  )}
                </div>
              </CardHeader>
              <CardContent className="pb-3 px-0">
                <p className="text-[15px] font-sans text-gray-700">{address.phone}</p>
                <p className="text-[15px] font-sans text-gray-700">{address.street}</p>
                <p className="text-[15px] font-sans text-gray-700">{`${address.city}, ${address.state} ${address.zipCode}`}</p>
                <p className="text-[15px] font-sans text-gray-700">{address.country}</p>
              </CardContent>
              <CardFooter className="flex justify-between items-center pt-2 border-t min-h-[50px] px-0">
                <div className="flex gap-2 items-center">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-[15px] font-sans"
                    onClick={() => setEditingAddress(address)}
                  >
                    <FiEdit className="mr-1 !w-[20px] !h-[20px]" /> Edit
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-[15px] font-sans"
                    onClick={() => handleRemoveAddress(address.id)}
                  >
                    <FiTrash2 className="mr-1 !w-[20px] !h-[20px]" /> Delete
                  </Button>
                </div>
                {!address.isDefault && (
                  <Button
                    variant="outline"
                    className='w-[130px] h-[35px] px-4 bg-[#0496FF] text-white hover:bg-blue-600 text-[15px] font-sans'
                    onClick={() => handleSetDefault(address.id)}
                  >
                    Set as Default
                  </Button>
                )}
              </CardFooter>
            </Card>
          ))}
        </div>

        {addresses.length === 0 && !showAddForm && (
          <div className="text-center p-8 bg-white rounded-lg shadow-sm">
            <h3 className="text-[18px] font-sans font-medium mb-2">No addresses saved</h3>
            <p className="text-[15px] font-sans text-gray-500 mb-4">Add your first shipping address to speed up checkout.</p>
            <Button
              onClick={() => setShowAddForm(true)}
              className="bg-[#0496FF] text-white hover:bg-blue-600 text-[15px] font-sans font-medium"
            >
              Add Address
            </Button>
          </div>
        )}
      </div>
    </UserLayout>
  );
};

export default Addresses;