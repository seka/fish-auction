'use client';

import { Suspense } from 'react';
import { Spinner } from '@atoms';
import { AdminResetPasswordContainer } from './components/AdminResetPasswordContainer';

export default function AdminResetPasswordPage() {
  return (
    <Suspense fallback={<Spinner />}>
      <AdminResetPasswordContainer />
    </Suspense>
  );
}
