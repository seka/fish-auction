import { useQuery } from '@tanstack/react-query';
import { getCurrentBuyer } from '@/src/api/buyer_auth';

// Check if user is logged in by calling the backend
const checkAuth = async (): Promise<boolean> => {
    const buyer = await getCurrentBuyer();
    return buyer !== null;
};

export const useAuth = () => {
    const { data: isLoggedIn = false, isLoading: isChecking } = useQuery({
        queryKey: ['auth', 'me'],
        queryFn: checkAuth,
        staleTime: 5 * 60 * 1000, // 5 minutes - auth状態は頻繁に変わらない
        retry: 1, // 認証チェックは1回だけリトライ
    });

    return { isLoggedIn, isChecking };
};
