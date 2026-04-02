'use client';

import { UseFormRegister, FieldErrors, UseFormHandleSubmit } from 'react-hook-form';
import { Box, Button, Text, Stack } from '@atoms';
import { css } from 'styled-system/css';

interface AdminResetPasswordFormProps {
  register: UseFormRegister<{ new_password: string; confirm_password: string }>;
  handleSubmit: UseFormHandleSubmit<{ new_password: string; confirm_password: string }>;
  onSubmit: (data: { new_password: string; confirm_password: string }) => Promise<void>;
  errors: FieldErrors<{ new_password: string; confirm_password: string }>;
  isSubmitting: boolean;
  t: (key: string, values?: Record<string, string | number>) => string;
}

export const AdminResetPasswordForm = ({
  register,
  handleSubmit,
  onSubmit,
  errors,
  isSubmitting,
  t,
}: AdminResetPasswordFormProps) => {
  return (
    <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
      <Text
        variant="h3"
        textAlign="center"
        mb="6"
        className={css({ color: 'indigo.700', fontWeight: 'bold' })}
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
              {...register('new_password')}
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
            />
            {errors.new_password && (
              <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                {errors.new_password.message}
              </Text>
            )}
            <Text className={css({ color: 'gray.500', fontSize: 'xs', mt: '1.5' })}>
              {t('Validation.password_complexity_hint')}
            </Text>
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
              {...register('confirm_password')}
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
            />
            {errors.confirm_password && (
              <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>
                {errors.confirm_password.message}
              </Text>
            )}
          </Box>

          <Button type="submit" w="full" disabled={isSubmitting} variant="primary">
            {isSubmitting
              ? t('Auth.ResetPassword.changing')
              : t('Auth.ResetPassword.submit_change')}
          </Button>
        </Stack>
      </form>
    </Box>
  );
};
