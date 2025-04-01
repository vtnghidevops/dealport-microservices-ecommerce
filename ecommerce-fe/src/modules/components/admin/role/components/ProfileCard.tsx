// components/admin/role/components/ProfileCard.tsx
import React from 'react';
import { AdminRole } from "../models/adminRole.model"
import { SocialMediaLinks } from './SocialMediaLinks';
import { MdOutlineContentCopy } from "react-icons/md";
import { PasswordChangeForm } from './PasswordChangeForm';
interface ProfileCardProps {
  adminData: AdminRole;
}

export const ProfileCard: React.FC<ProfileCardProps> = ({ adminData }) => {
  return (
    <div className="flex flex-col bg-white w-[360px] h-[356px] rounded-lg drop-shadow-md filter p-[1rem]">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-[18px] font-bold text-cyprus">Profile</h3>
        <div className="flex space-x-2 w-[64px] h-[32px] -mr-[1rem]">
          <button className="text-neutral-500 hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
          </button>
          <button className="text-neutral-600 hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
            </svg>
          </button>
        </div>
      </div>
      
      <div className="flex flex-col items-center justify-center">
        <div className="relative mb-4 w-[96px] h-[96px]">
          <img 
            src={adminData.profileImage} 
            alt={`${adminData.firstName} ${adminData.lastName}`} 
            className="w-full h-full rounded-full object-cover"
          />
        </div>
        
        <h4 className="text-lg font-medium">{adminData.firstName} {adminData.lastName}</h4>
        <button className='flex gap-8 items-center mb-[1rem]'>
          <p className="text-neutral-500 text-sm mb-4">{adminData.email}</p>
          <MdOutlineContentCopy className='text-primary'></MdOutlineContentCopy>
        </button>
        
        <div className=''>
          <SocialMediaLinks socialMedia={adminData.socialMedia} />
        </div>
      </div>
      
    </div>
  );
};