self.addEventListener('install', (event) => {
    self.skipWaiting();
});

self.addEventListener('activate', (event) => {
    event.waitUntil(clients.claim());
});

const pushChannel = new BroadcastChannel('push-notification-channel');

self.addEventListener('push', function (event) {
    if (!event.data) {
        console.log('Push event but no data');
        return;
    }

    let data = {};
    try {
        data = event.data.json();
    } catch (e) {
        data = { title: 'New Notification', body: event.data.text() };
    }

    const options = {
        body: data.body || 'You have a new notification.',
        icon: '/icons/icon-192x192.png', // Ensure these icons exist or provide defaults
        badge: '/icons/badge-72x72.png',
        data: {
            url: data.url || '/'
        }
    };

    const notificationPromise = self.registration.showNotification(data.title || 'Fish Auction', options);

    // Send a message via BroadcastChannel (reliable multi-tab communication)
    pushChannel.postMessage({
        type: 'PUSH_NOTIFICATION',
        title: data.title,
        body: data.body,
        url: data.url
    });
    console.log('Push message sent via BroadcastChannel');

    event.waitUntil(notificationPromise);
});

self.addEventListener('notificationclick', function (event) {
    event.notification.close();

    // Open the URL from the notification data
    if (event.notification.data && event.notification.data.url) {
        event.waitUntil(
            clients.openWindow(event.notification.data.url)
        );
    } else {
        event.waitUntil(
            clients.openWindow('/')
        );
    }
});
