import { z } from 'zod';
import { useTranslations } from 'next-intl';

export type ValidationT = ReturnType<typeof useTranslations<'Validation'>>;

export const getPasswordComplexitySchema = (t: ValidationT) =>
  z
    .string()
    .min(1, t('required', { field: t('field_name.password') }))
    .min(8, t('password_too_short', { min: 8 }))
    .max(72, t('password_too_long', { max: 72 }))
    .regex(/[A-Z]/, t('password_uppercase'))
    .regex(/[a-z]/, t('password_lowercase'))
    .regex(/[0-9]/, t('password_number'))
    .regex(/^[\x20-\x7E]*$/, t('password_invalid_chars'));
