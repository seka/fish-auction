import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useRegisterFisherman, useFishermen, useRegisterBuyer, useBuyers, useRegisterItem } from './useAdmin';
import { registerFisherman, getFishermen, registerBuyer, getBuyers, registerItem } from '@/src/api/admin';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactNode } from 'react';

// Mock API
vi.mock('@/src/api/admin', () => ({
    registerFisherman: vi.fn(),
    getFishermen: vi.fn(),
    registerBuyer: vi.fn(),
    registerItem: vi.fn(),
    getBuyers: vi.fn(),
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

describe('useAdmin Hooks', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        queryClient.clear();
    });

    describe('useRegisterFisherman', () => {
        it('calls registerFisherman API on mutate', async () => {
            (registerFisherman as any).mockResolvedValueOnce({ id: 1, name: 'Fisherman 1' });

            const { result } = renderHook(() => useRegisterFisherman(), { wrapper });

            await result.current.registerFisherman({ name: 'Fisherman 1' });

            expect(registerFisherman).toHaveBeenCalledWith('Fisherman 1');
        });
    });

    describe('useFishermen', () => {
        it('fetches fishermen', async () => {
            const mockData = [{ id: 1, name: 'Fisherman 1' }];
            (getFishermen as any).mockResolvedValueOnce(mockData);

            const { result } = renderHook(() => useFishermen(), { wrapper });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.fishermen).toEqual(mockData);
        });
    });

    describe('useRegisterBuyer', () => {
        it('calls registerBuyer API on mutate', async () => {
            (registerBuyer as any).mockResolvedValueOnce({ id: 2, name: 'Buyer 1' });

            const { result } = renderHook(() => useRegisterBuyer(), { wrapper });

            await result.current.registerBuyer({ name: 'Buyer 1' } as any);

            expect(registerBuyer).toHaveBeenCalledWith({ name: 'Buyer 1' });
        });
    });

    describe('useBuyers', () => {
        it('fetches buyers', async () => {
            const mockData = [{ id: 2, name: 'Buyer 1' }];
            (getBuyers as any).mockResolvedValueOnce(mockData);

            const { result } = renderHook(() => useBuyers(), { wrapper });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.buyers).toEqual(mockData);
        });
    });

    describe('useRegisterItem', () => {
        it('calls registerItem API on mutate and invalidates queries', async () => {
            (registerItem as any).mockResolvedValueOnce({ id: 100, name: 'Item 1' });
            const invalidateQueriesSpy = vi.spyOn(queryClient, 'invalidateQueries');

            const { result } = renderHook(() => useRegisterItem(), { wrapper });

            await result.current.registerItem({ name: 'Item 1', start_price: 100 } as any);

            expect(registerItem).toHaveBeenCalledWith(expect.objectContaining({ name: 'Item 1' }));

            await waitFor(() => expect(invalidateQueriesSpy).toHaveBeenCalledWith({ queryKey: ['items'] }));
        });
    });
});
