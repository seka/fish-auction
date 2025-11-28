import { z } from 'zod';

export const loginSchema = z.object({
    password: z.string().min(1, 'パスワードを入力してください'),
});

export type LoginFormData = z.infer<typeof loginSchema>;
