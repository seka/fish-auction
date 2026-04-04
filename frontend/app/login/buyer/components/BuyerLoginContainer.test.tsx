import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BuyerLoginContainer } from './BuyerLoginContainer';
import { useRouter } from 'next/navigation';
import { loginBuyer } from '@/src/data/api/buyer_auth';
import { useQueryClient } from '@tanstack/react-query';
import { authKeys } from '@/src/data/queries/auth/keys';

// Mocks
vi.mock('next/navigation', () => ({
  useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

vi.mock('@/src/data/api/buyer_auth', () => ({
  loginBuyer: vi.fn(),
}));

vi.mock('@tanstack/react-query', () => ({
  useQueryClient: vi.fn(),
}));

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

describe('BuyerLoginContainer', () => {
  const mockRouter = { push: vi.fn() };
  const mockInvalidateQueries = vi.fn();
  const mockQueryClient = { invalidateQueries: mockInvalidateQueries };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useRouter).mockReturnValue(mockRouter as unknown as ReturnType<typeof useRouter>);
    vi.mocked(useQueryClient).mockReturnValue(mockQueryClient as unknown as ReturnType<typeof useQueryClient>);
  });

  it('renders login form', () => {
    render(<BuyerLoginContainer />);
    expect(screen.getByPlaceholderText('Common.email')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Common.password')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Public.Login.submit' })).toBeInTheDocument();
  });

  it('calls loginBuyer and invalidates authKeys.me() on success', async () => {
    vi.mocked(loginBuyer).mockResolvedValue({ id: 1, email: 'test@example.com' } as unknown as Awaited<ReturnType<typeof loginBuyer>>);

    render(<BuyerLoginContainer />);

    fireEvent.change(screen.getByPlaceholderText('Common.email'), {
      target: { value: 'test@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Common.password'), {
      target: { value: 'password' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Public.Login.submit' }));

    await waitFor(() => {
      expect(loginBuyer).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password',
      });
    });

    await waitFor(() => {
      expect(mockInvalidateQueries).toHaveBeenCalledWith({
        queryKey: authKeys.me(),
      });
    });

    await waitFor(() => {
      expect(mockRouter.push).toHaveBeenCalledWith('/auctions');
    });
  });

  it('displays error message on login failure', async () => {
    vi.mocked(loginBuyer).mockResolvedValue(null as unknown as Awaited<ReturnType<typeof loginBuyer>>);

    render(<BuyerLoginContainer />);

    fireEvent.change(screen.getByPlaceholderText('Common.email'), {
      target: { value: 'test@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Common.password'), {
      target: { value: 'wrong-password' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Public.Login.submit' }));

    await waitFor(() => {
      expect(screen.getByText('Public.Login.error_credentials')).toBeInTheDocument();
    });

    expect(mockInvalidateQueries).not.toHaveBeenCalled();
    expect(mockRouter.push).not.toHaveBeenCalled();
  });
});
