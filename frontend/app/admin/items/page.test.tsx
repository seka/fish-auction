import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminItemsPage from './page';
import { useItemPage } from './_hooks/useItemPage';

// Mock hook
vi.mock('./_hooks/useItemPage');
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

describe('AdminItemsPage', () => {
  const mockOnSubmit = vi.fn((e) => e.preventDefault());
  const mockRegister = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useItemPage).mockReturnValue({
      state: {
        auctions: [
          {
            id: 1,
            auctionDate: '2023-01-01',
            venueId: 1,
            status: 'scheduled',
            createdAt: '2023-01-01',
            updatedAt: '2023-01-01',
          } as unknown,
        ],
        fishermen: [
          { id: 1, name: 'Fisherman 1', fishermanId: 'F001', createdAt: '2023-01-01' } as unknown,
        ],
        isCreating: false,
        isDeleting: false,
        isUpdating: false,
        isSorting: false,
        isItemsLoading: false,
        filterAuctionId: undefined,
        editingItem: null,
        items: [],
        message: '',
      },
      form: {
        register: mockRegister,
        errors: {},
      } as unknown,
      actions: {
        onSubmit: mockOnSubmit,
      } as unknown,
      t: ((key: string) => key) as unknown,
    } as unknown as ReturnType<typeof useItemPage>);
  });

  it('renders form elements', () => {
    render(<AdminItemsPage />);
    expect(screen.getByText('Admin.Items.title')).toBeInTheDocument();
    expect(screen.getByText('Admin.Items.auction')).toBeInTheDocument();
    expect(screen.getByText('Admin.Items.fisherman')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Admin.Items.placeholder_fish_type')).toBeInTheDocument();
  });

  it('calls submit action', () => {
    render(<AdminItemsPage />);
    fireEvent.click(screen.getByRole('button', { name: 'Common.register' }));
    expect(mockOnSubmit).toHaveBeenCalled();
  });
});
