import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { itemSchema, ItemFormData } from '@/src/models/schemas/admin';
import { useItemMutation, useItemQuery } from '@/src/repositories/item';
import { useFishermanQuery } from '@/src/repositories/fisherman';
import { useAuctionQuery } from '@/src/repositories/auction';
import { AuctionItem } from '@/src/models';

export const useItemPage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');

    const searchParams = useSearchParams();
    const initialAuctionId = searchParams.get('auctionId');

    const [filterAuctionId, setFilterAuctionId] = useState<number | undefined>(
        initialAuctionId ? parseInt(initialAuctionId) : undefined
    );
    const [editingItem, setEditingItem] = useState<AuctionItem | null>(null);

    const { fishermen } = useFishermanQuery();
    const { auctions } = useAuctionQuery({});
    const { data: items, isLoading: isItemsLoading } = useItemQuery(filterAuctionId);
    const { createItem, isCreating, deleteItem, isDeleting, updateItem, isUpdating, reorderItems, isSorting } = useItemMutation();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<ItemFormData>({
        resolver: zodResolver(itemSchema),
        defaultValues: {
            auctionId: initialAuctionId || '',
        },
    });

    useEffect(() => {
        if (initialAuctionId) {
            setValue('auctionId', initialAuctionId);
            setFilterAuctionId(parseInt(initialAuctionId));
        }
    }, [initialAuctionId, setValue]);

    const onSubmit = async (data: ItemFormData) => {
        try {
            if (editingItem) {
                await updateItem({
                    id: editingItem.id,
                    auctionId: parseInt(data.auctionId),
                    fishermanId: parseInt(data.fishermanId),
                    fishType: data.fishType,
                    quantity: parseInt(data.quantity),
                    unit: data.unit,
                    status: editingItem.status,
                });
                setMessage(t('Admin.Items.success_update'));
                setEditingItem(null);
            } else {
                await createItem({
                    auctionId: parseInt(data.auctionId),
                    fishermanId: parseInt(data.fishermanId),
                    fishType: data.fishType,
                    quantity: parseInt(data.quantity),
                    unit: data.unit,
                });
                setMessage(t('Admin.Items.success_register'));
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage(editingItem ? t('Admin.Items.fail_update') : t('Admin.Items.fail_register'));
        }
    };

    const onEdit = (item: AuctionItem) => {
        setEditingItem(item);
        setValue('auctionId', String(item.auctionId));
        setValue('fishermanId', String(item.fishermanId));
        setValue('fishType', item.fishType);
        setValue('quantity', String(item.quantity));
        setValue('unit', item.unit);
    };

    const onCancelEdit = () => {
        setEditingItem(null);
        reset();
        // If we have a filterAuctionId, keep it in the form
        if (filterAuctionId) {
            setValue('auctionId', String(filterAuctionId));
        }
    };

    const onDragEnd = async (event: any) => {
        const { active, over } = event;
        if (!over || active.id === over.id) return;

        if (!items || !filterAuctionId) return;

        const oldIndex = items.findIndex((item) => item.id === active.id);
        const newIndex = items.findIndex((item) => item.id === over.id);

        if (oldIndex !== -1 && newIndex !== -1) {
            try {
                const newItemsOrder = [...items];
                const [movedItem] = newItemsOrder.splice(oldIndex, 1);
                newItemsOrder.splice(newIndex, 0, movedItem);

                await reorderItems({
                    auctionId: filterAuctionId,
                    ids: newItemsOrder.map(item => item.id),
                });
                setMessage(t('Admin.Items.success_sort'));
            } catch (e) {
                console.error(e);
                setMessage(t('Admin.Items.fail_sort'));
            }
        }
    };

    const onDelete = async (id: number) => {
        if (!window.confirm(t('Admin.Items.confirm_delete'))) return;
        try {
            await deleteItem(id);
            setMessage(t('Admin.Items.success_delete'));
        } catch (e) {
            console.error(e);
            setMessage(t('Admin.Items.fail_delete'));
        }
    };

    return {
        state: {
            message,
            fishermen,
            auctions,
            items: items || [],
            isCreating,
            isDeleting,
            isUpdating,
            isSorting,
            isItemsLoading,
            filterAuctionId,
            editingItem,
        },
        form: {
            register,
            errors,
            reset,
        },
        actions: {
            onSubmit: handleSubmit(onSubmit),
            onEdit,
            onCancelEdit,
            onDelete,
            onDragEnd,
            setFilterAuctionId,
        },
        t,
    };
};
