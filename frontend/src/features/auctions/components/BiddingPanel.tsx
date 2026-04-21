'use client';

import { UseFormReturn } from 'react-hook-form';
import { BidFormData } from '@schemas/auction';
import { AuctionItem, Auction } from '../types';
import { Box, Text, Card, Stack, Input, Button } from '@atoms';
import { css } from 'styled-system/css';

import { useTranslations } from 'next-intl';

interface BiddingPanelProps {
  selectedItem: AuctionItem | null;
  auction: Auction;
  auctionActive: boolean;
  bidForm: UseFormReturn<BidFormData>;
  onSubmitBid: (data: BidFormData) => Promise<void>;
  isBidLoading: boolean;
  t: ReturnType<typeof useTranslations>;
}

export const BiddingPanel = ({
  selectedItem,
  auction,
  auctionActive,
  bidForm,
  onSubmitBid,
  isBidLoading,
  t,
}: BiddingPanelProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = bidForm;

  return (
    <Box gridColumn={{ base: '1', lg: 'span 1' }}>
      <Card
        padding="lg"
        shadow="lg"
        className={css({
          borderWidth: '1px',
          borderColor: 'gray.200',
          position: { lg: 'sticky' },
          top: '6',
        })}
      >
        <Text
          fontSize="xl"
          fontWeight="bold"
          className={css({ color: 'gray.800' })}
          borderBottom="1px solid"
          borderColor="gray.200"
          pb="2"
          mb="6"
        >
          {t('Public.AuctionDetail.bidding_panel')}
        </Text>
        {selectedItem ? (
          <form onSubmit={handleSubmit(onSubmitBid)}>
            <Stack spacing="6">
              <Box p="5" bg="gray.50" borderRadius="lg" borderWidth="1px" borderColor="gray.200">
                <Text fontSize="sm" className={css({ color: 'gray.600' })} mb="1">
                  {t('Public.AuctionDetail.selected_item')}
                </Text>
                <Text fontWeight="bold" fontSize="2xl" className={css({ color: 'gray.900' })}>
                  {selectedItem.fishType}
                </Text>
                <Text fontSize="lg" className={css({ color: 'gray.700' })}>
                  {selectedItem.quantity.label}
                </Text>
                {selectedItem.bidding.highestBid && (
                  <Text
                    fontSize="sm"
                    mt="2"
                    className={css({ color: 'orange.600' })}
                    fontWeight="bold"
                  >
                    {t('Public.AuctionDetail.current_max_bid', {
                      price: selectedItem.price.label,
                    })}
                    {selectedItem.bidding.highestBidderName && (
                      <Text as="span" ml="2" className={css({ color: 'gray.700' })}>
                        {t('Public.AuctionDetail.bidder_name', {
                          name: selectedItem.bidding.highestBidderName,
                        })}
                      </Text>
                    )}
                  </Text>
                )}
              </Box>

              {auctionActive ? (
                <>
                  <Box>
                    <Text
                      as="label"
                      display="block"
                      fontSize="sm"
                      fontWeight="bold"
                      className={css({ color: 'gray.700' })}
                      mb="1"
                    >
                      {t('Public.AuctionDetail.bid_amount_label')}
                    </Text>
                    <Text fontSize="xs" className={css({ color: 'gray.500', mb: '2' })}>
                      {t('Public.AuctionDetail.next_min_bid', {
                        price: selectedItem.bidding.nextMinBid.label,
                      })}
                    </Text>
                    <Box position="relative">
                      <Box
                        position="absolute"
                        top="50%"
                        left="3"
                        transform="translateY(-50%)"
                        pointerEvents="none"
                      >
                        <Text fontSize="sm" className={css({ color: 'gray.600' })}>
                          ¥
                        </Text>
                      </Box>
                      <Input
                        type="number"
                        {...register('price')}
                        placeholder={selectedItem.bidding.nextMinBid.value.toString()}
                        className={css({ pl: '7' })}
                      />
                    </Box>
                    {errors.price && (
                      <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                        {errors.price.message}
                      </Text>
                    )}
                  </Box>

                  <Button
                    type="submit"
                    disabled={isBidLoading}
                    width="full"
                    size="lg"
                    className={css({
                      bg: 'red.600',
                      _hover: { bg: 'red.700', transform: 'scale(1.02)' },
                      color: 'white',
                      shadow: 'md',
                      transition: 'all 0.2s',
                    })}
                  >
                    {isBidLoading
                      ? t('Public.AuctionDetail.bidding_process')
                      : t('Public.AuctionDetail.bid_button')}
                  </Button>
                </>
              ) : (
                <Box
                  textAlign="center"
                  py="6"
                  bg="yellow.50"
                  borderRadius="lg"
                  borderWidth="1px"
                  borderColor="yellow.200"
                >
                  <Text className={css({ color: 'yellow.800' })} fontWeight="bold" mb="2">
                    {t('Public.AuctionDetail.out_of_hours_title')}
                  </Text>
                  {auction.duration.label && (
                    <Text fontSize="sm" className={css({ color: 'yellow.700' })}>
                      {t('Public.AuctionDetail.out_of_hours_msg', {
                        duration: auction.duration.label,
                      })}
                    </Text>
                  )}
                </Box>
              )}
            </Stack>
          </form>
        ) : (
          <Box textAlign="center" py="12" color="gray.400">
            <Text className={css({ whiteSpace: 'pre-line' })}>
              {t('Public.AuctionDetail.select_instruction')}
            </Text>
          </Box>
        )}
      </Card>
    </Box>
  );
};
