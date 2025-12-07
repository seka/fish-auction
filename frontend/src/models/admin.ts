export interface Fisherman {
    id?: number;
    name: string;
}

export interface Buyer {
    id?: number;
    name: string;
}

export interface RegisterItemParams {
    auctionId: number;
    fishermanId: number;
    fishType: string;
    quantity: number;
    unit: string;
}
