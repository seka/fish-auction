import { z } from 'zod';
import { getPasswordComplexitySchema, ValidationT } from './fields/password';

export const getResetPasswordSchema = (t: ValidationT) =>
  z
    .object({
      new_password: getPasswordComplexitySchema(t),
      confirm_password: z.string().min(1, t('required', { field: t('Auth.ResetPassword.label_confirm_password') })),
    })
    .refine((data) => data.new_password === data.confirm_password, {
      message: t('password_mismatch'),
      path: ['confirm_password'],
    });
