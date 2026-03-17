import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionsListPage from './page';
import { useQuery, UseQueryResult } from '@tanstack/react-query';
import { usePublicVenues } from './_hooks/usePublicVenues';
import { Auction } from '@/src/models';

// Mocks
vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

vi.mock('./_hooks/usePublicVenues', () => ({
  usePublicVenues: vi.fn(),
}));

vi.mock('@/src/api/auction', () => ({
  getAuctions: vi.fn(),
}));

describe('AuctionsListPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(usePublicVenues).mockReturnValue({
      venues: [{ id: 1, name: 'Venue A', createdAt: new Date().toISOString() }],
    });
  });

  it('renders loading state', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: undefined,
      isLoading: true,
      isFetching: true,
      error: null,
      refetch: vi.fn(),
    } as unknown as UseQueryResult<Auction[], Error>);
    render(<AuctionsListPage />);
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
  });

  it('renders empty state when no auctions', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: [],
      isLoading: false,
      isFetching: false,
      error: null,
      refetch: vi.fn(),
    } as unknown as UseQueryResult<Auction[], Error>);
    render(<AuctionsListPage />);
    expect(screen.getByText('Public.Auctions.no_auctions')).toBeInTheDocument();
  });

  it('renders auctions list', () => {
    const mockAuctions = [
      {
        id: 1,
        status: 'in_progress',
        auctionDate: '2023-12-01',
        startTime: '10:00:00',
        endTime: '12:00:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: 2,
        status: 'scheduled',
        auctionDate: '2023-12-02',
        startTime: '10:00:00',
        endTime: '12:00:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ];
    vi.mocked(useQuery).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      isFetching: false,
      error: null,
      refetch: vi.fn(),
    } as unknown as UseQueryResult<Auction[], Error>);

    render(<AuctionsListPage />);

    expect(screen.getAllByText(/Venue A/).length).toBeGreaterThan(0);
    expect(screen.getByText('AuctionStatus.in_progress')).toBeInTheDocument(); // Mock translation key
  });

  it('sorts auctions correctly (in_progress first, then date)', () => {
    const mockAuctions = [
      {
        id: 1,
        status: 'scheduled',
        auctionDate: '2023-12-05',
        startTime: '10:00',
        endTime: '12:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: 2,
        status: 'in_progress',
        auctionDate: '2023-12-05',
        startTime: '10:00',
        endTime: '12:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      }, // Should be first
      {
        id: 3,
        status: 'scheduled',
        auctionDate: '2023-12-01',
        startTime: '10:00',
        endTime: '12:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      }, // Should be second (earlier than 2023-12-05)
    ];
    vi.mocked(useQuery).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      isFetching: false,
      error: null,
      refetch: vi.fn(),
    } as unknown as UseQueryResult<Auction[], Error>);
    render(<AuctionsListPage />);

    const cards = screen
      .getAllByRole('link')
      .filter((link) => link.getAttribute('href')?.startsWith('/auctions/'));

    expect(cards[0]).toHaveTextContent('AuctionStatus.in_progress'); // ID 2
    expect(cards[1]).toHaveTextContent('2023-12-01');
    expect(cards[2]).toHaveTextContent('2023-12-05');
  });

  it('resolves and displays venue name', () => {
    vi.mocked(usePublicVenues).mockReturnValue({
      venues: [{ id: 99, name: 'Special Venue', createdAt: new Date().toISOString() }],
    });
    const mockAuctions = [
      {
        id: 1,
        status: 'scheduled',
        auctionDate: '2023-12-01',
        startTime: '10:00',
        endTime: '12:00',
        venueId: 99,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ];
    vi.mocked(useQuery).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      isFetching: false,
      error: null,
      refetch: vi.fn(),
    } as unknown as UseQueryResult<Auction[], Error>);
    render(<AuctionsListPage />);
    expect(screen.getByText('Special Venue')).toBeInTheDocument();
  });

  it('renders auctions with various statuses', () => {
    const mockAuctions = [
      {
        id: 1,
        status: 'cancelled',
        auctionDate: '2023-12-05',
        startTime: '10:00',
        endTime: '12:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: 2,
        status: 'completed',
        auctionDate: '2023-12-04',
        startTime: '10:00',
        endTime: '12:00',
        venueId: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ];
    vi.mocked(useQuery).mockReturnValue({
      data: mockAuctions,
      isLoading: false,
      isFetching: false,
      error: null,
      refetch: vi.fn(),
    } as unknown as UseQueryResult<Auction[], Error>);
    render(<AuctionsListPage />);
    expect(screen.getByText(/AuctionStatus.cancelled/)).toBeInTheDocument();
    expect(screen.getByText(/AuctionStatus.completed/)).toBeInTheDocument();
  });
});
