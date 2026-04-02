'use client';

import Link from 'next/link';
import Image from 'next/image';
import { usePathname } from 'next/navigation';
import { Box, Button, HStack, Text } from '@atoms';
import { useTranslations } from 'next-intl';
import { css } from 'styled-system/css';
import { ReactNode } from 'react';

export interface PublicNavbarProps {
  isLoggedIn: boolean;
  buyerName?: string | null;
  onLogout: () => Promise<void>;
}

export const PublicNavbar = ({ isLoggedIn, buyerName, onLogout }: PublicNavbarProps) => {
  const pathname = usePathname();
  const t = useTranslations();

  if (pathname?.startsWith('/admin')) {
    return null;
  }

  return (
    <Box
      as="header"
      position="sticky"
      top="0"
      zIndex="sticky"
      w="full"
      bg="white/90"
      backdropFilter="blur(8px)"
      shadow="sm"
      borderBottom="1px solid"
      borderColor="gray.100"
    >
      <Box
        maxW="7xl"
        mx="auto"
        px={{ base: '4', md: '8' }}
        h="16"
        display="flex"
        alignItems="center"
        justifyContent="space-between"
      >
        <Link href="/" className={css({ textDecoration: 'none', _hover: { opacity: 0.8 } })}>
          <HStack spacing="0">
            <Image src="/logo_icon.png" alt="FISHING AUCTION Logo" width={50} height={50} />
            <Text
              fontSize="lg"
              fontWeight="bold"
              className={css({ color: 'indigo.900' })}
              display={{ base: 'none', sm: 'block' }}
            >
              {t('Common.app_name')}
            </Text>
          </HStack>
        </Link>

        <HStack spacing="6">
          <Box display={{ base: 'none', md: 'block' }}>
            <HStack spacing="6">
              <NavLink href="/auctions">{t('Navbar.active_auctions')}</NavLink>
              {isLoggedIn && <NavLink href="/mypage">{t('Navbar.mypage')}</NavLink>}
            </HStack>
          </Box>

          <HStack spacing="3">
            {isLoggedIn ? (
              <HStack spacing="3">
                <Text fontSize="sm" fontWeight="medium" className={css({ color: 'gray.600' })}>
                  {buyerName} {t('Navbar.honorific')}
                </Text>
                <Button size="sm" variant="outline" onClick={onLogout}>
                  {t('Navbar.logout')}
                </Button>
              </HStack>
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
const NavLink = ({ href, children }: { href: string; children: ReactNode }) => {
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
          bg: 'indigo.50',
        },
      })}
    >
      {children}
    </Link>
  );
};
