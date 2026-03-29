import { z } from 'zod';

export const passwordComplexitySchema = z
  .string()
  .min(1, 'パスワードを入力してください')
  .min(8, 'パスワードは8文字以上である必要があります')
  .max(72, 'パスワードは72文字以内である必要があります')
  .regex(/[A-Z]/, '大文字を1文字以上含めてください')
  .regex(/[a-z]/, '小文字を1文字以上含めてください')
  .regex(/[0-9]/, '数字を1文字以上含めてください')
  .regex(/^[\x20-\x7E]*$/, '使用できない文字が含まれています');
