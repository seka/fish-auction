'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import Link from 'next/link';

import { loginBuyer } from '@/src/api/buyer_auth';
import { AuctionItem } from '@/src/models';
import { bidSchema, BidFormData } from '@/src/models/schemas/auction';
import { buyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { useAuctionData } from './_hooks/useAuctionData';
import { useBidMutation } from './_hooks/useBidMutation';
import { useAuth } from './_hooks/useAuth';
import { isAuctionActive, formatTime } from '@/src/utils/auction';

export default function AuctionRoomPage() {
    const params = useParams();
    const router = useRouter();
    const auctionId = Number(params.id);

    const [selectedItem, setSelectedItem] = useState<AuctionItem | null>(null);
    const [message, setMessage] = useState('');
    const [loginError, setLoginError] = useState('');

    const { auction, items, isLoading, refetchItems } = useAuctionData(auctionId);
    const { submitBid, isLoading: isBidLoading } = useBidMutation();
    const { isLoggedIn, isChecking } = useAuth();

    // Check if auction is active (within bidding hours)
    const auctionActive = auction ? isAuctionActive(auction) : false;

    const { register, handleSubmit, reset, formState: { errors } } = useForm<BidFormData>({
        resolver: zodResolver(bidSchema),
    });

    const { register: registerLogin, handleSubmit: handleSubmitLogin, formState: { errors: loginErrors, isSubmitting: isLoggingIn } } = useForm<BuyerLoginFormData>({
        resolver: zodResolver(buyerLoginSchema),
    });

    // Reset selected item if it disappears from list or status changes (optional)
    useEffect(() => {
        if (selectedItem) {
            const current = items.find(i => i.id === selectedItem.id);
            if (current && current.status !== selectedItem.status) {
                setSelectedItem(current);
            }
        }
    }, [items, selectedItem]);

    if (isNaN(auctionId)) {
        return <div>Invalid Auction ID</div>;
    }

    if (isChecking) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <div className="text-xl text-gray-600">読み込み中...</div>
            </div>
        );
    }

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <div className="text-xl text-gray-600">読み込み中...</div>
            </div>
        );
    }

    if (!auction) {
        return <div>Auction not found</div>;
    }

    const onSubmitLogin = async (data: BuyerLoginFormData) => {
        setLoginError('');
        const buyer = await loginBuyer(data);
        if (buyer) {
            // Reload page to update auth state and stay on current page
            window.location.reload();
        } else {
            setLoginError('メールアドレスまたはパスワードが間違っています');
        }
    };

    const onSubmitBid = async (data: BidFormData) => {
        if (!selectedItem) return;

        const success = await submitBid({
            item_id: selectedItem.id,
            buyer_id: 0, // Backend handles this from context
            price: parseInt(data.price),
        });

        if (success) {
            setMessage(`落札成功！ (${selectedItem.fish_type})`);
            setSelectedItem(null);
            reset();
            refetchItems();
            // Clear message after 3 seconds
            setTimeout(() => setMessage(''), 3000);
        } else {
            setMessage('入札に失敗しました');
        }
    };

    // Show login form if not logged in
    if (!isLoggedIn) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
                <div className="max-w-md w-full space-y-8">
                    <div>
                        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                            セリ会場へのログイン
                        </h2>
                        <p className="mt-2 text-center text-sm text-gray-600">
                            入札するにはログインが必要です
                        </p>
                    </div>
                    <form className="mt-8 space-y-6" onSubmit={handleSubmitLogin(onSubmitLogin)}>
                        <div className="rounded-md shadow-sm -space-y-px">
                            <div>
                                <label htmlFor="email" className="sr-only">メールアドレス</label>
                                <input
                                    id="email"
                                    type="email"
                                    {...registerLogin('email')}
                                    className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                    placeholder="メールアドレス"
                                />
                                {loginErrors.email && <p className="text-red-500 text-xs mt-1">{loginErrors.email.message}</p>}
                            </div>
                            <div>
                                <label htmlFor="password" className="sr-only">パスワード</label>
                                <input
                                    id="password"
                                    type="password"
                                    {...registerLogin('password')}
                                    className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                    placeholder="パスワード"
                                />
                                {loginErrors.password && <p className="text-red-500 text-xs mt-1">{loginErrors.password.message}</p>}
                            </div>
                        </div>

                        {loginError && <div className="text-red-500 text-sm text-center">{loginError}</div>}

                        <div>
                            <button
                                type="submit"
                                disabled={isLoggingIn}
                                className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                            >
                                ログイン
                            </button>
                        </div>
                        <div className="text-center">
                            <Link href="/signup" className="text-sm text-indigo-600 hover:text-indigo-500">
                                アカウントをお持ちでない方はこちら
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 p-4 md:p-8">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
                    <div>
                        <div className="flex items-center gap-3 mb-1">
                            <Link href="/auctions" className="text-gray-500 hover:text-gray-700">
                                &larr; 一覧へ
                            </Link>
                            <span className={`px-3 py-1 rounded-full text-sm font-bold ${auction.status === 'in_progress'
                                ? 'bg-orange-100 text-orange-700 animate-pulse'
                                : 'bg-blue-100 text-blue-700'
                                }`}>
                                {auction.status === 'in_progress' ? '🔥 開催中' : auction.status}
                            </span>
                        </div>
                        <h1 className="text-3xl font-bold text-gray-900">
                            セリ会場 #{auction.id}
                        </h1>
                        <p className="text-gray-500">
                            {auction.auction_date} {auction.start_time?.substring(0, 5)} - {auction.end_time?.substring(0, 5)}
                        </p>
                    </div>
                    <div className="text-right hidden md:block">
                        <p className="text-sm text-gray-500">自動更新中 (5秒)</p>
                    </div>
                </div>

                {message && (
                    <div className="bg-green-50 border-l-4 border-green-500 text-green-700 p-4 mb-6 rounded shadow-sm animate-bounce" role="alert">
                        <p className="font-bold">{message}</p>
                    </div>
                )}

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Item List */}
                    <div className="lg:col-span-2 space-y-4">
                        <h2 className="text-xl font-bold mb-4 text-gray-800 border-b pb-2">出品リスト</h2>
                        {items.length === 0 ? (
                            <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
                                <p className="text-gray-500">現在、出品されている商品はありません。</p>
                            </div>
                        ) : (
                            items.map((item) => (
                                <div
                                    key={item.id}
                                    className={`p-6 border rounded-xl cursor-pointer transition-all duration-200 ${selectedItem?.id === item.id
                                        ? 'border-orange-500 bg-orange-50 shadow-md transform scale-[1.01]'
                                        : 'bg-white hover:shadow-md border-gray-200'
                                        }`}
                                    onClick={() => setSelectedItem(item)}
                                >
                                    <div className="flex justify-between items-center">
                                        <div className="flex items-center space-x-4">
                                            <div className="bg-blue-100 text-blue-800 font-bold px-3 py-1 rounded text-xs">
                                                ID: {item.id}
                                            </div>
                                            <div>
                                                <h3 className="text-xl font-bold text-gray-900">{item.fish_type}</h3>
                                                <p className="text-gray-600 mt-1">
                                                    <span className="font-bold text-lg">{item.quantity}</span> {item.unit}
                                                    <span className="text-sm ml-2 text-gray-400">(漁師ID: {item.fisherman_id})</span>
                                                </p>
                                                {item.highest_bid && (
                                                    <p className="text-sm mt-1 text-orange-600 font-semibold">
                                                        現在の最高額: ¥{item.highest_bid.toLocaleString()}
                                                        {item.highest_bidder_name && (
                                                            <span className="ml-2 text-gray-600">({item.highest_bidder_name} さん)</span>
                                                        )}
                                                    </p>
                                                )}
                                            </div>
                                        </div>
                                        <span className={`px-4 py-2 rounded-full text-sm font-bold shadow-sm ${item.status === 'Pending'
                                            ? 'bg-green-100 text-green-800'
                                            : 'bg-gray-100 text-gray-600'
                                            }`}>
                                            {item.status === 'Pending' ? '入札受付中' : item.status}
                                        </span>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>

                    {/* Bidding Panel */}
                    <div className="lg:col-span-1">
                        <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 sticky top-6">
                            <h2 className="text-xl font-bold mb-6 text-gray-800 border-b pb-2">入札パネル</h2>
                            {selectedItem ? (
                                <form onSubmit={handleSubmit(onSubmitBid)} className="space-y-6">
                                    <div className="p-5 bg-gray-50 rounded-lg border border-gray-200">
                                        <p className="text-sm text-gray-500 mb-1">選択中の商品</p>
                                        <p className="font-bold text-2xl text-gray-900">{selectedItem.fish_type}</p>
                                        <p className="text-lg text-gray-700">{selectedItem.quantity} {selectedItem.unit}</p>
                                        <p className="text-sm text-gray-500 mt-2">ステータス: {selectedItem.status}</p>
                                        {selectedItem.highest_bid && (
                                            <p className="text-sm mt-2 text-orange-600 font-bold">
                                                現在の最高額: ¥{selectedItem.highest_bid.toLocaleString()}
                                                {selectedItem.highest_bidder_name && (
                                                    <span className="ml-2 text-gray-600">({selectedItem.highest_bidder_name} さん)</span>
                                                )}
                                            </p>
                                        )}
                                    </div>

                                    {selectedItem.status === 'Pending' ? (
                                        !auctionActive ? (
                                            <div className="space-y-4">
                                                <div className="text-center py-6 bg-yellow-50 rounded-lg border border-yellow-200">
                                                    <p className="text-yellow-800 font-bold mb-2">⏰ 入札受付時間外</p>
                                                    {auction.start_time && auction.end_time && (
                                                        <p className="text-sm text-yellow-700">
                                                            受付時間: {formatTime(auction.start_time)} ~ {formatTime(auction.end_time)}
                                                        </p>
                                                    )}
                                                </div>
                                            </div>
                                        ) : (
                                            <>
                                                <div>
                                                    <label className="block text-sm font-bold text-gray-700 mb-1">入札価格 (円)</label>
                                                    <div className="relative rounded-md shadow-sm">
                                                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                                            <span className="text-gray-500 sm:text-sm">¥</span>
                                                        </div>
                                                        <input
                                                            type="number"
                                                            {...register('price')}
                                                            className="block w-full rounded-md border-gray-300 pl-7 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border"
                                                            placeholder="0"
                                                        />
                                                    </div>
                                                    {errors.price && (
                                                        <p className="text-red-500 text-sm mt-1">{errors.price.message}</p>
                                                    )}
                                                </div>

                                                <button
                                                    type="submit"
                                                    disabled={isBidLoading}
                                                    className="w-full flex justify-center py-4 px-4 border border-transparent rounded-md shadow-md text-lg font-bold text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors transform hover:scale-[1.02] disabled:opacity-50"
                                                >
                                                    {isBidLoading ? '処理中...' : '落札する！'}
                                                </button>
                                            </>
                                        )
                                    ) : (
                                        <div className="text-center py-4 bg-gray-100 rounded text-gray-500">
                                            この商品は既に入札が終了しています
                                        </div>
                                    )}
                                </form>
                            ) : (
                                <div className="text-center py-12 text-gray-400">
                                    <p>左のリストから<br />商品を選択してください</p>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
