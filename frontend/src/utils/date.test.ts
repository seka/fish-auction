import { describe, it, expect } from 'vitest';
import { formatDate, formatDateTimeForInput, toJSTISOString } from './date';

describe('utils/date', () => {
  describe('formatDate', () => {
    it('UTC 日時を JST の YYYY/MM/DD 形式に変換する', () => {
      // UTC 2026-03-15T15:00:00Z = JST 2026-03-16T00:00:00+09:00
      const date = new Date('2026-03-15T15:00:00Z');
      expect(formatDate(date)).toBe('2026/03/16');
    });

    it('日付変更前の UTC 時刻は JST 同日として変換する', () => {
      // UTC 2026-03-15T00:00:00Z = JST 2026-03-15T09:00:00+09:00
      const date = new Date('2026-03-15T00:00:00Z');
      expect(formatDate(date)).toBe('2026/03/15');
    });

    it('月・日が1桁の場合ゼロ埋めする', () => {
      // UTC 2026-01-01T00:00:00Z = JST 2026-01-01T09:00:00+09:00
      const date = new Date('2026-01-01T00:00:00Z');
      expect(formatDate(date)).toBe('2026/01/01');
    });
  });

  describe('formatDateTimeForInput', () => {
    it('UTC 日時を JST の datetime-local 形式に変換する', () => {
      // UTC 2026-03-15T00:00:00Z = JST 2026-03-15T09:00
      const date = new Date('2026-03-15T00:00:00Z');
      expect(formatDateTimeForInput(date)).toBe('2026-03-15T09:00');
    });

    it('UTC 深夜は JST 翌日として変換する', () => {
      // UTC 2026-03-15T15:30:00Z = JST 2026-03-16T00:30
      const date = new Date('2026-03-15T15:30:00Z');
      expect(formatDateTimeForInput(date)).toBe('2026-03-16T00:30');
    });

    it('時・分が1桁の場合ゼロ埋めする', () => {
      // UTC 2026-03-15T00:05:00Z = JST 2026-03-15T09:05
      const date = new Date('2026-03-15T00:05:00Z');
      expect(formatDateTimeForInput(date)).toBe('2026-03-15T09:05');
    });
  });

  describe('toJSTISOString', () => {
    it('datetime-local 値を RFC3339 JST 文字列に変換する', () => {
      expect(toJSTISOString('2026-03-15T09:00')).toBe('2026-03-15T09:00:00+09:00');
    });

    it('深夜0時をまたぐ値も正しく変換する', () => {
      expect(toJSTISOString('2026-03-16T00:30')).toBe('2026-03-16T00:30:00+09:00');
    });

    it('分が00の場合も正しく変換する', () => {
      expect(toJSTISOString('2026-12-31T23:59')).toBe('2026-12-31T23:59:00+09:00');
    });
  });
});
