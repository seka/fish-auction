import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthorizableAdminSidebar } from './AuthorizableAdminSidebar';
import { useAdminLogoutMutation } from '@/src/features/auth';
import { useRouter } from 'next/navigation';

// Mock dependencies
vi.mock('@/src/features/auth', () => ({
  useAdminLogoutMutation: vi.fn(),
}));

vi.mock('next/navigation', () => ({
  useRouter: vi.fn(),
}));

vi.mock('@/src/core/components/organisms/AdminSidebar', () => ({
  AdminSidebar: vi.fn(({ onLogout }) => (
    <div data-testid="admin-sidebar">
      <button onClick={onLogout}>Logout</button>
    </div>
  )),
}));

describe('AuthorizableAdminSidebar', () => {
  const mockMutateAsync = vi.fn();
  const mockPush = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useAdminLogoutMutation).mockReturnValue({
      mutateAsync: mockMutateAsync,
    } as any);
    vi.mocked(useRouter).mockReturnValue({
      push: mockPush,
    } as any);
  });

  it('calls logout mutation and redirects on handleLogout', async () => {
    render(<AuthorizableAdminSidebar />);
    
    // Trigger logout
    fireEvent.click(screen.getByText('Logout'));

    // Verify mutation was called
    expect(mockMutateAsync).toHaveBeenCalledTimes(1);

    // Wait for the async logout to finish and verification redirection
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/login/admin');
    });
  });
});
