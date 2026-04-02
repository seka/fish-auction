import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthorizablePublicNavbar } from './AuthorizablePublicNavbar';
import { useBuyerAuth } from '../queries/useAuth';
import { usePathname } from 'next/navigation';

// Mock dependecies
vi.mock('../queries/useAuth');
vi.mock('next/navigation', () => ({
  usePathname: vi.fn(),
}));
vi.mock('@organisms', () => ({
  PublicNavbar: vi.fn(({ isLoggedIn, isLoading, buyerName, onLogout }) => (
    <div data-testid="public-navbar">
      <span>{isLoading ? 'Loading...' : 'Loaded'}</span>
      <span>{isLoggedIn ? `Logged In as ${buyerName}` : 'Logged Out'}</span>
      <button onClick={onLogout}>Logout</button>
    </div>
  )),
}));

describe('AuthorizablePublicNavbar', () => {
  const mockLogout = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(usePathname).mockReturnValue('/');
  });

  it('does not render on admin pages', () => {
    vi.mocked(usePathname).mockReturnValue('/admin/dashboard');
    vi.mocked(useBuyerAuth).mockReturnValue({
      buyer: null,
      isLoggedIn: false,
      isLoading: false,
      logout: mockLogout,
    });

    const { container } = render(<AuthorizablePublicNavbar />);
    expect(container).toBeEmptyDOMElement();
  });

  it('passes isLoading: true when auth is loading', () => {
    vi.mocked(useBuyerAuth).mockReturnValue({
      buyer: null,
      isLoggedIn: false,
      isLoading: true,
      logout: mockLogout,
    });

    render(<AuthorizablePublicNavbar />);

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('passes isLoggedIn: true and buyerName when logged in', () => {
    vi.mocked(useBuyerAuth).mockReturnValue({
      buyer: { name: 'Test Buyer' },
      isLoggedIn: true,
      isLoading: false,
      logout: mockLogout,
    });

    render(<AuthorizablePublicNavbar />);

    expect(screen.getByText('Loaded')).toBeInTheDocument();
    expect(screen.getByText('Logged In as Test Buyer')).toBeInTheDocument();
  });

  it('passes isLoggedIn: false when logged out', () => {
    vi.mocked(useBuyerAuth).mockReturnValue({
      buyer: null,
      isLoggedIn: false,
      isLoading: false,
      logout: mockLogout,
    });

    render(<AuthorizablePublicNavbar />);

    expect(screen.getByText('Loaded')).toBeInTheDocument();
    expect(screen.getByText('Logged Out')).toBeInTheDocument();
  });

  it('calls logout when onLogout is triggered', () => {
    vi.mocked(useBuyerAuth).mockReturnValue({
      buyer: { name: 'Test Buyer' },
      isLoggedIn: true,
      isLoading: false,
      logout: mockLogout,
    });

    render(<AuthorizablePublicNavbar />);
    fireEvent.click(screen.getByText('Logout'));

    expect(mockLogout).toHaveBeenCalledTimes(1);
  });
});
