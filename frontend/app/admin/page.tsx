'use client';

import { useState } from 'react';

export default function AdminPage() {
    const [fishermanName, setFishermanName] = useState('');
    const [buyerName, setBuyerName] = useState('');
    const [item, setItem] = useState({ fishermanId: '', fishType: '', quantity: '', unit: '' });
    const [message, setMessage] = useState('');

    const registerFisherman = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await fetch('/api/fishermen', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: fishermanName }),
            });
            if (res.ok) {
                setMessage('Fisherman registered!');
                setFishermanName('');
            } else {
                setMessage('Failed to register fisherman');
            }
        } catch (error) {
            setMessage('Error registering fisherman');
        }
    };

    const registerBuyer = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await fetch('/api/buyers', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: buyerName }),
            });
            if (res.ok) {
                setMessage('Buyer registered!');
                setBuyerName('');
            } else {
                setMessage('Failed to register buyer');
            }
        } catch (error) {
            setMessage('Error registering buyer');
        }
    };

    const registerItem = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await fetch('/api/items', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    fisherman_id: parseInt(item.fishermanId),
                    fish_type: item.fishType,
                    quantity: parseInt(item.quantity),
                    unit: item.unit,
                }),
            });
            if (res.ok) {
                setMessage('Item registered!');
                setItem({ fishermanId: '', fishType: '', quantity: '', unit: '' });
            } else {
                setMessage('Failed to register item');
            }
        } catch (error) {
            setMessage('Error registering item');
        }
    };

    return (
        <div className="p-8 max-w-4xl mx-auto">
            <h1 className="text-3xl font-bold mb-8">Admin Dashboard</h1>

            {message && (
                <div className="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 mb-6" role="alert">
                    <p>{message}</p>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                {/* Register Fisherman */}
                <div className="bg-white p-6 rounded-lg shadow-md">
                    <h2 className="text-xl font-semibold mb-4">Register Fisherman</h2>
                    <form onSubmit={registerFisherman} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Name</label>
                            <input
                                type="text"
                                value={fishermanName}
                                onChange={(e) => setFishermanName(e.target.value)}
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                        >
                            Register
                        </button>
                    </form>
                </div>

                {/* Register Buyer */}
                <div className="bg-white p-6 rounded-lg shadow-md">
                    <h2 className="text-xl font-semibold mb-4">Register Buyer</h2>
                    <form onSubmit={registerBuyer} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Name</label>
                            <input
                                type="text"
                                value={buyerName}
                                onChange={(e) => setBuyerName(e.target.value)}
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                        >
                            Register
                        </button>
                    </form>
                </div>

                {/* Register Item */}
                <div className="bg-white p-6 rounded-lg shadow-md md:col-span-2">
                    <h2 className="text-xl font-semibold mb-4">Register Auction Item</h2>
                    <form onSubmit={registerItem} className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Fisherman ID</label>
                            <input
                                type="number"
                                value={item.fishermanId}
                                onChange={(e) => setItem({ ...item, fishermanId: e.target.value })}
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Fish Type</label>
                            <input
                                type="text"
                                value={item.fishType}
                                onChange={(e) => setItem({ ...item, fishType: e.target.value })}
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Quantity</label>
                            <input
                                type="number"
                                value={item.quantity}
                                onChange={(e) => setItem({ ...item, quantity: e.target.value })}
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Unit</label>
                            <input
                                type="text"
                                value={item.unit}
                                onChange={(e) => setItem({ ...item, unit: e.target.value })}
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                placeholder="kg, box, etc."
                                required
                            />
                        </div>
                        <div className="md:col-span-2">
                            <button
                                type="submit"
                                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500"
                            >
                                Register Item
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}
