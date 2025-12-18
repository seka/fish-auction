import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useMyPage } from './useMyPage';
import { getMyPurchases, getMyAuctions } from '@/src/api/buyer_mypage';
import { logoutBuyer } from '@/src/api/buyer_auth';
import { useRouter } from 'next/navigation';

// Mocks
vi.mock('next/navigation', () => ({
    useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

vi.mock('@/src/api/buyer_mypage', () => ({
    getMyPurchases: vi.fn(),
    getMyAuctions: vi.fn(),
}));

vi.mock('@/src/api/buyer_auth', () => ({
    logoutBuyer: vi.fn(),
}));

// Fetch mock needs to be global
global.fetch = vi.fn();

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);

describe('useMyPage', () => {
    const mockPush = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useRouter as any).mockReturnValue({ push: mockPush });
        (getMyPurchases as any).mockResolvedValue([]);
        (getMyAuctions as any).mockResolvedValue([]);
        queryClient.clear();
    });

    it('returns initial state', async () => {
        const { result } = renderHook(() => useMyPage(), { wrapper });

        expect(result.current.activeTab).toBe('purchases');
        expect(result.current.currentPassword).toBe('');
        await waitFor(() => {
            expect(getMyPurchases).toHaveBeenCalled();
        });
    });

    it('handles logout', async () => {
        (logoutBuyer as any).mockResolvedValue(true);
        const { result } = renderHook(() => useMyPage(), { wrapper });

        await act(async () => {
            await result.current.handleLogout();
        });

        expect(logoutBuyer).toHaveBeenCalled();
        expect(mockPush).toHaveBeenCalledWith('/login/buyer');
    });

    it('handles password update validation', async () => {
        const { result } = renderHook(() => useMyPage(), { wrapper });

        // Mismatch scenario
        act(() => {
            result.current.setNewPassword('password123');
            result.current.setConfirmPassword('password124');
        });

        await act(async () => {
            await result.current.handleUpdatePassword({ preventDefault: vi.fn() } as any);
        });

        expect(result.current.message).toEqual({ type: 'error', text: '新しいパスワードが一致しません。' });

        // Short password scenario
        act(() => {
            result.current.setNewPassword('short');
            result.current.setConfirmPassword('short');
        });

        await act(async () => {
            await result.current.handleUpdatePassword({ preventDefault: vi.fn() } as any);
        });

        expect(result.current.message).toEqual({ type: 'error', text: 'パスワードは8文字以上である必要があります。' });
    });

    it('handles successful password update', async () => {
        // Mock fetch success
        (global.fetch as any).mockResolvedValue({
            ok: true,
            json: async () => ({}),
        });

        const { result } = renderHook(() => useMyPage(), { wrapper });

        act(() => {
            result.current.setCurrentPassword('currentpass');
            result.current.setNewPassword('newpassword123');
            result.current.setConfirmPassword('newpassword123');
        });

        await act(async () => {
            await result.current.handleUpdatePassword({ preventDefault: vi.fn() } as any);
        });

        expect(global.fetch).toHaveBeenCalledWith('/api/proxy/api/buyers/password', expect.objectContaining({
            method: 'PUT',
            body: JSON.stringify({
                current_password: 'currentpass',
                new_password: 'newpassword123',
            }),
        }));

        expect(result.current.message).toEqual({ type: 'success', text: 'パスワードを更新しました。' });
        expect(result.current.currentPassword).toBe('');
    });
});
