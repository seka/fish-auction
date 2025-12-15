'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Box, Stack, Text, HStack } from '@/src/core/ui';
import { css } from 'styled-system/css';

// 共通のスタイル定義 (Recipe的アプローチ)
const sidebarItemStyles = {
    base: css.raw({
        display: 'block',
        py: '3',
        px: '4',
        borderRadius: 'md',
        fontSize: 'sm',
        fontWeight: 'medium',
        color: 'indigo.100',
        cursor: 'pointer',
        transition: 'all 0.2s',
        w: 'full',
        textAlign: 'left',
        _hover: {
            bg: 'indigo.800',
            color: 'white'
        },
    }),
    active: css.raw({
        bg: 'indigo.800',
        color: 'white',
        fontWeight: 'bold',
        position: 'relative',
        _before: {
            content: '""',
            position: 'absolute',
            left: '0',
            top: '0',
            bottom: '0',
            width: '4px',
            bg: 'indigo.400',
            borderTopLeftRadius: 'md',
            borderBottomLeftRadius: 'md',
        }
    })
};

type SidebarItemProps = {
    children: React.ReactNode;
    href?: string;         // Linkとして使う場合
    onClick?: () => void;  // Buttonとして使う場合
    icon?: string;         // アイコン (絵文字など)
    isActive?: boolean;    // 明示的にActiveにする場合 (基本はhrefで自動判定)
};

const SidebarItem = ({ children, href, onClick, icon, isActive: explicitActive }: SidebarItemProps) => {
    const pathname = usePathname();

    // hrefがある場合は、現在のパスがhrefで始まっているかを判定 (サブパスも含めるため startsWith を使用)
    const isActive = explicitActive ?? (href ? (href === '/admin' ? pathname === '/admin' : pathname.startsWith(href)) : false);

    const className = css(
        sidebarItemStyles.base,
        isActive ? sidebarItemStyles.active : {}
    );

    const content = (
        <HStack spacing="3">
            {icon && <span className={css({ w: '6', textAlign: 'center' })}>{icon}</span>}
            <span>{children}</span>
        </HStack>
    );

    if (href) {
        return (
            <Link href={href} className={className}>
                {content}
            </Link>
        );
    }

    return (
        <button onClick={onClick} className={className}>
            {content}
        </button>
    );
};

export const Sidebar = () => {
    return (
        <Box w="64" bg="indigo.900" color="white" flexShrink={0} shadow="xl" display="flex" flexDirection="column" h="full">
            <Box p="6" bg="indigo.950">
                <Text as="h2" fontSize="xl" fontWeight="bold" letterSpacing="wider" className={css({ color: 'white' })}>管理画面</Text>
                <Text fontSize="xs" className={css({ color: 'indigo.300' })} mt="1">Fish Auction Admin</Text>
            </Box>

            <Stack as="nav" mt="6" px="2" spacing="1" flex="1">
                <SidebarItem href="/" icon="↩️">
                    トップに戻る
                </SidebarItem>

                <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

                <SidebarItem href="/admin" icon="📊">
                    ダッシュボード
                </SidebarItem>

                <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

                <SidebarItem href="/admin/fishermen" icon="👨‍🌾">
                    漁師管理
                </SidebarItem>

                <SidebarItem href="/admin/buyers" icon="👔">
                    中買人管理
                </SidebarItem>

                <SidebarItem href="/admin/venues" icon="🏢">
                    会場管理
                </SidebarItem>

                <SidebarItem href="/admin/auctions" icon="📅">
                    セリ管理
                </SidebarItem>

                <SidebarItem href="/admin/items" icon="🐟">
                    出品管理
                </SidebarItem>

                <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

                <SidebarItem href="/admin/invoice" icon="💰">
                    請求書発行
                </SidebarItem>

                <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

                <SidebarItem href="/admin/settings" icon="⚙️">
                    設定
                </SidebarItem>
            </Stack>

            {/* Footer / User info could go here */}
            <Box p="4" bg="indigo.950" fontSize="xs" className={css({ color: 'indigo.400', textAlign: 'center' })}>
                &copy; Fish Auction System
            </Box>
        </Box>
    );
};
