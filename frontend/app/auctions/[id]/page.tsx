'use client';

import { useParams } from 'next/navigation';
import { AuctionDetailContainer } from './components/AuctionDetailContainer';

export default function AuctionDetailPage() {
  const params = useParams();
  const id = params.id as string;

  return <AuctionDetailContainer id={id} />;
}
