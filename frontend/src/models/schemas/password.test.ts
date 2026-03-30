import { describe, it, expect } from 'vitest';
import { getPasswordComplexitySchema } from './fields/password';
import { getResetPasswordSchema } from './password';

const tMock = ((key: string, values?: any) => {
  if (key === 'required') return `${values?.field}を入力してください`;
  if (key === 'password_too_short') return `パスワードは${values?.min}文字以上である必要があります`;
  if (key === 'password_too_long') return `パスワードは${values?.max}文字以内である必要があります`;
  if (key === 'password_uppercase') return '大文字を1文字以上含めてください';
  if (key === 'password_lowercase') return '小文字を1文字以上含めてください';
  if (key === 'password_number') return '数字を1文字以上含めてください';
  if (key === 'password_invalid_chars') return '使用できない文字が含まれています';
  if (key === 'password_mismatch') return 'パスワードが一致しません';
  if (key === 'field_name.password') return 'パスワード';
  if (key === 'Auth.ResetPassword.label_confirm_password') return '確認用パスワード';
  return key;
}) as any;

const passwordComplexitySchema = getPasswordComplexitySchema(tMock);
const resetPasswordSchema = getResetPasswordSchema(tMock);

describe('passwordComplexitySchema', () => {
  it('should accept valid passwords', () => {
    const validPasswords = [
      'Abcdefgh1',
      'Password123!',
      '1234abcdA',
      'A1b2C3d4',
      'LongPasswordWith1Upper1Lower1Digit',
    ];
    validPasswords.forEach((pwd) => {
      const result = passwordComplexitySchema.safeParse(pwd);
      expect(result.success).toBe(true);
    });
  });

  it('should reject empty passwords', () => {
    const result = passwordComplexitySchema.safeParse('');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('パスワードを入力してください');
    }
  });

  it('should reject passwords shorter than 8 characters', () => {
    const result = passwordComplexitySchema.safeParse('Abc1');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('パスワードは8文字以上である必要があります');
    }
  });

  it('should reject passwords without uppercase letters', () => {
    const result = passwordComplexitySchema.safeParse('abcdefgh1');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('大文字を1文字以上含めてください');
    }
  });

  it('should reject passwords without lowercase letters', () => {
    const result = passwordComplexitySchema.safeParse('ABCDEFGH1');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('小文字を1文字以上含めてください');
    }
  });

  it('should reject passwords without digits', () => {
    const result = passwordComplexitySchema.safeParse('Abcdefgh');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('数字を1文字以上含めてください');
    }
  });

  it('should reject passwords with non-printable ASCII characters', () => {
    const result = passwordComplexitySchema.safeParse('Abcdefgh1\n');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('使用できない文字が含まれています');
    }
  });

  it('should reject passwords with multi-byte characters', () => {
    const result = passwordComplexitySchema.safeParse('Abcdefgh1あ');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('使用できない文字が含まれています');
    }
  });

  it('should accept 72 characters password', () => {
    const longPwd = 'Ab1' + 'a'.repeat(68) + '1'; // 3 + 68 + 1 = 72
    const result = passwordComplexitySchema.safeParse(longPwd);
    expect(result.success).toBe(true);
  });

  it('should reject 73 characters password', () => {
    const tooLongPwd = 'Ab1' + 'a'.repeat(69) + '1'; // 3 + 69 + 1 = 73
    const result = passwordComplexitySchema.safeParse(tooLongPwd);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('パスワードは72文字以内である必要があります');
    }
  });
});

describe('resetPasswordSchema', () => {
  it('should accept when passwords match', () => {
    const data = {
      new_password: 'Password123!',
      confirm_password: 'Password123!',
    };
    const result = resetPasswordSchema.safeParse(data);
    expect(result.success).toBe(true);
  });

  it('should reject when passwords mismatch', () => {
    const data = {
      new_password: 'Password123!',
      confirm_password: 'Different123!',
    };
    const result = resetPasswordSchema.safeParse(data);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('パスワードが一致しません');
      expect(result.error.issues[0].path).toContain('confirm_password');
    }
  });

  it('should reject when new_password is invalid', () => {
    const data = {
      new_password: 'short',
      confirm_password: 'short',
    };
    const result = resetPasswordSchema.safeParse(data);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('パスワードは8文字以上である必要があります');
    }
  });
});
