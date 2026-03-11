'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { usePushSubscription } from '../../../hooks/usePushSubscription';
import { useAuth } from '../../../hooks/useAuth';
import { useToast } from '../../../hooks/useToast';

export const useSetup = () => {
  const { isSupported, subscription, subscribeToPush } = usePushSubscription();
  const { isLoggedIn } = useAuth();
  const { showToast } = useToast();
  const pathname = usePathname();

  useEffect(() => {
    // Skip for admin pages
    if (pathname?.startsWith('/admin')) {
      return;
    }

    console.log('PushInitializer checking status:', {
      isSupported,
      isLoggedIn,
      hasSubscription: !!subscription,
      permission: typeof window !== 'undefined' ? Notification.permission : 'unknown',
      hasVapidKey: !!process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY,
    });

    if (isSupported && isLoggedIn && !subscription) {
      // Attempt to subscribe if granted, or if default (to trigger prompt)
      if (Notification.permission === 'granted' || Notification.permission === 'default') {
        console.log('Attempting subscription...');
        subscribeToPush();
      } else {
        console.log('Permission is:', Notification.permission);
      }
    }
  }, [isSupported, isLoggedIn, subscription, subscribeToPush, pathname]);

  useEffect(() => {
    if (!('serviceWorker' in navigator)) return;

    const handleMessage = (event: MessageEvent) => {
      console.log('PushInitializer: Received message', event.data);
      if (event.data && event.data.type === 'PUSH_NOTIFICATION') {
        showToast(
          event.data.title || 'Notification',
          event.data.body || '',
          'info',
          event.data.url,
        );
      }
    };

    // Listen for messages from the service worker via BroadcastChannel (modern multi-tab communication)
    const channel = new BroadcastChannel('push-notification-channel');
    channel.onmessage = handleMessage;

    return () => {
      channel.close();
    };
  }, [showToast]);

  const hasPermission =
    typeof window !== 'undefined' ? Notification.permission === 'granted' : true;

  const shouldShowPrompt =
    isSupported && isLoggedIn && !hasPermission && !pathname?.startsWith('/admin');

  return {
    shouldShowPrompt,
    subscribeToPush,
  };
};
