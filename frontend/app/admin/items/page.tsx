'use client';

import { Suspense } from 'react';
import { Spinner } from '@atoms';
import { ItemManagementContainer } from './components/ItemManagementContainer';

export default function AdminItemsPage() {
  return (
    <Suspense fallback={<Spinner />}>
      <ItemManagementContainer />
    </Suspense>
  );
}
