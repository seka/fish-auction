'use client';

import { Box, Card, HStack, Text, Stack, EmptyState } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface Purchase {
  id: number;
  itemId: number;
  createdAt: string;
  fishType: string;
  quantity: number;
  unit: string;
  auctionId: number;
  auctionDate: string;
  price: number;
}

interface PurchaseHistoryProps {
  purchases: Purchase[];
}

export const PurchaseHistory = ({ purchases }: PurchaseHistoryProps) => {
  const t = useTranslations();

  return (
    <Stack spacing="4">
      <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.800' })}>
        {t('Public.MyPage.purchase_history')}
      </Text>
      {purchases.length === 0 ? (
        <EmptyState
          message={t('Public.MyPage.no_history')}
          icon={
            <span role="img" aria-label="invoice">
              🧾
            </span>
          }
        />
      ) : (
        purchases.map((purchase) => (
          <Card
            key={purchase.id}
            padding="lg"
            className={css({
              _hover: { shadow: 'md' },
              transition: 'all 0.2s',
              borderWidth: '1px',
              borderColor: 'gray.200',
              bg: 'white',
            })}
          >
            <Box display="flex" justifyContent="space-between" alignItems="start">
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
                    ID: {purchase.itemId}
                  </Box>
                  <Text fontSize="xs" className={css({ color: 'gray.500' })}>
                    {new Date(purchase.createdAt).toLocaleDateString('ja-JP')}{' '}
                    {new Date(purchase.createdAt).toLocaleTimeString('ja-JP')}
                  </Text>
                </HStack>
                <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.900' })} mb="1">
                  {purchase.fishType}
                </Text>
                <Text className={css({ color: 'gray.700' })} mb="2">
                  {t('Public.MyPage.quantity')}:{' '}
                  <Text as="span" fontWeight="bold">
                    {purchase.quantity}
                  </Text>{' '}
                  {purchase.unit}
                </Text>
                <Text fontSize="sm" className={css({ color: 'gray.500' })}>
                  {t('Public.MyPage.auction_id')}: {purchase.auctionId} | {t('Public.MyPage.date')}:{' '}
                  {purchase.auctionDate}
                </Text>
              </Box>
              <Box textAlign="right">
                <Text fontSize="2xl" fontWeight="bold" className={css({ color: 'green.600' })}>
                  ¥{purchase.price.toLocaleString()}
                </Text>
              </Box>
            </Box>
          </Card>
        ))
      )}
    </Stack>
  );
};
