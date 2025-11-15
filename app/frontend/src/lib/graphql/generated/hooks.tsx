import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  JSON: { input: any; output: any; }
  Time: { input: any; output: any; }
};

export type AnalyticsSummary = {
  __typename?: 'AnalyticsSummary';
  eventsByType: Scalars['JSON']['output'];
  topEvents: Array<EventCount>;
  totalEvents: Scalars['Int']['output'];
  uniqueUsers: Scalars['Int']['output'];
};

export type AuthPayload = {
  __typename?: 'AuthPayload';
  expiresAt: Scalars['Time']['output'];
  refreshToken: Scalars['String']['output'];
  token: Scalars['String']['output'];
  user: User;
};

export type CheckoutPayload = {
  __typename?: 'CheckoutPayload';
  sessionId: Scalars['String']['output'];
  url: Scalars['String']['output'];
};

export type CreateRoleInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
  permissions: Array<Scalars['String']['input']>;
};

export type EventCount = {
  __typename?: 'EventCount';
  count: Scalars['Int']['output'];
  eventName: Scalars['String']['output'];
};

export type Feature = {
  __typename?: 'Feature';
  createdAt: Scalars['Time']['output'];
  description: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
};

export type FeatureVariant = {
  __typename?: 'FeatureVariant';
  enabled: Scalars['Boolean']['output'];
  payload?: Maybe<Scalars['JSON']['output']>;
  variantName: Scalars['String']['output'];
};

export type LlmCallInput = {
  maxTokens?: InputMaybe<Scalars['Int']['input']>;
  messages: Scalars['JSON']['input'];
  model: Scalars['String']['input'];
  temperature?: InputMaybe<Scalars['Float']['input']>;
};

export type LlmResponse = {
  __typename?: 'LLMResponse';
  content: Scalars['String']['output'];
  cost: Scalars['Float']['output'];
  finishReason: Scalars['String']['output'];
  model: Scalars['String']['output'];
  tokensUsed: Scalars['Int']['output'];
};

export type LlmUsageStats = {
  __typename?: 'LLMUsageStats';
  callsByModel: Scalars['JSON']['output'];
  totalCalls: Scalars['Int']['output'];
  totalCost: Scalars['Float']['output'];
  totalTokens: Scalars['Int']['output'];
};

export type LoginInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  assignRole: User;
  callLLM: LlmResponse;
  callPrompt: PromptResponse;
  cancelSubscription: Subscription;
  changePassword: Scalars['Boolean']['output'];
  createRole: Role;
  createSubscriptionCheckout: CheckoutPayload;
  identifyUser: Scalars['Boolean']['output'];
  login: AuthPayload;
  logout: Scalars['Boolean']['output'];
  markNotificationRead: Scalars['Boolean']['output'];
  register: AuthPayload;
  removeRole: User;
  requestPasswordReset: Scalars['Boolean']['output'];
  resetPassword: Scalars['Boolean']['output'];
  sendNotification: Scalars['Boolean']['output'];
  trackEvent: Scalars['Boolean']['output'];
  updateNotificationPreferences: NotificationPreferences;
  updateProfile: User;
  updateSubscription: Subscription;
};


export type MutationAssignRoleArgs = {
  roleId: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationCallLlmArgs = {
  input: LlmCallInput;
};


export type MutationCallPromptArgs = {
  name: Scalars['String']['input'];
  variables: Scalars['JSON']['input'];
};


export type MutationChangePasswordArgs = {
  currentPassword: Scalars['String']['input'];
  newPassword: Scalars['String']['input'];
};


export type MutationCreateRoleArgs = {
  input: CreateRoleInput;
};


export type MutationCreateSubscriptionCheckoutArgs = {
  planId: Scalars['ID']['input'];
};


export type MutationIdentifyUserArgs = {
  properties: Scalars['JSON']['input'];
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationMarkNotificationReadArgs = {
  notificationId: Scalars['ID']['input'];
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationRemoveRoleArgs = {
  roleId: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationRequestPasswordResetArgs = {
  email: Scalars['String']['input'];
};


export type MutationResetPasswordArgs = {
  newPassword: Scalars['String']['input'];
  token: Scalars['String']['input'];
};


export type MutationSendNotificationArgs = {
  input: SendNotificationInput;
};


export type MutationTrackEventArgs = {
  input: TrackEventInput;
};


export type MutationUpdateNotificationPreferencesArgs = {
  input: NotificationPreferencesInput;
};


export type MutationUpdateProfileArgs = {
  input: UpdateProfileInput;
};


export type MutationUpdateSubscriptionArgs = {
  planId: Scalars['ID']['input'];
};

export type NotificationPreferences = {
  __typename?: 'NotificationPreferences';
  channels: Scalars['JSON']['output'];
  emailEnabled: Scalars['Boolean']['output'];
  inAppEnabled: Scalars['Boolean']['output'];
  pushEnabled: Scalars['Boolean']['output'];
};

export type NotificationPreferencesInput = {
  channels?: InputMaybe<Scalars['JSON']['input']>;
  emailEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  inAppEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  pushEnabled?: InputMaybe<Scalars['Boolean']['input']>;
};

export type NotificationToken = {
  __typename?: 'NotificationToken';
  expiresAt: Scalars['Time']['output'];
  socketUrl: Scalars['String']['output'];
  token: Scalars['String']['output'];
};

export type Plan = {
  __typename?: 'Plan';
  currency: Scalars['String']['output'];
  description?: Maybe<Scalars['String']['output']>;
  features: Array<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  interval: Scalars['String']['output'];
  isActive: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  price: Scalars['Float']['output'];
  stripePriceId: Scalars['String']['output'];
};

export type PromptMetadata = {
  __typename?: 'PromptMetadata';
  description: Scalars['String']['output'];
  model: Scalars['String']['output'];
  name: Scalars['String']['output'];
  temperature: Scalars['Float']['output'];
  variables: Array<Scalars['String']['output']>;
  version: Scalars['String']['output'];
};

export type PromptResponse = {
  __typename?: 'PromptResponse';
  content: Scalars['String']['output'];
  cost: Scalars['Float']['output'];
  finishReason: Scalars['String']['output'];
  model: Scalars['String']['output'];
  tokensUsed: Scalars['Int']['output'];
};

export type Query = {
  __typename?: 'Query';
  availableFeatures: Array<Feature>;
  availablePrompts: Array<PromptMetadata>;
  billingPortalUrl: Scalars['String']['output'];
  featureVariant?: Maybe<FeatureVariant>;
  isFeatureEnabled: Scalars['Boolean']['output'];
  me: User;
  myAnalytics: AnalyticsSummary;
  myLLMUsage: LlmUsageStats;
  myNotificationPreferences: NotificationPreferences;
  myPermissions: Array<Scalars['String']['output']>;
  mySubscription?: Maybe<Subscription>;
  notificationToken: NotificationToken;
  plans: Array<Plan>;
  promptDetails?: Maybe<PromptMetadata>;
  role?: Maybe<Role>;
  roles: Array<Role>;
  subscription?: Maybe<Subscription>;
  user?: Maybe<User>;
  users: UserConnection;
};


export type QueryFeatureVariantArgs = {
  featureName: Scalars['String']['input'];
  properties?: InputMaybe<Scalars['JSON']['input']>;
};


export type QueryIsFeatureEnabledArgs = {
  featureName: Scalars['String']['input'];
  properties?: InputMaybe<Scalars['JSON']['input']>;
};


export type QueryMyAnalyticsArgs = {
  endDate?: InputMaybe<Scalars['Time']['input']>;
  startDate?: InputMaybe<Scalars['Time']['input']>;
};


export type QueryPromptDetailsArgs = {
  name: Scalars['String']['input'];
};


export type QueryRoleArgs = {
  id: Scalars['ID']['input'];
};


export type QuerySubscriptionArgs = {
  id: Scalars['ID']['input'];
};


export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};


export type QueryUsersArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};

export type RegisterInput = {
  email: Scalars['String']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  password: Scalars['String']['input'];
  teamId?: InputMaybe<Scalars['String']['input']>;
};

export type Role = {
  __typename?: 'Role';
  createdAt: Scalars['Time']['output'];
  description?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  isSystem: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  permissions: Array<Scalars['String']['output']>;
};

export type SendNotificationInput = {
  data?: InputMaybe<Scalars['JSON']['input']>;
  message: Scalars['String']['input'];
  title: Scalars['String']['input'];
  type: Scalars['String']['input'];
  userId: Scalars['ID']['input'];
};

export type Subscription = {
  __typename?: 'Subscription';
  cancelAtPeriodEnd: Scalars['Boolean']['output'];
  createdAt: Scalars['Time']['output'];
  currentPeriodEnd: Scalars['Time']['output'];
  currentPeriodStart: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  plan: Plan;
  planId: Scalars['ID']['output'];
  status: Scalars['String']['output'];
  stripeSubscriptionId: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  user: User;
  userId: Scalars['ID']['output'];
};

export type TrackEventInput = {
  eventName: Scalars['String']['input'];
  properties?: InputMaybe<Scalars['JSON']['input']>;
};

export type UpdateProfileInput = {
  email?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  enabledFeatures: Array<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  name?: Maybe<Scalars['String']['output']>;
  permissions: Array<Scalars['String']['output']>;
  roles: Array<Role>;
  subscription?: Maybe<Subscription>;
  teamId?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['Time']['output'];
};

export type UserConnection = {
  __typename?: 'UserConnection';
  nodes: Array<User>;
  totalCount: Scalars['Int']['output'];
};

export type CallPromptMutationVariables = Exact<{
  name: Scalars['String']['input'];
  variables: Scalars['JSON']['input'];
}>;


export type CallPromptMutation = { __typename?: 'Mutation', callPrompt: { __typename?: 'PromptResponse', content: string, model: string, tokensUsed: number, cost: number } };

export type LoginMutationVariables = Exact<{
  input: LoginInput;
}>;


export type LoginMutation = { __typename?: 'Mutation', login: { __typename?: 'AuthPayload', token: string, refreshToken: string, expiresAt: any, user: { __typename?: 'User', id: string, email: string, name?: string | null, roles: Array<{ __typename?: 'Role', id: string, name: string }> } } };

export type RegisterMutationVariables = Exact<{
  input: RegisterInput;
}>;


export type RegisterMutation = { __typename?: 'Mutation', register: { __typename?: 'AuthPayload', token: string, refreshToken: string, expiresAt: any, user: { __typename?: 'User', id: string, email: string, name?: string | null, roles: Array<{ __typename?: 'Role', id: string, name: string }> } } };

export type ValidateTokenQueryVariables = Exact<{ [key: string]: never; }>;


export type ValidateTokenQuery = { __typename?: 'Query', me: { __typename?: 'User', id: string, email: string, name?: string | null, roles: Array<{ __typename?: 'Role', id: string, name: string }> } };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { __typename?: 'Mutation', logout: boolean };

export type MeQueryVariables = Exact<{ [key: string]: never; }>;


export type MeQuery = { __typename?: 'Query', me: { __typename?: 'User', id: string, email: string, name?: string | null, enabledFeatures: Array<string>, roles: Array<{ __typename?: 'Role', id: string, name: string }> } };

export type TrackEventMutationVariables = Exact<{
  input: TrackEventInput;
}>;


export type TrackEventMutation = { __typename?: 'Mutation', trackEvent: boolean };

export type IdentifyUserMutationVariables = Exact<{
  properties: Scalars['JSON']['input'];
}>;


export type IdentifyUserMutation = { __typename?: 'Mutation', identifyUser: boolean };


export const CallPromptDocument = gql`
    mutation CallPrompt($name: String!, $variables: JSON!) {
  callPrompt(name: $name, variables: $variables) {
    content
    model
    tokensUsed
    cost
  }
}
    `;
export type CallPromptMutationFn = Apollo.MutationFunction<CallPromptMutation, CallPromptMutationVariables>;

/**
 * __useCallPromptMutation__
 *
 * To run a mutation, you first call `useCallPromptMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCallPromptMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [callPromptMutation, { data, loading, error }] = useCallPromptMutation({
 *   variables: {
 *      name: // value for 'name'
 *      variables: // value for 'variables'
 *   },
 * });
 */
export function useCallPromptMutation(baseOptions?: Apollo.MutationHookOptions<CallPromptMutation, CallPromptMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CallPromptMutation, CallPromptMutationVariables>(CallPromptDocument, options);
      }
export type CallPromptMutationHookResult = ReturnType<typeof useCallPromptMutation>;
export type CallPromptMutationResult = Apollo.MutationResult<CallPromptMutation>;
export type CallPromptMutationOptions = Apollo.BaseMutationOptions<CallPromptMutation, CallPromptMutationVariables>;
export const LoginDocument = gql`
    mutation Login($input: LoginInput!) {
  login(input: $input) {
    token
    refreshToken
    user {
      id
      email
      name
      roles {
        id
        name
      }
    }
    expiresAt
  }
}
    `;
export type LoginMutationFn = Apollo.MutationFunction<LoginMutation, LoginMutationVariables>;

/**
 * __useLoginMutation__
 *
 * To run a mutation, you first call `useLoginMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLoginMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [loginMutation, { data, loading, error }] = useLoginMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useLoginMutation(baseOptions?: Apollo.MutationHookOptions<LoginMutation, LoginMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LoginMutation, LoginMutationVariables>(LoginDocument, options);
      }
export type LoginMutationHookResult = ReturnType<typeof useLoginMutation>;
export type LoginMutationResult = Apollo.MutationResult<LoginMutation>;
export type LoginMutationOptions = Apollo.BaseMutationOptions<LoginMutation, LoginMutationVariables>;
export const RegisterDocument = gql`
    mutation Register($input: RegisterInput!) {
  register(input: $input) {
    token
    refreshToken
    user {
      id
      email
      name
      roles {
        id
        name
      }
    }
    expiresAt
  }
}
    `;
export type RegisterMutationFn = Apollo.MutationFunction<RegisterMutation, RegisterMutationVariables>;

/**
 * __useRegisterMutation__
 *
 * To run a mutation, you first call `useRegisterMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRegisterMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [registerMutation, { data, loading, error }] = useRegisterMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useRegisterMutation(baseOptions?: Apollo.MutationHookOptions<RegisterMutation, RegisterMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RegisterMutation, RegisterMutationVariables>(RegisterDocument, options);
      }
export type RegisterMutationHookResult = ReturnType<typeof useRegisterMutation>;
export type RegisterMutationResult = Apollo.MutationResult<RegisterMutation>;
export type RegisterMutationOptions = Apollo.BaseMutationOptions<RegisterMutation, RegisterMutationVariables>;
export const ValidateTokenDocument = gql`
    query ValidateToken {
  me {
    id
    email
    name
    roles {
      id
      name
    }
  }
}
    `;

/**
 * __useValidateTokenQuery__
 *
 * To run a query within a React component, call `useValidateTokenQuery` and pass it any options that fit your needs.
 * When your component renders, `useValidateTokenQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useValidateTokenQuery({
 *   variables: {
 *   },
 * });
 */
export function useValidateTokenQuery(baseOptions?: Apollo.QueryHookOptions<ValidateTokenQuery, ValidateTokenQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ValidateTokenQuery, ValidateTokenQueryVariables>(ValidateTokenDocument, options);
      }
export function useValidateTokenLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ValidateTokenQuery, ValidateTokenQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ValidateTokenQuery, ValidateTokenQueryVariables>(ValidateTokenDocument, options);
        }
export function useValidateTokenSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<ValidateTokenQuery, ValidateTokenQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<ValidateTokenQuery, ValidateTokenQueryVariables>(ValidateTokenDocument, options);
        }
export type ValidateTokenQueryHookResult = ReturnType<typeof useValidateTokenQuery>;
export type ValidateTokenLazyQueryHookResult = ReturnType<typeof useValidateTokenLazyQuery>;
export type ValidateTokenSuspenseQueryHookResult = ReturnType<typeof useValidateTokenSuspenseQuery>;
export type ValidateTokenQueryResult = Apollo.QueryResult<ValidateTokenQuery, ValidateTokenQueryVariables>;
export const LogoutDocument = gql`
    mutation Logout {
  logout
}
    `;
export type LogoutMutationFn = Apollo.MutationFunction<LogoutMutation, LogoutMutationVariables>;

/**
 * __useLogoutMutation__
 *
 * To run a mutation, you first call `useLogoutMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLogoutMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [logoutMutation, { data, loading, error }] = useLogoutMutation({
 *   variables: {
 *   },
 * });
 */
export function useLogoutMutation(baseOptions?: Apollo.MutationHookOptions<LogoutMutation, LogoutMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LogoutMutation, LogoutMutationVariables>(LogoutDocument, options);
      }
export type LogoutMutationHookResult = ReturnType<typeof useLogoutMutation>;
export type LogoutMutationResult = Apollo.MutationResult<LogoutMutation>;
export type LogoutMutationOptions = Apollo.BaseMutationOptions<LogoutMutation, LogoutMutationVariables>;
export const MeDocument = gql`
    query Me {
  me {
    id
    email
    name
    enabledFeatures
    roles {
      id
      name
    }
  }
}
    `;

/**
 * __useMeQuery__
 *
 * To run a query within a React component, call `useMeQuery` and pass it any options that fit your needs.
 * When your component renders, `useMeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMeQuery({
 *   variables: {
 *   },
 * });
 */
export function useMeQuery(baseOptions?: Apollo.QueryHookOptions<MeQuery, MeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MeQuery, MeQueryVariables>(MeDocument, options);
      }
export function useMeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MeQuery, MeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MeQuery, MeQueryVariables>(MeDocument, options);
        }
export function useMeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MeQuery, MeQueryVariables>(MeDocument, options);
        }
export type MeQueryHookResult = ReturnType<typeof useMeQuery>;
export type MeLazyQueryHookResult = ReturnType<typeof useMeLazyQuery>;
export type MeSuspenseQueryHookResult = ReturnType<typeof useMeSuspenseQuery>;
export type MeQueryResult = Apollo.QueryResult<MeQuery, MeQueryVariables>;
export const TrackEventDocument = gql`
    mutation TrackEvent($input: TrackEventInput!) {
  trackEvent(input: $input)
}
    `;
export type TrackEventMutationFn = Apollo.MutationFunction<TrackEventMutation, TrackEventMutationVariables>;

/**
 * __useTrackEventMutation__
 *
 * To run a mutation, you first call `useTrackEventMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useTrackEventMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [trackEventMutation, { data, loading, error }] = useTrackEventMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useTrackEventMutation(baseOptions?: Apollo.MutationHookOptions<TrackEventMutation, TrackEventMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<TrackEventMutation, TrackEventMutationVariables>(TrackEventDocument, options);
      }
export type TrackEventMutationHookResult = ReturnType<typeof useTrackEventMutation>;
export type TrackEventMutationResult = Apollo.MutationResult<TrackEventMutation>;
export type TrackEventMutationOptions = Apollo.BaseMutationOptions<TrackEventMutation, TrackEventMutationVariables>;
export const IdentifyUserDocument = gql`
    mutation IdentifyUser($properties: JSON!) {
  identifyUser(properties: $properties)
}
    `;
export type IdentifyUserMutationFn = Apollo.MutationFunction<IdentifyUserMutation, IdentifyUserMutationVariables>;

/**
 * __useIdentifyUserMutation__
 *
 * To run a mutation, you first call `useIdentifyUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useIdentifyUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [identifyUserMutation, { data, loading, error }] = useIdentifyUserMutation({
 *   variables: {
 *      properties: // value for 'properties'
 *   },
 * });
 */
export function useIdentifyUserMutation(baseOptions?: Apollo.MutationHookOptions<IdentifyUserMutation, IdentifyUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<IdentifyUserMutation, IdentifyUserMutationVariables>(IdentifyUserDocument, options);
      }
export type IdentifyUserMutationHookResult = ReturnType<typeof useIdentifyUserMutation>;
export type IdentifyUserMutationResult = Apollo.MutationResult<IdentifyUserMutation>;
export type IdentifyUserMutationOptions = Apollo.BaseMutationOptions<IdentifyUserMutation, IdentifyUserMutationVariables>;