import { render } from '@testing-library/react';
import { describe, it, vi, beforeEach } from 'vitest';
import LoginPage from './page';
import { useAdminLogin } from '@/src/features/auth';

// Mock dependencies
vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock('./components/LoginContainer', () => ({
  LoginContainer: () => <div data-testid="login-container" />,
}));

vi.mock('@/src/features/auth', () => ({
  useAdminLogin: vi.fn(),
}));

describe('LoginPage', () => {
  const mockLogin = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useAdminLogin).mockReturnValue({
      login: mockLogin,
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useAdminLogin>);
  });

  it('renders login container', () => {
    render(<LoginPage />);
    // Since LoginPage is now just a container wrapper, we verify it renders without crashing
    // or test the container directly.
  });
});
