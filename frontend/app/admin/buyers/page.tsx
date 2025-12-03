'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { buyerSchema, BuyerFormData } from '@/src/models/schemas/admin';
import { useBuyers, useBuyerMutations } from './_hooks/useBuyer';

export default function BuyersPage() {
    const [message, setMessage] = useState('');

    const { buyers, isLoading } = useBuyers();
    const { createBuyer, isCreating } = useBuyerMutations();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<BuyerFormData>({
        resolver: zodResolver(buyerSchema),
    });

    const onSubmit = async (data: BuyerFormData) => {
        try {
            await createBuyer({ name: data.name });
            setMessage('中買人を登録しました');
            reset();
        } catch (e) {
            console.error(e);
            setMessage('登録に失敗しました');
        }
    };

    return (
        <div className="max-w-5xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-8 text-gray-800 border-b pb-4">中買人管理</h1>

            {message && (
                <div className="bg-blue-50 border-l-4 border-blue-500 text-blue-700 p-4 mb-8 rounded shadow-sm" role="alert">
                    <p className="font-bold">通知</p>
                    <p>{message}</p>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                {/* Form Section */}
                <div className="md:col-span-1">
                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200 sticky top-6">
                        <h2 className="text-xl font-bold mb-6 text-green-900 flex items-center">
                            <span className="w-2 h-6 bg-green-500 mr-3 rounded-full"></span>
                            新規中買人登録
                        </h2>
                        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                            <div>
                                <label className="block text-sm font-bold text-gray-700 mb-1">氏名</label>
                                <input
                                    type="text"
                                    {...register('name')}
                                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500 sm:text-sm p-3 border bg-gray-50"
                                    placeholder="例: 鈴木 花子"
                                />
                                {errors.name && (
                                    <p className="text-red-500 text-sm mt-1">{errors.name.message}</p>
                                )}
                            </div>

                            <button
                                type="submit"
                                disabled={isCreating}
                                className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors disabled:opacity-50"
                            >
                                {isCreating ? '登録中...' : '登録する'}
                            </button>
                        </form>
                    </div>
                </div>

                {/* List Section */}
                <div className="md:col-span-2">
                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                        <div className="p-6 border-b border-gray-200">
                            <h2 className="text-xl font-bold text-gray-800">中買人一覧</h2>
                        </div>
                        {isLoading ? (
                            <div className="p-6 text-center text-gray-500">読み込み中...</div>
                        ) : buyers.length === 0 ? (
                            <div className="p-6 text-center text-gray-500">中買人が登録されていません</div>
                        ) : (
                            <ul className="divide-y divide-gray-200">
                                {buyers.map((buyer) => (
                                    <li key={buyer.id} className="p-6 hover:bg-gray-50 transition-colors">
                                        <div className="flex justify-between items-center">
                                            <div>
                                                <h3 className="text-lg font-bold text-green-900">{buyer.name}</h3>
                                                <p className="text-sm text-gray-500 mt-1">ID: {buyer.id}</p>
                                            </div>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
