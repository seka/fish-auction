package request

// CreateItem holds data for item creation.
type CreateItem struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

// UpdateItem holds data for updating an item.
type UpdateItem struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

// UpdateItemSortOrder holds data for updating an item's sort order.
type UpdateItemSortOrder struct {
	SortOrder int `json:"sort_order"`
}

// ReorderItems holds data for reordering multiple items.
type ReorderItems struct {
	IDs []int `json:"ids"`
}
