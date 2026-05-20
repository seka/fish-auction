import { getAdminMe, adminSessionCookie } from '@/src/data/api/auth';
import { getBuyerMe, buyerSessionCookie } from '@/src/data/api/buyer_auth';
import { ApiError } from '@/src/core/api/client';

export { adminSessionCookie, buyerSessionCookie };

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
