import { gql } from '@apollo/client';

// Query to fetch analytics summary for the current user
export const MY_ANALYTICS_QUERY = gql`
  query MyAnalytics($startDate: Time, $endDate: Time) {
    myAnalytics(startDate: $startDate, endDate: $endDate) {
      totalEvents
      uniqueUsers
      eventsByType
      topEvents {
        eventName
        count
      }
    }
  }
`;
