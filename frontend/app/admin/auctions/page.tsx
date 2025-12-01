'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { auctionSchema, AuctionFormData } from '@/src/models/schemas/auction';
import { useAuctions, useAuctionMutations } from './_hooks/useAuction';
import { useVenues } from '../venues/_hooks/useVenue';
import { Auction } from '@/src/models/auction';

export default function AuctionsPage() {
    const [message, setMessage] = useState('');
    const [editingAuction, setEditingAuction] = useState<Auction | null>(null);
    const [filterVenueId, setFilterVenueId] = useState<number | undefined>(undefined);

    const { venues } = useVenues();
    const { auctions, isLoading } = useAuctions({ venue_id: filterVenueId });
    const { createAuction, updateAuction, updateStatus, deleteAuction, isCreating, isUpdating, isUpdatingStatus, isDeleting } = useAuctionMutations();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<AuctionFormData>({
        resolver: zodResolver(auctionSchema),
    });

    const onSubmit = async (data: AuctionFormData) => {
        try {
            // Convert strings to appropriate types if needed (though schema handles validation)
            // The API expects venue_id as number
            const payload = {
                ...data,
                venue_id: Number(data.venue_id),
            };

            if (editingAuction) {
                await updateAuction({ id: editingAuction.id, data: payload });
                setMessage('セリ情報を更新しました');
                setEditingAuction(null);
            } else {
                await createAuction(payload);
                setMessage('セリを作成しました');
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage('エラーが発生しました');
        }
    };

    const onEdit = (auction: Auction) => {
        setEditingAuction(auction);
        setValue('venue_id', auction.venue_id);
        setValue('auction_date', auction.auction_date);
        setValue('start_time', auction.start_time || '');
        setValue('end_time', auction.end_time || '');
        setValue('status', auction.status);
    };

    const onCancelEdit = () => {
        setEditingAuction(null);
        reset();
    };

    const onDelete = async (id: number) => {
        if (confirm('本当に削除しますか？')) {
            try {
                await deleteAuction(id);
                setMessage('セリを削除しました');
            } catch (e) {
                console.error(e);
                setMessage('削除に失敗しました');
            }
        }
    };

    const onStatusChange = async (id: number, status: string) => {
        try {
            await updateStatus({ id, status });
            setMessage(`ステータスを ${status} に更新しました`);
        } catch (e) {
            console.error(e);
            setMessage('ステータス更新に失敗しました');
        }
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'scheduled':
                return <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded">予定</span>;
            case 'in_progress':
                return <span className="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded animate-pulse">開催中</span>;
            case 'completed':
                return <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded">終了</span>;
            case 'cancelled':
                return <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded">中止</span>;
            default:
                return <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded">{status}</span>;
        }
    };

    return (
        <div className="max-w-6xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-8 text-gray-800 border-b pb-4">セリ管理</h1>

            {message && (
                <div className="bg-blue-50 border-l-4 border-blue-500 text-blue-700 p-4 mb-8 rounded shadow-sm" role="alert">
                    <p className="font-bold">通知</p>
                    <p>{message}</p>
                </div>
            )}

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Form Section */}
                <div className="lg:col-span-1">
                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200 sticky top-6">
                        <h2 className="text-xl font-bold mb-6 text-indigo-900 flex items-center">
                            <span className="w-2 h-6 bg-indigo-500 mr-3 rounded-full"></span>
                            {editingAuction ? 'セリ編集' : '新規セリ登録'}
                        </h2>
                        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                            <div>
                                <label className="block text-sm font-bold text-gray-700 mb-1">会場</label>
                                <select
                                    {...register('venue_id', { valueAsNumber: true })}
                                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                >
                                    <option value="">会場を選択してください</option>
                                    {venues.map((venue) => (
                                        <option key={venue.id} value={venue.id}>
                                            {venue.name}
                                        </option>
                                    ))}
                                </select>
                                {errors.venue_id && (
                                    <p className="text-red-500 text-sm mt-1">{errors.venue_id.message}</p>
                                )}
                            </div>
                            <div>
                                <label className="block text-sm font-bold text-gray-700 mb-1">開催日</label>
                                <input
                                    type="date"
                                    {...register('auction_date')}
                                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                />
                                {errors.auction_date && (
                                    <p className="text-red-500 text-sm mt-1">{errors.auction_date.message}</p>
                                )}
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-bold text-gray-700 mb-1">開始時間</label>
                                    <input
                                        type="time"
                                        {...register('start_time')}
                                        className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-bold text-gray-700 mb-1">終了時間</label>
                                    <input
                                        type="time"
                                        {...register('end_time')}
                                        className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-3 border bg-gray-50"
                                    />
                                </div>
                            </div>

                            <div className="flex gap-2 pt-4">
                                <button
                                    type="submit"
                                    disabled={isCreating || isUpdating}
                                    className="flex-1 flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors disabled:opacity-50"
                                >
                                    {editingAuction ? (isUpdating ? '更新中...' : '更新する') : (isCreating ? '登録中...' : '登録する')}
                                </button>
                                {editingAuction && (
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
                <div className="lg:col-span-2">
                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                        <div className="p-6 border-b border-gray-200 flex justify-between items-center flex-wrap gap-4">
                            <h2 className="text-xl font-bold text-gray-800">セリ一覧</h2>
                            <div className="flex items-center gap-2">
                                <label className="text-sm text-gray-600">会場絞り込み:</label>
                                <select
                                    value={filterVenueId || ''}
                                    onChange={(e) => setFilterVenueId(e.target.value ? Number(e.target.value) : undefined)}
                                    className="rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                >
                                    <option value="">すべて</option>
                                    {venues.map((venue) => (
                                        <option key={venue.id} value={venue.id}>
                                            {venue.name}
                                        </option>
                                    ))}
                                </select>
                            </div>
                        </div>
                        {isLoading ? (
                            <div className="p-6 text-center text-gray-500">読み込み中...</div>
                        ) : auctions.length === 0 ? (
                            <div className="p-6 text-center text-gray-500">セリが登録されていません</div>
                        ) : (
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">開催日 / 時間</th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">会場</th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ステータス</th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {auctions.map((auction) => {
                                            const venue = venues.find(v => v.id === auction.venue_id);
                                            return (
                                                <tr key={auction.id} className="hover:bg-gray-50">
                                                    <td className="px-6 py-4 whitespace-nowrap">
                                                        <div className="text-sm font-medium text-gray-900">{auction.auction_date}</div>
                                                        <div className="text-sm text-gray-500">
                                                            {auction.start_time ? auction.start_time.substring(0, 5) : '--:--'} - {auction.end_time ? auction.end_time.substring(0, 5) : '--:--'}
                                                        </div>
                                                    </td>
                                                    <td className="px-6 py-4 whitespace-nowrap">
                                                        <div className="text-sm text-gray-900">{venue?.name || `ID: ${auction.venue_id}`}</div>
                                                    </td>
                                                    <td className="px-6 py-4 whitespace-nowrap">
                                                        {getStatusBadge(auction.status)}
                                                    </td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                        <div className="flex justify-end gap-2">
                                                            {auction.status === 'scheduled' && (
                                                                <button
                                                                    onClick={() => onStatusChange(auction.id, 'in_progress')}
                                                                    disabled={isUpdatingStatus}
                                                                    className="text-green-600 hover:text-green-900 bg-green-50 px-2 py-1 rounded"
                                                                >
                                                                    開始
                                                                </button>
                                                            )}
                                                            {auction.status === 'in_progress' && (
                                                                <button
                                                                    onClick={() => onStatusChange(auction.id, 'completed')}
                                                                    disabled={isUpdatingStatus}
                                                                    className="text-blue-600 hover:text-blue-900 bg-blue-50 px-2 py-1 rounded"
                                                                >
                                                                    終了
                                                                </button>
                                                            )}
                                                            <button
                                                                onClick={() => onEdit(auction)}
                                                                className="text-indigo-600 hover:text-indigo-900 bg-indigo-50 px-2 py-1 rounded"
                                                            >
                                                                編集
                                                            </button>
                                                            <button
                                                                onClick={() => onDelete(auction.id)}
                                                                disabled={isDeleting}
                                                                className="text-red-600 hover:text-red-900 bg-red-50 px-2 py-1 rounded disabled:opacity-50"
                                                            >
                                                                削除
                                                            </button>
                                                        </div>
                                                    </td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
