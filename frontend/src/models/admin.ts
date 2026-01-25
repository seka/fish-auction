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

export interface UpdateItemParams extends RegisterItemParams {
    id: number;
    status: string;
}

export interface UpdateItemSortOrderParams {
    id: number;
    auctionId: number;
    sortOrder: number;
    newIndex: number;
}

export interface ReorderItemsParams {
    auctionId: number;
    ids: number[];
}
