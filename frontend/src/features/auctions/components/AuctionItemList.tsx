'use client';

import { Box, Text, Stack, Card, HStack } from '@atoms';
import { ItemStatusBadge } from '@molecules';
import { css } from 'styled-system/css';
import { AuctionItem } from '@/src/models';

import { useTranslations } from 'next-intl';

interface AuctionItemListProps {
  items: AuctionItem[];
  selectedItemId: number | null;
  onSelectItem: (id: number) => void;
  t: ReturnType<typeof useTranslations>;
}

export const AuctionItemList = ({
  items,
  selectedItemId,
  onSelectItem,
  t,
}: AuctionItemListProps) => {
  return (
    <Box gridColumn={{ base: '1', lg: 'span 2' }}>
      <Stack spacing="4">
        <Text
          fontSize="xl"
          fontWeight="bold"
          className={css({ color: 'gray.800' })}
          borderBottom="1px solid"
          borderColor="gray.200"
          pb="2"
        >
          {t('Public.AuctionDetail.item_list')}
        </Text>
        {items.length === 0 ? (
          <Box
            textAlign="center"
            py="12"
            bg="white"
            borderRadius="xl"
            border="1px dashed"
            borderColor="gray.300"
          >
            <Text className={css({ color: 'gray.600' })}>{t('Public.AuctionDetail.no_items')}</Text>
          </Box>
        ) : (
          items.map((item: AuctionItem) => (
            <Card
              key={item.id}
              padding="lg"
              className={css({
                borderWidth: '2px',
                borderColor: selectedItemId === item.id ? 'orange.500' : 'gray.200',
                bg: selectedItemId === item.id ? 'orange.50' : 'white',
                cursor: 'pointer',
                transition: 'all 0.2s',
                shadow: selectedItemId === item.id ? 'md' : 'none',
                transform: selectedItemId === item.id ? 'scale(1.01)' : 'none',
                _hover: { shadow: 'md' },
              })}
              onClick={() => onSelectItem(item.id)}
            >
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <HStack spacing="4">
                  <Box
                    bg="blue.100"
                    color="blue.800"
                    fontWeight="bold"
                    px="3"
                    py="1"
                    borderRadius="md"
                    fontSize="xs"
                  >
                    ID: {item.id}
                  </Box>
                  <Box>
                    <Text fontSize="xl" fontWeight="bold" className={css({ color: 'gray.900' })}>
                      {item.fishType}
                    </Text>
                    <Text className={css({ color: 'gray.700' })} mt="1">
                      <Text as="span" fontWeight="bold" fontSize="lg">
                        {item.quantity}
                      </Text>{' '}
                      {item.unit}
                      <Text as="span" fontSize="sm" ml="2" className={css({ color: 'gray.500' })}>
                        ({t('Public.AuctionDetail.fisherman_id', { id: item.fishermanId })})
                      </Text>
                    </Text>
                    {item.highestBid && (
                      <Text
                        fontSize="sm"
                        mt="1"
                        className={css({ color: 'orange.600' })}
                        fontWeight="semibold"
                      >
                        {t('Public.AuctionDetail.current_max_bid', {
                          price: item.highestBid.toLocaleString(),
                        })}
                        {item.highestBidderName && (
                          <Text as="span" ml="2" className={css({ color: 'gray.700' })}>
                            {t('Public.AuctionDetail.bidder_name', {
                              name: item.highestBidderName,
                            })}
                          </Text>
                        )}
                      </Text>
                    )}
                  </Box>
                </HStack>
                <ItemStatusBadge status={item.status} />
              </Box>
            </Card>
          ))
        )}
      </Stack>
    </Box>
  );
};
