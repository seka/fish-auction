import { useMutation } from '@tanstack/react-query';
import { submitBid } from '@/src/data/api/bid';

export const useBidMutation = () => {
  const mutation = useMutation({
    mutationFn: submitBid,
  });

  return {
    submitBid: mutation.mutateAsync,
    isLoading: mutation.isPending,
  };
};
