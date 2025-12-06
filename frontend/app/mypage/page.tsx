'use client';

import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { getMyPurchases, getMyAuctions, type Purchase, type AuctionSummary } from '@/src/api/buyer_mypage';
import { logoutBuyer } from '@/src/api/buyer_auth';
import Link from 'next/link';
import { Box, Text, Button, Card, Stack, HStack } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { useState } from 'react';

export default function MyPage() {
    const [activeTab, setActiveTab] = useState<'purchases' | 'auctions'>('purchases');
    const router = useRouter();

    // 購入履歴を取得
    const { data: purchases = [], isLoading: isPurchasesLoading } = useQuery({
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
                <Box display="flex" justifyContent="space-between" alignItems="center" mb="8">
                    <Box>
                        <Text as="h1" fontSize="3xl" fontWeight="bold" className={css({ color: 'gray.900' })}>
                            マイページ
                        </Text>
                        <Text className={css({ color: 'gray.500' })} mt="1">
                            購入履歴と参加したセリを確認できます
                        </Text>
                    </Box>
                    <HStack spacing="4">
                        <Link href="/auctions" className={css({ color: 'blue.600', _hover: { color: 'blue.700' }, fontWeight: 'medium' })}>
                            セリ一覧へ
                        </Link>
                        <Button
                            onClick={handleLogout}
                            className={css({ bg: 'gray.600', _hover: { bg: 'gray.700' }, color: 'white' })}
                        >
                            ログアウト
                        </Button>
                    </HStack>
                </Box>

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
                            購入履歴
                        </Box>
                        <Box
                            px="6"
                            py="3"
                            cursor="pointer"
                            borderBottom="2px solid"
                            borderColor={activeTab === 'auctions' ? 'blue.600' : 'transparent'}
                            color={activeTab === 'auctions' ? 'blue.600' : 'gray.500'}
                            fontWeight={activeTab === 'auctions' ? 'bold' : 'normal'}
                            onClick={() => setActiveTab('auctions')}
                            className={css({ transition: 'all 0.2s', _hover: { color: 'blue.600' } })}
                        >
                            参加中のセリ
                        </Box>
                    </HStack>
                </Box>

                {/* Content */}
                {activeTab === 'purchases' ? (
                    <Stack spacing="4">
                        <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
                            購入履歴
                        </Text>
                        {purchases.length === 0 ? (
                            <Box textAlign="center" py="12" bg="white" borderRadius="xl" border="1px dashed" borderColor="gray.300">
                                <Text className={css({ color: 'gray.500' })}>購入履歴がありません</Text>
                            </Box>
                        ) : (
                            purchases.map((purchase) => (
                                <Card
                                    key={purchase.id}
                                    p="6"
                                    borderWidth="1px"
                                    borderColor="gray.200"
                                    bg="white"
                                    className={css({ _hover: { shadow: 'md' }, transition: 'all 0.2s' })}
                                >
                                    <Box display="flex" justifyContent="space-between" alignItems="start">
                                        <Box>
                                            <HStack spacing="3" mb="2">
                                                <Box bg="blue.100" color="blue.800" fontWeight="bold" px="3" py="1" borderRadius="md" fontSize="xs">
                                                    ID: {purchase.item_id}
                                                </Box>
                                                <Text fontSize="xs" className={css({ color: 'gray.500' })}>
                                                    {new Date(purchase.created_at).toLocaleDateString('ja-JP')} {new Date(purchase.created_at).toLocaleTimeString('ja-JP')}
                                                </Text>
                                            </HStack>
                                            <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.900' })} mb="1">
                                                {purchase.fish_type}
                                            </Text>
                                            <Text className={css({ color: 'gray.700' })} mb="2">
                                                数量: <Text as="span" fontWeight="bold">{purchase.quantity}</Text> {purchase.unit}
                                            </Text>
                                            <Text fontSize="sm" className={css({ color: 'gray.500' })}>
                                                セリID: {purchase.auction_id} | 開催日: {purchase.auction_date}
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
                            参加中のセリ
                        </Text>
                        {auctions.length === 0 ? (
                            <Box textAlign="center" py="12" bg="white" borderRadius="xl" border="1px dashed" borderColor="gray.300">
                                <Text className={css({ color: 'gray.500' })}>参加中のセリがありません</Text>
                            </Box>
                        ) : (
                            auctions.map((auction) => (
                                <Card
                                    key={auction.id}
                                    p="6"
                                    borderWidth="1px"
                                    borderColor="gray.200"
                                    bg="white"
                                    className={css({ _hover: { shadow: 'md' }, transition: 'all 0.2s' })}
                                >
                                    <Box display="flex" justifyContent="space-between" alignItems="center">
                                        <Box>
                                            <HStack spacing="3" mb="2">
                                                <Box bg="blue.100" color="blue.800" fontWeight="bold" px="3" py="1" borderRadius="md" fontSize="xs">
                                                    セリ #{auction.id}
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
                                                    {auction.status === 'in_progress' ? '開催中' :
                                                        auction.status === 'completed' ? '終了' :
                                                            auction.status === 'scheduled' ? '予定' : auction.status}
                                                </Box>
                                            </HStack>
                                            <Text fontSize="lg" fontWeight="bold" className={css({ color: 'gray.900' })} mb="1">
                                                {auction.auction_date}
                                            </Text>
                                            {auction.start_time && auction.end_time && (
                                                <Text fontSize="sm" className={css({ color: 'gray.700' })}>
                                                    {auction.start_time.substring(0, 5)} - {auction.end_time.substring(0, 5)}
                                                </Text>
                                            )}
                                        </Box>
                                        <Link href={`/auctions/${auction.id}`}>
                                            <Button
                                                className={css({ bg: 'blue.600', _hover: { bg: 'blue.700' }, color: 'white' })}
                                            >
                                                詳細を見る
                                            </Button>
                                        </Link>
                                    </Box>
                                </Card>
                            ))
                        )}
                    </Stack>
                )}
            </Box>
        </Box>
    );
}
