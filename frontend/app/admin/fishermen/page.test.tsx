import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminFishermenPage from './page';
import { useFishermanManagement } from '@/src/features/admin/states/useFishermanManagement';

// Mock hook
vi.mock('@/src/features/admin/states/useFishermanManagement');
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

describe('AdminFishermenPage', () => {
  const mockOnSubmit = vi.fn((e) => e.preventDefault());
  const mockRegister = vi.fn();
  const tMock = Object.assign((key: string) => key, {
    rich: vi.fn(),
    markup: vi.fn(),
    raw: vi.fn(),
    has: vi.fn(),
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useFishermanManagement).mockReturnValue({
      state: {
        fishermen: [
          { id: 1, name: 'Fisherman 1' },
          { id: 2, name: 'Fisherman 2' },
        ],
        isLoading: false,
        isCreating: false,
        isDeleting: false,
        message: '',
      },
      form: {
        register: mockRegister,
        formState: { errors: { name: undefined } },
      } as unknown as ReturnType<typeof useFishermanManagement>['form'],
      actions: {
        onSubmit: mockOnSubmit,
        onDelete: vi.fn(),
      },
      t: tMock as unknown as ReturnType<typeof useFishermanManagement>['t'],
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
    vi.mocked(useFishermanManagement).mockReturnValue({
      state: { fishermen: [], isLoading: true, isCreating: false, isDeleting: false, message: '' },
      form: { register: mockRegister, formState: { errors: {} } } as unknown as ReturnType<
        typeof useFishermanManagement
      >['form'],
      actions: { onSubmit: mockOnSubmit, onDelete: vi.fn() },
      t: tMock as unknown as ReturnType<typeof useFishermanManagement>['t'],
    });
    render(<AdminFishermenPage />);
    expect(screen.getAllByText('Common.loading').length).toBeGreaterThan(0);
  });

  it('shows empty state', () => {
    vi.mocked(useFishermanManagement).mockReturnValue({
      state: { fishermen: [], isLoading: false, isCreating: false, isDeleting: false, message: '' },
      form: { register: mockRegister, formState: { errors: {} } } as unknown as ReturnType<
        typeof useFishermanManagement
      >['form'],
      actions: { onSubmit: mockOnSubmit, onDelete: vi.fn() },
      t: tMock as unknown as ReturnType<typeof useFishermanManagement>['t'],
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
    vi.mocked(useFishermanManagement).mockReturnValue({
      state: {
        fishermen: [{ id: 1, name: 'Fisherman 1' }],
        isLoading: false,
        isCreating: false,
        isDeleting: false,
        message: '',
      },
      form: { register: mockRegister, formState: { errors: {} } } as unknown as ReturnType<
        typeof useFishermanManagement
      >['form'],
      actions: { onSubmit: mockOnSubmit, onDelete: mockOnDelete },
      t: tMock as unknown as ReturnType<typeof useFishermanManagement>['t'],
    });
    render(<AdminFishermenPage />);
    fireEvent.click(screen.getAllByText('Common.delete')[0]);
    expect(mockOnDelete).toHaveBeenCalledWith(1);
  });
});
