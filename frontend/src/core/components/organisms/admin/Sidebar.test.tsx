import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Sidebar } from './Sidebar';

// Mock next/navigation
vi.mock('next/navigation', () => ({
  usePathname: vi.fn(() => '/admin'),
}));

// Mock next-intl
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

describe('Sidebar', () => {
  const mockLogout = vi.fn();

  it('renders all menu items', () => {
    render(<Sidebar onLogout={mockLogout} />);

    expect(screen.getByText('title')).toBeInTheDocument();
    expect(screen.getByText('back_to_top')).toBeInTheDocument();
    expect(screen.getByText('dashboard')).toBeInTheDocument();
    expect(screen.getByText('fishermen')).toBeInTheDocument();
  });

  it('highlights active link based on pathname', () => {
    render(<Sidebar onLogout={mockLogout} />);

    const dashboardLink = screen.getByRole('link', { name: /dashboard/ });
    expect(dashboardLink).toBeInTheDocument();
    expect(dashboardLink).toHaveAttribute('href', '/admin');
  });
});
