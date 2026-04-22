'use client';

import Link from 'next/link';
import { Box, HStack, Text, Card, Stack } from '@atoms';
import { AuctionStatusBadge } from './AuctionStatusBadge';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { Auction } from '../types';

interface AuctionCardProps {
  auction: Auction;
  venueName: string;
}

export const AuctionCard = ({ auction, venueName }: AuctionCardProps) => {
  const t = useTranslations();

  return (
    <Link
      href={`/auctions/${auction.id}`}
      className={css({
        textDecoration: 'none',
        display: 'block',
        transition: 'transform 0.2s',
        _hover: { transform: 'translateY(-4px)' },
      })}
    >
      <Card
        padding="none"
        overflow="hidden"
        className={css({
          h: 'full',
          border: '1px solid',
          borderColor: 'gray.200',
          _hover: { shadow: 'md', borderColor: 'indigo.200' },
        })}
      >
        <Box bg="indigo.600" h="1" w="full" />

        <Box p="6">
          <HStack justify="between" mb="4">
            <AuctionStatusBadge status={auction.status} />
            <Text fontSize="sm" className={css({ color: 'gray.500' })}>
              {auction.duration.startAt
                ? auction.duration.startAt.toLocaleDateString('sv-SE', { timeZone: 'Asia/Tokyo' })
                : '-'}
            </Text>
          </HStack>

          <Text
            as="h3"
            fontSize="xl"
            fontWeight="bold"
            className={css({ color: 'gray.900', mb: '2', lineClamp: 1 })}
          >
            {venueName}
          </Text>

          <Stack spacing="2" mt="4">
            <HStack className={css({ fontSize: 'sm', color: 'gray.600' })}>
              <span className={css({ w: '5', textAlign: 'center' })}>⏰</span>
              <Text>{auction.duration.label}</Text>
            </HStack>
          </Stack>
        </Box>

        <Box
          px="6"
          py="4"
          bg="gray.50"
          borderTop="1px solid"
          borderColor="gray.100"
          display="flex"
          justifyContent="flex-end"
        >
          <Text
            className={css({
              color: 'indigo.600',
              fontWeight: 'bold',
              fontSize: 'sm',
              display: 'flex',
              alignItems: 'center',
            })}
          >
            {t('Public.Auctions.enter_venue')} <span className={css({ ml: '1' })}>&rarr;</span>
          </Text>
        </Box>
      </Card>
    </Link>
  );
};
