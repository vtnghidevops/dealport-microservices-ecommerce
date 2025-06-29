import React from "react";

interface Props {
  paymentMethod: string;
  setPaymentMethod: (method: string) => void;
}

const PaymentOption: React.FC<Props> = ({ paymentMethod, setPaymentMethod }) => {
  const paymentMethods = [
    {
      id: "cod",
      label: "Cash on Delivery",
      imageSrc: "/images/payments/cash-money-payment-svgrepo-com.svg",
      alt: "Cash on Delivery Icon",
      disabled: false
    },
    {
      id: "momo",
      label: "Momo",
      imageSrc: "/images/payments/momo.svg",
      alt: "Momo Icon",
      disabled: false
    },
    {
      id: "vnpay",
      label: "VNPAY",
      imageSrc: "/images/payments/vnpay.svg",
      alt: "VNPAY Icon",
      disabled: true
    }
  ];

  return (
    <div className="p-6 bg-white rounded-lg shadow-sm mt-5">
      <h2 className="text-xl font-semibold mb-6">Payment Option</h2>
      <div className="flex flex-wrap gap-5">
        {paymentMethods.map((method) => (
          <div
            key={method.id}
            className={`w-[200px] h-[100px] flex flex-col items-center justify-center p-4 border-2 rounded-lg transition-all duration-200 relative
              ${method.disabled
                ? 'border-gray-200 bg-gray-50 cursor-not-allowed opacity-70'
                : paymentMethod === method.id
                  ? 'border-blue-500 bg-blue-50 text-blue-600 cursor-pointer'
                  : 'border-gray-200 hover:border-blue-400 hover:shadow-md cursor-pointer'
              }`}
            onClick={() => !method.disabled && setPaymentMethod(method.id)}
          >
            <div className="mb-3 h-[40px] w-[80px]">
              <img
                src={method.imageSrc}
                alt={method.alt}
                className={`w-full h-full object-contain transition-transform duration-200 ${!method.disabled && 'group-hover:scale-105'}`}
              />
            </div>
            <div className="text-sm font-medium">{method.label}</div>
            {method.disabled && (
              <div className="absolute top-1 right-2 bg-orange-500 text-white text-xs px-2 py-0.5 rounded-full">
                Coming Soon
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default PaymentOption;