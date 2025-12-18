import { render, screen, fireEvent, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionsPage from './page';
import { useAuctionPage } from './_hooks/useAuctionPage';

// Mock hook
vi.mock('./_hooks/useAuctionPage');
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

describe('AuctionsPage', () => {
    const mockOnSubmit = vi.fn((e) => e.preventDefault());
    const mockRegister = vi.fn();
    const mockOnStatusChange = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuctionPage as any).mockReturnValue({
            state: {
                venues: [{ id: 1, name: 'Venue 1' }],
                auctions: [
                    { id: 1, venueId: 1, auctionDate: '2023-01-01', status: 'scheduled' },
                    { id: 2, venueId: 1, auctionDate: '2023-01-02', status: 'in_progress' }
                ],
                isLoading: false,
                isCreating: false,
                message: '',
            },
            form: {
                register: mockRegister,
                errors: {},
            },
            actions: {
                onSubmit: mockOnSubmit,
                onStatusChange: mockOnStatusChange,
                setFilterVenueId: vi.fn(),
            },
            t: (key: string) => key,
        });
    });

    it('renders form and list', () => {
        render(<AuctionsPage />);
        expect(screen.getByText('Admin.Auctions.title')).toBeInTheDocument();
        // Venue 1 appears in the select dropdown and potentially in the list.
        expect(screen.getAllByText('Venue 1').length).toBeGreaterThan(0);
    });

    it('renders auctions with status badges', () => {
        render(<AuctionsPage />);
        // Use a function or regex to find text flexibly, or verify badge content structure
        // The badge might look different due to mockup status keys.
        // We mocked t to return key as is.
        // Status 'scheduled' -> t(AUCTION_STATUS_KEYS['scheduled']) -> 'scheduled' (if mocks work simply)
        // But let's check what's actually rendered.
        // If AUCTION_STATUS_KEYS['scheduled'] returns undefined, fallback is status string.

        // With mock t, we get the translation keys back.
        // scheduled -> 'AuctionStatus.scheduled'
        // in_progress -> 'AuctionStatus.in_progress'

        expect(screen.getByText('AuctionStatus.scheduled')).toBeInTheDocument();
        // in_progress has an icon, so use regex or flexible matcher
        expect(screen.getByText((content) => content.includes('AuctionStatus.in_progress'))).toBeInTheDocument();
    });

    it('calls status change action', () => {
        render(<AuctionsPage />);
        const startButton = screen.getByRole('button', { name: 'Admin.Auctions.start' });
        fireEvent.click(startButton);
        expect(mockOnStatusChange).toHaveBeenCalledWith(1, 'in_progress');
    });
});
