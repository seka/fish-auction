import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { venueSchema, VenueFormData } from '@/src/models/schemas/auction';
import { useVenueQuery } from '@/src/data/queries/adminVenue/useQuery';
import { useVenueMutation } from '@/src/data/queries/adminVenue/useMutation';
import { Venue } from '@/src/models/venue';

export const useVenueManagement = () => {
  const t = useTranslations();
  const [message, setMessage] = useState('');

  const { venues, isLoading } = useVenueQuery();
  const { createVenue, isCreating, deleteVenue, isDeleting, updateVenue, isUpdating } =
    useVenueMutation();

  const [editingVenue, setEditingVenue] = useState<Venue | null>(null);

  const form = useForm<VenueFormData>({
    resolver: zodResolver(venueSchema),
  });

  const { reset, handleSubmit, setValue } = form;

  const onSubmit = async (data: VenueFormData) => {
    try {
      if (editingVenue) {
        await updateVenue({ ...editingVenue, ...data });
        setMessage(t('Admin.Venues.success_update'));
        setEditingVenue(null);
      } else {
        await createVenue(data);
        setMessage(t('Admin.Venues.success_register'));
      }
      reset();
    } catch {
      setMessage(editingVenue ? t('Admin.Venues.fail_update') : t('Admin.Venues.fail_register'));
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
