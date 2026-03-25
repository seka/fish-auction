'use client';

import { Button, HStack, Text } from '@atoms';
import { Table, Thead, Tbody, Tr, Th, Td, AuctionStatusBadge } from '@molecules';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { Auction, Venue } from '@/src/models';

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
          {auctions.map((auction) => {
            const venue = venues.find((v) => v.id === auction.venueId);
            return (
              <Tr key={auction.id}>
                <Td>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    className={css({ color: 'gray.900' })}
                  >
                    {auction.auctionDate}
                  </Text>
                  <Text fontSize="sm" className={css({ color: 'gray.500' })}>
                    {auction.startTime ? auction.startTime.substring(0, 5) : '--:--'} -{' '}
                    {auction.endTime ? auction.endTime.substring(0, 5) : '--:--'}
                  </Text>
                </Td>
                <Td>
                  <Text fontSize="sm" className={css({ color: 'gray.900' })}>
                    {venue?.name || `ID: ${auction.venueId}`}
                  </Text>
                </Td>
                <Td>
                  <AuctionStatusBadge status={auction.status} />
                </Td>
                <Td className={css({ textAlign: 'right' })}>
                  <HStack justify="end" spacing="2">
                    {auction.status === 'scheduled' && (
                      <Button
                        size="sm"
                        onClick={() => onStatusChange(auction.id, 'in_progress')}
                        disabled={isUpdatingStatus}
                        className={css({
                          color: 'green.600',
                          bg: 'green.50',
                          borderColor: 'transparent',
                          _hover: { bg: 'green.100', color: 'green.900' },
                        })}
                      >
                        {t('Admin.Auctions.start')}
                      </Button>
                    )}
                    {auction.status === 'in_progress' && (
                      <Button
                        size="sm"
                        onClick={() => onStatusChange(auction.id, 'completed')}
                        disabled={isUpdatingStatus}
                        className={css({
                          color: 'blue.600',
                          bg: 'blue.50',
                          borderColor: 'transparent',
                          _hover: { bg: 'blue.100', color: 'blue.900' },
                        })}
                      >
                        {t('Admin.Auctions.finish')}
                      </Button>
                    )}
                    <Link href={`/admin/items?auctionId=${auction.id}`}>
                      <Button
                        size="sm"
                        variant="outline"
                        className={css({
                          borderStyle: 'dashed',
                          _hover: { borderColor: 'indigo.500', color: 'indigo.600' },
                        })}
                      >
                        📦 {t('Admin.Auctions.manage_items')}
                      </Button>
                    </Link>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => onEdit(auction)}
                    >
                      {t('Common.edit')}
                    </Button>
                    <Button
                      size="sm"
                      className={css({
                        bg: 'red.50',
                        color: 'red.600',
                        _hover: { bg: 'red.100' },
                      })}
                      onClick={() => onDelete(auction.id)}
                    >
                      {t('Common.delete')}
                    </Button>
                  </HStack>
                </Td>
              </Tr>
            );
          })}
        </Tbody>
      </Table>
    </Box>
  );
};

// Internal Box component to avoid import issues in this snippet
const Box = ({ children, ...props }: any) => <div {...props}>{children}</div>;
