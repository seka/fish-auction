import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminBuyersPage from './page';
import { useBuyerPage } from './_hooks/useBuyerPage';

// Mock hook
vi.mock('./_hooks/useBuyerPage');
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

describe('AdminBuyersPage', () => {
    const mockOnSubmit = vi.fn((e) => e.preventDefault());
    const mockRegister = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useBuyerPage as any).mockReturnValue({
            state: {
                buyers: [
                    { id: 1, name: 'Buyer 1' },
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
                onDelete: vi.fn(),
            },
            t: (key: string) => key,
        });
    });

    it('renders titles and form', () => {
        render(<AdminBuyersPage />);
        expect(screen.getByText('Admin.Buyers.title')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Admin.Buyers.placeholder_name')).toBeInTheDocument();
    });

    it('renders list of buyers', () => {
        render(<AdminBuyersPage />);
        expect(screen.getByText('Buyer 1')).toBeInTheDocument();
    });

    it('calls submit action', () => {
        render(<AdminBuyersPage />);
        fireEvent.click(screen.getByRole('button', { name: 'Common.register' }));
        expect(mockOnSubmit).toHaveBeenCalled();
    });

    it('calls delete action', () => {
        const mockOnDelete = vi.fn();
        (useBuyerPage as any).mockReturnValue({
            state: {
                buyers: [{ id: 1, name: 'Buyer 1' }],
                isLoading: false,
                isCreating: false,
                isDeleting: false,
            },
            form: { register: mockRegister, errors: {} },
            actions: { onSubmit: mockOnSubmit, onDelete: mockOnDelete },
            t: (key: string) => key,
        });
        render(<AdminBuyersPage />);
        fireEvent.click(screen.getAllByText('Common.delete')[0]);
        expect(mockOnDelete).toHaveBeenCalledWith(1);
    });
});
