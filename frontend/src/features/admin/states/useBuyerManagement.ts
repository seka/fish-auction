import { useState, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { getBuyerSchema, BuyerFormData } from '@schemas/admin';
import { useAdminBuyers, useAdminBuyerMutations } from '../queries/useBuyers';

export const useBuyerManagement = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const [message, setMessage] = useState('');

  const { buyers, isLoading } = useAdminBuyers();
  const { createBuyer, isCreating, deleteBuyer, isDeleting } = useAdminBuyerMutations();

  const schema = useMemo(() => getBuyerSchema(tValidation), [tValidation]);
  const form = useForm<BuyerFormData>({
    resolver: zodResolver(schema),
  });

  const { reset, handleSubmit } = form;

  const onSubmit = async (data: BuyerFormData) => {
    try {
      await createBuyer(data);
      setMessage(t('Admin.Buyers.success_register'));
      reset();
    } catch (e) {
      console.error(e);
      setMessage(t('Admin.Buyers.fail_register'));
    }
  };

  const onDelete = async (id: number) => {
    if (!window.confirm(t('Common.confirm_delete'))) return;
    try {
      await deleteBuyer(id);
      setMessage(t('Common.success_delete'));
    } catch (e) {
      console.error(e);
      setMessage(t('Common.error_occurred'));
    }
  };

  return {
    state: {
      message,
      buyers,
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
