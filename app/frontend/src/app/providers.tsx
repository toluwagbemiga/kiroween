'use client';

import React from 'react';
import { ApolloProvider } from '@apollo/client';
import { apolloClient } from '@/lib/apollo-client';
import { AuthProvider } from '@/contexts/AuthContext';
import { FeatureFlagProvider } from '@/contexts/FeatureFlagContext';
import { SocketProvider } from '@/lib/SocketProvider';
import { ToastContainer, useToast } from '@/components/ui';

/**
 * ToastProvider component that manages toast notifications
 */
const ToastProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { toasts } = useToast();

  return (
    <>
      {children}
      <ToastContainer toasts={toasts} position="top-right" />
    </>
  );
};

/**
 * Root providers component that wraps the entire application
 * Includes Apollo Client, Authentication, Feature Flags, Socket.IO, and Toast notifications
 */
export const Providers: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <ApolloProvider client={apolloClient}>
      <AuthProvider>
        <FeatureFlagProvider>
          <SocketProvider>
            <ToastProvider>{children}</ToastProvider>
          </SocketProvider>
        </FeatureFlagProvider>
      </AuthProvider>
    </ApolloProvider>
  );
};
