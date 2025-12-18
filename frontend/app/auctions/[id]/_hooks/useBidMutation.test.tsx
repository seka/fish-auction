import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { useBidMutation } from './useBidMutation';
import { submitBid } from '@/src/api/bid';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mocks
vi.mock('@/src/api/bid', () => ({
    submitBid: vi.fn(),
}));

const queryClient = new QueryClient();
const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient} > {children} </QueryClientProvider>
);

describe('useBidMutation', () => {
    it('submits bid', async () => {
        (submitBid as any).mockResolvedValue(true);

        const { result } = renderHook(() => useBidMutation(), { wrapper });

        await result.current.submitBid({ auctionId: 1, itemId: 101, price: 5000 });

        expect(submitBid).toHaveBeenCalledWith(
            { auctionId: 1, itemId: 101, price: 5000 },
            expect.anything()
        );
    });
});
