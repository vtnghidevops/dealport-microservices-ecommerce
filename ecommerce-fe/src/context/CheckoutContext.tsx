// src/context/CheckoutContext.tsx - phần cần sửa
import { createContext, useContext, useState, ReactNode } from 'react';

// Cập nhật các phương thức thanh toán được chấp nhận
type PaymentMethod = 'momo' | 'banking';

interface BillingInfo {
  firstName: string;
  lastName: string;
  companyName?: string;
  address: string;
  country: string;
  region: string;
  city: string;
  zipCode: string;
  email: string;
  phone: string;
  shipToDifferentAddress: boolean;
}

interface CheckoutContextType {
  billingInfo: BillingInfo;
  updateBillingInfo: (info: Partial<BillingInfo>) => void;
  paymentMethod: PaymentMethod;
  setPaymentMethod: (method: PaymentMethod) => void;
  // Bỏ cardInfo và updateCardInfo
  orderNotes: string;
  setOrderNotes: (notes: string) => void;
  validateCheckout: () => boolean;
}

const initialBillingInfo: BillingInfo = {
  firstName: '',
  lastName: '',
  companyName: '',
  address: '',
  country: '',
  region: '',
  city: '',
  zipCode: '',
  email: '',
  phone: '',
  shipToDifferentAddress: false
};

const CheckoutContext = createContext<CheckoutContextType | undefined>(undefined);

export const CheckoutProvider = ({ children }: { children: ReactNode }) => {
  const [billingInfo, setBillingInfo] = useState<BillingInfo>(initialBillingInfo);
  // Mặc định chọn Momo
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('momo');
  const [orderNotes, setOrderNotes] = useState('');

  const updateBillingInfo = (info: Partial<BillingInfo>) => {
    setBillingInfo({ ...billingInfo, ...info });
  };

  const validateCheckout = (): boolean => {
    // Basic validation
    if (!billingInfo.firstName || !billingInfo.lastName || !billingInfo.address || 
        !billingInfo.country || !billingInfo.city || !billingInfo.zipCode || 
        !billingInfo.email || !billingInfo.phone) {
      return false;
    }

    return true;
  };

  return (
    <CheckoutContext.Provider value={{
      billingInfo,
      updateBillingInfo,
      paymentMethod,
      setPaymentMethod,
      orderNotes,
      setOrderNotes,
      validateCheckout
    }}>
      {children}
    </CheckoutContext.Provider>
  );
};

export const useCheckout = () => {
  const context = useContext(CheckoutContext);
  if (context === undefined) {
    throw new Error('useCheckout must be used within a CheckoutProvider');
  }
  return context;
};