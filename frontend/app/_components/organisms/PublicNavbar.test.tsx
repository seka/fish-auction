import { render, screen, fireEvent } from '@testing-library/react';
import { PublicNavbar } from './PublicNavbar';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { usePathname, useRouter } from 'next/navigation';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import * as buyerAuth from '@/src/api/buyer_auth';

// Mocks
vi.mock('next/navigation', () => ({
  usePathname: vi.fn(),
  useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  useQueryClient: vi.fn(),
}));

vi.mock('@/src/api/buyer_auth', () => ({
  getCurrentBuyer: vi.fn(),
  logoutBuyer: vi.fn(),
}));

describe('PublicNavbar', () => {
  const mockPush = vi.fn();
  const mockSetQueryData = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useRouter).mockReturnValue({ push: mockPush } as any // eslint-disable-line @typescript-eslint/no-explicit-any
);
    vi.mocked(useQueryClient).mockReturnValue({ setQueryData: mockSetQueryData } as any // eslint-disable-line @typescript-eslint/no-explicit-any
);
    vi.mocked(usePathname).mockReturnValue('/');
  });

  it('renders correctly when not logged in', () => {
    vi.mocked(useQuery).mockReturnValue({ data: null } as any // eslint-disable-line @typescript-eslint/no-explicit-any
);

    render(<PublicNavbar />);

    expect(screen.getByText('Common.app_name')).toBeInTheDocument();
    expect(screen.getByText('Navbar.active_auctions')).toBeInTheDocument();
    expect(screen.getByText('Navbar.login')).toBeInTheDocument();
    expect(screen.queryByText('Navbar.logout')).not.toBeInTheDocument();
    expect(screen.queryByText('Navbar.mypage')).not.toBeInTheDocument();
  });

  it('renders correctly when logged in', () => {
    vi.mocked(useQuery).mockReturnValue({ data: { id: 1, name: 'Buyer' } } as any // eslint-disable-line @typescript-eslint/no-explicit-any
);

    render(<PublicNavbar />);

    expect(screen.getByText('Navbar.logout')).toBeInTheDocument();
    expect(screen.getByText('Navbar.mypage')).toBeInTheDocument();
    expect(screen.queryByText('Navbar.login')).not.toBeInTheDocument();
  });

  it('handles logout correctly', async () => {
    vi.mocked(useQuery).mockReturnValue({ data: { id: 1, name: 'Buyer' } } as any // eslint-disable-line @typescript-eslint/no-explicit-any
);
    const mockLogoutBuyer = vi.spyOn(buyerAuth, 'logoutBuyer').mockResolvedValue(undefined);

    render(<PublicNavbar />);

    const logoutButton = screen.getByText('Navbar.logout');
    fireEvent.click(logoutButton);

    expect(mockLogoutBuyer).toHaveBeenCalled();
    // Wait for async actions if necessary, but here spy serves.
    // In real execution, we might need waitFor if state update is involved,
    // but here handleLogout calls router.push immediately after await.
    // Since logoutBuyer is mocked to resolve immediately:
  });

  it('does not render on admin pages', () => {
    vi.mocked(usePathname).mockReturnValue('/admin/dashboard');
    const { container } = render(<PublicNavbar />);
    expect(container).toBeEmptyDOMElement();
  });
});
