import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { usePublicVenues } from './usePublicVenues';
import { getVenues } from '@/src/api/venue';

// Mocks
vi.mock('@/src/api/venue', () => ({
    getVenues: vi.fn(),
}));

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client= { queryClient } > { children } </QueryClientProvider>
);

describe('usePublicVenues', () => {
    it('fetches venues', async () => {
        const mockVenues = [{ id: 1, name: 'Venue A' }];
        (getVenues as any).mockResolvedValue(mockVenues);

        const { result } = renderHook(() => usePublicVenues(), { wrapper });

        await waitFor(() => {
            expect(result.current.venues).toEqual(mockVenues);
        });
    });
});
