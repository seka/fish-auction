'use client';

import { Box, Button, Stack, Text, Input, Card, HStack, Select } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { UseFormReturn } from 'react-hook-form';
import { AuctionFormInput } from '@schemas/auction';
import { Venue, Auction } from '../../types';

interface AuctionFormProps {
  form: UseFormReturn<AuctionFormInput>;
  onSubmit: (e: React.FormEvent) => void;
  onCancelEdit: () => void;
  isCreating: boolean;
  isUpdating: boolean;
  editingAuction: Auction | null;
  venues: Venue[];
}

export const AuctionForm = ({
  form,
  onSubmit,
  onCancelEdit,
  isCreating,
  isUpdating,
  editingAuction,
  venues,
}: AuctionFormProps) => {
  const t = useTranslations();
  const {
    register,
    formState: { errors },
  } = form;

  return (
    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
      <HStack mb="6">
        <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
        <Text as="h2" variant="h4" className={css({ color: 'indigo.900' })} fontWeight="bold">
          {editingAuction ? t('Admin.Auctions.edit_title') : t('Admin.Auctions.register_title')}
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
              {t('Admin.Auctions.venue')}
            </Text>
            <Select {...register('venueId')}>
              <option value="">{t('Admin.Auctions.placeholder_select_venue')}</option>
              {venues.map((venue) => (
                <option key={venue.id} value={venue.id}>
                  {venue.name}
                </option>
              ))}
            </Select>
            {errors.venueId && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.venueId.message)}
              </Text>
            )}
          </Box>
          <HStack spacing="4">
            <Box flex="1">
              <Text
                as="label"
                display="block"
                fontSize="sm"
                fontWeight="bold"
                className={css({ color: 'gray.700' })}
                mb="1"
              >
                {t('Admin.Auctions.start_time')}
              </Text>
              <Input type="datetime-local" {...register('startAt')} error={!!errors.startAt} />
              {errors.startAt && (
                <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                  {String(errors.startAt.message)}
                </Text>
              )}
            </Box>
            <Box flex="1">
              <Text
                as="label"
                display="block"
                fontSize="sm"
                fontWeight="bold"
                className={css({ color: 'gray.700' })}
                mb="1"
              >
                {t('Admin.Auctions.end_time')}
              </Text>
              <Input type="datetime-local" {...register('endAt')} error={!!errors.endAt} />
            </Box>
          </HStack>

          <HStack spacing="2" pt="4">
            <Button
              type="submit"
              disabled={isCreating || isUpdating}
              width="full"
              className={css({ flex: '1' })}
              variant="primary"
            >
              {editingAuction
                ? isUpdating
                  ? t('Common.loading')
                  : t('Common.update')
                : isCreating
                  ? t('Common.loading')
                  : t('Common.register')}
            </Button>
            {editingAuction && (
              <Button type="button" onClick={onCancelEdit} variant="outline">
                {t('Common.cancel')}
              </Button>
            )}
          </HStack>
        </Stack>
      </form>
    </Card>
  );
};
