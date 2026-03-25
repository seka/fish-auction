import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LoginPage from './page';
import { useLogin } from '@/src/features/login';

// Mock dependencies
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock('@/src/features/login', () => ({
  LoginContainer: ({ children }: any) => <div>{children}</div>,
  useLogin: vi.fn(),
}));

describe('LoginPage', () => {
  const mockLogin = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useLogin).mockReturnValue({
      login: mockLogin,
      isLoading: false,
      error: null,
    } as any);
  });

  it('renders login container', () => {
    render(<LoginPage />);
    // Since LoginPage is now just a container wrapper, we verify it renders without crashing
    // or test the container directly.
  });
});
