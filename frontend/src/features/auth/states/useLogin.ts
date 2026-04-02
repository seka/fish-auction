import { useAdminLoginMutation } from '../queries/useAuth';

export const useAdminLogin = () => {
  const mutation = useAdminLoginMutation();

  return {
    login: (email: string, password: string) => mutation.mutateAsync({ email, password }),
    isLoading: mutation.isPending,
    error: mutation.error,
  };
};
