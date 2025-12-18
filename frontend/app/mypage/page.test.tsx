import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import MyPage from './page';
import { useMyPage } from './_hooks/useMyPage';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';

// Mocks
vi.mock('next-intl', () => ({
    useTranslations: () => (key: string) => key,
}));

vi.mock('./_hooks/useMyPage', () => ({
    useMyPage: vi.fn(),
}));

// Mock Styled System components if necessary, but we seem to rely on implementation's imports which are real.
// If Styled System components cause issues in JSDOM (unlikely for basic ones), we might need to mock them,
// but usually they render as normal elements with classes.

describe('MyPage', () => {
    const mockHandleLogout = vi.fn();
    const mockSetActiveTab = vi.fn();
    const mockHandleUpdatePassword = vi.fn();
    const mockSetCurrentPassword = vi.fn();
    const mockSetNewPassword = vi.fn();
    const mockSetConfirmPassword = vi.fn();

    const defaultMockValues = {
        t: (key: string) => key,
        activeTab: 'purchases',
        setActiveTab: mockSetActiveTab,
        currentPassword: '',
        setCurrentPassword: mockSetCurrentPassword,
        newPassword: '',
        setNewPassword: mockSetNewPassword,
        confirmPassword: '',
        setConfirmPassword: mockSetConfirmPassword,
        message: null,
        isPasswordUpdating: false,
        purchases: [],
        auctions: [],
        isLoading: false,
        handleLogout: mockHandleLogout,
        handleUpdatePassword: mockHandleUpdatePassword,
    };

    beforeEach(() => {
        vi.clearAllMocks();
        (useMyPage as any).mockReturnValue(defaultMockValues);
    });

    it('renders loading state', () => {
        (useMyPage as any).mockReturnValue({ ...defaultMockValues, isLoading: true });
        render(<MyPage />);
        expect(screen.getByText('読み込み中...')).toBeInTheDocument();
    });

    it('renders initial state (purchases tab)', () => {
        render(<MyPage />);
        // Title
        expect(screen.getAllByText('Public.MyPage.purchase_history').length).toBeGreaterThan(0);
        // Empty state
        expect(screen.getByText('Public.MyPage.no_history')).toBeInTheDocument();
    });

    it('renders purchases list', () => {
        const mockPurchases = [
            {
                id: 1,
                itemId: 101,
                fishType: 'Tuna',
                quantity: 10,
                unit: 'kg',
                price: 50000,
                auctionId: 1,
                auctionDate: '2023-12-01',
                createdAt: '2023-12-01T10:00:00Z',
            },
        ];
        (useMyPage as any).mockReturnValue({
            ...defaultMockValues,
            activeTab: 'purchases',
            purchases: mockPurchases,
        });

        render(<MyPage />);
        expect(screen.getByText('Tuna')).toBeInTheDocument();
        expect(screen.getByText('¥50,000')).toBeInTheDocument();
    });

    it('switches tabs', () => {
        render(<MyPage />);

        const auctionsTab = screen.getByText('Public.MyPage.participating_auctions');
        fireEvent.click(auctionsTab);

        expect(mockSetActiveTab).toHaveBeenCalledWith('auctions');
    });

    it('renders participating auctions tab', () => {
        (useMyPage as any).mockReturnValue({
            ...defaultMockValues,
            activeTab: 'auctions',
        });
        render(<MyPage />);

        expect(screen.getAllByText('Public.MyPage.participating_auctions').length).toBeGreaterThan(0);
        expect(screen.getByText('Public.MyPage.no_participating')).toBeInTheDocument();
    });

    it('renders settings tab and password form', () => {
        (useMyPage as any).mockReturnValue({
            ...defaultMockValues,
            activeTab: 'settings',
        });
        render(<MyPage />);

        expect(screen.getByText('パスワード変更')).toBeInTheDocument();
        expect(screen.getByText('現在のパスワード')).toBeInTheDocument();
        expect(screen.getByText('新しいパスワード')).toBeInTheDocument();
        expect(screen.getByText('新しいパスワード（確認）')).toBeInTheDocument();
    });

    it('calls logout', () => {
        render(<MyPage />);
        const logoutButton = screen.getByText('Common.logout');
        fireEvent.click(logoutButton);
        expect(mockHandleLogout).toHaveBeenCalled();
    });

    it('updates active tab when clicked', () => {
        render(<MyPage />);
        const settingsTab = screen.getByText('設定');
        fireEvent.click(settingsTab);
        expect(mockSetActiveTab).toHaveBeenCalledWith('settings');
    });

    it('handles password update flow', () => {
        (useMyPage as any).mockReturnValue({
            ...defaultMockValues,
            activeTab: 'settings',
        });
        render(<MyPage />);

        const currentPasswordInput = screen.getByText('現在のパスワード').closest('div')?.querySelector('input');
        fireEvent.change(currentPasswordInput!, { target: { value: 'current123' } });
        expect(mockSetCurrentPassword).toHaveBeenCalledWith('current123');

        const newPasswordInput = screen.getByText('新しいパスワード').closest('div')?.querySelector('input');
        fireEvent.change(newPasswordInput!, { target: { value: 'password123' } });
        expect(mockSetNewPassword).toHaveBeenCalledWith('password123');

        const confirmPasswordInput = screen.getByText('新しいパスワード（確認）').closest('div')?.querySelector('input');
        fireEvent.change(confirmPasswordInput!, { target: { value: 'password123' } });
        expect(mockSetConfirmPassword).toHaveBeenCalledWith('password123');

        const submitButton = screen.getByText('パスワードを変更する');
        fireEvent.submit(submitButton.closest('form')!);
        expect(mockHandleUpdatePassword).toHaveBeenCalled();
    });
});
