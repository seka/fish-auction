import { render, screen, fireEvent } from '@testing-library/react';
import { PublicNavbar } from './PublicNavbar';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { usePathname, useRouter } from 'next/navigation';
import { useQuery, useQueryClient, UseQueryResult } from '@tanstack/react-query';
import * as buyerAuth from '@/src/data/api/buyer_auth';
import { Buyer } from '@entities';

// Mocks
vi.mock('next/navigation', () => ({
  usePathname: vi.fn(),
  useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  useQueryClient: vi.fn(),
}));

vi.mock('@/src/api/buyer_auth', () => ({
  getCurrentBuyer: vi.fn(),
  logoutBuyer: vi.fn(),
}));

describe('PublicNavbar', () => {
  const mockPush = vi.fn();
  const mockSetQueryData = vi.fn();

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
    vi.mocked(useQueryClient).mockReturnValue({
      setQueryData: mockSetQueryData,
      getQueryData: vi.fn(),
      invalidateQueries: vi.fn(),
      clear: vi.fn(),
    } as unknown as ReturnType<typeof useQueryClient>);
    vi.mocked(usePathname).mockReturnValue('/');
  });

  const mockQuerySuccess = (data: Buyer | null): UseQueryResult<Buyer | null, Error> =>
    ({
      data,
      isLoading: false,
      isFetching: false,
      isSuccess: true,
      isError: false,
      error: null,
      status: 'success',
      refetch: vi.fn(),
      isPending: false,
      isPlaceholderData: false,
      isRefetching: false,
      isStale: false,
      dataUpdatedAt: 0,
      errorUpdatedAt: 0,
      failureCount: 0,
      failureReason: null,
      isFetched: true,
      isFetchedAfterMount: true,
      isInitialLoading: false,
      isLoadingError: false,
      isPaused: false,
      isRefetchError: false,
      fetchStatus: 'idle',
      promise: Promise.resolve(data),
    }) as unknown as UseQueryResult<Buyer | null, Error>;

  it('renders correctly when not logged in', () => {
    vi.mocked(useQuery).mockReturnValue(mockQuerySuccess(null));

    render(<PublicNavbar />);

    expect(screen.getByText('Common.app_name')).toBeInTheDocument();
    expect(screen.getByText('Navbar.active_auctions')).toBeInTheDocument();
    expect(screen.getByText('Navbar.login')).toBeInTheDocument();
    expect(screen.queryByText('Navbar.logout')).not.toBeInTheDocument();
    expect(screen.queryByText('Navbar.mypage')).not.toBeInTheDocument();
  });

  it('renders correctly when logged in', () => {
    vi.mocked(useQuery).mockReturnValue(mockQuerySuccess({ id: 1, name: 'Buyer' }));

    render(<PublicNavbar />);

    expect(screen.getByText('Navbar.logout')).toBeInTheDocument();
    expect(screen.getByText('Navbar.mypage')).toBeInTheDocument();
    expect(screen.queryByText('Navbar.login')).not.toBeInTheDocument();
  });

  it('handles logout correctly', async () => {
    vi.mocked(useQuery).mockReturnValue(mockQuerySuccess({ id: 1, name: 'Buyer' }));
    const mockLogoutBuyer = vi.spyOn(buyerAuth, 'logoutBuyer').mockResolvedValue(true);

    render(<PublicNavbar />);

    const logoutButton = screen.getByText('Navbar.logout');
    fireEvent.click(logoutButton);

    expect(mockLogoutBuyer).toHaveBeenCalled();
  });

  it('does not render on admin pages', () => {
    vi.mocked(usePathname).mockReturnValue('/admin/dashboard');
    const { container } = render(<PublicNavbar />);
    expect(container).toBeEmptyDOMElement();
  });
});
