// components/admin/role/components/ProfileUpdateForm.tsx
import React, { useState } from 'react';
import { AdminRole } from '../models/adminRole.model'
import { FiEdit } from "react-icons/fi";
interface ProfileUpdateFormProps {
  adminData: AdminRole;
  onUpdateProfile: (data: Partial<AdminRole>) => Promise<boolean>;
  setAdminData: (data: AdminRole) => void;
}

export const ProfileUpdateForm: React.FC<ProfileUpdateFormProps> = ({ 
  adminData, 
  onUpdateProfile,
  setAdminData 
}) => {
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [formData, setFormData] = useState<AdminRole>(adminData);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    const success = await onUpdateProfile(formData);
    if (success) {
      setIsEditing(false);
    }
  };
  

  const handleCancel = () => {
    setFormData(adminData);
    setIsEditing(false);
  };

  return (
    <div className="bg-white rounded-lg shadow p-[1.5rem] w-[738px] h-[856px] drop-shadow-md filter ml-[1.5rem]">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-[18px] text-cyprus font-bold">Profile Update</h3>
        {!isEditing && (
          <button
            onClick={() => setIsEditing(true)}
            className="flex items-center justify-center gap-4 w-[84px] h-[32px] py-4 pl-8 pr-12 bg-white border border-gray-300 text-neutral-500 font-medium rounded hover:bg-neutral-50 transition-colors"
          >
            <FiEdit></FiEdit>
            <span>Edit</span>
          </button>
        )}
      </div>

      <div className="mb-[1rem] flex items-center mt-[1rem] w-[340px] h-[64px]">
        <div className="relative flex items-center gap-16">
          <img
            src={adminData.profileImage}
            alt={`${adminData.firstName} ${adminData.lastName}`}
            className="w-[64px] h-[64px] rounded-full object-cover"
          />
          {isEditing && (
            <div className="w-[252px] gap-12 flex mt-2 justify-center space-x-2">
              <button className="w-[120px] h-[42px] bg-ocean-green hover:bg-green-600 text-white text-[15px] font-bold px-3 py-1 rounded-lg">
                Upload New
              </button>
              <button className="w-[120px] h-[42px] bg-white border border-gray-300 text-neutral-600 text-[15px] font-medium px-3 py-1 rounded-lg hover:bg-gray-100">
                Delete
              </button>
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-[20px]">
        <div>
          <label htmlFor='firstName' className="block text-[15px] text-cyprus mb-2">
            First Name
          </label>
          <input
            type="text"
            id="firstName"
            className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={formData.firstName}
            onChange={handleChange}
            readOnly={!isEditing}
          />
        </div>

        <div>
          <label htmlFor="lastName" className="block text-[15px] text-cyprus mb-2">
            Last Name
          </label>
          <input
            type="text"
            id="lastName"
            className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={formData.lastName}
            onChange={handleChange}
            readOnly={!isEditing}
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-[15px] text-cyprus mb-2">Password</label>
          <input
            type="password"
            className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value="************"
            readOnly
            id="password"
          />
        </div>

        <div>
          <label htmlFor="phoneNumber" className="block text-[15px] text-cyprus mb-2">
            Phone Number
          </label>
          <div className="flex">
            <div className="border border-gray-300 rounded-l-md px-2 py-2 bg-white flex items-center">
              <img src="/us-flag.svg" alt="US" className="h-5 w-5" />
              <span className="text-sm ml-1">+1</span>
            </div>
            <input
              type="text"
              id="phoneNumber"
              className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={formData.phoneNumber}
              onChange={handleChange}
              readOnly={!isEditing}
            />
          </div>
        </div>

        <div>
          <label htmlFor='email' className="block text-[15px] text-cyprus mb-2">E-mail</label>
          <input
            type="email"
            id="email"
            className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={formData.email}
            onChange={handleChange}
            readOnly={!isEditing}
          />
        </div>

        <div>
          <label htmlFor="dateOfBirth" className="block text-[15px] text-cyprus mb-2">
            Date of Birth
          </label>
          <input
            type="text"
            id="dateOfBirth"
            className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={formData.dateOfBirth}
            onChange={handleChange}
            readOnly={!isEditing}
          />
        </div>

        <div className="col-span-2">
          <label htmlFor="location" className="block text-[15px] text-cyprus mb-2">Location</label>
          <input
            type="text"
            id="location"
            className="text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={formData.location}
            onChange={handleChange}
            readOnly={!isEditing}
          />
        </div>

        <div className="col-span-2">
          <label htmlFor="creditCard" className="block text-[15px] text-cyprus mb-2">
            Credit Card
          </label>
          <div className="relative">
            <div className="mr-5 absolute left-3 top-1/2 transform -translate-y-1/2">
              <svg
                className="h-[24px] w-[38px]"
                xmlns="http://www.w3.org/2000/svg"
                x="0px"
                y="0px"
                width="100"
                height="100"
                viewBox="0 0 48 48"
              >
                <linearGradient
                  id="bvEjD8y7f6LTu_wtNrsRga_tPLW4S6pddUz_gr1"
                  x1="43.913"
                  x2="9.262"
                  y1="24"
                  y2="24"
                  gradientUnits="userSpaceOnUse"
                >
                  <stop offset="0" stop-color="#fce724"></stop>
                  <stop offset=".214" stop-color="#fdd01e"></stop>
                  <stop offset=".661" stop-color="#fe930e"></stop>
                  <stop offset="1" stop-color="#ff6201"></stop>
                </linearGradient>
                <circle
                  cx="34"
                  cy="24"
                  r="14"
                  fill="url(#bvEjD8y7f6LTu_wtNrsRga_tPLW4S6pddUz_gr1)"
                ></circle>
                <path
                  fill="#ff9045"
                  d="M20,24c0,4.529,2.16,8.544,5.495,11.103C28.278,32.223,30,28.311,30,24 s-1.722-8.223-4.505-11.103C22.16,15.456,20,19.471,20,24z"
                  opacity=".05"
                ></path>
                <path
                  fill="#ff9045"
                  d="M20,24c0,4.385,2.019,8.296,5.175,10.863C27.908,32.052,29.6,28.223,29.6,24 s-1.692-8.052-4.425-10.863C22.019,15.704,20,19.615,20,24z"
                  opacity=".1"
                ></path>
                <path
                  fill="#ff9045"
                  d="M20,24c0,4.242,1.894,8.035,4.874,10.603C27.546,31.862,29.2,28.124,29.2,24 s-1.654-7.862-4.326-10.603C21.894,15.965,20,19.758,20,24z"
                  opacity=".2"
                ></path>
                <path
                  fill="#ff9045"
                  d="M20,24c0,4.099,1.772,7.774,4.58,10.335C27.187,31.666,28.8,28.022,28.8,24 s-1.613-7.666-4.22-10.335C21.772,16.226,20,19.901,20,24z"
                  opacity=".3"
                ></path>
                <path
                  fill="#ff9045"
                  d="M20,24c0,3.956,1.647,7.521,4.285,10.067C26.827,31.47,28.4,27.92,28.4,24 c0-3.92-1.573-7.47-4.115-10.067C21.647,16.479,20,20.044,20,24z"
                  opacity=".4"
                ></path>
                <path
                  fill="#fe8c0c"
                  d="M28,24c0-4.892-2.514-9.185-6.317-11.684C18.789,15.366,17,19.474,17,24 s1.789,8.634,4.683,11.684C25.486,33.185,28,28.892,28,24z"
                  opacity=".1"
                ></path>
                <path
                  fill="#fe8c0c"
                  d="M28,24c0-4.621-2.25-8.702-5.704-11.247C19.487,15.675,17.75,19.634,17.75,24 s1.737,8.325,4.546,11.247C25.75,32.702,28,28.621,28,24z"
                  opacity=".2"
                ></path>
                <path
                  fill="#fe8c0c"
                  d="M28,24c0-4.352-1.991-8.229-5.108-10.792C20.179,16,18.5,19.804,18.5,24 s1.679,8,4.392,10.792C26.009,32.229,28,28.352,28,24z"
                  opacity=".3"
                ></path>
                <path
                  fill="#fe8c0c"
                  d="M28,24c0-4.083-1.757-7.744-4.546-10.299C20.857,16.36,19.25,19.992,19.25,24 c0,4.008,1.607,7.64,4.204,10.299C26.243,31.744,28,28.083,28,24z"
                  opacity=".4"
                ></path>
                <path
                  fill="#ff9045"
                  d="M20,24c0,3.812,1.528,7.264,4,9.789c2.472-2.525,4-5.977,4-9.789s-1.528-7.264-4-9.789	C21.528,16.736,20,20.188,20,24z"
                ></path>
                <linearGradient
                  id="bvEjD8y7f6LTu_wtNrsRgb_tPLW4S6pddUz_gr2"
                  x1="0"
                  x2="41.005"
                  y1="24"
                  y2="24"
                  gradientUnits="userSpaceOnUse"
                >
                  <stop
                    offset="0"
                    stop-color="#f4805d"
                    stop-opacity=".4"
                  ></stop>
                  <stop
                    offset="1"
                    stop-color="#ffd8bb"
                    stop-opacity=".4"
                  ></stop>
                </linearGradient>
                <circle
                  cx="14"
                  cy="24"
                  r="14"
                  fill="url(#bvEjD8y7f6LTu_wtNrsRgb_tPLW4S6pddUz_gr2)"
                ></circle>
                <linearGradient
                  id="bvEjD8y7f6LTu_wtNrsRgc_tPLW4S6pddUz_gr3"
                  x1="22.885"
                  x2="4.248"
                  y1="14.038"
                  y2="34.934"
                  gradientUnits="userSpaceOnUse"
                >
                  <stop offset="0" stop-color="ivory" stop-opacity=".2"></stop>
                  <stop
                    offset="1"
                    stop-color="#ff9a83"
                    stop-opacity=".4"
                  ></stop>
                </linearGradient>
                <path
                  fill="url(#bvEjD8y7f6LTu_wtNrsRgc_tPLW4S6pddUz_gr3)"
                  d="M14,10.5c7.444,0,13.5,6.056,13.5,13.5	S21.444,37.5,14,37.5S0.5,31.444,0.5,24S6.556,10.5,14,10.5 M14,10C6.268,10,0,16.268,0,24s6.268,14,14,14s14-6.268,14-14	S21.732,10,14,10L14,10z"
                ></path>
              </svg>
            </div>
            <input
              type="text"
              id="creditCard"
              className="pl-[50px] text-cyprus w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={formData.creditCard}
              onChange={handleChange}
              readOnly={!isEditing}
            />
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-[24px] w-[38px] text-gray-400"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
          </div>
        </div>

        <div className="col-span-2 relative h-[147px]">
          <label htmlFor="biography" className="block text-[15px] text-cyprus mb-2">
            Biography
          </label>
          <textarea
            id="biography"
            className="text-cyprus h-[117px] w-full border bg-neutral-50 border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
            rows={3}
            placeholder="Enter a biography about you"
            value={formData.biography || ""}
            onChange={handleChange}
            readOnly={!isEditing}
          ></textarea>
          {!isEditing && (
            <div className="absolute bottom-2 right-4 col-span-2 flex justify-end space-x-2">
              <button className=" text-neutral-500 px-4 py-2 rounded hover:text-cyprus transition-colors">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-[24px] w-[24px]"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                </svg>
              </button>
              <button className="text-neutral-500  px-4 py-2 rounded hover:text-cyprus transition-colors">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-[24px] w-[24px]"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M4 2a2 2 0 00-2 2v11a3 3 0 106 0V4a2 2 0 00-2-2H4zm1 14a1 1 0 100-2 1 1 0 000 2zm5-1.757l4.9-4.9a2 2 0 000-2.828L13.485 5.1a2 2 0 00-2.828 0L10 5.757v8.486zM16 18H9.071l6-6H16a2 2 0 012 2v2a2 2 0 01-2 2z"
                    clipRule="evenodd"
                  />
                </svg>
              </button>
            </div>
          )}
        </div>

        {isEditing && (
          <div className="col-span-2 flex justify-end space-x-2">
            <button
              onClick={handleCancel}
              className="w-[110px] bg-white border border-gray-300 text-cyprus px-4 py-2 rounded hover:bg-gray-100 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="w-[110px] bg-ocean-green hover:bg-green-600 text-white px-4 py-2 rounded transition-colors"
            >
              Save
            </button>
          </div>
        )}
      </div>
    </div>
  );
};