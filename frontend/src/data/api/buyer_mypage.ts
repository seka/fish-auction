import { apiClient } from '@/src/core/api/client';
import { Purchase, AuctionSummary } from '@entities';

export const getMyPurchases = async (): Promise<Purchase[]> => {
  try {
    return await apiClient.get<Purchase[]>('/api/buyer/me/purchases');
  } catch (e) {
    console.error(e);
    return [];
  }
};

export const getMyAuctions = async (): Promise<AuctionSummary[]> => {
  try {
    return await apiClient.get<AuctionSummary[]>('/api/buyer/me/auctions');
  } catch (e) {
    console.error(e);
    return [];
  }
};
