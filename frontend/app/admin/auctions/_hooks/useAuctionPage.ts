import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { auctionSchema, AuctionFormData } from '@/src/models/schemas/auction';
import { useAuctionQuery, useAuctionMutation } from '@/src/repositories/auction';
import { useVenueQuery } from '@/src/repositories/venue';
import { Auction } from '@/src/models/auction';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';

export const useAuctionPage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');
    const [editingAuction, setEditingAuction] = useState<Auction | null>(null);
    const [filterVenueId, setFilterVenueId] = useState<number | undefined>(undefined);

    const { venues } = useVenueQuery();
    const { auctions, isLoading } = useAuctionQuery({ venueId: filterVenueId });
    const { createAuction, updateAuction, updateStatus, deleteAuction, isCreating, isUpdating, isUpdatingStatus, isDeleting } = useAuctionMutation();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<AuctionFormData>({
        resolver: zodResolver(auctionSchema),
    });

    const onSubmit = async (data: AuctionFormData) => {
        try {
            const payload = {
                ...data,
                venueId: Number(data.venueId),
            };

            if (editingAuction) {
                await updateAuction({ id: editingAuction.id, data: payload });
                setMessage(t('Admin.Auctions.success_update'));
                setEditingAuction(null);
            } else {
                await createAuction(payload);
                setMessage(t('Admin.Auctions.success_create'));
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage(t('Common.error_occurred'));
        }
    };

    const onEdit = (auction: Auction) => {
        setEditingAuction(auction);
        setValue('venueId', auction.venueId);
        setValue('auctionDate', auction.auctionDate);
        setValue('startTime', auction.startTime || '');
        setValue('endTime', auction.endTime || '');
        setValue('status', auction.status);
    };

    const onCancelEdit = () => {
        setEditingAuction(null);
        reset();
    };

    const onDelete = async (id: number) => {
        if (confirm(t('Common.confirm_delete'))) {
            try {
                await deleteAuction(id);
                setMessage(t('Admin.Auctions.success_delete'));
            } catch (e) {
                console.error(e);
                setMessage(t('Admin.Auctions.fail_delete'));
            }
        }
    };

    const onStatusChange = async (id: number, status: string) => {
        try {
            await updateStatus({ id, status });
            setMessage(t('Admin.Auctions.success_status_update'));
        } catch (e) {
            console.error(e);
            setMessage(t('Admin.Auctions.fail_status_update'));
        }
    };

    return {
        state: {
            message,
            editingAuction,
            filterVenueId,
            auctions,
            venues,
            isLoading,
            isCreating,
            isUpdating,
            isUpdatingStatus,
            isDeleting,
        },
        actions: {
            setFilterVenueId,
            onEdit,
            onCancelEdit,
            onDelete,
            onStatusChange,
            onSubmit: handleSubmit(onSubmit),
        },
        form: {
            register,
            errors,
        },
        t, // Helper to access translations if needed
    };
};
