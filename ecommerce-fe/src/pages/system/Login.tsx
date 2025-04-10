import React, { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FcGoogle } from "react-icons/fc";
import { FaApple } from "react-icons/fa";
import { Link } from "react-router-dom";

interface LoginFormData {
  email: string;
  password: string;
}

const Login: React.FC = () => {
  const [formData, setFormData] = useState<LoginFormData>({
    email: "",
    password: ""
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleLogin = (): void => {
    console.log("Login data", formData);
    // TODO: Gọi API backend để đăng nhập
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <div className="flex justify-between items-center pb-2 mb-4">
        <button className="w-1/2 text-[18px] pb-3 border-b-4 border-orange-500">
          Sign In
        </button>
        <Link 
          to="/register" 
          className="flex text-[18px] pb-3 items-center justify-center text-neutral-500 w-1/2 hover:text-neutral-800"
        >
          Sign Up
        </Link>
      </div>

      <div className="space-y-4 mt-5">
        <Input
          type="email"
          name="email"
          placeholder="Email Address"
          value={formData.email}
          onChange={handleChange}
          className="mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        
        <div className="relative mb-5">
          <Input
            type="password"
            name="password"
            placeholder="Password"
            value={formData.password}
            onChange={handleChange}
            className="mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          />
          <Link
            to="/forgot-password"
            className="absolute -bottom-[35px] right-2 text-sm text-blue-500 hover:text-blue-700"
          >
            Forgot Password
          </Link>
        </div>

        <Button
          onClick={handleLogin}
          className="font-sans !mt-[2.5rem] w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px]"
        >
          SIGN IN
        </Button>

        <div className="text-center text-gray-500">or</div>

        <Button
          variant="outline"
          className="!mb-5 h-[44px] w-full flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
        >
          <FcGoogle className="!h-[20px] !w-[20px]" /> Login with Google
        </Button>

        <Button
          variant="outline"
          className="w-full h-[44px] flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
        >
          <FaApple className="!h-[20px] !w-[20px]" /> Login with Apple
        </Button>
      </div>
    </div>
  );
};

export default Login;