export interface Venue {
    id: number;
    name: string;
    location?: string;
    description?: string;
    createdAt: string;
}

export interface Auction {
    id: number;
    venueId: number;
    auctionDate: string;
    startTime?: string;
    endTime?: string;
    status: 'scheduled' | 'in_progress' | 'completed' | 'cancelled';
    createdAt: string;
    updatedAt: string;
}

export interface AuctionItem {
    id: number;
    auctionId: number;
    fishermanId: number;
    fishType: string;
    quantity: number;
    unit: string;
    status: string;
    highestBid?: number;
    highestBidderId?: number;
    highestBidderName?: string;
    createdAt: string;
}

export interface Bid {
    itemId: number;
    buyerId: number;
    price: number;
}
