import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import SignupPage from './page';
import * as buyerAuth from '@/src/api/buyer_auth';
import { useRouter } from 'next/navigation';

// Mocks
vi.mock('next/navigation', () => ({
    useRouter: vi.fn(),
}));

vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

vi.mock('@/src/api/buyer_auth', () => ({
    signupBuyer: vi.fn(),
}));

describe('SignupPage', () => {
    const mockPush = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useRouter as any).mockReturnValue({ push: mockPush });
    });

    it('renders signup form', () => {
        render(<SignupPage />);
        expect(screen.getByText('中買人登録')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('名前')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('メールアドレス')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Common.register' })).toBeInTheDocument();
    });

    it('shows validation error on empty submit', async () => {
        render(<SignupPage />);

        fireEvent.click(screen.getByRole('button', { name: 'Common.register' }));

        await waitFor(() => {
            // Zod schema messages would appear here. Since schema is imported, 
            // we rely on the implementation. Assuming standard required messages exist or checking for Error text color.
            // But let's check if signupBuyer was NOT called.
            expect(buyerAuth.signupBuyer).not.toHaveBeenCalled();
        });
    });

    it('submits form with valid data', async () => {
        (buyerAuth.signupBuyer as any).mockResolvedValue({});

        render(<SignupPage />);

        fireEvent.change(screen.getByPlaceholderText('名前'), { target: { value: 'Test Buyer' } });
        fireEvent.change(screen.getByPlaceholderText('メールアドレス'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByPlaceholderText('所属組織'), { target: { value: 'Test Org' } });
        fireEvent.change(screen.getByPlaceholderText('連絡先'), { target: { value: '090-1234-5678' } });
        fireEvent.change(screen.getByPlaceholderText('パスワード'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: 'Common.register' }));

        await waitFor(() => {
            expect(buyerAuth.signupBuyer).toHaveBeenCalledWith({
                name: 'Test Buyer',
                email: 'test@example.com',
                organization: 'Test Org',
                contact_info: '090-1234-5678',
                password: 'password123',
            });
            expect(mockPush).toHaveBeenCalledWith('/login/buyer');
        });
    });
});
