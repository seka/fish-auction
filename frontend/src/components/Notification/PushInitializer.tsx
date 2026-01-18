'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { usePushSubscription } from '../../hooks/usePushSubscription';
import { useAuth } from '../../hooks/useAuth';
import { useToast } from '../../hooks/useToast';
import { css } from 'styled-system/css';

export const PushInitializer = () => {
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
            hasVapidKey: !!process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY
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
                    event.data.url
                );
            }
        };

        // Listen for messages from the service worker via BroadcastChannel (modern multi-tab communication)
        const channel = new BroadcastChannel('push-notification-channel');
        channel.onmessage = handleMessage;

        // Also check if controller is already sending messages
        if (navigator.serviceWorker.controller) {
            console.log('SW controller is active');
        }

        return () => {
            channel.close();
        };
    }, [showToast]);

    const hasPermission = typeof window !== 'undefined' ? Notification.permission === 'granted' : true;

    if (isSupported && isLoggedIn && !hasPermission && !pathname?.startsWith('/admin')) {
        return (
            <div className={css({
                position: 'fixed',
                bottom: '4',
                right: '4',
                bg: 'white',
                border: '1px solid',
                borderColor: 'orange.400',
                p: '4',
                borderRadius: 'lg',
                boxShadow: 'xl',
                zIndex: 2147483647,
                display: 'flex',
                flexDirection: 'column',
                gap: '3',
                maxWidth: '300px'
            })}>
                <div className={css({ display: 'flex', gap: '2', alignItems: 'flex-start' })}>
                    <span className={css({ fontSize: 'xl' })}>🔔</span>
                    <div className={css({ display: 'flex', flexDirection: 'column', gap: '1' })}>
                        <p className={css({ fontSize: 'sm', fontWeight: 'bold', color: 'gray.800' })}>
                            通知を有効にしてください
                        </p>
                        <p className={css({ fontSize: 'xs', color: 'gray.600', lineHeight: '1.4' })}>
                            セリの高値更新情報をリアルタイムで受け取るには通知の許可が必要です。
                        </p>
                    </div>
                </div>
                <button
                    onClick={() => subscribeToPush()}
                    className={css({
                        bg: 'orange.500',
                        _hover: { bg: 'orange.600' },
                        color: 'white',
                        py: '2',
                        px: '4',
                        borderRadius: 'md',
                        fontSize: 'xs',
                        fontWeight: 'bold',
                        cursor: 'pointer',
                        transition: 'background 0.2s',
                        textAlign: 'center'
                    })}
                >
                    通知を有効にする
                </button>
            </div>
        );
    }

    return null;
};
