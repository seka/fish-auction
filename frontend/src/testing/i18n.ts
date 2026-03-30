import { ValidationT } from '@schema/fields/password';

/**
 * テスト用の翻訳関数モック。
 * 翻訳キーと、プレースホルダーとして渡された値を「key:KEY(PARAM1:VAL1,...)」の形式で返す。
 * これにより、特定の言語（日本語等）に依存せず、正しいキーが呼ばれているかを検証できる。
 */
export const tMock: ValidationT = ((key: string, values?: Record<string, unknown>) => {
  const params = values
    ? Object.entries(values)
        .map(([k, v]) => `${k}:${v}`)
        .join(',')
    : '';

  return `key:${key}${params ? `(${params})` : ''}`;
}) as unknown as ValidationT;
