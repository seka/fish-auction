export interface AuctionItem {
    id: number;
    fisherman_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
    status: string;
    created_at: string;
}

export interface Bid {
    item_id: number;
    buyer_id: number;
    price: number;
}
