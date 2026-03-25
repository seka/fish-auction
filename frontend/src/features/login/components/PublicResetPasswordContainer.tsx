'use client';

import { useState, useEffect, Suspense } from 'react';
import { useForm } from 'react-hook-form';
import { useRouter, useSearchParams } from 'next/navigation';
import {
  confirmPasswordReset,
  verifyResetToken,
  ResetPasswordConfirmRequest,
} from '@/src/data/api/auth_reset';
import { Box, Button, Text, Stack } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

const PublicResetPasswordForm = () => {
  const t = useTranslations();
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');

  const [isVerifying, setIsVerifying] = useState(true);
  const [isValidToken, setIsValidToken] = useState(false);
  const [isComplete, setIsComplete] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<Omit<ResetPasswordConfirmRequest, 'token'> & { confirm_password: string }>();
  const newPassword = watch('new_password');

  useEffect(() => {
    if (!token) {
      setIsVerifying(false);
      return;
    }

    const verify = async () => {
      try {
        await verifyResetToken({ token });
        setIsValidToken(true);
      } catch (error) {
        console.error('Invalid token', error);
        setIsValidToken(false);
      } finally {
        setIsVerifying(false);
      }
    };
    verify();
  }, [token]);

  const onSubmit = async (
    data: Omit<ResetPasswordConfirmRequest, 'token'> & { confirm_password: string },
  ) => {
    if (!token) return;
    try {
      await confirmPasswordReset({ token, new_password: data.new_password });
      setIsComplete(true);
    } catch (error) {
      console.error('Failed to reset password', error);
      alert(t('Auth.ResetPassword.error_failed'));
    }
  };

  if (isVerifying) {
    return (
      <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
        <Text>{t('Auth.ResetPassword.verify_loading')}</Text>
      </Box>
    );
  }

  if (!token || !isValidToken) {
    return (
      <Box
        minH="screen"
        display="flex"
        alignItems="center"
        justifyContent="center"
        bg="gray.50"
        p="4"
      >
        <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
          <Text
            variant="h3"
            textAlign="center"
            mb="6"
            className={css({ color: 'red.600', fontWeight: 'bold' })}
          >
            {t('Auth.ResetPassword.invalid_link_title')}
          </Text>
          <Text textAlign="center" className={css({ color: 'gray.600' })}>
            {t('Auth.ResetPassword.invalid_link_desc')}
          </Text>
          <Button
            mt="6"
            w="full"
            variant="secondary"
            onClick={() => router.push('/login/forgot_password')}
          >
            {t('Auth.ResetPassword.request_page_button')}
          </Button>
        </Box>
      </Box>
    );
  }

  if (isComplete) {
    return (
      <Box
        minH="screen"
        display="flex"
        alignItems="center"
        justifyContent="center"
        bg="gray.50"
        p="4"
      >
        <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
          <Text
            variant="h3"
            textAlign="center"
            mb="6"
            className={css({ color: 'primary.600', fontWeight: 'bold' })}
          >
            {t('Auth.ResetPassword.complete_title')}
          </Text>
          <Text textAlign="center" className={css({ color: 'gray.600' })}>
            {t('Auth.ResetPassword.complete_desc')}
          </Text>
          <Button mt="6" w="full" onClick={() => router.push('/login/buyer')}>
            {t('Auth.ResetPassword.login_page_button')}
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
      bg="gray.50"
      p="4"
    >
      <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
        <Text
          variant="h3"
          textAlign="center"
          mb="6"
          className={css({ color: 'primary.600', fontWeight: 'bold' })}
        >
          {t('Auth.ResetPassword.new_password_title')}
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
                {t('Auth.ResetPassword.label_new_password')}
              </label>
              <input
                type="password"
                {...register('new_password', {
                  required: t('Validation.required', { field: t('Common.password') }),
                  minLength: { value: 8, message: t('Validation.min_length', { min: 8 }) },
                })}
                className={css({
                  w: 'full',
                  p: '2.5',
                  border: '1px solid',
                  borderColor: 'gray.300',
                  borderRadius: 'md',
                  _focus: {
                    borderColor: 'primary.500',
                    outline: 'none',
                    ring: '2px',
                    ringColor: 'primary.200',
                  },
                })}
              />
              {errors.new_password && (
                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                  {errors.new_password.message}
                </Text>
              )}
            </Box>

            <Box w="full">
              <label
                className={css({
                  display: 'block',
                  mb: '1.5',
                  fontWeight: 'medium',
                  color: 'gray.700',
                })}
              >
                {t('Auth.ResetPassword.label_confirm_password')}
              </label>
              <input
                type="password"
                {...register('confirm_password', {
                  validate: (value) => value === newPassword || t('Validation.password_mismatch'),
                })}
                className={css({
                  w: 'full',
                  p: '2.5',
                  border: '1px solid',
                  borderColor: 'gray.300',
                  borderRadius: 'md',
                  _focus: {
                    borderColor: 'primary.500',
                    outline: 'none',
                    ring: '2px',
                    ringColor: 'primary.200',
                  },
                })}
              />
              {errors.confirm_password && (
                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                  {errors.confirm_password.message}
                </Text>
              )}
            </Box>

            <Button type="submit" w="full" disabled={isSubmitting}>
              {isSubmitting
                ? t('Auth.ResetPassword.changing')
                : t('Auth.ResetPassword.submit_change')}
            </Button>
          </Stack>
        </form>
      </Box>
    </Box>
  );
};

export const PublicResetPasswordContainer = () => {
  const t = useTranslations();
  return (
    <Suspense
      fallback={
        <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.100">
          <Text>{t('Common.loading')}</Text>
        </Box>
      }
    >
      <PublicResetPasswordForm />
    </Suspense>
  );
};
