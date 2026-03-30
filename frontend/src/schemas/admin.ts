import { z } from 'zod';
import { getPasswordComplexitySchema, ValidationT } from './fields/password';
import { getEmailSchema } from './fields/email';
import { getQuantitySchema } from './fields/quantity';

export const getFishermanSchema = (t: ValidationT) =>
  z.object({
    name: z.string().min(1, t('required', { field: t('field_name.fisherman_name') })),
  });

export const getBuyerSchema = (t: ValidationT) =>
  z.object({
    name: z.string().min(1, t('required', { field: t('field_name.buyer_name') })),
    email: getEmailSchema(t),
    password: getPasswordComplexitySchema(t),
    organization: z.string().min(1, t('required', { field: t('field_name.organization') })),
    contactInfo: z.string().min(1, t('required', { field: t('field_name.contact_info') })),
  });

export const getItemSchema = (t: ValidationT) =>
  z.object({
    auctionId: z.string().min(1, t('select_required', { field: t('Items.auction') })),
    fishermanId: z.string().min(1, t('select_required', { field: t('Items.fisherman') })),
    fishType: z.string().min(1, t('required', { field: t('Items.fish_type') })),
    quantity: getQuantitySchema(t),
    unit: z.string().min(1, t('required', { field: t('Items.unit') })),
  });

export type FishermanFormData = z.infer<ReturnType<typeof getFishermanSchema>>;
export type BuyerFormData = z.infer<ReturnType<typeof getBuyerSchema>>;
export type ItemFormData = z.infer<ReturnType<typeof getItemSchema>>;
