import { z } from 'zod';
import { ValidationT } from './password';

export const getQuantitySchema = (t: ValidationT) =>
  z
    .string()
    .min(1, t('required', { field: t('Items.quantity') }))
    .refine((val) => !isNaN(Number(val)) && Number(val) > 0, {
      message: t('positive_number'),
    });
