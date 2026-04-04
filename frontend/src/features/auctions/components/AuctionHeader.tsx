'use client';

import Link from 'next/link';
import { Box, Text, HStack } from '@atoms';
import { AuctionStatusBadge } from './AuctionStatusBadge';
import { css } from 'styled-system/css';
import { Auction } from '../types';

import { useTranslations } from 'next-intl';

interface AuctionHeaderProps {
  auction: Auction;
  t: ReturnType<typeof useTranslations>;
}

export const AuctionHeader = ({ auction, t }: AuctionHeaderProps) => {
  return (
    <Box
      display="flex"
      flexDirection={{ base: 'column', md: 'row' }}
      justifyContent="space-between"
      alignItems={{ base: 'start', md: 'center' }}
      mb="8"
      gap="4"
    >
      <Box>
        <HStack spacing="3" mb="1">
          <Link
            href="/auctions"
            className={css({ color: 'gray.500', _hover: { color: 'gray.700' } })}
          >
            &larr; {t('Public.AuctionDetail.back_to_list')}
          </Link>
          <AuctionStatusBadge status={auction.status} />
        </HStack>
        <Text as="h1" fontSize="3xl" fontWeight="bold" className={css({ color: 'gray.900' })}>
          {t('Public.AuctionDetail.auction_venue_title', { id: auction.id })}
        </Text>
        <Text className={css({ color: 'gray.600' })}>
          {auction.duration?.dateLabel} {auction.duration?.label}
        </Text>
      </Box>
      <Box textAlign="right" display={{ base: 'none', md: 'block' }}>
        <Text fontSize="sm" className={css({ color: 'gray.600' })}>
          {t('Public.AuctionDetail.auto_refresh')}
        </Text>
      </Box>
    </Box>
  );
};
