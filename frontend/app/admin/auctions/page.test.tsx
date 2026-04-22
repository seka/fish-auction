import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AuctionsPage from './page';
import { useAuctionManagement } from '@/src/features/admin/states/useAuctionManagement';
import { selectAuctionStatus, selectTimeLabel } from '@/src/features/admin/selectors/selectAuction';
import { Auction } from '@/src/features/admin/types/auction';
import { AuctionStatus } from '@/src/data/entities/auction';

// Helper to create mock Auction
const createMockAuction = (
  id: number,
  venueId: number,
  startAt: string | null,
  status: AuctionStatus,
  endAt: string | null = null,
): Auction => {
  const statusObj = selectAuctionStatus(status);
  const startDate = startAt ? new Date(startAt) : null;
  const endDate = endAt ? new Date(endAt) : null;
  return {
    id,
    venueId,
    status: statusObj,
    duration: {
      startAt: startDate,
      endAt: endDate,
      label: selectTimeLabel(startDate, endDate),
    },
    actions: {
      canStart: statusObj.isScheduled,
      canFinish: statusObj.isInProgress,
    },
    createdAt: '2023-01-01',
    updatedAt: '2023-01-01',
  };
};

// Mock hook
vi.mock('@/src/features/admin/states/useAuctionManagement');
vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

describe('AuctionsPage', () => {
  const mockOnSubmit = vi.fn((e) => e.preventDefault());
  const mockRegister = vi.fn();
  const mockOnStatusChange = vi.fn();
  const tMock = Object.assign((key: string) => key, {
    rich: vi.fn(),
    markup: vi.fn(),
    raw: vi.fn(),
    has: vi.fn(),
  });

  beforeEach(() => {
    vi.mocked(useAuctionManagement).mockReturnValue({
      state: {
        venues: [{ id: 1, name: 'Venue 1', createdAt: '2023-01-01' }],
        auctions: [
          createMockAuction(1, 1, null, 'scheduled'),
          createMockAuction(2, 1, '2023-01-02T10:00:00+09:00', 'in_progress', '2023-01-02T12:00:00+09:00'),
        ],
        isLoading: false,
        isCreating: false,
        isUpdating: false,
        isUpdatingStatus: false,
        isDeleting: false,
        message: '',
        editingAuction: null,
        filterVenueId: undefined,
      },
      form: {
        register: mockRegister,
        formState: { errors: {} },
      } as unknown as ReturnType<typeof useAuctionManagement>['form'],
      actions: {
        onSubmit: mockOnSubmit,
        onStatusChange: mockOnStatusChange,
        setFilterVenueId: vi.fn(),
        onEdit: vi.fn(),
        onCancelEdit: vi.fn(),
        onDelete: vi.fn(),
      },
      t: tMock as unknown as ReturnType<typeof useAuctionManagement>['t'],
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
    expect(screen.getByText('AuctionStatus.scheduled')).toBeInTheDocument();
    expect(
      screen.getByText((content) => content.includes('AuctionStatus.in_progress')),
    ).toBeInTheDocument();
  });

  it('calls status change action', () => {
    render(<AuctionsPage />);
    const startButton = screen.getByRole('button', { name: 'Admin.Auctions.start' });
    fireEvent.click(startButton);
    expect(mockOnStatusChange).toHaveBeenCalledWith(1, 'in_progress');
  });

  it('handles filter change', () => {
    const mockSetFilterVenueId = vi.fn();
    vi.mocked(useAuctionManagement).mockReturnValue({
      state: {
        venues: [
          { id: 1, name: 'Venue 1', createdAt: '2023-01-01' },
          { id: 2, name: 'Venue 2', createdAt: '2023-01-01' },
        ],
        auctions: [],
        isLoading: false,
        isCreating: false,
        isUpdating: false,
        isUpdatingStatus: false,
        isDeleting: false,
        message: '',
        editingAuction: null,
        filterVenueId: undefined,
      },
      form: {
        register: mockRegister,
        formState: { errors: {} },
      } as unknown as ReturnType<typeof useAuctionManagement>['form'],
      actions: {
        onSubmit: mockOnSubmit,
        onStatusChange: mockOnStatusChange,
        setFilterVenueId: mockSetFilterVenueId,
        onEdit: vi.fn(),
        onCancelEdit: vi.fn(),
        onDelete: vi.fn(),
      },
      t: tMock as unknown as ReturnType<typeof useAuctionManagement>['t'],
    });

    render(<AuctionsPage />);
    const select = screen.getByDisplayValue('Admin.Auctions.filter_all');
    fireEvent.change(select, { target: { value: '2' } });
    expect(mockSetFilterVenueId).toHaveBeenCalledWith(2);
  });

  it('calls edit and delete actions', () => {
    const mockOnDelete = vi.fn();
    const mockOnEdit = vi.fn();

    vi.mocked(useAuctionManagement).mockReturnValue({
      state: {
        venues: [{ id: 1, name: 'Venue 1', createdAt: '2023-01-01' }],
        auctions: [createMockAuction(1, 1, null, 'scheduled')],
        isLoading: false,
        isCreating: false,
        isUpdating: false,
        isUpdatingStatus: false,
        isDeleting: false,
        message: '',
        editingAuction: null,
        filterVenueId: undefined,
      },
      form: {
        register: mockRegister,
        formState: { errors: {} },
      } as unknown as ReturnType<typeof useAuctionManagement>['form'],
      actions: {
        onSubmit: mockOnSubmit,
        onStatusChange: mockOnStatusChange,
        setFilterVenueId: vi.fn(),
        onDelete: mockOnDelete,
        onEdit: mockOnEdit,
        onCancelEdit: vi.fn(),
      },
      t: tMock as unknown as ReturnType<typeof useAuctionManagement>['t'],
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
    vi.mocked(useAuctionManagement).mockReturnValue({
      state: {
        venues: [{ id: 1, name: 'Venue 1', createdAt: '2023-01-01' }],
        auctions: [createMockAuction(2, 1, '2023-01-02T10:00:00+09:00', 'in_progress', '2023-01-02T12:00:00+09:00')],
        isLoading: false,
        isCreating: false,
        isUpdating: false,
        isUpdatingStatus: false,
        isDeleting: false,
        message: '',
        editingAuction: null,
        filterVenueId: undefined,
      },
      form: {
        register: mockRegister,
        formState: { errors: {} },
      } as unknown as ReturnType<typeof useAuctionManagement>['form'],
      actions: {
        onSubmit: mockOnSubmit,
        onStatusChange: mockOnStatusChange,
        setFilterVenueId: vi.fn(),
        onDelete: vi.fn(),
        onEdit: vi.fn(),
        onCancelEdit: vi.fn(),
      },
      t: tMock as unknown as ReturnType<typeof useAuctionManagement>['t'],
    });

    render(<AuctionsPage />);
    const finishButton = screen.getByText('Admin.Auctions.finish');
    fireEvent.click(finishButton);
    expect(mockOnStatusChange).toHaveBeenCalledWith(2, 'completed');
  });
});
