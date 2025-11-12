import React from 'react';
import { cn } from '@/lib/utils';
import {
  CheckCircleIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
  XCircleIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline';

export interface ToastProps {
  id: string;
  title?: string;
  message: string;
  variant?: 'success' | 'error' | 'warning' | 'info';
  duration?: number;
  onClose?: () => void;
}

const Toast: React.FC<ToastProps> = ({
  id,
  title,
  message,
  variant = 'info',
  onClose,
}) => {
  const variants = {
    success: {
      icon: CheckCircleIcon,
      styles: 'bg-green-500/20 border-green-500/30 text-green-300',
      iconColor: 'text-green-400',
    },
    error: {
      icon: XCircleIcon,
      styles: 'bg-red-500/20 border-red-500/30 text-red-300',
      iconColor: 'text-red-400',
    },
    warning: {
      icon: ExclamationCircleIcon,
      styles: 'bg-yellow-500/20 border-yellow-500/30 text-yellow-300',
      iconColor: 'text-yellow-400',
    },
    info: {
      icon: InformationCircleIcon,
      styles: 'bg-blue-500/20 border-blue-500/30 text-blue-300',
      iconColor: 'text-blue-400',
    },
  };

  const { icon: Icon, styles, iconColor } = variants[variant];

  return (
    <div
      role="alert"
      aria-live="polite"
      className={cn(
        'flex items-start gap-3 p-4 rounded-lg',
        'backdrop-blur-sm border shadow-lg',
        'animate-slide-in-right',
        'min-w-[320px] max-w-md',
        styles
      )}
    >
      <Icon className={cn('h-5 w-5 flex-shrink-0 mt-0.5', iconColor)} />
      
      <div className="flex-1 min-w-0">
        {title && (
          <h4 className="font-semibold text-white mb-1">{title}</h4>
        )}
        <p className="text-sm">{message}</p>
      </div>

      {onClose && (
        <button
          onClick={onClose}
          className="flex-shrink-0 text-white/60 hover:text-white transition-colors"
          aria-label="Close notification"
        >
          <XMarkIcon className="h-5 w-5" />
        </button>
      )}
    </div>
  );
};

// Toast Container Component
export interface ToastContainerProps {
  toasts: ToastProps[];
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center' | 'bottom-center';
}

const ToastContainer: React.FC<ToastContainerProps> = ({
  toasts,
  position = 'top-right',
}) => {
  const positions = {
    'top-right': 'top-4 right-4',
    'top-left': 'top-4 left-4',
    'bottom-right': 'bottom-4 right-4',
    'bottom-left': 'bottom-4 left-4',
    'top-center': 'top-4 left-1/2 -translate-x-1/2',
    'bottom-center': 'bottom-4 left-1/2 -translate-x-1/2',
  };

  return (
    <div
      className={cn(
        'fixed z-50 flex flex-col gap-2',
        positions[position]
      )}
      aria-live="polite"
      aria-atomic="true"
    >
      {toasts.map((toast) => (
        <Toast key={toast.id} {...toast} />
      ))}
    </div>
  );
};

// Toast Hook
interface ToastOptions {
  title?: string;
  message: string;
  variant?: 'success' | 'error' | 'warning' | 'info';
  duration?: number;
}

let toastId = 0;

export const useToast = () => {
  const [toasts, setToasts] = React.useState<ToastProps[]>([]);

  const addToast = React.useCallback((options: ToastOptions) => {
    const id = `toast-${++toastId}`;
    const duration = options.duration ?? 5000;

    const toast: ToastProps = {
      id,
      ...options,
      onClose: () => removeToast(id),
    };

    setToasts((prev) => [...prev, toast]);

    if (duration > 0) {
      setTimeout(() => removeToast(id), duration);
    }

    return id;
  }, []);

  const removeToast = React.useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  const success = React.useCallback(
    (message: string, title?: string, duration?: number) =>
      addToast({ message, title, variant: 'success', duration }),
    [addToast]
  );

  const error = React.useCallback(
    (message: string, title?: string, duration?: number) =>
      addToast({ message, title, variant: 'error', duration }),
    [addToast]
  );

  const warning = React.useCallback(
    (message: string, title?: string, duration?: number) =>
      addToast({ message, title, variant: 'warning', duration }),
    [addToast]
  );

  const info = React.useCallback(
    (message: string, title?: string, duration?: number) =>
      addToast({ message, title, variant: 'info', duration }),
    [addToast]
  );

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info,
  };
};

export { Toast, ToastContainer };
