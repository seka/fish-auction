import { z } from 'zod';

export const bidSchema = z.object({
    price: z.string()
        .min(1, '価格を入力してください')
        .refine((val) => !isNaN(Number(val)) && Number(val) > 0, {
            message: '価格は正の数値で入力してください',
        }),
});

export type BidFormData = z.infer<typeof bidSchema>;
