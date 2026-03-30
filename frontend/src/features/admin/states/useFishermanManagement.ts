import { useState, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { getFishermanSchema, FishermanFormData } from '@schema/admin';
import { useAdminFishermen, useAdminFishermanMutations } from '../queries/useFishermen';

export const useFishermanManagement = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const [message, setMessage] = useState('');

  const { fishermen, isLoading } = useAdminFishermen();
  const { createFisherman, isCreating, deleteFisherman, isDeleting } = useAdminFishermanMutations();

  const schema = useMemo(() => getFishermanSchema(tValidation), [tValidation]);
  const form = useForm<FishermanFormData>({
    resolver: zodResolver(schema),
  });

  const { reset, handleSubmit } = form;

  const onSubmit = async (data: FishermanFormData) => {
    try {
      await createFisherman({ name: data.name });
      setMessage(t('Admin.Fishermen.success_register'));
      reset();
    } catch (e) {
      console.error(e);
      setMessage(t('Admin.Fishermen.fail_register'));
    }
  };

  const onDelete = async (id: number) => {
    if (!window.confirm(t('Common.confirm_delete'))) return;
    try {
      await deleteFisherman(id);
      setMessage(t('Common.success_delete'));
    } catch (e) {
      console.error(e);
      setMessage(t('Common.error_occurred'));
    }
  };

  return {
    state: {
      message,
      fishermen,
      isLoading,
      isCreating,
      isDeleting,
    },
    form,
    actions: {
      onSubmit: handleSubmit(onSubmit),
      onDelete,
    },
    t,
  };
};
