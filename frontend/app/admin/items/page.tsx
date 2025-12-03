'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { itemSchema, ItemFormData } from '@/src/models/schemas/admin';
import { useItemMutations } from './_hooks/useItem';
import { useFishermen } from '../fishermen/_hooks/useFisherman';
import { useAuctions } from '../auctions/_hooks/useAuction';

export default function ItemsPage() {
    const [message, setMessage] = useState('');

    const { fishermen } = useFishermen();
    const { auctions } = useAuctions({});
    const { createItem, isCreating } = useItemMutations();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<ItemFormData>({
        resolver: zodResolver(itemSchema),
    });

    const onSubmit = async (data: ItemFormData) => {
        try {
            await createItem({
                auction_id: parseInt(data.auctionId),
                fisherman_id: parseInt(data.fishermanId),
                fish_type: data.fishType,
                quantity: parseInt(data.quantity),
                unit: data.unit,
            });
            setMessage('出品を登録しました');
            reset();
        } catch (e) {
            console.error(e);
            setMessage('登録に失敗しました');
        }
    };

    return (
        <div className="max-w-6xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-8 text-gray-800 border-b pb-4">出品管理</h1>

            {message && (
                <div className="bg-blue-50 border-l-4 border-blue-500 text-blue-700 p-4 mb-8 rounded shadow-sm" role="alert">
                    <p className="font-bold">通知</p>
                    <p>{message}</p>
                </div>
            )}

            <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
                <h2 className="text-xl font-bold mb-6 text-orange-900 flex items-center">
                    <span className="w-2 h-6 bg-orange-500 mr-3 rounded-full"></span>
                    新規出品登録
                </h2>
                <form onSubmit={handleSubmit(onSubmit)} className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="md:col-span-2">
                        <label className="block text-sm font-bold text-gray-700 mb-1">セリ</label>
                        <select
                            {...register('auctionId')}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                        >
                            <option value="">セリを選択してください</option>
                            {auctions.map((auction) => (
                                <option key={auction.id} value={auction.id}>
                                    {auction.auction_date} {auction.start_time?.substring(0, 5)} - {auction.end_time?.substring(0, 5)} (ID: {auction.id})
                                </option>
                            ))}
                        </select>
                        {errors.auctionId && (
                            <p className="text-red-500 text-sm mt-1">{errors.auctionId.message}</p>
                        )}
                    </div>
                    <div>
                        <label className="block text-sm font-bold text-gray-700 mb-1">漁師</label>
                        <select
                            {...register('fishermanId')}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                        >
                            <option value="">漁師を選択してください</option>
                            {fishermen.map((fisherman) => (
                                <option key={fisherman.id} value={fisherman.id}>
                                    {fisherman.name}
                                </option>
                            ))}
                        </select>
                        {errors.fishermanId && (
                            <p className="text-red-500 text-sm mt-1">{errors.fishermanId.message}</p>
                        )}
                    </div>
                    <div>
                        <label className="block text-sm font-bold text-gray-700 mb-1">魚種</label>
                        <input
                            type="text"
                            {...register('fishType')}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                            placeholder="例: マグロ"
                        />
                        {errors.fishType && (
                            <p className="text-red-500 text-sm mt-1">{errors.fishType.message}</p>
                        )}
                    </div>
                    <div>
                        <label className="block text-sm font-bold text-gray-700 mb-1">数量</label>
                        <input
                            type="number"
                            {...register('quantity')}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                            placeholder="例: 10"
                        />
                        {errors.quantity && (
                            <p className="text-red-500 text-sm mt-1">{errors.quantity.message}</p>
                        )}
                    </div>
                    <div>
                        <label className="block text-sm font-bold text-gray-700 mb-1">単位</label>
                        <input
                            type="text"
                            {...register('unit')}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                            placeholder="例: kg, 匹, 箱"
                        />
                        {errors.unit && (
                            <p className="text-red-500 text-sm mt-1">{errors.unit.message}</p>
                        )}
                    </div>
                    <div className="md:col-span-2 pt-4">
                        <button
                            type="submit"
                            disabled={isCreating}
                            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500 transition-colors disabled:opacity-50"
                        >
                            {isCreating ? '出品中...' : '出品する'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
