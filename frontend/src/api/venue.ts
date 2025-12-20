import { apiClient } from '@/src/core/api/client';
import { Venue } from '@/src/models/venue';
import { VenueFormData } from '@/src/models/schemas/auction';

export const createVenue = async (data: VenueFormData): Promise<Venue> => {
    return apiClient.post<Venue>('/api/admin/venues', data);
};

export const getVenues = async (): Promise<Venue[]> => {
    return apiClient.get<Venue[]>('/api/venues');
};

export const getVenue = async (id: number): Promise<Venue> => {
    return apiClient.get<Venue>(`/api/venues/${id}`);
};

export const updateVenue = async (id: number, data: VenueFormData): Promise<void> => {
    return apiClient.put<void>(`/api/admin/venues/${id}`, data);
};

export const deleteVenue = async (id: number): Promise<void> => {
    return apiClient.delete(`/api/admin/venues/${id}`);
};
