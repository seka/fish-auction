'use client';

import { Box, Text, Card } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { AuctionCard } from './AuctionCard';
import { Auction, Venue } from '../types';

interface AuctionsListProps {
  auctions: Auction[];
  venues: Venue[];
}

export const AuctionsList = ({ auctions, venues }: AuctionsListProps) => {
  const t = useTranslations();

  const getVenueName = (id: number) => venues?.find((v) => v.id === id)?.name || `ID: ${id}`;

  if (!auctions || auctions.length === 0) {
    return (
      <Card padding="md">
        <Box py="12" textAlign="center">
          <Text fontSize="xl" className={css({ color: 'gray.500' })}>
            {t('Public.Auctions.no_auctions')}
          </Text>
        </Box>
      </Card>
    );
  }

  return (
    <Box
      display="grid"
      gridTemplateColumns={{ base: '1fr', md: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }}
      gap="6"
    >
      {auctions.map((auction) => (
        <AuctionCard key={auction.id} auction={auction} venueName={getVenueName(auction.venueId)} />
      ))}
    </Box>
  );
};
