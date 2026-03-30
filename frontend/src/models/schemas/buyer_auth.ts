import { z } from 'zod';
import { ValidationT } from './fields/password';

export const getBuyerLoginSchema = (t: ValidationT) =>
  z.object({
    email: z.email(t('invalid_email')),
    password: z.string().min(1, t('required', { field: t('field_name.password') })),
  });

export type BuyerLoginFormData = z.infer<ReturnType<typeof getBuyerLoginSchema>>;
