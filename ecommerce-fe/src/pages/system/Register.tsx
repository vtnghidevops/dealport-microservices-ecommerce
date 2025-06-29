import { useState, FormEvent, useCallback, useMemo } from "react";
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
    firstName: "",
    lastName: "",
    username: "",
    email: "",
    phone: "",
    password: "",
    confirmPassword: "",
    acceptTerms: false,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  // const [errors, setErrors] = useState<Record<string, string>>({});
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isPasswordFocused, setIsPasswordFocused] = useState(false);

  // Memoize the input change handler to prevent recreating on every render
  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>): void => {
      const { name, value, type, checked } = e.target;
      setFormData((prev) => ({
        ...prev,
        [name]: type === "checkbox" ? checked : value,
      }));
    },
    []
  );

  // Memoize password validation function
  const validatePassword = useCallback((password: string): boolean => {
    const hasUpperCase = /[A-Z]/.test(password);
    const hasLowerCase = /[a-z]/.test(password);
    const hasNumbers = /\d/.test(password);
    const hasSpecialChar = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(
      password
    );

    return hasUpperCase && hasLowerCase && hasNumbers && hasSpecialChar;
  }, []);

  // Memoize social signup handler
  const handleSocialSignup = useCallback(
    (provider: string) => {
      toast({
        title: "Coming Soon",
        description: `${
          provider.charAt(0).toUpperCase() + provider.slice(1)
        } signup will be available soon!`,
        variant: "default",
      });
    },
    [toast]
  );

  // Memoize toggle password visibility handlers
  const togglePasswordVisibility = useCallback(() => {
    setShowPassword((prev) => !prev);
  }, []);

  const toggleConfirmPasswordVisibility = useCallback(() => {
    setShowConfirmPassword((prev) => !prev);
  }, []);

  // Memoize focus handlers
  const handlePasswordFocus = useCallback(() => {
    setIsPasswordFocused(true);
  }, []);

  const handlePasswordBlur = useCallback(() => {
    setIsPasswordFocused(false);
  }, []);

  const handleRegister = async (e: FormEvent): Promise<void> => {
    e.preventDefault();

    // Validation
    if (!formData.acceptTerms) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please agree to terms and conditions",
      });
      return;
    }
    if (formData.password !== formData.confirmPassword) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Passwords do not match",
      });
      return;
    }
    if (formData.password.length < 8) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Password must be at least 8 characters",
      });
      return;
    }

    if (!validatePassword(formData.password)) {
      toast({
        variant: "destructive",
        title: "Password Format Error",
        description:
          "Password must contain uppercase, lowercase, numbers, and special characters",
      });
      return;
    }

    if (
      !formData.email ||
      !formData.firstName ||
      !formData.lastName ||
      !formData.username
    ) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please fill in all required fields",
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
          title: "Registration Submitted",
          description:
            "Please verify your email with the OTP code sent to your inbox.",
        });

        // Redirect to OTP verification page instead of login
        navigate("/verify-otp", {
          state: {
            email: formData.email,
            purpose: "registration",
            expiresIn: 10, // 10 minutes
          },
        });
      } else if (authState.error) {
        toast({
          variant: "destructive",
          title: "Registration Failed",
          description: authState.error || "Please try again.",
        });
      }
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Registration Failed",
        description: error.message || "Please try again.",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Memoize the password requirements text to prevent recreation on each render
  const passwordRequirementsText = useMemo(
    () =>
      isPasswordFocused && (
        <div className="text-xs text-gray-600 mt-1 mb-2">
          Password must contain at least 8 characters with uppercase (viết hoa),
          lowercase (viết thường), numbers (số), and special characters (ký tự
          đặc biệt).
        </div>
      ),
    [isPasswordFocused]
  );

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <div className="flex justify-between items-center pb-2 mb-4">
        <Link
          to="/login"
          className="flex text-[18px] pb-3 items-center justify-center text-neutral-500 w-1/2 hover:text-neutral-800"
        >
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
          placeholder="Phone Number"
          value={formData.phone || ""}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          disabled={isSubmitting}
        />
        <div className="relative mb-3">
          <Input
            name="password"
            type={showPassword ? "text" : "password"}
            placeholder="Enter Your Password (8+ characters)"
            value={formData.password}
            onChange={handleInputChange}
            className="h-[44px] focus:border-2 focus:border-blue-400 pr-10"
            disabled={isSubmitting}
            required
            minLength={8}
            onFocus={handlePasswordFocus}
            onBlur={handlePasswordBlur}
            autoComplete="new-password"
            data-lpignore="true"
            data-form-type="other"
          />
          <button
            type="button"
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-700 z-10"
            tabIndex={-1}
            onClick={togglePasswordVisibility}
          >
            {showPassword ? <FiEyeOff /> : <FiEye />}
          </button>
        </div>
        {passwordRequirementsText}
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
            autoComplete="new-password"
            data-lpignore="true"
            data-form-type="other"
          />
          <button
            type="button"
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-700 z-10"
            tabIndex={-1}
            onClick={toggleConfirmPasswordVisibility}
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
            Are you agree to SapoGo{" "}
            <a className="text-blue-500 cursor-pointer hover:text-blue-700">
              Terms of Condition
            </a>{" "}
            and{" "}
            <a className="text-blue-500 cursor-pointer hover:text-blue-700">
              Privacy Policy
            </a>
            .
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
          onClick={() => handleSocialSignup("google")}
        >
          <FcGoogle className="!h-[20px] !w-[20px]" /> Sign up with Google
        </Button>

        <Button
          type="button"
          variant="outline"
          className="w-full h-[44px] flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
          disabled={isSubmitting}
          onClick={() => handleSocialSignup("apple")}
        >
          <FaApple className="!h-[20px] !w-[20px]" /> Sign up with Apple
        </Button>
      </form>
    </div>
  );
};

export default Register;
