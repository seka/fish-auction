import { useState, useEffect } from 'react';

// Check if user is logged in by calling the backend
export const useAuth = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [isChecking, setIsChecking] = useState(true);

    useEffect(() => {
        // Check authentication status via API
        fetch('/api/buyers/me', {
            credentials: 'include', // Include cookies
        })
            .then(response => {
                if (response.ok) {
                    setIsLoggedIn(true);
                } else {
                    setIsLoggedIn(false);
                }
            })
            .catch(() => {
                setIsLoggedIn(false);
            })
            .finally(() => {
                setIsChecking(false);
            });
    }, []);

    return { isLoggedIn, isChecking };
};
