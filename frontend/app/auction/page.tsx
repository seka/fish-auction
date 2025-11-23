'use client';

import { useState, useEffect } from 'react';

interface AuctionItem {
    id: number;
    fisherman_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
    status: string;
    created_at: string;
}

export default function AuctionPage() {
    const [items, setItems] = useState<AuctionItem[]>([]);
    const [selectedItem, setSelectedItem] = useState<AuctionItem | null>(null);
    const [buyerId, setBuyerId] = useState('');
    const [price, setPrice] = useState('');
    const [message, setMessage] = useState('');

    const fetchItems = async () => {
        try {
            const res = await fetch('/api/items?status=Pending');
            if (res.ok) {
                const data = await res.json();
                setItems(data || []);
            }
        } catch (error) {
            console.error('Failed to fetch items', error);
        }
    };

    useEffect(() => {
        fetchItems();
        const interval = setInterval(fetchItems, 5000); // Poll every 5s
        return () => clearInterval(interval);
    }, []);

    const handleBid = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedItem) return;

        try {
            const res = await fetch('/api/bid', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    item_id: selectedItem.id,
                    buyer_id: parseInt(buyerId),
                    price: parseInt(price),
                }),
            });

            if (res.ok) {
                setMessage(`Successfully bought ${selectedItem.fish_type}!`);
                setSelectedItem(null);
                setBuyerId('');
                setPrice('');
                fetchItems();
            } else {
                setMessage('Failed to submit bid');
            }
        } catch (error) {
            setMessage('Error submitting bid');
        }
    };

    return (
        <div className="p-8 max-w-6xl mx-auto">
            <h1 className="text-3xl font-bold mb-8">Auction Floor</h1>

            {message && (
                <div className="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-6" role="alert">
                    <p>{message}</p>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Item List */}
                <div className="md:col-span-2 space-y-4">
                    <h2 className="text-xl font-semibold mb-4">Available Items</h2>
                    {items.length === 0 ? (
                        <p className="text-gray-500">No items currently for auction.</p>
                    ) : (
                        items.map((item) => (
                            <div
                                key={item.id}
                                className={`p-4 border rounded-lg cursor-pointer hover:shadow-md transition-shadow ${selectedItem?.id === item.id ? 'border-blue-500 bg-blue-50' : 'bg-white'
                                    }`}
                                onClick={() => setSelectedItem(item)}
                            >
                                <div className="flex justify-between items-center">
                                    <div>
                                        <h3 className="text-lg font-bold">{item.fish_type}</h3>
                                        <p className="text-gray-600">
                                            {item.quantity} {item.unit} (Fisherman ID: {item.fisherman_id})
                                        </p>
                                    </div>
                                    <span className="px-3 py-1 bg-green-100 text-green-800 rounded-full text-sm">
                                        {item.status}
                                    </span>
                                </div>
                            </div>
                        ))
                    )}
                </div>

                {/* Bidding Panel */}
                <div className="bg-white p-6 rounded-lg shadow-md h-fit sticky top-6">
                    <h2 className="text-xl font-semibold mb-4">Place Bid</h2>
                    {selectedItem ? (
                        <form onSubmit={handleBid} className="space-y-4">
                            <div className="p-4 bg-gray-50 rounded mb-4">
                                <p className="font-bold">{selectedItem.fish_type}</p>
                                <p>{selectedItem.quantity} {selectedItem.unit}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700">Buyer ID</label>
                                <input
                                    type="number"
                                    value={buyerId}
                                    onChange={(e) => setBuyerId(e.target.value)}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                    required
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700">Bid Price (JPY)</label>
                                <input
                                    type="number"
                                    value={price}
                                    onChange={(e) => setPrice(e.target.value)}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                    required
                                />
                            </div>

                            <button
                                type="submit"
                                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                            >
                                Submit Bid
                            </button>
                        </form>
                    ) : (
                        <p className="text-gray-500 text-center py-8">Select an item to bid</p>
                    )}
                </div>
            </div>
        </div>
    );
}
