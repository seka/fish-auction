import { apiClient } from '@/src/core/api/client';

export interface Purchase {
    id: number;
    item_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
    price: number;
    buyer_id: number;
    auction_id: number;
    auction_date: string;
    created_at: string;
}

export interface AuctionSummary {
    id: number;
    venue_id: number;
    auction_date: string;
    start_time: string | null;
    end_time: string | null;
    status: string;
    created_at: string;
    updated_at: string;
}

export const getMyPurchases = async (): Promise<Purchase[]> => {
    try {
        return await apiClient.get<Purchase[]>('/api/buyers/me/purchases');
    } catch (e) {
        console.error(e);
        return [];
    }
};

export const getMyAuctions = async (): Promise<AuctionSummary[]> => {
    try {
        return await apiClient.get<AuctionSummary[]>('/api/buyers/me/auctions');
    } catch (e) {
        console.error(e);
        return [];
    }
};
