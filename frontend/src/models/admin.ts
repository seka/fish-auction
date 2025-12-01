export interface Fisherman {
    id?: number;
    name: string;
}

export interface Buyer {
    id?: number;
    name: string;
}

export interface RegisterItemParams {
    auction_id: number;
    fisherman_id: number;
    fish_type: string;
    quantity: number;
    unit: string;
}
