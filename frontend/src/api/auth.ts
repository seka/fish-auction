import { apiClient } from '@/src/core/api/client';

export const login = async (password: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/auth/login', { password });
        return true;
    } catch (e) {
        console.error(e);
        return false;
    }
};
