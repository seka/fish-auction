'use client';

import { Box, Button, Stack, Text, Input } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface PasswordResetFormProps {
  state: {
    currentPassword: string;
    newPassword: string;
    confirmPassword: string;
    message: { type: 'success' | 'error'; text: string } | null;
    isLoading: boolean;
  };
  actions: {
    setCurrentPassword: (val: string) => void;
    setNewPassword: (val: string) => void;
    setConfirmPassword: (val: string) => void;
    handleSubmit: (e: React.FormEvent) => void;
  };
}

export const PasswordResetForm = ({ state, actions }: PasswordResetFormProps) => {
  const t = useTranslations();

  return (
    <Box bg="white" p="6" borderRadius="lg" shadow="sm" maxW="md">
      <Text
        as="h2"
        fontSize="lg"
        fontWeight="semibold"
        mb="4"
        className={css({ color: 'gray.800' })}
      >
        {t('Auth.ResetPassword.title')}
      </Text>

      {state.message && (
        <Box
          p="3"
          mb="4"
          borderRadius="md"
          bg={state.message.type === 'success' ? 'green.50' : 'red.50'}
          color={state.message.type === 'success' ? 'green.700' : 'red.700'}
          border="1px solid"
          borderColor={state.message.type === 'success' ? 'green.200' : 'red.200'}
        >
          {state.message.text}
        </Box>
      )}

      <form onSubmit={actions.handleSubmit}>
        <Stack spacing="4" alignItems="stretch">
          <Box>
            <label
              htmlFor="current-password"
              className={css({
                display: 'block',
                mb: '1',
                fontSize: 'sm',
                fontWeight: 'medium',
                color: 'gray.700',
              })}
            >
              {t('Validation.field_name.password')}
            </label>
            <Input
              id="current-password"
              type="password"
              value={state.currentPassword}
              onChange={(e) => actions.setCurrentPassword(e.target.value)}
              required
              className={css({ w: 'full' })}
            />
          </Box>

          <Box>
            <label
              htmlFor="new-password"
              className={css({
                display: 'block',
                mb: '1',
                fontSize: 'sm',
                fontWeight: 'medium',
                color: 'gray.700',
              })}
            >
              {t('Auth.ResetPassword.label_new_password')}
            </label>
            <Input
              id="new-password"
              type="password"
              value={state.newPassword}
              onChange={(e) => actions.setNewPassword(e.target.value)}
              required
              minLength={8}
              className={css({ w: 'full' })}
            />
            <Text className={css({ color: 'gray.500', fontSize: 'xs', mt: '1.5' })}>
              {t('Validation.password_complexity_hint')}
            </Text>
          </Box>

          <Box>
            <label
              htmlFor="confirm-password"
              className={css({
                display: 'block',
                mb: '1',
                fontSize: 'sm',
                fontWeight: 'medium',
                color: 'gray.700',
              })}
            >
              {t('Auth.ResetPassword.label_confirm_password')}
            </label>
            <Input
              id="confirm-password"
              type="password"
              value={state.confirmPassword}
              onChange={(e) => actions.setConfirmPassword(e.target.value)}
              required
              minLength={8}
              className={css({ w: 'full' })}
            />
          </Box>

          <Button
            type="submit"
            disabled={state.isLoading}
            className={css({
              bg: 'indigo.600',
              color: 'white',
              _hover: { bg: 'indigo.700' },
              _disabled: { opacity: 0.6, cursor: 'not-allowed' },
            })}
          >
            {state.isLoading ? t('Common.loading') : t('Public.MyPage.submit_password_update')}
          </Button>
        </Stack>
      </form>
    </Box>
  );
};
