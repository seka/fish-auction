'use client';

import { Box, Stack, Text } from '@atoms';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

interface Fisherman {
  id?: number;
  name: string;
}

interface FishermanListProps {
  fishermen: Fisherman[];
  isLoading: boolean;
  isDeleting: boolean;
  onDelete: (id: number) => void;
}

import { FishermanListItem } from './FishermanListItem';

export const FishermanList = ({
  fishermen,
  isLoading,
  isDeleting,
  onDelete,
}: FishermanListProps) => {
  const t = useTranslations();

  if (isLoading) {
    return (
      <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>
        {t('Common.loading')}
      </Box>
    );
  }

  if (fishermen.length === 0) {
    return (
      <Box p="12" textAlign="center">
        <Text color="muted">{t('Common.no_data')}</Text>
      </Box>
    );
  }

  return (
    <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
      {fishermen.map((fisherman) => (
        <FishermanListItem
          key={fisherman.id}
          fisherman={fisherman}
          isDeleting={isDeleting}
          onDelete={onDelete}
          t={t}
        />
      ))}
    </Stack>
  );
};
