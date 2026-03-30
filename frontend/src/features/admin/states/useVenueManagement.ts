import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { getVenueSchema, VenueFormData } from '@/src/models/schemas/auction';
import { useAdminVenues, useAdminVenueMutations } from '../queries/useVenues';
import { Venue } from '@/src/models/venue';

export const useVenueManagement = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const [message, setMessage] = useState('');

  const { venues, isLoading } = useAdminVenues();
  const { createVenue, isCreating, deleteVenue, isDeleting, updateVenue, isUpdating } =
    useAdminVenueMutations();

  const [editingVenue, setEditingVenue] = useState<Venue | null>(null);

  const form = useForm<VenueFormData>({
    resolver: zodResolver(getVenueSchema(tValidation)),
  });

  const { reset, handleSubmit, setValue } = form;

  const onSubmit = async (data: VenueFormData) => {
    try {
      if (editingVenue) {
        await updateVenue({ id: editingVenue.id, data });
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
    if (!window.confirm(t('Common.confirm_delete'))) return;
    try {
      await deleteVenue(id);
      setMessage(t('Common.success_delete'));
    } catch {
      setMessage(t('Common.error_occurred'));
    }
  };

  return {
    state: {
      message,
      venues,
      isLoading,
      isCreating,
      isDeleting,
      isUpdating,
      editingVenue,
    },
    form,
    actions: {
      onSubmit: handleSubmit(onSubmit),
      onEdit,
      onCancelEdit,
      onDelete,
    },
    t,
  };
};
