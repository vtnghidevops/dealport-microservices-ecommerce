import React, { useState, useEffect, useRef } from "react";
import { MdOutlineArrowDropDown } from "react-icons/md";
import { FaUser, FaShoppingCart, FaRegHeart } from "react-icons/fa";
import { TbMinusVertical, TbSettings } from "react-icons/tb";
import { FiSearch, FiExternalLink, FiEye, FiEyeOff } from "react-icons/fi";
import { IoMenu, IoTicketOutline } from "react-icons/io5";
import { TiHome } from "react-icons/ti";
import { CiShop, CiHeart } from "react-icons/ci";
import { LuUsersRound } from "react-icons/lu";
import { HiOutlineStar } from "react-icons/hi";
import { FaLocationDot } from "react-icons/fa6";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from '@/hooks/useAuth';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import LoggedInUserMenu from './LoggedInUserMenu';
import { useCart } from "@/hooks/useCart";
import { extractErrorMessage } from '@/utils/error-handler';

const Logo: React.FC = () => {
  return (
    <a href="/">
      <img src="/images/common/logo.png" alt="logo" className="h-7" />
    </a>
  );
};

const SearchBar: React.FC = () => {
  return (
    <div className="flex items-center justify-between bg-aqua-spring rounded-full px-4 py-2 w-[25rem]">
      <input
        type="text"
        placeholder="What you're looking for"
        className="bg-transparent focus:outline-none ml-[1rem] w-full"
      />
      <button className="flex items-center justify-center rounded-[20px] text-black w-[9rem] bg-white mr-1  h-[2rem] hover:shadow-md hover:scale-105 transition-transform transition-shadow duration-300">
        <FiSearch className="mr-1 " />
        Search
      </button>
    </div>
  );
};

/**
 * PART 1: USER MENU WHEN NOT LOGGED IN
 * Display a simple Login button that redirects to the login page when clicked
 */
const LoggedOutUserMenu: React.FC = () => {
  const navigate = useNavigate();
  const [showLoginForm, setShowLoginForm] = useState(false);
  const [formData, setFormData] = useState({
    email: "",
    password: ""
  });
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const { toast } = useToast();

  const loginFormRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Access Header's forceRerender via React Context
  const headerContext = React.useContext(HeaderContext);

  const [showPassword, setShowPassword] = useState(false);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        loginFormRef.current &&
        buttonRef.current &&
        !loginFormRef.current.contains(event.target as Node) &&
        !buttonRef.current.contains(event.target as Node)
      ) {
        setShowLoginForm(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [loginFormRef, buttonRef]);

  const handleMouseLeave = () => {
    // Don't hide form during login attempt
    if (isLoading) return;
    setShowLoginForm(false);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!formData.email || !formData.password) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Please fill in all required fields"
      });
      return;
    }

    setIsLoading(true);
    const result = await login(formData.email, formData.password);

    if (result.success) {
      toast({
        variant: "success",
        title: "Success",
        description: "Login successful!"
      });
      setShowLoginForm(false);
      window.dispatchEvent(new Event('storage'));
      if (headerContext) {
        headerContext.setForceRerender(headerContext.forceRerender + 1);
      }
      if (window.location.pathname === '/login') {
        navigate('/');
      }
    } else {
      const friendlyMsg = extractErrorMessage(result.error);
      toast({
        variant: "destructive",
        title: "Login Failed",
        description: friendlyMsg
      });
    }
    setIsLoading(false);
  };

  return (
    <div className="relative" ref={containerRef} onMouseLeave={handleMouseLeave}>
      <button
        ref={buttonRef}
        className="text-[1rem] h-full flex items-center hover:text-green-600 transition-colors"
        onMouseEnter={() => setShowLoginForm(true)}
        onClick={() => setShowLoginForm(true)}
      >
        <FaUser className="mr-1" />
        <MdOutlineArrowDropDown />
        <div className="absolute -left-[80%] w-[100px] h-[25px] bottom-[-25px] bg-transparent z-50"></div>
      </button>

      {/* Login form popup */}
      {showLoginForm && (
        <div
          ref={loginFormRef}
          className="absolute -right-[70px] mt-[25px] w-[424px] bg-white rounded-xl shadow-lg z-50 border"
        >
          {/* Triangle pointer connecting to button */}
          <div className="absolute -top-2 right-[76px] w-4 h-4 bg-white border-t border-l border-gray-200 transform rotate-45"></div>
          <div className="px-[2rem] pt-[1.5rem] text-center">
            <h3 className="text-lg font-medium text-gray-800">You have an account! Login Now</h3>
            <p className="text-sm text-gray-500 mt-1">Enter your credentials to access your account</p>
          </div>
          <div className="p-[2rem] pt-[1rem]">
            <form onSubmit={handleSubmit} className="space-y-4">
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
                  onClick={() => setShowPassword((v) => !v)}
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
                disabled={isLoading}
                className="font-sans !mt-[2.5rem] w-full bg-orange-500 text-white hover:bg-orange-600 h-[44px]"
              >
                {isLoading ? "SIGNING IN..." : "SIGN IN"}
              </Button>

              <div className="text-sm text-center mt-4">
                <span className="text-gray-600">Don't have an account? </span>
                <Link to="/register" className="text-blue-500 hover:text-blue-700">
                  Register here
                </Link>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

const NavMenu: React.FC = () => {
  const categories = [
    { name: "Electronics", href: "/category/electronics", icon: <FiExternalLink className="mr-2" /> },
    { name: "Clothing", href: "/category/clothing", icon: <CiShop className="mr-2" /> },
    { name: "Home & Garden", href: "/category/home-garden", icon: <TiHome className="mr-2" /> },
    { name: "Beauty", href: "/category/beauty", icon: <HiOutlineStar className="mr-2" /> },
    { name: "Sports", href: "/category/sports", icon: <IoTicketOutline className="mr-2" /> }
  ];

  return (
    <div className="relative group">
      <button
        className="flex items-center"
      >
        <IoMenu className="text-xl mr-1" />
        Menu
        <span className="mx-4 text-gray-400">|</span>
      </button>

      {/* Invisible connector to prevent mouse-out */}
      <div className="absolute left-[-40px] w-[120px] h-[25px] bottom-[-25px] bg-transparent z-50"></div>

      <div className="w-[200px] absolute left-0 mt-4 bg-white rounded-md shadow-lg py-2 z-50 border border-gray-200 hidden group-hover:block transition-all duration-200">
        {/* Triangle pointer connecting to button */}
        <div className="absolute -top-2 left-6 w-4 h-4 bg-white border-t border-l border-gray-200 transform rotate-45"></div>

        {categories.map((category, index) => (
          <Link
            key={index}
            to={category.href}
            className="text-[16px] pl-5 px-4 py-2 text-gray-700 hover:bg-gray-100 flex items-center"
          >
            {category.icon} {category.name}
          </Link>
        ))}
        <div className="border-t border-gray-100 my-1"></div>
        <Link
          to="/new-arrivals"
          className="pl-5 px-4 py-2 text-[16px] font-medium text-green-600 hover:bg-gray-100 flex items-center"
        >
          <CiHeart className="mr-2" /> New Arrivals
        </Link>
        <Link
          to="/sale"
          className="pl-5 px-4 py-2 text-[16px] font-medium text-red-600 hover:bg-gray-100 flex items-center"
        >
          <TbSettings className="mr-2" /> On Sale
        </Link>
      </div>
    </div>
  );
};

// Create Context for sharing state between Header and subcomponents
export const HeaderContext = React.createContext<{
  forceRerender: number;
  setForceRerender: React.Dispatch<React.SetStateAction<number>>;
  checkAuthStatus: () => void;
} | null>(null);

const Header: React.FC = () => {
  const { cartItems } = useCart();
  const cartItemCount = cartItems.length;
  const { authState } = useAuth();
  const { isAuthenticated, isLoading, user } = authState;
  const [forceRerender, setForceRerender] = useState(0);
  // const [showDebug, setShowDebug] = useState(false);
  // const navigate = useNavigate();

  // Check localStorage directly
  useEffect(() => {
    checkLoginStatus();

    // Listen for storage events which might indicate login state changes
    const handleStorageChange = () => {
      console.log("Storage change detected - refreshing header");
      setForceRerender(prev => prev + 1);
      checkLoginStatus();
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, [isAuthenticated]);

  const checkLoginStatus = () => {
    const token = localStorage.getItem('token');
    const userString = localStorage.getItem('user');

    // If we have a token in localStorage but isAuthenticated is false, force a reload
    if (token && userString && !isAuthenticated) {
      console.log("Header detected token in localStorage but not in context - forcing context update");
      // Instead of reload, try to update the auth state context by dispatching an event
      window.dispatchEvent(new Event('storage'));
    }

    console.log("Header directly checking localStorage:", {
      hasToken: !!token,
      hasUser: !!userString,
      isAuthenticated,
      userData: userString ? JSON.parse(userString) : null
    });

    setForceRerender(prev => prev + 1);
  };

  // Debug code để kiểm tra trạng thái đăng nhập
  // console.log("Header Auth State:", {
  //   isAuthenticated,
  //   isLoading,
  //   hasUser: !!user,
  //   user: user,
  //   forceRerender
  // });

  // Determine if logged in by checking both context and localStorage
  const token = localStorage.getItem('token');
  const userString = localStorage.getItem('user');
  const isLoggedIn = isAuthenticated || (!!token && !!userString);

  // If we have data in localStorage but auth context hasn't updated yet, parse from localStorage
  let parsedUser = user;
  if (!user && userString) {
    try {
      parsedUser = JSON.parse(userString);
      console.log("Using user data from localStorage:", parsedUser);
    } catch (e) {
      console.error("Failed to parse user from localStorage:", e);
    }
  }

  // Debug the menu state
  // console.log("Menu rendering state:", {
  //   isLoggedIn,
  //   token: token ? "exists" : "none",
  //   userInStorage: userString ? "exists" : "none"
  // });

  return (
    <HeaderContext.Provider value={{
      forceRerender,
      setForceRerender,
      checkAuthStatus: checkLoginStatus
    }}>
      {/* {process.env.NODE_ENV === 'development' && showDebug && (
        <div className="bg-yellow-100 p-2 text-xs border-b border-yellow-300">
          <div className="flex justify-between items-center">
            <div>
              <span className="font-bold">Auth Debug:</span>
              {isAuthenticated ?
                <span className="text-green-700 mx-1">Authenticated</span> :
                <span className="text-red-700 mx-1">Not Authenticated</span>}
              {isLoading && <span className="text-blue-700 mx-1">Loading</span>}
              {!!user && <span className="mx-1">User: {user.email}</span>}
              Token: {localStorage.getItem('token') ? '✓' : '✗'}
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => window.dispatchEvent(new Event('storage'))}
                className="bg-blue-500 text-white px-2 rounded text-xs"
              >
                Force Update
              </button>
              <button
                onClick={() => (window as any).debugAuth?.()}
                className="bg-blue-500 text-white px-2 rounded text-xs"
              >
                Console Debug
              </button>
              <button
                onClick={() => setShowDebug(false)}
                className="bg-gray-500 text-white px-2 rounded text-xs"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )} */}

      <header className="bg-white sticky top-0 z-50 flex justify-between items-center p-4 border-b-2 border-gray-200 h-[5rem] w-full">
        <div className="flex items-center space-x-4 ml-[80px] min-h-56 w-1/2">
          <Logo />
          {/* {process.env.NODE_ENV === 'development' && !showDebug && (
            <button
              onClick={() => setShowDebug(true)}
              className="text-xs text-gray-400 hover:text-gray-600 absolute top-1 left-1"
            >
              Debug
            </button>
          )} */}
          <span className="flex items-center space-x-1 h-[3rem]">
            <TbMinusVertical className="h-[100%] text-[30px]" />
            <a href="/" className="flex items-center space-x-2 max-w-[9rem]">
              <FaLocationDot className="text-[20px] text-red-500" />
              <span className="text-[15px]">Deliver to Your address</span>
            </a>
          </span>
          <span className="flex items-center space-x-1 !-ml-1">
            <TbMinusVertical className="h-[100%] text-[30px]" />
            <a href="/" className="text-[15px]">
              EN
            </a>
            <MdOutlineArrowDropDown />
          </span>
        </div>

        <div className="flex items-center justify-end mr-[80px] space-x-5 min-h-56 w-1/2 p-1 h-[3rem]">
          <SearchBar />

          {/* Hiển thị UserMenu tùy theo trạng thái đăng nhập - thêm debug class */}
          <div className={`auth-state-${isLoggedIn ? 'authenticated' : 'unauthenticated'} rerender-${forceRerender}`}>
            {isLoading ? (
              <div className="h-[40px] flex items-center">
                <span className="text-gray-500">Loading...</span>
              </div>
            ) : isLoggedIn ? (
              <LoggedInUserMenu key={`loggedin-${forceRerender}`} />
            ) : (
              <LoggedOutUserMenu />
            )}
          </div>

          <Link
            to={isLoggedIn ? "/user/orders" : "/login"}
            className="flex items-center justify-center w-[4rem] h-[3rem] hover:text-green-600 transition-colors relative"
          >
            <FaShoppingCart className="mr-2 text-[1rem]" />
            Cart
            {cartItemCount > 0 && (
              <span className="absolute -top-1 -right-3 bg-red-500 text-white text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center">
                {cartItemCount}
              </span>
            )}
          </Link>
        </div>
      </header>

      <nav className="flex justify-between items-center py-2 border-b-2 h-[3rem] max-w-[1440px]">
        {/* Left section */}
        <div className="flex items-center space-x-10 ml-[5.4rem]">
          <NavMenu />

          <a href="/explore" className="flex items-center ml-[1rem] relative group">
            <CiShop className="mr-1" /> Explore
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/deals" className="flex items-center relative group">
            <IoTicketOutline className="mr-1" /> Deals
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/saved" className="flex items-center relative group">
            <FaRegHeart className="mr-1" /> Saved
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
        </div>

        {/* Right section */}
        <div className="flex items-center space-x-10 mr-[5.4rem]">
          <a href="/" className="flex items-center relative group">
            <TiHome className="mr-1" /> Home
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/products" className="flex items-center relative group">
            <CiShop className="mr-1" /> Product
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/about" className="flex items-center relative group">
            <LuUsersRound className="mr-1" /> About Us
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/contact" className="flex items-center relative group">
            <FaLocationDot className="mr-1" /> Contact
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
        </div>
      </nav>
    </HeaderContext.Provider>
  );
}
export default Header;