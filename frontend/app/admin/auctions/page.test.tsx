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
    it('handles filter change', () => {
        const mockSetFilterVenueId = vi.fn();
        (useAuctionPage as any).mockReturnValue({
            state: {
                venues: [{ id: 1, name: 'Venue 1' }, { id: 2, name: 'Venue 2' }],
                auctions: [],
                isLoading: false,
                isCreating: false,
                message: '',
                filterVenueId: undefined,
            },
            form: {
                register: mockRegister,
                errors: {},
            },
            actions: {
                onSubmit: mockOnSubmit,
                onStatusChange: mockOnStatusChange,
                setFilterVenueId: mockSetFilterVenueId,
            },
            t: (key: string) => key,
        });

        render(<AuctionsPage />);
        // Use combobox or find by display value. Since we mocked select to show options:
        // Select is usually implemented with native <select> in this project based on previous file reads?
        // Let's assume standard behavior or use getByRole('combobox').
        // Wait, the component file uses `Select` from `src/core/ui`.
        // Let's try finding the select by role or label.
        // Label is 'Admin.Auctions.filter_venue'

        // Note: The label text is used as label for select?
        // <Text as="label" ...>{t('Admin.Auctions.filter_venue')}</Text>
        // <Select ...>
        // Depending on accessibility, it might not be linked. Let's rely on display value 'Admin.Auctions.filter_all' which is the default option.

        const select = screen.getByDisplayValue('Admin.Auctions.filter_all');
        fireEvent.change(select, { target: { value: '2' } });
        expect(mockSetFilterVenueId).toHaveBeenCalledWith(2);
    });

    it('calls edit and delete actions', () => {
        const mockOnDelete = vi.fn();
        const mockOnEdit = vi.fn();

        (useAuctionPage as any).mockReturnValue({
            state: {
                venues: [{ id: 1, name: 'Venue 1' }],
                auctions: [
                    { id: 1, venueId: 1, auctionDate: '2023-01-01', status: 'scheduled' }
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
                onDelete: mockOnDelete,
                onEdit: mockOnEdit,
            },
            t: (key: string) => key,
        });

        render(<AuctionsPage />);

        const editButton = screen.getByText('Common.edit');
        fireEvent.click(editButton);
        expect(mockOnEdit).toHaveBeenCalled();

        const deleteButton = screen.getByText('Common.delete');
        fireEvent.click(deleteButton);
        expect(mockOnDelete).toHaveBeenCalledWith(1);
    });

    it('calls status change (finish)', () => {
        (useAuctionPage as any).mockReturnValue({
            state: {
                venues: [{ id: 1, name: 'Venue 1' }],
                auctions: [
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
                onDelete: vi.fn(),
                onEdit: vi.fn(),
            },
            t: (key: string) => key,
        });

        render(<AuctionsPage />);
        const finishButton = screen.getByText('Admin.Auctions.finish');
        fireEvent.click(finishButton);
        expect(mockOnStatusChange).toHaveBeenCalledWith(2, 'completed');
    });
});
