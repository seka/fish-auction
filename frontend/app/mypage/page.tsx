'use client';

import { AUCTION_STATUS_KEYS, AuctionStatus } from '@/src/core/assets/status';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import Link from 'next/link';
import { Box, Text, Button, Card, Stack, HStack, Input } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { EmptyState } from '../_components/atoms/EmptyState';

import { useMyPage } from './_hooks/useMyPage';

export default function MyPage() {
    const {
        t,
        activeTab,
        setActiveTab,
        currentPassword,
        setCurrentPassword,
        newPassword,
        setNewPassword,
        confirmPassword,
        setConfirmPassword,
        message,
        isPasswordUpdating,
        purchases,
        auctions,
        isLoading,
        handleLogout,
        handleUpdatePassword,
    } = useMyPage();

    if (isLoading) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
                <Text fontSize="xl" className={css({ color: 'gray.700' })}>読み込み中...</Text>
            </Box>
        );
    }

    return (
        <Box minH="screen" bg="gray.50" py="8" px="4">
            <Box maxW="7xl" mx="auto">
                {/* Header */}
                <HStack justify="between" alignItems="center" mb="8">
                    <Box>
                        <Text as="h1" fontSize="3xl" fontWeight="bold" className={css({ color: 'gray.900' })}>
                            {t(COMMON_TEXT_KEYS.mypage)}
                        </Text>
                        <Text className={css({ color: 'gray.500' })} mt="1">
                            {t(COMMON_TEXT_KEYS.mypage_description)}
                        </Text>
                    </Box>
                    <HStack spacing="4">
                        <Link href="/auctions" className={css({ color: 'blue.600', _hover: { color: 'blue.700' }, fontWeight: 'medium' })}>
                            {t(COMMON_TEXT_KEYS.auction_list)}
                        </Link>
                        <Button
                            onClick={handleLogout}
                            className={css({ bg: 'gray.600', _hover: { bg: 'gray.700' }, color: 'white' })}
                        >
                            {t(COMMON_TEXT_KEYS.logout)}
                        </Button>
                    </HStack>
                </HStack>

                {/* Tabs */}
                <Box borderBottom="1px solid" borderColor="gray.200" mb="6">
                    <HStack spacing="0">
                        <Box
                            px="6"
                            py="3"
                            cursor="pointer"
                            borderBottom="2px solid"
                            borderColor={activeTab === 'purchases' ? 'blue.600' : 'transparent'}
                            color={activeTab === 'purchases' ? 'blue.600' : 'gray.500'}
                            fontWeight={activeTab === 'purchases' ? 'bold' : 'normal'}
                            onClick={() => setActiveTab('purchases')}
                            className={css({ transition: 'all 0.2s', _hover: { color: 'blue.600' } })}
                        >
                            {t('Public.MyPage.purchase_history')}
                        </Box>
                        <Box
                            as="button"
                            px="4"
                            py="2"
                            fontWeight={activeTab === 'auctions' ? 'bold' : 'medium'}
                            color={activeTab === 'auctions' ? 'indigo.600' : 'gray.500'}
                            borderBottom={activeTab === 'auctions' ? '2px solid' : 'none'}
                            borderColor="indigo.600"
                            onClick={() => setActiveTab('auctions')}
                            className={css({ transition: 'all 0.2s', _hover: { color: 'indigo.600' } })}
                        >
                            {t('Public.MyPage.participating_auctions')}
                        </Box>
                        <Box
                            as="button"
                            px="4"
                            py="2"
                            fontWeight={activeTab === 'settings' ? 'bold' : 'medium'}
                            color={activeTab === 'settings' ? 'indigo.600' : 'gray.500'}
                            borderBottom={activeTab === 'settings' ? '2px solid' : 'none'}
                            borderColor="indigo.600"
                            onClick={() => setActiveTab('settings')}
                            className={css({ transition: 'all 0.2s', _hover: { color: 'indigo.600' } })}
                        >
                            設定
                        </Box>
                    </HStack>
                </Box>

                {/* Content */}
                {activeTab === 'purchases' ? (
                    <Stack spacing="4">
                        <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
                            {t('Public.MyPage.purchase_history')}
                        </Text>
                        {purchases.length === 0 ? (
                            <EmptyState
                                message={t('Public.MyPage.no_history')}
                                icon={<span role="img" aria-label="invoice">🧾</span>}
                            />
                        ) : (
                            purchases.map((purchase) => (
                                <Card
                                    key={purchase.id}
                                    padding="lg"
                                    className={css({ _hover: { shadow: 'md' }, transition: 'all 0.2s', borderWidth: '1px', borderColor: 'gray.200', bg: 'white' })}
                                >
                                    <Box display="flex" justifyContent="space-between" alignItems="start">
                                        <Box>
                                            <HStack spacing="3" mb="2">
                                                <Box bg="blue.100" color="blue.800" fontWeight="bold" px="3" py="1" borderRadius="md" fontSize="xs">
                                                    ID: {purchase.itemId}
                                                </Box>
                                                <Text fontSize="xs" className={css({ color: 'gray.500' })}>
                                                    {new Date(purchase.createdAt).toLocaleDateString('ja-JP')} {new Date(purchase.createdAt).toLocaleTimeString('ja-JP')}
                                                </Text>
                                            </HStack>
                                            <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.900' })} mb="1">
                                                {purchase.fishType}
                                            </Text>
                                            <Text className={css({ color: 'gray.700' })} mb="2">
                                                {t('Public.MyPage.quantity')}: <Text as="span" fontWeight="bold">{purchase.quantity}</Text> {purchase.unit}
                                            </Text>
                                            <Text fontSize="sm" className={css({ color: 'gray.500' })}>
                                                {t('Public.MyPage.auction_id')}: {purchase.auctionId} | {t('Public.MyPage.date')}: {purchase.auctionDate}
                                            </Text>
                                        </Box>
                                        <Box textAlign="right">
                                            <Text fontSize="2xl" fontWeight="bold" className={css({ color: 'green.600' })}>
                                                ¥{purchase.price.toLocaleString()}
                                            </Text>
                                        </Box>
                                    </Box>
                                </Card>
                            ))
                        )}
                    </Stack>
                ) : activeTab === 'settings' ? (
                    <Box bg="white" p="6" borderRadius="lg" shadow="sm" maxW="md">
                        <Text as="h2" fontSize="lg" fontWeight="semibold" mb="4" className={css({ color: 'gray.800' })}>
                            パスワード変更
                        </Text>

                        {message && (
                            <Box
                                p="3"
                                mb="4"
                                borderRadius="md"
                                bg={message.type === 'success' ? 'green.50' : 'red.50'}
                                color={message.type === 'success' ? 'green.700' : 'red.700'}
                                border="1px solid"
                                borderColor={message.type === 'success' ? 'green.200' : 'red.200'}
                            >
                                {message.text}
                            </Box>
                        )}

                        <form onSubmit={handleUpdatePassword}>
                            <Stack spacing="4" alignItems="stretch">
                                <Box>
                                    <Text as="label" display="block" mb="1" fontSize="sm" fontWeight="medium" className={css({ color: 'gray.700' })}>
                                        現在のパスワード
                                    </Text>
                                    <Input
                                        type="password"
                                        value={currentPassword}
                                        onChange={(e) => setCurrentPassword(e.target.value)}
                                        required
                                        className={css({ w: 'full' })}
                                    />
                                </Box>

                                <Box>
                                    <Text as="label" display="block" mb="1" fontSize="sm" fontWeight="medium" className={css({ color: 'gray.700' })}>
                                        新しいパスワード
                                    </Text>
                                    <Input
                                        type="password"
                                        value={newPassword}
                                        onChange={(e) => setNewPassword(e.target.value)}
                                        required
                                        minLength={8}
                                        className={css({ w: 'full' })}
                                    />
                                </Box>

                                <Box>
                                    <Text as="label" display="block" mb="1" fontSize="sm" fontWeight="medium" className={css({ color: 'gray.700' })}>
                                        新しいパスワード（確認）
                                    </Text>
                                    <Input
                                        type="password"
                                        value={confirmPassword}
                                        onChange={(e) => setConfirmPassword(e.target.value)}
                                        required
                                        minLength={8}
                                        className={css({ w: 'full' })}
                                    />
                                </Box>

                                <Button
                                    type="submit"
                                    disabled={isPasswordUpdating}
                                    className={css({
                                        bg: 'indigo.600',
                                        color: 'white',
                                        _hover: { bg: 'indigo.700' },
                                        _disabled: { opacity: 0.6, cursor: 'not-allowed' }
                                    })}
                                >
                                    {isPasswordUpdating ? '更新中...' : 'パスワードを変更する'}
                                </Button>
                            </Stack>
                        </form>
                    </Box>
                ) : (
                    <Stack spacing="4">
                        <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
                            {t('Public.MyPage.participating_auctions')}
                        </Text>
                        {auctions.length === 0 ? (
                            <EmptyState
                                message={t('Public.MyPage.no_participating')}
                                icon={<span role="img" aria-label="auction">🏷️</span>}
                                action={{
                                    label: t(COMMON_TEXT_KEYS.auction_list),
                                    onClick: () => window.location.href = '/auctions'
                                }}
                            />
                        ) : (
                            auctions.map((auction) => (
                                <Link key={auction.id} href={`/auctions/${auction.id}`}>
                                    <Card
                                        padding="lg"
                                        className={css({ _hover: { shadow: 'md' }, transition: 'all 0.2s', borderWidth: '1px', borderColor: 'gray.200', bg: 'white', display: 'block' })}
                                    >
                                        <Box display="flex" justifyContent="space-between" alignItems="center">
                                            <Box>
                                                <HStack spacing="3" mb="2">
                                                    <Box bg="blue.100" color="blue.800" fontWeight="bold" px="3" py="1" borderRadius="md" fontSize="xs">
                                                        {t('Public.MyPage.auction_id')} #{auction.id}
                                                    </Box>
                                                    <Box
                                                        px="3"
                                                        py="1"
                                                        borderRadius="full"
                                                        fontSize="xs"
                                                        fontWeight="bold"
                                                        bg={auction.status === 'in_progress' ? 'orange.100' : auction.status === 'completed' ? 'gray.100' : 'blue.100'}
                                                        color={auction.status === 'in_progress' ? 'orange.700' : auction.status === 'completed' ? 'gray.700' : 'blue.700'}
                                                    >
                                                        {auction.status === 'in_progress' ? t(AUCTION_STATUS_KEYS['in_progress']) : t(AUCTION_STATUS_KEYS[auction.status as AuctionStatus])}
                                                    </Box>
                                                </HStack>
                                                <Text fontSize="lg" fontWeight="bold" className={css({ color: 'gray.900' })} mb="1">
                                                    {auction.auctionDate}
                                                </Text>
                                                {auction.startTime && auction.endTime && (
                                                    <Text fontSize="sm" className={css({ color: 'gray.700' })}>
                                                        {auction.startTime.substring(0, 5)} - {auction.endTime.substring(0, 5)}
                                                    </Text>
                                                )}
                                            </Box>
                                            <Button
                                                variant="primary"
                                                size="sm"
                                            >
                                                {t('Public.MyPage.view_detail')}
                                            </Button>
                                        </Box>
                                    </Card>
                                </Link>
                            ))
                        )}
                    </Stack>
                )
                }

            </Box >
        </Box >
    );
}
