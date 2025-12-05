'use client';

import React from 'react';
import { DashboardLayout, PageContainer } from '@/components/layout';
import { Card, CardHeader, CardTitle, CardContent, Button, Badge, StatCardSkeleton, ActivityItemSkeleton } from '@/components/ui';
import { withAuth } from '@/components/auth';
import { ErrorBoundary } from '@/components/ErrorBoundary';
import { useMyAnalyticsQuery, useMySubscriptionQuery, useUsersQuery } from '@/lib/graphql/generated/hooks';
import { AnalyticsChart } from '@/components/charts';

export const dynamic = 'force-dynamic';
import {
  UserGroupIcon,
  CreditCardIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

const DashboardPage: React.FC = () => {
  // Fetch analytics data
  const { data: analyticsData, loading: analyticsLoading, error: analyticsError } = useMyAnalyticsQuery({
    variables: {
      startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(), // Last 30 days
      endDate: new Date().toISOString(),
    },
  });

  // Fetch subscription data
  const { data: subscriptionData, loading: subscriptionLoading, error: subscriptionError } = useMySubscriptionQuery();

  // Fetch users count
  const { data: usersData, loading: usersLoading, error: usersError } = useUsersQuery({
    variables: {
      limit: 1,
      offset: 0,
    },
  });

  // Calculate stats from real data
  const stats = [
    {
      name: 'Total Users',
      value: usersLoading ? '...' : usersData?.users?.totalCount?.toString() || '0',
      change: '+12.5%',
      trend: 'up',
      icon: UserGroupIcon,
      loading: usersLoading,
      error: usersError,
    },
    {
      name: 'Total Events',
      value: analyticsLoading ? '...' : analyticsData?.myAnalytics?.totalEvents?.toLocaleString() || '0',
      change: '+8.2%',
      trend: 'up',
      icon: ChartBarIcon,
      loading: analyticsLoading,
      error: analyticsError,
    },
    {
      name: 'Unique Users',
      value: analyticsLoading ? '...' : analyticsData?.myAnalytics?.uniqueUsers?.toString() || '0',
      change: '+4.3%',
      trend: 'up',
      icon: UserGroupIcon,
      loading: analyticsLoading,
      error: analyticsError,
    },
    {
      name: 'Subscription',
      value: subscriptionLoading ? '...' : subscriptionData?.mySubscription?.plan?.name || 'Free',
      change: subscriptionData?.mySubscription?.status || 'inactive',
      trend: subscriptionData?.mySubscription?.status === 'active' ? 'up' : 'down',
      icon: CreditCardIcon,
      loading: subscriptionLoading,
      error: subscriptionError,
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
        {/* Dashboard Content */}
        <div>
          {/* Stats Grid */}
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map((stat) => {
              const Icon = stat.icon;
              return (
                <Card key={stat.name} variant="glass" hover>
                  <ErrorBoundary>
                    {stat.loading ? (
                      <StatCardSkeleton />
                    ) : stat.error ? (
                      <CardContent className="p-6">
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <p className="text-sm font-medium text-gray-300">{stat.name}</p>
                            <p className="mt-2 text-sm text-red-400">Error loading data</p>
                          </div>
                          <div className="p-3 rounded-lg bg-red-500/20">
                            <Icon className="h-6 w-6 text-red-400" />
                          </div>
                        </div>
                      </CardContent>
                    ) : (
                      <CardContent className="p-6">
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <p className="text-sm font-medium text-gray-300">{stat.name}</p>
                            <p className="mt-2 text-3xl font-bold text-white">{stat.value}</p>
                            <div className="mt-2 flex items-center">
                              <Badge 
                                variant={stat.trend === 'up' ? 'success' : 'warning'} 
                                size="sm"
                              >
                                {stat.change}
                              </Badge>
                            </div>
                          </div>
                          <div className="p-3 rounded-lg bg-primary-500/20">
                            <Icon className="h-6 w-6 text-primary-400" />
                          </div>
                        </div>
                      </CardContent>
                    )}
                  </ErrorBoundary>
                </Card>
              );
            })}
          </div>

          {/* Event Trends Chart */}
          <Card className="mt-6">
            <CardHeader>
              <CardTitle>Event Distribution</CardTitle>
            </CardHeader>
            <CardContent>
              <ErrorBoundary>
                {analyticsLoading ? (
                  <div className="h-[300px] flex items-center justify-center">
                    <div className="animate-pulse text-gray-400">Loading chart...</div>
                  </div>
                ) : analyticsError ? (
                  <div className="h-[300px] flex items-center justify-center">
                    <p className="text-sm text-red-400">Error loading chart data</p>
                  </div>
                ) : analyticsData?.myAnalytics?.topEvents?.length ? (
                  <AnalyticsChart
                    data={analyticsData.myAnalytics.topEvents.map(event => ({
                      name: event.eventName,
                      value: event.count,
                    }))}
                    type="bar"
                    dataKey="value"
                    xAxisKey="name"
                    height={300}
                  />
                ) : (
                  <div className="h-[300px] flex items-center justify-center">
                    <p className="text-sm text-gray-400">No event data available</p>
                  </div>
                )}
              </ErrorBoundary>
            </CardContent>
          </Card>

          {/* Recent Activity */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Top Events</CardTitle>
              </CardHeader>
              <CardContent>
                <ErrorBoundary>
                  {analyticsLoading ? (
                    <div className="space-y-4">
                      {[1, 2, 3, 4].map((i) => (
                        <ActivityItemSkeleton key={i} />
                      ))}
                    </div>
                  ) : analyticsError ? (
                    <div className="text-center py-8">
                      <p className="text-sm text-red-400">Error loading events</p>
                    </div>
                  ) : analyticsData?.myAnalytics?.topEvents?.length ? (
                    <div className="space-y-4">
                      {analyticsData.myAnalytics.topEvents.slice(0, 5).map((event, i) => (
                        <div key={i} className="flex items-center space-x-4">
                          <div className="h-10 w-10 rounded-full bg-primary-500/20 flex items-center justify-center">
                            <ChartBarIcon className="h-5 w-5 text-primary-400" />
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium text-white">{event.eventName}</p>
                            <p className="text-xs text-gray-400">{event.count} occurrences</p>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <p className="text-sm text-gray-400">No events tracked yet</p>
                    </div>
                  )}
                </ErrorBoundary>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Subscription Status</CardTitle>
              </CardHeader>
              <CardContent>
                <ErrorBoundary>
                  {subscriptionLoading ? (
                    <div className="space-y-3">
                      <ActivityItemSkeleton />
                      <ActivityItemSkeleton />
                    </div>
                  ) : subscriptionError ? (
                    <div className="text-center py-4">
                      <p className="text-sm text-red-400">Error loading subscription</p>
                    </div>
                  ) : subscriptionData?.mySubscription ? (
                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium text-gray-300">Current Plan</p>
                          <p className="text-lg font-bold text-white">
                            {subscriptionData.mySubscription.plan.name}
                          </p>
                        </div>
                        <Badge 
                          variant={subscriptionData.mySubscription.status === 'active' ? 'success' : 'warning'}
                          size="md"
                        >
                          {subscriptionData.mySubscription.status}
                        </Badge>
                      </div>
                      <div className="pt-2 border-t border-gray-700">
                        <p className="text-xs text-gray-400 mb-1">Price</p>
                        <p className="text-sm font-medium text-white">
                          {subscriptionData.mySubscription.plan.currency.toUpperCase()} {subscriptionData.mySubscription.plan.price} / {subscriptionData.mySubscription.plan.interval}
                        </p>
                      </div>
                      <div className="pt-2">
                        <Button variant="primary" size="sm" className="w-full">
                          <CreditCardIcon className="h-4 w-4 mr-2" />
                          Manage Subscription
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      <p className="text-sm text-gray-400">No active subscription</p>
                      <Button variant="primary" size="sm" className="w-full">
                        <CreditCardIcon className="h-4 w-4 mr-2" />
                        View Plans
                      </Button>
                    </div>
                  )}
                </ErrorBoundary>
              </CardContent>
            </Card>
          </div>
        </div>
      </PageContainer>
    </DashboardLayout>
  );
};

export default withAuth(DashboardPage);
