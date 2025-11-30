import { apiClient } from '@/src/core/api/client';
import { BuyerSignupFormData, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { Buyer } from '@/src/models';

export const signupBuyer = async (data: BuyerSignupFormData): Promise<Buyer> => {
    return await apiClient.post<Buyer>('/api/buyers', data);
};

export const loginBuyer = async (data: BuyerLoginFormData): Promise<Buyer | null> => {
    try {
        return await apiClient.post<Buyer>('/api/buyers/login', data);
    } catch (e) {
        console.error(e);
        return null;
    }
};

export const logoutBuyer = async (): Promise<boolean> => {
    try {
        await apiClient.post('/api/buyers/logout', {});
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};
