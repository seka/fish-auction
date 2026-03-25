'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';

export const useSettingsManagement = () => {
  const t = useTranslations();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);

    if (newPassword !== confirmPassword) {
      setMessage({ type: 'error', text: t('Validation.password_mismatch') });
      return;
    }

    if (newPassword.length < 8) {
      setMessage({ type: 'error', text: t('Validation.password_too_short', { min: 8 }) });
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
