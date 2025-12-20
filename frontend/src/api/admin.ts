import { apiClient } from '@/src/core/api/client';
import { RegisterItemParams, Fisherman, Buyer } from '@/src/models';

export const registerFisherman = async (name: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/admin/fishermen', { name });
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};

export const registerBuyer = async (name: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/admin/buyers', { name });
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
