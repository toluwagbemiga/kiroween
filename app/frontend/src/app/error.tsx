'use client';

import React, { useEffect } from 'react';
import { Button, Card, CardContent } from '@/components/ui';
import { ExclamationCircleIcon } from '@heroicons/react/24/outline';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Log the error to an error reporting service
    console.error('Application error:', error);
  }, [error]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardContent className="p-8 text-center">
          <div className="inline-flex items-center justify-center h-16 w-16 rounded-full bg-red-500/20 mb-4">
            <ExclamationCircleIcon className="h-8 w-8 text-red-400" />
          </div>
          
          <h1 className="text-2xl font-bold text-white mb-2">Something went wrong!</h1>
          <p className="text-gray-300 mb-6">
            An unexpected error occurred. Please try again or contact support if the problem persists.
          </p>
          
          {error.message && (
            <div className="mb-6 p-3 rounded-lg bg-red-500/10 border border-red-500/20">
              <p className="text-sm text-red-300 font-mono">{error.message}</p>
            </div>
          )}
          
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Button variant="primary" onClick={reset}>
              Try Again
            </Button>
            <Button variant="secondary" onClick={() => window.location.href = '/dashboard'}>
              Go to Dashboard
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
