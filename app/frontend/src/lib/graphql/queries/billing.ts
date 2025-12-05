import { gql } from '@apollo/client';

// Query to fetch all available subscription plans
export const PLANS_QUERY = gql`
  query Plans {
    plans {
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
`;

// Query to fetch current user's subscription
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

// Query to get billing portal URL
export const BILLING_PORTAL_URL_QUERY = gql`
  query BillingPortalUrl {
    billingPortalUrl
  }
`;
