'use client';

import React, { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Card, CardContent, Button } from '@/components/ui';
import { CheckCircleIcon } from '@heroicons/react/24/outline';
import { useQuery } from '@apollo/client';
import { MY_SUBSCRIPTION_QUERY } from '@/lib/graphql/queries/billing';

export default function BillingSuccessPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const sessionId = searchParams.get('session_id');

  // Refetch subscription to get updated data
  const { refetch } = useQuery(MY_SUBSCRIPTION_QUERY);

  useEffect(() => {
    // Refetch subscription data after successful payment
    if (sessionId) {
      refetch();
    }
  }, [sessionId, refetch]);

  return (
    <DashboardLayout
      title="Payment Successful"
      description="Your subscription has been activated"
    >
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardContent className="p-12 text-center">
            <div className="flex justify-center mb-6">
              <div className="p-4 bg-green-500/20 rounded-full">
                <CheckCircleIcon className="h-16 w-16 text-green-400" />
              </div>
            </div>

            <h1 className="text-3xl font-bold text-white mb-4">
              Payment Successful!
            </h1>

            <p className="text-gray-300 mb-8">
              Thank you for subscribing. Your payment has been processed successfully
              and your subscription is now active.
            </p>

            {sessionId && (
              <p className="text-sm text-gray-400 mb-8">
                Session ID: {sessionId}
              </p>
            )}

            <div className="flex items-center justify-center gap-4">
              <Button
                variant="primary"
                onClick={() => router.push('/billing')}
              >
                View Subscription
              </Button>
              <Button
                variant="secondary"
                onClick={() => router.push('/dashboard')}
              >
                Go to Dashboard
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  );
}
