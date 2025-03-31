// // components/admin/product/add/AddProduct.tsx
// import React, { useState, ChangeEvent } from "react";
// import { Product } from "../models/product.model";
// import { CiCirclePlus } from "react-icons/ci";
// export const AddProduct: React.FC = () => {
//   const [product, setProduct] = useState<Product>({
//     id: "",
//     name: "",
//     description: "",
//     price: 0,
//     discountedPrice: 0,
//     saleAmount: 0,
//     taxIncluded: true,
//     expirationStart: "",
//     expirationEnd: "",
//     stockQuantity: "Unlimited",
//     stockStatus: "In Stock",
//     highlighted: false,
//     images: [],
//     categories: [],
//     tags: [],
//     color: "",
//   });

//   const [selectedImages, setSelectedImages] = useState<string[]>([]);
//   const [mainImage, setMainImage] = useState<string | null>(null);

//   const handleInputChange = (
//     e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
//   ) => {
//     const { name, value, type } = e.target as HTMLInputElement;

//     if (type === "checkbox") {
//       const checked = (e.target as HTMLInputElement).checked;
//       setProduct((prev) => ({ ...prev, [name]: checked }));
//     } else {
//       setProduct((prev) => ({ ...prev, [name]: value }));
//     }
//   };

//   const handleRadioChange = (e: ChangeEvent<HTMLInputElement>) => {
//     setProduct((prev) => ({
//       ...prev,
//       [e.target.name]: e.target.value === "Yes",
//     }));
//   };

//   const handleImageUpload = (e: ChangeEvent<HTMLInputElement>) => {
//     const files = e.target.files;
//     if (!files) return;

//     const newImages = Array.from(files).map((file) =>
//       URL.createObjectURL(file)
//     );
//     setSelectedImages((prev) => [...prev, ...newImages]);

//     if (!mainImage && newImages.length > 0) {
//       setMainImage(newImages[0]);
//     }
//   };

//   const handleColorSelect = (color: string) => {
//     setProduct((prev) => ({ ...prev, color }));
//   };

//   const handlePublish = () => {
//     console.log("Publishing product:", product);
//     // Integrate with your backend service here
//     alert("Product published successfully!");
//   };

//   const handleSaveDraft = () => {
//     console.log("Saving draft:", product);
//     // Save draft logic here
//     alert("Draft saved successfully!");
//   };

//   const calculateSavings = () => {
//     const originalPrice = parseFloat(String(product.price)) || 0;
//     const discountPrice = parseFloat(String(product.discountedPrice)) || 0;

//     // Kiểm tra nếu discountPrice có giá trị và nhỏ hơn originalPrice
//     if (discountPrice && discountPrice < originalPrice) {
//       const savings = originalPrice - discountPrice;
//       return savings.toFixed(2); // Làm tròn đến 2 chữ số thập phân
//     }

//     return "0.00"; // Giá trị mặc định nếu không có giảm giá
//   };

//   return (
//     <div className="container mx-auto">
//       <div className="flex items-center justify-between mb-6">
//         <h1 className="text-[18px] font-bold text-cyprus">Add New Product</h1>
//         <div className="flex items-center space-x-2">
//           <div className="mr-[0.5rem] bg-white border rounded-lg border-gray-200  relative w-[315p] h-[48px] px-12 py-6">
//             <input
//               type="text"
//               placeholder="Search product for add"
//               className="pl-10 pr-4 py-2 text-neutral-600 w-72 placeholder:text-neutral-600 focus:outline-none"
//             />
//             <svg
//               className="w-5 h-5 absolute top-[54%] left-3 -translate-y-1/2 text-neutral-600"
//               viewBox="0 0 20 20"
//               fill="currentColor"
//             >
//               <path
//                 fillRule="evenodd"
//                 d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
//                 clipRule="evenodd"
//               />
//             </svg>
//           </div>
//           <button
//             onClick={handlePublish}
//             className="py-6 px-12 bg-ocean-green border rounded-lg border-gray-200  text-white relative w-[315p] h-[48px] hover:bg-green-600 flex items-center"
//           >
//             Publish Product
//           </button>
//           <button
//             onClick={handleSaveDraft}
//             className="py-6 px-12 bg-white border-neutral-300 rounded-lg border text-cyprus relative w-[315p] h-[48px] flex items-center hover:bg-gray-50"
//           >
//             <svg
//               className="w-5 h-5 mr-1"
//               fill="none"
//               stroke="currentColor"
//               viewBox="0 0 24 24"
//               xmlns="http://www.w3.org/2000/svg"
//             >
//               <path
//                 strokeLinecap="round"
//                 strokeLinejoin="round"
//                 strokeWidth={2}
//                 d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
//               />
//             </svg>
//             Save to draft
//           </button>
//           <button className="flex items-center justify-center text-neutral-500 border border-neutral-300 bg-white rounded-lg w-[48px] h-[48px] hover:bg-gray-100 p-2">
//             <CiCirclePlus className="p-4 w-full h-full"></CiCirclePlus>
//           </button>
//         </div>
//       </div>

//       <div className="flex flex-wrap -mx-3 mt-[1.5rem] drop-shadow-sm filter ">
//         <div className="w-[611px] h-[1053px] md:w-7/12 px-3 mb-6">
//           <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6">
//             <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
//               Basic Details
//             </h2>

//             <div className="mb-4">
//               <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//                 Product Name
//               </label>
//               <input
//                 type="text"
//                 name="name"
//                 value={product.name}
//                 onChange={handleInputChange}
//                 placeholder="iPhone 15"
//                 className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
//               />
//             </div>

//             <div className="mb-4">
//               <label className="block text-cyprus font-bold text-[15px] mb-12">
//                 Product Description
//               </label>
//               <textarea
//                 name="description"
//                 value={product.description}
//                 onChange={handleInputChange}
//                 placeholder="The iPhone 15 delivers cutting-edge performance with the A16 Bionic chip, an immersive Super Retina XDR display, advanced dual-camera system, and exceptional battery life, all encased in stunning aerospace-grade aluminum."
//                 rows={5}
//                 className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
//               ></textarea>
//               <div className="flex justify-end mt-2 space-x-2 border-b border-neutral-300">
//                 <button className="p-2 hover:bg-gray-100 rounded">
//                   <svg
//                     className="w-5 h-5 text-gray-500"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
//                     />
//                   </svg>
//                 </button>
//                 <button className="p-2 hover:bg-gray-100 rounded">
//                   <svg
//                     className="w-5 h-5 text-gray-500"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
//                     />
//                   </svg>
//                 </button>
//               </div>
//             </div>
//             <div className="bg-white rounded-lg p-6 mb-6">
//               <h2 className="block text-cyprus font-bold text-[22px] mb-12 mt-[1rem]">
//                 Pricing
//               </h2>

//               <div className="mb-4">
//                 <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//                   Product Price
//                 </label>
//                 <div className="flex">
//                   <div className="flex-grow">
//                     <input
//                       type="text"
//                       name="price"
//                       value={product.price || ""}
//                       onChange={handleInputChange}
//                       placeholder="$999.89"
//                       className="w-full border text-cyprus bg-neutral-50 border-gray-200 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
//                     />
//                   </div>
//                 </div>
//               </div>

//               <div className="grid grid-cols-2 gap-4 mb-4">
//                 <div>
//                   <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//                     Discounted Price{" "}
//                     <span className="text-neutral-500">(Optional)</span>
//                   </label>
//                   <div className="border border-gray-200 bg-neutral-50 rounded-lg flex items-center w-[272px] gap-4">
//                     <span className="ml-2 flex items-center justify-center bg-aqua-spring w-[32px] h-[32px] text-black rounded-lg mr-2">
//                       $
//                     </span>
//                     <input
//                       type="text"
//                       name="discountedPrice"
//                       value={product.discountedPrice || ""}
//                       onChange={handleInputChange}
//                       placeholder="$99"
//                       className="w-[108px] flex-grow py-[10px] text-cyprus bg-neutral-50 px-12 focus:outline-none focus:border-ocean-green"
//                     />
//                     <label className="block text-[15px] font-bold text-cyprus mb-1 w-[103px] h-[18px]">
//                       Sale ={" "}
//                       <span className="text-cyprus">${calculateSavings()}</span>
//                     </label>
//                   </div>
//                 </div>
//               </div>

//               <div className="mb-4">
//                 <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//                   Expiration
//                 </label>
//                 <div className="grid grid-cols-2 gap-4">
//                   <div className="relative">
//                     <input
//                       type="text" // Thay đổi từ "date" thành "text"
//                       onFocus={(e) => (e.target.type = "date")}
//                       onBlur={(e) => {
//                         if (!e.target.value) e.target.type = "text";
//                       }}
//                       name="expirationStart"
//                       value={product.expirationStart}
//                       onChange={handleInputChange}
//                       placeholder="Start"
//                       className="w-full border border-gray-200 rounded-lg p-2 pr-10 text-cyprus bg-neutral-50 focus:outline-none focus:border-ocean-green"
//                     />
//                     <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500">
//                       <svg
//                         className="w-5 h-5"
//                         fill="none"
//                         stroke="currentColor"
//                         viewBox="0 0 24 24"
//                       >
//                         <path
//                           strokeLinecap="round"
//                           strokeLinejoin="round"
//                           strokeWidth={2}
//                           d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
//                         />
//                       </svg>
//                     </div>
//                   </div>
//                   <div className="relative">
//                     <input
//                       type="text" // Thay đổi từ "date" thành "text"
//                       onFocus={(e) => (e.target.type = "date")}
//                       onBlur={(e) => {
//                         if (!e.target.value) e.target.type = "text";
//                       }}
//                       name="expirationEnd"
//                       value={product.expirationEnd}
//                       onChange={handleInputChange}
//                       placeholder="End"
//                       className="w-full border border-gray-200 rounded-lg p-2 pr-10 bg-neutral-50 focus:outline-none focus:border-ocean-green"
//                     />
//                     <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500">
//                       <svg
//                         className="w-5 h-5"
//                         fill="none"
//                         stroke="currentColor"
//                         viewBox="0 0 24 24"
//                       >
//                         <path
//                           strokeLinecap="round"
//                           strokeLinejoin="round"
//                           strokeWidth={2}
//                           d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
//                         />
//                       </svg>
//                     </div>
//                   </div>
//                 </div>
//               </div>
//             </div>
//             <div className="bg-white rounded-lg p-6 mb-6">
//               <h2 className="block text-cyprus font-bold text-[22px] mb-12 mt-[1rem]">
//                 Inventory
//               </h2>
//               <div className="flex items-center justify-between">
//                 <div className="mb-4 w-[272px] ">
//                   <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//                     Stock Quantity
//                   </label>
//                   <input
//                     type="text"
//                     name="stockQuantity"
//                     value={product.stockQuantity}
//                     onChange={handleInputChange}
//                     placeholder="Unlimited"
//                     className="mb-4 w-full border border-gray-200 rounded-lg p-2 bg-neutral-50 focus:outline-none focus:border-ocean-green"
//                   />
//                   <div className="flex items-center mt-2">
//                     <label className="flex items-center cursor-pointer">
//                       <div
//                         className={`mr-2 relative inline-block w-[48px] h-[24px] transition-all duration-200 ease-in-out rounded-full ${
//                           product.stockQuantity === "Unlimited"
//                             ? "bg-green-500"
//                             : "bg-gray-300"
//                         }`}
//                       >
//                         <input
//                           type="checkbox"
//                           className="opacity-0 absolute  w-full h-full"
//                           checked={product.stockQuantity === "Unlimited"}
//                           onChange={(e) =>
//                             setProduct((prev) => ({
//                               ...prev,
//                               stockQuantity: e.target.checked
//                                 ? "Unlimited"
//                                 : "0",
//                             }))
//                           }
//                         />
//                         <span
//                           className={`absolute left-4 top-1/2 -translate-y-1/2 bg-white w-[15px] h-[15px] rounded-full transition-transform duration-200 ease-in-out transform ${
//                             product.stockQuantity === "Unlimited"
//                               ? "translate-x-5"
//                               : "translate-x-0"
//                           }`}
//                         ></span>
//                       </div>
//                       <span className="text-sm text-gray-700">Unlimited</span>
//                     </label>
//                   </div>
//                 </div>

//                 <div className="w-[272px] mb-[2.1rem]">
//                   <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
//                     Stock Status
//                   </label>
//                   <div className="relative">
//                     <select
//                       name="stockStatus"
//                       value={product.stockStatus}
//                       onChange={(e) =>
//                         setProduct((prev) => ({
//                           ...prev,
//                           stockStatus: e.target.value,
//                         }))
//                       }
//                       className="text-cyprus w-full border border-gray-200 bg-neutral-50 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:border-ocean-green"
//                     >
//                       <option value="In Stock">In Stock</option>
//                       <option value="Out of Stock">Out of Stock</option>
//                       <option value="Pre-order">Pre-order</option>
//                     </select>
//                     <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
//                       <svg
//                         className="w-5 h-5"
//                         fill="none"
//                         stroke="currentColor"
//                         viewBox="0 0 24 24"
//                       >
//                         <path
//                           strokeLinecap="round"
//                           strokeLinejoin="round"
//                           strokeWidth={2}
//                           d="M19 9l-7 7-7-7"
//                         />
//                       </svg>
//                     </div>
//                   </div>
//                 </div>
//               </div>

//               <div className="mb-4 mt-[0.5rem]">
//                 <label className="flex items-center cursor-pointer">
//                   <input
//                     type="checkbox"
//                     name="highlighted"
//                     checked={product.highlighted}
//                     onChange={(e) =>
//                       setProduct((prev) => ({
//                         ...prev,
//                         highlighted: e.target.checked,
//                       }))
//                     }
//                     className="h-[20px] w-[20px] text-ocean-green focus:ring-green-500 border-gray-300 rounded"
//                   />
//                   <span className="ml-2 text-[15px] text-neutral-500 ">
//                     Highlight this product in a featured section.
//                   </span>
//                 </label>
//               </div>
//             </div>
//             <div className="flex justify-end mt-6 space-x-3">
//               <button
//                 onClick={handleSaveDraft}
//                 className="px-5 py-2 border border-gray-200 rounded-lg flex items-center hover:bg-gray-50"
//               >
//                 <svg
//                   className="w-5 h-5 mr-2"
//                   fill="none"
//                   stroke="currentColor"
//                   viewBox="0 0 24 24"
//                 >
//                   <path
//                     strokeLinecap="round"
//                     strokeLinejoin="round"
//                     strokeWidth={2}
//                     d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
//                   />
//                 </svg>
//                 Save to draft
//               </button>
//               <button
//                 onClick={handlePublish}
//                 className="px-5 py-2 bg-ocean-green hover:bg-green-600 text-white rounded-lg"
//               >
//                 Publish Product
//               </button>
//             </div>
//           </div>
//         </div>

//         <div className="w-full md:w-5/12 px-3">
//           <div className="bg-white rounded-lg p-6 shadow-sm mb-6">
//             <h2 className="font-medium text-lg mb-4">Upload Product Image</h2>
//             <p className="text-sm text-gray-500 mb-3">Product Image</p>

//             <div className="mb-4 border border-gray-200 rounded-lg p-4 flex flex-col items-center">
//               {mainImage ? (
//                 <div className="relative w-full h-40 mb-3">
//                   <img
//                     src={mainImage}
//                     className="w-full h-full object-contain"
//                     alt="Product"
//                   />
//                 </div>
//               ) : (
//                 <div className="w-full h-40 mb-3 flex items-center justify-center bg-gray-50">
//                   <svg
//                     className="w-16 h-16 text-gray-300"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
//                     />
//                   </svg>
//                 </div>
//               )}

//               <div className="flex justify-between w-full mb-4">
//                 <button
//                   className="px-3 py-1 border border-gray-200 rounded-lg bg-gray-50 text-sm flex items-center"
//                   onClick={() =>
//                     document.getElementById("file-upload")?.click()
//                   }
//                 >
//                   <svg
//                     className="w-4 h-4 mr-1"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
//                     />
//                   </svg>
//                   Browse
//                   <input
//                     id="file-upload"
//                     type="file"
//                     className="hidden"
//                     onChange={handleImageUpload}
//                     multiple
//                     accept="image/*"
//                   />
//                 </button>

//                 {mainImage && (
//                   <button
//                     className="px-3 py-1 border border-gray-200 rounded-lg bg-gray-50 text-sm flex items-center"
//                     onClick={() => {
//                       setMainImage(null);
//                       setSelectedImages([]);
//                     }}
//                   >
//                     <svg
//                       className="w-4 h-4 mr-1"
//                       fill="none"
//                       stroke="currentColor"
//                       viewBox="0 0 24 24"
//                     >
//                       <path
//                         strokeLinecap="round"
//                         strokeLinejoin="round"
//                         strokeWidth={2}
//                         d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
//                       />
//                     </svg>
//                     Replace
//                   </button>
//                 )}
//               </div>

//               <div className="flex w-full space-x-2 overflow-x-auto pb-2">
//                 {selectedImages.length > 0 &&
//                   selectedImages.map((img, idx) => (
//                     <div
//                       key={idx}
//                       className="flex-shrink-0 w-16 h-16 relative rounded border border-gray-200"
//                     >
//                       <img
//                         src={img}
//                         className="w-full h-full object-cover rounded"
//                         alt={`Product ${idx + 1}`}
//                       />
//                       <button
//                         className="absolute top-0 right-0 bg-white rounded-full p-0.5 transform translate-x-1/3 -translate-y-1/3 shadow"
//                         onClick={() => {
//                           const newImages = [...selectedImages];
//                           newImages.splice(idx, 1);
//                           setSelectedImages(newImages);
//                           if (mainImage === img) {
//                             setMainImage(
//                               newImages.length > 0 ? newImages[0] : null
//                             );
//                           }
//                         }}
//                       >
//                         <svg
//                           className="w-3 h-3 text-gray-500"
//                           fill="none"
//                           stroke="currentColor"
//                           viewBox="0 0 24 24"
//                         >
//                           <path
//                             strokeLinecap="round"
//                             strokeLinejoin="round"
//                             strokeWidth={2}
//                             d="M6 18L18 6M6 6l12 12"
//                           />
//                         </svg>
//                       </button>
//                     </div>
//                   ))}

//                 <div
//                   className="flex-shrink-0 w-16 h-16 border border-dashed border-gray-300 rounded flex items-center justify-center cursor-pointer hover:bg-gray-50"
//                   onClick={() =>
//                     document.getElementById("file-upload")?.click()
//                   }
//                 >
//                   <svg
//                     className="w-6 h-6 text-gray-400"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M12 6v6m0 0v6m0-6h6m-6 0H6"
//                     />
//                   </svg>
//                   <span className="sr-only">Add Image</span>
//                 </div>
//               </div>
//             </div>
//           </div>

//           <div className="bg-white rounded-lg p-6 shadow-sm">
//             <h2 className="font-medium text-lg mb-4">Categories</h2>

//             <div className="mb-4">
//               <label className="block text-sm font-medium mb-1">
//                 Product Categories
//               </label>
//               <div className="relative">
//                 <select className="w-full border border-gray-200 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500">
//                   <option>Select your product</option>
//                 </select>
//                 <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
//                   <svg
//                     className="w-5 h-5"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M19 9l-7 7-7-7"
//                     />
//                   </svg>
//                 </div>
//               </div>
//             </div>

//             <div className="mb-4">
//               <label className="block text-sm font-medium mb-1">
//                 Product Tag
//               </label>
//               <div className="relative">
//                 <select className="w-full border border-gray-200 rounded-lg p-2 pr-10 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500">
//                   <option>Select your product</option>
//                 </select>
//                 <div className="absolute inset-y-0 right-0 flex items-center px-3 text-gray-500 pointer-events-none">
//                   <svg
//                     className="w-5 h-5"
//                     fill="none"
//                     stroke="currentColor"
//                     viewBox="0 0 24 24"
//                   >
//                     <path
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                       strokeWidth={2}
//                       d="M19 9l-7 7-7-7"
//                     />
//                   </svg>
//                 </div>
//               </div>
//             </div>

//             <div>
//               <label className="block text-sm font-medium mb-3">
//                 Select your color
//               </label>
//               <div className="flex space-x-2">
//                 {["#D9F3D9", "#FFD8D8", "#D8E7F3", "#FFF8D8", "#2F3133"].map(
//                   (color) => (
//                     <button
//                       key={color}
//                       className={`w-8 h-8 rounded-md ${
//                         product.color === color
//                           ? "ring-2 ring-offset-2 ring-blue-500"
//                           : ""
//                       }`}
//                       style={{ backgroundColor: color }}
//                       onClick={() => handleColorSelect(color)}
//                     ></button>
//                   )
//                 )}
//               </div>
//             </div>
//           </div>
//         </div>
//       </div>
//     </div>
//   );
// };
// components/admin/product/add/AddProduct.tsx
import React from "react";
import { ProductHeader } from "./components/ProductHeader";
import { BasicDetails } from "./components/BasicDetails";
import { ProductPricing } from "./components/ProductPricing";
import { ProductInventory } from "./components/ProductInventory";
import { ProductImages } from "./components/ProductImages";
import { ProductCategories } from "./components/ProductCategories";
import { useProductForm } from "./hooks/useProductForm";
import { useImageUpload } from "./hooks/useImageUpload";

export const AddProduct: React.FC = () => {
  const {
    product,
    handleInputChange,
    handleRadioChange,
    handleColorSelect,
    handlePublish,
    handleSaveDraft,
    calculateSavings
  } = useProductForm();

  const {
    selectedImages,
    mainImage,
    handleImageUpload,
    removeImage,
    setMainImage
  } = useImageUpload();

  return (
    <div className="container mx-auto">
      <ProductHeader 
        onPublish={handlePublish} 
        onSaveDraft={handleSaveDraft} 
      />

      <div className="flex flex-wrap -mx-3 mt-[1.5rem] drop-shadow-sm filter">
        <div className="w-[611px] h-auto md:w-7/12 px-3 mb-6">
          <BasicDetails 
            product={product} 
            onChange={handleInputChange} 
          />
          
          <ProductPricing 
            product={product}
            onChange={handleInputChange}
            discountAmount={calculateSavings()}
          />
          
          <ProductInventory 
            product={product}
            onChange={handleInputChange}
            onRadioChange={handleRadioChange}
            handlePublish={handlePublish}
            handleSaveDraft={handleSaveDraft}
          />

        </div>

        <div className="w-full md:w-5/12 px-3">
          <ProductImages 
            selectedImages={selectedImages}
            mainImage={mainImage}
            onImageUpload={handleImageUpload}
            onRemoveImage={removeImage}
            onSetMainImage={setMainImage}
          />
          
          <ProductCategories 
            product={product}
            onColorSelect={handleColorSelect}
          />
        </div>
      </div>
    </div>
  );
};