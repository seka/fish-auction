import { logoutBuyer } from '@/src/data/api/buyer_auth';
import { useQueryClient } from '@tanstack/react-query';
import { authKeys } from '@/src/data/queries/auth/keys';
import { useRouter } from 'next/navigation';

/**
 * マイページ用認証関連アクション
 */
export const useMyPageAuth = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  const handleLogout = async () => {
    const success = await logoutBuyer();
    if (success) {
      await queryClient.invalidateQueries({ queryKey: authKeys.me() });
      router.push('/login/buyer');
    }
    return success;
  };

  return {
    handleLogout,
  };
};
