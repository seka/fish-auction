'use client';

import { useState, useRef, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { loginBuyer } from '@/src/data/api/buyer_auth';
import { getBidSchema, BidFormData } from '@/src/models/schemas/auction';
import { getBuyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { useAuctionDetailData, useBidSubmit } from '../queries/useAuctions';
import { useAuthQuery } from '@/src/data/queries/auth/useQuery';
import { getMinimumBidIncrement } from '@/src/utils/auction';
import { useTranslations } from 'next-intl';
import { useQueryClient } from '@tanstack/react-query';
import { authKeys } from '@/src/data/queries/auth/keys';
import { AuctionItem } from '@/src/models';

export const useAuctionDetail = (auctionId: number) => {
  const t = useTranslations();
  const tValidation = useTranslations('Validation');
  const queryClient = useQueryClient();
  const [selectedItemId, setSelectedItemId] = useState<number | null>(null);
  const [message, setMessage] = useState('');
  const [loginError, setLoginError] = useState('');
  const messageTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Clean up message timeout on unmount
  useEffect(() => {
    const currentTimeout = messageTimeoutRef.current;
    return () => {
      if (currentTimeout) {
        clearTimeout(currentTimeout);
      }
    };
  }, [messageTimeoutRef]);

  const { auction, items, isLoading: isDataLoading, refetchItems } = useAuctionDetailData(auctionId);
  const { submitBid, isLoading: isBidLoading } = useBidSubmit();
  const { isLoggedIn, isChecking } = useAuthQuery();
  const isLoading = isDataLoading;

  const bidForm = useForm<BidFormData>({
    resolver: zodResolver(getBidSchema(tValidation)),
  });

  const loginForm = useForm<BuyerLoginFormData>({
    resolver: zodResolver(getBuyerLoginSchema(tValidation)),
  });

  const selectedItem = items?.find((i: AuctionItem) => i.id === selectedItemId) || null;
  const auctionActive = auction?.isActive || false;

  const onSelectItem = (id: number) => {
    setSelectedItemId(id);
    bidForm.reset();
    setMessage('');
    if (messageTimeoutRef.current) {
      clearTimeout(messageTimeoutRef.current);
      messageTimeoutRef.current = null;
    }
  };

  const onSubmitLogin = async (data: BuyerLoginFormData) => {
    setLoginError('');
    const buyer = await loginBuyer(data);
    if (buyer) {
      await queryClient.invalidateQueries({ queryKey: authKeys.me() });
    } else {
      setLoginError(t('Public.Login.error_credentials'));
    }
  };

  const onSubmitBid = async (data: BidFormData) => {
    if (!selectedItem) return;

    const currentPrice = selectedItem.highestBid || 0;
    const minIncrement = getMinimumBidIncrement(currentPrice);
    const inputPrice = parseInt(data.price);

    if (inputPrice < currentPrice + minIncrement) {
      setMessage(
        t('Public.AuctionDetail.fail_bid_low_price', {
          min: (currentPrice + minIncrement).toLocaleString(),
        }),
      );
      return;
    }

    const success = await submitBid({
      itemId: selectedItem.id,
      buyerId: 0,
      price: inputPrice,
    });

    if (success) {
      setMessage(t('Public.AuctionDetail.success_bid', { item: selectedItem.fishType }));
      setSelectedItemId(null);
      bidForm.reset();
      refetchItems();

      // Clear existing timeout if any
      if (messageTimeoutRef.current) {
        clearTimeout(messageTimeoutRef.current);
      }

      messageTimeoutRef.current = setTimeout(() => {
        setMessage('');
        messageTimeoutRef.current = null;
      }, 3000);
    } else {
      setMessage(t('Public.AuctionDetail.fail_bid'));
    }
  };

  return {
    auction,
    items,
    isLoading,
    isChecking,
    isLoggedIn,
    selectedItem,
    selectedItemId,
    auctionActive,
    message,
    loginError,
    isBidLoading,
    bidForm,
    loginForm,
    onSelectItem,
    onSubmitLogin,
    onSubmitBid,
    t,
  };
};
