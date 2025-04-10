import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

const ResetPassword: React.FC = () => {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const handleReset = () => {
    if (password !== confirmPassword) {
      alert("Passwords do not match!");
      return;
    }
    console.log("Resetting password to", password);
    // TODO: Gửi mật khẩu mới đến backend
  };

  return (
    <div className="my-[2rem] max-w-sm mx-auto p-[2rem] shadow-lg rounded-xl border w-[424px]">
      <h2 className="text-[22px] font-semibold mb-4">Reset Password</h2>
      <p className="text-sm text-gray-600 mb-[1rem]">
        Create a new password for your Dealport account.
        Please choose a strong password to protect your account.
      </p>

      <div className="space-y-4">
        <Input
          type="password"
          placeholder="8+ characters"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        <Input
          type="password"
          placeholder="Confirm Password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        <Button 
          onClick={handleReset} 
          className="!mb-3 w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px] font-sans"
        >
          RESET PASSWORD
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
}
export default ResetPassword;