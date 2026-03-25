'use client';

import { Box, Button, Stack, Text, Input, Card, HStack } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { UseFormReturn } from 'react-hook-form';

import { FishermanFormData } from '@/src/models/schemas/admin';

interface FishermanFormProps {
  form: UseFormReturn<FishermanFormData>;
  onSubmit: (e: React.FormEvent) => void;
  isCreating: boolean;
}

export const FishermanForm = ({ form, onSubmit, isCreating }: FishermanFormProps) => {
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
          {t('Admin.Fishermen.register_title')}
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
              {t('Admin.Fishermen.name')}
            </Text>
            <Input
              type="text"
              {...register('name')}
              placeholder={t('Admin.Fishermen.placeholder_name')}
              error={!!errors.name}
            />
            {errors.name && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.name.message)}
              </Text>
            )}
          </Box>

          <Button
            type="submit"
            disabled={isCreating}
            width="full"
            className={css({ flex: '1' })}
            variant="primary"
          >
            {isCreating ? t('Common.loading') : t('Common.register')}
          </Button>
        </Stack>
      </form>
    </Card>
  );
};
