'use client';

import { Box, Button, HStack, Text } from '@atoms';
import { Tr, Td } from '@molecules';
import { css } from 'styled-system/css';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { AuctionItem, Fisherman } from '@/src/models';

interface SortableRowProps {
  item: AuctionItem;
  fisherman: Fisherman | undefined;
  onEdit: (item: AuctionItem) => void;
  onDelete: (id: number) => void;
  isDeleting: boolean;
  t: (key: string) => string;
}

export const SortableRow = ({
  item,
  fisherman,
  onEdit,
  onDelete,
  isDeleting,
  t,
}: SortableRowProps) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: item.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    zIndex: isDragging ? 1 : 0,
    position: 'relative' as const,
    backgroundColor: isDragging ? '#f9fafb' : 'white',
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <Tr ref={setNodeRef} style={style}>
      <Td>
        <HStack spacing="2">
          <Box
            {...attributes}
            {...listeners}
            className={css({ cursor: 'grab', color: 'gray.400', _active: { cursor: 'grabbing' } })}
          >
            <span role="img" aria-label="drag">
              ⠿
            </span>
          </Box>
          <Text fontSize="sm" fontWeight="bold" className={css({ color: 'gray.500' })}>
            #{item.sortOrder}
          </Text>
        </HStack>
      </Td>
      <Td>
        <Text fontSize="sm" fontWeight="medium" className={css({ color: 'gray.900' })}>
          {item.fishType}
        </Text>
      </Td>
      <Td>
        <Text fontSize="sm" className={css({ color: 'gray.900' })}>
          {fisherman?.name || `ID: ${item.fishermanId}`}
        </Text>
      </Td>
      <Td>
        <Text fontSize="sm" className={css({ color: 'gray.900' })}>
          {item.quantity} {item.unit}
        </Text>
      </Td>
      <Td className={css({ textAlign: 'right' })}>
        <HStack justify="end" spacing="2">
          <Button size="sm" variant="outline" onClick={() => onEdit(item)}>
            {t('Common.edit')}
          </Button>
          <Button
            size="sm"
            className={css({ bg: 'red.50', color: 'red.600', _hover: { bg: 'red.100' } })}
            onClick={() => onDelete(item.id)}
            disabled={isDeleting}
          >
            {t('Common.delete')}
          </Button>
        </HStack>
      </Td>
    </Tr>
  );
};
