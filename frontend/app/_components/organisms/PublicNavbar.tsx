'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Box, HStack, Button, Text } from '@/src/core/ui'; // Button, Text等は src/core/ui からインポート
import { useTranslations } from 'next-intl';
import { css } from 'styled-system/css';

export const PublicNavbar = () => {
    const pathname = usePathname();
    const t = useTranslations();

    // Placeholder for auth state - to be implemented properly later
    const isLoggedIn = false;
    const handleLogout = () => {
        // Implement logout logic
        console.log('Logout');
    };

    // 管理画面では表示しない
    if (pathname?.startsWith('/admin')) {
        return null;
    }

    // ログイン画面でもシンプルにするため非表示などの検討余地はあるが、
    // いったんナビゲーションはあっても便利なので表示する方針とする。
    // 必要であれば if (pathname?.startsWith('/login')) return null; 等を追加。

    return (
        <Box
            as="header"
            position="sticky"
            top="0"
            zIndex="sticky"
            w="full"
            bg="white/90" // 透過設定
            backdropFilter="blur(8px)"
            shadow="sm"
            borderBottom="1px solid"
            borderColor="gray.100"
        >
            <Box maxW="7xl" mx="auto" px={{ base: '4', md: '8' }} h="16" display="flex" alignItems="center" justifyContent="space-between">
                <Link href="/" className={css({ textDecoration: 'none', _hover: { opacity: 0.8 } })}>
                    <HStack spacing="2">
                        {/* 簡易的なロゴ表示 (画像があればImageコンポーネントに差し替え) */}
                        <Box fontSize="2xl" role="img" aria-label="Logo">
                            🐟
                        </Box>
                        <Text fontSize="lg" fontWeight="bold" className={css({ color: 'indigo.900' })} display={{ base: 'none', sm: 'block' }}>
                            {t('Common.app_name')}
                        </Text>
                    </HStack>
                </Link>

                <HStack spacing="6">
                    {/* デスクトップナビゲーション */}
                    <Box display={{ base: 'none', md: 'block' }}>
                        <HStack spacing="6">
                            <NavLink href="/auctions">{t('Navbar.active_auctions')}</NavLink>
                            {isLoggedIn && (
                                <NavLink href="/mypage">{t('Navbar.mypage')}</NavLink>
                            )}
                        </HStack>
                    </Box>

                    {/* 認証ボタン */}
                    {/* ログイン状態に応じた出し分けが必要だが、ここはいったんリンクベースで配置 */}
                    {/* 必要に応じて useAuth フックなどで状態監視して出し分ける */}
                    <HStack spacing="3">
                        {isLoggedIn ? (
                            <Button size="sm" variant="outline" onClick={handleLogout}>
                                {t('Navbar.logout')}
                            </Button>
                        ) : (
                            <Link href="/login/buyer">
                                <Button size="sm" variant="primary">
                                    {t('Navbar.login')}
                                </Button>
                            </Link>
                        )}
                    </HStack>
                </HStack>
            </Box>
        </Box>
    );
};

// 内部用 NavLink コンポーネント
const NavLink = ({ href, children }: { href: string; children: React.ReactNode }) => {
    const pathname = usePathname();
    const isActive = pathname === href || (href !== '/' && pathname?.startsWith(href));

    return (
        <Link
            href={href}
            className={css({
                px: '4',
                py: '2',
                borderRadius: 'full',
                fontSize: 'sm',
                fontWeight: 'medium',
                color: isActive ? 'indigo.700' : 'gray.600',
                bg: isActive ? 'indigo.50' : 'transparent',
                transition: 'all 0.2s',
                _hover: {
                    color: 'indigo.700',
                    bg: 'indigo.50'
                }
            })}
        >
            {children}
        </Link>
    );
};
