'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useItems, useSubmitBid, useBuyers } from './_hooks/useAuction';
import { AuctionItem } from '@/src/models';
import { bidSchema, BidFormData } from '@/src/models/schemas/auction';

export default function AuctionPage() {
    const [selectedItem, setSelectedItem] = useState<AuctionItem | null>(null);
    const [message, setMessage] = useState('');

    const { items, refetch } = useItems({
        status: 'Pending',
        pollingInterval: 5000
    });
    const { submitBid, isLoading: isBidLoading } = useSubmitBid();
    const { buyers } = useBuyers();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<BidFormData>({
        resolver: zodResolver(bidSchema),
    });

    const onSubmitBid = async (data: BidFormData) => {
        if (!selectedItem) return;

        const success = await submitBid({
            item_id: selectedItem.id,
            buyer_id: 0, // Backend handles this from context
            price: parseInt(data.price),
        });

        if (success) {
            setMessage(`Successfully bought ${selectedItem.fish_type}!`);
            setSelectedItem(null);
            reset();
            refetch();
        } else {
            setMessage('Failed to submit bid');
        }
    };

    return (
        <div className="p-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-end mb-8">
                <h1 className="text-3xl font-bold text-gray-900">本日の出品一覧</h1>
                <p className="text-gray-500 text-sm">5秒ごとに自動更新されます</p>
            </div>

            {message && (
                <div className="bg-green-50 border-l-4 border-green-500 text-green-700 p-4 mb-6 rounded shadow-sm animate-pulse" role="alert">
                    <p className="font-bold">落札成功！</p>
                    <p>{message}</p>
                </div>
            )}

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Item List */}
                <div className="lg:col-span-2 space-y-4">
                    <h2 className="text-xl font-bold mb-4 text-gray-800 border-b pb-2">出品リスト</h2>
                    {items.length === 0 ? (
                        <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
                            <p className="text-gray-500">現在、セリに出されている商品はありません。</p>
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
                                        </div>
                                    </div>
                                    <span className="px-4 py-2 bg-green-100 text-green-800 rounded-full text-sm font-bold shadow-sm">
                                        {item.status === 'Pending' ? '入札受付中' : item.status}
                                    </span>
                                </div>
                            </div>
                        ))
                    )}
                </div>

                {/* Bidding Panel */}
                <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200 h-fit sticky top-24">
                    <h2 className="text-xl font-bold mb-6 text-gray-800 border-b pb-2">入札パネル</h2>
                    {selectedItem ? (
                        <form onSubmit={handleSubmit(onSubmitBid)} className="space-y-6">
                            <div className="p-5 bg-gray-50 rounded-lg border border-gray-200">
                                <p className="text-sm text-gray-500 mb-1">選択中の商品</p>
                                <p className="font-bold text-2xl text-gray-900">{selectedItem.fish_type}</p>
                                <p className="text-lg text-gray-700">{selectedItem.quantity} {selectedItem.unit}</p>
                            </div>

                            {/* Buyer selection removed */}

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
                        </form>
                    ) : (
                        <div className="text-center py-12 text-gray-400">
                            <p>左のリストから<br />商品を選択してください</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
