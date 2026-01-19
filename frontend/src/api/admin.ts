import { apiClient } from '@/src/core/api/client';
import { RegisterItemParams, Fisherman, Buyer } from '@/src/models';
import { BuyerFormData } from '@/src/models/schemas/admin';

export const registerFisherman = async (name: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/admin/fishermen', { name });
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};

export const registerBuyer = async (data: BuyerFormData): Promise<boolean> => {
    try {
        await apiClient.post('/api/admin/buyers', data);
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};

export const registerItem = async (item: RegisterItemParams): Promise<boolean> => {
    try {
        await apiClient.post('/api/admin/items', item);
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};

export const getFishermen = async (): Promise<Fisherman[]> => {
    return apiClient.get<Fisherman[]>('/api/fishermen');
};

export const getBuyers = async (): Promise<Buyer[]> => {
    return apiClient.get<Buyer[]>('/api/buyers');
};

export const deleteFisherman = async (id: number): Promise<void> => {
    await apiClient.delete(`/api/admin/fishermen/${id}`);
};

export const deleteBuyer = async (id: number): Promise<void> => {
    await apiClient.delete(`/api/admin/buyers/${id}`);
};
