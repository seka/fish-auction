import { useState } from 'react';

interface SubmitBidParams {
    item_id: number;
    buyer_id: number;
    price: number;
}

export function useSubmitBid() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const submitBid = async (params: SubmitBidParams): Promise<boolean> => {
        setIsLoading(true);
        setError(null);

        try {
            const res = await fetch('/api/bid', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(params),
            });

            if (res.ok) {
                setIsLoading(false);
                return true;
            } else {
                setError('Failed to submit bid');
                setIsLoading(false);
                return false;
            }
        } catch (err) {
            setError('Error submitting bid');
            setIsLoading(false);
            return false;
        }
    };

    return {
        submitBid,
        isLoading,
        error,
    };
}
