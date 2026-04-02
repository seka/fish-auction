'use client';

import { ToastProvider as BaseToastProvider } from '@/src/bootstraps/ToastProvider/useToast';
import { ReactNode } from 'react';

export function ToastProvider({ children }: { children: ReactNode }) {
  return <BaseToastProvider>{children}</BaseToastProvider>;
}
