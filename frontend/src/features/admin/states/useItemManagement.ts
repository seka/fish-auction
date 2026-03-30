import { useState, useMemo } from 'react';
import { useSearchParams } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { DragEndEvent } from '@dnd-kit/core';
import { getItemSchema, ItemFormData } from '@schemas/admin';
import { useAdminItems, useAdminItemMutations } from '../queries/useItems';
import { useAdminFishermen } from '../queries/useFishermen';
import { useAdminAuctions } from '../queries/useAuctions';
import { AuctionItem } from '../types';

export const useItemManagement = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const [message, setMessage] = useState('');

  const searchParams = useSearchParams();
  const initialAuctionId = searchParams.get('auctionId');

  const [filterAuctionId, setFilterAuctionId] = useState<number | undefined>(
    initialAuctionId ? parseInt(initialAuctionId) : undefined,
  );
  const [prevInitialId, setPrevInitialId] = useState(initialAuctionId);
  const [editingItem, setEditingItem] = useState<AuctionItem | null>(null);

  const { fishermen } = useAdminFishermen();
  const { auctions } = useAdminAuctions({});
  const { data: items, isLoading: isItemsLoading } = useAdminItems(filterAuctionId);
  const {
    createItem,
    isCreating,
    deleteItem,
    isDeleting,
    updateItem,
    isUpdating,
    reorderItems,
    isSorting,
  } = useAdminItemMutations();

  const schema = useMemo(() => getItemSchema(tValidation), [tValidation]);
  const form = useForm<ItemFormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      auctionId: initialAuctionId || '',
    },
  });

  if (initialAuctionId !== prevInitialId) {
    setPrevInitialId(initialAuctionId);
    setFilterAuctionId(initialAuctionId ? parseInt(initialAuctionId) : undefined);
    if (initialAuctionId) {
      form.setValue('auctionId', initialAuctionId);
    }
  }

  const { reset, handleSubmit, setValue } = form;

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
    } catch {
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

  const onDragEnd = async (event: DragEndEvent) => {
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
          ids: newItemsOrder.map((item) => item.id),
        });
        setMessage(t('Admin.Items.success_sort'));
      } catch {
        setMessage(t('Admin.Items.success_sort'));
      }
    }
  };

  const onDelete = async (id: number) => {
    if (!window.confirm(t('Admin.Items.confirm_delete'))) return;
    try {
      await deleteItem(id);
      setMessage(t('Admin.Items.success_delete'));
    } catch {
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
    form,
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
