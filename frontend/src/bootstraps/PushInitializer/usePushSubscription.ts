'use client';

import { useState, useEffect } from 'react';
import { urlBase64ToUint8Array } from '@/src/utils/webPush';

const VAPID_PUBLIC_KEY = process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY;

export const usePushSubscription = () => {
  const [subscription, setSubscription] = useState<PushSubscription | null>(null);
  const [isSupported, setIsSupported] = useState(false);

  useEffect(() => {
    const registerServiceWorker = async () => {
      try {
        const registration = await navigator.serviceWorker.register('/sw.js');
        console.log('Service Worker registered:', registration);

        const sub = await registration.pushManager.getSubscription();
        setSubscription(sub);
      } catch (error) {
        console.error('Service Worker registration failed:', error);
      }
    };

    if ('serviceWorker' in navigator && 'PushManager' in window) {
      // ブラウザ機能検出はSSRでは行えないため、hydration後にeffect内で反映する
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setIsSupported(true);
      registerServiceWorker();
    }
  }, []);

  const subscribeToPush = async () => {
    if (!isSupported || !VAPID_PUBLIC_KEY) {
      console.error('Push notifications are not supported or VAPID key is missing');
      return;
    }

    try {
      const registration = await navigator.serviceWorker.ready;
      const sub = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(VAPID_PUBLIC_KEY),
      });

      setSubscription(sub);
      await sendSubscriptionToBackend(sub);
      console.log('Subscribed to push notifications:', sub);
    } catch (error) {
      console.error('Failed to subscribe to push notifications:', error);
    }
  };

  const sendSubscriptionToBackend = async (subscription: PushSubscription) => {
    try {
      // Assuming we have a way to make authenticated requests.
      // If we are using fetch directly, we rely on cookies being sent automatically.
      const response = await fetch('/api/buyer/push/subscribe', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscription),
      });

      if (!response.ok) {
        throw new Error('Failed to send subscription to backend');
      }
    } catch (error) {
      console.error('Error sending subscription to backend:', error);
    }
  };

  return {
    subscription,
    isSupported,
    subscribeToPush,
  };
};
