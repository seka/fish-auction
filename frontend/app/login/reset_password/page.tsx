'use client';

import { Suspense } from 'react';
import { Spinner } from '@atoms';
import { PublicResetPasswordContainer } from './components/PublicResetPasswordContainer';

export default function ResetPasswordPage() {
  return (
    <Suspense fallback={<Spinner />}>
      <PublicResetPasswordContainer />
    </Suspense>
  );
}
