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
            await createBuyer({ name: data.name });
            setMessage('中買人を登録しました');
            reset();
        } catch (e) {
            console.error(e);
            setMessage('登録に失敗しました');
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
