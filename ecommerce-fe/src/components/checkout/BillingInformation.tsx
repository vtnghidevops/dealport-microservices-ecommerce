import React, { useState } from "react";

const BillingInformation: React.FC = () => {
  const [firstNameFocus, setFirstNameFocus] = useState(false);
  const [lastNameFocus, setLastNameFocus] = useState(false);
  const [companyNameFocus, setCompanyNameFocus] = useState(false);
  const [addressFocus, setAddressFocus] = useState(false);
  const [zipCodeFocus, setZipCodeFocus] = useState(false);
  const [emailFocus, setEmailFocus] = useState(false);
  const [phoneNumberFocus, setPhoneNumberFocus] = useState(false);

  const getInputValue = (id: string): string => 
    (document.getElementById(id) as HTMLInputElement)?.value || '';

  return (
    <div className="p-6 bg-white">
      <h2 className="text-[18px] font-semibold mb-5 text-gray-800">Billing Information</h2>
      <form className="grid gap-5">
        {/* User Name Section */}
        <div className="flex items-center gap-8">
          <div className="relative w-[30%]">
            <label
              className={`left-2 absolute text-sm transition-all duration-200 ${
                firstNameFocus || getInputValue("first-name")
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

          <div className="relative w-[30%]">
            <label
              className={`absolute text-sm transition-all duration-200 ${
                lastNameFocus || getInputValue("last-name")
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

          <div className="relative w-[40%]">
            <label
              className={`absolute text-sm transition-all duration-200 ${
                companyNameFocus || getInputValue("company-name")
                  ? "-top-2.5 left-2 bg-white px-2 text-blue-600"
                  : "left-5 top-1/2 -translate-y-1/2 text-gray-500"
              }`}
              htmlFor="company-name"
            >
              Company Name (Optional)
            </label>
            <input
              type="text"
              id="company-name"
              className="w-full p-3 border-2 rounded-lg outline-none transition-colors
                       focus:border-blue-500 hover:border-gray-400"
              onFocus={() => setCompanyNameFocus(true)}
              onBlur={(e) => setCompanyNameFocus(!!e.target.value)}
            />
          </div>
        </div>

        {/* Address Section */}
        <div className="relative">
          <label
            className={`absolute text-sm transition-all duration-200 ${
              addressFocus || getInputValue("address")
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
              {/* Add state options */}
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
              {/* Add city options */}
            </select>
            <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
              <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>

          <div className="relative">
            <label
              className={`absolute text-sm transition-all duration-200 ${
                zipCodeFocus || getInputValue("zip-code")
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
              className={`absolute text-sm transition-all duration-200 ${
                emailFocus || getInputValue("email")
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
              className={`absolute text-sm transition-all duration-200 ${
                phoneNumberFocus || getInputValue("phone-number")
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
      </form>
    </div>
  );
};

export default BillingInformation;  