'use client';

import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { getMyPurchases, getMyAuctions } from '@/src/api/buyer_mypage';
import { logoutBuyer } from '@/src/api/buyer_auth';
import { AUCTION_STATUS_KEYS, AuctionStatus } from '@/src/core/assets/status';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { Box, Text, Button, Card, Stack, HStack } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { useState } from 'react';

export default function MyPage() {
    const t = useTranslations();
    const [activeTab, setActiveTab] = useState<'purchases' | 'auctions'>('purchases');
    const router = useRouter();

    // 購入履歴を取得
    // Fetch purchase history using the query hook
    const {
        data: purchases = [],
        isLoading: isPurchasesLoading
    } = useQuery({
        queryKey: ['purchases'],
        queryFn: getMyPurchases,
    });

    // 参加中のセリを取得
    const { data: auctions = [], isLoading: isAuctionsLoading } = useQuery({
        queryKey: ['auctions', 'my'],
        queryFn: getMyAuctions,
    });

    // どちらかのデータがロード中の場合はローディング表示
    const isLoading = isPurchasesLoading || isAuctionsLoading;

    const handleLogout = async () => {
        const success = await logoutBuyer();
        if (success) {
            router.push('/login/buyer');
        }
    };

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
                    </HStack>
                </Box>

                {/* Content */}
                {activeTab === 'purchases' ? (
                    <Stack spacing="4">
                        <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
                            {t('Public.MyPage.purchase_history')}
                        </Text>
                        {purchases.length === 0 ? (
                            <Box textAlign="center" py="12" bg="white" borderRadius="xl" border="1px dashed" borderColor="gray.300">
                                <Text className={css({ color: 'gray.500' })}>{t('Public.MyPage.no_history')}</Text>
                            </Box>
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
                ) : (
                    <Stack spacing="4">
                        <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
                            {t('Public.MyPage.participating_auctions')}
                        </Text>
                        {auctions.length === 0 ? (
                            <Box textAlign="center" py="12" bg="white" borderRadius="xl" border="1px dashed" borderColor="gray.300">
                                <Text className={css({ color: 'gray.500' })}>{t('Public.MyPage.no_participating')}</Text>
                            </Box>
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
                )}

            </Box >
        </Box >
    );
}
