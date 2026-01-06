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

    event.waitUntil(
        self.registration.showNotification(data.title || 'Fish Auction', options)
    );
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
