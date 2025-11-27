import { useMutation } from '@tanstack/react-query';

interface RegisterFishermanParams {
    name: string;
}

export function useRegisterFisherman() {
    const { mutateAsync: registerFisherman, isPending: isLoading, error } = useMutation({
        mutationFn: async (params: RegisterFishermanParams) => {
            const res = await fetch('/api/fishermen', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: params.name }),
            });

            if (!res.ok) {
                throw new Error('Failed to register fisherman');
            }
            return true;
        },
    });

    return {
        registerFisherman: async (params: RegisterFishermanParams) => {
            try {
                return await registerFisherman(params);
            } catch (e) {
                return false;
            }
        },
        isLoading,
        error: error ? (error as Error).message : null,
    };
}
