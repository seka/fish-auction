'use client';

import { Box, Text, Button, Stack, Input } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface PasswordMessage {
  text: string;
  type: 'info' | 'error' | 'success';
}

interface SettingsFormProps {
  currentPassword: string;
  setCurrentPassword: (value: string) => void;
  newPassword: string;
  setNewPassword: (value: string) => void;
  confirmPassword: string;
  setConfirmPassword: (value: string) => void;
  passwordMessage: PasswordMessage | null;
  handleUpdatePassword: (e: React.FormEvent) => void;
  isPasswordUpdating: boolean;
}

export const SettingsForm = ({
  currentPassword,
  setCurrentPassword,
  newPassword,
  setNewPassword,
  confirmPassword,
  setConfirmPassword,
  passwordMessage,
  handleUpdatePassword,
  isPasswordUpdating,
}: SettingsFormProps) => {
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
        {t('Public.MyPage.password_change_title')}
      </Text>

      {passwordMessage && (
        <Box
          p="3"
          mb="4"
          borderRadius="md"
          bg={passwordMessage.type === 'success' ? 'green.50' : 'red.50'}
          color={passwordMessage.type === 'success' ? 'green.700' : 'red.700'}
          border="1px solid"
          borderColor={passwordMessage.type === 'success' ? 'green.200' : 'red.200'}
        >
          {passwordMessage.text}
        </Box>
      )}

      <form onSubmit={handleUpdatePassword}>
        <Stack spacing="4" alignItems="stretch">
          <Box>
            <Text
              as="label"
              display="block"
              mb="1"
              fontSize="sm"
              fontWeight="medium"
              className={css({ color: 'gray.700' })}
            >
              {t('Validation.field_name.password')}
            </Text>
            <Input
              type="password"
              value={currentPassword}
              onChange={(e) => setCurrentPassword(e.target.value)}
              required
              className={css({ w: 'full' })}
            />
          </Box>

          <Box>
            <Text
              as="label"
              display="block"
              mb="1"
              fontSize="sm"
              fontWeight="medium"
              className={css({ color: 'gray.700' })}
            >
              {t('Auth.ResetPassword.label_new_password')}
            </Text>
            <Input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
              minLength={8}
              className={css({ w: 'full' })}
            />
          </Box>

          <Box>
            <Text
              as="label"
              display="block"
              mb="1"
              fontSize="sm"
              fontWeight="medium"
              className={css({ color: 'gray.700' })}
            >
              {t('Auth.ResetPassword.label_confirm_password')}
            </Text>
            <Input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              minLength={8}
              className={css({ w: 'full' })}
            />
          </Box>

          <Button
            type="submit"
            disabled={isPasswordUpdating}
            className={css({
              bg: 'indigo.600',
              color: 'white',
              _hover: { bg: 'indigo.700' },
              _disabled: { opacity: 0.6, cursor: 'not-allowed' },
            })}
          >
            {isPasswordUpdating
              ? t('Public.MyPage.updating')
              : t('Public.MyPage.submit_password_update')}
          </Button>
        </Stack>
      </form>
    </Box>
  );
};
