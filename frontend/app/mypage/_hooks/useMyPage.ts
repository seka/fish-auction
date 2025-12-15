import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { getMyPurchases, getMyAuctions } from '@/src/api/buyer_mypage';
import { logoutBuyer } from '@/src/api/buyer_auth';

export const useMyPage = () => {
    const t = useTranslations();
    const router = useRouter();
    const [activeTab, setActiveTab] = useState<'purchases' | 'auctions' | 'settings'>('purchases');

    // Password state
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [isPasswordUpdating, setIsPasswordUpdating] = useState(false);

    // Fetch purchase history
    const {
        data: purchases = [],
        isLoading: isPurchasesLoading
    } = useQuery({
        queryKey: ['purchases'],
        queryFn: getMyPurchases,
    });

    // Fetch participating auctions
    const {
        data: auctions = [],
        isLoading: isAuctionsLoading
    } = useQuery({
        queryKey: ['auctions', 'my'],
        queryFn: getMyAuctions,
    });

    const isLoading = isPurchasesLoading || isAuctionsLoading;

    const handleLogout = async () => {
        const success = await logoutBuyer();
        if (success) {
            router.push('/login/buyer');
        }
    };

    const handleUpdatePassword = async (e: React.FormEvent) => {
        e.preventDefault();
        setMessage(null);

        if (newPassword !== confirmPassword) {
            setMessage({ type: 'error', text: '新しいパスワードが一致しません。' });
            return;
        }

        if (newPassword.length < 8) {
            setMessage({ type: 'error', text: 'パスワードは8文字以上である必要があります。' });
            return;
        }

        setIsPasswordUpdating(true);

        try {
            const res = await fetch('/api/proxy/api/buyers/password', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    current_password: currentPassword,
                    new_password: newPassword,
                }),
            });

            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.error || 'パスワードの更新に失敗しました。');
            }

            setMessage({ type: 'success', text: 'パスワードを更新しました。' });
            setCurrentPassword('');
            setNewPassword('');
            setConfirmPassword('');
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message });
        } finally {
            setIsPasswordUpdating(false);
        }
    };

    return {
        t,
        activeTab,
        setActiveTab,
        currentPassword,
        setCurrentPassword,
        newPassword,
        setNewPassword,
        confirmPassword,
        setConfirmPassword,
        message,
        isPasswordUpdating,
        purchases,
        auctions,
        isLoading,
        handleLogout,
        handleUpdatePassword,
    };
};
