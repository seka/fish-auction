import { useInvoiceQuery } from '@/src/data/queries/adminInvoice/useQuery';

/**
 * 管理画面用の請求書一覧クエリフック
 * src/data/queries/adminInvoice のベースフックをラップし、
 * 必要に応じて管理画面固有の加工・選択ロジックをここに追加します。
 */
export const useInvoices = () => {
  const query = useInvoiceQuery();
  
  // 今後、管理画面用にデータのフィルタリングや変換が必要になった場合は
  // TanStack Query の 'select' オプション等を追加してここで加工します。
  
  return query;
};
