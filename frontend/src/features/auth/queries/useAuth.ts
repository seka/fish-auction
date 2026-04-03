import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { login, logout as logoutAdmin } from '@/src/data/api/auth';
import { getCurrentBuyer, loginBuyer, logoutBuyer } from '@/src/data/api/buyer_auth';
import { authKeys } from '@/src/data/queries/auth/keys';
import { useRouter } from 'next/navigation';

/**
 * ログインミューテーションフック（管理者用）
 */
export const useAdminLoginMutation = () => {
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      login(email, password),
  });
};

/**
 * ログアウトミューテーションフック（管理者用）
 */
export const useAdminLogoutMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: logoutAdmin,
    onSuccess: () => {
      queryClient.clear();
    },
  });
};

/**
 * ログインミューテーションフック（バイヤー用）
 */
export const useBuyerLoginMutation = () => {
  return useMutation({
    mutationFn: loginBuyer,
  });
};

/**
 * ログアウトミューテーションフック（バイヤー用）
 */
export const useBuyerLogoutMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: logoutBuyer,
    onSuccess: () => {
      queryClient.setQueryData(authKeys.me(), null);
    },
  });
};

/**
 * 認証状態とログアウト処理を管理するフック（バイヤー用）
 */
export const useBuyerAuth = () => {
  const router = useRouter();
  const logoutMutation = useBuyerLogoutMutation();

  const { data: buyer, isLoading } = useQuery({
    queryKey: authKeys.me(),
    queryFn: getCurrentBuyer,
    retry: false,
  });

  const isLoggedIn = !!buyer;

  const logout = async () => {
    await logoutMutation.mutateAsync();
    router.push('/login/buyer');
  };

  return {
    buyer,
    isLoggedIn,
    isLoading,
    logout,
  };
};
