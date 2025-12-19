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

// 現在ログインしているバイヤー情報を取得（認証チェック）
export const getCurrentBuyer = async (): Promise<Buyer | null> => {
    try {
        return await apiClient.get<Buyer>('/api/buyer/me');
    } catch (e) {
        // 未認証の場合は 401 エラーになるので null を返す
        return null;
    }
};
