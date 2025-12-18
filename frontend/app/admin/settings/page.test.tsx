import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AdminSettingsPage from './page';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('AdminSettingsPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders form', () => {
        render(<AdminSettingsPage />);
        expect(screen.getByText('設定')).toBeInTheDocument();
        expect(screen.getByLabelText('現在のパスワード')).toBeInTheDocument();
    });

    it('shows error if passwords match', async () => {
        render(<AdminSettingsPage />);

        fireEvent.change(screen.getByLabelText('新しいパスワード'), { target: { value: 'password123' } });
        fireEvent.change(screen.getByLabelText('新しいパスワード（確認）'), { target: { value: 'password456' } });

        const form = screen.getByRole('button', { name: 'パスワードを変更する' }).closest('form');
        fireEvent.submit(form!);

        expect(await screen.findByText('新しいパスワードが一致しません。')).toBeInTheDocument();
    });

    it('calls API on valid submission', async () => {
        mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({}) });
        render(<AdminSettingsPage />);

        // Use less strict matching or select by index if labels are identical (they shouldn't be based on code)
        // Labels are: '現在のパスワード', '新しいパスワード', '新しいパスワード（確認）'
        // Check if there are multiple matches. '新しいパスワード' might match '新しいパスワード（確認）' partially?
        // Let's use exact: true which is default but maybe there's whitespace.

        fireEvent.change(screen.getByLabelText('現在のパスワード'), { target: { value: 'oldpass' } });
        fireEvent.change(screen.getByLabelText('新しいパスワード'), { target: { value: 'newpassword123' } });
        fireEvent.change(screen.getByLabelText('新しいパスワード（確認）'), { target: { value: 'newpassword123' } });

        const form = screen.getByRole('button', { name: 'パスワードを変更する' }).closest('form');
        fireEvent.submit(form!);

        await waitFor(() => {
            expect(mockFetch).toHaveBeenCalledWith('/api/proxy/api/admin/password', expect.any(Object));
            expect(screen.getByText('パスワードを更新しました。')).toBeInTheDocument();
        });
    });
});
