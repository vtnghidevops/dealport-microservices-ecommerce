// components/user/AddressForm.tsx
import React, { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { UserAddress } from '@/types/user.model';

type AddressFormData = Omit<UserAddress, 'id'> & { id?: string };

interface AddressFormProps {
  address?: UserAddress;
  onSubmit: (address: AddressFormData) => void;
  onCancel: () => void;
}

const AddressForm: React.FC<AddressFormProps> = ({
  address,
  onSubmit,
  onCancel
}) => {
  const [formData, setFormData] = useState<AddressFormData>(
    address || {
      name: '',
      phone: '',
      street: '',
      city: '',
      state: '',
      zipCode: '',
      country: 'Vietnam',
      isDefault: false
    }
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 mb-4">
        <div className='mb-6'>
          <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">Full Name</label>
          <Input
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            className="text-[15px] font-sans"
          />
        </div>

        <div className='mb-5'>
          <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">Phone Number</label>
          <Input
            name="phone"
            value={formData.phone}
            onChange={handleChange}
            required
            className="text-[15px] font-sans"
          />
        </div>
      </div>

      <div className="!mb-5">
        <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">Street Address</label>
        <Input
          name="street"
          value={formData.street}
          onChange={handleChange}
          required
          className="text-[15px] font-sans"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 !mb-5">
        <div>
          <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">City</label>
          <Input
            name="city"
            value={formData.city}
            onChange={handleChange}
            required
            className="text-[15px] font-sans"
          />
        </div>

        <div>
          <label className="block text-[15px] font-sans font-medium text-gray-700 mb-6">State/Province</label>
          <Input
            name="state"
            value={formData.state}
            onChange={handleChange}
            required
            className="text-[15px] font-sans"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 mb-4">
        <div>
          <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">ZIP/Postal Code</label>
          <Input
            name="zipCode"
            value={formData.zipCode}
            onChange={handleChange}
            required
            className="text-[15px] font-sans"
          />
        </div>

        <div>
          <label className="block text-[15px] font-sans font-medium text-gray-700 mb-1">Country</label>
          <Select
            value={formData.country}
            onValueChange={(value: string) => setFormData(prev => ({ ...prev, country: value }))}
          >
            <SelectTrigger className="text-[15px] font-sans">
              <SelectValue placeholder="Select a country" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Vietnam">Vietnam</SelectItem>
              <SelectItem value="United States">United States</SelectItem>
              <SelectItem value="China">China</SelectItem>
              <SelectItem value="Japan">Japan</SelectItem>
              <SelectItem value="South Korea">South Korea</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="flex items-center space-x-2 mb-6 !mt-5">
        <Checkbox
          id="default-address"
          checked={formData.isDefault}
          onCheckedChange={(checked) =>
            setFormData(prev => ({ ...prev, isDefault: checked === true }))
          }
          className='!w-[15px] !h-[15px] data-[state=checked]:bg-[#0496FF] data-[state=checked]:text-white border-gray-400'
        />
        <label htmlFor="default-address" className="text-[15px] font-sans font-medium text-gray-700">
          Set as default address
        </label>
      </div>

      <div className="flex justify-end gap-3">
        <Button
          type="button"
          variant="ghost"
          onClick={onCancel}
          className="hover:bg-neutral-300 bg-neutral-200 text-[15px] font-sans font-medium text-gray-700 !h-[40px] !w-[80px]"
        >
          Cancel
        </Button>
        <Button
          type="submit"
          className="hover:bg-blue-500 bg-[#0496FF] text-[15px] font-sans font-medium text-white !h-[40px] !w-[100px]"
        >
          {address ? 'Update' : 'Save'}

        </Button>
      </div>
    </form>
  );
};

export default AddressForm;