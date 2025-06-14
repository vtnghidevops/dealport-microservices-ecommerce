import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { resetPassword } from "@/services/auth/auth.service";
import { useToast } from "@/hooks/use-toast";
import { Loader2, Eye, EyeOff } from "lucide-react";
import { extractErrorMessage } from "@/utils/error-handler";

interface LocationState {
  email?: string;
  token?: string;
  verified?: boolean;
}

const ResetPassword: React.FC = () => {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [token, setToken] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { toast } = useToast();

  useEffect(() => {
    // Extract token and email from URL parameters
    const queryParams = new URLSearchParams(location.search);
    const tokenParam = queryParams.get("token");
    const emailParam = queryParams.get("email");

    // Also check location state (coming from OTPVerification)
    const state = location.state as LocationState;

    // Set email - prefer URL parameter over state
    if (emailParam) {
      setEmail(emailParam);
    } else if (state?.email) {
      setEmail(state.email);
    }

    // Set token - prefer URL parameter over state
    if (tokenParam) {
      setToken(tokenParam);
    } else if (state?.token) {
      setToken(state.token);
    }

    console.log(
      "Reset Password Component - Email:",
      emailParam || state?.email
    );
    console.log(
      "Reset Password Component - Token exists:",
      !!tokenParam || !!state?.token
    );

    // Check if we have both email and token
    const hasToken = !!tokenParam || !!state?.token;
    const hasEmail = !!emailParam || !!state?.email;

    if (!hasToken || !hasEmail) {
      toast({
        title: "Invalid Reset Link",
        description: "The password reset link is invalid or expired.",
        variant: "destructive",
      });
    }
  }, [location, toast]);

  const validatePassword = (password: string): boolean => {
    // At least 8 characters
    if (password.length < 8) {
      toast({
        title: "Password too short",
        description: "Password must be at least 8 characters long.",
        variant: "destructive",
      });
      return false;
    }

    // Should contain at least one number
    if (!/\d/.test(password)) {
      toast({
        title: "Password too weak",
        description: "Password must contain at least one number.",
        variant: "destructive",
      });
      return false;
    }

    // Should contain at least one uppercase letter
    if (!/[A-Z]/.test(password)) {
      toast({
        title: "Password too weak",
        description: "Password must contain at least one uppercase letter.",
        variant: "destructive",
      });
      return false;
    }

    return true;
  };

  const handleReset = async () => {
    if (password !== confirmPassword) {
      toast({
        title: "Passwords do not match",
        description: "Please make sure both passwords are identical.",
        variant: "destructive",
      });
      return;
    }

    if (!validatePassword(password)) {
      return;
    }

    if (!email || !token) {
      toast({
        title: "Missing information",
        description: "Reset link is invalid or expired.",
        variant: "destructive",
      });
      return;
    }

    try {
      setLoading(true);

      await resetPassword(email, password, token);

      toast({
        variant: "success",
        title: "Password updated",
        description:
          "Your password has been reset successfully. You can now log in with your new password.",
      });

      // Redirect to login page after successful password reset
      navigate("/login");
    } catch (error: unknown) {
      console.error("Password reset failed:", error);
      toast({
        title: "Password reset failed",
        description:
          extractErrorMessage(error) ||
          "An error occurred during password reset. Please try again.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <h2 className="text-[22px] font-semibold mb-4">Reset Password</h2>
      <p className="text-sm text-gray-600 mb-[1rem]">
        Create a new password for your SapoGo account. Please choose a strong
        password to protect your account.
      </p>

      <div className="space-y-4">
        <div>
          <div className="relative">
            <Input
              type={showPassword ? "text" : "password"}
              placeholder="8+ characters"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mb-1 h-[44px] focus:border-2 focus:border-blue-400 pr-10"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-gray-700"
            >
              {showPassword ? (
                <EyeOff className="h-5 w-5" />
              ) : (
                <Eye className="h-5 w-5" />
              )}
            </button>
          </div>
          <p className="text-xs text-gray-500 mt-1">
            Password must be at least 8 characters with one number and one
            uppercase letter.
          </p>
        </div>

        <div className="relative">
          <Input
            type={showConfirmPassword ? "text" : "password"}
            placeholder="Confirm Password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400 pr-10"
          />
          <button
            type="button"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-gray-700"
          >
            {showConfirmPassword ? (
              <EyeOff className="h-5 w-5" />
            ) : (
              <Eye className="h-5 w-5" />
            )}
          </button>
        </div>

        <Button
          onClick={handleReset}
          className="!mb-3 w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px] font-sans"
          disabled={loading}
        >
          {loading ? (
            <span className="flex items-center justify-center">
              <Loader2 className="mr-2 h-4 w-4 animate-spin" /> RESETTING...
            </span>
          ) : (
            "RESET PASSWORD"
          )}
        </Button>

        <div className="flex justify-center pt-4">
          <Link
            to="/login"
            className="text-blue-500 hover:text-blue-700 text-sm"
          >
            Back to Sign In
          </Link>
        </div>
      </div>
    </div>
  );
};
export default ResetPassword;
