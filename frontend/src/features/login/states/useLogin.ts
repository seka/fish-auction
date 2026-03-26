import { useLoginMutation } from '../queries/useAuth';

export const useLogin = () => {
  const mutation = useLoginMutation();

  return {
    login: (email: string, password: string) => mutation.mutateAsync({ email, password }),
    isLoading: mutation.isPending,
    error: mutation.error,
  };
};
