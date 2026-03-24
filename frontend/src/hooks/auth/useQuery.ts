import { useQuery } from '@tanstack/react-query';
import { getCurrentBuyer } from '@/src/api/buyer_auth';
import { authKeys } from './keys';

// Check if user is logged in by calling the backend
const checkAuth = async (): Promise<boolean> => {
  try {
    const buyer = await getCurrentBuyer();
    return buyer !== null;
  } catch {
    return false;
  }
};

export const useAuthQuery = () => {
  const { data: isLoggedIn = false, isLoading: isChecking } = useQuery({
    queryKey: authKeys.me(),
    queryFn: checkAuth,
    staleTime: 5 * 60 * 1000, // 5 minutes - auth状態は頻繁に変わらない
    retry: 1, // 認証チェックは1回だけリトライ
  });

  return { isLoggedIn, isChecking };
};
