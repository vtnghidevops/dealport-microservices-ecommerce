import React, { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useToast } from "@/hooks/use-toast";
import paymentService from '@/services/user/payment.service'
import { BsCheckCircleFill, BsXCircleFill } from "react-icons/bs";

const PaymentResult: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { toast } = useToast();

  const [loading, setLoading] = useState(true);
  const [paymentSuccess, setPaymentSuccess] = useState(false);
  const [orderId, setOrderId] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    // Get all URL parameters as an object
    const params = Object.fromEntries(searchParams.entries());

    // Determine payment gateway from URL or params
    const paymentGateway = params.paymentGateway ||
      (params.partnerCode ? 'momo' :
        (params.vnp_TxnRef ? 'vnpay' : null));

    if (!paymentGateway) {
      setLoading(false);
      setPaymentSuccess(false);
      setMessage("Không thể xác định cổng thanh toán. Vui lòng liên hệ hỗ trợ.");
      return;
    }

    const verifyPayment = async () => {
      try {
        let result;

        // Verify payment based on the payment gateway
        if (paymentGateway === 'momo') {
          result = await paymentService.verifyMomoPayment(params);
        // } else if (paymentGateway === 'vnpay') {
        //   result = await paymentService.verifyVnpayPayment(params);
        // } 
        }
        else {
          throw new Error("Cổng thanh toán không được hỗ trợ");
        }

        // Set state based on result
        setPaymentSuccess(result.success);
        setOrderId(result.orderId);
        setMessage(result.message);

        // Show toast message
        if (result.success) {
          toast({
            title: "Thanh toán thành công",
            description: "Đơn hàng của bạn đã được thanh toán",
            variant: "success"
          });

          // Redirect to success page after a short delay
          setTimeout(() => {
            navigate(`/user/checkout/success?orderId=${result.orderId}&method=${paymentGateway}`);
          }, 2000);
        } else {
          toast({
            title: "Thanh toán thất bại",
            description: result.message || "Đã xảy ra lỗi trong quá trình thanh toán",
            variant: "destructive"
          });
        }
      } catch (error) {
        setPaymentSuccess(false);
        setMessage(error instanceof Error ? error.message : "Đã xảy ra lỗi không xác định");

        toast({
          title: "Lỗi xác thực thanh toán",
          description: error instanceof Error ? error.message : "Đã xảy ra lỗi không xác định",
          variant: "destructive"
        });
      } finally {
        setLoading(false);
      }
    };

    verifyPayment();
  }, [searchParams, navigate, toast]);

  return (
    <div className="max-w-[1440px] min-w-[1024px] px-[5rem] py-[3rem] mx-auto">
      <div className="max-w-2xl mx-auto bg-white p-8 rounded-lg shadow-sm text-center">
        {loading ? (
          <div className="py-10">
            <div className="animate-spin mx-auto h-12 w-12 border-4 border-blue-500 border-t-transparent rounded-full"></div>
            <p className="mt-4 text-lg">Đang xác thực thanh toán...</p>
          </div>
        ) : paymentSuccess ? (
          <div className="py-8">
            <BsCheckCircleFill className="text-green-500 text-6xl mx-auto mb-4" />
            <h1 className="text-2xl font-bold text-gray-800 mb-2">Thanh toán thành công!</h1>
            <p className="text-gray-600 mb-6">{message || "Đơn hàng của bạn đã được thanh toán thành công."}</p>
            <p className="text-gray-500">Mã đơn hàng: <span className="font-medium">{orderId}</span></p>
            <p className="text-gray-500 mt-6">Bạn sẽ được chuyển hướng đến trang xác nhận đơn hàng...</p>
          </div>
        ) : (
          <div className="py-8">
            <BsXCircleFill className="text-red-500 text-6xl mx-auto mb-4" />
            <h1 className="text-2xl font-bold text-gray-800 mb-2">Thanh toán thất bại</h1>
            <p className="text-gray-600 mb-6">{message || "Đã xảy ra lỗi trong quá trình thanh toán. Vui lòng thử lại."}</p>

            <div className="flex gap-4 justify-center mt-8">
              <button
                onClick={() => navigate("/user/cart")}
                className="px-6 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
              >
                Quay lại giỏ hàng
              </button>
              <button
                onClick={() => navigate("/")}
                className="px-6 py-2 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
              >
                Tiếp tục mua sắm
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default PaymentResult; 