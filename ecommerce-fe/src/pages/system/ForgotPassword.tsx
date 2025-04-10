import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState("");

  const handleSendCode = () => {
    console.log("Sending reset code to", email);
    // TODO: Gửi mã đến email
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <h2 className="text-[22px] font-semibold mb-4">Forgot Password</h2>
      
      <p className="text-sm text-gray-600 mb-[1rem]">
        Enter the email address associated with your Dealport account.
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
        >
          SEND CODE
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