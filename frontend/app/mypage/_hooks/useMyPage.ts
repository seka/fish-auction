import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { getMyPurchases, getMyAuctions } from '@/src/api/buyer_mypage';
import { logoutBuyer } from '@/src/api/buyer_auth';
import { authKeys } from '@/src/hooks/auth/keys';
import { auctionKeys } from '@/src/hooks/auction/keys';
import { buyerKeys } from '@/src/hooks/buyer/keys';

export const useMyPage = () => {
  const t = useTranslations();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'purchases' | 'auctions' | 'settings'>('purchases');

  // Password state
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordMessage, setPasswordMessage] = useState({ text: '', type: 'info' as 'info' | 'error' | 'success' });

  // Fetch purchases
  const { data: purchases = [], isLoading: isPurchasesLoading } = useQuery({
    queryKey: buyerKeys.purchases,
    queryFn: getMyPurchases,
  });

  // Fetch participating auctions
  const { data: auctions = [], isLoading: isAuctionsLoading } = useQuery({
    queryKey: auctionKeys.lists(),
    queryFn: getMyAuctions,
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
    purchases,
    auctions,
    isLoading: isPurchasesLoading || isAuctionsLoading,
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
    },
  };
};
