'use client';

import {
  createContext,
  useContext,
  useState,
  useCallback,
  ReactNode,
  useRef,
  useEffect,
} from 'react';

type ToastType = 'info' | 'success' | 'warning' | 'error';

export interface Toast {
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
  const timeoutRefs = useRef<{ [key: string]: NodeJS.Timeout }>({});

  // Clean up all timeouts on unmount
  useEffect(() => {
    const currentTimeouts = timeoutRefs.current;
    return () => {
      Object.values(currentTimeouts).forEach(clearTimeout);
    };
  }, [timeoutRefs]);

  const showToast = useCallback(
    (title: string, body: string, type: ToastType = 'info', url?: string) => {
      console.log('showToast called:', { title, body, type, url });
      const id = Math.random().toString(36).substring(2, 9);
      setToasts((prev) => {
        const next = [...prev, { id, title, body, type, url }];
        console.log('Current toasts state:', next);
        return next;
      });

      // Auto remove after 5 seconds
      const timeoutId = setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
        delete timeoutRefs.current[id];
      }, 5000);

      timeoutRefs.current[id] = timeoutId;
    },
    [],
  );

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
    if (timeoutRefs.current[id]) {
      clearTimeout(timeoutRefs.current[id]);
      delete timeoutRefs.current[id];
    }
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
