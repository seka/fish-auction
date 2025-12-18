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
});
