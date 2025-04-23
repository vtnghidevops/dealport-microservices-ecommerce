import React, { useState } from "react";
import { IoLocationSharp } from "react-icons/io5";
import { MdOutlineArrowDropDown } from "react-icons/md";
import { FaUser, FaShoppingCart, FaRegHeart } from "react-icons/fa";
import { TbMinusVertical, TbSettings } from "react-icons/tb";
import { FiSearch, FiExternalLink } from "react-icons/fi";
import { IoMenu, IoTicketOutline } from "react-icons/io5";
import { IoIosLogOut } from "react-icons/io";
import { useCart } from '@/hooks/useCart';
import { TiHome } from "react-icons/ti";
import { IoCartOutline } from "react-icons/io5";
import { CiShop, CiHeart } from "react-icons/ci";
import { LuUsersRound } from "react-icons/lu";
import { HiOutlineStar } from "react-icons/hi";
import { FaLocationDot } from "react-icons/fa6";

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

const UserMenu: React.FC = () => {
  return (
    <div className="relative group">
      <button
        className="text-[1rem] h-full flex items-center hover:text-green-600 transition-colors"
      >
        <FaUser className="mr-1" />
        <MdOutlineArrowDropDown />
      </button>

      {/* Invisible connector to prevent mouse-out */}
      <div className="absolute left-[-40px] w-[120px] h-[25px] bottom-[-25px] bg-transparent z-50"></div>


      <div className="absolute -right-[70px] mt-[25px] w-[200px] h-[200px] bg-white rounded-md shadow-lg z-50 border border-gray-200 hidden group-hover:flex flex-col items-center justify-center transition-all duration-200">
        {/* Triangle pointer connecting to button */}
        <div className="absolute -top-2 right-[76px] w-4 h-4 bg-white border-t border-l border-gray-200 transform rotate-45"></div>

        <a
          href="/user/profile"
          className="border-b border-gray-200 w-full h-[50px] flex justify-start items-center pl-5 px-4 py-2 text-[16px] text-gray-700 hover:bg-gray-100"
        >
          <FaUser className="mr-2" /> Profile
        </a>
        <a
          href="/user/orders"
          className="border-b border-gray-200 w-full h-[50px] justify-start  pl-5 px-4 py-2 text-[16px] text-gray-700 hover:bg-gray-100 flex items-center"
        >
          <IoCartOutline className="mr-2" /> My Orders
        </a>
        <a
          href="/user/wishlist"
          className="border-b border-gray-200 w-full h-[50px] justify-start  pl-5 px-4 py-2 text-[16px] text-gray-700 hover:bg-gray-100 flex items-center"
        >
          <FaRegHeart className="mr-2" /> Wishlist
        </a>
        <a
          href="/logout"
          className="border-b border-gray-200 w-full h-[50px] justify-start pl-5 px-4 py-2 text-[16px] text-error hover:bg-gray-100 flex items-center"
        >
          <IoIosLogOut className="mr-2" /> Logout
        </a>
      </div>
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
          <a
            key={index}
            href={category.href}
            className="text-[16px] pl-5 px-4 py-2  text-gray-700 hover:bg-gray-100 flex items-center"
          >
            {category.icon} {category.name}
          </a>
        ))}
        <div className="border-t border-gray-100 my-1"></div>
        <a
          href="/new-arrivals"
          className="pl-5 px-4 py-2 text-[16px] font-medium text-green-600 hover:bg-gray-100 flex items-center"
        >
          <CiHeart className="mr-2" /> New Arrivals
        </a>
        <a
          href="/sale"
          className="pl-5 px-4 py-2 text-[16px] font-medium text-red-600 hover:bg-gray-100 flex items-center"
        >
          <TbSettings className="mr-2" /> On Sale
        </a>
      </div>
    </div>
  );
};

const Header: React.FC = () => {
  const { cartItems } = useCart();
  const cartItemCount = cartItems.length;

  return (
    <>
      <header className="bg-white sticky top-0 z-50 flex justify-between items-center p-4 border-b-2 border-gray-200 h-[5rem] w-full">
        <div className="flex items-center space-x-4 ml-[80px] min-h-56 w-1/2">
          <Logo />
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
          <UserMenu />
          <a
            href="user/orders"
            className="flex items-center justify-center w-[4rem] h-[3rem] hover:text-green-600 transition-colors relative"
          >
            <FaShoppingCart className="mr-2 text-[1rem]" />
            Cart
            {cartItemCount > 0 && (
              <span className="absolute -top-1 -right-3 bg-red-500 text-white text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center">
                {cartItemCount}
              </span>
            )}
          </a>
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
    </>
  );
}
export default Header;