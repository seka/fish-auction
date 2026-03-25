import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { fishermanSchema, FishermanFormData } from '@/src/models/schemas/admin';
import { useFishermanQuery } from '@/src/data/queries/adminFisherman/useQuery';
import { useFishermanMutation } from '@/src/data/queries/adminFisherman/useMutation';

export const useFishermanManagement = () => {
  const t = useTranslations();
  const [message, setMessage] = useState('');

  const { fishermen, isLoading } = useFishermanQuery();
  const { createFisherman, isCreating, deleteFisherman, isDeleting } = useFishermanMutation();

  const form = useForm<FishermanFormData>({
    resolver: zodResolver(fishermanSchema),
  });

  const { reset, handleSubmit, formState: { errors } } = form;

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
