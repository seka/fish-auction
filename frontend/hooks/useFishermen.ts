import { useState } from 'react';

interface RegisterFishermanParams {
    name: string;
}

export function useRegisterFisherman() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const registerFisherman = async (params: RegisterFishermanParams): Promise<boolean> => {
        setIsLoading(true);
        setError(null);

        try {
            const res = await fetch('/api/fishermen', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: params.name }),
            });

            if (res.ok) {
                setIsLoading(false);
                return true;
            } else {
                setError('Failed to register fisherman');
                setIsLoading(false);
                return false;
            }
        } catch (err) {
            setError('Error registering fisherman');
            setIsLoading(false);
            return false;
        }
    };

    return {
        registerFisherman,
        isLoading,
        error,
    };
}
