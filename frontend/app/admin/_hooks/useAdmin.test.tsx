import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  useRegisterFisherman,
  useFishermen,
  useRegisterBuyer,
  useBuyers,
  useRegisterItem,
} from './useAdmin';
import {
  registerFisherman,
  getFishermen,
  registerBuyer,
  getBuyers,
  registerItem,
} from '@/src/api/admin';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactNode } from 'react';
import { itemKeys } from '@/src/hooks/item/keys';

// Mock API
vi.mock('@/src/api/admin', () => ({
  registerFisherman: vi.fn(),
  getFishermen: vi.fn(),
  registerBuyer: vi.fn(),
  getBuyers: vi.fn(),
  registerItem: vi.fn(),
}));

const createWrapper = () => {
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
  wrapper.displayName = 'QueryClientWrapper';
  return wrapper;
};

describe('useAdmin hooks', () => {
  let queryClient: QueryClient;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let invalidateQueriesSpy: any;

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = new QueryClient();
    invalidateQueriesSpy = vi.spyOn(queryClient, 'invalidateQueries');
  });

  describe('useRegisterFisherman', () => {
    it('should register a fisherman', async () => {
      vi.mocked(registerFisherman).mockResolvedValue(true);
      const { result } = renderHook(() => useRegisterFisherman(), { wrapper: createWrapper() });

      await result.current.registerFisherman({ name: 'Fisher 1' });
      expect(registerFisherman).toHaveBeenCalledWith('Fisher 1');
    });
  });

  describe('useRegisterBuyer', () => {
    it('should register a buyer', async () => {
      vi.mocked(registerBuyer).mockResolvedValue(true);
      const { result } = renderHook(() => useRegisterBuyer(), { wrapper: createWrapper() });

      await result.current.registerBuyer({
        name: 'Buyer 1',
        email: 'buyer@example.com',
        password: 'password123',
        organization: 'Org 1',
        contactInfo: 'Address 1',
      });
      expect(registerBuyer).toHaveBeenCalledWith({
        name: 'Buyer 1',
        email: 'buyer@example.com',
        password: 'password123',
        organization: 'Org 1',
        contactInfo: 'Address 1',
      });
    });
  });

  describe('useRegisterItem', () => {
    it('should register an item and invalidate queries', async () => {
      vi.mocked(registerItem).mockResolvedValue(true);

      // We need to use the same queryClient instance in the hook
      const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      );
      wrapper.displayName = 'RegisterItemWrapper';

      const { result } = renderHook(() => useRegisterItem(), { wrapper });

      await result.current.registerItem({
        auctionId: 1,
        fishType: 'Item 1',
        quantity: 10,
        unit: 'kg',
        fishermanId: 1,
      });

      expect(registerItem).toHaveBeenCalledWith(expect.objectContaining({ fishType: 'Item 1' }));

      await waitFor(() =>
        expect(invalidateQueriesSpy).toHaveBeenCalledWith({ queryKey: itemKeys.all }),
      );
    });
  });

  describe('useFishermen', () => {
    it('should fetch fishermen', async () => {
      const mockData = [{ id: 1, name: 'Fisher 1' }];
      vi.mocked(getFishermen).mockResolvedValue(mockData);

      const { result } = renderHook(() => useFishermen(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.fishermen).toEqual(mockData));
      expect(getFishermen).toHaveBeenCalled();
    });
  });

  describe('useBuyers', () => {
    it('should fetch buyers', async () => {
      const mockData = [{ id: 1, name: 'Buyer 1' }];
      vi.mocked(getBuyers).mockResolvedValue(mockData);

      const { result } = renderHook(() => useBuyers(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.buyers).toEqual(mockData));
      expect(getBuyers).toHaveBeenCalled();
    });
  });
});
