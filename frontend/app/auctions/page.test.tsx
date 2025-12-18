import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionsListPage from './page';
import { useQuery } from '@tanstack/react-query';
import { usePublicVenues } from './_hooks/usePublicVenues';

// Mocks
vi.mock('@tanstack/react-query', () => ({
    useQuery: vi.fn(),
}));

vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
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
        (usePublicVenues as any).mockReturnValue({ venues: [{ id: 1, name: 'Venue A' }] });
    });

    it('renders loading state', () => {
        (useQuery as any).mockReturnValue({ data: undefined, isLoading: true });
        render(<AuctionsListPage />);
        expect(screen.getByText('読み込み中...')).toBeInTheDocument();
    });

    it('renders empty state when no auctions', () => {
        (useQuery as any).mockReturnValue({ data: [], isLoading: false });
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
            },
            {
                id: 2,
                status: 'scheduled',
                auctionDate: '2023-12-02',
                startTime: '10:00:00',
                endTime: '12:00:00',
                venueId: 1,
            },
        ];
        (useQuery as any).mockReturnValue({ data: mockAuctions, isLoading: false });

        render(<AuctionsListPage />);

        expect(screen.getAllByText(/Venue A/).length).toBeGreaterThan(0);
        expect(screen.getByText('🔥 AuctionStatus.in_progress')).toBeInTheDocument(); // Mock translation key
    });
    it('sorts auctions correctly (in_progress first, then date)', () => {
        const mockAuctions = [
            { id: 1, status: 'scheduled', auctionDate: '2023-12-05', startTime: '10:00', endTime: '12:00', venueId: 1 },
            { id: 2, status: 'in_progress', auctionDate: '2023-12-05', startTime: '10:00', endTime: '12:00', venueId: 1 }, // Should be first
            { id: 3, status: 'scheduled', auctionDate: '2023-12-01', startTime: '10:00', endTime: '12:00', venueId: 1 }, // Should be second (earlier than 2023-12-05)
        ];
        (useQuery as any).mockReturnValue({ data: mockAuctions, isLoading: false });
        render(<AuctionsListPage />);

        const auctionLinks = screen.getAllByRole('link');
        // Check order. The component creates a link for each auction in grid + back to top link.
        // Links to auction details have href `/auctions/${id}`.

        // Filter out "Back to top", assuming it doesn't match the same structure or verify text content inside cards.
        // Let's rely on finding text that appears in order.

        // However, `getAllByText` might return multiple elements if dates are same.
        // Best to find unique text or IDs.
        // The component renders status badges.
        // We know 'in_progress' renders with '🔥'.

        // Let's check if the first card contains '🔥'.
        const cards = screen.getAllByRole('link').filter(link => link.getAttribute('href')?.startsWith('/auctions/'));

        expect(cards[0]).toHaveTextContent('🔥'); // ID 2
        // Next should be ID 3 (2023-12-01) vs ID 1 (2023-12-05).
        // Sorting logic: if status not 'in_progress', sort by date+time ascending.
        // 2023-12-01 < 2023-12-05. So ID 3 should come before ID 1.


        // Check content.
        expect(cards[1]).toHaveTextContent('2023-12-01');
        expect(cards[2]).toHaveTextContent('2023-12-05');
    });

    it('resolves and displays venue name', () => {
        (usePublicVenues as any).mockReturnValue({ venues: [{ id: 99, name: 'Special Venue' }] });
        const mockAuctions = [
            { id: 1, status: 'scheduled', auctionDate: '2023-12-01', startTime: '10:00', endTime: '12:00', venueId: 99 },
        ];
        (useQuery as any).mockReturnValue({ data: mockAuctions, isLoading: false });
        render(<AuctionsListPage />);

        expect(screen.getByText('Special Venue')).toBeInTheDocument();
    });

    it('renders auctions with various statuses', () => {
        const mockAuctions = [
            { id: 1, status: 'cancelled', auctionDate: '2023-12-05', startTime: '10:00', endTime: '12:00', venueId: 1 },
            { id: 2, status: 'completed', auctionDate: '2023-12-04', startTime: '10:00', endTime: '12:00', venueId: 1 },
        ];
        (useQuery as any).mockReturnValue({ data: mockAuctions, isLoading: false });
        render(<AuctionsListPage />);

        expect(screen.getByText(/AuctionStatus.cancelled/)).toBeInTheDocument();
        expect(screen.getByText(/AuctionStatus.completed/)).toBeInTheDocument();
    });
});
