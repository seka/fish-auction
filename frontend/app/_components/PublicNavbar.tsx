'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Box, HStack, Button, Text } from '@/src/core/ui'; // Button, Text等は src/core/ui からインポート
import { useTranslations } from 'next-intl';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { css } from 'styled-system/css';

export const PublicNavbar = () => {
    const t = useTranslations();
    const pathname = usePathname();

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
                {/* Logo / Brand */}
                <Link href="/" className={css({ textDecoration: 'none', _hover: { opacity: 0.8 }, transition: 'opacity 0.2s' })}>
                    <HStack spacing="3">
                        {/* 簡易的なロゴ表示 (画像があればImageコンポーネントに差し替え) */}
                        <Box bg="gradient-to-br from-blue.500 to-indigo.600" w="8" h="8" borderRadius="md" display="flex" alignItems="center" justifyContent="center" color="white" fontWeight="bold">
                            🐟
                        </Box>
                        <Text fontWeight="bold" fontSize="lg" className={css({ color: 'gray.900', letterSpacing: 'tight' })}>
                            漁港のせりシステム
                        </Text>
                    </HStack>
                </Link>

                {/* Navigation Links */}
                <HStack spacing="1" display={{ base: 'none', md: 'flex' }}>
                    <NavLink href="/">{t(COMMON_TEXT_KEYS.home)}</NavLink>
                    <NavLink href="/auctions">{t(COMMON_TEXT_KEYS.auction_venue)}</NavLink>
                    <NavLink href="/mypage">{t(COMMON_TEXT_KEYS.mypage)}</NavLink>
                </HStack>

                {/* Mobile Menu Button (Future work if needed) */}
                {/* <Box display={{ base: 'block', md: 'none' }}>
                    <Button variant="ghost" size="sm">Menu</Button>
                </Box> */}

                {/* Action Buttons */}
                <HStack spacing="4">
                    {/* ログイン状態に応じた出し分けが必要だが、ここはいったんリンクベースで配置 */}
                    {/* 必要に応じて useAuth フックなどで状態監視して出し分ける */}
                    <Link href="/login/buyer">
                        <Button size="sm" className={css({ bg: 'gray.600', _hover: { bg: 'gray.700' }, color: 'white', fontWeight: 'medium' })}>
                            {t(COMMON_TEXT_KEYS.login)}
                        </Button>
                    </Link>
                    <Link href="/signup">
                        <Button size="sm" className={css({ bg: 'indigo.600', color: 'white', _hover: { bg: 'indigo.700' }, fontWeight: 'bold', px: '6' })}>
                            {t(COMMON_TEXT_KEYS.signup)}
                        </Button>
                    </Link>
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
