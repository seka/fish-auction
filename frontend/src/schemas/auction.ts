import { z } from 'zod';
import { ValidationT } from './fields/password';
import { getPriceSchema } from './fields/price';

export const getVenueSchema = (t: ValidationT) =>
  z.object({
    id: z.number().optional(),
    name: z.string().min(1, t('required', { field: t('Admin.Venues.name') })),
    location: z.string().optional(),
    description: z.string().optional(),
    createdAt: z.string().optional(),
  });

export const getAuctionSchema = (t: ValidationT) =>
  z.object({
    id: z.number().optional(),
    venueId: z
      .union([z.string(), z.number()])
      .transform(Number)
      .refine((n) => n >= 1, t('select_required', { field: t('Admin.Auctions.venue') })),
    startAt: z.string().min(1, t('required', { field: t('Admin.Auctions.start_time') })),
    endAt: z.string().optional(),
    status: z.enum(['scheduled', 'in_progress', 'completed', 'cancelled']).optional(),
    createdAt: z.string().optional(),
    updatedAt: z.string().optional(),
  });

export const getBidSchema = (t: ValidationT) =>
  z.object({
    price: getPriceSchema(t),
  });

export type VenueFormData = z.infer<ReturnType<typeof getVenueSchema>>;
export type AuctionFormInput = z.input<ReturnType<typeof getAuctionSchema>>;
export type AuctionFormData = z.output<ReturnType<typeof getAuctionSchema>>;
export type BidFormData = z.infer<ReturnType<typeof getBidSchema>>;
