// components/user/UserAvatar.tsx
import React, { useState } from 'react';
import { User } from '@/types/user.model';
import { Button } from '@/components/ui/button';
import { FiCamera } from 'react-icons/fi';

interface UserAvatarProps {
  user: User | null;
  editable?: boolean;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  onImageChange?: (file: File) => void;
}

const UserAvatar: React.FC<UserAvatarProps> = ({
  user,
  editable = false,
  size = 'md',
  onImageChange
}) => {
  const [preview, setPreview] = useState<string | null>(null);
  const fileInputRef = React.useRef<HTMLInputElement>(null);

  const sizeClasses = {
    sm: 'w-12 h-12',
    md: 'w-20 h-20',
    lg: 'w-32 h-32',
    xl: 'w-[120px] h-[120px]'

  };

  const handleImageClick = () => {
    if (editable && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const reader = new FileReader();

      reader.onloadend = () => {
        setPreview(reader.result as string);
      };

      reader.readAsDataURL(file);

      if (onImageChange) {
        onImageChange(file);
      }
    }
  };

  // Get initials for fallback avatar
  const getInitials = () => {
    if (!user) return '?';

    const firstName = user.profile?.firstName || '';
    const lastName = user.profile?.lastName || '';

    return `${firstName[0] || ''}${lastName[0] || ''}`.toUpperCase();
  };

  const avatarUrl = preview || user?.profile?.avatar || "/images/system/default-avatars.png";
  const userName = user ? `${user.profile?.firstName || ''} ${user.profile?.lastName || ''}`.trim() : '';

  return (
    <div className="flex flex-col items-center">
      <div
        className={`${sizeClasses[size]} relative rounded-full overflow-hidden ${editable ? 'cursor-pointer' : ''
          }`}
        onClick={handleImageClick}
      >
        {avatarUrl ? (
          <img
            src={avatarUrl}
            alt={userName || 'User avatar'}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className={`${sizeClasses[size]} flex items-center justify-center bg-blue-100 text-blue-600 text-2xl font-semibold`}>
            {getInitials()}
          </div>
        )}

        {editable && (
          <div className="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity">
            <FiCamera className="text-white text-xl" />
          </div>
        )}
      </div>
    </div>
  );
};

export default UserAvatar;