'use client';

import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import { io, Socket } from 'socket.io-client';
import { useAuth, getToken } from '@/contexts/AuthContext';
import { useToast } from '@/components/ui/Toast';

interface SocketContextType {
  socket: Socket | null;
  isConnected: boolean;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  isConnected: false,
});

export const useSocket = () => useContext(SocketContext);

interface NotificationPayload {
  id: string;
  title?: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
  data?: Record<string, any>;
  timestamp: string;
}

/**
 * SocketProvider component that manages Socket.IO connection
 * Connects to the notifications-service and handles real-time notifications
 */
export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, user } = useAuth();
  const { success, error, warning, info } = useToast();
  const [socket, setSocket] = useState<Socket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<Socket | null>(null);

  useEffect(() => {
    // Only connect if user is authenticated
    if (!isAuthenticated || !user) {
      // Disconnect if socket exists
      if (socketRef.current) {
        socketRef.current.disconnect();
        socketRef.current = null;
        setSocket(null);
        setIsConnected(false);
      }
      return;
    }

    // Get JWT token for authentication
    const token = getToken();
    if (!token) {
      console.warn('[Socket] No token available for Socket.IO connection');
      return;
    }

    // Socket.IO server URL from environment variable
    const socketUrl = process.env.NEXT_PUBLIC_NOTIFICATIONS_URL || 'http://localhost:8085';

    console.log('[Socket] Connecting to notifications service:', socketUrl);

    // Create Socket.IO connection with JWT authentication
    const newSocket = io(socketUrl, {
      auth: {
        token: token,
      },
      transports: ['websocket', 'polling'], // Try WebSocket first, fallback to polling
      reconnection: true,
      reconnectionDelay: 1000,
      reconnectionDelayMax: 5000,
      reconnectionAttempts: 5,
    });

    // Connection event handlers
    newSocket.on('connect', () => {
      console.log('[Socket] Connected to notifications service');
      setIsConnected(true);
    });

    newSocket.on('disconnect', (reason) => {
      console.log('[Socket] Disconnected from notifications service:', reason);
      setIsConnected(false);
    });

    newSocket.on('connect_error', (err) => {
      console.error('[Socket] Connection error:', err.message);
      setIsConnected(false);
    });

    newSocket.on('reconnect', (attemptNumber) => {
      console.log('[Socket] Reconnected after', attemptNumber, 'attempts');
      setIsConnected(true);
    });

    newSocket.on('reconnect_attempt', (attemptNumber) => {
      console.log('[Socket] Reconnection attempt', attemptNumber);
    });

    newSocket.on('reconnect_error', (err) => {
      console.error('[Socket] Reconnection error:', err.message);
    });

    newSocket.on('reconnect_failed', () => {
      console.error('[Socket] Reconnection failed after maximum attempts');
    });

    // Listen for new notifications
    newSocket.on('new_notification', (payload: NotificationPayload) => {
      console.log('[Socket] Received notification:', payload);

      // Display toast notification based on type
      const title = payload.title || 'Notification';
      const message = payload.message;

      switch (payload.type) {
        case 'success':
          success(message, title);
          break;
        case 'error':
          error(message, title);
          break;
        case 'warning':
          warning(message, title);
          break;
        case 'info':
        default:
          info(message, title);
          break;
      }

      // You can also trigger other actions here, like:
      // - Update notification count in header
      // - Play a sound
      // - Show a browser notification
      // - Update a notifications list
    });

    // Listen for notification read confirmations
    newSocket.on('notification_read', (data: { notificationId: string }) => {
      console.log('[Socket] Notification marked as read:', data.notificationId);
      // Update UI to reflect read status
    });

    // Listen for broadcast messages
    newSocket.on('broadcast', (payload: NotificationPayload) => {
      console.log('[Socket] Received broadcast:', payload);
      info(payload.message, payload.title || 'Announcement');
    });

    socketRef.current = newSocket;
    setSocket(newSocket);

    // Cleanup on unmount or when auth changes
    return () => {
      console.log('[Socket] Cleaning up socket connection');
      newSocket.disconnect();
      socketRef.current = null;
    };
  }, [isAuthenticated, user, success, error, warning, info]);

  const value: SocketContextType = {
    socket,
    isConnected,
  };

  return <SocketContext.Provider value={value}>{children}</SocketContext.Provider>;
};

/**
 * Hook to send notifications programmatically
 */
export const useSendNotification = () => {
  const { socket, isConnected } = useSocket();

  const sendNotification = (
    userId: string,
    title: string,
    message: string,
    type: 'info' | 'success' | 'warning' | 'error' = 'info',
    data?: Record<string, any>
  ) => {
    if (!socket || !isConnected) {
      console.warn('[Socket] Cannot send notification: not connected');
      return false;
    }

    socket.emit('send_notification', {
      userId,
      title,
      message,
      type,
      data,
    });

    return true;
  };

  return { sendNotification, isConnected };
};
