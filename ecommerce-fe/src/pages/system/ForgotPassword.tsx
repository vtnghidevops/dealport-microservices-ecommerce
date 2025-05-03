import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Link, useNavigate } from "react-router-dom";
import { requestOTP } from "@/services/auth/auth.service";
import { useToast } from "@/hooks/use-toast";

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();

  const handleSendCode = async () => {
    if (!email) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please enter your email address"
      });
      return;
    }

    setLoading(true);
    try {
      const response = await requestOTP(email, "password_reset");

      // Navigate to OTP verification page with the email and purpose
      navigate("/verify-otp", {
        state: {
          email,
          purpose: "password_reset",
          expiresIn: response.expiresIn || 10
        }
      });

      toast({
        variant: "success",
        title: "OTP Sent",
        description: "A verification code has been sent to your email"
      });
    } catch (error: any) {
      // console.error("Failed to send reset code", error);
      toast({
        variant: "destructive",
        title: "Account not found",
        description: error.response?.data?.message || "Please try again later."
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <h2 className="text-[22px] font-semibold mb-4">Forgot Password</h2>

      <p className="text-sm text-gray-600 mb-[1rem]">
        Enter the email address associated with your account.
      </p>

      <div className="space-y-4">
        <Input
          type="email"
          placeholder="Email Address"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />

        <Button
          onClick={handleSendCode}
          className="!mb-3 w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px] font-sans"
          disabled={loading}
        >
          {loading ? "SENDING..." : "SEND CODE"}
        </Button>

        <div className="flex justify-between items-center pt-4">
          <Link
            to="/login"
            className="text-blue-500 hover:text-blue-700 text-sm"
          >
            Back to Sign In
          </Link>
          <Link
            to="/register"
            className="text-blue-500 hover:text-blue-700 text-sm"
          >
            Create Account
          </Link>
        </div>

        <div className="text-sm text-gray-600 !mt-3">
          You may contact{" "}
          <button className="text-orange-500 hover:text-orange-600">
            Customer Service
          </button>{" "}
          for help restoring access to your account.
        </div>
      </div>
    </div>
  );
}
export default ForgotPassword;