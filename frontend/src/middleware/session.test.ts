import { describe, it, expect, vi, beforeEach } from 'vitest';
import { checkAdminSession, checkBuyerSession } from './session';
import { ApiError, cookieHeader } from '@/src/core/api/client';
import { getAdminMe } from '@/src/data/api/auth';
import { getBuyerMe } from '@/src/data/api/buyer_auth';

vi.mock('@/src/data/api/auth');
vi.mock('@/src/data/api/buyer_auth');

const cookie = cookieHeader('admin_session=abc');

describe('checkAdminSession', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('セッションが有効なとき valid を返す', async () => {
    vi.mocked(getAdminMe).mockResolvedValue({ id: 1, email: 'admin@example.com' });
    expect(await checkAdminSession(cookie)).toBe('valid');
  });

  it('401 のとき invalid を返す', async () => {
    vi.mocked(getAdminMe).mockRejectedValue(new ApiError(401, 'Unauthorized'));
    expect(await checkAdminSession(cookie)).toBe('invalid');
  });

  it('ApiError(5xx) のとき error を返す', async () => {
    vi.mocked(getAdminMe).mockRejectedValue(new ApiError(500, 'Internal Server Error'));
    expect(await checkAdminSession(cookie)).toBe('error');
  });

  it('ネットワークエラーのとき error を返す', async () => {
    vi.mocked(getAdminMe).mockRejectedValue(new Error('network error'));
    expect(await checkAdminSession(cookie)).toBe('error');
  });
});

describe('checkBuyerSession', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('セッションが有効なとき valid を返す', async () => {
    vi.mocked(getBuyerMe).mockResolvedValue({ id: 1, name: 'Test Buyer' });
    expect(await checkBuyerSession(cookie)).toBe('valid');
  });

  it('401 のとき invalid を返す', async () => {
    vi.mocked(getBuyerMe).mockRejectedValue(new ApiError(401, 'Unauthorized'));
    expect(await checkBuyerSession(cookie)).toBe('invalid');
  });

  it('ApiError(5xx) のとき error を返す', async () => {
    vi.mocked(getBuyerMe).mockRejectedValue(new ApiError(500, 'Internal Server Error'));
    expect(await checkBuyerSession(cookie)).toBe('error');
  });

  it('ネットワークエラーのとき error を返す', async () => {
    vi.mocked(getBuyerMe).mockRejectedValue(new Error('network error'));
    expect(await checkBuyerSession(cookie)).toBe('error');
  });
});
