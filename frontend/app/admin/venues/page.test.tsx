import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminVenuesPage from './page';
import { useVenuePage } from './_hooks/useVenuePage';

// Mock hook
vi.mock('./_hooks/useVenuePage');
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

describe('AdminVenuesPage', () => {
    const mockOnSubmit = vi.fn((e) => e.preventDefault());
    const mockRegister = vi.fn();
    const mockOnEdit = vi.fn();
    const mockOnDelete = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useVenuePage as any).mockReturnValue({
            state: {
                venues: [
                    { id: 1, name: 'Venue 1', location: 'Loc 1', description: 'Desc 1' },
                ],
                isLoading: false,
                isCreating: false,
                message: '',
                editingVenue: null,
            },
            form: {
                register: mockRegister,
                errors: {},
            },
            actions: {
                onSubmit: mockOnSubmit,
                onEdit: mockOnEdit,
                onDelete: mockOnDelete,
            },
            t: (key: string) => key,
        });
    });

    it('renders titles and form', () => {
        render(<AdminVenuesPage />);
        expect(screen.getByText('Admin.Venues.title')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Admin.Venues.placeholder_name')).toBeInTheDocument();
    });

    it('renders list of venues', () => {
        render(<AdminVenuesPage />);
        expect(screen.getByText('Venue 1')).toBeInTheDocument();
        expect(screen.getByText('Loc 1')).toBeInTheDocument();
    });

    it('calls edit and delete actions', () => {
        render(<AdminVenuesPage />);
        const editButton = screen.getByRole('button', { name: 'Common.edit' });
        const deleteButton = screen.getByRole('button', { name: 'Common.delete' });

        fireEvent.click(editButton);
        expect(mockOnEdit).toHaveBeenCalled();

        fireEvent.click(deleteButton);
        expect(mockOnDelete).toHaveBeenCalled();
    });
});
