import { gql } from '@apollo/client';

// Mutation to track a custom event
export const TRACK_EVENT_MUTATION = gql`
  mutation TrackEvent($input: TrackEventInput!) {
    trackEvent(input: $input)
  }
`;

// Mutation to identify user with properties
export const IDENTIFY_USER_MUTATION = gql`
  mutation IdentifyUser($properties: JSON!) {
    identifyUser(properties: $properties)
  }
`;
