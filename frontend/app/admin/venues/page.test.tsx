import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminVenuesPage from './page';
import { useVenueManagement } from '@/src/features/admin/states/useVenueManagement';

// Mock hook
vi.mock('@/src/features/admin/states/useVenueManagement');
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

describe('AdminVenuesPage', () => {
  const mockOnSubmit = vi.fn((e) => e.preventDefault());
  const mockRegister = vi.fn();
  const mockOnEdit = vi.fn();
  const mockOnDelete = vi.fn();
  const tMock = Object.assign((key: string) => key, {
    rich: vi.fn(),
    markup: vi.fn(),
    raw: vi.fn(),
    has: vi.fn(),
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useVenueManagement).mockReturnValue({
      state: {
        venues: [
          {
            id: 1,
            name: 'Venue 1',
            location: 'Loc 1',
            description: 'Desc 1',
            createdAt: new Date().toISOString(),
          },
        ],
        isLoading: false,
        isCreating: false,
        isUpdating: false,
        isDeleting: false,
        message: '',
        editingVenue: null,
      },
      form: {
        register: mockRegister,
        errors: {},
      } as unknown as ReturnType<typeof useVenueManagement>['form'],
      actions: {
        onSubmit: mockOnSubmit,
        onEdit: mockOnEdit,
        onCancelEdit: vi.fn(),
        onDelete: mockOnDelete,
      },
      t: tMock as unknown as ReturnType<typeof useVenueManagement>['t'],
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
