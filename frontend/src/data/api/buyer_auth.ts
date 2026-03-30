import { apiClient } from '@/src/core/api/client';
import { BuyerLoginFormData } from '@schema/buyer_auth';
import { Buyer } from '@entities';

export const loginBuyer = async (data: BuyerLoginFormData): Promise<Buyer | null> => {
  try {
    return await apiClient.post<Buyer>('/api/buyer/login', data);
  } catch (e) {
    console.error(e);
    return null;
  }
};

export const logoutBuyer = async (): Promise<boolean> => {
  try {
    await apiClient.post('/api/buyer/logout', {});
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
  } catch {
    // 未認証の場合は 401 エラーになるので null を返す
    return null;
  }
};
