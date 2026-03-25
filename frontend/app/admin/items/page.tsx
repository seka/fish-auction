'use client';

import { Suspense } from 'react';
import { ItemManagementContainer } from '@/src/features/admin';
import { Box, Text } from '@atoms';

export default function AdminItemsPage() {
  return (
    <Suspense
      fallback={
        <Box maxW="6xl" mx="auto" p="6" textAlign="center">
          <Text>Loading...</Text>
        </Box>
      }
    >
      <ItemManagementContainer />
    </Suspense>
  );
}
