import { z } from 'zod';

export const fishermanSchema = z.object({
    name: z.string().min(1, '漁師名を入力してください'),
});

export const buyerSchema = z.object({
    name: z.string().min(1, '中買人名を入力してください'),
});

export const itemSchema = z.object({
    fishermanId: z.string().min(1, '漁師を選択してください'),
    fishType: z.string().min(1, '魚種を入力してください'),
    quantity: z.string()
        .min(1, '数量を入力してください')
        .refine((val) => !isNaN(Number(val)) && Number(val) > 0, {
            message: '数量は正の数値で入力してください',
        }),
    unit: z.string().min(1, '単位を入力してください'),
});

export type FishermanFormData = z.infer<typeof fishermanSchema>;
export type BuyerFormData = z.infer<typeof buyerSchema>;
export type ItemFormData = z.infer<typeof itemSchema>;
