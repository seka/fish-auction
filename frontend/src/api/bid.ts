import { apiClient } from '@/src/core/api/client';
import { Bid } from '@/src/models';

export const submitBid = async (bid: Bid): Promise<boolean> => {
    try {
        await apiClient.post('/api/buyer/bids', bid);
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};
