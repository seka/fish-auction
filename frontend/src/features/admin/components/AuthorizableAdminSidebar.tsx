'use client';

import { useAdminLogoutMutation } from '@/src/features/auth';
import { AdminSidebar } from '@/src/core/components/organisms/admin/AdminSidebar';

/**
 * 認証ロジックを注入した管理者用サイドバー
 */
export const AuthorizableAdminSidebar = () => {
  const logoutMutation = useAdminLogoutMutation();

  const handleLogout = async () => {
    await logoutMutation.mutateAsync();
    window.location.href = '/login/admin';
  };

  return <AdminSidebar onLogout={handleLogout} />;
};
