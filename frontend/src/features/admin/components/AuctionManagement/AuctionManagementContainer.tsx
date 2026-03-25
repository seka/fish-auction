'use client';

import { useAuctionManagement } from '../../hooks/useAuctionManagement';
import { AuctionList } from './AuctionList';
import { AuctionForm } from './AuctionForm';
import { Box, Card, Text, HStack, Select, EmptyState } from '@atoms';
import { css } from 'styled-system/css';

export const AuctionManagementContainer = () => {
  const { state, form, actions, t } = useAuctionManagement();

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
        {t('Admin.Auctions.title')}
      </Text>

      {state.message && (
        <Box
          bg="blue.50"
          borderLeft="4px solid"
          borderColor="blue.500"
          className={css({ color: 'blue.700' })}
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
          <AuctionForm
            form={form}
            onSubmit={actions.onSubmit}
            onCancelEdit={actions.onCancelEdit}
            isCreating={state.isCreating}
            isUpdating={state.isUpdating}
            editingAuction={state.editingAuction}
            venues={state.venues}
          />
        </Box>

        {/* List Section */}
        <Box className={css({ md: { gridColumn: '2 / 3' } })}>
          <Card padding="none" overflow="hidden">
            <Box
              p="6"
              borderBottom="1px solid"
              borderColor="gray.200"
              bg="white"
              display="flex"
              justifyContent="space-between"
              alignItems="center"
              flexWrap="wrap"
              gap="4"
            >
              <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">
                {t('Admin.Auctions.list_title')}
              </Text>
              <HStack spacing="2">
                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>
                  {t('Admin.Auctions.filter_venue')}
                </Text>
                <Select
                  value={state.filterVenueId || ''}
                  onChange={(e) =>
                    actions.setFilterVenueId(e.target.value ? Number(e.target.value) : undefined)
                  }
                  className={css({ width: 'auto', py: '1' })}
                >
                  <option value="">{t('Admin.Auctions.filter_all')}</option>
                  {state.venues.map((venue) => (
                    <option key={venue.id} value={venue.id}>
                      {venue.name}
                    </option>
                  ))}
                </Select>
              </HStack>
            </Box>
            {state.isLoading ? (
              <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
                {t('Common.loading')}
              </Box>
            ) : state.auctions.length === 0 ? (
              <EmptyState
                message={t('Common.no_data')}
                icon={
                  <span role="img" aria-label="auction">
                    🏷️
                  </span>
                }
              />
            ) : (
              <AuctionList
                auctions={state.auctions}
                venues={state.venues}
                isLoading={state.isLoading}
                isUpdatingStatus={state.isUpdatingStatus}
                onEdit={actions.onEdit}
                onDelete={actions.onDelete}
                onStatusChange={actions.onStatusChange}
              />
            )}
          </Card>
        </Box>
      </Box>
    </Box>
  );
};
