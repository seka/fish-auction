export const translateAuctionStatus = (status: string): string => {
    switch (status) {
        case 'scheduled':
            return '開催予定';
        case 'in_progress':
            return '開催中';
        case 'completed':
            return '終了';
        case 'cancelled':
            return '中止';
        default:
            return status;
    }
};

export const translateItemStatus = (status: string): string => {
    switch (status) {
        case 'Pending':
            return '入札受付中';
        case 'Sold':
            return '落札済';
        case 'Unsold':
            return '不落';
        case 'Bidding':
            return '入札中';
        default:
            return status;
    }
};
