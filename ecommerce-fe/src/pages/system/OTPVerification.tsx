import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { verifyOTP, resendOTP } from "@/services/auth/auth.service";
import { useToast } from "@/hooks/use-toast";

const OTPVerification: React.FC = () => {
  const [otp, setOtp] = useState("");
  const [email, setEmail] = useState("");
  const [remainingTime, setRemainingTime] = useState(0);
  const [purpose, setPurpose] = useState<string>(""); // registration or password_reset
  const [loading, setLoading] = useState(false);
  const [resending, setResending] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { toast } = useToast();

  useEffect(() => {
    // Get data from location state
    const state = location.state as {
      email?: string;
      purpose?: string;
      expiresIn?: number;
    } | null;

    if (state) {
      setEmail(state.email || "");
      setPurpose(state.purpose || "");
      setRemainingTime(state.expiresIn || 0);
    } else {
      // If no state is provided, redirect to homepage
      toast({
        variant: "destructive",
        title: "Invalid Navigation",
        description: "Please follow the normal registration or password reset flow"
      });
      navigate("/");
    }
  }, [location, navigate, toast]);

  useEffect(() => {
    if (remainingTime <= 0) return;

    const timer = setInterval(() => {
      setRemainingTime(prev => Math.max(0, prev - 1));
    }, 60000); // decrease every minute

    return () => clearInterval(timer);
  }, [remainingTime]);

  const handleVerify = async () => {
    if (otp.length !== 6) {
      toast({
        variant: "destructive",
        title: "Invalid OTP",
        description: "Please enter a valid 6-digit OTP"
      });
      return;
    }

    setLoading(true);
    try {
      const response = await verifyOTP(email, otp, purpose);
      console.log("OTP Verification Response:", response);

      if (response.success || response.data?.success) {
        toast({
          variant: "success",
          title: "Verification Success",
          description: purpose === "registration"
            ? "Your account has been verified. You can now log in."
            : "Your identity has been verified. You can now reset your password."
        });

        // For password reset, extract token and navigate to reset password page
        if (purpose === "password_reset") {
          // Extract token from the response data
          let token = "";

          // Handle different response formats
          if (response.data && response.data.token) {
            token = response.data.token;
          } else if (response.data && response.data.reset_token) {
            token = response.data.reset_token;
          } else if (response.token) {
            token = response.token;
          } else if (response.reset_token) {
            token = response.reset_token;
          }

          console.log("Token from verify OTP:", token);

          if (!token) {
            console.error("No token returned from verification!");
            toast({
              variant: "destructive",
              title: "Error",
              description: "Could not complete verification. Please try again."
            });
            return;
          }

          // Normalize email to lowercase for consistency
          const normalizedEmail = email.toLowerCase().trim();

          // Navigate with both URL params and location state
          // URL params allow bookmarking or accessing directly
          // Location state is used by the ResetPassword component if available
          navigate(`/reset-password?email=${encodeURIComponent(normalizedEmail)}&token=${encodeURIComponent(token)}`, {
            state: {
              email: normalizedEmail,
              verified: true,
              token: token
            }
          });
        } else {
          // For registration, just navigate to login page
          navigate('/login', {
            state: {
              message: "Your account has been verified. You can now log in."
            }
          });
        }
      } else {
        console.error("OTP verification failed", response);
        toast({
          variant: "destructive",
          title: "Verification Failed",
          description: response.response?.data?.message || "Invalid or expired OTP. Please try again."
        });
      }
    } catch (error: any) {
      console.error("OTP verification failed", error);
      toast({
        variant: "destructive",
        title: "Verification Failed",
        description: error.response?.data?.message || "Invalid or expired OTP. Please try again."
      });
    } finally {
      setLoading(false);
    }
  };

  const handleResendOTP = async () => {
    setResending(true);
    try {
      const response = await resendOTP(email, purpose);
      setRemainingTime(response.expiresIn || 10);

      toast({
        variant: "success",
        title: "OTP Resent",
        description: "A new verification code has been sent to your email"
      });
    } catch (error: any) {
      console.error("Failed to resend OTP", error);
      toast({
        variant: "destructive",
        title: "Failed to Resend OTP",
        description: error.response?.data?.message || "Please try again later."
      });
    } finally {
      setResending(false);
    }
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <h2 className="text-[22px] font-semibold mb-4">OTP Verification</h2>

      <p className="text-sm text-gray-600 mb-[1rem]">
        Enter the verification code sent to {email || "your email"}
      </p>

      <div className="space-y-4">
        <Input
          type="text"
          placeholder="Enter 6-digit code"
          value={otp}
          onChange={(e) => setOtp(e.target.value.replace(/[^0-9]/g, '').slice(0, 6))}
          className="mb-3 h-[44px] focus:border-2 focus:border-blue-400"
          maxLength={6}
        />

        <Button
          onClick={handleVerify}
          className="!mb-3 w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px] font-sans"
          disabled={otp.length !== 6 || loading}
        >
          {loading ? "VERIFYING..." : "VERIFY"}
        </Button>

        <div className="text-center text-sm">
          {remainingTime > 0 ? (
            <p className="text-gray-600">
              Code expires in {remainingTime} {remainingTime === 1 ? 'minute' : 'minutes'}
            </p>
          ) : (
            <p className="text-red-500">Code has expired</p>
          )}
        </div>

        <div className="flex justify-between items-center pt-4">
          <button
            onClick={handleResendOTP}
            className="text-blue-500 hover:text-blue-700 text-sm disabled:text-gray-400"
            disabled={remainingTime > 0 || resending}
          >
            {resending ? "Sending..." : "Resend Code"}
          </button>

          <Link
            to={purpose === "password_reset" ? "/forgot-password" : "/register"}
            className="text-blue-500 hover:text-blue-700 text-sm"
          >
            Back
          </Link>
        </div>
      </div>
    </div>
  );
};

export default OTPVerification; 