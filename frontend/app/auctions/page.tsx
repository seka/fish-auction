'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { getAuctions } from '@/src/api/auction';
import { getVenues } from '@/src/api/venue';
import { Box, Stack, HStack, Text, Card } from '@/src/core/ui';
import { css } from 'styled-system/css';

const usePublicVenues = () => {
    const { data: venues } = useQuery({
        queryKey: ['public_venues'],
        queryFn: getVenues,
    });
    return { venues };
};

export default function AuctionsListPage() {
    // Fetch all auctions
    const { data: allAuctions, isLoading } = useQuery({
        queryKey: ['public_auctions_list'],
        queryFn: () => getAuctions(),
    });

    const { venues } = usePublicVenues();

    if (isLoading) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
                <Text fontSize="xl" color="muted">読み込み中...</Text>
            </Box>
        );
    }

    // Filter for active auctions (Scheduled or In Progress)
    const activeAuctions = allAuctions?.filter(a =>
        a.status === 'scheduled' || a.status === 'in_progress'
    ).sort((a, b) => {
        // Sort: In Progress first, then by date/time
        if (a.status === 'in_progress' && b.status !== 'in_progress') return -1;
        if (a.status !== 'in_progress' && b.status === 'in_progress') return 1;
        return new Date(`${a.auctionDate}T${a.startTime}`).getTime() - new Date(`${b.auctionDate}T${b.startTime}`).getTime();
    }) || [];

    const getVenueName = (id: number) => venues?.find(v => v.id === id)?.name || `会場ID: ${id}`;

    return (
        <Box minH="screen" bg="gray.50" p="8">
            <Box maxW="5xl" mx="auto">
                <HStack justify="between" mb="8">
                    <Text as="h1" variant="h2" color="default">開催中のセリ一覧</Text>
                    <HStack gap="4">
                        <Link href="/mypage" className={css({ color: 'indigo.600', _hover: { color: 'indigo.800' }, fontWeight: 'medium' })}>
                            マイページ
                        </Link>
                        <Link href="/" className={css({ color: 'indigo.600', _hover: { color: 'indigo.800' }, fontWeight: 'medium' })}>
                            &larr; トップに戻る
                        </Link>
                    </HStack>
                </HStack>

                {activeAuctions.length === 0 ? (
                    <Card padding="lg" className={css({ textAlign: 'center' })}>
                        <Text fontSize="xl" className={css({ color: 'gray.500' })}>現在開催予定のセリはありません</Text>
                    </Card>
                ) : (
                    <Box display="grid" gridTemplateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap="6">
                        {activeAuctions.map((auction) => (
                            <Link
                                key={auction.id}
                                href={`/auctions/${auction.id}`}
                                className={css({ display: 'block', _hover: { textDecoration: 'none' } })}
                            >
                                <Card
                                    variant="interactive"
                                    className={css({
                                        height: 'full',
                                        transition: 'all 0.2s',
                                        borderColor: auction.status === 'in_progress' ? 'orange.400' : 'transparent',
                                        borderWidth: auction.status === 'in_progress' ? '2px' : '1px',
                                        ring: auction.status === 'in_progress' ? '4px' : '0',
                                        ringColor: auction.status === 'in_progress' ? 'orange.50' : 'transparent',
                                        _hover: {
                                            borderColor: auction.status === 'in_progress' ? 'orange.500' : 'indigo.200',
                                        }
                                    })}
                                >
                                    <HStack justify="between" align="start" mb="4">
                                        <Box>
                                            <Box
                                                display="inline-block"
                                                px="3"
                                                py="1"
                                                borderRadius="full"
                                                fontSize="sm"
                                                fontWeight="bold"
                                                mb="2"
                                                bg={auction.status === 'in_progress' ? 'orange.100' : 'blue.100'}
                                                color={auction.status === 'in_progress' ? 'orange.700' : 'blue.700'}
                                                animation={auction.status === 'in_progress' ? 'pulse 2s infinite' : 'none'}
                                            >
                                                {auction.status === 'in_progress' ? '🔥 開催中' : '📅 開催予定'}
                                            </Box>
                                            <Text variant="h3" color="default" className={css({ _groupHover: { color: 'indigo.700' }, transition: 'colors' })}>
                                                {getVenueName(auction.venueId)}
                                            </Text>
                                        </Box>
                                        <Box textAlign="right">
                                            <Text fontSize="2xl" fontWeight="bold" color="default">
                                                {auction.startTime?.substring(0, 5)}
                                            </Text>
                                            <Text fontSize="sm" className={css({ color: 'gray.500' })}>
                                                {auction.auctionDate}
                                            </Text>
                                        </Box>
                                    </HStack>

                                    <HStack justify="between" mt="4" pt="4" borderTop="1px solid" borderColor="gray.100">
                                        <Text fontSize="sm" color="muted">
                                            終了予定: {auction.endTime?.substring(0, 5)}
                                        </Text>
                                        <Text className={css({ color: 'indigo.600', fontWeight: 'bold', display: 'flex', alignItems: 'center', _groupHover: { transform: 'translateX(4px)' }, transition: 'transform' })}>
                                            会場へ入る <span className={css({ ml: '1' })}>&rarr;</span>
                                        </Text>
                                    </HStack>
                                </Card>
                            </Link>
                        ))}
                    </Box>
                )}
            </Box>
        </Box>
    );
}
