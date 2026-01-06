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
                setMessage(t('Admin.Venues.success_update'));
                setEditingVenue(null);
            } else {
                await createVenue(data);
                setMessage(t('Admin.Venues.success_create'));
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage(t('Common.error_occurred'));
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
        if (confirm(t('Common.confirm_delete'))) {
            try {
                await deleteVenue(id);
                setMessage(t('Admin.Venues.success_delete'));
            } catch (e: any) {
                console.error(e);
                let errorMsg = t('Admin.Venues.fail_delete');
                if (e.name === 'ApiError') {
                    if (e.status === 409) {
                        errorMsg = t('Admin.Venues.error_delete_conflict');
                    } else if (e.status === 500 || e.message === 'An internal error occurred') {
                        errorMsg = t('Common.error_occurred');
                    } else if (e.message) {
                        errorMsg = e.message;
                    }
                }
                setMessage(errorMsg);
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
