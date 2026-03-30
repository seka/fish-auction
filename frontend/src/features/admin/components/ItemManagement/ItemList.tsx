'use client';

import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import { Box, EmptyState, Text } from '@atoms';
import { Table, Thead, Tbody, Tr, Th } from '@molecules';
import { css } from 'styled-system/css';
import { AuctionItem, Fisherman } from '../../types';
import { ItemListCell } from './ItemListCell';

interface ItemListProps {
  items: AuctionItem[];
  fishermen: Fisherman[];
  isItemsLoading: boolean;
  filterAuctionId?: number;
  isDeleting: boolean;
  onDragEnd: (event: DragEndEvent) => void;
  onEdit: (item: AuctionItem) => void;
  onDelete: (id: number) => void;
  t: (key: string, values?: Record<string, string | number>) => string;
}

export const ItemList = ({
  items,
  fishermen,
  isItemsLoading,
  filterAuctionId,
  isDeleting,
  onDragEnd,
  onEdit,
  onDelete,
  t,
}: ItemListProps) => {
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

  if (isItemsLoading && items.length === 0) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  if (!filterAuctionId) {
    return (
      <Box p="10" textAlign="center">
        <Text className={css({ color: 'gray.500' })}>
          {t('Admin.Items.placeholder_select_auction')}
        </Text>
      </Box>
    );
  }

  if (items.length === 0) {
    return (
      <EmptyState
        message={t('Common.no_data')}
        icon={
          <span role="img" aria-label="item">
            🐟
          </span>
        }
      />
    );
  }

  return (
    <Box overflowX="auto">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={onDragEnd}
        modifiers={[restrictToVerticalAxis]}
      >
        <Table>
          <Thead>
            <Tr>
              <Th width="120px">{t('Admin.Items.sort_order')}</Th>
              <Th>{t('Admin.Items.fish_type')}</Th>
              <Th>{t('Admin.Items.fisherman')}</Th>
              <Th>{t('Admin.Items.quantity')}</Th>
              <Th className={css({ textAlign: 'right' })}>{t('Admin.Auctions.action')}</Th>
            </Tr>
          </Thead>
          <Tbody>
            <SortableContext items={items.map((i) => i.id)} strategy={verticalListSortingStrategy}>
              {items.map((item) => {
                const fisherman = fishermen.find((f) => f.id === item.fishermanId);
                return (
                  <ItemListCell
                    key={item.id}
                    item={item}
                    fisherman={fisherman}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    isDeleting={isDeleting}
                    t={t}
                  />
                );
              })}
            </SortableContext>
          </Tbody>
        </Table>
      </DndContext>
    </Box>
  );
};
