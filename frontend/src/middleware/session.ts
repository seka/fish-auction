import { getAdminMe } from '@/src/data/api/auth';
import { getBuyerMe } from '@/src/data/api/buyer_auth';
import { ApiError } from '@/src/core/api/client';

export const ADMIN_SESSION_COOKIE = 'admin_session';
export const BUYER_SESSION_COOKIE = 'buyer_session';

export const checkAdminSession = async (cookie: CookieHeader): Promise<boolean> => {
  try {
    await getAdminMe(cookie);
    return true;
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) return false;
    return true;
  }
};

export const checkBuyerSession = async (cookie: CookieHeader): Promise<boolean> => {
  try {
    await getBuyerMe(cookie);
    return true;
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) return false;
    return true;
  }
};
