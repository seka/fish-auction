'use client';

import { Box } from '@atoms';
import { Table, Thead, Tbody, Tr, Th } from '@molecules';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { Auction, Venue } from '@/src/models';
import { AuctionListCell } from './AuctionListCell';

interface AuctionListProps {
  auctions: Auction[];
  venues: Venue[];
  isLoading: boolean;
  isUpdatingStatus: boolean;
  onEdit: (auction: Auction) => void;
  onDelete: (id: number) => void;
  onStatusChange: (id: number, status: string) => void;
}

export const AuctionList = ({
  auctions,
  venues,
  isLoading,
  isUpdatingStatus,
  onEdit,
  onDelete,
  onStatusChange,
}: AuctionListProps) => {
  const t = useTranslations();

  if (isLoading) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  return (
    <Box overflowX="auto">
      <Table>
        <Thead>
          <Tr>
            <Th>{t('Admin.Auctions.date_time')}</Th>
            <Th>{t('Admin.Auctions.venue')}</Th>
            <Th>{t('Admin.Auctions.status')}</Th>
            <Th className={css({ textAlign: 'right' })}>{t('Admin.Auctions.action')}</Th>
          </Tr>
        </Thead>
        <Tbody>
          {auctions.map((auction) => (
            <AuctionListCell
              key={auction.id}
              auction={auction}
              venue={venues.find((v) => v.id === auction.venueId)}
              isUpdatingStatus={isUpdatingStatus}
              onEdit={onEdit}
              onDelete={onDelete}
              onStatusChange={onStatusChange}
              t={t}
            />
          ))}
        </Tbody>
      </Table>
    </Box>
  );
};
