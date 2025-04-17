// pages/user/SecuritySettings.tsx
import React, { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import UserLayout from '@/components/layouts/UserLayout';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useToast } from '@/hooks/use-toast';
import { Switch } from '@/components/ui/switch';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { FiLock, FiShield, FiAlertCircle } from 'react-icons/fi';

const SecuritySettings: React.FC = () => {
  const { authState } = useAuth();
  const { toast } = useToast();
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [passwordVisible, setPasswordVisible] = useState(false);
  
  // Mock security settings
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false);
  const [emailNotifications, setEmailNotifications] = useState(true);
  
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData(prev => ({ ...prev, [name]: value }));
  };
  
  const handlePasswordSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validate password
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      toast({
        title: 'Passwords do not match',
        description: 'Please make sure your new password and confirmation match.',
        variant: 'destructive'
      });
      return;
    }
    
    if (passwordData.newPassword.length < 8) {
      toast({
        title: 'Password too short',
        description: 'Your password must be at least 8 characters long.',
        variant: 'destructive'
      });
      return;
    }
    
    // Mock password change success
    toast({
      title: 'Password Updated',
      description: 'Your password has been successfully changed.',
    });
    
    // Reset form
    setPasswordData({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    });
  };
  
  const handleToggleTwoFactor = (checked: boolean) => {
    setTwoFactorEnabled(checked);
    toast({
      title: `Two-factor Authentication ${checked ? 'Enabled' : 'Disabled'}`,
      description: checked 
        ? 'Your account is now more secure with 2FA.'
        : 'Two-factor authentication has been turned off.',
    });
  };
  
  return (
    <UserLayout>
      <div className="max-w-4xl mx-auto p-6">
        <h1 className="text-2xl font-bold mb-6">Security Settings</h1>
        
        <div className="space-y-6">
          {/* Password Change */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <FiLock className="text-blue-500" />
                <CardTitle>Change Password</CardTitle>
              </div>
              <CardDescription>
                Update your password to ensure your account stays secure
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handlePasswordSubmit}>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium mb-1">
                      Current Password
                    </label>
                    <div className="relative">
                      <Input
                        type={passwordVisible ? 'text' : 'password'}
                        name="currentPassword"
                        value={passwordData.currentPassword}
                        onChange={handlePasswordChange}
                        required
                      />
                    </div>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium mb-1">
                      New Password
                    </label>
                    <div className="relative">
                      <Input
                        type={passwordVisible ? 'text' : 'password'}
                        name="newPassword"
                        value={passwordData.newPassword}
                        onChange={handlePasswordChange}
                        required
                      />
                    </div>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium mb-1">
                      Confirm New Password
                    </label>
                    <div className="relative">
                      <Input
                        type={passwordVisible ? 'text' : 'password'}
                        name="confirmPassword"
                        value={passwordData.confirmPassword}
                        onChange={handlePasswordChange}
                        required
                      />
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-2">
                    <Switch
                      id="show-password"
                      checked={passwordVisible}
                      onCheckedChange={setPasswordVisible}
                    />
                    <label htmlFor="show-password" className="text-sm">
                      Show password
                    </label>
                  </div>
                </div>
                
                <Button type="submit" className="mt-4">
                  Update Password
                </Button>
              </form>
            </CardContent>
          </Card>
          
          {/* Two-Factor Authentication */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <FiShield className="text-blue-500" />
                <CardTitle>Two-Factor Authentication</CardTitle>
              </div>
              <CardDescription>
                Add an extra layer of security to your account
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Two-factor authentication</p>
                  <p className="text-sm text-gray-500">
                    Secure your account with 2FA
                  </p>
                </div>
                <Switch 
                  checked={twoFactorEnabled}
                  onCheckedChange={handleToggleTwoFactor}
                />
              </div>
            </CardContent>
          </Card>
          
          {/* Account Activity */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <FiAlertCircle className="text-blue-500" />
                <CardTitle>Account Activity</CardTitle>
              </div>
              <CardDescription>
                Monitor and manage your account activity
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between mb-4">
                <div>
                  <p className="font-medium">Email notifications for new logins</p>
                  <p className="text-sm text-gray-500">
                    Receive an email when a new device logs into your account
                  </p>
                </div>
                <Switch 
                  checked={emailNotifications}
                  onCheckedChange={setEmailNotifications}
                />
              </div>
              
              <Button variant="outline">
                View Login History
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </UserLayout>
  );
};

export default SecuritySettings;