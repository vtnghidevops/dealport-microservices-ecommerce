import { useState, FormEvent } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FcGoogle } from "react-icons/fc";
import { FaApple } from "react-icons/fa";
import { Link, useNavigate } from "react-router-dom";
import { UserRegistrationData } from "@/types/user.model";
import { useAuth } from "@/hooks/useAuth";
import { useToast } from "@/hooks/use-toast";
import { FiEye, FiEyeOff } from "react-icons/fi";

const Register: React.FC = () => {
  const navigate = useNavigate();
  const { register, authState } = useAuth();
  const { toast } = useToast();
  const [formData, setFormData] = useState<UserRegistrationData>({
    firstName: '',
    lastName: '',
    username: '',
    email: '',
    phone: '',
    password: '',
    confirmPassword: '',
    acceptTerms: false,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleRegister = async (e: FormEvent): Promise<void> => {
    e.preventDefault();

    // Validation
    if (!formData.acceptTerms) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please agree to terms and conditions"
      });
      return;
    }
    if (formData.password !== formData.confirmPassword) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Passwords do not match"
      });
      return;
    }
    if (formData.password.length < 8) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Password must be at least 8 characters"
      });
      return;
    }
    if (!formData.email || !formData.firstName || !formData.lastName || !formData.username) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please fill in all required fields"
      });
      return;
    }

    setIsSubmitting(true);

    try {
      // register bây giờ trả về boolean - true nếu thành công
      const success = await register(formData);

      if (success) {
        toast({
          variant: "success",
          title: "Registration Successful!",
          description: "Your account has been created. Please log in to continue."
        });

        // Chuyển hướng đến trang đăng nhập thay vì trang chủ
        navigate("/login");
      } else if (authState.error) {
        toast({
          variant: "destructive",
          title: "Registration Failed",
          description: authState.error || "Please try again."
        });
      }
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Registration Failed",
        description: error.message || "Please try again."
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const validate = () => {
    const errors: Record<string, string> = {};
    if (!formData.firstName) {
      errors.firstName = 'First name is required';
    }
    if (!formData.lastName) {
      errors.lastName = 'Last name is required';
    }
    if (!formData.username) {
      errors.username = 'Username is required';
    }
    if (!formData.email) {
      errors.email = 'Email is required';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      errors.email = 'Email is invalid';
    }
    setErrors(errors);
    return Object.keys(errors).length === 0;
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <div className="flex justify-between items-center pb-2 mb-4">
        <Link to="/login" className="flex text-[18px] pb-3 items-center justify-center text-neutral-500 w-1/2 hover:text-neutral-800">
          Sign In
        </Link>
        <button className="w-1/2 text-[18px] pb-3 border-b-4 border-orange-500">
          Sign Up
        </button>
      </div>

      <form onSubmit={handleRegister} className="space-y-4 mt-5">
        <Input
          name="firstName"
          placeholder="First Name"
          value={formData.firstName}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isSubmitting}
          required
        />
        <Input
          name="lastName"
          placeholder="Last Name"
          value={formData.lastName}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isSubmitting}
          required
        />
        <Input
          name="username"
          placeholder="Username"
          value={formData.username}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isSubmitting}
          required
        />
        <Input
          name="email"
          type="email"
          placeholder="Email Address"
          value={formData.email}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isSubmitting}
          required
        />
        <Input
          name="phone"
          placeholder="Phone Number (Optional)"
          value={formData.phone || ""}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isSubmitting}
        />
        <div className="relative mb-3">
          <Input
            name="password"
            type={showPassword ? "text" : "password"}
            placeholder="8+ characters"
            value={formData.password}
            onChange={handleInputChange}
            className="h-[44px] focus:border-2 focus:border-blue-400 pr-10"
            disabled={isSubmitting}
            required
            minLength={8}
          />
          <button
            type="button"
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-700"
            tabIndex={-1}
            onClick={() => setShowPassword((v) => !v)}
          >
            {showPassword ? <FiEyeOff /> : <FiEye />}
          </button>
        </div>
        <div className="relative mb-3 !mt-3">
          <Input
            name="confirmPassword"
            type={showConfirmPassword ? "text" : "password"}
            placeholder="Confirm Password"
            value={formData.confirmPassword}
            onChange={handleInputChange}
            className="h-[44px] focus:border-2 focus:border-blue-400 pr-10"
            disabled={isSubmitting}
            required
          />
          <button
            type="button"
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-700"
            tabIndex={-1}
            onClick={() => setShowConfirmPassword((v) => !v)}
          >
            {showConfirmPassword ? <FiEyeOff /> : <FiEye />}
          </button>
        </div>

        <label className="text-sm flex gap-2 !mt-5">
          <input
            type="checkbox"
            name="acceptTerms"
            className="accent-orange-500 h-5 w-5 checked:bg-orange-500 checked:text-white"
            checked={formData.acceptTerms}
            onChange={handleInputChange}
            disabled={isSubmitting}
          />
          <span className="text-gray-600">
            Are you agree to Dealport{" "}
            <a className="text-blue-500 cursor-pointer hover:text-blue-700">Terms of Condition</a>{" "}
            and{" "}
            <a className="text-blue-500 cursor-pointer hover:text-blue-700">Privacy Policy</a>.
          </span>
        </label>

        <Button
          type="submit"
          className="font-sans !mt-[2rem] w-full bg-orange-500 text-white hover:bg-orange-600"
          disabled={isSubmitting}
        >
          {isSubmitting ? "SIGNING UP..." : "SIGN UP"}
        </Button>

        <div className="text-center text-gray-500">or</div>

        <Button
          type="button"
          variant="outline"
          className="!mb-5 h-[44px] w-full flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
          disabled={isSubmitting}
        >
          <FcGoogle className="!h-[20px] !w-[20px]" /> Sign up with Google
        </Button>

        <Button
          type="button"
          variant="outline"
          className="w-full h-[44px] flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
          disabled={isSubmitting}
        >
          <FaApple className="!h-[20px] !w-[20px]" /> Sign up with Apple
        </Button>
      </form>
    </div>
  );
}

export default Register;