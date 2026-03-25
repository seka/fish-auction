'use client';

import { use } from 'react';
import { AuctionDetailContainer } from '@/src/features/auctions';

export default function AuctionDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);

  return <AuctionDetailContainer id={id} />;
}
