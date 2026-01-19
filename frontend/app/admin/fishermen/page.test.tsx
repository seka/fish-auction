import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminFishermenPage from './page';
import { useFishermanPage } from './_hooks/useFishermanPage';

// Mock hook
vi.mock('./_hooks/useFishermanPage');
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

describe('AdminFishermenPage', () => {
    const mockOnSubmit = vi.fn((e) => e.preventDefault());
    const mockRegister = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useFishermanPage as any).mockReturnValue({
            state: {
                fishermen: [
                    { id: 1, name: 'Fisherman 1' },
                    { id: 2, name: 'Fisherman 2' },
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
        render(<AdminFishermenPage />);
        expect(screen.getByText('Admin.Fishermen.title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Fishermen.register_title')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Admin.Fishermen.placeholder_name')).toBeInTheDocument();
    });

    it('renders list of fishermen', () => {
        render(<AdminFishermenPage />);
        expect(screen.getByText('Fisherman 1')).toBeInTheDocument();
        expect(screen.getByText('Fisherman 2')).toBeInTheDocument();
    });

    it('shows loading state', () => {
        (useFishermanPage as any).mockReturnValue({
            state: { fishermen: [], isLoading: true },
            form: { register: mockRegister, errors: {} },
            actions: { onSubmit: mockOnSubmit },
            t: (key: string) => key,
        });
        render(<AdminFishermenPage />);
        expect(screen.getAllByText('Common.loading').length).toBeGreaterThan(0);
    });

    it('shows empty state', () => {
        (useFishermanPage as any).mockReturnValue({
            state: { fishermen: [], isLoading: false },
            form: { register: mockRegister, errors: {} },
            actions: { onSubmit: mockOnSubmit },
            t: (key: string) => key,
        });
        render(<AdminFishermenPage />);
        expect(screen.getByText('Common.no_data')).toBeInTheDocument();
    });

    it('submits form', () => {
        render(<AdminFishermenPage />);
        const button = screen.getByRole('button', { name: 'Common.register' });
        fireEvent.click(button);
        expect(mockOnSubmit).toHaveBeenCalled();
    });

    it('calls delete action', () => {
        const mockOnDelete = vi.fn();
        (useFishermanPage as any).mockReturnValue({
            state: {
                fishermen: [{ id: 1, name: 'Fisherman 1' }],
                isLoading: false,
                isCreating: false,
                isDeleting: false,
            },
            form: { register: mockRegister, errors: {} },
            actions: { onSubmit: mockOnSubmit, onDelete: mockOnDelete },
            t: (key: string) => key,
        });
        render(<AdminFishermenPage />);
        fireEvent.click(screen.getAllByText('Common.delete')[0]);
        expect(mockOnDelete).toHaveBeenCalledWith(1);
    });
});
