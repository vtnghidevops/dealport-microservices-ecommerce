// components/user/UserAvatar.tsx
import React, { useState } from 'react';
import { User } from '@/types/user.model';
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

    // If user has username but no profile info
    if (user.username && (!user.profile?.firstName && !user.profile?.lastName)) {
      return user.username.substring(0, 2).toUpperCase();
    }

    const firstName = user.profile?.firstName || '';
    const lastName = user.profile?.lastName || '';

    if (!firstName && !lastName) {
      // If no name data at all, use first 1-2 chars of email
      return user.email ? user.email.substring(0, 2).toUpperCase() : '?';
    }

    return `${firstName[0] || ''}${lastName[0] || ''}`.toUpperCase();
  };

  // Determine avatar URL with appropriate fallbacks
  const avatarUrl = preview || user?.profile?.avatar || null;

  // Get user's name for alt text with fallbacks
  const userName = user ?
    (user.profile?.firstName || user.profile?.lastName) ?
      `${user.profile?.firstName || ''} ${user.profile?.lastName || ''}`.trim() :
      user.username || user.email || 'User' :
    'User';

  return (
    <div className="flex flex-col items-center">
      <div
        className={`${sizeClasses[size]} relative rounded-full overflow-hidden ${editable ? 'cursor-pointer' : ''}`}
        onClick={handleImageClick}
      >
        {avatarUrl ? (
          <img
            src={avatarUrl}
            alt={userName}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className={`${sizeClasses[size]} flex items-center justify-center bg-blue-100 text-blue-600 text-2xl font-semibold`}>
            {getInitials()}
          </div>
        )}

        {editable && (
          <>
            <div className="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity">
              <FiCamera className="text-white text-xl" />
            </div>
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleFileChange}
              accept="image/*"
              className="hidden"
            />
          </>
        )}
      </div>

      {/* Optional: Display user name below avatar */}
      {size === 'lg' || size === 'xl' ? (
        <p className="mt-2 text-sm font-medium text-gray-700">{userName}</p>
      ) : null}
    </div>
  );
};

export default UserAvatar;