import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import InvoicePage from './page';
import { useInvoiceManagement } from '@/src/features/admin/states/useInvoiceManagement';
import { useTranslations } from 'next-intl';

// Mock hook
vi.mock('@/src/features/admin/states/useInvoiceManagement');
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

describe('InvoicePage', () => {
  const mockSetSelectedInvoice = vi.fn();
  const tMock = Object.assign((key: string) => key, {
    rich: vi.fn(),
    markup: vi.fn(),
    raw: vi.fn(),
    has: vi.fn(),
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useInvoiceManagement).mockReturnValue({
      state: {
        invoices: [{ buyerId: 1, buyerName: 'Buyer 1', totalAmount: 1000 }],
        isLoading: false,
        selectedInvoice: null,
      },
      actions: {
        setSelectedInvoice: mockSetSelectedInvoice,
      },
      t: tMock as ReturnType<typeof useTranslations>,
    });
  });

  it('renders invoice list', () => {
    render(<InvoicePage />);
    expect(screen.getByText('Admin.Invoice.title')).toBeInTheDocument();
    expect(screen.getByText('Buyer 1')).toBeInTheDocument();
    expect(screen.getByText('¥1,000')).toBeInTheDocument();
  });

  it('opens detail modal on row click', () => {
    render(<InvoicePage />);
    const row = screen.getByText('Buyer 1').closest('tr');
    fireEvent.click(row!);
    expect(mockSetSelectedInvoice).toHaveBeenCalled();
  });

  it('renders detail modal when selected', () => {
    vi.mocked(useInvoiceManagement).mockReturnValue({
      state: {
        invoices: [],
        isLoading: false,
        selectedInvoice: { buyerId: 1, buyerName: 'Buyer 1', totalAmount: 1000 },
      },
      actions: {
        setSelectedInvoice: mockSetSelectedInvoice,
      },
      t: tMock as ReturnType<typeof useTranslations>,
    });
    render(<InvoicePage />);
    expect(screen.getByText('Admin.Invoice.modal_title')).toBeInTheDocument();
    expect(screen.getByText('Buyer 1')).toBeInTheDocument();
  });

  it('closes detail modal', () => {
    vi.mocked(useInvoiceManagement).mockReturnValue({
      state: {
        invoices: [],
        isLoading: false,
        selectedInvoice: { buyerId: 1, buyerName: 'Buyer 1', totalAmount: 1000 },
      },
      actions: {
        setSelectedInvoice: mockSetSelectedInvoice,
      },
      t: tMock as ReturnType<typeof useTranslations>,
    });
    render(<InvoicePage />);
    const closeButton = screen.getByText('Admin.Invoice.close');
    fireEvent.click(closeButton);
    expect(mockSetSelectedInvoice).toHaveBeenCalledWith(null);
  });

  it('renders empty state', () => {
    vi.mocked(useInvoiceManagement).mockReturnValue({
      state: {
        invoices: [],
        isLoading: false,
        selectedInvoice: null,
      },
      actions: {
        setSelectedInvoice: mockSetSelectedInvoice,
      },
      t: tMock as ReturnType<typeof useTranslations>,
    });
    render(<InvoicePage />);
    expect(screen.getByText('Admin.Invoice.no_data')).toBeInTheDocument();
  });
});
