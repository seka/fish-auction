import { apiClient } from '@/src/core/api/client';

export const login = async (email: string, password: string): Promise<boolean> => {
  try {
    await apiClient.post('/api/login', { email, password });
    return true;
  } catch (e) {
    console.error('Login error:', e);
    return false;
  }
};
export const logout = async (): Promise<boolean> => {
  try {
    await apiClient.post('/api/admin/logout', {});
    return true;
  } catch (e) {
    console.error('Logout error:', e);
    return false;
  }
};
