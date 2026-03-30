import { z } from 'zod';
import { ValidationT } from './fields/password';

export const getLoginSchema = (t: ValidationT) =>
  z.object({
    email: z.string().email(t('invalid_email')),
    password: z.string().min(1, t('required', { field: t('field_name.password') })),
  });

export type LoginFormData = z.infer<ReturnType<typeof getLoginSchema>>;
