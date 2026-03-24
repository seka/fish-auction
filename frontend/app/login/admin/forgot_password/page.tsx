'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useRouter } from 'next/navigation';
import { requestAdminPasswordReset, ResetPasswordRequest } from '@/src/data/api/admin_auth_reset';
import { Box, Button, Text, Stack } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

export default function AdminForgotPasswordPage() {
  const t = useTranslations();
  const router = useRouter();
  const [isSubmitted, setIsSubmitted] = useState(false);
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ResetPasswordRequest>();

  const onSubmit = async (data: ResetPasswordRequest) => {
    try {
      await requestAdminPasswordReset(data);
      setIsSubmitted(true);
    } catch (error) {
      console.error('Failed to request password reset', error);
      setIsSubmitted(true); // Treat as success to prevent email enumeration
    }
  };

  if (isSubmitted) {
    return (
      <Box
        minH="screen"
        display="flex"
        alignItems="center"
        justifyContent="center"
        bg="gray.100"
        p="4"
      >
        <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
          <Text
            variant="h3"
            textAlign="center"
            mb="6"
            className={css({ color: 'indigo.700', fontWeight: 'bold' })}
          >
            {t('Auth.ForgotPassword.success_title')}
          </Text>
          <Text textAlign="center" className={css({ color: 'gray.600' })}>
            {t('Auth.ForgotPassword.admin_success_description')}
          </Text>
          <Button mt="6" w="full" variant="secondary" onClick={() => router.push('/login')}>
            {t('Auth.ForgotPassword.back_to_admin_login')}
          </Button>
        </Box>
      </Box>
    );
  }

  return (
    <Box
      minH="screen"
      display="flex"
      alignItems="center"
      justifyContent="center"
      bg="gray.100"
      p="4"
    >
      <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
        <Text
          variant="h3"
          textAlign="center"
          mb="6"
          className={css({ color: 'indigo.700', fontWeight: 'bold' })}
        >
          {t('Auth.ForgotPassword.admin_title')}
        </Text>
        <Text mb="6" className={css({ color: 'gray.600', textAlign: 'center' })}>
          {t('Auth.ForgotPassword.description')}
        </Text>

        <form onSubmit={handleSubmit(onSubmit)}>
          <Stack gap="4">
            <Box w="full">
              <label
                className={css({
                  display: 'block',
                  mb: '1.5',
                  fontWeight: 'medium',
                  color: 'gray.700',
                })}
              >
                {t('Common.email')}
              </label>
              <input
                {...register('email', {
                  required: t('Validation.required', { field: t('Common.email') }),
                  pattern: {
                    value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                    message: t('Validation.invalid_email'),
                  },
                })}
                className={css({
                  w: 'full',
                  p: '2.5',
                  border: '1px solid',
                  borderColor: 'gray.300',
                  borderRadius: 'md',
                  _focus: {
                    borderColor: 'indigo.500',
                    outline: 'none',
                    ring: '2px',
                    ringColor: 'indigo.200',
                  },
                })}
                placeholder="admin@fish-auction.com"
              />
              {errors.email && (
                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                  {errors.email.message}
                </Text>
              )}
            </Box>

            <Button type="submit" w="full" disabled={isSubmitting} variant="primary">
              {isSubmitting ? t('Auth.ForgotPassword.sending') : t('Auth.ForgotPassword.submit')}
            </Button>

            <Button
              variant="outline"
              width="full"
              onClick={() => router.push('/login')}
              style={{ border: 'none' }}
            >
              {t('Common.cancel')}
            </Button>
          </Stack>
        </form>
      </Box>
    </Box>
  );
}
