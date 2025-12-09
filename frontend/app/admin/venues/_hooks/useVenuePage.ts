import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { venueSchema, VenueFormData } from '@/src/models/schemas/auction';
import { useVenueQuery, useVenueMutation } from '@/src/repositories/venue';
import { Venue } from '@/src/models/venue';

export const useVenuePage = () => {
    const t = useTranslations();
    const [message, setMessage] = useState('');
    const [editingVenue, setEditingVenue] = useState<Venue | null>(null);

    const { venues, isLoading } = useVenueQuery();
    const { createVenue, updateVenue, deleteVenue, isCreating, isUpdating, isDeleting } = useVenueMutation();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<VenueFormData>({
        resolver: zodResolver(venueSchema),
    });

    const onSubmit = async (data: VenueFormData) => {
        try {
            if (editingVenue) {
                await updateVenue({ id: editingVenue.id, data });
                setMessage('会場を更新しました');
                setEditingVenue(null);
            } else {
                await createVenue(data);
                setMessage('会場を作成しました');
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage('エラーが発生しました');
        }
    };

    const onEdit = (venue: Venue) => {
        setEditingVenue(venue);
        setValue('name', venue.name);
        setValue('location', venue.location || '');
        setValue('description', venue.description || '');
    };

    const onCancelEdit = () => {
        setEditingVenue(null);
        reset();
    };

    const onDelete = async (id: number) => {
        if (confirm('本当に削除しますか？')) {
            try {
                await deleteVenue(id);
                setMessage('会場を削除しました');
            } catch (e) {
                console.error(e);
                setMessage('削除に失敗しました');
            }
        }
    };

    return {
        state: {
            message,
            venues,
            isLoading,
            editingVenue,
            isCreating,
            isUpdating,
            isDeleting,
        },
        form: {
            register,
            errors,
        },
        actions: {
            onSubmit: handleSubmit(onSubmit),
            onEdit,
            onCancelEdit,
            onDelete,
        },
        t,
    };
};
