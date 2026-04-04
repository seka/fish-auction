import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useMyPageAuth } from '../queries/useAuth';

import { MyPageTab } from '../types';

export const useMyPage = () => {
  const t = useTranslations();
  const { handleLogout } = useMyPageAuth();
  const [activeTab, setActiveTab] = useState<MyPageTab>('purchases');

  // Password state
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordMessage, setPasswordMessage] = useState({
    text: '',
    type: 'info' as 'info' | 'error' | 'success',
  });

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      setPasswordMessage({ text: t('MyPage.Settings.password_mismatch'), type: 'error' });
      return;
    }
    // Simulation
    setPasswordMessage({ text: t('MyPage.Settings.password_updated'), type: 'success' });
    setCurrentPassword('');
    setNewPassword('');
    setConfirmPassword('');
  };

  return {
    t,
    activeTab,
    setActiveTab,
    handleLogout,
    passwordState: {
      currentPassword,
      setCurrentPassword,
      newPassword,
      setNewPassword,
      confirmPassword,
      setConfirmPassword,
      passwordMessage,
      handleUpdatePassword,
      isPasswordUpdating: false,
    },
  };
};
