'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { getAuctions } from '@/src/api/auction';
import { getVenues } from '@/src/api/venue';

const usePublicVenues = () => {
    const { data: venues } = useQuery({
        queryKey: ['public_venues'],
        queryFn: getVenues,
    });
    return { venues };
};

export default function AuctionsListPage() {
    // Fetch all auctions
    const { data: allAuctions, isLoading } = useQuery({
        queryKey: ['public_auctions_list'],
        queryFn: () => getAuctions(),
    });

    const { venues } = usePublicVenues();

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <div className="text-xl text-gray-600">読み込み中...</div>
            </div>
        );
    }

    // Filter for active auctions (Scheduled or In Progress)
    const activeAuctions = allAuctions?.filter(a =>
        a.status === 'scheduled' || a.status === 'in_progress'
    ).sort((a, b) => {
        // Sort: In Progress first, then by date/time
        if (a.status === 'in_progress' && b.status !== 'in_progress') return -1;
        if (a.status !== 'in_progress' && b.status === 'in_progress') return 1;
        return new Date(`${a.auction_date}T${a.start_time}`).getTime() - new Date(`${b.auction_date}T${b.start_time}`).getTime();
    }) || [];

    const getVenueName = (id: number) => venues?.find(v => v.id === id)?.name || `会場ID: ${id}`;

    return (
        <div className="min-h-screen bg-gray-50 p-8">
            <div className="max-w-5xl mx-auto">
                <div className="flex items-center justify-between mb-8">
                    <h1 className="text-3xl font-bold text-gray-900">開催中のセリ一覧</h1>
                    <Link href="/" className="text-indigo-600 hover:text-indigo-800 font-medium">
                        &larr; トップに戻る
                    </Link>
                </div>

                {activeAuctions.length === 0 ? (
                    <div className="bg-white rounded-xl shadow-sm p-12 text-center">
                        <p className="text-xl text-gray-500">現在開催予定のセリはありません</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        {activeAuctions.map((auction) => (
                            <Link
                                key={auction.id}
                                href={`/auctions/${auction.id}`}
                                className="block group"
                            >
                                <div className={`bg-white rounded-xl shadow-sm border-2 transition-all duration-200 p-6 hover:shadow-md ${auction.status === 'in_progress'
                                    ? 'border-orange-400 ring-4 ring-orange-50'
                                    : 'border-transparent hover:border-indigo-200'
                                    }`}>
                                    <div className="flex justify-between items-start mb-4">
                                        <div>
                                            <span className={`inline-block px-3 py-1 rounded-full text-sm font-bold mb-2 ${auction.status === 'in_progress'
                                                ? 'bg-orange-100 text-orange-700 animate-pulse'
                                                : 'bg-blue-100 text-blue-700'
                                                }`}>
                                                {auction.status === 'in_progress' ? '🔥 開催中' : '📅 開催予定'}
                                            </span>
                                            <h2 className="text-xl font-bold text-gray-900 group-hover:text-indigo-700 transition-colors">
                                                {getVenueName(auction.venue_id)}
                                            </h2>
                                        </div>
                                        <div className="text-right">
                                            <div className="text-2xl font-bold text-gray-900">
                                                {auction.start_time?.substring(0, 5)}
                                            </div>
                                            <div className="text-sm text-gray-500">
                                                {auction.auction_date}
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
                                        <span className="text-gray-600 text-sm">
                                            終了予定: {auction.end_time?.substring(0, 5)}
                                        </span>
                                        <span className="text-indigo-600 font-bold group-hover:translate-x-1 transition-transform flex items-center">
                                            会場へ入る <span className="ml-1">&rarr;</span>
                                        </span>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
