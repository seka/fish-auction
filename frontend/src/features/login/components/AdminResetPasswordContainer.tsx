'use client';

import { useState, useEffect, Suspense, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useRouter, useSearchParams } from 'next/navigation';
import {
  verifyAdminResetToken,
  confirmAdminPasswordReset,
  ResetPasswordConfirmRequest,
} from '@/src/data/api/admin_auth_reset';
import { Box, Text } from '@atoms';
import { useTranslations } from 'next-intl';
import { AdminResetPasswordForm } from './AdminResetPassword/AdminResetPasswordForm';
import { AdminResetPasswordStatus } from './AdminResetPassword/AdminResetPasswordStatus';

import { zodResolver } from '@hookform/resolvers/zod';
import { getResetPasswordSchema } from '@schema/password';

const AdminResetPasswordContent = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');

  const [isVerifying, setIsVerifying] = useState(true);
  const [isValidToken, setIsValidToken] = useState(false);
  const [isComplete, setIsComplete] = useState(false);

  const schema = useMemo(() => getResetPasswordSchema(tValidation), [tValidation]);
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm({
    resolver: zodResolver(schema),
  });

  useEffect(() => {
    if (!token) {
      setIsVerifying(false);
      return;
    }

    const verify = async () => {
      try {
        await verifyAdminResetToken({ token });
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
      await confirmAdminPasswordReset({ token, new_password: data.new_password });
      setIsComplete(true);
    } catch (error) {
      console.error('Failed to reset password', error);
      alert(t('Auth.ResetPassword.error_failed'));
    }
  };

  const handleBackToRequest = () => router.push('/login/admin/forgot_password');
  const handleBackToLogin = () => router.push('/login');

  if (isVerifying) {
    return <AdminResetPasswordStatus type="verifying" t={t} />;
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
      {!token || !isValidToken ? (
        <AdminResetPasswordStatus type="invalid" onButtonClick={handleBackToRequest} t={t} />
      ) : isComplete ? (
        <AdminResetPasswordStatus type="complete" onButtonClick={handleBackToLogin} t={t} />
      ) : (
        <AdminResetPasswordForm
          register={register}
          handleSubmit={handleSubmit}
          onSubmit={onSubmit}
          errors={errors}
          isSubmitting={isSubmitting}
          t={t}
        />
      )}
    </Box>
  );
};

export const AdminResetPasswordContainer = () => {
  const t = useTranslations();
  return (
    <Suspense
      fallback={
        <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.100">
          <Text>{t('Common.loading')}</Text>
        </Box>
      }
    >
      <AdminResetPasswordContent />
    </Suspense>
  );
};
