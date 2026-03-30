import { z } from 'zod';
import { ValidationT } from './password';

export const getPriceSchema = (t: ValidationT) =>
  z
    .string()
    .min(1, t('required', { field: t('field_name.price') }))
    .refine((val) => !isNaN(Number(val)) && Number(val) > 0, {
      message: t('positive_number'),
    });
