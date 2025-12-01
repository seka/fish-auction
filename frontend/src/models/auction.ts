export interface Venue {
    id: number;
    name: string;
    location?: string;
    description?: string;
    created_at: string;
}

export interface Auction {
    id: number;
    venue_id: number;
    auction_date: string;
    start_time?: string;
    end_time?: string;
    status: 'scheduled' | 'in_progress' | 'completed' | 'cancelled';
    created_at: string;
    updated_at: string;
}

export interface AuctionItem {
    id: number;
    auction_id: number;
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
