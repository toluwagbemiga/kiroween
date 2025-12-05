import { gql } from '@apollo/client';

// Mutation to update notification preferences
export const UPDATE_NOTIFICATION_PREFERENCES_MUTATION = gql`
  mutation UpdateNotificationPreferences($input: NotificationPreferencesInput!) {
    updateNotificationPreferences(input: $input) {
      emailEnabled
      pushEnabled
      inAppEnabled
      channels
    }
  }
`;

// Mutation to mark notification as read
export const MARK_NOTIFICATION_READ_MUTATION = gql`
  mutation MarkNotificationRead($notificationId: ID!) {
    markNotificationRead(notificationId: $notificationId)
  }
`;

// Mutation to send a notification (admin only)
export const SEND_NOTIFICATION_MUTATION = gql`
  mutation SendNotification($input: SendNotificationInput!) {
    sendNotification(input: $input)
  }
`;
