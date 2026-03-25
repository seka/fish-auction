'use client';

import { useState, useEffect, Suspense } from 'react';
import { useForm } from 'react-hook-form';
import { useRouter, useSearchParams } from 'next/navigation';
import {
  confirmPasswordReset,
  verifyResetToken,
  ResetPasswordConfirmRequest,
} from '@/src/data/api/auth_reset';
import { Box, Text } from '@atoms';
import { useTranslations } from 'next-intl';
import { PublicResetPasswordForm } from './PublicResetPassword/PublicResetPasswordForm';
import { PublicResetPasswordStatus } from './PublicResetPassword/PublicResetPasswordStatus';

const PublicResetPasswordContent = () => {
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
  const newPasswordValue = watch('new_password');

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

  const handleBackToRequest = () => router.push('/login/forgot_password');
  const handleBackToLogin = () => router.push('/login/buyer');

  if (isVerifying) {
    return <PublicResetPasswordStatus type="verifying" t={t} />;
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
      {!token || !isValidToken ? (
        <PublicResetPasswordStatus type="invalid" onButtonClick={handleBackToRequest} t={t} />
      ) : isComplete ? (
        <PublicResetPasswordStatus type="complete" onButtonClick={handleBackToLogin} t={t} />
      ) : (
        <PublicResetPasswordForm
          register={register}
          handleSubmit={handleSubmit}
          onSubmit={onSubmit}
          errors={errors}
          isSubmitting={isSubmitting}
          newPasswordValue={newPasswordValue}
          t={t}
        />
      )}
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
      <PublicResetPasswordContent />
    </Suspense>
  );
};
