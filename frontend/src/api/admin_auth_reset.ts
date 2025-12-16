import { apiClient } from '@/src/core/api/client';

export interface ResetPasswordRequest {
    email: string;
}

export interface ResetPasswordVerifyRequest {
    token: string;
}

export interface ResetPasswordConfirmRequest {
    token: string;
    new_password: string;
}

export const requestAdminPasswordReset = async (data: ResetPasswordRequest): Promise<void> => {
    await apiClient.post('/api/admin/password-reset/request', data);
};

export const verifyAdminResetToken = async (data: ResetPasswordVerifyRequest): Promise<void> => {
    await apiClient.post('/api/admin/password-reset/verify', data);
};

export const confirmAdminPasswordReset = async (data: ResetPasswordConfirmRequest): Promise<void> => {
    await apiClient.post('/api/admin/password-reset/confirm', data);
};
