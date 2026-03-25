'use client';

import { useFishermanManagement } from '../../hooks/useFishermanManagement';
import { FishermanList } from './FishermanList';
import { FishermanForm } from './FishermanForm';
import { Box, Card, Text } from '@atoms';
import { css } from 'styled-system/css';

export const FishermanManagementContainer = () => {
  const { state, form, actions, t } = useFishermanManagement();

  return (
    <Box maxW="5xl" mx="auto" p="6">
      <Text
        as="h1"
        variant="h2"
        className={css({ color: 'gray.800' })}
        mb="8"
        pb="4"
        borderBottom="1px solid"
        borderColor="gray.200"
      >
        {t('Admin.Fishermen.title')}
      </Text>

      {state.message && (
        <Box
          bg="blue.50"
          borderLeft="4px solid"
          borderColor="blue.500"
          color="blue.700"
          p="4"
          mb="8"
          borderRadius="sm"
          shadow="sm"
          role="alert"
        >
          <Text fontWeight="bold">{t('Common.notification')}</Text>
          <Text>{state.message}</Text>
        </Box>
      )}

      <Box
        display="grid"
        gridTemplateColumns={{ base: '1fr', md: '3fr 1fr' }}
        gap="8"
        className={css({ md: { gridTemplateColumns: '1fr 2fr' } })}
      >
        {/* Form Section */}
        <Box className={css({ md: { gridColumn: '1 / 2' } })}>
          <FishermanForm form={form} onSubmit={actions.onSubmit} isCreating={state.isCreating} />
        </Box>

        {/* List Section */}
        <Box className={css({ md: { gridColumn: '2 / 3' } })}>
          <Card padding="none" overflow="hidden">
            <Box p="6" borderBottom="1px solid" borderColor="gray.200">
              <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">
                {t('Admin.Fishermen.list_title')}
              </Text>
            </Box>
            <FishermanList
              fishermen={state.fishermen}
              isLoading={state.isLoading}
              isDeleting={state.isDeleting}
              onDelete={actions.onDelete}
            />
          </Card>
        </Box>
      </Box>
    </Box>
  );
};
