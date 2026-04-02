'use client';

import { Box, Button, Text } from '@atoms';
import { css } from 'styled-system/css';

interface PublicResetPasswordStatusProps {
  type: 'verifying' | 'invalid' | 'complete';
  onButtonClick?: () => void;
  t: (key: string) => string;
}

export const PublicResetPasswordStatus = ({
  type,
  onButtonClick,
  t,
}: PublicResetPasswordStatusProps) => {
  if (type === 'verifying') {
    return (
      <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
        <Text>{t('Auth.ResetPassword.verify_loading')}</Text>
      </Box>
    );
  }

  if (type === 'invalid') {
    return (
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
        <Button mt="6" w="full" variant="secondary" onClick={onButtonClick}>
          {t('Auth.ResetPassword.request_page_button')}
        </Button>
      </Box>
    );
  }

  if (type === 'complete') {
    return (
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
        <Button mt="6" w="full" onClick={onButtonClick}>
          {t('Auth.ResetPassword.login_page_button')}
        </Button>
      </Box>
    );
  }

  return null;
};
