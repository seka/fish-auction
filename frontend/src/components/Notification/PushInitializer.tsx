'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { usePushSubscription } from '../../hooks/usePushSubscription';
import { useAuth } from '../../hooks/useAuth';

export const PushInitializer = () => {
    const { isSupported, subscription, subscribeToPush } = usePushSubscription();
    const { isLoggedIn } = useAuth();
    const pathname = usePathname();

    useEffect(() => {
        // Skip for admin pages
        if (pathname?.startsWith('/admin')) {
            return;
        }

        if (isSupported && isLoggedIn && !subscription) {
            // Check if permission is already granted
            if (Notification.permission === 'granted') {
                subscribeToPush();
            }
        }
    }, [isSupported, isLoggedIn, subscription, subscribeToPush, pathname]);

    // For now, render nothing. 
    // Ideally we could show a toaster or button if permission is default.
    return null;
};
