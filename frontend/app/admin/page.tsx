'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useRegisterFisherman, useRegisterBuyer, useRegisterItem, useFishermen } from './_hooks/useAdmin';
import { useAuctions } from './auctions/_hooks/useAuction';
import { fishermanSchema, buyerSchema, itemSchema, FishermanFormData, BuyerFormData, ItemFormData } from '@/src/models/schemas/admin';

export default function AdminPage() {
    const [message, setMessage] = useState('');

    const { registerFisherman, isLoading: isFishermanLoading } = useRegisterFisherman();
    const { registerBuyer, isLoading: isBuyerLoading } = useRegisterBuyer();
    const { registerItem, isLoading: isItemLoading } = useRegisterItem();
    const { fishermen } = useFishermen();
    const { auctions } = useAuctions({ status: 'scheduled' }); // Only show scheduled auctions

    const { register: registerFishermanForm, handleSubmit: handleSubmitFisherman, reset: resetFisherman, formState: { errors: fishermanErrors } } = useForm<FishermanFormData>({
        resolver: zodResolver(fishermanSchema),
    });
    const { register: registerBuyerForm, handleSubmit: handleSubmitBuyer, reset: resetBuyer, formState: { errors: buyerErrors } } = useForm<BuyerFormData>({
        resolver: zodResolver(buyerSchema),
    });
    const { register: registerItemForm, handleSubmit: handleSubmitItem, reset: resetItem, formState: { errors: itemErrors } } = useForm<ItemFormData>({
        resolver: zodResolver(itemSchema),
    });

    const onRegisterFisherman = async (data: FishermanFormData) => {
        const success = await registerFisherman({ name: data.name });
        if (success) {
            setMessage('Fisherman registered!');
            resetFisherman();
        } else {
            setMessage('Failed to register fisherman');
        }
    };

    const onRegisterBuyer = async (data: BuyerFormData) => {
        const success = await registerBuyer({ name: data.name });
        if (success) {
            setMessage('Buyer registered!');
            resetBuyer();
        } else {
            setMessage('Failed to register buyer');
        }
    };

    const onRegisterItem = async (data: ItemFormData) => {
        const success = await registerItem({
            auction_id: parseInt(data.auctionId),
            fisherman_id: parseInt(data.fishermanId),
            fish_type: data.fishType,
            quantity: parseInt(data.quantity),
            unit: data.unit,
        });
        if (success) {
            setMessage('Item registered!');
            resetItem();
        } else {
            setMessage('Failed to register item');
        }
    };

    return (
        <div className="max-w-5xl mx-auto">
            <h1 className="text-3xl font-bold mb-8 text-gray-800 border-b pb-4">管理ダッシュボード</h1>

            {message && (
                <div className="bg-blue-50 border-l-4 border-blue-500 text-blue-700 p-4 mb-8 rounded shadow-sm" role="alert">
                    <p className="font-bold">通知</p>
                    <p>{message}</p>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                {/* Register Fisherman */}
                <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                    <h2 className="text-xl font-bold mb-6 text-indigo-900 flex items-center">
                        <span className="w-2 h-6 bg-indigo-500 mr-3 rounded-full"></span>
                        漁師登録
                    </h2>
                    <form onSubmit={handleSubmitFisherman(onRegisterFisherman)} className="space-y-4">
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">氏名</label>
                            <input
                                type="text"
                                {...registerFishermanForm('name')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: 山田 太郎"
                            />
                            {fishermanErrors.name && (
                                <p className="text-red-500 text-sm mt-1">{fishermanErrors.name.message}</p>
                            )}
                        </div>
                        <button
                            type="submit"
                            disabled={isFishermanLoading}
                            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors disabled:opacity-50"
                        >
                            {isFishermanLoading ? '登録中...' : '登録する'}
                        </button>
                    </form>
                </div>

                {/* Register Buyer */}
                <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                    <h2 className="text-xl font-bold mb-6 text-green-900 flex items-center">
                        <span className="w-2 h-6 bg-green-500 mr-3 rounded-full"></span>
                        中買人登録
                    </h2>
                    <form onSubmit={handleSubmitBuyer(onRegisterBuyer)} className="space-y-4">
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">屋号・氏名</label>
                            <input
                                type="text"
                                {...registerBuyerForm('name')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: 魚市場 花子"
                            />
                            {buyerErrors.name && (
                                <p className="text-red-500 text-sm mt-1">{buyerErrors.name.message}</p>
                            )}
                        </div>
                        <button
                            type="submit"
                            disabled={isBuyerLoading}
                            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors disabled:opacity-50"
                        >
                            {isBuyerLoading ? '登録中...' : '登録する'}
                        </button>
                    </form>
                </div>

                {/* Register Item */}
                <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-200 md:col-span-2">
                    <h2 className="text-xl font-bold mb-6 text-orange-900 flex items-center">
                        <span className="w-2 h-6 bg-orange-500 mr-3 rounded-full"></span>
                        出品登録 (セリ対象)
                    </h2>
                    <form onSubmit={handleSubmitItem(onRegisterItem)} className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="md:col-span-2">
                            <label className="block text-sm font-bold text-gray-700 mb-1">セリ (開催予定)</label>
                            <select
                                {...registerItemForm('auctionId')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                            >
                                <option value="">セリを選択してください</option>
                                {auctions.map((auction) => (
                                    <option key={auction.id} value={auction.id}>
                                        {auction.auction_date} {auction.start_time?.substring(0, 5)} - {auction.end_time?.substring(0, 5)} (ID: {auction.id})
                                    </option>
                                ))}
                            </select>
                            {itemErrors.auctionId && (
                                <p className="text-red-500 text-sm mt-1">{itemErrors.auctionId.message}</p>
                            )}
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">漁師</label>
                            <select
                                {...registerItemForm('fishermanId')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                            >
                                <option value="">漁師を選択してください</option>
                                {fishermen.map((fisherman) => (
                                    <option key={fisherman.id} value={fisherman.id}>
                                        {fisherman.name}
                                    </option>
                                ))}
                            </select>
                            {itemErrors.fishermanId && (
                                <p className="text-red-500 text-sm mt-1">{itemErrors.fishermanId.message}</p>
                            )}
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">魚種</label>
                            <input
                                type="text"
                                {...registerItemForm('fishType')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: マグロ"
                            />
                            {itemErrors.fishType && (
                                <p className="text-red-500 text-sm mt-1">{itemErrors.fishType.message}</p>
                            )}
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">数量</label>
                            <input
                                type="number"
                                {...registerItemForm('quantity')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: 10"
                            />
                            {itemErrors.quantity && (
                                <p className="text-red-500 text-sm mt-1">{itemErrors.quantity.message}</p>
                            )}
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">単位</label>
                            <input
                                type="text"
                                {...registerItemForm('unit')}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: kg, 匹, 箱"
                            />
                            {itemErrors.unit && (
                                <p className="text-red-500 text-sm mt-1">{itemErrors.unit.message}</p>
                            )}
                        </div>
                        <div className="md:col-span-2 pt-4">
                            <button
                                type="submit"
                                disabled={isItemLoading}
                                className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500 transition-colors disabled:opacity-50"
                            >
                                {isItemLoading ? '出品中...' : '出品する'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}
