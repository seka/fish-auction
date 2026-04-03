'use client';

import { usePathname } from 'next/navigation';

/**
 * 管理画面のナビゲーション状態を管理するフック
 */
export const useAdminNavigation = () => {
  const pathname = usePathname();

  const getIsActive = (href?: string, explicitActive?: boolean) => {
    if (explicitActive !== undefined) return explicitActive;
    if (!href) return false;

    // トップレベルのパスは完全一致で判定する
    if (href === '/' || href === '/admin') {
      return pathname === href;
    }

    return pathname.startsWith(href);
  };

  return {
    pathname,
    getIsActive,
  };
};
