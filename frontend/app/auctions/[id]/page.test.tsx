import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionDetailPage from './page';
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
      return { id: '1' };
    },
  };
});

vi.mock('./_hooks/useAuctionDetailPage', () => ({
  useAuctionDetailPage: vi.fn(),
}));

vi.mock('@/src/api/buyer_auth', () => ({
  loginBuyer: vi.fn(),
}));

import { useAuctionDetailPage } from './_hooks/useAuctionDetailPage';

describe('AuctionDetailPage', () => {
  const mockOnSelectItem = vi.fn();
  const mockOnSubmitLogin = vi.fn();
  const mockOnSubmitBid = vi.fn();
  const t = (key: string) => key;
  t.rich = (key: string) => key;

  const defaultMockValue = {
    auction: {
      id: 1,
      venueId: 1,
      auctionDate: '2026-03-11',
      startTime: '10:00:00',
      endTime: '12:00:00',
      status: 'in_progress',
    },
    items: [
      {
        id: 101,
        fishType: 'Tuna',
        quantity: 10,
        unit: 'kg',
        status: 'Pending',
      },
      {
        id: 102,
        fishType: 'Salmon',
        quantity: 5,
        unit: 'kg',
        status: 'Sold',
      },
    ],
    isLoading: false,
    isChecking: false,
    isLoggedIn: true,
    selectedItem: null,
    selectedItemId: null,
    auctionActive: true,
    message: '',
    loginError: '',
    isBidLoading: false,
    bidForm: {
      register: vi.fn(),
      handleSubmit: (fn: (data: { price: string }) => void) => (e: React.BaseSyntheticEvent) => {
        e?.preventDefault();
        fn({ price: '1000' });
      },
      formState: { errors: {} },
    },
    loginForm: {
      register: vi.fn(),
      handleSubmit: (fn: (data: { email: string }) => void) => (e: React.BaseSyntheticEvent) => {
        e?.preventDefault();
        fn({ email: 'user@example.com' });
      },
      formState: { errors: {} },
    },
    onSelectItem: mockOnSelectItem,
    onSubmitLogin: mockOnSubmitLogin,
    onSubmitBid: mockOnSubmitBid,
    t,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useAuctionDetailPage).mockReturnValue(defaultMockValue as unknown as ReturnType<typeof useAuctionDetailPage>);
  });

  const params = Promise.resolve({ id: '1' });

  it('renders loading state', () => {
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      isLoading: true,
      auction: null,
    } as unknown as ReturnType<typeof useAuctionDetailPage>);
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders login form if not logged in', () => {
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      isLoggedIn: false,
    } as unknown as ReturnType<typeof useAuctionDetailPage>);
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
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      selectedItem: defaultMockValue.items[0],
      selectedItemId: 101,
    } as unknown as ReturnType<typeof useAuctionDetailPage>);

    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );

    // Bidding panel should show selected item
    expect(screen.getByText('Public.AuctionDetail.selected_item')).toBeInTheDocument();

    // Submit bid
    const bidButton = screen.getByRole('button', { name: 'Public.AuctionDetail.bid_button' });
    fireEvent.click(bidButton);

    expect(mockOnSubmitBid).toHaveBeenCalled();
  });

  it('shows error message if provided', () => {
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      message: 'Success Message',
    } as unknown as ReturnType<typeof useAuctionDetailPage>);
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Success Message')).toBeInTheDocument();
  });

  it('renders checking state', () => {
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      isChecking: true,
      isLoggedIn: false,
    } as unknown as ReturnType<typeof useAuctionDetailPage>);
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders no data when auction not found', () => {
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      auction: null,
      isLoading: false,
    } as unknown as ReturnType<typeof useAuctionDetailPage>);
    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.no_data')).toBeInTheDocument();
  });

  it('handles login flow', async () => {
    vi.mocked(useAuctionDetailPage).mockReturnValue({
      ...defaultMockValue,
      isLoggedIn: false,
    } as unknown as ReturnType<typeof useAuctionDetailPage>);

    render(
      <ToastProvider>
        <AuctionDetailPage params={params} />
      </ToastProvider>,
    );

    const loginButton = screen.getByRole('button', { name: 'Public.Login.submit' });
    fireEvent.click(loginButton);

    expect(mockOnSubmitLogin).toHaveBeenCalled();
  });
});
