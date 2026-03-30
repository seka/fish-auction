import { ValidationT } from './fields/password';

/**
 * テスト用の翻訳関数モック。
 * 翻訳キーと、プレースホルダーとして渡された値を文字列として返す。
 */
export const tMock: ValidationT = ((key: string, values?: Record<string, unknown>) => {
  if (key === 'required') return `${values?.field}を入力してください`;
  if (key === 'select_required') return `${values?.field}を選択してください`;
  if (key === 'invalid_email') return '有効なメールアドレスを入力してください';
  if (key === 'password_too_short') return `パスワードは${values?.min}文字以上である必要があります`;
  if (key === 'password_too_long') return `パスワードは${values?.max}文字以内である必要があります`;
  if (key === 'password_uppercase') return '大文字を1文字以上含めてください';
  if (key === 'password_lowercase') return '小文字を1文字以上含めてください';
  if (key === 'password_number') return '数字を1文字以上含めてください';
  if (key === 'password_invalid_chars') return '使用できない文字が含まれています';
  if (key === 'password_mismatch') return 'パスワードが一致しません';
  if (key === 'positive_number') return '正の数値を入力してください';
  
  // フィールド名などのマッピング
  const fieldNames: Record<string, string> = {
    'Admin.Venues.name': '会場名',
    'Admin.Auctions.venue': '会場',
    'Admin.Auctions.date': '開催日',
    'field_name.fisherman_name': '漁師名',
    'field_name.buyer_name': '中買人名',
    'field_name.organization': '団体・組織',
    'field_name.contact_info': '連絡先',
    'Items.auction': 'セリ',
    'Items.fisherman': '漁師',
    'Items.fish_type': '魚種',
    'Items.quantity': '数量',
    'Items.unit': '単位',
    'field_name.password': 'パスワード',
    'field_name.price': '価格',
    'field_name.email': 'メールアドレス',
    'Auth.ResetPassword.label_confirm_password': '確認用パスワード',
  };

  return fieldNames[key] || key;
}) as unknown as ValidationT;
