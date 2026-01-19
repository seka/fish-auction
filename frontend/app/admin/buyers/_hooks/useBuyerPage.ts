import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { buyerSchema, BuyerFormData } from '@/src/models/schemas/admin';
import { useBuyerQuery, useBuyerMutation } from '@/src/repositories/buyer';

export const useBuyerPage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');

    const { buyers, isLoading } = useBuyerQuery();
    const { createBuyer, isCreating, deleteBuyer, isDeleting } = useBuyerMutation();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<BuyerFormData>({
        resolver: zodResolver(buyerSchema),
    });

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
        form: {
            register,
            errors,
        },
        actions: {
            onSubmit: handleSubmit(onSubmit),
            onDelete,
        },
        t,
    };
};
