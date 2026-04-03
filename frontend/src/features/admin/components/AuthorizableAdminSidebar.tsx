'use client';

import { useAdminLogoutMutation } from '@/src/features/auth';
import { AdminSidebar } from '@organisms';
import { useRouter } from 'next/navigation';
import { useAdminNavigation } from '../states/useAdminNavigation';

/**
 * 認証ロジックを注入した管理者用サイドバー
 */
export const AuthorizableAdminSidebar = () => {
  const router = useRouter();
  const logoutMutation = useAdminLogoutMutation();
  const { getIsActive } = useAdminNavigation();

  const handleLogout = async () => {
    await logoutMutation.mutateAsync();
    // note: mutationHook 側で queryClient.clear() を行っているため
    // フルリロード (window.location.href) をしなくても、クリーンな状態で遷移できる
    router.push('/login/admin');
  };

  return <AdminSidebar onLogout={handleLogout} getIsActive={getIsActive} />;
};
