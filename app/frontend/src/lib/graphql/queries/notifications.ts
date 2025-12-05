import { gql } from '@apollo/client';

// Query to get notification preferences
export const MY_NOTIFICATION_PREFERENCES_QUERY = gql`
  query MyNotificationPreferences {
    myNotificationPreferences {
      emailEnabled
      pushEnabled
      inAppEnabled
      channels
    }
  }
`;

// Query to get notification token for Socket.IO
export const NOTIFICATION_TOKEN_QUERY = gql`
  query NotificationToken {
    notificationToken {
      token
      socketUrl
      expiresAt
    }
  }
`;
