'use client';

import React from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Card, CardContent, Button, Loading } from '@/components/ui';
import { EmptyState } from '@/components/layout/EmptyState';
import {
  MY_NOTIFICATION_PREFERENCES_QUERY,
} from '@/lib/graphql/queries/notifications';
import {
  UPDATE_NOTIFICATION_PREFERENCES_MUTATION,
} from '@/lib/graphql/mutations/notifications';
import {
  MyNotificationPreferencesQuery,
  UpdateNotificationPreferencesMutation,
  UpdateNotificationPreferencesMutationVariables,
} from '@/lib/graphql/generated/graphql';
import { useToast } from '@/components/ui';
import { BellIcon } from '@heroicons/react/24/outline';

export default function NotificationsPage() {
  const { success, error: showError } = useToast();

  const { data, loading, error } = useQuery<MyNotificationPreferencesQuery>(
    MY_NOTIFICATION_PREFERENCES_QUERY
  );

  const [updatePreferences, { loading: updating }] = useMutation<
    UpdateNotificationPreferencesMutation,
    UpdateNotificationPreferencesMutationVariables
  >(UPDATE_NOTIFICATION_PREFERENCES_MUTATION, {
    onCompleted: () => {
      success('Notification preferences updated', 'Success');
    },
    onError: (error) => {
      showError(error.message || 'Failed to update preferences', 'Error');
    },
  });

  const handleToggle = async (field: 'emailEnabled' | 'pushEnabled' | 'inAppEnabled') => {
    if (!data?.myNotificationPreferences) return;

    const current = data.myNotificationPreferences;
    await updatePreferences({
      variables: {
        input: {
          emailEnabled: field === 'emailEnabled' ? !current.emailEnabled : current.emailEnabled,
          pushEnabled: field === 'pushEnabled' ? !current.pushEnabled : current.pushEnabled,
          inAppEnabled: field === 'inAppEnabled' ? !current.inAppEnabled : current.inAppEnabled,
        },
      },
    });
  };

  if (error) {
    return (
      <DashboardLayout
        title="Notifications"
        description="Manage your notification preferences"
      >
        <EmptyState
          title="Error loading preferences"
          description={error.message}
        />
      </DashboardLayout>
    );
  }

  const preferences = data?.myNotificationPreferences;

  return (
    <DashboardLayout
      title="Notifications"
      description="Manage your notification preferences"
    >
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loading size="lg" />
        </div>
      ) : !preferences ? (
        <EmptyState
          icon={<BellIcon className="h-12 w-12" />}
          title="No preferences found"
          description="Unable to load notification preferences"
        />
      ) : (
        <div className="max-w-2xl space-y-6">
          <Card>
            <CardContent className="p-6">
              <h2 className="text-xl font-semibold text-white mb-6">
                Notification Channels
              </h2>

              <div className="space-y-4">
                {/* Email Notifications */}
                <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg">
                  <div>
                    <h3 className="text-white font-medium">Email Notifications</h3>
                    <p className="text-sm text-gray-400">
                      Receive notifications via email
                    </p>
                  </div>
                  <Button
                    variant={preferences.emailEnabled ? 'primary' : 'secondary'}
                    size="sm"
                    onClick={() => handleToggle('emailEnabled')}
                    disabled={updating}
                  >
                    {preferences.emailEnabled ? 'Enabled' : 'Disabled'}
                  </Button>
                </div>

                {/* Push Notifications */}
                <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg">
                  <div>
                    <h3 className="text-white font-medium">Push Notifications</h3>
                    <p className="text-sm text-gray-400">
                      Receive push notifications in your browser
                    </p>
                  </div>
                  <Button
                    variant={preferences.pushEnabled ? 'primary' : 'secondary'}
                    size="sm"
                    onClick={() => handleToggle('pushEnabled')}
                    disabled={updating}
                  >
                    {preferences.pushEnabled ? 'Enabled' : 'Disabled'}
                  </Button>
                </div>

                {/* In-App Notifications */}
                <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg">
                  <div>
                    <h3 className="text-white font-medium">In-App Notifications</h3>
                    <p className="text-sm text-gray-400">
                      Show notifications within the application
                    </p>
                  </div>
                  <Button
                    variant={preferences.inAppEnabled ? 'primary' : 'secondary'}
                    size="sm"
                    onClick={() => handleToggle('inAppEnabled')}
                    disabled={updating}
                  >
                    {preferences.inAppEnabled ? 'Enabled' : 'Disabled'}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </DashboardLayout>
  );
}
