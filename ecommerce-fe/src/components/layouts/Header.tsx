import React from "react";
import { IoLocationSharp } from "react-icons/io5";
import { MdOutlineArrowDropDown } from "react-icons/md";
import { FaUser, FaShoppingCart } from "react-icons/fa";
import { TbMinusVertical } from "react-icons/tb";
import { FiSearch } from "react-icons/fi";
import { IoMenu } from "react-icons/io5";

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

const Header: React.FC = () => {
  return (
    <>
      <header className="flex justify-between items-center p-4 border-b-2 border-gray-200 h-[5rem] w-full">
        <div className="flex items-center space-x-4 ml-[80px] min-h-56 w-1/2">
          <Logo />
          <span className="flex items-center space-x-1 h-[3rem]">
            <TbMinusVertical className="h-[100%] text-[30px]" />
            <a href="/" className="flex items-center space-x-2 max-w-[9rem]">
              <IoLocationSharp className="text-[40px]" />
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
          <a
            href="/"
            className="text-[1rem] h-[100%] flex items-center hover:text-green-600 transition-colors"
          >
            <FaUser />
          </a>
          <a
            href="/"
            className="flex items-center justify-center w-[4rem] h-[3rem] hover:text-green-600 transition-colors"
          >
            <FaShoppingCart className="mr-2 text-[1rem]" />
            Cart
          </a>
        </div>
      </header>
      <nav className="flex justify-between items-center py-2 border-b-2 h-[3rem] max-w-[1440px]">
        {/* Left section */}
        <div className="flex items-center space-x-10 ml-[5.4rem]">
          <a href="/" className="flex items-center">
            <IoMenu className="text-xl mr-1" />
            Menu
            <span className="mx-4 text-gray-400">|</span>
          </a>

          <a href="/" className="flex items-center ml-[1rem] relative group">
            Explore
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/" className="flex items-center relative group">
            Deals
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/" className="flex items-center relative group">
            Saved
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
        </div>

        {/* Right section */}
        <div className="flex items-center space-x-10 mr-[5.4rem]">
          <a href="/" className="flex items-center relative group">
            Home
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/" className="flex items-center relative group">
            Product
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/" className="flex items-center relative group">
            About Us
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
          <a href="/" className="flex items-center relative group">
            Contact
            <span className="absolute h-[3px] w-0 bg-success bottom-[-5px] left-0 transition-all duration-300 group-hover:w-full"></span>
          </a>
        </div>
      </nav>
    </>
  );
}
export default Header;