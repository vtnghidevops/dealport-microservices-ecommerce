// components/admin/role/components/PasswordInput.tsx
import React from 'react';
import { FaRegEye, FaRegEyeSlash } from "react-icons/fa";
export interface PasswordInputProps {
  type: string;
  placeholder: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  showPassword: boolean;
  toggleShow: () => void;
  disabled?: boolean;
}

export const PasswordInput: React.FC<PasswordInputProps> = ({
  type,
  placeholder,
  value,
  onChange,
  showPassword,
  toggleShow,
  disabled = false
}) => {
  return (
    <div className="relative">
      <input
        type={showPassword ? "text" : "password"}
        className="text-cyprus w-full border border-gray-300 rounded-md px-3 py-2 pr-10 focus:outline-none focus:border-aqua-spring"
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        disabled={disabled}
      />
      <button
        type="button"
        className="absolute inset-y-0 right-0 pr-3 flex items-center"
        onClick={toggleShow}
        disabled={disabled}
      >
        {showPassword ? (
          <FaRegEye></FaRegEye>
        ) : (
          <FaRegEyeSlash></FaRegEyeSlash>
        )}
      </button>
    </div>
  );
};