import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LoginPage from './page';
import { useLogin } from './_hooks/useAuth';

// Mock dependencies
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

vi.mock('next/navigation', () => ({
    useRouter: () => ({
        push: vi.fn(),
    }),
}));

vi.mock('./_hooks/useAuth', () => ({
    useLogin: vi.fn(),
}));

describe('LoginPage', () => {
    const mockLogin = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useLogin as any).mockReturnValue({
            login: mockLogin,
            isLoading: false,
        });
    });

    it('renders login form', () => {
        render(<LoginPage />);

        expect(screen.getByPlaceholderText('Common.email')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Common.password')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Common.submit' })).toBeInTheDocument();
    });

    it('submits form with valid data', async () => {
        mockLogin.mockResolvedValue(true);
        render(<LoginPage />);

        const emailInput = screen.getByPlaceholderText('Common.email');
        const passwordInput = screen.getByPlaceholderText('Common.password');
        const submitButton = screen.getByRole('button', { name: 'Common.submit' });

        fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
        fireEvent.change(passwordInput, { target: { value: 'password123' } });
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(mockLogin).toHaveBeenCalledWith('test@example.com', 'password123');
        });
    });

    it('displays error on login failure', async () => {
        mockLogin.mockResolvedValue(false);
        render(<LoginPage />);

        const emailInput = screen.getByPlaceholderText('Common.email');
        const passwordInput = screen.getByPlaceholderText('Common.password');
        const submitButton = screen.getByRole('button', { name: 'Common.submit' });

        fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
        fireEvent.change(passwordInput, { target: { value: 'wrongpass' } });
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText('Admin.Login.error_invalid_password')).toBeInTheDocument();
        });
    });
});
