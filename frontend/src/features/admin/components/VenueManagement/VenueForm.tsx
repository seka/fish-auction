'use client';

import { Box, Button, Stack, Text, Input, Card, HStack } from '@atoms';
import { css } from 'styled-system/css';
import { styled } from 'styled-system/jsx';
import { useTranslations } from 'next-intl';
import { VenueFormData } from '@/src/models/schemas/auction';
import { Venue } from '@/src/models';

// ... (Textarea definition)

interface VenueFormProps {
  form: UseFormReturn<VenueFormData>;
  onSubmit: (e: React.FormEvent) => void;
  onCancelEdit: () => void;
  isCreating: boolean;
  isUpdating: boolean;
  editingVenue: Venue | null;
}

export const VenueForm = ({
  form,
  onSubmit,
  onCancelEdit,
  isCreating,
  isUpdating,
  editingVenue,
}: VenueFormProps) => {
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
          {editingVenue ? t('Admin.Venues.edit_title') : t('Admin.Venues.register_title')}
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
              {t('Admin.Venues.name')}
            </Text>
            <Input
              type="text"
              {...register('name')}
              placeholder={t('Admin.Venues.placeholder_name')}
              error={!!errors.name}
            />
            {errors.name && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.name.message)}
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
              {t('Admin.Venues.location')}
            </Text>
            <Input
              type="text"
              {...register('location')}
              placeholder={t('Admin.Venues.placeholder_location')}
            />
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
              {t('Admin.Venues.description')}
            </Text>
            <Textarea
              {...register('description')}
              rows={3}
              placeholder={t('Admin.Venues.placeholder_description')}
            />
          </Box>

          <HStack spacing="2">
            <Button
              type="submit"
              disabled={isCreating || isUpdating}
              width="full"
              className={css({ flex: '1' })}
              variant="primary"
            >
              {editingVenue
                ? isUpdating
                  ? t('Common.loading')
                  : t('Common.update')
                : isCreating
                  ? t('Common.loading')
                  : t('Common.register')}
            </Button>
            {editingVenue && (
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
