import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionDetailPage from './page';
import { useAuctionData } from './_hooks/useAuctionData';
import { useBidMutation } from './_hooks/useBidMutation';
import { useAuth } from '@/src/hooks/useAuth';
import { isAuctionActive } from '@/src/utils/auction';
import { ToastProvider } from '@/src/hooks/useToast';

// Mock dependencies
vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => {
    const t = (key: string) => (namespace ? `${namespace}.${key}` : key);
    t.raw = (key: string) => (namespace ? `${namespace}.${key}` : key);
    t.rich = (key: string, _values: Record<string, unknown>) => {
      const base = namespace ? `${namespace}.${key}` : key;
      return base;
    };
    return t;
  },
}));

vi.mock('react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('react')>();
  return {
    ...actual,
    use: (_: unknown) => {
      // Unwrapping promise for params
      // In test environment, we can just return a resolved object if passed as promise
      // or we might need to simulate it.
      // Simplified: Expecting params to be passed not as promise in test render if possible,
      // or handle it here.
      // Actually, in Next.js 15, params is a Promise. We need to handle this.
      // But 'use' hook unwrapping is internal to React.
      // Check if we can just pass an object that matches usage.
      // If 'promise' has 'then', it's a promise.
      // For testing, let's assume we pass { id: '1' } directly and mock 'use' to return it.
      return { id: '1' };
    },
  };
});

vi.mock('./_hooks/useAuctionData');
vi.mock('./_hooks/useBidMutation');
vi.mock('@/src/hooks/useAuth');
vi.mock('@/src/utils/auction');
vi.mock('@/src/api/buyer_auth', () => ({
  loginBuyer: vi.fn(),
}));

describe('AuctionDetailPage', () => {
  const mockRefetch = vi.fn();
  const mockSubmitBid = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(useAuctionData).mockReturnValue({
      auction: {
        id: 1,
        venueId: 1,
        auctionDate: '2026-03-11',
        startTime: '10:00:00',
        endTime: '12:00:00',
        status: 'in_progress',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      items: [
        {
          id: 101,
          auctionId: 1,
          fishermanId: 1,
          fishType: 'Tuna',
          quantity: 10,
          unit: 'kg',
          status: 'Pending',
          sortOrder: 1,
          createdAt: new Date().toISOString(),
        },
        {
          id: 102,
          auctionId: 1,
          fishermanId: 1,
          fishType: 'Salmon',
          quantity: 5,
          unit: 'kg',
          status: 'Sold',
          sortOrder: 2,
          createdAt: new Date().toISOString(),
        },
      ],
      isLoading: false,
      refetchItems: mockRefetch,
    });

    vi.mocked(useBidMutation).mockReturnValue({
      submitBid: mockSubmitBid,
      isLoading: false,
    });

    vi.mocked(useAuth).mockReturnValue({
      isLoggedIn: true,
      isChecking: false,
    });

    vi.mocked(isAuctionActive).mockReturnValue(true);
  });

  // Note: params prop in Next 15 is Promise.
  // We mocked 'react'.use so we just need to pass a Promise.
  const params = Promise.resolve({ id: '1' });

  it('renders loading state', () => {
    vi.mocked(useAuctionData).mockReturnValue({
      auction: undefined,
      items: [],
      isLoading: true,
      refetchItems: vi.fn(),
    });
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders login form if not logged in', () => {
    vi.mocked(useAuth).mockReturnValue({ isLoggedIn: false, isChecking: false });
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Public.AuctionDetail.login_title')).toBeInTheDocument();
  });

  it('renders auction details when logged in', () => {
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Public.AuctionDetail.auction_venue_title')).toBeInTheDocument();
    expect(screen.getByText('Tuna')).toBeInTheDocument();
    expect(screen.getByText('Salmon')).toBeInTheDocument();
  });

  it('selects an item and allows bidding', async () => {
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );

    // Select Tuna
    const tunaCard = screen.getByText('Tuna');
    fireEvent.click(tunaCard);

    // Check panel updates
    expect(screen.getByText('Public.AuctionDetail.selected_item')).toBeInTheDocument();

    // Bid amount input
    const bidInput = screen.getByRole('spinbutton');
    fireEvent.change(bidInput, { target: { value: '1000' } });

    // Submit
    mockSubmitBid.mockResolvedValue(true);
    const bidButton = screen.getByRole('button', { name: 'Public.AuctionDetail.bid_button' });
    fireEvent.click(bidButton);

    await waitFor(() => {
      expect(mockSubmitBid).toHaveBeenCalledWith(
        expect.objectContaining({ itemId: 101, price: 1000 }),
      );
    });
  });

  it('shows error message on failed bid', async () => {
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );

    fireEvent.click(screen.getByText('Tuna'));
    fireEvent.change(screen.getByRole('spinbutton'), { target: { value: '1000' } });

    mockSubmitBid.mockResolvedValue(false);
    fireEvent.click(screen.getByRole('button', { name: 'Public.AuctionDetail.bid_button' }));

    await waitFor(() => {
      expect(screen.getByText('Public.AuctionDetail.fail_bid')).toBeInTheDocument();
    });
  });
  it('renders checking state', () => {
    vi.mocked(useAuth).mockReturnValue({ isLoggedIn: false, isChecking: true });
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getAllByText('Common.loading').length).toBeGreaterThan(0);
  });

  it('renders no data when auction not found', () => {
    vi.mocked(useAuctionData).mockReturnValue({
      auction: undefined,
      items: [],
      isLoading: false,
      refetchItems: vi.fn(),
    });
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.no_data')).toBeInTheDocument();
  });

  it('renders item ended message for non-pending items', () => {
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    fireEvent.click(screen.getByText('Salmon')); // Status: Sold
    expect(screen.getByText('Public.AuctionDetail.item_ended')).toBeInTheDocument();
  });

  it('renders out of hours message when auction is inactive', () => {
    vi.mocked(isAuctionActive).mockReturnValue(false);
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    fireEvent.click(screen.getByText('Tuna')); // Status: Pending
    expect(screen.getByText('Public.AuctionDetail.out_of_hours_title')).toBeInTheDocument();
  });

  // Test login success separately?
  // onSubmitLogin calls loginBuyer and reloads.
  it('handles login flow', async () => {
    // Need to mock window.location.reload
    const originalLocation = window.location;
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { ...originalLocation, reload: vi.fn() },
    });

    vi.mocked(useAuth).mockReturnValue({ isLoggedIn: false, isChecking: false });
    const { loginBuyer } = await import('@/src/api/buyer_auth');
    vi.mocked(loginBuyer).mockResolvedValue({
      id: 1,
      name: 'Buyer',
    });

    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );

    fireEvent.change(screen.getByPlaceholderText('Common.email'), {
      target: { value: 'user@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Common.password'), {
      target: { value: 'password' },
    });

    fireEvent.submit(screen.getByRole('button', { name: 'Public.Login.submit' })); // Or click button in form

    await waitFor(() => {
      expect(loginBuyer).toHaveBeenCalled();
      expect(window.location.reload).toHaveBeenCalled();
    });

    // Cleanup
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation,
    });
  });
});
