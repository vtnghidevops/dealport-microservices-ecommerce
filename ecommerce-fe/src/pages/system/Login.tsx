import React, { useState, useCallback } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FcGoogle } from "react-icons/fc";
import { FaApple } from "react-icons/fa";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useCart } from "@/hooks/useCart";
import { useToast } from "@/hooks/use-toast";
import { FiEye, FiEyeOff } from "react-icons/fi";
import { extractErrorMessage } from "@/utils/error-handler";

interface LoginFormData {
  email: string;
  password: string;
}

const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const { refreshCartTTL } = useCart();
  const [formData, setFormData] = useState<LoginFormData>({
    email: "",
    password: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();
  const [showPassword, setShowPassword] = useState(false);

  // Memoize input change handler
  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>): void => {
      const { name, value } = e.target;
      setFormData((prev) => ({
        ...prev,
        [name]: value,
      }));
    },
    []
  );

  // Memoize toggle password visibility handler
  const togglePasswordVisibility = useCallback(() => {
    setShowPassword((prev) => !prev);
  }, []);

  // Memoize social login handler
  const handleSocialLogin = useCallback(
    (provider: string) => {
      toast({
        title: "Coming Soon",
        description: `${
          provider.charAt(0).toUpperCase() + provider.slice(1)
        } login will be available soon!`,
        variant: "default",
      });
    },
    [toast]
  );

  const handleLogin = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();

    // Validation
    if (!formData.email || !formData.password) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please fill in all required fields",
      });
      return;
    }

    setIsLoading(true);

    const result = await login(formData.email, formData.password);

    if (result.success) {
      // Sau khi đăng nhập thành công, refresh cart TTL
      try {
        await refreshCartTTL();
      } catch (err) {
        console.error("Failed to refresh cart TTL after login:", err);
      }

      toast({
        variant: "success",
        title: "Success",
        description: "Login successful!",
      });
      navigate("/");
    } else {
      const friendlyMsg = extractErrorMessage(result.error);
      toast({
        variant: "destructive",
        title: "Login Failed",
        description: friendlyMsg,
      });
    }
    setIsLoading(false);
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

      <form onSubmit={handleLogin} className="space-y-4 mt-5">
        <Input
          type="email"
          name="email"
          placeholder="Email Address"
          value={formData.email}
          onChange={handleChange}
          className="mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isLoading}
          required
        />

        <div className="relative mb-5">
          <Input
            type={showPassword ? "text" : "password"}
            name="password"
            placeholder="Password"
            value={formData.password}
            onChange={handleChange}
            className="mb-3 h-[44px] focus:border-2 focus:border-blue-400 pr-10"
            disabled={isLoading}
            required
          />
          <button
            type="button"
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-700"
            tabIndex={-1}
            onClick={togglePasswordVisibility}
          >
            {showPassword ? <FiEyeOff /> : <FiEye />}
          </button>
          <Link
            to="/forgot-password"
            className="absolute -bottom-[35px] right-2 text-sm text-blue-500 hover:text-blue-700"
          >
            Forgot Password
          </Link>
        </div>

        <Button
          type="submit"
          className="font-sans !mt-[2.5rem] w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px]"
          disabled={isLoading}
        >
          {isLoading ? "SIGNING IN..." : "SIGN IN"}
        </Button>

        <div className="text-center text-gray-500">or</div>

        <Button
          type="button"
          variant="outline"
          className="!mb-5 h-[44px] w-full flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
          disabled={isLoading}
          onClick={() => handleSocialLogin("google")}
        >
          <FcGoogle className="!h-[20px] !w-[20px]" /> Login with Google
        </Button>

        <Button
          type="button"
          variant="outline"
          className="w-full h-[44px] flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
          disabled={isLoading}
          onClick={() => handleSocialLogin("apple")}
        >
          <FaApple className="!h-[20px] !w-[20px]" /> Login with Apple
        </Button>
      </form>
    </div>
  );
};

export default Login;
