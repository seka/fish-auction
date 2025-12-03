import { Auction } from '@/src/models/auction';

/**
 * Check if an auction is currently active (within bidding hours)
 */
export const isAuctionActive = (auction: Auction): boolean => {
    if (!auction.start_time || !auction.end_time) {
        // If no time is set, consider it always active
        return true;
    }

    const now = new Date();
    const auctionDate = new Date(auction.auction_date);

    // Parse start and end times
    const [startHour, startMin] = auction.start_time.split(':').map(Number);
    const [endHour, endMin] = auction.end_time.split(':').map(Number);

    // Create start and end datetime objects
    const startDateTime = new Date(auctionDate);
    startDateTime.setHours(startHour, startMin, 0, 0);

    const endDateTime = new Date(auctionDate);
    endDateTime.setHours(endHour, endMin, 0, 0);

    return now >= startDateTime && now <= endDateTime;
};

/**
 * Format time for display (HH:MM)
 */
export const formatTime = (time?: string): string => {
    if (!time) return '';
    return time.substring(0, 5); // Extract HH:MM from HH:MM:SS
};
