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
    const { createBuyer, isCreating } = useBuyerMutation();

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

    return {
        state: {
            message,
            buyers,
            isLoading,
            isCreating,
        },
        form: {
            register,
            errors,
        },
        actions: {
            onSubmit: handleSubmit(onSubmit),
        },
        t,
    };
};
