import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionDetailPage from './page';
import { ToastProvider } from '@/src/bootstraps/ToastProvider/useToast';

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

vi.mock('next/navigation', () => ({
  useParams: vi.fn(() => ({ id: '1' })),
  useRouter: vi.fn(() => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
  })),
  useSearchParams: vi.fn(() => ({
    get: vi.fn(),
  })),
}));

vi.mock('@/src/features/auctions/states/useAuctionDetail', () => ({
  useAuctionDetail: vi.fn(),
}));

vi.mock('@/src/api/buyer_auth', () => ({
  loginBuyer: vi.fn(),
}));

import { useAuctionDetail } from '@/src/features/auctions/states/useAuctionDetail';


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
      duration: {
        dateLabel: '2026-03-11',
        label: '10:00 ~ 12:00',
      },
    },
    items: [
      {
        id: 101,
        auctionId: 1,
        fishermanId: 1,
        fishType: 'Tuna',
        quantity: { value: 10, label: '10 kg' },
        unit: 'kg',
        price: { value: 3000, label: '¥3,000' },
        status: {
          value: 'Pending',
          labelKey: 'ItemStatus.Pending',
          variant: 'neutral',
          isPending: true,
          isBidding: false,
          isSold: false,
          isUnsold: false,
        },
        bidding: {
          highestBid: null,
          highestBidderId: null,
          highestBidderName: null,
          nextMinBid: { value: 1000, label: '¥1,000' },
        },
      },
      {
        id: 102,
        auctionId: 1,
        fishermanId: 2,
        fishType: 'Salmon',
        quantity: { value: 5, label: '5 kg' },
        unit: 'kg',
        price: { value: 5000, label: '¥5,000' },
        status: {
          value: 'Sold',
          labelKey: 'ItemStatus.Sold',
          variant: 'success',
          isPending: false,
          isBidding: false,
          isSold: true,
          isUnsold: false,
        },
        bidding: {
          highestBid: 5000,
          highestBidderId: 99,
          highestBidderName: 'テスト太郎',
          nextMinBid: { value: 5500, label: '¥5,500' },
        },
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
    vi.mocked(useAuctionDetail).mockReturnValue(
      defaultMockValue as unknown as ReturnType<typeof useAuctionDetail>,
    );
  });


  it('renders loading state', () => {
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      isLoading: true,
      auction: null,
    } as unknown as ReturnType<typeof useAuctionDetail>);
    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders login form if not logged in', () => {
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      isLoggedIn: false,
    } as unknown as ReturnType<typeof useAuctionDetail>);
    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );
    expect(screen.getByText('Public.AuctionDetail.login_title')).toBeInTheDocument();
  });

  it('renders auction details when logged in', () => {
    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );
    expect(screen.getByText('Public.AuctionDetail.auction_venue_title')).toBeInTheDocument();
    expect(screen.getByText('Tuna')).toBeInTheDocument();
    expect(screen.getByText('Salmon')).toBeInTheDocument();
  });

  it('selects an item and allows bidding', async () => {
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      selectedItem: defaultMockValue.items[0],
      selectedItemId: 101,
    } as unknown as ReturnType<typeof useAuctionDetail>);

    render(
      <ToastProvider>
        <AuctionDetailPage />
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
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      message: 'Success Message',
    } as unknown as ReturnType<typeof useAuctionDetail>);
    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );
    expect(screen.getByText('Success Message')).toBeInTheDocument();
  });

  it('renders checking state', () => {
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      isChecking: true,
      isLoggedIn: false,
    } as unknown as ReturnType<typeof useAuctionDetail>);
    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders no data when auction not found', () => {
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      auction: null,
      isLoading: false,
    } as unknown as ReturnType<typeof useAuctionDetail>);
    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );
    expect(screen.getByText('Common.no_data')).toBeInTheDocument();
  });

  it('handles login flow', async () => {
    vi.mocked(useAuctionDetail).mockReturnValue({
      ...defaultMockValue,
      isLoggedIn: false,
    } as unknown as ReturnType<typeof useAuctionDetail>);

    render(
      <ToastProvider>
        <AuctionDetailPage />
      </ToastProvider>,
    );

    const loginButton = screen.getByRole('button', { name: 'Public.Login.submit' });
    fireEvent.click(loginButton);

    expect(mockOnSubmitLogin).toHaveBeenCalled();
  });
});
