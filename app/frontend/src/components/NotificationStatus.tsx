'use client';

import React from 'react';
import { useSocket } from '@/lib/SocketProvider';
import { Badge } from '@/components/ui';

/**
 * NotificationStatus component that shows Socket.IO connection status
 * Useful for debugging and showing users when real-time notifications are active
 */
export const NotificationStatus: React.FC<{ className?: string }> = ({ className }) => {
  const { isConnected } = useSocket();

  if (!isConnected) {
    return null; // Don't show anything when disconnected (optional)
  }

  return (
    <Badge variant="success" size="sm" className={className}>
      <span className="flex items-center gap-1">
        <span className="h-2 w-2 rounded-full bg-green-400 animate-pulse" />
        Live
      </span>
    </Badge>
  );
};

/**
 * Detailed connection status for debugging
 */
export const DetailedNotificationStatus: React.FC = () => {
  const { isConnected, socket } = useSocket();

  return (
    <div className="fixed bottom-4 right-4 z-50 p-3 rounded-lg bg-white/10 backdrop-blur-md border border-white/20 text-white text-sm">
      <div className="flex items-center gap-2">
        <div
          className={`h-3 w-3 rounded-full ${
            isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
          }`}
        />
        <span className="font-medium">
          {isConnected ? 'Connected' : 'Disconnected'}
        </span>
      </div>
      {socket && (
        <div className="mt-2 text-xs text-gray-400">
          <div>Socket ID: {socket.id || 'N/A'}</div>
          <div>Transport: {socket.io.engine?.transport?.name || 'N/A'}</div>
        </div>
      )}
    </div>
  );
};
