import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getCurrentBuyer, loginBuyer } from './buyer_auth';
import { apiClient } from '@/src/core/api/client';

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
    await getCurrentBuyer();
    expect(apiClient.get).toHaveBeenCalledWith('/api/buyer/me');
  });

  it('loginBuyer calls singular /api/buyers/login (Wait, checking handler... handler is registered at /api/buyers/login)', async () => {
    // Backend: s.router.HandleFunc("/api/buyers/login", ...)
    // This was actually CORRECT as plural in backend handler registry for login.
    // Let's verify backend server.go again.
    await loginBuyer({ email: 'test@example.com', password: 'password' });
    expect(apiClient.post).toHaveBeenCalledWith('/api/buyers/login', {
      email: 'test@example.com',
      password: 'password',
    });
  });
});
