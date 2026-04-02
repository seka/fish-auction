'use client';

import { PublicNavbar } from '@organisms';
import { useBuyerAuth } from '../queries/useAuth';

/**
 * 認証ロジックを注入したナビゲーションバー
 */
export const AuthorizablePublicNavbar = () => {
  const { buyer, isLoggedIn, isLoading, logout } = useBuyerAuth();

  return (
    <PublicNavbar
      isLoggedIn={isLoggedIn}
      isLoading={isLoading}
      buyerName={buyer?.name}
      onLogout={logout}
    />
  );
};
