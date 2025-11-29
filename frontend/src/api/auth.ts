import { apiClient } from '@/src/core/api/client';

export const login = async (password: string): Promise<boolean> => {
    try {
        await apiClient.post('/api/login', { password });
        return true;
    } catch (e) {
        console.error('Login error:', e);
        return false;
    }
};
