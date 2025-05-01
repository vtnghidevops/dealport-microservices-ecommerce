import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Calendar } from '@/components/ui/calendar';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { format } from 'date-fns';
import { CalendarIcon, ArrowLeft, Check } from 'lucide-react';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import couponService from '@/services/admin/coupon.service';
import { useToast } from '@/hooks/use-toast';
import AdminHeader from '@/components/admin/layout/AdminHeader';

interface CouponFormData {
  code: string;
  discountType: 'percentage' | 'fixed';
  discount: number;
  minOrderAmount: number;
  maxDiscount?: number;
  maxUsage: number;
  validFrom: Date;
  validTo: Date;
  isActive: boolean;
  description?: string;
}

const CouponCreate: React.FC = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [formData, setFormData] = useState<CouponFormData>({
    code: '',
    discountType: 'percentage',
    discount: 10,
    minOrderAmount: 0,
    maxDiscount: undefined,
    maxUsage: 100,
    validFrom: new Date(),
    validTo: new Date(new Date().setMonth(new Date().getMonth() + 1)),
    isActive: true,
    description: '',
  });

  // Style cơ bản cho input
  const inputBaseClass = "w-full border rounded-md h-[48px] px-3 focus:outline-none focus:border-ocean-green bg-white";
  const textareaBaseClass = "w-full border rounded-md px-3 py-2 focus:outline-none focus:border-ocean-green bg-white min-h-[100px]";

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;

    // Special handling for numeric fields
    if (name === 'discount' || name === 'minOrderAmount' || name === 'maxUsage' || name === 'maxDiscount') {
      // Allow empty value (will be converted later)
      if (value === '') {
        const newValue = name === 'maxDiscount' ? undefined : 0;
        setFormData({
          ...formData,
          [name]: newValue
        });
        return;
      }

      // Accept only numeric input with optional decimal point
      const numericRegex = /^[0-9]*\.?[0-9]*$/;
      if (numericRegex.test(value)) {
        // Convert to number for internal state
        setFormData({
          ...formData,
          [name]: value === '' ? (name === 'maxDiscount' ? undefined : 0) : Number(value)
        });
      }
      // If not numeric, ignore the change
      return;
    }

    // For non-numeric fields
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSelectChange = (name: string, value: string) => {
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleDateChange = (name: string, date: Date | undefined) => {
    if (date) {
      setFormData({
        ...formData,
        [name]: date,
      });
    }
  };

  const handleCheckboxChange = (name: string, checked: boolean) => {
    setFormData({
      ...formData,
      [name]: checked,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      // Form validation
      if (!formData.code) {
        throw new Error('Coupon code cannot be empty');
      }

      if (formData.discount <= 0) {
        throw new Error('Discount value must be greater than 0');
      }

      if (formData.discountType === 'percentage' && formData.discount > 100) {
        throw new Error('Percentage discount cannot exceed 100%');
      }

      if (formData.validFrom > formData.validTo) {
        throw new Error('Start date must be before end date');
      }

      // Format data for API
      const couponData = {
        code: formData.code,
        discountType: formData.discountType,
        discount: formData.discount,
        minOrderAmount: formData.minOrderAmount,
        maxDiscount: formData.maxDiscount,
        maxUsage: formData.maxUsage,
        validFrom: format(formData.validFrom, 'yyyy-MM-dd'),
        validTo: format(formData.validTo, 'yyyy-MM-dd'),
        isActive: formData.isActive,
        description: formData.description || undefined,
      };

      // Call API to create coupon
      await couponService.createCoupon(couponData);

      toast({
        variant: "success",
        title: 'Success',
        description: 'New coupon created successfully',
      });

      // Redirect to coupon list
      navigate('/admin/coupons');
    } catch (error) {
      console.error('Error creating coupon:', error);
      toast({
        title: 'Error',
        description: error instanceof Error ? error.message : 'Unable to create coupon',
        variant: 'destructive',
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Create New Coupon" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            <div className="flex justify-between items-center mb-4 h-[3rem] w-[11rem] p-3 border-2 shadow-sm hover:bg-aqua-spring duration-200 transfrom border-neutral-100 bg-white rounded-md">
              <Button
                variant="ghost"
                className="flex items-center text-gray-600 text-[15px] font-medium"
                onClick={() => navigate('/admin/coupons')}
              >
                <ArrowLeft className="mr-2 h-5 w-5" />
                Back to Coupons
              </Button>
            </div>

            <div className="!p-5 w-[1116px] bg-white rounded-lg shadow mb-6 drop-shadow filter">
              <div className="mb-5">
                <h2 className="text-xl font-semibold">Coupon Information</h2>
                <p className="text-gray-500 text-[15px] mt-1">Create a new coupon with specific discount conditions</p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                  {/* Coupon code */}
                  <div className="space-y-2">
                    <Label htmlFor="code" className="font-medium text-cyprus">Coupon Code</Label>
                    <Input
                      id="code"
                      name="code"
                      placeholder="SUMMER2025"
                      className={inputBaseClass}
                      value={formData.code}
                      onChange={handleChange}
                      required
                    />
                  </div>

                  {/* Discount type */}
                  <div className="space-y-2">
                    <Label htmlFor="discountType" className="font-medium text-cyprus">Discount Type</Label>
                    <div className="grid grid-cols-2 gap-3 h-[48px]">
                      <button
                        type="button"
                        className={`flex items-center justify-center rounded-md border-2 ${formData.discountType === 'percentage'
                          ? 'border-green-500 bg-green-50 text-green-700'
                          : 'border-gray-300 bg-white text-gray-700 hover:border-gray-400'
                          }`}
                        onClick={() => handleSelectChange('discountType', 'percentage')}
                      >
                        <span className="font-medium">Percentage (%)</span>
                      </button>
                      <button
                        type="button"
                        className={`flex items-center justify-center rounded-md border-2 ${formData.discountType === 'fixed'
                          ? 'border-green-500 bg-green-50 text-green-700'
                          : 'border-gray-300 bg-white text-gray-700 hover:border-gray-400'
                          }`}
                        onClick={() => handleSelectChange('discountType', 'fixed')}
                      >
                        <span className="font-medium">Fixed Amount</span>
                      </button>
                    </div>
                  </div>

                  {/* Discount value */}
                  <div className="space-y-2">
                    <Label htmlFor="discount" className="font-medium text-cyprus">
                      {formData.discountType === 'percentage' ? 'Discount Percentage (%)' : 'Discount Amount'}
                    </Label>
                    <Input
                      id="discount"
                      name="discount"
                      type="text"
                      className={inputBaseClass}
                      value={formData.discount}
                      onChange={handleChange}
                      placeholder={formData.discountType === 'percentage' ? "e.g. 10" : "e.g. 1000"}
                      required
                    />
                    {formData.discountType === 'percentage' && formData.discount > 100 && (
                      <p className="text-red-500 text-sm mt-1">Percentage discount cannot exceed 100%</p>
                    )}
                  </div>

                  {/* Max discount (only for percentage) */}
                  {formData.discountType === 'percentage' && (
                    <div className="space-y-2">
                      <Label htmlFor="maxDiscount" className="font-medium text-cyprus">Max Discount (leave empty for no limit)</Label>
                      <Input
                        id="maxDiscount"
                        name="maxDiscount"
                        type="text"
                        className={inputBaseClass}
                        value={formData.maxDiscount || ''}
                        onChange={handleChange}
                        placeholder="e.g. 50000"
                      />
                    </div>
                  )}

                  {/* Minimum order amount */}
                  <div className="space-y-2">
                    <Label htmlFor="minOrderAmount" className="font-medium text-cyprus">Minimum Order Amount</Label>
                    <Input
                      id="minOrderAmount"
                      name="minOrderAmount"
                      type="text"
                      className={inputBaseClass}
                      value={formData.minOrderAmount}
                      onChange={handleChange}
                      placeholder="e.g. 0"
                      required
                    />
                  </div>

                  {/* Maximum usage count */}
                  <div className="space-y-2">
                    <Label htmlFor="maxUsage" className="font-medium text-cyprus">Maximum Usage Count</Label>
                    <Input
                      id="maxUsage"
                      name="maxUsage"
                      type="text"
                      className={inputBaseClass}
                      value={formData.maxUsage}
                      onChange={handleChange}
                      placeholder="e.g. 100"
                      required
                    />
                  </div>

                  {/* Start date */}
                  <div className="space-y-2">
                    <Label htmlFor="validFrom" className="font-medium text-cyprus">Start Date</Label>
                    <div className="relative">
                      <input
                        type="date"
                        id="validFrom"
                        className={inputBaseClass}
                        value={format(formData.validFrom, 'yyyy-MM-dd')}
                        onChange={(e) => {
                          if (e.target.value) {
                            handleDateChange('validFrom', new Date(e.target.value));
                          }
                        }}
                      />
                    </div>
                  </div>

                  {/* End date */}
                  <div className="space-y-2">
                    <Label htmlFor="validTo" className="font-medium text-cyprus">End Date</Label>
                    <div className="relative">
                      <input
                        type="date"
                        id="validTo"
                        className={inputBaseClass}
                        value={format(formData.validTo, 'yyyy-MM-dd')}
                        onChange={(e) => {
                          if (e.target.value) {
                            handleDateChange('validTo', new Date(e.target.value));
                          }
                        }}
                      />
                    </div>
                  </div>

                  {/* Status */}
                  <div className="flex items-center space-y-0">
                    <div className="flex items-center space-x-2">
                      <div className="relative" onClick={() => handleCheckboxChange('isActive', !formData.isActive)}>
                        <div className={`h-5 w-5 rounded border-2 cursor-pointer transition-colors flex items-center justify-center ${formData.isActive ? 'bg-aqua-spring border-aqua-spring' : 'bg-white border-gray-300'}`}>
                          {formData.isActive && (
                            <Check className="h-3.5 w-3.5 text-black" strokeWidth={3} />
                          )}
                        </div>
                        <input
                          type="checkbox"
                          id="isActive"
                          className="sr-only"
                          checked={formData.isActive}
                          onChange={(e) => handleCheckboxChange('isActive', e.target.checked)}
                        />
                      </div>
                      <Label
                        htmlFor="isActive"
                        className="font-medium cursor-pointer"
                        onClick={() => handleCheckboxChange('isActive', !formData.isActive)}
                      >
                        Enable this coupon
                      </Label>
                    </div>
                  </div>
                </div>

                {/* Description */}
                <div className="space-y-2 mt-2">
                  <Label htmlFor="description" className="font-medium text-cyprus">Description (optional)</Label>
                  <Textarea
                    id="description"
                    name="description"
                    placeholder="Describe this coupon"
                    className={textareaBaseClass}
                    value={formData.description || ''}
                    onChange={handleChange}
                    rows={3}
                  />
                </div>

                <div className="flex justify-end space-x-4 pt-4 mt-4 border-t border-gray-100">
                  <Button
                    type="button"
                    variant="outline"
                    className="h-10 px-5"
                    onClick={() => navigate('/admin/coupons')}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    disabled={isSubmitting}
                    className="h-10 px-5 bg-ocean-green text-white hover:bg-green-700"
                  >
                    {isSubmitting ? 'Creating...' : 'Create Coupon'}
                  </Button>
                </div>
              </form>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default CouponCreate; 