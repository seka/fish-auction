import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuctionDetail } from './useAuctionDetail';
import { loginBuyer } from '@/src/data/api/buyer_auth';
import { useAuctionDetailData, useBidSubmit } from '../queries/useAuctions';
import { useAuthQuery } from '@/src/data/queries/auth/useQuery';
import { useQueryClient } from '@tanstack/react-query';
import { authKeys } from '@/src/data/queries/auth/keys';

// Mocks
vi.mock('@/src/data/api/buyer_auth', () => ({
  loginBuyer: vi.fn(),
}));

vi.mock('../queries/useAuctions', () => ({
  useAuctionDetailData: vi.fn(),
  useBidSubmit: vi.fn(),
}));

vi.mock('@/src/data/queries/auth/useQuery', () => ({
  useAuthQuery: vi.fn(),
}));

vi.mock('@tanstack/react-query', () => ({
  useQueryClient: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

describe('useAuctionDetail', () => {
  const mockInvalidateQueries = vi.fn();
  const mockQueryClient = { invalidateQueries: mockInvalidateQueries };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useQueryClient).mockReturnValue(
      mockQueryClient as unknown as ReturnType<typeof useQueryClient>,
    );
    vi.mocked(useAuctionDetailData).mockReturnValue({
      auction: { id: 1, isActive: true } as unknown as ReturnType<
        typeof useAuctionDetailData
      >['auction'],
      items: [],
      isLoading: false,
      refetchItems: vi.fn(),
    });
    vi.mocked(useBidSubmit).mockReturnValue({
      submitBid: vi.fn(),
      isLoading: false,
    });
    vi.mocked(useAuthQuery).mockReturnValue({
      isLoggedIn: false,
      isChecking: false,
    } as unknown as ReturnType<typeof useAuthQuery>);
  });

  it('calls loginBuyer and invalidates authKeys.me() on successful login', async () => {
    vi.mocked(loginBuyer).mockResolvedValue({
      id: 1,
      email: 'test@example.com',
    } as unknown as Awaited<ReturnType<typeof loginBuyer>>);

    const { result } = renderHook(() => useAuctionDetail(1));

    await act(async () => {
      await result.current.onSubmitLogin({
        email: 'test@example.com',
        password: 'password',
      });
    });

    expect(loginBuyer).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password',
    });

    expect(mockInvalidateQueries).toHaveBeenCalledWith({
      queryKey: authKeys.me(),
    });
  });

  it('sets login error on failed login', async () => {
    vi.mocked(loginBuyer).mockResolvedValue(
      null as unknown as Awaited<ReturnType<typeof loginBuyer>>,
    );

    const { result } = renderHook(() => useAuctionDetail(1));

    await act(async () => {
      await result.current.onSubmitLogin({
        email: 'test@example.com',
        password: 'wrong-password',
      });
    });

    expect(result.current.loginError).toBe('Public.Login.error_credentials');
    expect(mockInvalidateQueries).not.toHaveBeenCalled();
  });
});
