import React, { useEffect, useState, useRef, useContext } from 'react';
import { FaUser, FaRegHeart } from "react-icons/fa";
import { IoCartOutline } from "react-icons/io5";
// import { IoIosLogOut } from "react-icons/io";
import { TbSettings } from "react-icons/tb";
import { MdOutlineArrowDropDown } from "react-icons/md";
import { FiLogOut } from "react-icons/fi";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from '@/hooks/useAuth';
import { useToast } from '@/hooks/use-toast';
import { HeaderContext } from './Header';

const LoggedInUserMenu: React.FC = () => {
  const { authState, logout, logoutFromAllDevices } = useAuth();
  const { toast } = useToast();
  const navigate = useNavigate();
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const headerContext = useContext(HeaderContext);

  const user = authState.user;
  const isAdmin = user?.role === 'admin';

  // console.log("LoggedInUserMenu - User data:", user);
  // console.log("LoggedInUserMenu - Is admin:", isAdmin);
  // console.log("LoggedInUserMenu - User profile structure:", user?.profile);

  // Add debugging log for component rendering
  // useEffect(() => {
  //   console.log("LoggedInUserMenu rendered - User data:", user);
  //   console.log("LoggedInUserMenu - Component mounted with auth state:", authState);
  // }, []);

  // Handle click outside to close menu
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        menuRef.current &&
        buttonRef.current &&
        !menuRef.current.contains(event.target as Node) &&
        !buttonRef.current.contains(event.target as Node)
      ) {
        setShowMenu(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleLogout = () => {
    console.log("LoggedInUserMenu: handleLogout called");
    logout();
    console.log("LoggedInUserMenu: logout function completed");
    toast({
      variant: "success",
      title: "Success",
      description: "Logged out successfully from current session"
    });
    setShowMenu(false);

    // Update header state on logout
    if (headerContext) {
      headerContext.setForceRerender(headerContext.forceRerender + 1);
    }

    navigate('/');
  };

  const handleLogoutAllDevices = () => {
    console.log("LoggedInUserMenu: handleLogoutAllDevices called");
    logoutFromAllDevices();
    console.log("LoggedInUserMenu: logoutFromAllDevices function completed");
    toast({
      variant: "success",
      title: "Success",
      description: "Logged out successfully from all devices"
    });
    setShowMenu(false);

    // Update header state on logout
    if (headerContext) {
      headerContext.setForceRerender(headerContext.forceRerender + 1);
    }

    navigate('/');
  };

  const handleAdminDashboard = () => {
    navigate('/admin/dashboard');
    setShowMenu(false);
  };

  const toggleMenu = () => {
    setShowMenu(prev => !prev);
  };

  return (
    <div className="relative">
      <button
        ref={buttonRef}
        className="text-[1rem] h-full flex items-center hover:text-green-600 transition-colors"
        onClick={toggleMenu}
        onMouseEnter={() => setShowMenu(true)}
      >
        <FaUser className="mr-1" />
        <span className="mr-1">{user?.profile?.firstName || "User"}</span>
        <MdOutlineArrowDropDown />
      </button>

      {/* Invisible connector to prevent mouse-out */}
      <div
        className="absolute left-[-40px] w-[120px] h-[25px] bottom-[-25px] bg-transparent z-50"
        onMouseEnter={() => setShowMenu(true)}
      ></div>

      {/* Menu dropdown */}
      {showMenu && (
        <div
          ref={menuRef}
          className="absolute -right-[70px] mt-[25px] w-[250px] bg-white rounded-md shadow-lg z-50 border border-gray-200 flex flex-col items-center justify-center transition-all duration-200"
          onMouseLeave={() => setShowMenu(false)}
        >
          {/* Triangle pointer connecting to button */}
          <div className="absolute -top-2 right-[76px] w-4 h-4 bg-white border-t border-l border-gray-200 transform rotate-45"></div>

          <div className="border-b border-gray-200 w-full py-3 px-5 text-center">
            <p className="font-medium text-[15px]">Welcome, <span className="text-success">{user?.profile?.firstName || "User"}</span></p>
            <p className="text-xs text-gray-500">{user?.email || ""}</p>
            <p className="text-xs text-gray-500">Role: <span className="font-semibold">{user?.role || "user"}</span></p>
          </div>

          <Link
            to="/user/profile"
            className="border-b border-gray-200 w-full h-[50px] flex justify-start items-center pl-5 px-4 py-2 text-[16px] text-gray-700 hover:bg-gray-100"
            onClick={() => setShowMenu(false)}
          >
            <FaUser className="mr-2" /> Profile
          </Link>

          {/* Hiển thị My Orders và Wishlist chỉ khi không phải admin */}
          {!isAdmin && (
            <>
              <Link
                to="/user/orders"
                className="border-b border-gray-200 w-full h-[50px] justify-start  pl-5 px-4 py-2 text-[16px] text-gray-700 hover:bg-gray-100 flex items-center"
                onClick={() => setShowMenu(false)}
              >
                <IoCartOutline className="mr-2" /> My Orders
              </Link>

              <Link
                to="/user/wishlist"
                className="border-b border-gray-200 w-full h-[50px] justify-start  pl-5 px-4 py-2 text-[16px] text-gray-700 hover:bg-gray-100 flex items-center"
                onClick={() => setShowMenu(false)}
              >
                <FaRegHeart className="mr-2" /> Wishlist
              </Link>
            </>
          )}

          {/* Admin Management option - only shown if user is admin */}
          {isAdmin && (
            <button
              onClick={handleAdminDashboard}
              className="border-b border-gray-200 w-full h-[50px] justify-start pl-5 px-4 py-2 text-[16px] text-primary hover:bg-gray-100 flex items-center"
            >
              <TbSettings className="mr-2" /> Admin Management
            </button>
          )}

          <button
            onClick={handleLogout}
            className="border-b border-gray-200 w-full h-[50px] justify-start pl-5 px-4 py-2 text-[16px] text-error hover:bg-gray-100 flex items-center"
          >
            <FiLogOut className="mr-2" /> Logout Current Session
          </button>

          <button
            onClick={handleLogoutAllDevices}
            className="w-full h-[50px] justify-start pl-5 px-4 py-2 text-[16px] text-error hover:bg-gray-100 flex items-center"
          >
            <FiLogOut className="mr-2" /> Logout All Devices
          </button>
        </div>
      )}
    </div>
  );
};

export default LoggedInUserMenu; 