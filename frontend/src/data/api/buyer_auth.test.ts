import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getCurrentBuyer, loginBuyer, logoutBuyer } from './buyer_auth';
import { apiClient } from '@/src/core/api/client';
import { Buyer } from '@/src/models';

// Mock the apiClient
vi.mock('@/src/core/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

describe('Buyer Auth API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('getCurrentBuyer calls singular /api/buyer/me', async () => {
    const mockBuyer: Buyer = { id: 1, name: "Test Buyer" };
    (apiClient.get as any).mockResolvedValue(mockBuyer);

    await getCurrentBuyer();

    expect(apiClient.get).toHaveBeenCalledWith('/api/buyer/me');
  });

  it('loginBuyer calls singular /api/buyer/login', async () => {
    const credentials = { email: "test@example.com", password: "password" };
    const mockBuyer: Buyer = { id: 1, name: "Test Buyer" };
    (apiClient.post as any).mockResolvedValue(mockBuyer);

    await loginBuyer(credentials);

    expect(apiClient.post).toHaveBeenCalledWith('/api/buyer/login', credentials);
  });

  it('logoutBuyer calls singular /api/buyer/logout', async () => {
    (apiClient.post as any).mockResolvedValue({});

    await logoutBuyer();

    expect(apiClient.post).toHaveBeenCalledWith('/api/buyer/logout', {});
  });
});
