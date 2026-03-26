'use client';

import { useSettingsManagement } from '../../states/useSettingsManagement';
import { PasswordResetForm } from './PasswordResetForm';
import { Box, Text } from '@atoms';
import { css } from 'styled-system/css';

export const SettingsManagementContainer = () => {
  const { state, actions, t } = useSettingsManagement();

  return (
    <Box p="6">
      <Text
        as="h1"
        fontSize="2xl"
        fontWeight="bold"
        mb="6"
        className={css({ color: 'indigo.900' })}
      >
        {t('Admin.Settings.title')}
      </Text>

      <PasswordResetForm state={state} actions={actions} />
    </Box>
  );
};
