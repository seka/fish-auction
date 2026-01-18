'use client';

import { createContext, useContext, useState, useCallback, ReactNode } from 'react';

type ToastType = 'info' | 'success' | 'warning' | 'error';

interface Toast {
    id: string;
    title: string;
    body: string;
    type: ToastType;
    url?: string;
}

interface ToastContextType {
    toasts: Toast[];
    showToast: (title: string, body: string, type?: ToastType, url?: string) => void;
    removeToast: (id: string) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const ToastProvider = ({ children }: { children: ReactNode }) => {
    const [toasts, setToasts] = useState<Toast[]>([]);

    const showToast = useCallback((title: string, body: string, type: ToastType = 'info', url?: string) => {
        console.log('showToast called:', { title, body, type, url });
        const id = Math.random().toString(36).substring(2, 9);
        setToasts((prev) => {
            const next = [...prev, { id, title, body, type, url }];
            console.log('Current toasts state:', next);
            return next;
        });

        // Auto remove after 5 seconds
        setTimeout(() => {
            setToasts((prev) => prev.filter((t) => t.id !== id));
        }, 5000);
    }, []);

    const removeToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
    }, []);

    return (
        <ToastContext.Provider value={{ toasts, showToast, removeToast }}>
            {children}
        </ToastContext.Provider>
    );
};

export const useToast = () => {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
};
