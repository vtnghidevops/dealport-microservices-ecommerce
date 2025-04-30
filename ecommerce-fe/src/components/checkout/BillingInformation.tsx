import React, { useState, useEffect } from "react";
import { BillingInfo, ShippingInfo } from '@/types/checkout.model';

interface BillingInformationProps {
  onBillingInfoChange?: (info: Partial<BillingInfo>) => void;
  onShippingInfoChange?: (info: Partial<ShippingInfo>) => void;
}

const BillingInformation: React.FC<BillingInformationProps> = ({
  onBillingInfoChange,
  onShippingInfoChange
}) => {
  const [firstNameFocus, setFirstNameFocus] = useState(false);
  const [lastNameFocus, setLastNameFocus] = useState(false);
  const [companyNameFocus, setCompanyNameFocus] = useState(false);
  const [addressFocus, setAddressFocus] = useState(false);
  const [zipCodeFocus, setZipCodeFocus] = useState(false);
  const [emailFocus, setEmailFocus] = useState(false);
  const [phoneNumberFocus, setPhoneNumberFocus] = useState(false);
  const [differentShippingAddress, setDifferentShippingAddress] = useState(false);

  // Shipping address form states
  const [shippingFirstNameFocus, setShippingFirstNameFocus] = useState(false);
  const [shippingLastNameFocus, setShippingLastNameFocus] = useState(false);
  const [shippingCompanyNameFocus, setShippingCompanyNameFocus] = useState(false);
  const [shippingAddressFocus, setShippingAddressFocus] = useState(false);
  const [shippingZipCodeFocus, setShippingZipCodeFocus] = useState(false);

  const getInputValue = (id: string): string =>
    (document.getElementById(id) as HTMLInputElement)?.value || '';

  // Gather billing info and update parent component
  useEffect(() => {
    if (!onBillingInfoChange) return;

    const updateBillingInfo = () => {
      const billingInfo: Partial<BillingInfo> = {
        firstName: getInputValue("first-name"),
        lastName: getInputValue("last-name"),
        address: getInputValue("address"),
        country: getInputValue("country"),
        region: getInputValue("region-state"),
        city: getInputValue("city"),
        zipCode: getInputValue("zip-code"),
        email: getInputValue("email"),
        phone: getInputValue("phone-number")
      };

      onBillingInfoChange(billingInfo);
    };

    // Add event listeners to form inputs
    const inputs = document.querySelectorAll("input, select");
    inputs.forEach(input => {
      input.addEventListener('change', updateBillingInfo);
      input.addEventListener('blur', updateBillingInfo);
    });

    // Initial update
    updateBillingInfo();

    // Cleanup
    return () => {
      inputs.forEach(input => {
        input.removeEventListener('change', updateBillingInfo);
        input.removeEventListener('blur', updateBillingInfo);
      });
    };
  }, [onBillingInfoChange]);

  // Update shipping info when checkbox changes
  useEffect(() => {
    if (!onShippingInfoChange) return;

    const shippingInfo: Partial<ShippingInfo> = {
      shipToDifferentAddress: differentShippingAddress,
      shippingMethod: "standard" // Default shipping method
    };

    // If different shipping address is selected, add shipping address details
    if (differentShippingAddress) {
      shippingInfo.firstName = getInputValue("shipping-first-name");
      shippingInfo.lastName = getInputValue("shipping-last-name");
      shippingInfo.address = getInputValue("shipping-address");
      shippingInfo.country = getInputValue("shipping-country");
      shippingInfo.region = getInputValue("shipping-region-state");
      shippingInfo.city = getInputValue("shipping-city");
      shippingInfo.zipCode = getInputValue("shipping-zip-code");
    }

    onShippingInfoChange(shippingInfo);
  }, [differentShippingAddress, onShippingInfoChange]);

  return (
    <div className="p-6 bg-white">
      <h2 className="text-[18px] font-semibold mb-5 text-gray-800">Billing Information</h2>
      <form className="grid gap-5">
        {/* User Name Section */}
        <div className="flex items-center gap-8">
          <div className="relative w-[50%]">
            <label
              className={`left-2 absolute text-sm transition-all duration-200 ${firstNameFocus || getInputValue("first-name")
                ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                }`}
              htmlFor="first-name"
            >
              First name
            </label>
            <input
              type="text"
              id="first-name"
              className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                       focus:border-blue-500 hover:border-gray-400"
              onFocus={() => setFirstNameFocus(true)}
              onBlur={(e) => setFirstNameFocus(!!e.target.value)}
            />
          </div>

          <div className="relative w-[50%]">
            <label
              className={`absolute text-sm transition-all duration-200 ${lastNameFocus || getInputValue("last-name")
                ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                }`}
              htmlFor="last-name"
            >
              Last name
            </label>
            <input
              type="text"
              id="last-name"
              className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                       focus:border-blue-500 hover:border-gray-400"
              onFocus={() => setLastNameFocus(true)}
              onBlur={(e) => setLastNameFocus(!!e.target.value)}
            />
          </div>
        </div>

        {/* Address Section */}
        <div className="relative">
          <label
            className={`absolute text-sm transition-all duration-200 ${addressFocus || getInputValue("address")
              ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
              : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
              }`}
            htmlFor="address"
          >
            Address
          </label>
          <input
            type="text"
            id="address"
            className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                     focus:border-blue-500 hover:border-gray-400"
            onFocus={() => setAddressFocus(true)}
            onBlur={(e) => setAddressFocus(!!e.target.value)}
          />
        </div>

        {/* Location Details */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="relative">
            <select
              id="country"
              className="w-full p-3 border-2 rounded-lg outline-none appearance-none
                       bg-white transition-colors focus:border-blue-500 hover:border-gray-400"
            >
              <option value="">Select Country</option>
              <option value="VN">Vietnam</option>
              <option value="US">United States</option>
            </select>
            <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
              <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>

          <div className="relative">
            <select
              id="region-state"
              className="w-full p-3 border-2 rounded-lg outline-none appearance-none
                       bg-white transition-colors focus:border-blue-500 hover:border-gray-400"
            >
              <option value="">Select State</option>
              <option value="HCM">Ho Chi Minh</option>
              <option value="HN">Ha Noi</option>
              <option value="DN">Da Nang</option>
            </select>
            <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
              <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>

          <div className="relative">
            <select
              id="city"
              className="w-full p-3 border-2 rounded-lg outline-none appearance-none
                       bg-white transition-colors focus:border-blue-500 hover:border-gray-400"
            >
              <option value="">Select City</option>
              <option value="D1">District 1</option>
              <option value="D2">District 2</option>
              <option value="D3">District 3</option>
            </select>
            <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
              <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>

          <div className="relative">
            <label
              className={`absolute text-sm transition-all duration-200 ${zipCodeFocus || getInputValue("zip-code")
                ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                }`}
              htmlFor="zip-code"
            >
              Zip Code
            </label>
            <input
              type="text"
              id="zip-code"
              className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                       focus:border-blue-500 hover:border-gray-400"
              onFocus={() => setZipCodeFocus(true)}
              onBlur={(e) => setZipCodeFocus(!!e.target.value)}
            />
          </div>
        </div>

        {/* Contact Information */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="relative">
            <label
              className={`absolute text-sm transition-all duration-200 ${emailFocus || getInputValue("email")
                ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                }`}
              htmlFor="email"
            >
              Email
            </label>
            <input
              type="email"
              id="email"
              className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                       focus:border-blue-500 hover:border-gray-400"
              onFocus={() => setEmailFocus(true)}
              onBlur={(e) => setEmailFocus(!!e.target.value)}
            />
          </div>

          <div className="relative">
            <label
              className={`absolute text-sm transition-all duration-200 ${phoneNumberFocus || getInputValue("phone-number")
                ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                }`}
              htmlFor="phone-number"
            >
              Phone Number
            </label>
            <input
              type="tel"
              id="phone-number"
              className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                       focus:border-blue-500 hover:border-gray-400"
              onFocus={() => setPhoneNumberFocus(true)}
              onBlur={(e) => setPhoneNumberFocus(!!e.target.value)}
            />
          </div>
        </div>

        {/* Ship to different address checkbox */}
        <div className="mt-4">
          <div className="flex items-center">
            <input
              type="checkbox"
              id="ship-to-different-address"
              className="w-5 h-5 text-blue-600 rounded focus:ring-blue-500 border-gray-300"
              checked={differentShippingAddress}
              onChange={(e) => setDifferentShippingAddress(e.target.checked)}
            />
            <label htmlFor="ship-to-different-address" className="ml-2 text-[16px] font-medium text-gray-700">
              Ship to a different address?
            </label>
          </div>
        </div>

        {/* Shipping Address Section (shown conditionally) */}
        {differentShippingAddress && (
          <div className="mt-4 p-5 border border-gray-200 rounded-lg">
            <h3 className="text-base font-semibold mb-4">Shipping Address</h3>

            {/* Shipping Name Fields */}
            <div className="flex items-center gap-8 mb-5">
              <div className="relative w-[50%]">
                <label
                  className={`absolute text-sm transition-all duration-200 ${shippingFirstNameFocus || getInputValue("shipping-first-name")
                    ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                    : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                    }`}
                  htmlFor="shipping-first-name"
                >
                  First name
                </label>
                <input
                  type="text"
                  id="shipping-first-name"
                  className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                           focus:border-blue-500 hover:border-gray-400"
                  onFocus={() => setShippingFirstNameFocus(true)}
                  onBlur={(e) => setShippingFirstNameFocus(!!e.target.value)}
                />
              </div>

              <div className="relative w-[50%]">
                <label
                  className={`absolute text-sm transition-all duration-200 ${shippingLastNameFocus || getInputValue("shipping-last-name")
                    ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                    : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                    }`}
                  htmlFor="shipping-last-name"
                >
                  Last name
                </label>
                <input
                  type="text"
                  id="shipping-last-name"
                  className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                           focus:border-blue-500 hover:border-gray-400"
                  onFocus={() => setShippingLastNameFocus(true)}
                  onBlur={(e) => setShippingLastNameFocus(!!e.target.value)}
                />
              </div>
            </div>

            {/* Company Name Field */}
            <div className="relative mb-5">
              <label
                className={`absolute text-sm transition-all duration-200 ${shippingCompanyNameFocus || getInputValue("shipping-company-name")
                  ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                  : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                  }`}
                htmlFor="shipping-company-name"
              >
                Company Name (Optional)
              </label>
              <input
                type="text"
                id="shipping-company-name"
                className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                         focus:border-blue-500 hover:border-gray-400"
                onFocus={() => setShippingCompanyNameFocus(true)}
                onBlur={(e) => setShippingCompanyNameFocus(!!e.target.value)}
              />
            </div>

            {/* Shipping Address */}
            <div className="relative mb-5">
              <label
                className={`absolute text-sm transition-all duration-200 ${shippingAddressFocus || getInputValue("shipping-address")
                  ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                  : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                  }`}
                htmlFor="shipping-address"
              >
                Address
              </label>
              <input
                type="text"
                id="shipping-address"
                className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                         focus:border-blue-500 hover:border-gray-400"
                onFocus={() => setShippingAddressFocus(true)}
                onBlur={(e) => setShippingAddressFocus(!!e.target.value)}
              />
            </div>

            {/* Shipping Location Details */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              <div className="relative">
                <select
                  id="shipping-country"
                  className="w-full p-3 border-2 rounded-lg outline-none appearance-none
                           bg-white transition-colors focus:border-blue-500 hover:border-gray-400"
                >
                  <option value="">Select Country</option>
                  <option value="VN">Vietnam</option>
                  <option value="US">United States</option>
                </select>
                <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </div>

              <div className="relative">
                <select
                  id="shipping-region-state"
                  className="w-full p-3 border-2 rounded-lg outline-none appearance-none
                           bg-white transition-colors focus:border-blue-500 hover:border-gray-400"
                >
                  <option value="">Select State</option>
                  <option value="HCM">Ho Chi Minh</option>
                  <option value="HN">Ha Noi</option>
                </select>
                <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </div>

              <div className="relative">
                <select
                  id="shipping-city"
                  className="w-full p-3 border-2 rounded-lg outline-none appearance-none
                           bg-white transition-colors focus:border-blue-500 hover:border-gray-400"
                >
                  <option value="">Select City</option>
                  <option value="D1">District 1</option>
                  <option value="D2">District 2</option>
                </select>
                <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </div>

              <div className="relative">
                <label
                  className={`absolute text-sm transition-all duration-200 ${shippingZipCodeFocus || getInputValue("shipping-zip-code")
                    ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                    : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
                    }`}
                  htmlFor="shipping-zip-code"
                >
                  Zip Code
                </label>
                <input
                  type="text"
                  id="shipping-zip-code"
                  className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                           focus:border-blue-500 hover:border-gray-400"
                  onFocus={() => setShippingZipCodeFocus(true)}
                  onBlur={(e) => setShippingZipCodeFocus(!!e.target.value)}
                />
              </div>
            </div>
          </div>
        )}

        {/* Order Notes */}
        <div className="mt-4">
          <label htmlFor="order-notes" className="block text-sm font-medium text-gray-700 mb-2">
            Order Notes (Optional)
          </label>
          <textarea
            id="order-notes"
            rows={3}
            placeholder="Notes about your order, e.g. special notes for delivery."
            className="w-full p-3 border-2 rounded-lg outline-none resize-none focus:border-blue-500"
          ></textarea>
        </div>
      </form>
    </div>
  );
};

export default BillingInformation;  