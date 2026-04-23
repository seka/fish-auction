/**
 * Dateオブジェクトを YYYY/MM/DD 形式の文字列に変換する（日本時間固定）
 */
export const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    timeZone: 'Asia/Tokyo',
  }).format(date);
};

/**
 * datetime-local 入力値 "YYYY-MM-DDTHH:MM" を RFC3339 JST 文字列 "YYYY-MM-DDTHH:MM:00+09:00" に変換する
 */
export const toJSTISOString = (localDatetime: string): string => `${localDatetime}:00+09:00`;

/**
 * Dateオブジェクトを datetime-local 入力形式 "YYYY-MM-DDTHH:mm" に変換する（日本時間固定）
 */
export const formatDateTimeForInput = (date: Date): string => {
  const parts = new Intl.DateTimeFormat('ja-JP', {
    timeZone: 'Asia/Tokyo',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).formatToParts(date);

  const getPart = (type: string) => parts.find((p) => p.type === type)?.value || '00';

  const y = getPart('year');
  const m = getPart('month');
  const d = getPart('day');
  const hh = getPart('hour');
  const mm = getPart('minute');

  return `${y}-${m}-${d}T${hh}:${mm}`;
};
