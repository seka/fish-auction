'use client';

import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import { Box, Card, Text, HStack, Select, EmptyState } from '@atoms';
import { Table, Thead, Tbody, Tr, Th } from '@molecules';
import { css } from 'styled-system/css';
import { useItemManagement } from '../../hooks/useItemManagement';
import { SortableRow } from './SortableRow';
import { ItemForm } from './ItemForm';

export const ItemManagementContainer = () => {
  const { state, form, actions, t } = useItemManagement();

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  return (
    <Box maxW="6xl" mx="auto" p="6">
      <Text
        as="h1"
        variant="h2"
        className={css({ color: 'gray.800' })}
        mb="8"
        pb="4"
        borderBottom="1px solid"
        borderColor="gray.200"
      >
        {t('Admin.Items.title')}
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

      <Box display="grid" gridTemplateColumns={{ base: '1fr', md: '1fr 2fr' }} gap="8">
        {/* Form Section */}
        <Box>
          <ItemForm
            form={form}
            onSubmit={actions.onSubmit}
            onCancelEdit={actions.onCancelEdit}
            isCreating={state.isCreating}
            isUpdating={state.isUpdating}
            editingItem={state.editingItem}
            auctions={state.auctions}
            fishermen={state.fishermen}
          />
        </Box>

        {/* List Section */}
        <Box>
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
                {t('Admin.Items.list_title')}
              </Text>
              <HStack spacing="2">
                <Text as="label" fontSize="sm" className={css({ color: 'gray.600' })}>
                  {t('Admin.Items.filter_auction')}
                </Text>
                <Select
                  value={state.filterAuctionId || ''}
                  onChange={(e) =>
                    actions.setFilterAuctionId(e.target.value ? Number(e.target.value) : undefined)
                  }
                  className={css({ width: 'auto', py: '1' })}
                >
                  <option value="">{t('Admin.Items.filter_all')}</option>
                  {state.auctions.map((auction) => (
                    <option key={auction.id} value={auction.id}>
                      {auction.auctionDate} (ID: {auction.id})
                    </option>
                  ))}
                </Select>
              </HStack>
            </Box>

            {state.isItemsLoading && state.items.length === 0 ? (
              <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
                {t('Common.loading')}
              </Box>
            ) : !state.filterAuctionId ? (
              <Box p="10" textAlign="center">
                <Text className={css({ color: 'gray.500' })}>
                  {t('Admin.Items.placeholder_select_auction')}
                </Text>
              </Box>
            ) : state.items.length === 0 ? (
              <EmptyState
                message={t('Common.no_data')}
                icon={
                  <span role="img" aria-label="item">
                    🐟
                  </span>
                }
              />
            ) : (
              <Box overflowX="auto">
                <DndContext
                  sensors={sensors}
                  collisionDetection={closestCenter}
                  onDragEnd={actions.onDragEnd}
                  modifiers={[restrictToVerticalAxis]}
                >
                  <Table>
                    <Thead>
                      <Tr>
                        <Th width="120px">{t('Admin.Items.sort_order')}</Th>
                        <Th>{t('Admin.Items.fish_type')}</Th>
                        <Th>{t('Admin.Items.fisherman')}</Th>
                        <Th>{t('Admin.Items.quantity')}</Th>
                        <Th className={css({ textAlign: 'right' })}>
                          {t('Admin.Auctions.action')}
                        </Th>
                      </Tr>
                    </Thead>
                    <Tbody>
                      <SortableContext
                        items={state.items.map((i) => i.id)}
                        strategy={verticalListSortingStrategy}
                      >
                        {state.items.map((item) => {
                          const fisherman = state.fishermen.find((f) => f.id === item.fishermanId);
                          return (
                            <SortableRow
                              key={item.id}
                              item={item}
                              fisherman={fisherman}
                              onEdit={actions.onEdit}
                              onDelete={actions.onDelete}
                              isDeleting={state.isDeleting}
                              t={t}
                            />
                          );
                        })}
                      </SortableContext>
                    </Tbody>
                  </Table>
                </DndContext>
              </Box>
            )}
          </Card>
        </Box>
      </Box>
    </Box>
  );
};
