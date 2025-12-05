import { gql } from '@apollo/client';

// Mutation to create a Stripe checkout session
export const CREATE_SUBSCRIPTION_CHECKOUT_MUTATION = gql`
  mutation CreateSubscriptionCheckout($planId: ID!) {
    createSubscriptionCheckout(planId: $planId) {
      sessionId
      url
    }
  }
`;

// Mutation to update subscription to a different plan
export const UPDATE_SUBSCRIPTION_MUTATION = gql`
  mutation UpdateSubscription($planId: ID!) {
    updateSubscription(planId: $planId) {
      id
      userId
      planId
      status
      currentPeriodStart
      currentPeriodEnd
      cancelAtPeriodEnd
      plan {
        id
        name
        description
        price
        currency
        interval
        features
      }
    }
  }
`;

// Mutation to cancel subscription
export const CANCEL_SUBSCRIPTION_MUTATION = gql`
  mutation CancelSubscription {
    cancelSubscription {
      id
      userId
      planId
      status
      currentPeriodStart
      currentPeriodEnd
      cancelAtPeriodEnd
      plan {
        id
        name
      }
    }
  }
`;
