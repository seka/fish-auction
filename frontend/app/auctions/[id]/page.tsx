'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import Link from 'next/link';

import { loginBuyer } from '@/src/api/buyer_auth';
import { AuctionItem } from '@/src/models';
import { bidSchema, BidFormData } from '@/src/models/schemas/auction';
import { buyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { useAuctionData } from './_hooks/useAuctionData';
import { useBidMutation } from './_hooks/useBidMutation';
import { useAuth } from './_hooks/useAuth';
import { isAuctionActive, formatTime } from '@/src/utils/auction';
import { AUCTION_STATUS_KEYS, ITEM_STATUS_KEYS, AuctionStatus } from '@/src/core/assets/status';
import { useTranslations } from 'next-intl';
import { Box, Text, Button, Input, Card, Stack, HStack } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function AuctionDetailPage({ params }: { params: { id: string } }) {
    const t = useTranslations();
    const router = useRouter();
    const auctionId = Number(params.id);

    const [selectedItem, setSelectedItem] = useState<AuctionItem | null>(null);
    const [message, setMessage] = useState('');
    const [loginError, setLoginError] = useState('');

    const { auction, items, isLoading, refetchItems } = useAuctionData(auctionId);
    const { submitBid, isLoading: isBidLoading } = useBidMutation();
    const { isLoggedIn, isChecking } = useAuth();

    // Check if auction is active (within bidding hours)
    const auctionActive = auction ? isAuctionActive(auction) : false;

    const { register, handleSubmit, reset, formState: { errors } } = useForm<BidFormData>({
        resolver: zodResolver(bidSchema),
    });

    const { register: registerLogin, handleSubmit: handleSubmitLogin, formState: { errors: loginErrors, isSubmitting: isLoggingIn } } = useForm<BuyerLoginFormData>({
        resolver: zodResolver(buyerLoginSchema),
    });

    // Reset selected item if it disappears from list or status changes (optional)
    useEffect(() => {
        if (selectedItem) {
            const current = items.find(i => i.id === selectedItem.id);
            if (current && current.status !== selectedItem.status) {
                setSelectedItem(current);
            }
        }
    }, [items, selectedItem]);

    if (isNaN(auctionId)) {
        return <Box>Invalid Auction ID</Box>;
    }

    if (isChecking) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
                <Text fontSize="xl" className={css({ color: 'gray.700' })}>読み込み中...</Text>
            </Box>
        );
    }

    if (isLoading) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
                <Text fontSize="xl" className={css({ color: 'gray.700' })}>読み込み中...</Text>
            </Box>
        );
    }

    if (!auction) {
        return <Box>Auction not found</Box>;
    }

    const onSubmitLogin = async (data: BuyerLoginFormData) => {
        setLoginError('');
        const buyer = await loginBuyer(data);
        if (buyer) {
            // Reload page to update auth state and stay on current page
            window.location.reload();
        } else {
            setLoginError('メールアドレスまたはパスワードが間違っています');
        }
    };

    const onSubmitBid = async (data: BidFormData) => {
        if (!selectedItem) return;

        const success = await submitBid({
            itemId: selectedItem.id,
            buyerId: 0, // Backend handles this from context
            price: parseInt(data.price),
        });

        if (success) {
            setMessage(`落札成功！ (${selectedItem.fishType})`);
            setSelectedItem(null);
            reset();
            refetchItems();
            // Clear message after 3 seconds
            setTimeout(() => setMessage(''), 3000);
        } else {
            setMessage('入札に失敗しました');
        }
    };

    // Show login form if not logged in
    if (!isLoggedIn) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" py="12" px="4">
                <Box maxW="md" w="full">
                    <Stack spacing="8">
                        <Box textAlign="center">
                            <Text as="h2" fontSize="3xl" fontWeight="extrabold" className={css({ color: 'gray.900' })}>
                                セリ会場へのログイン
                            </Text>
                            <Text mt="2" fontSize="sm" className={css({ color: 'gray.700' })}>
                                入札するにはログインが必要です
                            </Text>
                        </Box>
                        <form onSubmit={handleSubmitLogin(onSubmitLogin)}>
                            <Stack spacing="6">
                                <Stack spacing="0">
                                    <Box>
                                        <label htmlFor="email" className={css({ srOnly: true })}>メールアドレス</label>
                                        <Input
                                            id="email"
                                            type="email"
                                            {...registerLogin('email')}
                                            placeholder="メールアドレス"
                                            className={css({ borderBottomLeftRadius: '0', borderBottomRightRadius: '0' })}
                                        />
                                        {loginErrors.email && <Text className={css({ color: 'red.500' })} fontSize="xs" mt="1">{loginErrors.email.message}</Text>}
                                    </Box>
                                    <Box>
                                        <label htmlFor="password" className={css({ srOnly: true })}>パスワード</label>
                                        <Input
                                            id="password"
                                            type="password"
                                            {...registerLogin('password')}
                                            placeholder="パスワード"
                                            className={css({ borderTopLeftRadius: '0', borderTopRightRadius: '0', borderTop: 'none' })}
                                        />
                                        {loginErrors.password && <Text className={css({ color: 'red.500' })} fontSize="xs" mt="1">{loginErrors.password.message}</Text>}
                                    </Box>
                                </Stack>

                                {loginError && <Text className={css({ color: 'red.500' })} fontSize="sm" textAlign="center">{loginError}</Text>}

                                <Button
                                    type="submit"
                                    disabled={isLoggingIn}
                                    w="full"
                                    className={css({ bg: 'indigo.600', _hover: { bg: 'indigo.700' }, color: 'white' })}
                                >
                                    ログイン
                                </Button>
                                <Box textAlign="center">
                                    <Link href="/signup" className={css({ fontSize: 'sm', color: 'indigo.600', _hover: { color: 'indigo.500' } })}>
                                        アカウントをお持ちでない方はこちら
                                    </Link>
                                </Box>
                            </Stack>
                        </form>
                    </Stack>
                </Box>
            </Box>
        );
    }

    return (
        <Box minH="screen" bg="gray.50" p={{ base: '4', md: '8' }}>
            <Box maxW="7xl" mx="auto">
                {/* Header */}
                <Box display="flex" flexDirection={{ base: 'column', md: 'row' }} justifyContent="space-between" alignItems={{ base: 'start', md: 'center' }} mb="8" gap="4">
                    <Box>
                        <HStack spacing="3" mb="1">
                            <Link href="/auctions" className={css({ color: 'gray.500', _hover: { color: 'gray.700' } })}>
                                &larr; 一覧へ
                            </Link>
                            <Box
                                px="3"
                                py="1"
                                borderRadius="full"
                                fontSize="sm"
                                fontWeight="bold"
                                bg={auction.status === 'in_progress' ? 'orange.100' : 'blue.100'}
                                color={auction.status === 'in_progress' ? 'orange.700' : 'blue.700'}
                                className={auction.status === 'in_progress' ? css({ animation: 'pulse 2s infinite' }) : ''}
                            >
                                {auction.status === 'in_progress' ? '🔥 ' + t(AUCTION_STATUS_KEYS['in_progress']) : t(AUCTION_STATUS_KEYS[auction.status as AuctionStatus])}
                            </Box>
                        </HStack>
                        <Text as="h1" fontSize="3xl" fontWeight="bold" className={css({ color: 'gray.900' })}>
                            セリ会場 #{auction.id}
                        </Text>
                        <Text className={css({ color: 'gray.600' })}>
                            {auction.auctionDate} {auction.startTime?.substring(0, 5)} - {auction.endTime?.substring(0, 5)}
                        </Text>
                    </Box>
                    <Box textAlign="right" display={{ base: 'none', md: 'block' }}>
                        <Text fontSize="sm" className={css({ color: 'gray.600' })}>自動更新中 (5秒)</Text>
                    </Box>
                </Box>

                {message && (
                    <Card
                        mb="6"
                        p="4"
                        borderLeft="4px solid"
                        borderColor="green.500"
                        bg="green.50"
                        className={css({ animation: 'bounce 1s infinite' })}
                    >
                        <Text fontWeight="bold" className={css({ color: 'green.700' })}>{message}</Text>
                    </Card>
                )}

                <Box display="grid" gridTemplateColumns={{ base: 'repeat(1, 1fr)', lg: 'repeat(3, 1fr)' }} gap="8">
                    {/* Item List */}
                    <Box gridColumn={{ base: '1', lg: 'span 2' }}>
                        <Stack spacing="4">
                            <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })} borderBottom="1px solid" borderColor="gray.200" pb="2">
                                出品リスト
                            </Text>
                            {items.length === 0 ? (
                                <Box textAlign="center" py="12" bg="white" borderRadius="xl" border="1px dashed" borderColor="gray.300">
                                    <Text className={css({ color: 'gray.600' })}>現在、出品されている商品はありません。</Text>
                                </Box>
                            ) : (
                                items.map((item) => (
                                    <Card
                                        key={item.id}
                                        p="6"
                                        borderWidth="2px"
                                        borderColor={selectedItem?.id === item.id ? 'orange.500' : 'gray.200'}
                                        bg={selectedItem?.id === item.id ? 'orange.50' : 'white'}
                                        cursor="pointer"
                                        transition="all 0.2s"
                                        className={selectedItem?.id === item.id ? css({ shadow: 'md', transform: 'scale(1.01)' }) : css({ _hover: { shadow: 'md' } })}
                                        onClick={() => setSelectedItem(item)}
                                    >
                                        <Box display="flex" justifyContent="space-between" alignItems="center">
                                            <HStack spacing="4">
                                                <Box bg="blue.100" color="blue.800" fontWeight="bold" px="3" py="1" borderRadius="md" fontSize="xs">
                                                    ID: {item.id}
                                                </Box>
                                                <Box>
                                                    <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.900' })}>{item.fishType}</Text>
                                                    <Text className={css({ color: 'gray.700' })} mt="1">
                                                        <Text as="span" fontWeight="bold" fontSize="lg">{item.quantity}</Text> {item.unit}
                                                        <Text as="span" fontSize="sm" ml="2" className={css({ color: 'gray.500' })}>(漁師ID: {item.fishermanId})</Text>
                                                    </Text>
                                                    {item.highestBid && (
                                                        <Text fontSize="sm" mt="1" className={css({ color: 'orange.600' })} fontWeight="semibold">
                                                            現在の最高額: ¥{item.highestBid.toLocaleString()}
                                                            {item.highestBidderName && (
                                                                <Text as="span" ml="2" className={css({ color: 'gray.700' })}>({item.highestBidderName} さん)</Text>
                                                            )}
                                                        </Text>
                                                    )}
                                                </Box>
                                            </HStack>
                                            <Box
                                                px="4"
                                                py="2"
                                                borderRadius="full"
                                                fontSize="sm"
                                                fontWeight="bold"
                                                bg={item.status === 'Pending' ? 'green.100' : 'gray.100'}
                                                color={item.status === 'Pending' ? 'green.800' : 'gray.600'}
                                                shadow="sm"
                                            >
                                                {item.status === 'Pending' ? t(ITEM_STATUS_KEYS['Pending']) : (ITEM_STATUS_KEYS[item.status] ? t(ITEM_STATUS_KEYS[item.status]) : item.status)}
                                            </Box>
                                        </Box>
                                    </Card>
                                ))
                            )}
                        </Stack>
                    </Box>

                    {/* Bidding Panel */}
                    <Box gridColumn={{ base: '1', lg: 'span 1' }}>
                        <Card p="6" shadow="lg" borderWidth="1px" borderColor="gray.200" position={{ lg: 'sticky' }} top="6">
                            <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })} borderBottom="1px solid" borderColor="gray.200" pb="2" mb="6">
                                入札パネル
                            </Text>
                            {selectedItem ? (
                                <form onSubmit={handleSubmit(onSubmitBid)}>
                                    <Stack spacing="6">
                                        <Box p="5" bg="gray.50" borderRadius="lg" borderWidth="1px" borderColor="gray.200">
                                            <Text fontSize="sm" className={css({ color: 'gray.600' })} mb="1">選択中の商品</Text>
                                            <Text fontWeight="bold" fontSize="2xl" className={css({ color: 'gray.900' })}>{selectedItem.fishType}</Text>
                                            <Text fontSize="lg" className={css({ color: 'gray.700' })}>{selectedItem.quantity} {selectedItem.unit}</Text>
                                            <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="2">ステータス: {ITEM_STATUS_KEYS[selectedItem.status] ? t(ITEM_STATUS_KEYS[selectedItem.status]) : selectedItem.status}</Text>
                                            {selectedItem.highestBid && (
                                                <Text fontSize="sm" mt="2" className={css({ color: 'orange.600' })} fontWeight="bold">
                                                    現在の最高額: ¥{selectedItem.highestBid.toLocaleString()}
                                                    {selectedItem.highestBidderName && (
                                                        <Text as="span" ml="2" className={css({ color: 'gray.700' })}>({selectedItem.highestBidderName} さん)</Text>
                                                    )}
                                                </Text>
                                            )}
                                        </Box>

                                        {selectedItem.status === 'Pending' ? (
                                            !auctionActive ? (
                                                <Box textAlign="center" py="6" bg="yellow.50" borderRadius="lg" borderWidth="1px" borderColor="yellow.200">
                                                    <Text className={css({ color: 'yellow.800' })} fontWeight="bold" mb="2">⏰ 入札受付時間外</Text>
                                                    {auction.startTime && auction.endTime && (
                                                        <Text fontSize="sm" className={css({ color: 'yellow.700' })}>
                                                            受付時間: {formatTime(auction.startTime)} ~ {formatTime(auction.endTime)}
                                                        </Text>
                                                    )}
                                                </Box>
                                            ) : (
                                                <>
                                                    <Box>
                                                        <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                                            入札価格 (円)
                                                        </Text>
                                                        <Box position="relative">
                                                            <Box position="absolute" top="50%" left="3" transform="translateY(-50%)" pointerEvents="none">
                                                                <Text fontSize="sm" className={css({ color: 'gray.600' })}>¥</Text>
                                                            </Box>
                                                            <Input
                                                                type="number"
                                                                {...register('price')}
                                                                placeholder="0"
                                                                className={css({ pl: '7' })}
                                                            />
                                                        </Box>
                                                        {errors.price && (
                                                            <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.price.message}</Text>
                                                        )}
                                                    </Box>

                                                    <Button
                                                        type="submit"
                                                        disabled={isBidLoading}
                                                        w="full"
                                                        size="lg"
                                                        className={css({
                                                            bg: 'red.600',
                                                            _hover: { bg: 'red.700', transform: 'scale(1.02)' },
                                                            color: 'white',
                                                            shadow: 'md',
                                                            transition: 'all 0.2s'
                                                        })}
                                                    >
                                                        {isBidLoading ? '処理中...' : '落札する！'}
                                                    </Button>
                                                </>
                                            )
                                        ) : (
                                            <Box textAlign="center" py="4" bg="gray.100" borderRadius="md" color="gray.500">
                                                この商品は既に入札が終了しています
                                            </Box>
                                        )}
                                    </Stack>
                                </form>
                            ) : (
                                <Box textAlign="center" py="12" color="gray.400">
                                    <Text>左のリストから<br />商品を選択してください</Text>
                                </Box>
                            )}
                        </Card>
                    </Box>
                </Box>
            </Box>
        </Box>
    );
}
