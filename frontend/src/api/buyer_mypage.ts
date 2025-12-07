import { apiClient } from '@/src/core/api/client';

export interface Purchase {
    id: number;
    itemId: number;
    fishType: string;
    quantity: number;
    unit: string;
    price: number;
    buyerId: number;
    auctionId: number;
    auctionDate: string;
    createdAt: string;
}

export interface AuctionSummary {
    id: number;
    venueId: number;
    auctionDate: string;
    startTime: string | null;
    endTime: string | null;
    status: string;
    createdAt: string;
    updatedAt: string;
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
