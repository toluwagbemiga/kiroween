'use client';

import React from 'react';
import Link from 'next/link';
import { Button, Card, CardContent } from '@/components/ui';
import { ShieldExclamationIcon } from '@heroicons/react/24/outline';

export const dynamic = 'force-dynamic';

export default function UnauthorizedPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardContent className="p-8 text-center">
          <div className="inline-flex items-center justify-center h-16 w-16 rounded-full bg-red-500/20 mb-4">
            <ShieldExclamationIcon className="h-8 w-8 text-red-400" />
          </div>
          
          <h1 className="text-2xl font-bold text-white mb-2">Access Denied</h1>
          <p className="text-gray-300 mb-6">
            You don't have permission to access this page. Please contact your administrator if you believe this is an error.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Link href="/dashboard">
              <Button variant="primary">Go to Dashboard</Button>
            </Link>
            <Link href="/login">
              <Button variant="secondary">Sign Out</Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
