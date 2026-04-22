/**
 * 価格に応じた最小入札単位を取得する
 */
export const selectMinimumBidIncrement = (currentPrice: number): number => {
  if (currentPrice < 1000) return 100;
  if (currentPrice < 10000) return 500;
  if (currentPrice < 100000) return 1000;
  return 5000;
};

/**
 * 次の最低入札額を計算する
 */
export const selectNextMinimumBid = (currentHighestBid: number): number => {
  return currentHighestBid + selectMinimumBidIncrement(currentHighestBid);
};
