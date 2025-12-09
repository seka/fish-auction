import { useQuery } from '@tanstack/react-query';
import { getVenues } from '@/src/api/venue';

export const usePublicVenues = () => {
    const { data: venues } = useQuery({
        queryKey: ['public_venues'],
        queryFn: getVenues,
    });
    return { venues };
};
