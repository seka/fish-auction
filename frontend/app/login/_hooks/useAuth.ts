import { useMutation } from '@tanstack/react-query';
import { login } from '@/src/api/auth';

export const useLogin = () => {
    const mutation = useMutation({
        mutationFn: ({ email, password }: { email: string; password: string }) => login(email, password),
    });

    return {
        login: (email: string, password: string) => mutation.mutateAsync({ email, password }),
        isLoading: mutation.isPending,
        error: mutation.error,
    };
};
