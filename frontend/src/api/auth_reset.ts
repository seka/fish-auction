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

export const requestPasswordReset = async (data: ResetPasswordRequest): Promise<void> => {
    await apiClient.post('/api/auth/password-reset/request', data);
};

export const verifyResetToken = async (data: ResetPasswordVerifyRequest): Promise<void> => {
    await apiClient.post('/api/auth/password-reset/verify', data);
};

export const confirmPasswordReset = async (data: ResetPasswordConfirmRequest): Promise<void> => {
    await apiClient.post('/api/auth/password-reset/confirm', data);
};
