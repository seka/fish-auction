'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { getAuctions } from '@/src/api/auction';
import { AUCTION_STATUS_KEYS, AuctionStatus } from '@/src/core/assets/status';
import { useTranslations } from 'next-intl';
import { Box, Stack, HStack, Text, Card } from '@/src/core/ui';
import { css } from 'styled-system/css';

import { usePublicVenues } from './_hooks/usePublicVenues';

export default function AuctionsListPage() {
    const t = useTranslations();
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
    const auctions = allAuctions?.filter(a =>
        a.status === 'scheduled' || a.status === 'in_progress' || a.status === 'completed' // Show completed too for public view? Or just active? Let's show all for now or stick to active. Previous code had active.
    ).sort((a, b) => {
        // Sort: In Progress first, then by date/time
        if (a.status === 'in_progress' && b.status !== 'in_progress') return -1;
        if (a.status !== 'in_progress' && b.status === 'in_progress') return 1;
        return new Date(`${a.auctionDate}T${a.startTime}`).getTime() - new Date(`${b.auctionDate}T${b.startTime}`).getTime();
    }) || [];

    const getVenueName = (id: number) => venues?.find(v => v.id === id)?.name || `ID: ${id}`;

    // Helper functions for status styles
    const getStatusColor = (status: string) => {
        switch (status) {
            case 'in_progress': return 'orange.500';
            case 'scheduled': return 'blue.500';
            case 'completed': return 'gray.500';
            case 'cancelled': return 'red.500';
            default: return 'gray.200';
        }
    };

    const getStatusBg = (status: string) => {
        switch (status) {
            case 'in_progress': return 'orange.100';
            case 'scheduled': return 'blue.100';
            case 'completed': return 'gray.100';
            case 'cancelled': return 'red.100';
            default: return 'gray.100';
        }
    };

    const getStatusTextColor = (status: string) => {
        switch (status) {
            case 'in_progress': return 'orange.800';
            case 'scheduled': return 'blue.800';
            case 'completed': return 'gray.800';
            case 'cancelled': return 'red.800';
            default: return 'gray.800';
        }
    };

    return (
        <Box maxW="7xl" mx="auto" px={{ base: '4', md: '8' }} py="8">
            <Stack spacing="8">
                {/* Header */}
                <Box>
                    <Text as="h1" variant="h2" color="default">{t('Public.Auctions.title')}</Text>
                    <HStack mt="2" className={css({ fontSize: 'sm', color: 'gray.500' })}>
                        <Link href="/" className={css({ _hover: { color: 'indigo.600', textDecoration: 'underline' } })}>{t('Common.home_title')}</Link>
                        <Text color="gray.300" mx="2">/</Text>
                        <Text>{t('Public.Auctions.title')}</Text>
                    </HStack>
                </Box>

                {/* Back to Top */}
                <Box>
                    <Link href="/" className={css({ display: 'inline-flex', alignItems: 'center', color: 'indigo.600', fontWeight: 'medium', _hover: { textDecoration: 'underline' } })}>
                        &larr; {t('Common.back_to_top')}
                    </Link>
                </Box>

                {/* Auctions Grid */}
                {(!auctions || auctions.length === 0) ? (
                    <Card padding="md">
                        <Box py="12" textAlign="center">
                            <Text fontSize="xl" className={css({ color: 'gray.500' })}>{t('Public.Auctions.no_auctions')}</Text>
                        </Box>
                    </Card>
                ) : (
                    <Box display="grid" gridTemplateColumns={{ base: '1fr', md: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }} gap="6">
                        {auctions.map((auction) => (
                            <Link key={auction.id} href={`/auctions/${auction.id}`} className={css({ textDecoration: 'none', display: 'block', transition: 'transform 0.2s', _hover: { transform: 'translateY(-4px)' } })}>
                                <Card padding="none" overflow="hidden" className={css({ h: 'full', border: '1px solid', borderColor: 'gray.200', _hover: { shadow: 'md', borderColor: 'indigo.200' } })}>
                                    {/* Status Badge Strip */}
                                    <Box bg={getStatusColor(auction.status)} h="2" w="full" />

                                    <Box p="6">
                                        <HStack justify="between" mb="4">
                                            <span className={css({
                                                fontSize: 'xs',
                                                fontWeight: 'bold',
                                                px: '2.5',
                                                py: '1',
                                                borderRadius: 'full',
                                                bg: getStatusBg(auction.status),
                                                color: getStatusTextColor(auction.status)
                                            })}>
                                                {auction.status === 'in_progress' ? '🔥 ' + t(AUCTION_STATUS_KEYS['in_progress']) : t(AUCTION_STATUS_KEYS[auction.status as AuctionStatus] || 'scheduled')}
                                            </span>
                                            <Text fontSize="sm" className={css({ color: 'gray.500' })}>{auction.auctionDate}</Text>
                                        </HStack>

                                        <Text as="h3" fontSize="xl" fontWeight="bold" className={css({ color: 'gray.900', mb: '2', lineClamp: 1 })}>
                                            {getVenueName(auction.venueId)}
                                        </Text>

                                        <Stack spacing="2" mt="4">
                                            <HStack className={css({ fontSize: 'sm', color: 'gray.600' })}>
                                                <span className={css({ w: '5', textAlign: 'center' })}>⏰</span>
                                                <Text>
                                                    {auction.startTime ? auction.startTime.substring(0, 5) : '--:--'} - {auction.endTime ? auction.endTime.substring(0, 5) : '--:--'}
                                                </Text>
                                            </HStack>
                                            {auction.endTime && (
                                                <HStack className={css({ fontSize: 'sm', color: 'gray.600' })}>
                                                    <span className={css({ w: '5', textAlign: 'center' })}>🏁</span>
                                                    <Text>{t('Public.Auctions.end_time_prefix')} {auction.endTime?.substring(0, 5)}</Text>
                                                </HStack>
                                            )}
                                        </Stack>
                                    </Box>

                                    <Box px="6" py="4" bg="gray.50" borderTop="1px solid" borderColor="gray.100" display="flex" justifyContent="flex-end">
                                        <Text className={css({ color: 'indigo.600', fontWeight: 'bold', fontSize: 'sm', display: 'flex', alignItems: 'center' })}>
                                            {t('Public.Auctions.enter_venue')} <span className={css({ ml: '1' })}>&rarr;</span>
                                        </Text>
                                    </Box>
                                </Card>
                            </Link>
                        ))}
                    </Box>
                )}
            </Stack>
        </Box>
    );
}
