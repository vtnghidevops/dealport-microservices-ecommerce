import { useState, FormEvent } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FcGoogle } from "react-icons/fc";
import { FaApple } from "react-icons/fa";
import { Link } from "react-router-dom";

interface RegisterFormData {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

const Register: React.FC = () => {
  const [formData, setFormData] = useState<RegisterFormData>({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [agree, setAgree] = useState<boolean>(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleRegister = (e: FormEvent): void => {
    e.preventDefault();
    if (!agree) {
      alert("Please agree to terms.");
      return;
    }
    if (formData.password !== formData.confirmPassword) {
      alert("Passwords do not match.");
      return;
    }
    console.log("Registering", {
      name: formData.name,
      email: formData.email,
      password: formData.password
    });
    // TODO: Call API to create account
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
          name="name"
          placeholder="Name"
          value={formData.name}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        <Input
          name="email"
          type="email"
          placeholder="Email Address"
          value={formData.email}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        <Input
          name="password"
          type="password"
          placeholder="8+ characters"
          value={formData.password}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        <Input
          name="confirmPassword"
          type="password"
          placeholder="Confirm Password"
          value={formData.confirmPassword}
          onChange={handleInputChange}
          className="!mb-3 h-[44px] focus:border-2 focus:border-blue-400"
        />
        
        <label className="text-sm flex gap-2">
          <input
            type="checkbox"
            className="accent-orange-500 h-5 w-5 checked:bg-orange-500 checked:text-white"
            checked={agree}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void => setAgree(e.target.checked)}
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
        >
          SIGN UP
        </Button>

        <div className="text-center text-gray-500">or</div>

        <Button
          type="button"
          variant="outline"
          className="!mb-5 h-[44px] w-full flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
        >
          <FcGoogle className="!h-[20px] !w-[20px]" /> Sign up with Google
        </Button>

        <Button
          type="button"
          variant="outline"
          className="w-full h-[44px] flex items-center gap-2 justify-center hover:bg-gray-100 transition-colors"
        >
          <FaApple className="!h-[20px] !w-[20px]" /> Sign up with Apple
        </Button>
      </form>
    </div>
  );
}

export default Register;