'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { loginBuyer } from '@/src/api/buyer_auth';
import { bidSchema, BidFormData } from '@/src/models/schemas/auction';
import { buyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { useAuctionData } from './useAuctionData';
import { useBidMutation } from './useBidMutation';
import { useAuthQuery } from '@/src/hooks/auth/useQuery';
import { isAuctionActive, getMinimumBidIncrement } from '@/src/utils/auction';
import { useTranslations } from 'next-intl';
import { useQueryClient } from '@tanstack/react-query';
import { authKeys } from '@/src/hooks/auth/keys';
import { AuctionItem } from '@/src/models';

export const useAuctionDetailPage = (auctionId: number) => {
  const t = useTranslations();
  const queryClient = useQueryClient();
  const [selectedItemId, setSelectedItemId] = useState<number | null>(null);
  const [message, setMessage] = useState('');
  const [loginError, setLoginError] = useState('');

  const { auction, items, isLoading, refetchItems } = useAuctionData(auctionId);
  const { submitBid, isLoading: isBidLoading } = useBidMutation();
  const { isLoggedIn, isChecking } = useAuthQuery();

  const bidForm = useForm<BidFormData>({
    resolver: zodResolver(bidSchema),
  });

  const loginForm = useForm<BuyerLoginFormData>({
    resolver: zodResolver(buyerLoginSchema),
  });

  const selectedItem = items?.find((i: AuctionItem) => i.id === selectedItemId) || null;
  const auctionActive = auction ? isAuctionActive(auction) : false;

  const onSelectItem = (id: number) => {
    setSelectedItemId(id);
    bidForm.reset();
    setMessage('');
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
      setTimeout(() => setMessage(''), 3000);
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
