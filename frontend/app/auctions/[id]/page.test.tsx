import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionDetailPage from './page';
import { useAuctionData } from './_hooks/useAuctionData';
import { useBidMutation } from './_hooks/useBidMutation';
import { useAuth } from './_hooks/useAuth';
import { isAuctionActive } from '@/src/utils/auction';

// Mock dependencies
vi.mock('next-intl', () => ({
    useTranslations: () => {
        const t = (key: string) => key;
        t.raw = (key: string) => key;
        return t;
    },
}));

vi.mock('react', async (importOriginal) => {
    const actual = await importOriginal<typeof import('react')>();
    return {
        ...actual,
        use: (promise: Promise<any>) => {
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
        }
    };
});

vi.mock('./_hooks/useAuctionData');
vi.mock('./_hooks/useBidMutation');
vi.mock('./_hooks/useAuth');
vi.mock('@/src/utils/auction');
vi.mock('@/src/api/buyer_auth', () => ({
    loginBuyer: vi.fn(),
}));

describe('AuctionDetailPage', () => {
    const mockRefetch = vi.fn();
    const mockSubmitBid = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();

        (useAuctionData as any).mockReturnValue({
            auction: { id: 1, status: 'in_progress', startTime: '10:00:00', endTime: '12:00:00' },
            items: [
                { id: 101, fishType: 'Tuna', quantity: 10, unit: 'kg', status: 'Pending' },
                { id: 102, fishType: 'Salmon', quantity: 5, unit: 'kg', status: 'Sold' }
            ],
            isLoading: false,
            refetchItems: mockRefetch,
        });

        (useBidMutation as any).mockReturnValue({
            submitBid: mockSubmitBid,
            isLoading: false,
        });

        (useAuth as any).mockReturnValue({
            isLoggedIn: true,
            isChecking: false,
        });

        (isAuctionActive as any).mockReturnValue(true);
    });

    // Note: params prop in Next 15 is Promise. 
    // We mocked 'react'.use so we just need to pass a Promise.
    const params = Promise.resolve({ id: '1' });

    it('renders loading state', () => {
        (useAuctionData as any).mockReturnValue({ isLoading: true });
        render(<AuctionDetailPage params={params} />);
        expect(screen.getByText('Common.loading')).toBeInTheDocument();
    });

    it('renders login form if not logged in', () => {
        (useAuth as any).mockReturnValue({ isLoggedIn: false, isChecking: false });
        render(<AuctionDetailPage params={params} />);
        expect(screen.getByText('Public.AuctionDetail.login_title')).toBeInTheDocument();
    });

    it('renders auction details when logged in', () => {
        render(<AuctionDetailPage params={params} />);
        expect(screen.getByText('Public.AuctionDetail.auction_venue_title')).toBeInTheDocument();
        expect(screen.getByText('Tuna')).toBeInTheDocument();
        expect(screen.getByText('Salmon')).toBeInTheDocument();
    });

    it('selects an item and allows bidding', async () => {
        render(<AuctionDetailPage params={params} />);

        // Select Tuna
        const tunaCard = screen.getByText('Tuna');
        fireEvent.click(tunaCard);

        // Check panel updates
        expect(screen.getByText('Public.AuctionDetail.selected_item')).toBeInTheDocument();

        // Bid amount input
        const bidInput = screen.getByPlaceholderText('0');
        fireEvent.change(bidInput, { target: { value: '1000' } });

        // Submit
        mockSubmitBid.mockResolvedValue(true);
        const bidButton = screen.getByRole('button', { name: 'Public.AuctionDetail.bid_button' });
        fireEvent.click(bidButton);

        await waitFor(() => {
            expect(mockSubmitBid).toHaveBeenCalledWith(expect.objectContaining({ itemId: 101, price: 1000 }));
        });
    });

    it('shows error message on failed bid', async () => {
        render(<AuctionDetailPage params={params} />);

        fireEvent.click(screen.getByText('Tuna'));
        fireEvent.change(screen.getByPlaceholderText('0'), { target: { value: '1000' } });

        mockSubmitBid.mockResolvedValue(false);
        fireEvent.click(screen.getByRole('button', { name: 'Public.AuctionDetail.bid_button' }));

        await waitFor(() => {
            expect(screen.getByText('Public.AuctionDetail.fail_bid')).toBeInTheDocument();
        });
    });
    it('renders checking state', () => {
        (useAuth as any).mockReturnValue({ isLoggedIn: false, isChecking: true });
        render(<AuctionDetailPage params={params} />);
        expect(screen.getAllByText('Common.loading').length).toBeGreaterThan(0);
    });

    it('renders no data when auction not found', () => {
        (useAuctionData as any).mockReturnValue({ auction: null, isLoading: false });
        render(<AuctionDetailPage params={params} />);
        expect(screen.getByText('Common.no_data')).toBeInTheDocument();
    });

    it('renders item ended message for non-pending items', () => {
        render(<AuctionDetailPage params={params} />);
        fireEvent.click(screen.getByText('Salmon')); // Status: Sold
        expect(screen.getByText('Public.AuctionDetail.item_ended')).toBeInTheDocument();
    });

    it('renders out of hours message when auction is inactive', () => {
        (isAuctionActive as any).mockReturnValue(false);
        render(<AuctionDetailPage params={params} />);
        fireEvent.click(screen.getByText('Tuna')); // Status: Pending
        expect(screen.getByText('Public.AuctionDetail.out_of_hours_title')).toBeInTheDocument();
    });

    // Test login success separately?
    // onSubmitLogin calls loginBuyer and reloads.
    it('handles login flow', async () => {
        // Need to mock window.location.reload
        const originalLocation = window.location;
        // @ts-ignore
        delete window.location;
        window.location = { ...originalLocation, reload: vi.fn() };

        (useAuth as any).mockReturnValue({ isLoggedIn: false, isChecking: false });
        const { loginBuyer } = await import('@/src/api/buyer_auth');
        (loginBuyer as any).mockResolvedValue({ id: 1 });

        render(<AuctionDetailPage params={params} />);

        fireEvent.change(screen.getByPlaceholderText('Common.email'), { target: { value: 'user@example.com' } });
        fireEvent.change(screen.getByPlaceholderText('Common.password'), { target: { value: 'password' } });

        fireEvent.submit(screen.getByRole('button', { name: 'Public.Login.submit' })); // Or click button in form

        await waitFor(() => {
            expect(loginBuyer).toHaveBeenCalled();
            expect(window.location.reload).toHaveBeenCalled();
        });

        // Cleanup
        window.location = originalLocation;
    });
});
