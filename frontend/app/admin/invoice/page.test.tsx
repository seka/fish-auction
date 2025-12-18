import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import InvoicePage from './page';
import { useInvoicePage } from './_hooks/useInvoicePage';

// Mock hook
vi.mock('./_hooks/useInvoicePage');
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

describe('InvoicePage', () => {
    const mockSetSelectedInvoice = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useInvoicePage as any).mockReturnValue({
            state: {
                invoices: [
                    { buyerId: 1, buyerName: 'Buyer 1', totalAmount: 1000 },
                ],
                isLoading: false,
                selectedInvoice: null,
            },
            actions: {
                setSelectedInvoice: mockSetSelectedInvoice,
            },
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
        (useInvoicePage as any).mockReturnValue({
            state: {
                invoices: [],
                isLoading: false,
                selectedInvoice: { buyerId: 1, buyerName: 'Buyer 1', totalAmount: 1000 },
            },
            actions: {
                setSelectedInvoice: mockSetSelectedInvoice,
            },
        });
        render(<InvoicePage />);
        expect(screen.getByText('Admin.Invoice.modal_title')).toBeInTheDocument();
        expect(screen.getByText('Buyer 1')).toBeInTheDocument();
    });
});
