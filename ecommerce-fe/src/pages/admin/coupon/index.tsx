import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { PlusCircle, Search, MoreVertical, Edit, Trash2, Eye } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Pagination } from '@/components/ui/pagination';
import couponService, { Coupon } from '@/services/user/coupon.service';
import { useToast } from '@/hooks/use-toast';

const CouponManagement: React.FC = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [coupons, setCoupons] = useState<Coupon[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [apiError, setApiError] = useState<string | null>(null);

  // Lấy danh sách coupon từ API
  useEffect(() => {
    const fetchCoupons = async () => {
      setIsLoading(true);
      setApiError(null);

      try {
        const response = await couponService.getCoupons(currentPage, 10);
        setCoupons(response.coupons);
        setTotalPages(Math.ceil(response.total / 10));
      } catch (error) {
        console.error('Error fetching coupons:', error);
        setApiError('Không thể tải danh sách coupon. Vui lòng thử lại sau.');
        toast({
          title: 'Lỗi',
          description: 'Không thể tải danh sách coupon',
          variant: 'destructive'
        });

        // Sử dụng dữ liệu giả tạm thời khi API thất bại
        const mockCoupons: Coupon[] = [
          {
            id: '1',
            code: 'SUMMER2025',
            discount: 20,
            discountType: 'percentage',
            minOrderAmount: 100,
            maxUsage: 100,
            usageCount: 45,
            validFrom: '2025-06-01',
            validTo: '2025-08-31',
            isActive: true
          },
          {
            id: '2',
            code: 'WELCOME10',
            discount: 10,
            discountType: 'fixed',
            minOrderAmount: 50,
            maxUsage: 1000,
            usageCount: 789,
            validFrom: '2025-01-01',
            validTo: '2025-12-31',
            isActive: true
          },
          {
            id: '3',
            code: 'FLASH50',
            discount: 50,
            discountType: 'percentage',
            minOrderAmount: 200,
            maxUsage: 50,
            usageCount: 50,
            validFrom: '2025-04-15',
            validTo: '2025-04-20',
            isActive: false
          },
          {
            id: '4',
            code: 'FREESHIP',
            discount: 15,
            discountType: 'fixed',
            minOrderAmount: 80,
            maxUsage: 200,
            usageCount: 120,
            validFrom: '2025-03-01',
            validTo: '2025-09-30',
            isActive: true
          }
        ];

        setCoupons(mockCoupons);
        setTotalPages(Math.ceil(mockCoupons.length / 10));
      } finally {
        setIsLoading(false);
      }
    };

    fetchCoupons();
  }, [currentPage, toast]);

  // Lọc coupon theo term tìm kiếm
  const filteredCoupons = coupons.filter(coupon =>
    coupon.code.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Xử lý xóa coupon
  const handleDeleteCoupon = async (id: string) => {
    if (window.confirm('Bạn có chắc chắn muốn xóa coupon này không?')) {
      try {
        await couponService.deleteCoupon(id);
        setCoupons(prevCoupons => prevCoupons.filter(coupon => coupon.id !== id));

        toast({
          title: 'Thành công',
          description: 'Xóa coupon thành công',
        });
      } catch (error) {
        console.error('Error deleting coupon:', error);
        toast({
          title: 'Lỗi',
          description: 'Không thể xóa coupon. Vui lòng thử lại sau.',
          variant: 'destructive'
        });
      }
    }
  };

  // Xử lý áp dụng coupon
  const handleTestCoupon = async (code: string) => {
    try {
      await couponService.applyCoupon(code);
      toast({
        title: 'Thành công',
        description: `Đã áp dụng coupon ${code} vào giỏ hàng`,
      });
    } catch (error) {
      toast({
        title: 'Lỗi',
        description: error instanceof Error ? error.message : 'Không thể áp dụng coupon',
        variant: 'destructive'
      });
    }
  };

  // Định dạng ngày
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('vi-VN');
  };

  // Kiểm tra xem coupon đã hết hạn chưa
  const isExpired = (validTo: string) => {
    return new Date(validTo) < new Date();
  };

  // Lấy trạng thái coupon
  const getCouponStatus = (coupon: Coupon) => {
    if (!coupon.isActive) return 'inactive';
    if (isExpired(coupon.validTo)) return 'expired';
    if (coupon.usageCount >= coupon.maxUsage) return 'used';
    return 'active';
  };

  return (
    <div className="container mx-auto p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Quản lý Coupon</h1>
        <Button
          onClick={() => navigate('/admin/coupons/create')}
          className="flex items-center gap-2"
        >
          <PlusCircle size={16} />
          Tạo Coupon Mới
        </Button>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-between items-center mb-4">
          <div className="relative w-64">
            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Tìm kiếm theo mã coupon"
              className="pl-8"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
        </div>

        {isLoading ? (
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
          </div>
        ) : apiError ? (
          <div className="text-center py-8 text-red-500">
            {apiError}
          </div>
        ) : (
          <>
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Mã Coupon</TableHead>
                    <TableHead>Giảm giá</TableHead>
                    <TableHead>Đơn hàng tối thiểu</TableHead>
                    <TableHead>Sử dụng</TableHead>
                    <TableHead>Thời hạn</TableHead>
                    <TableHead>Trạng thái</TableHead>
                    <TableHead className="text-right">Thao tác</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredCoupons.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="text-center py-4">
                        Không tìm thấy coupon nào
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredCoupons.map((coupon) => (
                      <TableRow key={coupon.id}>
                        <TableCell className="font-medium">{coupon.code}</TableCell>
                        <TableCell>
                          {coupon.discountType === 'percentage'
                            ? `${coupon.discount}%`
                            : `${coupon.discount.toLocaleString('vi-VN')}đ`
                          }
                        </TableCell>
                        <TableCell>{coupon.minOrderAmount.toLocaleString('vi-VN')}đ</TableCell>
                        <TableCell>{coupon.usageCount}/{coupon.maxUsage}</TableCell>
                        <TableCell>
                          {formatDate(coupon.validFrom)} - {formatDate(coupon.validTo)}
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={
                              getCouponStatus(coupon) === 'active' ? 'default' :
                                getCouponStatus(coupon) === 'inactive' ? 'outline' :
                                  getCouponStatus(coupon) === 'expired' ? 'destructive' : 'secondary'
                            }
                          >
                            {getCouponStatus(coupon) === 'active' && 'Hoạt động'}
                            {getCouponStatus(coupon) === 'inactive' && 'Tạm ngưng'}
                            {getCouponStatus(coupon) === 'expired' && 'Hết hạn'}
                            {getCouponStatus(coupon) === 'used' && 'Hết lượt'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <MoreVertical className="h-4 w-4" />
                                <span className="sr-only">Mở menu</span>
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem onClick={() => navigate(`/admin/coupons/${coupon.id}`)}>
                                <Eye className="mr-2 h-4 w-4" />
                                <span>Chi tiết</span>
                              </DropdownMenuItem>
                              <DropdownMenuItem onClick={() => navigate(`/admin/coupons/${coupon.id}/edit`)}>
                                <Edit className="mr-2 h-4 w-4" />
                                <span>Chỉnh sửa</span>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => handleTestCoupon(coupon.code)}
                              >
                                <span>Thử áp dụng</span>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                className="text-red-600"
                                onClick={() => handleDeleteCoupon(coupon.id)}
                              >
                                <Trash2 className="mr-2 h-4 w-4" />
                                <span>Xóa</span>
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {totalPages > 1 && (
              <div className="flex justify-center mt-4">
                <Pagination
                  currentPage={currentPage}
                  totalPages={totalPages}
                  onPageChange={setCurrentPage}
                />
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default CouponManagement; 