'use client';

import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { DateRangePicker, DateRange } from '@/components/analytics/DateRangePicker';
import { EventTable } from '@/components/analytics/EventTable';
import { MetricCard } from '@/components/charts/MetricCard';
import { AnalyticsChart } from '@/components/charts/AnalyticsChart';
import { Card, CardContent, Loading } from '@/components/ui';
import { EmptyState } from '@/components/layout/EmptyState';
import { MY_ANALYTICS_QUERY } from '@/lib/graphql/queries/analytics';
import { MyAnalyticsQuery, MyAnalyticsQueryVariables } from '@/lib/graphql/generated/graphql';
import { subDays, startOfDay, endOfDay } from 'date-fns';
import { ChartBarIcon } from '@heroicons/react/24/outline';

export default function AnalyticsPage() {
  const [dateRange, setDateRange] = useState<DateRange>({
    startDate: startOfDay(subDays(new Date(), 30)),
    endDate: endOfDay(new Date()),
  });

  const { data, loading, error } = useQuery<MyAnalyticsQuery, MyAnalyticsQueryVariables>(
    MY_ANALYTICS_QUERY,
    {
      variables: {
        startDate: dateRange.startDate.toISOString(),
        endDate: dateRange.endDate.toISOString(),
      },
      fetchPolicy: 'cache-and-network',
    }
  );

  if (error) {
    return (
      <DashboardLayout
        title="Analytics"
        description="Track and analyze your usage"
      >
        <EmptyState
          title="Error loading analytics"
          description={error.message}
        />
      </DashboardLayout>
    );
  }

  const analytics = data?.myAnalytics;

  return (
    <DashboardLayout
      title="Analytics"
      description="Track and analyze your usage"
      actions={
        <DateRangePicker value={dateRange} onChange={setDateRange} />
      }
    >
      {loading && !data ? (
        <div className="flex items-center justify-center py-12">
          <Loading size="lg" />
        </div>
      ) : !analytics ? (
        <EmptyState
          icon={<ChartBarIcon className="h-12 w-12" />}
          title="No analytics data"
          description="Start using the platform to see analytics"
        />
      ) : (
        <div className="space-y-6">
          {/* Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <MetricCard
              title="Total Events"
              value={analytics.totalEvents}
              trend="neutral"
            />
            <MetricCard
              title="Unique Users"
              value={analytics.uniqueUsers}
              trend="neutral"
            />
            <MetricCard
              title="Event Types"
              value={Object.keys(analytics.eventsByType || {}).length}
              trend="neutral"
            />
          </div>

          {/* Charts */}
          {analytics.topEvents && analytics.topEvents.length > 0 && (
            <Card>
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold text-white mb-4">
                  Top Events
                </h3>
                <AnalyticsChart
                  data={analytics.topEvents.map((event) => ({
                    name: event.eventName,
                    value: event.count,
                  }))}
                  type="bar"
                />
              </CardContent>
            </Card>
          )}

          {/* Event Table */}
          {analytics.topEvents && analytics.topEvents.length > 0 && (
            <Card>
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold text-white mb-4">
                  Event Breakdown
                </h3>
                <EventTable events={analytics.topEvents} />
              </CardContent>
            </Card>
          )}
        </div>
      )}
    </DashboardLayout>
  );
}
