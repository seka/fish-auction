import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import AdminDashboard from './page';

// Mock next-intl
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

describe('AdminDashboard', () => {
    it('renders title and subtitle', () => {
        render(<AdminDashboard />);
        expect(screen.getByText('Admin.Dashboard.title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Dashboard.subtitle')).toBeInTheDocument();
    });

    it('renders all menu items', () => {
        render(<AdminDashboard />);

        // Simpler check using key text content which we mocked.
        // We know these keys are rendered.
        expect(screen.getByText('Admin.Dashboard.fishermen_title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Dashboard.buyers_title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Dashboard.items_title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Dashboard.venues_title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Dashboard.auctions_title')).toBeInTheDocument();
        expect(screen.getByText('Admin.Dashboard.invoice_title')).toBeInTheDocument();
    });
});
