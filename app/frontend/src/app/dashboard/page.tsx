'use client';

import React from 'react';
import { DashboardLayout, PageContainer } from '@/components/layout';
import { Card, CardHeader, CardTitle, CardContent, Button, Badge } from '@/components/ui';
import { withAuth } from '@/components/auth';
import { Feature } from '@/components/Feature';

export const dynamic = 'force-dynamic';
import {
  UserGroupIcon,
  CreditCardIcon,
  ChartBarIcon,
  ArrowTrendingUpIcon,
} from '@heroicons/react/24/outline';

const DashboardPage: React.FC = () => {
  const stats = [
    {
      name: 'Total Users',
      value: '2,543',
      change: '+12.5%',
      trend: 'up',
      icon: UserGroupIcon,
    },
    {
      name: 'Revenue',
      value: '$45,231',
      change: '+8.2%',
      trend: 'up',
      icon: CreditCardIcon,
    },
    {
      name: 'Active Sessions',
      value: '1,234',
      change: '+4.3%',
      trend: 'up',
      icon: ChartBarIcon,
    },
    {
      name: 'Conversion Rate',
      value: '3.24%',
      change: '+0.5%',
      trend: 'up',
      icon: ArrowTrendingUpIcon,
    },
  ];

  return (
    <DashboardLayout
      title="Dashboard"
      description="Welcome back! Here's what's happening with your application."
      actions={
        <>
          <Button variant="secondary" size="md">
            Download Report
          </Button>
          <Button variant="primary" size="md">
            Create New
          </Button>
        </>
      }
    >
      <PageContainer>
        {/* Feature-flagged Bento Grid Layout */}
        <Feature 
          name="bento-grid"
          fallback={
            <Card>
              <CardContent className="p-6">
                <p className="text-gray-300">
                  The new dashboard layout is not available yet. Check back soon!
                </p>
              </CardContent>
            </Card>
          }
        >
          {/* Stats Grid */}
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map((stat) => {
              const Icon = stat.icon;
              return (
                <Card key={stat.name} variant="glass" hover>
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <p className="text-sm font-medium text-gray-300">{stat.name}</p>
                        <p className="mt-2 text-3xl font-bold text-white">{stat.value}</p>
                        <div className="mt-2 flex items-center">
                          <Badge variant="success" size="sm">
                            {stat.change}
                          </Badge>
                        </div>
                      </div>
                      <div className="p-3 rounded-lg bg-primary-500/20">
                        <Icon className="h-6 w-6 text-primary-400" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>

          {/* Recent Activity */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Recent Activity</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="flex items-center space-x-4">
                      <div className="h-10 w-10 rounded-full bg-primary-500/20 flex items-center justify-center">
                        <UserGroupIcon className="h-5 w-5 text-primary-400" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-white">New user registered</p>
                        <p className="text-xs text-gray-400">2 minutes ago</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <Button variant="secondary" size="md" className="w-full justify-start">
                    <UserGroupIcon className="h-5 w-5 mr-2" />
                    Manage Users
                  </Button>
                  <Button variant="secondary" size="md" className="w-full justify-start">
                    <CreditCardIcon className="h-5 w-5 mr-2" />
                    View Billing
                  </Button>
                  <Button variant="secondary" size="md" className="w-full justify-start">
                    <ChartBarIcon className="h-5 w-5 mr-2" />
                    Analytics Dashboard
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </Feature>
      </PageContainer>
    </DashboardLayout>
  );
};

export default withAuth(DashboardPage);
