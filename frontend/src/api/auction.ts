import { apiClient } from '@/src/core/api/client';
import { Auction, AuctionItem } from '@/src/models/auction';
import { AuctionFormData } from '@/src/models/schemas/auction';

export const createAuction = async (data: AuctionFormData): Promise<Auction> => {
    return apiClient.post<Auction>('/api/admin/auctions', data);
};

export const getAuctions = async (params?: { venueId?: number; date?: string; status?: string }): Promise<Auction[]> => {
    const query = new URLSearchParams();
    if (params?.venueId) query.append('venue_id', params.venueId.toString());
    if (params?.date) query.append('date', params.date);
    if (params?.status) query.append('status', params.status);

    const queryString = query.toString();
    const url = queryString ? `/api/auctions?${queryString}` : '/api/auctions';
    return apiClient.get<Auction[]>(url);
};

export const getAuction = async (id: number): Promise<Auction> => {
    return apiClient.get<Auction>(`/api/auctions/${id}`);
};

export const getAuctionItems = async (id: number): Promise<AuctionItem[]> => {
    return apiClient.get<AuctionItem[]>(`/api/auctions/${id}/items`);
};

export const updateAuction = async (id: number, data: AuctionFormData): Promise<void> => {
    return apiClient.put<void>(`/api/admin/auctions/${id}`, data);
};

export const updateAuctionStatus = async (id: number, status: string): Promise<void> => {
    return apiClient.patch<void>(`/api/admin/auctions/${id}/status`, { status });
};

export const deleteAuction = async (id: number): Promise<void> => {
    return apiClient.delete(`/api/admin/auctions/${id}`);
};
