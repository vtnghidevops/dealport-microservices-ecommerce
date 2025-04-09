import React, { useState } from "react";
// import { MdOutlineLocalShipping } from "react-icons/md";
// import { SlBadge } from "react-icons/sl";
// import { BsHeadset } from "react-icons/bs";
// import { CiCreditCard2 } from "react-icons/ci";
// import { LiaHandshakeSolid } from "react-icons/lia";
interface ProductInformationProps {
  product: {
    description: string;
    features: {
      id: string;
      value: string;
    }[];
    shippingInfo: {
      courier: string;
      local: string;
      ups: string;
      global: string;
    };
  };
}

const ProductInformation: React.FC<ProductInformationProps> = ({ product }) => {
  const [activeTab, setActiveTab] = useState<string>("description");

  return (
    <div className="bg-white p-6 rounded-lg shadow-sm mt-[1rem] mb-[1rem]">
      {/* Tabs */}
      <div className="flex justify-center gap-9 border-b mb-4 overflow-x-auto">
        <button
          className={`px-4 py-2 text-sm font-medium whitespace-nowrap ${
            activeTab === "description"
              ? "text-gray-900 font-sans border-b-2 border-orange-500"
              : "text-gray-500 font-sans hover:text-orange-500"
          }`}
          onClick={() => setActiveTab("description")}
        >
          DESCRIPTION
        </button>
        <button
          className={`px-4 py-2 text-sm font-medium whitespace-nowrap ${
            activeTab === "additional"
              ? "text-gray-900 font-sans border-b-2 border-orange-500"
              : "text-gray-500 font-sans hover:text-orange-500"
          }`}
          onClick={() => setActiveTab("additional")}
        >
          ADDITIONAL INFORMATION
        </button>
        <button
          className={`px-4 py-2 text-sm font-medium whitespace-nowrap ${
            activeTab === "specification"
              ? "text-orange-500 font-sans border-b-2 border-orange-500"
              : "text-gray-500 font-sans hover:text-orange-500"
          }`}
          onClick={() => setActiveTab("specification")}
        >
          SPECIFICATION
        </button>
        <button
          className={`px-4 py-2 text-sm font-medium whitespace-nowrap ${
            activeTab === "review"
              ? "text-gray-900 font-sans border-b-2 border-orange-500"
              : "text-gray-500 font-sans hover:text-orange-500"
          }`}
          onClick={() => setActiveTab("review")}
        >
          REVIEW
        </button>
      </div>

      {/* Content based on active tab */}
      {activeTab === "description" && (
        <div className="flex gap-16 p-[1.5rem] ">
          <div className="col-span-2 max-w-[37rem]">
            <h3 className="text-lg font-semibold mb-3">Description</h3>
            <p className="text-gray-600">{product.description}</p>
            <p className="text-gray-600 mt-4">
              Even the most ambitious projects are easily handled with up to 10
              CPU cores, up to 16 GPU cores, a 16-core Neural Engine, and
              dedicated encode and decode media engines that support H.264,
              HEVC, and ProRes codecs.
            </p>
          </div>
          <div className="flex gap-[3rem]">
            <div>
              <h3 className="text-lg font-semibold mb-3">Feature</h3>
              <ul className="space-y-2">
                {product.features.map((feature, index) => (
                  <li key={index} className="flex items-center text-gray-600">
                    <span className="flex-shrink-0 w-6 h-6 flex items-center justify-center rounded-full bg-orange-100 text-orange-500 mr-2">
                      ✓
                    </span>
                    {feature.value}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="text-lg font-semibold mt-6 mb-3">
                Shipping Information
              </h3>
              <ul className="space-y-2 text-gray-600">
                <li className="flex items-center gap-1">
                  <span className="text-gray-900 font-medium">Courier:</span>{" "}
                  {product.shippingInfo.courier}
                </li>
                <li className="flex items-center gap-1">
                  <span className="text-gray-900 font-medium">
                    Local Shipping:
                  </span>
                  {product.shippingInfo.local}
                </li>
                <li className="flex items-center gap-1">
                  <span className="text-gray-900 font-medium">
                    UPS Ground Shipping:
                  </span>
                  {product.shippingInfo.ups}
                </li>
                <li className="flex items-center gap-1">
                  <span className="text-gray-900 font-medium">
                    Unishop Global Export:
                  </span>
                  {product.shippingInfo.global}
                </li>
              </ul>
            </div>
          </div>
        </div>
      )}

      {activeTab === "additional" && (
        <div className="p-[1.5rem]">
          <h3 className="text-lg font-semibold mb-3">Additional Information</h3>
          <table className="w-full border-collapse">
            <tbody>
              <tr className="border-b">
                <td className="py-2 pr-4 font-medium text-cyprus w-1/4">
                  Weight
                </td>
                <td className="py-2">1.2 kg</td>
              </tr>
              <tr className="border-b">
                <td className="py-2 pr-4 font-medium text-cyprus">
                  Dimensions
                </td>
                <td className="py-2">30.41 × 21.24 × 1.61 cm</td>
              </tr>
              <tr className="border-b">
                <td className="py-2 pr-4 font-medium text-cyprus">Color</td>
                <td className="py-2">Space Gray, Silver</td>
              </tr>
              <tr className="border-b">
                <td className="py-2 pr-4 font-medium text-cyprus">Storage</td>
                <td className="py-2">256GB, 512GB, 1TB, 2TB</td>
              </tr>
              <tr className="border-b">
                <td className="py-2 pr-4 font-medium text-cyprus">Memory</td>
                <td className="py-2">8GB, 16GB, 32GB</td>
              </tr>
            </tbody>
          </table>
        </div>
      )}

      {activeTab === "specification" && (
        <div className="p-[1.5rem]">
          <h3 className="text-lg font-semibold mb-3">
            Technical Specifications
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h4 className="font-bold mb-2 text-cyprus">Processor</h4>
              <ul className="list-disc list-inside text-gray-600 space-y-1 mb-4">
                <li>Apple M1 Pro or M1 Max chip</li>
                <li>Up to 10-core CPU</li>
                <li>Up to 32-core GPU</li>
                <li>16-core Neural Engine</li>
              </ul>

              <h4 className="font-bold mb-2 text-cyprus">Storage</h4>
              <ul className="list-disc list-inside text-gray-600 space-y-1 mb-4">
                <li>256GB, 512GB, 1TB, 2TB, 4TB, or 8TB SSD</li>
              </ul>

              <h4 className="font-bold mb-2 text-cyprus">Memory</h4>
              <ul className="list-disc list-inside text-gray-600 space-y-1">
                <li>16GB, 32GB, or 64GB unified memory</li>
              </ul>
            </div>

            <div>
              <h4 className="font-bold mb-2 text-cyprus">Display</h4>
              <ul className="list-disc list-inside text-gray-600 space-y-1 mb-4">
                <li>Liquid Retina XDR display</li>
                <li>14.2-inch or 16.2-inch diagonal</li>
                <li>3024 × 1964 or 3456 × 2234 native resolution</li>
                <li>1000 nits sustained brightness, 1600 nits peak</li>
              </ul>

              <h4 className="font-bold mb-2 text-cyprus">Battery</h4>
              <ul className="list-disc list-inside text-gray-600 space-y-1 mb-4">
                <li>Up to 17 hours video playback</li>
                <li>Up to 11 hours wireless web</li>
                <li>Fast charging capability</li>
              </ul>
            </div>
          </div>
        </div>
      )}

      {activeTab === "review" && (
        <div className="p-[1.5rem]">
          <div className="flex items-center mb-6">
            <div className="mr-4 w-1/4 flex flex-col justify-center items-center">
              <div className="text-5xl font-bold text-cyprus">4.7</div>
              <div className="flex text-[#FF9017] mt-1">
                <span>★</span>
                <span>★</span>
                <span>★</span>
                <span>★</span>
                <span>★</span>
              </div>
              <div className="text-sm text-gray-500 mt-1">21,671 Ratings</div>
              <button className="w-[160px] mt-5 bg-[#0496FF] text-white px-6 py-2 rounded-lg hover:bg-blue-500">
                Write a Review
              </button>
            </div>

            <div className="flex-1 w-3/4">
              <div className="flex items-center mb-1">
                <span className="w-[50px] text-sm text-gray-600">5 Star</span>
                <div className="flex-1 h-2 mx-2 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="bg-[#FF9017] h-full rounded-full"
                    style={{ width: "75%" }}
                  ></div>
                </div>
                <span className="w-8 text-sm text-gray-600">75%</span>
              </div>
              <div className="flex items-center mb-1">
                <span className="w-[50px] text-sm text-gray-600">4 Star</span>
                <div className="flex-1 h-2 mx-2 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="bg-[#FF9017] h-full rounded-full"
                    style={{ width: "20%" }}
                  ></div>
                </div>
                <span className="w-8 text-sm text-gray-600">20%</span>
              </div>
              <div className="flex items-center mb-1">
                <span className="w-[50px] text-sm text-gray-600">3 Star</span>
                <div className="flex-1 h-2 mx-2 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="bg-[#FF9017] h-full rounded-full"
                    style={{ width: "3%" }}
                  ></div>
                </div>
                <span className="w-8 text-sm text-gray-600">3%</span>
              </div>
              <div className="flex items-center mb-1">
                <span className="w-[50px] text-sm text-gray-600">2 Star</span>
                <div className="flex-1 h-2 mx-2 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="bg-[#FF9017] h-full rounded-full"
                    style={{ width: "1%" }}
                  ></div>
                </div>
                <span className="w-8 text-sm text-gray-600">1%</span>
              </div>
              <div className="flex items-center">
                <span className="w-[50px] text-sm text-gray-600">1 Star</span>
                <div className="flex-1 h-2 mx-2 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="bg-[#FF9017] h-full rounded-full"
                    style={{ width: "1%" }}
                  ></div>
                </div>
                <span className="w-8 text-sm text-gray-600">1%</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProductInformation;