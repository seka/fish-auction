import { useMutation } from '@tanstack/react-query';

interface RegisterBuyerParams {
    name: string;
}

export function useRegisterBuyer() {
    const { mutateAsync: registerBuyer, isPending: isLoading, error } = useMutation({
        mutationFn: async (params: RegisterBuyerParams) => {
            const res = await fetch('/api/buyers', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: params.name }),
            });

            if (!res.ok) {
                throw new Error('Failed to register buyer');
            }
            return true;
        },
    });

    return {
        registerBuyer: async (params: RegisterBuyerParams) => {
            try {
                return await registerBuyer(params);
            } catch (e) {
                return false;
            }
        },
        isLoading,
        error: error ? (error as Error).message : null,
    };
}
