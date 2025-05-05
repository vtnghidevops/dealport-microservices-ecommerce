import React, { useState, useEffect } from "react";
import { NavLink, useNavigate, useLocation } from "react-router-dom";
import { TiHome } from "react-icons/ti";
import { IoCartOutline } from "react-icons/io5";
import { LuUsersRound } from "react-icons/lu";
import { IoTicketOutline } from "react-icons/io5";
import { HiOutlineDuplicate } from "react-icons/hi";
import { SlCreditCard } from "react-icons/sl";
import { HiOutlineStar } from "react-icons/hi";
import { IoIosAddCircleOutline } from "react-icons/io";
import { IoImageOutline } from "react-icons/io5";
import { GrStorage } from "react-icons/gr";
import { MdOutlineReviews } from "react-icons/md";
import { RiUser6Line } from "react-icons/ri";
import { TbSettings } from "react-icons/tb";
import { BsArrowBarLeft, BsArrowBarRight } from "react-icons/bs";
import { IoIosLogOut } from "react-icons/io";
import { CiShop } from "react-icons/ci";
import { FiExternalLink } from "react-icons/fi";
interface SidebarProps {
  isOpen: boolean;
  onToggle?: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen = true, onToggle }) => {
  const [sidebarOpen, setSidebarOpen] = useState(isOpen);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    setSidebarOpen(isOpen);
  }, [isOpen]);

  const menuItems = [
    {
      name: "Main Menu",
      children: [
        { name: "Dashboard", icon: <TiHome />, link: "/admin/dashboard" },
        {
          name: "Order Management",
          icon: <IoCartOutline />,
          link: "/admin/orders",
        },
        { name: "Customers", icon: <LuUsersRound />, link: "/admin/customers" },
        {
          name: "Coupon Code",
          icon: <IoTicketOutline />,
          link: "/admin/coupons",
        },
        {
          name: "Categories",
          icon: <HiOutlineDuplicate />,
          link: "/admin/categories",
        },
        // {
        //   name: "Transaction",
        //   icon: <SlCreditCard />,
        //   link: "/admin/transactions",
        // },
        { name: "Brand", icon: <HiOutlineStar />, link: "/admin/brands" },
      ],
    },
    {
      name: "Product",
      children: [
        {
          name: "Add Products",
          icon: <IoIosAddCircleOutline />,
          link: "/admin/add-product",
        },
        // {
        //   name: "Product Media",
        //   icon: <IoImageOutline />,
        //   link: "/admin/products/media",
        // },
        {
          name: "Product List",
          icon: <GrStorage />,
          link: "/admin/products/list",
        },
        {
          name: "Product Reviews",
          icon: <MdOutlineReviews />,
          link: "/admin/products/reviews",
        },
      ],
    },
    {
      name: "Admin",
      children: [
        {
          name: "Admin role",
          icon: <RiUser6Line />,
          link: "/admin/role",
        },
        {
          name: "Control Authority",
          icon: <TbSettings />,
          link: "/admin/authority",
        },
      ],
    },
  ];

  // Kiểm tra xem đường dẫn có khớp với link menu không
  const isLinkActive = (link: string): boolean => {
    // Kiểm tra đúng URL
    if (location.pathname === link) return true;

    // Kiểm tra URL con (subpath)
    // Ví dụ: /admin/orders/123 cũng sẽ active cho /admin/orders
    if (link !== "/admin/dashboard" && location.pathname.startsWith(link)) return true;

    return false;
  };

  const toggleSidebar = () => {
    const newState = !sidebarOpen;
    setSidebarOpen(newState);

    if (onToggle) {
      onToggle();
    }
  };

  const handleLogoClick = () => {
    navigate('/');
  };

  return (
    <aside
      className={`bg-white w-[260px] border-r border-t border-gray-300 shadow-lg transition-all duration-300 ease-in-out flex flex-col ${sidebarOpen ? "translate-x-0" : "-translate-x-[200px] w-[60px]"
        } md:translate-x-0 fixed md:relative z-10`}
    >
      <div
        className={`flex mt-[1rem] items-center justify-between h-16 w-full px-3 py-5 flex-shrink-0
        ${sidebarOpen ? "relative" : ""}`}
      >
        <img
          src="/images/common/logo.png"
          alt="Logo"
          className={`text-2xl font-bold text-gray-800 max-w-[8rem] cursor-pointer transition-opacity duration-300 ${!sidebarOpen && "opacity-0"
            }`}
          onClick={handleLogoClick}
        />
        <span
          className={`cursor-pointer text-2xl ${!sidebarOpen ? "absolute left-[30%]" : ""
            }`}
          onClick={toggleSidebar}
        >
          {sidebarOpen ? <BsArrowBarLeft /> : <BsArrowBarRight />}
        </span>
      </div>

      {/* Main navigation - using flex-grow to fill available space */}
      <nav className="py-4 w-full flex-grow overflow-y-auto ">
        <ul className="mt-[10px]">
          {menuItems.map((item, index) => (
            <li
              key={index}
              className="flex justify-center flex-col px-[1rem] mt-5"
            >
              <div
                className={`flex items-center text-gray-600 hover:text-blue-500 ${!sidebarOpen && "hidden"
                  }`}
              >
                <span className="text-[15px] text-neutral-500">
                  {item.name}
                </span>
              </div>
              {item.children && (
                <ul className="flex flex-col items-center justify-center ">
                  {item.children.map((child, childIndex) => {
                    const isActive = isLinkActive(child.link);
                    return (
                      <li
                        key={childIndex}
                        className={`py-1 my-3 mb-0 w-full justify-center items-center flex ${sidebarOpen
                          ? "h-[2.5rem] rounded-[8px]"
                          : "h-[1.8rem] rounded-[10px]"
                          } ${isActive
                            ? "bg-ocean-green text-white"
                            : ""
                          }`}
                        title={!sidebarOpen ? child.name : ""}
                      >
                        <NavLink
                          to={child.link}
                          className={`w-full flex items-center px-[0.5rem] ${!sidebarOpen ? "justify-center" : "justify-start"
                            } ${isActive
                              ? "text-white"
                              : "text-neutral-500 hover:text-blue-500"
                            }`}
                        >
                          <span className="text-[18px] body-text flex items-center justify-center">
                            {child.icon}
                          </span>
                          {sidebarOpen && (
                            <span className="ml-2 body-text flex items-center justify-center">
                              {child.name}
                            </span>
                          )}
                        </NavLink>
                      </li>
                    );
                  })}
                </ul>
              )}
            </li>
          ))}
        </ul>
        {/* User profile and action buttons - at the bottom with flex-shrink-0 */}
        <div className="border-t border-gray-300 py-3 px-4 flex-shrink-0 mt-[3rem]">
          <div className="flex items-center mb-4 px-3">
            <img
              src="/images/common/avatars/admin.png"
              alt="Admin"
              className={`object-cover border border-gray-300 rounded-full ${sidebarOpen
                ? "h-[40px] w-[40px] mr-3"
                : "h-[1.8rem] w-[1.8rem] mx-auto"
                }`}
            />
            {sidebarOpen && (
              <div className="w-[80%] flex justify-between items-center">
                <div className="flex flex-col justify-center w-[80%]">
                  <p className="font-medium text-[16px] text-gray-800">
                    Dealport
                  </p>
                  <p className="text-xs text-[16px] text-[#737373] truncate">
                    tannghi.devops@example.com
                  </p>
                </div>
                <a
                  href="/admin/logout"
                  className="ml-[1rem] block w-[20%] text-[#737373] hover:text-red-500 text-[20px] transition-all duration-100"
                >
                  <IoIosLogOut></IoIosLogOut>
                </a>
              </div>
            )}
          </div>

          {/* Your Shop Section */}
          <div
            className={`flex items-center ${sidebarOpen ? "px-[20px]" : ""
              } py-[12px] w-full h-[48px] rounded-md border shadow-[0_-6px_10px_-4px_rgba(209,213,219,0.3)] drop-shadow-lg mt-[1rem]`}
          >
            <div
              className={`text-[#0F3641] text-xl ${!sidebarOpen && "mx-auto"}`}
            >
              <CiShop />
            </div>
            {sidebarOpen && (
              <div className="flex justify-between items-center w-full">
                <a
                  href="/shop"
                  className="font-medium text-[16px] text-[#0F3641] ml-3"
                >
                  Your Shop
                </a>
                <a
                  href="/shop/edit"
                  className="text-[#737373] hover:text-gray-700 text-[20px]"
                >
                  <FiExternalLink size={16} />
                </a>
              </div>
            )}
          </div>
        </div>
      </nav>
    </aside>
  );
};

export default Sidebar;
