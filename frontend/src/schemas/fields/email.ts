import { z } from 'zod';
import { ValidationT } from './password';

export const getEmailSchema = (t: ValidationT) =>
  z.string().email(t('invalid_email'));
