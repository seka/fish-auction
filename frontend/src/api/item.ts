import { apiClient } from '@/src/core/api/client';
import { AuctionItem } from '@/src/models';

export const getItems = async (status?: string): Promise<AuctionItem[]> => {
    const url = status ? `/api/items?status=${status}` : '/api/items';
    return apiClient.get<AuctionItem[]>(url);
};
