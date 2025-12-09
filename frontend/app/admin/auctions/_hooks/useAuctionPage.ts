import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { auctionSchema, AuctionFormData } from '@/src/models/schemas/auction';
import { useAuctionQuery, useAuctionMutation } from '@/src/repositories/auction';
import { useVenues } from '../../venues/_hooks/useVenue';
import { Auction } from '@/src/models/auction';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';

export const useAuctionPage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');
    const [editingAuction, setEditingAuction] = useState<Auction | null>(null);
    const [filterVenueId, setFilterVenueId] = useState<number | undefined>(undefined);

    const { venues } = useVenues();
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
                setMessage('セリ情報を更新しました');
                setEditingAuction(null);
            } else {
                await createAuction(payload);
                setMessage('セリを作成しました');
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage('エラーが発生しました');
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
        if (confirm('本当に削除しますか？')) {
            try {
                await deleteAuction(id);
                setMessage('セリを削除しました');
            } catch (e) {
                console.error(e);
                setMessage('削除に失敗しました');
            }
        }
    };

    const onStatusChange = async (id: number, status: string) => {
        try {
            await updateStatus({ id, status });
            setMessage(`ステータスを ${status} に更新しました`);
        } catch (e) {
            console.error(e);
            setMessage('ステータス更新に失敗しました');
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
