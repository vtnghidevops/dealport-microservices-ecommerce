import React from "react";

const Footer: React.FC = () => {
  return (
    <footer className="bg-aqua-spring py-10">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-start max-w-6xl mx-auto py-10">
          {/* Company Info */}
          <div className="h-[182px] w-[161px] group">
            <h3 className="font-semibold text-lg mb-4 group-hover:text-cyprus transition-colors">
              Company Info
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  About Us
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Careers
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Press Releases
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Sustainability Practices
                </a>
              </li>
            </ul>
          </div>

          {/* Customer Support */}
          <div className="h-[182px] w-[161px] group">
            <h3 className="font-semibold text-lg mb-4 group-hover:text-cyprus transition-colors">
              Customer Support
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Contact Us
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Help Center (FAQs)
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Track My Order
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Return & Refund Policy
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Shipping Information
                </a>
              </li>
            </ul>
          </div>

          {/* Logo section in the middle */}
          <div className="mt-[5rem] ">
            <a href="/">
              <img
                src="/images/common/logo.png"
                alt="SapoGo"
                className="flex justify-center items-center h-[4rem] w-[240px]"
              />
            </a>
          </div>

          {/* Explore */}
          <div className="mb-6 group">
            <h3 className="font-semibold text-lg mb-4 group-hover:text-cyprus transition-colors">
              Explore
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Categories
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Bestsellers
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  New Arrivals
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Deals & Promotions
                </a>
              </li>
            </ul>
          </div>

          {/* Legal */}
          <div className="group">
            <h3 className="font-semibold text-lg mb-4 group-hover:text-cyprus transition-colors">
              Legal
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Terms & Conditions
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Privacy Policy
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Cookie Policy
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Accessibility Statement
                </a>
              </li>
              <li>
                <a
                  href="#"
                  className="text-gray-700 hover:text-cyprus transition-colors duration-300 inline-block transform  "
                >
                  Return & Refund Policy
                </a>
              </li>
            </ul>
          </div>
        </div>

        {/* Newsletter */}
        <div className="mt-[3rem] w-full flex flex-col justify-center items-center relative">
          <h3 className="mb-[0.5rem] body-text text-ocean-green">
            Newsletter Signup
          </h3>
          <div className="flex items-center justify-between bg-surf-crest rounded-full px-4 py-2 w-[25rem]">
            <input
              placeholder="Enter your email address"
              className="transparent-autofill text-black placeholder-black/70 w-[70%] bg-transparent focus:outline-none ml-[1rem]"
              type="text"
            ></input>
            <button className="flex items-center rounded-[20px] text-black w-[7rem] mr-[0.5rem] bg-white justify-center h-[2rem]  hover:shadow-md hover:scale-105 transition-transform transition-shadow duration-300">
              Subscribe
            </button>
          </div>
          {/* Contact us */}
          <div className="absolute md:mt-0 w-[150px] h-[5rem] right-[10.5%] mt-[0.5rem] p-1">
            <h3 className="font-medium mb-2">Connect with us</h3>
            <div className="flex space-x-4 w-[7rem]">
              <a href="#" className="text-blue-600 hover:text-blue-800">
                <svg
                  className="w-[2rem] h-[2rem]"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    fillRule="evenodd"
                    d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z"
                    clipRule="evenodd"
                  />
                </svg>
              </a>
              <a href="#" className="text-pink-600 hover:text-pink-800">
                <svg
                  className="w-[2rem] h-[2rem]"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    fillRule="evenodd"
                    d="M12.315 2c2.43 0 2.784.013 3.808.06 1.064.049 1.791.218 2.427.465a4.902 4.902 0 011.772 1.153 4.902 4.902 0 011.153 1.772c.247.636.416 1.363.465 2.427.048 1.067.06 1.407.06 4.123v.08c0 2.643-.012 2.987-.06 4.043-.049 1.064-.218 1.791-.465 2.427a4.902 4.902 0 01-1.153 1.772 4.902 4.902 0 01-1.772 1.153c-.636.247-1.363.416-2.427.465-1.067.048-1.407.06-4.123.06h-.08c-2.643 0-2.987-.012-4.043-.06-1.064-.049-1.791-.218-2.427-.465a4.902 4.902 0 01-1.772-1.153 4.902 4.902 0 01-1.153-1.772c-.247-.636-.416-1.363-.465-2.427-.047-1.024-.06-1.379-.06-3.808v-.63c0-2.43.013-2.784.06-3.808.049-1.064.218-1.791.465-2.427a4.902 4.902 0 011.153-1.772A4.902 4.902 0 015.45 2.525c.636-.247 1.363-.416 2.427-.465C8.901 2.013 9.256 2 11.685 2h.63zm-.081 1.802h-.468c-2.456 0-2.784.011-3.807.058-.975.045-1.504.207-1.857.344-.467.182-.8.398-1.15.748-.35.35-.566.683-.748 1.15-.137.353-.3.882-.344 1.857-.047 1.023-.058 1.351-.058 3.807v.468c0 2.456.011 2.784.058 3.807.045.975.207 1.504.344 1.857.182.466.399.8.748 1.15.35.35.683.566 1.15.748.353.137.882.3 1.857.344 1.054.048 1.37.058 4.041.058h.08c2.597 0 2.917-.01 3.96-.058.976-.045 1.505-.207 1.858-.344.466-.182.8-.398 1.15-.748.35-.35.566-.683.748-1.15.137-.353.3-.882.344-1.857.048-1.055.058-1.37.058-4.041v-.08c0-2.597-.01-2.917-.058-3.96-.045-.976-.207-1.505-.344-1.858a3.097 3.097 0 00-.748-1.15 3.098 3.098 0 00-1.15-.748c-.353-.137-.882-.3-1.857-.344-1.023-.047-1.351-.058-3.807-.058zM12 6.865a5.135 5.135 0 110 10.27 5.135 5.135 0 010-10.27zm0 1.802a3.333 3.333 0 100 6.666 3.333 3.333 0 000-6.666zm5.338-3.205a1.2 1.2 0 110 2.4 1.2 1.2 0 010-2.4z"
                    clipRule="evenodd"
                  />
                </svg>
              </a>
              <a href="#" className="text-gray-800 hover:text-gray-600">
                <svg
                  className="w-[2rem] h-[2rem]"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84" />
                </svg>
              </a>
              <a href="#" className="text-blue-800 hover:text-blue-600">
                <svg
                  className="w-[2rem] h-[2rem]"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 01-2.063-2.065 2.064 2.064 0 112.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
                </svg>
              </a>
            </div>
          </div>
        </div>

        {/* Social Media and Copyright */}
        <div className="relative mt-10 flex flex-col md:flex-row justify-between items-center border-t border-green-200 pt-6">
          <div className="ml-[8rem] text-sm text-[16px] text-gray-600">
            © 2025 SapoGo. All rights reserved
          </div>

          <div className="mr-[8rem] mt-4 md:mt-0 text-[16px] text-sm text-gray-600">
            Trusted Seller Certifications by SSL Secure
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
