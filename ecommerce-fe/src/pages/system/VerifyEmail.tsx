import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

const VerifyEmail: React.FC = () => {
  const [code, setCode] = useState("");

  const handleVerify = () => {
    console.log("Verifying code:", code);
    // TODO: Gọi API xác minh mã
  };

  const handleResend = () => {
    console.log("Resending code...");
    // TODO: Gửi lại mã xác thực
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <h2 className="text-[22px] font-semibold mb-4">Verify Your Email</h2>
      <p className="text-sm text-gray-600 mb-[1rem]">
        Please enter the verification code that was sent to your email address.
        The code will expire in 10 minutes.
      </p>

      <div className="space-y-4">
        <div className="relative">
          <Input
            placeholder="Enter verification code"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            className="mb-3 h-[44px] focus:border-2 focus:border-blue-400 pr-24"
          />
          <button
            onClick={handleResend}
            className="absolute top-[10px] right-3 text-sm text-blue-500 hover:text-blue-700"
          >
            Resend Code
          </button>
        </div>

        <Button 
          onClick={handleVerify}
          className="!mb-3 w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px] font-sans"
        >
          VERIFY EMAIL
        </Button>

        <div className="flex justify-between items-center pt-4">
          <Link 
            to="/login" 
            className="text-blue-500 hover:text-blue-700 text-sm"
          >
            Back to Sign In
          </Link>
          <button 
            className="text-blue-500 hover:text-blue-700 text-sm"
            onClick={handleResend}
          >
            Need help?
          </button>
        </div>
      </div>
    </div>
  );
}
export default VerifyEmail;