'use client';

import { Box, Stack, HStack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { AuctionsList } from '@/src/features/auctions/components';
import { usePublicAuctions, useVenues } from '@/src/features/auctions/queries';
import { Auction } from '@/src/features/auctions/types';

export const AuctionsContainer = () => {
  const t = useTranslations();
  const { data: allAuctions, isLoading } = usePublicAuctions();
  const { venues = [] } = useVenues();

  if (isLoading) {
    return (
      <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
        <Text fontSize="xl" color="muted">
          {t('Common.loading')}
        </Text>
      </Box>
    );
  }

  // Filter and Sort logic
  const auctions = (allAuctions || [])
    .filter(
      (a: Auction) =>
        a.status.value === 'scheduled' ||
        a.status.value === 'in_progress' ||
        a.status.value === 'completed' ||
        a.status.value === 'cancelled',
    )
    .sort((a: Auction, b: Auction) => {
      if (a.status.value === 'in_progress' && b.status.value !== 'in_progress') return -1;
      if (a.status.value !== 'in_progress' && b.status.value === 'in_progress') return 1;

      const aTime = a.duration.startTime
        ? `${a.duration.dateLabel}T${a.duration.startTime}`
        : `${a.duration.dateLabel}T00:00:00`;
      const bTime = b.duration.startTime
        ? `${b.duration.dateLabel}T${b.duration.startTime}`
        : `${b.duration.dateLabel}T00:00:00`;

      return new Date(aTime).getTime() - new Date(bTime).getTime();
    });

  return (
    <Box maxW="7xl" mx="auto" px={{ base: '4', md: '8' }} py="8">
      <Stack spacing="8">
        {/* Header */}
        <Box>
          <Text as="h1" variant="h2" color="default">
            {t('Public.Auctions.title')}
          </Text>
          <HStack mt="2" className={css({ fontSize: 'sm', color: 'gray.500' })}>
            <Text className={css({ color: 'gray.300' })} mx="2">
              /
            </Text>
            <Text>{t('Public.Auctions.title')}</Text>
          </HStack>
        </Box>

        {/* Back to Top */}
        <Box>
          <Link
            href="/"
            className={css({
              display: 'inline-flex',
              alignItems: 'center',
              color: 'indigo.600',
              fontWeight: 'medium',
              _hover: { textDecoration: 'underline' },
            })}
          >
            &larr; {t('Common.back_to_top')}
          </Link>
        </Box>

        <AuctionsList auctions={auctions} venues={venues} />
      </Stack>
    </Box>
  );
};
