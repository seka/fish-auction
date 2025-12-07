import { z } from 'zod';

export const venueSchema = z.object({
    id: z.number().optional(),
    name: z.string().min(1, '会場名を入力してください'),
    location: z.string().optional(),
    description: z.string().optional(),
    createdAt: z.string().optional(),
});

export const auctionSchema = z.object({
    id: z.number().optional(),
    venueId: z.number().min(1, '会場を選択してください'),
    auctionDate: z.string().min(1, '開催日を入力してください'),
    startTime: z.string().optional(),
    endTime: z.string().optional(),
    status: z.enum(['scheduled', 'in_progress', 'completed', 'cancelled']).optional(),
    createdAt: z.string().optional(),
    updatedAt: z.string().optional(),
});

export const bidSchema = z.object({
    price: z.string()
        .min(1, '価格を入力してください')
        .refine((val) => !isNaN(Number(val)) && Number(val) > 0, {
            message: '価格は正の数値で入力してください',
        }),
});

export type VenueFormData = z.infer<typeof venueSchema>;
export type AuctionFormData = z.infer<typeof auctionSchema>;
export type BidFormData = z.infer<typeof bidSchema>;
