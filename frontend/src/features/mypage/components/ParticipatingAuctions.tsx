'use client';

import { Box, Card, HStack, Text, Button, Stack, EmptyState } from '@atoms';
import { AuctionStatusBadge } from '@molecules';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useParticipatingAuctions } from '@/src/data/queries/buyerAuction/useQuery';
import { AuctionStatus } from '@entities/auction';

export const ParticipatingAuctions = () => {
  const t = useTranslations();
  const { auctions, isLoading } = useParticipatingAuctions();

  if (isLoading) {
    return <Text>{t('Common.loading')}</Text>;
  }

  return (
    <Stack spacing="4">
      <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
        {t('Public.MyPage.participating_auctions')}
      </Text>
      {auctions.length === 0 ? (
        <EmptyState
          message={t('Public.MyPage.no_participating')}
          icon={
            <span role="img" aria-label="auction">
              🏷️
            </span>
          }
          action={{
            label: t('Common.auction_list'),
            onClick: () => (window.location.href = '/auctions'),
          }}
        />
      ) : (
        auctions.map((auction) => (
          <Link key={auction.id} href={`/auctions/${auction.id}`}>
            <Card
              padding="lg"
              className={css({
                _hover: { shadow: 'md' },
                transition: 'all 0.2s',
                borderWidth: '1px',
                borderColor: 'gray.200',
                bg: 'white',
                display: 'block',
              })}
            >
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <HStack spacing="3" mb="2">
                    <Box
                      bg="blue.100"
                      color="blue.800"
                      fontWeight="bold"
                      px="3"
                      py="1"
                      borderRadius="md"
                      fontSize="xs"
                    >
                      {t('Public.MyPage.auction_id')} #{auction.id}
                    </Box>
                    <AuctionStatusBadge status={auction.status as AuctionStatus} />
                  </HStack>
                  <Text
                    fontSize="lg"
                    fontWeight="bold"
                    className={css({ color: 'gray.900' })}
                    mb="1"
                  >
                    {auction.auctionDate}
                  </Text>
                  {auction.startTime && auction.endTime && (
                    <Text fontSize="sm" className={css({ color: 'gray.700' })}>
                      {auction.startTime.substring(0, 5)} - {auction.endTime.substring(0, 5)}
                    </Text>
                  )}
                </Box>
                <Button variant="primary" size="sm">
                  {t('Public.MyPage.view_detail')}
                </Button>
              </Box>
            </Card>
          </Link>
        ))
      )}
    </Stack>
  );
};
