import { z } from 'zod';

export const buyerSignupSchema = z.object({
    name: z.string().min(1, '名前を入力してください'),
    email: z.string().email('有効なメールアドレスを入力してください'),
    password: z.string().min(6, 'パスワードは6文字以上で入力してください'),
    organization: z.string().min(1, '組織名を入力してください'),
    contact_info: z.string().min(1, '連絡先を入力してください'),
});

export type BuyerSignupFormData = z.infer<typeof buyerSignupSchema>;

export const buyerLoginSchema = z.object({
    email: z.string().email('有効なメールアドレスを入力してください'),
    password: z.string().min(1, 'パスワードを入力してください'),
});

export type BuyerLoginFormData = z.infer<typeof buyerLoginSchema>;
