import { apiClient } from '@/src/core/api/client';
import { Admin } from '@entities';

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

export const adminSessionCookie = 'admin_session';

// src/middleware から呼ばれる。エラーはスローし、呼び出し元で fail-open を判断する。
export const getAdminMe = async (cookie: CookieHeader): Promise<Admin> => {
  return apiClient.get<Admin>('/api/admin/me', { cookie });
};
