import React from "react";

interface Props {
  paymentMethod: string;
  setPaymentMethod: (method: string) => void;
}

const PaymentOption: React.FC<Props> = ({ paymentMethod, setPaymentMethod }) => {
  const paymentMethods = [
    {
      id: "cash",
      label: "Cash on Delivery",
      imageSrc: "/images/payments/cash-money-payment-svgrepo-com.svg",
      alt: "Cash on Delivery Icon"
    },
    {
      id: "momo",
      label: "Momo",
      imageSrc: "/images/payments/momo.svg",
      alt: "Momo Payment Icon"

    },
    {
      id: "banking",
      label: "Banking",
      imageSrc: "/images/payments/banking-bank-svgrepo-com.svg",
      alt: "Banking Icon"
    },
  ];

  return (
    <div className="p-6 bg-white rounded-lg shadow-sm mt-5">
      <h2 className="text-xl font-semibold mb-6">Payment Option</h2>
      <div className="flex gap-8">
        {paymentMethods.map((method) => (
          <div
            key={method.id}
            className={`w-[200px] h-[100px] flex flex-col items-center justify-center p-4 border-2 rounded-lg cursor-pointer transition-all duration-200
              ${paymentMethod === method.id 
                ? 'border-blue-500 bg-blue-50 text-blue-600' 
                : 'border-gray-200 hover:border-blue-400 hover:shadow-md'
              }`}
            onClick={() => setPaymentMethod(method.id)}
          >
            <div className="mb-3 h-[40px] w-[40px]">
            <img
                src={method.imageSrc}
                alt={method.alt}
                className="object-contain transition-transform duration-200 group-hover:scale-105"
              />
            </div>
            <div className="text-sm font-medium">{method.label}</div>
            <input
              type="radio"
              name="payment"
              value={method.id}
              checked={paymentMethod === method.id}
              className="hidden"
              onChange={() => {}}
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default PaymentOption;