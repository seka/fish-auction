import { useMutation } from '@tanstack/react-query';
import { login } from '../api';

export const useLogin = () => {
    const mutation = useMutation({
        mutationFn: (password: string) => login(password),
    });

    return {
        login: mutation.mutateAsync,
        isLoading: mutation.isPending,
        error: mutation.error,
    };
};
