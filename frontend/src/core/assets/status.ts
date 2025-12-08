
// 翻訳キーを管理する型と定数
// 型定義は src/models/auction.ts 等と共有するか、ここで import する
// 今回は単純化のため文字列リテラル型で定義するが、実運用的にはモデルからImport推奨

export type AuctionStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export const AUCTION_STATUS_KEYS: Record<AuctionStatus, string> = {
    scheduled: 'AuctionStatus.scheduled',
    in_progress: 'AuctionStatus.in_progress',
    completed: 'AuctionStatus.completed',
    cancelled: 'AuctionStatus.cancelled',
};

// ItemStatus は Pending 以外が動的な可能性があるが、既知のもの定義
export type ItemStatus = 'Pending' | 'Sold' | 'Unsold' | 'Bidding' | string;

export const ITEM_STATUS_KEYS: Record<string, string> = {
    Pending: 'ItemStatus.Pending',
    Sold: 'ItemStatus.Sold',
    Unsold: 'ItemStatus.Unsold',
    Bidding: 'ItemStatus.Bidding',
};

// 従来の変換関数は廃止し、コンポーネント側で useTranslations() を使用する形に移行する
