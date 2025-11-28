import { apiClient } from '@/src/core/api/client';
import { AuctionItem, Bid } from '@/src/models';

export const getItems = async (status?: string): Promise<AuctionItem[]> => {
    const url = status ? `/api/items?status=${status}` : '/api/items';
    return apiClient.get<AuctionItem[]>(url);
};

export const submitBid = async (bid: Bid): Promise<boolean> => {
    try {
        await apiClient.post('/api/bid', bid);
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};
