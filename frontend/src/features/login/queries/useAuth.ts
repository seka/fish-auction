import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { login } from '@/src/data/api/auth';
import { getCurrentBuyer, logoutBuyer } from '@/src/data/api/buyer_auth';
import { authKeys } from '@/src/data/queries/auth/keys';
import { useRouter } from 'next/navigation';

/**
 * ログインミューテーションフック（管理者用）
 */
export const useLoginMutation = () => {
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      login(email, password),
  });
};

/**
 * 認証状態とログアウト処理を管理するフック（バイヤー用）
 */
export const useBuyerAuth = () => {
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: buyer, isLoading } = useQuery({
    queryKey: authKeys.me(),
    queryFn: getCurrentBuyer,
    retry: false,
  });

  const isLoggedIn = !!buyer;

  const logout = async () => {
    await logoutBuyer();
    queryClient.setQueryData(authKeys.me(), null);
    router.push('/login/buyer');
  };

  return {
    buyer,
    isLoggedIn,
    isLoading,
    logout,
  };
};
