import { apiClient } from '@/src/shared/api/client';
import { RegisterItemParams } from '@/src/shared/models';

export const registerFisherman = async (name: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/fishermen', { name });
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};

export const registerBuyer = async (name: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/buyers', { name });
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};

export const registerItem = async (item: RegisterItemParams): Promise<boolean> => {
    try {
        await apiClient.post('/api/items', item);
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};
