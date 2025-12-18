import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuctionData } from './useAuctionData';
import { getAuction, getAuctionItems } from '@/src/api/auction';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactNode } from 'react';

vi.mock('@/src/api/auction', () => ({
    getAuction: vi.fn(),
    getAuctionItems: vi.fn(),
}));

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

const wrapper = ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);

describe('useAuctionData', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        queryClient.clear();
    });

    it('fetches auction and items', async () => {
        const mockAuction = { id: 1, status: 'scheduled' };
        const mockItems = [{ id: 101, name: 'Tuna' }];

        (getAuction as any).mockResolvedValueOnce(mockAuction);
        (getAuctionItems as any).mockResolvedValueOnce(mockItems);

        const { result } = renderHook(() => useAuctionData(1), { wrapper });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.auction).toEqual(mockAuction);
        expect(result.current.items).toEqual(mockItems);
    });
});
