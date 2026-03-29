'use client';

import { Box, Button, Stack, Text, Input, Card, HStack } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';
import { UseFormReturn } from 'react-hook-form';
import { BuyerFormData } from '@/src/models/schemas/admin';

interface BuyerFormProps {
  form: UseFormReturn<BuyerFormData>;
  onSubmit: (e: React.FormEvent) => void;
  isCreating: boolean;
}

export const BuyerForm = ({ form, onSubmit, isCreating }: BuyerFormProps) => {
  const t = useTranslations();
  const {
    register,
    formState: { errors },
  } = form;

  return (
    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
      <HStack mb="6">
        <Box w="2" h="6" bg="green.500" mr="3" borderRadius="full" />
        <Text as="h2" variant="h4" className={css({ color: 'green.900' })} fontWeight="bold">
          {t('Admin.Buyers.register_title')}
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
              {t('Admin.Buyers.name')}
            </Text>
            <Input
              type="text"
              {...register('name')}
              placeholder={t('Admin.Buyers.placeholder_name')}
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
              {t('Validation.field_name.email')}
            </Text>
            <Input
              type="email"
              {...register('email')}
              placeholder={t('Common.email')}
              error={!!errors.email}
            />
            {errors.email && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.email.message)}
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
              {t('Validation.field_name.password')}
            </Text>
            <Input
              type="password"
              {...register('password')}
              placeholder="********"
              error={!!errors.password}
            />
            {errors.password && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.password.message)}
              </Text>
            )}
            <Text className={css({ color: 'gray.500', fontSize: 'xs', mt: '1.5' })}>
              {t('Validation.password_complexity_hint')}
            </Text>
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
              {t('Validation.field_name.organization')}
            </Text>
            <Input
              type="text"
              {...register('organization')}
              placeholder={t('Admin.Buyers.placeholder_name_company')}
              error={!!errors.organization}
            />
            {errors.organization && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.organization.message)}
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
              {t('Validation.field_name.contact_info')}
            </Text>
            <Input
              type="text"
              {...register('contactInfo')}
              placeholder="03-1234-5678"
              error={!!errors.contactInfo}
            />
            {errors.contactInfo && (
              <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">
                {String(errors.contactInfo.message)}
              </Text>
            )}
          </Box>

          <Button type="submit" disabled={isCreating} width="full" variant="primary">
            {isCreating ? t('Common.loading') : t('Common.register')}
          </Button>
        </Stack>
      </form>
    </Card>
  );
};
