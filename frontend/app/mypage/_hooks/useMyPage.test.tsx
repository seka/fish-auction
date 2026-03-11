import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useMyPage } from './useMyPage';
import { getMyPurchases, getMyAuctions } from '@/src/api/buyer_mypage';
import { logoutBuyer } from '@/src/api/buyer_auth';
import { useRouter } from 'next/navigation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';

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
    vi.mocked(useRouter).mockReturnValue({
      push: mockPush,
      replace: vi.fn(),
      back: vi.fn(),
      forward: vi.fn(),
      refresh: vi.fn(),
      prefetch: vi.fn(),
    });
    vi.mocked(getMyPurchases).mockResolvedValue([]);
    vi.mocked(getMyAuctions).mockResolvedValue([]);
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
    vi.mocked(logoutBuyer).mockResolvedValue(true);
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

    const mockEvent = {
        preventDefault: vi.fn(),
    } as unknown as React.FormEvent<HTMLFormElement>;

    await act(async () => {
      await result.current.handleUpdatePassword(mockEvent);
    });

    expect(result.current.message).toEqual({
      type: 'error',
      text: 'Validation.password_mismatch',
    });
  });

  it('handles successful password update', async () => {
    // Mock fetch success
    vi.mocked(global.fetch).mockResolvedValue({
      ok: true,
      json: async () => ({}),
    } as Response);

    const { result } = renderHook(() => useMyPage(), { wrapper });

    act(() => {
      result.current.setCurrentPassword('currentpass');
      result.current.setNewPassword('newpassword123');
      result.current.setConfirmPassword('newpassword123');
    });

    const mockEvent = {
        preventDefault: vi.fn(),
    } as unknown as React.FormEvent<HTMLFormElement>;

    await act(async () => {
      await result.current.handleUpdatePassword(mockEvent);
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/proxy/api/buyers/password',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({
          current_password: 'currentpass',
          new_password: 'newpassword123',
        }),
      }),
    );

    expect(result.current.message).toEqual({ type: 'success', text: 'パスワードを更新しました。' });
    expect(result.current.currentPassword).toBe('');
  });
});
