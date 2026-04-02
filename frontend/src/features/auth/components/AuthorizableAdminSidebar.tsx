'use client';

import { useAdminLogoutMutation } from '../queries/useAuth';
import { Sidebar } from '@/src/core/components/organisms/admin/Sidebar';

/**
 * 認証ロジックを注入した管理者用サイドバー
 */
export const AuthorizableAdminSidebar = () => {
  const logoutMutation = useAdminLogoutMutation();

  const handleLogout = async () => {
    await logoutMutation.mutateAsync();
    window.location.href = '/login/admin';
  };

  return <Sidebar onLogout={handleLogout} />;
};
