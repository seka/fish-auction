import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { fishermanSchema, FishermanFormData } from '@/src/models/schemas/admin';
import { useFishermanQuery, useFishermanMutation } from '@/src/repositories/fisherman';

export const useFishermanPage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');

    const { fishermen, isLoading } = useFishermanQuery();
    const { createFisherman, isCreating } = useFishermanMutation();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<FishermanFormData>({
        resolver: zodResolver(fishermanSchema),
    });

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

    return {
        state: {
            message,
            fishermen,
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
