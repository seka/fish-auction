import { describe, it, expect } from 'vitest';
import { getPasswordComplexitySchema } from './fields/password';
import { getResetPasswordSchema } from './password';
import { tMock } from './test-utils';

const t = tMock;
const passwordComplexitySchema = getPasswordComplexitySchema(t);
const resetPasswordSchema = getResetPasswordSchema(t);

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
      expect(result.error.issues[0].message).toBe('key:required(field:key:field_name.password)');
    }
  });

  it('should reject passwords shorter than 8 characters', () => {
    const result = passwordComplexitySchema.safeParse('Abc1');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('key:password_too_short(min:8)');
    }
  });

  it('should reject passwords without uppercase letters', () => {
    const result = passwordComplexitySchema.safeParse('abcdefgh1');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('key:password_uppercase');
    }
  });

  it('should reject passwords without lowercase letters', () => {
    const result = passwordComplexitySchema.safeParse('ABCDEFGH1');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('key:password_lowercase');
    }
  });

  it('should reject passwords without digits', () => {
    const result = passwordComplexitySchema.safeParse('Abcdefgh');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('key:password_number');
    }
  });

  it('should reject passwords with non-printable ASCII characters', () => {
    const result = passwordComplexitySchema.safeParse('Abcdefgh1\n');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('key:password_invalid_chars');
    }
  });

  it('should reject passwords with multi-byte characters', () => {
    const result = passwordComplexitySchema.safeParse('Abcdefgh1あ');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe('key:password_invalid_chars');
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
      expect(result.error.issues[0].message).toBe('key:password_too_long(max:72)');
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
      expect(result.error.issues[0].message).toBe('key:password_mismatch');
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
      expect(result.error.issues[0].message).toBe('key:password_too_short(min:8)');
    }
  });
});
