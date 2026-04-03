'use client';

import Link from 'next/link';
import { Box, Stack, Text, HStack } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

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
      color: 'white',
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
    },
  }),
};

type AdminSidebarItemProps = {
  children: React.ReactNode;
  href?: string; // Linkとして使う場合
  onClick?: () => void; // Buttonとして使う場合
  icon?: string; // アイコン (絵文字など)
  isActive?: boolean; // アクティブ状態ならtrue
};

const AdminSidebarItem = ({ children, href, onClick, icon, isActive }: AdminSidebarItemProps) => {
  const className = css(sidebarItemStyles.base, isActive ? sidebarItemStyles.active : {});

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

export interface AdminSidebarProps {
  onLogout: () => Promise<void>;
  getIsActive: (href?: string, explicitActive?: boolean) => boolean;
}

export const AdminSidebar = ({ onLogout, getIsActive }: AdminSidebarProps) => {
  const t = useTranslations('Admin.Sidebar');

  return (
    <Box
      w="64"
      bg="indigo.900"
      color="white"
      flexShrink={0}
      shadow="xl"
      display="flex"
      flexDirection="column"
      h="full"
    >
      <Box p="6" bg="indigo.950">
        <Text
          as="h2"
          fontSize="xl"
          fontWeight="bold"
          letterSpacing="wider"
          className={css({ color: 'white' })}
        >
          {t('title')}
        </Text>
        <Text fontSize="xs" className={css({ color: 'indigo.300' })} mt="1">
          FISHING AUCTION Admin
        </Text>
      </Box>

      <Stack as="nav" mt="6" px="2" spacing="1" flex="1">
        <AdminSidebarItem href="/" icon="↩️" isActive={getIsActive('/')}>
          {t('back_to_top')}
        </AdminSidebarItem>

        <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

        <AdminSidebarItem href="/admin" icon="📊" isActive={getIsActive('/admin')}>
          {t('dashboard')}
        </AdminSidebarItem>

        <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

        <AdminSidebarItem
          href="/admin/fishermen"
          icon="👨‍🌾"
          isActive={getIsActive('/admin/fishermen')}
        >
          {t('fishermen')}
        </AdminSidebarItem>

        <AdminSidebarItem href="/admin/buyers" icon="👔" isActive={getIsActive('/admin/buyers')}>
          {t('buyers')}
        </AdminSidebarItem>

        <AdminSidebarItem href="/admin/venues" icon="🏢" isActive={getIsActive('/admin/venues')}>
          {t('venues')}
        </AdminSidebarItem>

        <AdminSidebarItem
          href="/admin/auctions"
          icon="📅"
          isActive={getIsActive('/admin/auctions')}
        >
          {t('auctions')}
        </AdminSidebarItem>

        <AdminSidebarItem href="/admin/items" icon="🐟" isActive={getIsActive('/admin/items')}>
          {t('items')}
        </AdminSidebarItem>

        <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

        <AdminSidebarItem href="/admin/invoice" icon="💰" isActive={getIsActive('/admin/invoice')}>
          {t('invoice')}
        </AdminSidebarItem>

        <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

        <AdminSidebarItem
          href="/admin/settings"
          icon="⚙️"
          isActive={getIsActive('/admin/settings')}
        >
          {t('settings')}
        </AdminSidebarItem>

        <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>

        <AdminSidebarItem onClick={onLogout} icon="🚪" isActive={false}>
          {t('logout') || 'ログアウト'}
        </AdminSidebarItem>
      </Stack>

      {/* Footer / User info could go here */}
      <Box
        p="4"
        bg="indigo.950"
        fontSize="xs"
        className={css({ color: 'indigo.400', textAlign: 'center' })}
      >
        &copy; FISHING AUCTION System
      </Box>
    </Box>
  );
};
