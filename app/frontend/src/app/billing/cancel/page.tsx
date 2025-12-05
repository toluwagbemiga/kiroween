'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Card, CardContent, Button } from '@/components/ui';
import { XCircleIcon } from '@heroicons/react/24/outline';

export default function BillingCancelPage() {
  const router = useRouter();

  return (
    <DashboardLayout
      title="Payment Canceled"
      description="Your subscription payment was canceled"
    >
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardContent className="p-12 text-center">
            <div className="flex justify-center mb-6">
              <div className="p-4 bg-red-500/20 rounded-full">
                <XCircleIcon className="h-16 w-16 text-red-400" />
              </div>
            </div>

            <h1 className="text-3xl font-bold text-white mb-4">
              Payment Canceled
            </h1>

            <p className="text-gray-300 mb-8">
              Your payment was canceled and no charges were made. You can try again
              whenever you're ready.
            </p>

            <div className="flex items-center justify-center gap-4">
              <Button
                variant="primary"
                onClick={() => router.push('/billing')}
              >
                View Plans
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
