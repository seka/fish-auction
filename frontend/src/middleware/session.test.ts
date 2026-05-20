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

  it('セッションが有効なとき true を返す', async () => {
    vi.mocked(getAdminMe).mockResolvedValue({ id: 1, email: 'admin@example.com' });
    expect(await checkAdminSession(cookie)).toBe(true);
  });

  it('401 のとき false を返す', async () => {
    vi.mocked(getAdminMe).mockRejectedValue(new ApiError(401, 'Unauthorized'));
    expect(await checkAdminSession(cookie)).toBe(false);
  });

  it('401 以外のエラーのとき true を返す（fail-open）', async () => {
    vi.mocked(getAdminMe).mockRejectedValue(new ApiError(500, 'Internal Server Error'));
    expect(await checkAdminSession(cookie)).toBe(true);
  });
});

describe('checkBuyerSession', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('セッションが有効なとき true を返す', async () => {
    vi.mocked(getBuyerMe).mockResolvedValue({ id: 1, name: 'Test Buyer' });
    expect(await checkBuyerSession(cookie)).toBe(true);
  });

  it('401 のとき false を返す', async () => {
    vi.mocked(getBuyerMe).mockRejectedValue(new ApiError(401, 'Unauthorized'));
    expect(await checkBuyerSession(cookie)).toBe(false);
  });

  it('401 以外のエラーのとき true を返す（fail-open）', async () => {
    vi.mocked(getBuyerMe).mockRejectedValue(new ApiError(500, 'Internal Server Error'));
    expect(await checkBuyerSession(cookie)).toBe(true);
  });
});
