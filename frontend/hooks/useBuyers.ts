import { useState } from 'react';

interface RegisterBuyerParams {
    name: string;
}

export function useRegisterBuyer() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const registerBuyer = async (params: RegisterBuyerParams): Promise<boolean> => {
        setIsLoading(true);
        setError(null);

        try {
            const res = await fetch('/api/buyers', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: params.name }),
            });

            if (res.ok) {
                setIsLoading(false);
                return true;
            } else {
                setError('Failed to register buyer');
                setIsLoading(false);
                return false;
            }
        } catch (err) {
            setError('Error registering buyer');
            setIsLoading(false);
            return false;
        }
    };

    return {
        registerBuyer,
        isLoading,
        error,
    };
}
