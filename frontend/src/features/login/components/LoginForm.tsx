'use client';

import { useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { getLoginSchema, LoginFormData } from '@/src/models/schemas/auth';
import { css } from 'styled-system/css';
import { Box, Text, Button, Input, Card, Stack } from '@atoms';
import { useTranslations } from 'next-intl';
import Link from 'next/link';

interface LoginFormProps {
  onSubmit: (data: LoginFormData) => Promise<void>;
  isLoading: boolean;
  error: string;
}

export const LoginForm = ({ onSubmit, isLoading, error }: LoginFormProps) => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const schema = useMemo(() => getLoginSchema(tValidation), [tValidation]);
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(schema),
  });

  return (
    <Card width="full" maxW="md" padding="lg" shadow="lg">
      <Box textAlign="center" mb="8">
        <Stack spacing="4">
          <Text as="h1" variant="h3" fontWeight="bold" className={css({ color: 'indigo.700' })}>
            {t('Admin.Login.title')}
          </Text>
          <Text className={css({ color: 'gray.600' })}>{t('Admin.Login.description')}</Text>
        </Stack>
      </Box>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Stack spacing="6">
          <Box>
            <label htmlFor="email" className={css({ srOnly: true })}>
              {t('Common.email')}
            </label>
            <Input
              id="email"
              type="email"
              {...register('email')}
              placeholder={t('Common.email')}
              bg="white"
            />
            {errors.email && (
              <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                {errors.email.message}
              </Text>
            )}
          </Box>
          <Box>
            <label htmlFor="password" className={css({ srOnly: true })}>
              {t('Common.password')}
            </label>
            <Input
              id="password"
              type="password"
              {...register('password')}
              placeholder={t('Common.password')}
              bg="white"
            />
            {errors.password && (
              <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                {errors.password.message}
              </Text>
            )}
          </Box>
          {error && (
            <Box bg="red.50" p="3" borderRadius="md">
              <Text className={css({ color: 'red.600', fontSize: 'sm', textAlign: 'center' })}>
                {error}
              </Text>
            </Box>
          )}
          <Button type="submit" width="full" size="lg" disabled={isLoading} variant="primary">
            {isLoading ? t('Admin.Login.logging_in') : t('Common.submit')}
          </Button>
          <Box textAlign="center">
            <Link
              href="/login/admin/forgot_password"
              className={css({
                fontSize: 'sm',
                color: 'gray.500',
                _hover: { textDecoration: 'underline' },
                cursor: 'pointer',
                display: 'block',
                mb: '2',
              })}
            >
              {t('Auth.ForgotPassword.link_text')}
            </Link>
          </Box>
        </Stack>
      </form>
    </Card>
  );
};
