'use client';

import { useState } from 'react';
import { useRegisterFisherman } from '@/hooks/useFishermen';
import { useRegisterBuyer } from '@/hooks/useBuyers';
import { useRegisterItem } from '@/hooks/useItems';

export default function AdminPage() {
    const [fishermanName, setFishermanName] = useState('');
    const [buyerName, setBuyerName] = useState('');
    const [item, setItem] = useState({ fishermanId: '', fishType: '', quantity: '', unit: '' });
    const [message, setMessage] = useState('');

    const { registerFisherman, isLoading: isFishermanLoading } = useRegisterFisherman();
    const { registerBuyer, isLoading: isBuyerLoading } = useRegisterBuyer();
    const { registerItem, isLoading: isItemLoading } = useRegisterItem();

    const handleRegisterFisherman = async (e: React.FormEvent) => {
        e.preventDefault();
        const success = await registerFisherman({ name: fishermanName });
        if (success) {
            setMessage('Fisherman registered!');
            setFishermanName('');
        } else {
            setMessage('Failed to register fisherman');
        }
    };

    const handleRegisterBuyer = async (e: React.FormEvent) => {
        e.preventDefault();
        const success = await registerBuyer({ name: buyerName });
        if (success) {
            setMessage('Buyer registered!');
            setBuyerName('');
        } else {
            setMessage('Failed to register buyer');
        }
    };

    const handleRegisterItem = async (e: React.FormEvent) => {
        e.preventDefault();
        const success = await registerItem({
            fisherman_id: parseInt(item.fishermanId),
            fish_type: item.fishType,
            quantity: parseInt(item.quantity),
            unit: item.unit,
        });
        if (success) {
            setMessage('Item registered!');
            setItem({ fishermanId: '', fishType: '', quantity: '', unit: '' });
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
                    <form onSubmit={handleRegisterFisherman} className="space-y-4">
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">氏名</label>
                            <input
                                type="text"
                                value={fishermanName}
                                onChange={(e) => setFishermanName(e.target.value)}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: 山田 太郎"
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
                        >
                            登録する
                        </button>
                    </form>
                </div>

                {/* Register Buyer */}
                <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                    <h2 className="text-xl font-bold mb-6 text-green-900 flex items-center">
                        <span className="w-2 h-6 bg-green-500 mr-3 rounded-full"></span>
                        中買人登録
                    </h2>
                    <form onSubmit={handleRegisterBuyer} className="space-y-4">
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">屋号・氏名</label>
                            <input
                                type="text"
                                value={buyerName}
                                onChange={(e) => setBuyerName(e.target.value)}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: すしざんまい"
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors"
                        >
                            登録する
                        </button>
                    </form>
                </div>

                {/* Register Item */}
                <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-200 md:col-span-2">
                    <h2 className="text-xl font-bold mb-6 text-orange-900 flex items-center">
                        <span className="w-2 h-6 bg-orange-500 mr-3 rounded-full"></span>
                        出品登録 (セリ対象)
                    </h2>
                    <form onSubmit={handleRegisterItem} className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">漁師ID</label>
                            <input
                                type="number"
                                value={item.fishermanId}
                                onChange={(e) => setItem({ ...item, fishermanId: e.target.value })}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="IDを入力"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">魚種</label>
                            <input
                                type="text"
                                value={item.fishType}
                                onChange={(e) => setItem({ ...item, fishType: e.target.value })}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: マグロ"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">数量</label>
                            <input
                                type="number"
                                value={item.quantity}
                                onChange={(e) => setItem({ ...item, quantity: e.target.value })}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: 10"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-bold text-gray-700 mb-1">単位</label>
                            <input
                                type="text"
                                value={item.unit}
                                onChange={(e) => setItem({ ...item, unit: e.target.value })}
                                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-orange-500 focus:ring-orange-500 sm:text-sm p-3 border bg-gray-50"
                                placeholder="例: kg, 匹, 箱"
                                required
                            />
                        </div>
                        <div className="md:col-span-2 pt-4">
                            <button
                                type="submit"
                                className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500 transition-colors"
                            >
                                出品する
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}
