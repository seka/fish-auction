import { z } from 'zod';
import { passwordComplexitySchema } from './fields/password';

export const resetPasswordSchema = z
  .object({
    new_password: passwordComplexitySchema,
    confirm_password: z.string().min(1, '確認用パスワードを入力してください'),
  })
  .refine((data) => data.new_password === data.confirm_password, {
    message: 'パスワードが一致しません',
    path: ['confirm_password'],
  });
