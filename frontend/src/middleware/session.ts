import { getAdminMe, adminSessionCookie } from '@/src/data/api/auth';
import { getBuyerMe, buyerSessionCookie } from '@/src/data/api/buyer_auth';
import { ApiError } from '@/src/core/api/client';

export { adminSessionCookie, buyerSessionCookie };

export type SessionResult = 'valid' | 'invalid' | 'error';

export const checkAdminSession = async (cookie: CookieHeader): Promise<SessionResult> => {
  try {
    await getAdminMe(cookie);
    return 'valid';
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) return 'invalid';
    return 'error';
  }
};

export const checkBuyerSession = async (cookie: CookieHeader): Promise<SessionResult> => {
  try {
    await getBuyerMe(cookie);
    return 'valid';
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) return 'invalid';
    return 'error';
  }
};
