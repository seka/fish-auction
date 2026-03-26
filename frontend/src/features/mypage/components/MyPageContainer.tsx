'use client';

import { Box, HStack, Text, Button } from '@atoms';
import { css } from 'styled-system/css';
import Link from 'next/link';
import { PurchaseHistory } from './PurchaseHistory';
import { ParticipatingAuctions } from './ParticipatingAuctions';
import { SettingsForm } from './SettingsForm';

import { useMyPage } from '../states/useMyPage';
import { MyPageTab } from '../types';

export const MyPageContainer = () => {
  const {
    t,
    activeTab,
    setActiveTab,
    handleLogout,
    passwordState,
  } = useMyPage();

  const {
    currentPassword,
    setCurrentPassword,
    newPassword,
    setNewPassword,
    confirmPassword,
    setConfirmPassword,
    passwordMessage,
    handleUpdatePassword,
    isPasswordUpdating,
  } = passwordState;

  return (
    <Box minH="screen" bg="gray.50" py="8" px="4">
      <Box maxW="7xl" mx="auto">
        {/* Header */}
        <HStack justify="between" alignItems="center" mb="8">
          <Box>
            <Text as="h1" fontSize="3xl" fontWeight="bold" className={css({ color: 'gray.900' })}>
              {t('Common.mypage')}
            </Text>
            <Text className={css({ color: 'gray.500' })} mt="1">
              {t('Common.mypage_description')}
            </Text>
          </Box>
          <HStack spacing="4">
            <Link
              href="/auctions"
              className={css({
                color: 'blue.600',
                _hover: { color: 'blue.700' },
                fontWeight: 'medium',
              })}
            >
              {t('Common.auction_list')}
            </Link>
            <Button
              onClick={handleLogout}
              className={css({ bg: 'gray.600', _hover: { bg: 'gray.700' }, color: 'white' })}
            >
              {t('Common.logout')}
            </Button>
          </HStack>
        </HStack>

        {/* Tabs */}
        <Box borderBottom="1px solid" borderColor="gray.200" mb="6">
          <HStack spacing="0">
            {[
              { id: 'purchases', label: t('Public.MyPage.purchase_history') },
              { id: 'auctions', label: t('Public.MyPage.participating_auctions') },
              { id: 'settings', label: t('Public.MyPage.settings') },
            ].map((tab) => (
              <Box
                key={tab.id}
                as="button"
                px="6"
                py="3"
                cursor="pointer"
                borderBottom="2px solid"
                borderColor={activeTab === tab.id ? 'blue.600' : 'transparent'}
                color={activeTab === tab.id ? 'blue.600' : 'gray.500'}
                fontWeight={activeTab === tab.id ? 'bold' : 'normal'}
                onClick={() => setActiveTab(tab.id as MyPageTab)}
                className={css({ transition: 'all 0.2s', _hover: { color: 'blue.600' } })}
              >
                {tab.label}
              </Box>
            ))}
          </HStack>
        </Box>

        {/* Content */}
        <Box>
          {activeTab === 'purchases' && <PurchaseHistory />}
          {activeTab === 'auctions' && <ParticipatingAuctions />}
          {activeTab === 'settings' && (
            <SettingsForm
              currentPassword={currentPassword}
              setCurrentPassword={setCurrentPassword}
              newPassword={newPassword}
              setNewPassword={setNewPassword}
              confirmPassword={confirmPassword}
              setConfirmPassword={setConfirmPassword}
              passwordMessage={passwordMessage}
              handleUpdatePassword={handleUpdatePassword}
              isPasswordUpdating={isPasswordUpdating}
            />
          )}
        </Box>
      </Box>
    </Box>
  );
};
