import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useMyPage } from './useMyPage';
import { useParticipatingAuctions } from '@/src/data/queries/buyerAuction/useQuery';
import { useMyPurchases } from '@/src/data/queries/buyerPurchase/useQuery';
import { logoutBuyer } from '@/src/data/api/buyer_auth';
import { useRouter } from 'next/navigation';
import { useQueryClient } from '@tanstack/react-query';
import { AppRouterInstance } from 'next/dist/shared/lib/app-router-context.shared-runtime';

// Mocks
vi.mock('next/navigation', () => ({
  useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock('@/src/data/queries/buyerAuction/useQuery', () => ({
  useParticipatingAuctions: vi.fn(),
}));

vi.mock('@/src/data/queries/buyerPurchase/useQuery', () => ({
  useMyPurchases: vi.fn(),
}));

vi.mock('@tanstack/react-query', () => ({
  useQueryClient: vi.fn(),
}));

vi.mock('@/src/data/api/buyer_auth', () => ({
  logoutBuyer: vi.fn(),
}));

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
    } as unknown as AppRouterInstance);
    vi.mocked(useMyPurchases).mockReturnValue({ purchases: [], isLoading: false, error: null });
    vi.mocked(useParticipatingAuctions).mockReturnValue({
      auctions: [],
      isLoading: false,
      error: null,
    });
    vi.mocked(useQueryClient).mockReturnValue({
      invalidateQueries: vi.fn(),
    } as unknown as ReturnType<typeof useQueryClient>);
  });

  it('returns initial state', async () => {
    const { result } = renderHook(() => useMyPage());

    expect(result.current.activeTab).toBe('purchases');
    expect(result.current.passwordState.currentPassword).toBe('');
    expect(useMyPurchases).toHaveBeenCalled();
  });

  it('handles logout', async () => {
    vi.mocked(logoutBuyer).mockResolvedValue(true);
    const { result } = renderHook(() => useMyPage());

    await act(async () => {
      await result.current.handleLogout();
    });

    expect(logoutBuyer).toHaveBeenCalled();
    expect(mockPush).toHaveBeenCalledWith('/login/buyer');
  });

  it('handles password update validation', async () => {
    const { result } = renderHook(() => useMyPage());

    // Mismatch scenario
    act(() => {
      result.current.passwordState.setNewPassword('password123');
      result.current.passwordState.setConfirmPassword('password124');
    });

    const mockEvent = {
      preventDefault: vi.fn(),
    } as unknown as React.FormEvent<HTMLFormElement>;

    await act(async () => {
      await result.current.passwordState.handleUpdatePassword(mockEvent);
    });

    expect(result.current.passwordState.passwordMessage).toEqual({
      type: 'error',
      text: 'MyPage.Settings.password_mismatch',
    });
  });

  it('handles successful password update', async () => {
    const { result } = renderHook(() => useMyPage());

    act(() => {
      result.current.passwordState.setCurrentPassword('currentpass');
      result.current.passwordState.setNewPassword('newpassword123');
      result.current.passwordState.setConfirmPassword('newpassword123');
    });

    const mockEvent = {
      preventDefault: vi.fn(),
    } as unknown as React.FormEvent<HTMLFormElement>;

    await act(async () => {
      await result.current.passwordState.handleUpdatePassword(mockEvent);
    });

    expect(result.current.passwordState.passwordMessage).toEqual({
      type: 'success',
      text: 'MyPage.Settings.password_updated',
    });
    expect(result.current.passwordState.currentPassword).toBe('');
  });
});
