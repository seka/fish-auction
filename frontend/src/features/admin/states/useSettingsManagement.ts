'use client';

import { useState, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { getPasswordComplexitySchema } from '@/src/models/schemas/fields/password';

export const useSettingsManagement = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const schema = useMemo(() => getPasswordComplexitySchema(tValidation), [tValidation]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);

    if (newPassword !== confirmPassword) {
      setMessage({ type: 'error', text: t('Validation.password_mismatch') });
      return;
    }

    const validation = schema.safeParse(newPassword);
    if (!validation.success) {
      setMessage({ type: 'error', text: validation.error.issues[0].message });
      return;
    }

    setIsLoading(true);

    try {
      const res = await fetch('/api/proxy/api/admin/password', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          current_password: currentPassword,
          new_password: newPassword,
        }),
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || t('Validation.password_update_failed'));
      }

      setMessage({ type: 'success', text: t('Validation.password_updated') });
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err: unknown) {
      setMessage({ type: 'error', text: err instanceof Error ? err.message : String(err) });
    } finally {
      setIsLoading(false);
    }
  };

  return {
    state: {
      currentPassword,
      newPassword,
      confirmPassword,
      message,
      isLoading,
    },
    actions: {
      setCurrentPassword,
      setNewPassword,
      setConfirmPassword,
      handleSubmit,
    },
    t,
  };
};
