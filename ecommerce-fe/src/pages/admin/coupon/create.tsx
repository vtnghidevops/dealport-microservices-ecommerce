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
import { CalendarIcon } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import couponService, { Coupon } from '@/services/user/coupon.service';
import { useToast } from '@/hooks/use-toast';

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

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: name === 'discount' || name === 'minOrderAmount' || name === 'maxUsage' || name === 'maxDiscount'
        ? Number(value)
        : value,
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
      // Kiểm tra form validation
      if (!formData.code) {
        throw new Error('Mã coupon không được để trống');
      }

      if (formData.discount <= 0) {
        throw new Error('Giá trị giảm giá phải lớn hơn 0');
      }

      if (formData.discountType === 'percentage' && formData.discount > 100) {
        throw new Error('Tỷ lệ giảm giá không thể lớn hơn 100%');
      }

      if (formData.validFrom > formData.validTo) {
        throw new Error('Ngày bắt đầu phải trước ngày kết thúc');
      }

      // Format dữ liệu gửi lên API
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

      // Gọi API tạo coupon
      await couponService.createCoupon(couponData);

      toast({
        title: 'Thành công',
        description: 'Tạo mã giảm giá mới thành công',
      });

      // Chuyển hướng về trang danh sách coupon
      navigate('/admin/coupons');
    } catch (error) {
      console.error('Error creating coupon:', error);
      toast({
        title: 'Lỗi',
        description: error instanceof Error ? error.message : 'Không thể tạo mã giảm giá',
        variant: 'destructive',
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="container mx-auto p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Tạo Mã Giảm Giá Mới</h1>
        <Button variant="outline" onClick={() => navigate('/admin/coupons')}>
          Quay lại
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Thông tin mã giảm giá</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Mã coupon */}
              <div className="space-y-2">
                <Label htmlFor="code">Mã coupon</Label>
                <Input
                  id="code"
                  name="code"
                  placeholder="SUMMER2025"
                  value={formData.code}
                  onChange={handleChange}
                  required
                />
              </div>

              {/* Loại giảm giá */}
              <div className="space-y-2">
                <Label htmlFor="discountType">Loại giảm giá</Label>
                <Select
                  value={formData.discountType}
                  onValueChange={(value) => handleSelectChange('discountType', value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Chọn loại giảm giá" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="percentage">Phần trăm (%)</SelectItem>
                    <SelectItem value="fixed">Số tiền cố định</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Giá trị giảm */}
              <div className="space-y-2">
                <Label htmlFor="discount">
                  {formData.discountType === 'percentage' ? 'Phần trăm giảm (%)' : 'Số tiền giảm'}
                </Label>
                <Input
                  id="discount"
                  name="discount"
                  type="number"
                  min={0}
                  max={formData.discountType === 'percentage' ? 100 : undefined}
                  step={formData.discountType === 'percentage' ? 1 : 1000}
                  value={formData.discount}
                  onChange={handleChange}
                  required
                />
              </div>

              {/* Giảm tối đa (chỉ cho % giảm) */}
              {formData.discountType === 'percentage' && (
                <div className="space-y-2">
                  <Label htmlFor="maxDiscount">Giảm tối đa (để trống nếu không giới hạn)</Label>
                  <Input
                    id="maxDiscount"
                    name="maxDiscount"
                    type="number"
                    min={0}
                    step={1000}
                    value={formData.maxDiscount || ''}
                    onChange={handleChange}
                  />
                </div>
              )}

              {/* Đơn hàng tối thiểu */}
              <div className="space-y-2">
                <Label htmlFor="minOrderAmount">Đơn hàng tối thiểu</Label>
                <Input
                  id="minOrderAmount"
                  name="minOrderAmount"
                  type="number"
                  min={0}
                  step={1000}
                  value={formData.minOrderAmount}
                  onChange={handleChange}
                  required
                />
              </div>

              {/* Số lượt sử dụng tối đa */}
              <div className="space-y-2">
                <Label htmlFor="maxUsage">Số lượt sử dụng tối đa</Label>
                <Input
                  id="maxUsage"
                  name="maxUsage"
                  type="number"
                  min={1}
                  value={formData.maxUsage}
                  onChange={handleChange}
                  required
                />
              </div>

              {/* Ngày bắt đầu */}
              <div className="space-y-2">
                <Label>Ngày bắt đầu</Label>
                <Popover>
                  <PopoverTrigger asChild>
                    <Button
                      variant="outline"
                      className={cn(
                        "w-full justify-start text-left font-normal",
                        !formData.validFrom && "text-muted-foreground"
                      )}
                    >
                      <CalendarIcon className="mr-2 h-4 w-4" />
                      {formData.validFrom ? (
                        format(formData.validFrom, 'dd/MM/yyyy')
                      ) : (
                        <span>Chọn ngày</span>
                      )}
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0">
                    <Calendar
                      mode="single"
                      selected={formData.validFrom}
                      onSelect={(date) => handleDateChange('validFrom', date)}
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
              </div>

              {/* Ngày kết thúc */}
              <div className="space-y-2">
                <Label>Ngày kết thúc</Label>
                <Popover>
                  <PopoverTrigger asChild>
                    <Button
                      variant="outline"
                      className={cn(
                        "w-full justify-start text-left font-normal",
                        !formData.validTo && "text-muted-foreground"
                      )}
                    >
                      <CalendarIcon className="mr-2 h-4 w-4" />
                      {formData.validTo ? (
                        format(formData.validTo, 'dd/MM/yyyy')
                      ) : (
                        <span>Chọn ngày</span>
                      )}
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0">
                    <Calendar
                      mode="single"
                      selected={formData.validTo}
                      onSelect={(date) => handleDateChange('validTo', date)}
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
              </div>

              {/* Trạng thái */}
              <div className="space-y-2 flex items-center">
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="isActive"
                    checked={formData.isActive}
                    onCheckedChange={(checked) =>
                      handleCheckboxChange('isActive', checked as boolean)
                    }
                  />
                  <Label htmlFor="isActive">Kích hoạt mã giảm giá</Label>
                </div>
              </div>
            </div>

            {/* Mô tả */}
            <div className="space-y-2">
              <Label htmlFor="description">Mô tả (tùy chọn)</Label>
              <Textarea
                id="description"
                name="description"
                placeholder="Mô tả về mã giảm giá này"
                value={formData.description || ''}
                onChange={handleChange}
                rows={3}
              />
            </div>

            <div className="flex justify-end space-x-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => navigate('/admin/coupons')}
              >
                Hủy
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? 'Đang tạo...' : 'Tạo mã giảm giá'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default CouponCreate; 