import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import Home from './page';

// Mock next-intl
vi.mock('next-intl', () => ({
    useTranslations: () => {
        const t = (key: string) => key;
        t.raw = (key: string) => key;
        return t;
    },
}));

describe('Home Page', () => {
    it('renders logo', () => {
        render(<Home />);
        const logo = screen.getByAltText(/FISHING AUCTION Logo/i);
        expect(logo).toBeInTheDocument();
    });

    it('renders links to admin and auctions', () => {
        render(<Home />);

        const adminLink = screen.getByRole('link', { name: /admin_panel/i });
        const auctionLink = screen.getByRole('link', { name: /auction_venue/i });

        expect(adminLink).toBeInTheDocument();
        expect(adminLink).toHaveAttribute('href', '/admin');

        expect(auctionLink).toBeInTheDocument();
        expect(auctionLink).toHaveAttribute('href', '/auctions');
    });
});
