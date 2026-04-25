'use client';

import { useState, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import { formatDateTimeForInput, toJSTISOString } from '@/src/utils/date';
import { getAuctionSchema, AuctionFormInput } from '@schemas/auction';
import { useAdminAuctions, useAdminAuctionMutations } from '../queries/useAuctions';
import { useAdminVenues } from '../queries/useVenues';
import { Auction } from '../types';
import { ApiError } from '@/src/core/api/client';

export const useAuctionManagement = () => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const [message, setMessage] = useState('');
  const [editingAuction, setEditingAuction] = useState<Auction | null>(null);
  const [filterVenueId, setFilterVenueId] = useState<number | undefined>(undefined);

  const { venues } = useAdminVenues();
  const { auctions, isLoading } = useAdminAuctions({ venueId: filterVenueId });
  const {
    createAuction,
    updateAuction,
    updateStatus,
    deleteAuction,
    isCreating,
    isUpdating,
    isUpdatingStatus,
    isDeleting,
  } = useAdminAuctionMutations();

  const schema = useMemo(() => getAuctionSchema(tValidation), [tValidation]);
  const form = useForm<AuctionFormInput>({
    resolver: zodResolver(schema),
  });

  const { reset, handleSubmit, setValue } = form;

  // Date オブジェクトを datetime-local 入力形式 "YYYY-MM-DDTHH:MM" (JST) に変換する
  const toDatetimeLocalString = (date: Date): string => formatDateTimeForInput(date);

  const onSubmit = async (data: AuctionFormInput) => {
    try {
      const payload = {
        ...data,
        venueId: Number(data.venueId),
        startAt: toJSTISOString(data.startAt),
        endAt: data.endAt ? toJSTISOString(data.endAt) : undefined,
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
      let errorMsg = t('Common.error_occurred');
      if (e instanceof ApiError) {
        if (e.status === 409) {
          errorMsg = t('Admin.Auctions.error_conflict');
        } else if (e.status === 500 || e.message === 'An internal error occurred') {
          errorMsg = t('Common.error_occurred');
        } else if (e.message) {
          errorMsg = e.message;
        }
      }
      setMessage(errorMsg);
    }
  };

  const onEdit = (auction: Auction) => {
    setEditingAuction(auction);
    setValue('venueId', auction.venueId);
    setValue('startAt', auction.duration.startAt ? toDatetimeLocalString(auction.duration.startAt) : '');
    setValue('endAt', auction.duration.endAt ? toDatetimeLocalString(auction.duration.endAt) : '');
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
      } catch {
        setMessage(t('Admin.Auctions.fail_delete'));
      }
    }
  };

  const onStatusChange = async (id: number, status: string) => {
    try {
      let startAt: string | undefined;
      if (status === 'in_progress') {
        const defaultValue = formatDateTimeForInput(new Date());
        const input = window.prompt('開始時刻を入力してください (YYYY-MM-DDTHH:mm)', defaultValue);
        if (input === null) {
          return;
        }
        if (!/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(input)) {
          setMessage('開始時刻の形式が不正です');
          return;
        }
        startAt = toJSTISOString(input);
      }

      await updateStatus({ id, status, startAt });
      setMessage(t('Admin.Auctions.success_status_update'));
    } catch {
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
    form,
    t,
  };
};
