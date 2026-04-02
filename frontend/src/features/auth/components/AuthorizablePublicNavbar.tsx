'use client';

import { PublicNavbar } from '@organisms';
import { usePathname } from 'next/navigation';
import { useBuyerAuth } from '../queries/useAuth';

/**
 * 認証ロジックを注入したナビゲーションバー
 */
export const AuthorizablePublicNavbar = () => {
  const pathname = usePathname();
  const { buyer, isLoggedIn, isLoading, logout } = useBuyerAuth();

  if (pathname?.startsWith('/admin')) {
    return null;
  }

  return (
    <PublicNavbar
      isLoggedIn={isLoggedIn}
      isLoading={isLoading}
      buyerName={buyer?.name}
      onLogout={logout}
    />
  );
};
