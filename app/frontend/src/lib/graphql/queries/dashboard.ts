import { gql } from '@apollo/client';

/**
 * Query to fetch analytics summary for the current user
 * Used in the dashboard to display event statistics
 */
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

/**
 * Query to fetch the current user's subscription information
 * Used in the dashboard to display billing status
 */
export const MY_SUBSCRIPTION_QUERY = gql`
  query MySubscription {
    mySubscription {
      id
      userId
      planId
      status
      currentPeriodStart
      currentPeriodEnd
      cancelAtPeriodEnd
      stripeSubscriptionId
      createdAt
      updatedAt
      plan {
        id
        name
        description
        price
        currency
        interval
        features
        stripePriceId
        isActive
      }
    }
  }
`;

/**
 * Query to fetch user statistics
 * Used in the dashboard to display total user count
 */
export const USERS_QUERY = gql`
  query Users($limit: Int, $offset: Int) {
    users(limit: $limit, offset: $offset) {
      totalCount
      nodes {
        id
        email
        name
        createdAt
        roles {
          id
          name
        }
      }
    }
  }
`;


