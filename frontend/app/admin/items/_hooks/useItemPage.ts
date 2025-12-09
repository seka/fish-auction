import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { itemSchema, ItemFormData } from '@/src/models/schemas/admin';
import { useItemMutation } from '@/src/repositories/item';
import { useFishermanQuery } from '@/src/repositories/fisherman';
import { useAuctionQuery } from '@/src/repositories/auction';

export const useItemPage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');

    const { fishermen } = useFishermanQuery();
    const { auctions } = useAuctionQuery({});
    const { createItem, isCreating } = useItemMutation();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<ItemFormData>({
        resolver: zodResolver(itemSchema),
    });

    const onSubmit = async (data: ItemFormData) => {
        try {
            await createItem({
                auctionId: parseInt(data.auctionId),
                fishermanId: parseInt(data.fishermanId),
                fishType: data.fishType,
                quantity: parseInt(data.quantity),
                unit: data.unit,
            });
            setMessage('出品を登録しました');
            reset();
        } catch (e) {
            console.error(e);
            setMessage('登録に失敗しました');
        }
    };

    return {
        state: {
            message,
            fishermen,
            auctions,
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
