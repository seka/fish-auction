import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminSettingsPage from './page';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock next-intl
vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

describe('AdminSettingsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders form', () => {
    render(<AdminSettingsPage />);
    expect(screen.getByText('Admin.Settings.title')).toBeInTheDocument();
    expect(screen.getByLabelText('Validation.field_name.password')).toBeInTheDocument();
  });

  it('shows error if passwords match', async () => {
    render(<AdminSettingsPage />);

    fireEvent.change(screen.getByLabelText('Auth.ResetPassword.label_new_password'), {
      target: { value: 'password123' },
    });
    fireEvent.change(screen.getByLabelText('Auth.ResetPassword.label_confirm_password'), {
      target: { value: 'password456' },
    });

    const form = screen
      .getByRole('button', { name: 'Public.MyPage.submit_password_update' })
      .closest('form');
    fireEvent.submit(form!);

    expect(await screen.findByText('Validation.password_mismatch')).toBeInTheDocument();
  });

  it('calls API on valid submission', async () => {
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({}) });
    render(<AdminSettingsPage />);

    // Use less strict matching or select by index if labels are identical (they shouldn't be based on code)
    // Labels are: '現在のパスワード', '新しいパスワード', '新しいパスワード（確認）'
    // Check if there are multiple matches. '新しいパスワード' might match '新しいパスワード（確認）' partially?
    // Let's use exact: true which is default but maybe there's whitespace.

    fireEvent.change(screen.getByLabelText('Validation.field_name.password'), {
      target: { value: 'oldpass' },
    });
    fireEvent.change(screen.getByLabelText('Auth.ResetPassword.label_new_password'), {
      target: { value: 'newpassword123' },
    });
    fireEvent.change(screen.getByLabelText('Auth.ResetPassword.label_confirm_password'), {
      target: { value: 'newpassword123' },
    });

    const form = screen
      .getByRole('button', { name: 'Public.MyPage.submit_password_update' })
      .closest('form');
    fireEvent.submit(form!);

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('/api/proxy/api/admin/password', expect.any(Object));
      expect(screen.getByText('Validation.password_updated')).toBeInTheDocument();
    });
  });
});
