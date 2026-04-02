import { render, screen, fireEvent } from '@testing-library/react';
import { PublicNavbar } from './PublicNavbar';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { usePathname } from 'next/navigation';

// Mocks
vi.mock('next/navigation', () => ({
  usePathname: vi.fn(),
  useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

describe('PublicNavbar', () => {
  const mockOnLogout = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(usePathname).mockReturnValue('/');
  });

  it('renders correctly when not logged in', () => {
    render(<PublicNavbar isLoggedIn={false} onLogout={mockOnLogout} />);

    expect(screen.getByText('Common.app_name')).toBeInTheDocument();
    expect(screen.getByText('Navbar.active_auctions')).toBeInTheDocument();
    expect(screen.getByText('Navbar.login')).toBeInTheDocument();
    expect(screen.queryByText('Navbar.logout')).not.toBeInTheDocument();
    expect(screen.queryByText('Navbar.mypage')).not.toBeInTheDocument();
  });

  it('renders correctly when logged in', () => {
    render(<PublicNavbar isLoggedIn={true} buyerName="Test User" onLogout={mockOnLogout} />);

    expect(screen.getByText('Test User Navbar.honorific')).toBeInTheDocument();
    expect(screen.getByText('Navbar.logout')).toBeInTheDocument();
    expect(screen.getByText('Navbar.mypage')).toBeInTheDocument();
    expect(screen.queryByText('Navbar.login')).not.toBeInTheDocument();
  });

  it('calls onLogout when logout button is clicked', async () => {
    render(<PublicNavbar isLoggedIn={true} buyerName="Test User" onLogout={mockOnLogout} />);

    const logoutButton = screen.getByText('Navbar.logout');
    fireEvent.click(logoutButton);

    expect(mockOnLogout).toHaveBeenCalled();
  });

  it('does not render on admin pages', () => {
    vi.mocked(usePathname).mockReturnValue('/admin/dashboard');
    const { container } = render(<PublicNavbar isLoggedIn={false} onLogout={mockOnLogout} />);
    expect(container).toBeEmptyDOMElement();
  });
});
