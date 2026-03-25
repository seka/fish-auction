'use client';

import { Box, Button, Stack, Text, Input, Card, HStack, Select } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { UseFormReturn } from 'react-hook-form';
import { ItemFormData } from '@/src/models/schemas/admin';
import { Auction, Fisherman, AuctionItem } from '@/src/models';

interface ItemFormProps {
  form: UseFormReturn<ItemFormData>;
  onSubmit: (e: React.FormEvent) => void;
  onCancelEdit: () => void;
  isCreating: boolean;
  isUpdating: boolean;
  editingItem: AuctionItem | null;
  auctions: Auction[];
  fishermen: Fisherman[];
}

export const ItemForm = ({
  form,
  onSubmit,
  onCancelEdit,
  isCreating,
  isUpdating,
  editingItem,
  auctions,
  fishermen,
}: ItemFormProps) => {
  const t = useTranslations();
  const {
    register,
    formState: { errors },
  } = form;

  return (
    <Card padding="lg">
      <HStack mb="6">
        <Box w="2" h="6" bg="orange.500" mr="3" borderRadius="full" />
        <Text as="h2" variant="h4" className={css({ color: 'orange.900' })} fontWeight="bold">
          {editingItem ? t('Admin.Items.edit_item') : t('Admin.Items.register_title')}
        </Text>
      </HStack>
      <form onSubmit={onSubmit}>
        <Stack spacing="4">
          <Box>
            <Text
              as="label"
              display="block"
              fontSize="sm"
              fontWeight="bold"
              className={css({ color: 'gray.700' })}
              mb="1"
            >
              {t('Admin.Items.auction')}
            </Text>
            <Select {...register('auctionId')}>
              <option value="">{t('Admin.Items.placeholder_select_auction')}</option>
              {auctions.map((auction) => (
                <option key={auction.id} value={auction.id}>
                  {auction.auctionDate} {auction.startTime?.substring(0, 5)} -{' '}
                  {auction.endTime?.substring(0, 5)} (ID: {auction.id})
                </option>
              ))}
            </Select>
            {errors.auctionId && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.auctionId.message)}
              </Text>
            )}
          </Box>
          <Box>
            <Text
              as="label"
              display="block"
              fontSize="sm"
              fontWeight="bold"
              className={css({ color: 'gray.700' })}
              mb="1"
            >
              {t('Admin.Items.fisherman')}
            </Text>
            <Select {...register('fishermanId')}>
              <option value="">{t('Admin.Items.placeholder_select_fisherman')}</option>
              {fishermen.map((fisherman) => (
                <option key={fisherman.id} value={fisherman.id}>
                  {fisherman.name}
                </option>
              ))}
            </Select>
            {errors.fishermanId && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.fishermanId.message)}
              </Text>
            )}
          </Box>
          <Box>
            <Text
              as="label"
              display="block"
              fontSize="sm"
              fontWeight="bold"
              className={css({ color: 'gray.700' })}
              mb="1"
            >
              {t('Admin.Items.fish_type')}
            </Text>
            <Input
              type="text"
              {...register('fishType')}
              placeholder={t('Admin.Items.placeholder_fish_type')}
              error={!!errors.fishType}
            />
            {errors.fishType && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.fishType.message)}
              </Text>
            )}
          </Box>
          <Box>
            <Text
              as="label"
              display="block"
              fontSize="sm"
              fontWeight="bold"
              className={css({ color: 'gray.700' })}
              mb="1"
            >
              {t('Admin.Items.quantity')}
            </Text>
            <Input
              type="number"
              {...register('quantity')}
              placeholder={t('Admin.Items.placeholder_quantity')}
              error={!!errors.quantity}
            />
            {errors.quantity && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.quantity.message)}
              </Text>
            )}
          </Box>
          <Box>
            <Text
              as="label"
              display="block"
              fontSize="sm"
              fontWeight="bold"
              className={css({ color: 'gray.700' })}
              mb="1"
            >
              {t('Admin.Items.unit')}
            </Text>
            <Input
              type="text"
              {...register('unit')}
              placeholder={t('Admin.Items.placeholder_unit')}
              error={!!errors.unit}
            />
            {errors.unit && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.unit.message)}
              </Text>
            )}
          </Box>
          <HStack spacing="2" mt="4">
            <Button
              type="submit"
              disabled={isCreating || isUpdating}
              width="full"
              variant="primary"
            >
              {isCreating || isUpdating
                ? t('Common.loading')
                : editingItem
                  ? t('Common.save')
                  : t('Common.register')}
            </Button>
            {editingItem && (
              <Button type="button" onClick={onCancelEdit} width="full" variant="outline">
                {t('Common.cancel')}
              </Button>
            )}
          </HStack>
        </Stack>
      </form>
    </Card>
  );
};
