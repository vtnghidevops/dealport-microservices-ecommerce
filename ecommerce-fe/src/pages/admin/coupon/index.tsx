import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { PlusCircle, Search, Edit, Trash2, Eye, Tag } from 'lucide-react';
import { Pagination } from '@/components/ui/pagination';
import couponService, { Coupon } from '@/services/admin/coupon.service';
import { useToast } from '@/hooks/use-toast';
import AdminHeader from '@/components/admin/layout/AdminHeader';

const CouponManagement: React.FC = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [coupons, setCoupons] = useState<Coupon[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [apiError, setApiError] = useState<string | null>(null);

  // Fetch coupons from API
  useEffect(() => {
    const fetchCoupons = async () => {
      setIsLoading(true);
      setApiError(null);

      try {
        const response = await couponService.getCoupons(currentPage, 10);
        //console.log('Coupon data from service:', response.coupons);
        setCoupons(response.coupons);
        setTotalPages(Math.ceil(response.total / 10));
      } catch (error) {
        //console.error('Error fetching coupons:', error);
        setApiError('Unable to load coupons. Please try again later.');
        toast({
          title: 'Error',
          description: 'Unable to load coupons',
          variant: 'destructive',
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchCoupons();
  }, [currentPage, toast]);

  // Filter coupons by search term
  const filteredCoupons = coupons.filter(coupon =>
    coupon.code.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Handle coupon deletion
  const handleDeleteCoupon = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this coupon?')) {
      try {
        await couponService.deleteCoupon(id);
        setCoupons(prevCoupons => prevCoupons.filter(coupon => coupon.id !== id));

        toast({
          variant: "success",
          title: 'Success',
          description: 'Coupon deleted successfully',
        });
      } catch (error) {
        console.error('Error deleting coupon:', error);
        toast({
          title: 'Error',
          description: 'Unable to delete coupon. Please try again later.',
          variant: 'destructive'
        });
      }
    }
  };

  // Handle applying a coupon to test it
  const handleTestCoupon = async (code: string) => {
    try {
      await couponService.applyCoupon(code);
      toast({
        variant: "success",
        title: 'Success',
        description: `Coupon ${code} applied to cart`,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: error instanceof Error ? error.message : 'Unable to apply coupon',
        variant: 'destructive'
      });
    }
  };

  // Format date
  const formatDate = (dateString: string) => {
    if (!dateString) return 'N/A';

    try {
      // Try to parse the date, handling both ISO format and RFC3339 format
      let date;
      if (dateString.includes('T')) {
        // RFC3339 or ISO format: "2025-05-01T00:00:00Z"
        date = new Date(dateString);
      } else {
        // Simple date format: "2025-05-01"
        const [year, month, day] = dateString.split('-').map(Number);
        date = new Date(year, month - 1, day);  // month is 0-based in JS Date
      }

      // Check if date is valid
      if (isNaN(date.getTime())) {
        console.error('Invalid date:', dateString);
        return 'Invalid Date';
      }

      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch (error) {
      console.error('Error formatting date:', dateString, error);
      return 'Invalid Date';
    }
  };

  // Check if a coupon is expired
  const isExpired = (validTo: string) => {
    if (!validTo) return false;
    try {
      const expiryDate = new Date(validTo);
      return !isNaN(expiryDate.getTime()) && expiryDate < new Date();
    } catch (error) {
      console.error('Error checking expiry:', validTo, error);
      return false;
    }
  };

  // Get coupon status
  const getCouponStatus = (coupon: Coupon) => {
    // Debugging
    console.log('Checking status for coupon:', coupon.code, {
      isActive: coupon.isActive,
      validTo: coupon.validTo,
      validFrom: coupon.validFrom,
      usageCount: coupon.usageCount,
      maxUsage: coupon.maxUsage
    });

    // Check for undefined values and use defaults if needed
    const isActive = coupon.isActive !== undefined ? coupon.isActive : true;
    const maxUsage = coupon.maxUsage || 0;
    const usageCount = coupon.usageCount || 0;

    if (isActive === false) {
      console.log(`Coupon ${coupon.code} is inactive`);
      return 'inactive';
    }

    if (isExpired(coupon.validTo)) {
      console.log(`Coupon ${coupon.code} is expired`);
      return 'expired';
    }

    if (usageCount >= maxUsage && maxUsage > 0) {
      console.log(`Coupon ${coupon.code} is fully used`);
      return 'used';
    }

    console.log(`Coupon ${coupon.code} is active`);
    return 'active';
  };

  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Coupon Management" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            {/* Button Row */}
            <div className="flex justify-end mb-5">
              <Button
                onClick={() => navigate('/admin/coupons/create')}
                className="text-white rounded-lg h-[48px] w-[160px] py-6 px-8 bg-ocean-green hover:bg-green-700 flex items-center justify-center"
              >
                <PlusCircle size={16} className="mr-2" />
                Add New Coupon
              </Button>
            </div>

            {/* Coupon Table */}
            <div className="w-[1116px] bg-white rounded-lg shadow p-[1rem] mb-6 mt-[1rem] drop-shadow filter">
              <div className="mb-5">
                <h2 className="text-xl font-semibold mb-3">Coupons</h2>

                <div className="relative w-[320px] h-[48px] flex items-center mb-5">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-muted-foreground" />
                  <Input
                    placeholder="Search by coupon code"
                    className="pl-10 h-full w-full border border-gray-300 rounded-md"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                  />
                </div>
              </div>

              {isLoading ? (
                <div className="flex justify-center items-center py-8">
                  <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
                </div>
              ) : apiError ? (
                <div className="text-center py-8 text-red-500">
                  {apiError}
                </div>
              ) : (
                <>
                  <div className="overflow-x-auto bg-white rounded-lg shadow">
                    <table className="min-w-full">
                      <thead>
                        <tr className="bg-aqua-spring h-[56px] border-b border-gray-200">
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Coupon Code
                          </th>
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Discount
                          </th>
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Min Order
                          </th>
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Usage
                          </th>
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Validity Period
                          </th>
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Status
                          </th>
                          <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus tracking-wider">
                            Actions
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {filteredCoupons.length === 0 ? (
                          <tr>
                            <td colSpan={7} className="text-center py-4">
                              No coupons found
                            </td>
                          </tr>
                        ) : (
                          filteredCoupons.map((coupon) => (
                            <tr
                              key={coupon.id}
                              className="hover:bg-gray-50 cursor-pointer h-[64px]"
                            >
                              <td className="text-center px-6 py-4 whitespace-nowrap text-[15px] text-black">
                                {coupon.code}
                              </td>
                              <td className="text-center px-6 py-4 whitespace-nowrap text-[15px] text-black">
                                {coupon.discountType === 'percentage'
                                  ? `${coupon.discount}%`
                                  : `$${coupon.discount?.toFixed(2) || '0.00'}`
                                }
                              </td>
                              <td className="text-center px-6 py-4 whitespace-nowrap text-[15px] text-black">
                                ${coupon.minOrderAmount !== undefined ? coupon.minOrderAmount.toFixed(2) : '0.00'}
                              </td>
                              <td className="text-center px-6 py-4 whitespace-nowrap text-[15px] text-black">
                                {coupon.usageCount !== undefined ? coupon.usageCount : 0}/{coupon.maxUsage || 0}
                              </td>
                              <td className="text-center px-6 py-4 whitespace-nowrap text-[15px] text-black">
                                {formatDate(coupon.validFrom)} - {formatDate(coupon.validTo)}
                              </td>
                              <td className="text-center px-6 py-4 whitespace-nowrap">
                                <span className={`!text-[15px] px-2 inline-flex text-xs leading-5 rounded-full
                                  ${getCouponStatus(coupon) === 'active' ? 'text-success' : ''} 
                                  ${getCouponStatus(coupon) === 'inactive' ? 'text-error' : ''}
                                  ${getCouponStatus(coupon) === 'expired' ? 'text-error' : ''}
                                  ${getCouponStatus(coupon) === 'used' ? 'text-warning' : ''}`}>
                                  • {getCouponStatus(coupon) === 'active' && 'Active'}
                                  {getCouponStatus(coupon) === 'inactive' && 'Inactive'}
                                  {getCouponStatus(coupon) === 'expired' && 'Expired'}
                                  {getCouponStatus(coupon) === 'used' && 'Fully Used'}
                                </span>
                              </td>
                              <td className="text-center px-6 py-4 whitespace-nowrap">
                                <button
                                  className="text-gray-500 hover:text-gray-700 mr-3"
                                  onClick={() => navigate(`/admin/coupons/${coupon.id}`)}
                                >
                                  <Eye className="h-5 w-5" />
                                </button>
                                <button
                                  className="text-gray-500 hover:text-gray-700 mr-3"
                                  onClick={() => navigate(`/admin/coupons/${coupon.id}/edit`)}
                                >
                                  <Edit className="h-5 w-5" />
                                </button>
                                <button
                                  className="text-gray-500 hover:text-gray-700 mr-3"
                                  onClick={() => handleTestCoupon(coupon.code)}
                                >
                                  <Tag className="h-5 w-5" />
                                </button>
                                <button
                                  className="text-gray-500 hover:text-red-700"
                                  onClick={() => handleDeleteCoupon(coupon.id)}
                                >
                                  <Trash2 className="h-5 w-5" />
                                </button>
                              </td>
                            </tr>
                          ))
                        )}
                      </tbody>
                    </table>
                  </div>

                  {totalPages > 1 && (
                    <div className="flex justify-center mt-[3rem]">
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
        </main>
      </div>
    </div>
  );
};

export default CouponManagement; 