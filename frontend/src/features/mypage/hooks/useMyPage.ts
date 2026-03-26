import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { logoutBuyer } from '@/src/data/api/buyer_auth';
import { authKeys } from '@/src/data/queries/auth/keys';

import { MyPageTab } from '../types';

export const useMyPage = () => {
  const t = useTranslations();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<MyPageTab>('purchases');

  // Password state
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordMessage, setPasswordMessage] = useState({
    text: '',
    type: 'info' as 'info' | 'error' | 'success',
  });

  const handleLogout = async () => {
    const success = await logoutBuyer();
    if (success) {
      await queryClient.invalidateQueries({ queryKey: authKeys.me() });
      router.push('/login/buyer');
    }
  };

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
