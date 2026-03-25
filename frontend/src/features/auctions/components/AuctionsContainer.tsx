'use client';

import { Box, Stack, HStack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { AuctionsList } from './AuctionsList';
import { useAuctionQuery } from '@/src/data/queries/publicAuction/useQuery';
import { useVenueQuery } from '@/src/data/queries/publicVenue/useQuery';
import { Auction } from '../types';

export const AuctionsContainer = () => {
  const t = useTranslations();
  const { auctions: allAuctions, isLoading } = useAuctionQuery();
  const { venues = [] } = useVenueQuery();

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
      (a) =>
        a.status === 'scheduled' ||
        a.status === 'in_progress' ||
        a.status === 'completed' ||
        a.status === 'cancelled',
    )
    .sort((a, b) => {
      if (a.status === 'in_progress' && b.status !== 'in_progress') return -1;
      if (a.status !== 'in_progress' && b.status === 'in_progress') return 1;
      return (
        new Date(`${a.auctionDate}T${a.startTime}`).getTime() -
        new Date(`${b.auctionDate}T${b.startTime}`).getTime()
      );
    }) as Auction[];

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
