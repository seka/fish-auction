'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useLogin } from '@/src/features/auth/states/useLogin';
import { LoginForm } from '@/src/features/auth/components/LoginForm';
import { LoginFormData } from '@schemas/auth';
import { Box } from '@atoms';
import { useTranslations } from 'next-intl';

export const LoginContainer = () => {
  const [error, setError] = useState('');
  const router = useRouter();
  const t = useTranslations();
  const { login, isLoading } = useLogin();

  const onSubmit = async (data: LoginFormData) => {
    setError('');

    try {
      const success = await login(data.email, data.password);

      if (success) {
        router.push('/admin');
      } else {
        setError(t('Admin.Login.error_invalid_password'));
      }
    } catch (_err: unknown) {
      setError(t('Admin.Login.error_invalid_password'));
    }
  };

  return (
    <Box display="flex" minH="screen" alignItems="center" justifyContent="center" bg="gray.100">
      <LoginForm onSubmit={onSubmit} isLoading={isLoading} error={error} />
    </Box>
  );
};
