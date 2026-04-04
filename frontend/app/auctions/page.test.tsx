import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionsListPage from './page';
import { usePublicAuctions } from '@/src/features/auctions/queries/usePublicAuctions';
import { useVenueQuery } from '@/src/data/queries/publicVenue/useQuery';
import { Auction } from '@/src/features/auctions';

// Mocks
vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

vi.mock('@/src/features/auctions/queries/usePublicAuctions', () => ({
  usePublicAuctions: vi.fn(),
}));

vi.mock('@/src/data/queries/publicVenue/useQuery', () => ({
  useVenueQuery: vi.fn(),
}));

describe('AuctionsListPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useVenueQuery).mockReturnValue({
      data: [{ id: 1, name: 'Venue A', createdAt: new Date().toISOString() }],
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useVenueQuery>);
  });

  it('renders loading state', () => {
    vi.mocked(usePublicAuctions).mockReturnValue({
      data: [],
      isLoading: true,
      error: null,
    } as unknown as ReturnType<typeof usePublicAuctions>);
    render(<AuctionsListPage />);
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders empty state when no auctions', () => {
    vi.mocked(usePublicAuctions).mockReturnValue({
      data: [],
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof usePublicAuctions>);
    render(<AuctionsListPage />);
    expect(screen.getByText('Public.Auctions.no_auctions')).toBeInTheDocument();
  });

  it('renders auctions list', () => {
    const mockAuctions: Auction[] = [
      {
        id: 1,
        status: {
          value: 'in_progress',
          labelKey: 'in_progress',
          variant: 'success',
          isScheduled: false,
          isInProgress: true,
          isCompleted: false,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-01T10:00:00+09:00'),
          endAt: new Date('2023-12-01T12:00:00+09:00'),
          dateLabel: '2023-12-01',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: true,
      },
      {
        id: 2,
        status: {
          value: 'scheduled',
          labelKey: 'scheduled',
          variant: 'info',
          isScheduled: true,
          isInProgress: false,
          isCompleted: false,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-02T10:00:00+09:00'),
          endAt: new Date('2023-12-02T12:00:00+09:00'),
          dateLabel: '2023-12-02',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: false,
      },
    ];
    vi.mocked(usePublicAuctions).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof usePublicAuctions>);

    render(<AuctionsListPage />);

    expect(screen.getAllByText(/Venue A/).length).toBeGreaterThan(0);
    expect(screen.getByText('AuctionStatus.in_progress')).toBeInTheDocument();
  });

  it('sorts auctions correctly (in_progress first, then date)', () => {
    const mockAuctions: Auction[] = [
      {
        id: 2,
        status: {
          value: 'in_progress',
          labelKey: 'in_progress',
          variant: 'success',
          isScheduled: false,
          isInProgress: true,
          isCompleted: false,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-05T10:00:00+09:00'),
          endAt: new Date('2023-12-05T12:00:00+09:00'),
          dateLabel: '2023-12-05',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: true,
      },
      {
        id: 3,
        status: {
          value: 'scheduled',
          labelKey: 'scheduled',
          variant: 'info',
          isScheduled: true,
          isInProgress: false,
          isCompleted: false,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-01T10:00:00+09:00'),
          endAt: new Date('2023-12-01T12:00:00+09:00'),
          dateLabel: '2023-12-01',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: false,
      },
      {
        id: 1,
        status: {
          value: 'scheduled',
          labelKey: 'scheduled',
          variant: 'info',
          isScheduled: true,
          isInProgress: false,
          isCompleted: false,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-05T10:00:00+09:00'),
          endAt: new Date('2023-12-05T12:00:00+09:00'),
          dateLabel: '2023-12-05',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: false,
      },
    ];
    vi.mocked(usePublicAuctions).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof usePublicAuctions>);
    render(<AuctionsListPage />);

    const cards = screen
      .getAllByRole('link')
      .filter((link) => link.getAttribute('href')?.startsWith('/auctions/'));

    expect(cards[0]).toHaveTextContent('AuctionStatus.in_progress'); // ID 2
    expect(cards[1]).toHaveTextContent('2023-12-01');
    expect(cards[2]).toHaveTextContent('2023-12-05');
  });

  it('resolves and displays venue name', () => {
    vi.mocked(useVenueQuery).mockReturnValue({
      data: [{ id: 99, name: 'Special Venue', createdAt: new Date().toISOString() }],
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useVenueQuery>);
    const mockAuctions: Auction[] = [
      {
        id: 1,
        status: {
          value: 'scheduled',
          labelKey: 'scheduled',
          variant: 'info',
          isScheduled: true,
          isInProgress: false,
          isCompleted: false,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-01T10:00:00+09:00'),
          endAt: new Date('2023-12-01T12:00:00+09:00'),
          dateLabel: '2023-12-01',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 99,
        isActive: false,
      },
    ];
    vi.mocked(usePublicAuctions).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof usePublicAuctions>);
    render(<AuctionsListPage />);
    expect(screen.getByText('Special Venue')).toBeInTheDocument();
  });

  it('renders auctions with various statuses', () => {
    const mockAuctions: Auction[] = [
      {
        id: 1,
        status: {
          value: 'cancelled',
          labelKey: 'cancelled',
          variant: 'error',
          isScheduled: false,
          isInProgress: false,
          isCompleted: false,
          isCancelled: true,
        },
        duration: {
          startAt: new Date('2023-12-05T10:00:00+09:00'),
          endAt: new Date('2023-12-05T12:00:00+09:00'),
          dateLabel: '2023-12-05',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: false,
      },
      {
        id: 2,
        status: {
          value: 'completed',
          labelKey: 'completed',
          variant: 'neutral',
          isScheduled: false,
          isInProgress: false,
          isCompleted: true,
          isCancelled: false,
        },
        duration: {
          startAt: new Date('2023-12-04T10:00:00+09:00'),
          endAt: new Date('2023-12-04T12:00:00+09:00'),
          dateLabel: '2023-12-04',
          startTime: '10:00:00',
          endTime: '12:00:00',
          label: '10:00 ~ 12:00',
        },
        venueId: 1,
        isActive: false,
      },
    ];
    vi.mocked(usePublicAuctions).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof usePublicAuctions>);
    render(<AuctionsListPage />);
    expect(screen.getByText(/AuctionStatus.cancelled/)).toBeInTheDocument();
    expect(screen.getByText(/AuctionStatus.completed/)).toBeInTheDocument();
  });
});
