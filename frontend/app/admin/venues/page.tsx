'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { venueSchema, VenueFormData } from '@/src/models/schemas/auction';
import { useVenues, useVenueMutations } from './_hooks/useVenue';
import { Venue } from '@/src/models/venue';

export default function VenuesPage() {
    const [message, setMessage] = useState('');
    const [editingVenue, setEditingVenue] = useState<Venue | null>(null);

    const { venues, isLoading } = useVenues();
    const { createVenue, updateVenue, deleteVenue, isCreating, isUpdating, isDeleting } = useVenueMutations();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<VenueFormData>({
        resolver: zodResolver(venueSchema),
    });

    const onSubmit = async (data: VenueFormData) => {
        try {
            if (editingVenue) {
                await updateVenue({ id: editingVenue.id, data });
                setMessage('会場を更新しました');
                setEditingVenue(null);
            } else {
                await createVenue(data);
                setMessage('会場を作成しました');
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage('エラーが発生しました');
        }
    };

    const onEdit = (venue: Venue) => {
        setEditingVenue(venue);
        setValue('name', venue.name);
        setValue('location', venue.location || '');
        setValue('description', venue.description || '');
    };

    const onCancelEdit = () => {
        setEditingVenue(null);
        reset();
    };

    const onDelete = async (id: number) => {
        if (confirm('本当に削除しますか？')) {
            try {
                await deleteVenue(id);
                setMessage('会場を削除しました');
            } catch (e) {
                console.error(e);
                setMessage('削除に失敗しました');
            }
        }
    };

    return (
        <div className="max-w-5xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-8 text-gray-800 border-b pb-4">会場管理</h1>

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
                        <h2 className="text-xl font-bold mb-6 text-indigo-900 flex items-center">
                            <span className="w-2 h-6 bg-indigo-500 mr-3 rounded-full"></span>
                            {editingVenue ? '会場編集' : '新規会場登録'}
                        </h2>
                        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                            <div>
                                <label className="block text-sm font-bold text-gray-700 mb-1">会場名</label>
                                <input
                                    type="text"
                                    {...register('name')}
                                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                    placeholder="例: 豊洲市場"
                                />
                                {errors.name && (
                                    <p className="text-red-500 text-sm mt-1">{errors.name.message}</p>
                                )}
                            </div>
                            <div>
                                <label className="block text-sm font-bold text-gray-700 mb-1">所在地</label>
                                <input
                                    type="text"
                                    {...register('location')}
                                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                    placeholder="例: 東京都江東区..."
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-bold text-gray-700 mb-1">説明</label>
                                <textarea
                                    {...register('description')}
                                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                    rows={3}
                                    placeholder="会場の説明..."
                                />
                            </div>

                            <div className="flex gap-2">
                                <button
                                    type="submit"
                                    disabled={isCreating || isUpdating}
                                    className="flex-1 flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors disabled:opacity-50"
                                >
                                    {editingVenue ? (isUpdating ? '更新中...' : '更新する') : (isCreating ? '登録中...' : '登録する')}
                                </button>
                                {editingVenue && (
                                    <button
                                        type="button"
                                        onClick={onCancelEdit}
                                        className="px-4 py-3 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                                    >
                                        キャンセル
                                    </button>
                                )}
                            </div>
                        </form>
                    </div>
                </div>

                {/* List Section */}
                <div className="md:col-span-2">
                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                        <div className="p-6 border-b border-gray-200">
                            <h2 className="text-xl font-bold text-gray-800">会場一覧</h2>
                        </div>
                        {isLoading ? (
                            <div className="p-6 text-center text-gray-500">読み込み中...</div>
                        ) : venues.length === 0 ? (
                            <div className="p-6 text-center text-gray-500">会場が登録されていません</div>
                        ) : (
                            <ul className="divide-y divide-gray-200">
                                {venues.map((venue) => (
                                    <li key={venue.id} className="p-6 hover:bg-gray-50 transition-colors">
                                        <div className="flex justify-between items-start">
                                            <div>
                                                <h3 className="text-lg font-bold text-indigo-900">{venue.name}</h3>
                                                {venue.location && (
                                                    <p className="text-sm text-gray-600 mt-1 flex items-center">
                                                        <span className="mr-2">📍</span>
                                                        {venue.location}
                                                    </p>
                                                )}
                                                {venue.description && (
                                                    <p className="text-sm text-gray-500 mt-2">{venue.description}</p>
                                                )}
                                            </div>
                                            <div className="flex space-x-2">
                                                <button
                                                    onClick={() => onEdit(venue)}
                                                    className="text-indigo-600 hover:text-indigo-900 text-sm font-medium px-3 py-1 rounded hover:bg-indigo-50"
                                                >
                                                    編集
                                                </button>
                                                <button
                                                    onClick={() => onDelete(venue.id)}
                                                    disabled={isDeleting}
                                                    className="text-red-600 hover:text-red-900 text-sm font-medium px-3 py-1 rounded hover:bg-red-50 disabled:opacity-50"
                                                >
                                                    削除
                                                </button>
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
