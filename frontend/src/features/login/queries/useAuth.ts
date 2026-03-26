import { useMutation } from '@tanstack/react-query';
import { login } from '@/src/data/api/auth';

/**
 * ログインミューテーションフック
 */
export const useLoginMutation = () => {
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      login(email, password),
  });
};
