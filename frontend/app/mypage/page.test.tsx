import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import MyPage from './page';
import { useMyPage } from '@/src/features/mypage/states/useMyPage';
import { useMyPurchases } from '@/src/data/queries/buyerPurchase/useQuery';
import { useParticipatingAuctions } from '@/src/data/queries/buyerAuction/useQuery';

// Mocks
vi.mock('next-intl', () => ({
  useTranslations: (namespace?: string) => (key: string) =>
    namespace ? `${namespace}.${key}` : key,
}));

vi.mock('@/src/features/mypage/states/useMyPage', () => ({
  useMyPage: vi.fn(),
}));

vi.mock('@/src/data/queries/buyerPurchase/useQuery', () => ({
  useMyPurchases: vi.fn(),
}));

vi.mock('@/src/data/queries/buyerAuction/useQuery', () => ({
  useParticipatingAuctions: vi.fn(),
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
    isLoading: false,
    handleLogout: mockHandleLogout,
    purchases: [],
    auctions: [],
    invoices: [],
    passwordState: {
      currentPassword: '',
      setCurrentPassword: mockSetCurrentPassword,
      newPassword: '',
      setNewPassword: mockSetNewPassword,
      confirmPassword: '',
      setConfirmPassword: mockSetConfirmPassword,
      passwordMessage: null,
      handleUpdatePassword: mockHandleUpdatePassword,
      isPasswordUpdating: false,
    },
  };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useMyPage).mockReturnValue(
      defaultMockValues as unknown as ReturnType<typeof useMyPage>,
    );
    vi.mocked(useMyPurchases).mockReturnValue({
      purchases: [],
      isLoading: false,
    } as unknown as ReturnType<typeof useMyPurchases>);
    vi.mocked(useParticipatingAuctions).mockReturnValue({
      auctions: [],
      isLoading: false,
    } as unknown as ReturnType<typeof useParticipatingAuctions>);
  });

  it('renders loading state', () => {
    vi.mocked(useMyPurchases).mockReturnValue({
      purchases: [],
      isLoading: true,
    } as unknown as ReturnType<typeof useMyPurchases>);
    render(<MyPage />);
    expect(screen.getByText('Common.loading')).toBeInTheDocument();
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
    vi.mocked(useMyPurchases).mockReturnValue({
      purchases: mockPurchases,
      isLoading: false,
    } as unknown as ReturnType<typeof useMyPurchases>);

    vi.mocked(useMyPage).mockReturnValue({
      ...defaultMockValues,
      activeTab: 'purchases',
      purchases: mockPurchases,
    } as unknown as ReturnType<typeof useMyPage>);

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
    vi.mocked(useMyPage).mockReturnValue({
      ...defaultMockValues,
      activeTab: 'auctions',
    } as unknown as ReturnType<typeof useMyPage>);
    render(<MyPage />);

    expect(screen.getAllByText('Public.MyPage.participating_auctions').length).toBeGreaterThan(0);
    expect(screen.getByText('Public.MyPage.no_participating')).toBeInTheDocument();
  });

  it('renders settings tab and password form', () => {
    vi.mocked(useMyPage).mockReturnValue({
      ...defaultMockValues,
      activeTab: 'settings',
    } as unknown as ReturnType<typeof useMyPage>);
    render(<MyPage />);

    expect(screen.getByText('Public.MyPage.password_change_title')).toBeInTheDocument();
    expect(screen.getByText('Validation.field_name.password')).toBeInTheDocument();
    expect(screen.getByText('Auth.ResetPassword.label_new_password')).toBeInTheDocument();
    expect(screen.getByText('Auth.ResetPassword.label_confirm_password')).toBeInTheDocument();
  });

  it('calls logout', () => {
    render(<MyPage />);
    const logoutButton = screen.getByText('Common.logout');
    fireEvent.click(logoutButton);
    expect(mockHandleLogout).toHaveBeenCalled();
  });

  it('updates active tab when clicked', () => {
    render(<MyPage />);
    const settingsTab = screen.getByText('Public.MyPage.settings');
    fireEvent.click(settingsTab);
    expect(mockSetActiveTab).toHaveBeenCalledWith('settings');
  });

  it('handles password update flow', () => {
    vi.mocked(useMyPage).mockReturnValue({
      ...defaultMockValues,
      activeTab: 'settings',
    } as unknown as ReturnType<typeof useMyPage>);
    render(<MyPage />);

    const currentPasswordInput = screen
      .getByText('Validation.field_name.password')
      .closest('div')
      ?.querySelector('input');
    fireEvent.change(currentPasswordInput!, { target: { value: 'current123' } });
    expect(mockSetCurrentPassword).toHaveBeenCalledWith('current123');

    const newPasswordInput = screen
      .getByText('Auth.ResetPassword.label_new_password')
      .closest('div')
      ?.querySelector('input');
    fireEvent.change(newPasswordInput!, { target: { value: 'password123' } });
    expect(mockSetNewPassword).toHaveBeenCalledWith('password123');

    const confirmPasswordInput = screen
      .getByText('Auth.ResetPassword.label_confirm_password')
      .closest('div')
      ?.querySelector('input');
    fireEvent.change(confirmPasswordInput!, { target: { value: 'password123' } });
    expect(mockSetConfirmPassword).toHaveBeenCalledWith('password123');

    const submitButton = screen.getByText('Public.MyPage.submit_password_update');
    fireEvent.submit(submitButton.closest('form')!);
    expect(mockHandleUpdatePassword).toHaveBeenCalled();
  });
});
