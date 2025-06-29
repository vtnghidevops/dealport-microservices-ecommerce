import React, { useState } from "react";
import { FiSearch, FiBell } from "react-icons/fi";
import { useNavigate } from "react-router-dom";

interface AdminHeaderProps {
  userName?: string;
  userAvatar?: string;
  title: string;
}

const AdminHeader: React.FC<AdminHeaderProps> = ({
  userName = "Admin",
  userAvatar = "/images/common/avatars/admin.png",
  title = "Dashboard",
}) => {
  const [searchQuery, setSearchQuery] = useState("");
  const navigate = useNavigate();

  const handleAvatarClick = () => {
    navigate('/admin/role');
  };

  return (
    <header className="bg-white h-[96px] px-[1rem] shadow-sm border-b border-gray-200 py-3 flex items-center justify-between border-t">
      {/* Left section */}
      <div className="flex items-center gap-4">
        <h1 className="font-bold text-[22px] text-gray-800 hidden md:block px-[1rem]">
          {title}
        </h1>
      </div>

      {/* Center section - Search */}
      <div className="mx-4 h-full w-[600px]">
        <div className="relative h-[48px] w-full flex items-center">
          <div className="absolute rounded-[25px] w-[390px] border border-gray-300 inset-y-0 left-0 pl-3 flex items-center">
            <input
              type="text"
              className="block ml-[1rem] w-[85%] h-[36px] pl-4 pr-10 py-2 rounded-[25px] bg-white placeholder-gray-500 outline-none transition focus:outline-none shadow-none"
              placeholder="Search data, users, or reports"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            ></input>
            <FiSearch className="h-[1.2rem] w-[1.2rem]" />
          </div>
          {/* Right section */}
          <div className="flex items-center gap-3 w-[180px] absolute right-0">
            {/* Notifications */}
            <button
              className="relative p-1 rounded-full text-gray-500 hover:text-gray-700 hover:bg-gray-100 focus:outline-none "
              aria-label="Notifications"
            >
              <FiBell className="h-5 w-5" />
              <span className="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-500 ring-2 ring-white"></span>
            </button>

            {/* Theme Toggle */}
            <button
              className="p-2 rounded-full text-gray-500 hover:text-gray-700 hover:bg-gray-200 focus:outline-none transition duration-200 ease-in-out"
              aria-label="Toggle theme"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z"
                  clipRule="evenodd"
                />
              </svg>
            </button>

            {/* User avatar */}
            <div className="relative ml-2">
              <button
                onClick={handleAvatarClick}
                className="flex items-center focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 rounded-full transition-all duration-300 hover:border-blue-400 cursor-pointer"
                aria-label="User profile"
                title="Go to profile"
              >
                <img
                  className="h-[2rem] w-[2rem] rounded-full border-2 border-gray-200 hover:border-blue-400"
                  src={userAvatar}
                  alt={userName}
                />
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default AdminHeader;
