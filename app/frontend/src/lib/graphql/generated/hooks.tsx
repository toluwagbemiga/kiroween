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

export type MyAnalyticsQueryVariables = Exact<{
  startDate?: InputMaybe<Scalars['Time']['input']>;
  endDate?: InputMaybe<Scalars['Time']['input']>;
}>;


export type MyAnalyticsQuery = { __typename?: 'Query', myAnalytics: { __typename?: 'AnalyticsSummary', totalEvents: number, uniqueUsers: number, eventsByType: any, topEvents: Array<{ __typename?: 'EventCount', eventName: string, count: number }> } };

export type MySubscriptionQueryVariables = Exact<{ [key: string]: never; }>;


export type MySubscriptionQuery = { __typename?: 'Query', mySubscription?: { __typename?: 'Subscription', id: string, userId: string, planId: string, status: string, currentPeriodStart: any, currentPeriodEnd: any, cancelAtPeriodEnd: boolean, stripeSubscriptionId: string, createdAt: any, updatedAt: any, plan: { __typename?: 'Plan', id: string, name: string, description?: string | null, price: number, currency: string, interval: string, features: Array<string>, stripePriceId: string, isActive: boolean } } | null };

export type UsersQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
}>;


export type UsersQuery = { __typename?: 'Query', users: { __typename?: 'UserConnection', totalCount: number, nodes: Array<{ __typename?: 'User', id: string, email: string, name?: string | null, createdAt: any, roles: Array<{ __typename?: 'Role', id: string, name: string }> }> } };

export type UpdateProfileMutationVariables = Exact<{
  input: UpdateProfileInput;
}>;


export type UpdateProfileMutation = { __typename?: 'Mutation', updateProfile: { __typename?: 'User', id: string, email: string, name?: string | null, updatedAt: any } };

export type AssignRoleMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  roleId: Scalars['ID']['input'];
}>;


export type AssignRoleMutation = { __typename?: 'Mutation', assignRole: { __typename?: 'User', id: string, email: string, name?: string | null, permissions: Array<string>, roles: Array<{ __typename?: 'Role', id: string, name: string, description?: string | null }> } };

export type RemoveRoleMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  roleId: Scalars['ID']['input'];
}>;


export type RemoveRoleMutation = { __typename?: 'Mutation', removeRole: { __typename?: 'User', id: string, email: string, name?: string | null, permissions: Array<string>, roles: Array<{ __typename?: 'Role', id: string, name: string, description?: string | null }> } };

export type CreateRoleMutationVariables = Exact<{
  input: CreateRoleInput;
}>;


export type CreateRoleMutation = { __typename?: 'Mutation', createRole: { __typename?: 'Role', id: string, name: string, description?: string | null, permissions: Array<string>, isSystem: boolean, createdAt: any } };

export type UsersListQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
}>;


export type UsersListQuery = { __typename?: 'Query', users: { __typename?: 'UserConnection', totalCount: number, nodes: Array<{ __typename?: 'User', id: string, email: string, name?: string | null, createdAt: any, updatedAt: any, roles: Array<{ __typename?: 'Role', id: string, name: string, description?: string | null }> }> } };

export type UserQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type UserQuery = { __typename?: 'Query', user?: { __typename?: 'User', id: string, email: string, name?: string | null, teamId?: string | null, createdAt: any, updatedAt: any, permissions: Array<string>, roles: Array<{ __typename?: 'Role', id: string, name: string, description?: string | null, permissions: Array<string> }> } | null };

export type RolesQueryVariables = Exact<{ [key: string]: never; }>;


export type RolesQuery = { __typename?: 'Query', roles: Array<{ __typename?: 'Role', id: string, name: string, description?: string | null, permissions: Array<string>, isSystem: boolean, createdAt: any }> };

export type PlansQueryVariables = Exact<{ [key: string]: never; }>;


export type PlansQuery = { __typename?: 'Query', plans: Array<{ __typename?: 'Plan', id: string, name: string, description?: string | null, price: number, currency: string, interval: string, features: Array<string>, stripePriceId: string, isActive: boolean }> };

export type BillingPortalUrlQueryVariables = Exact<{ [key: string]: never; }>;


export type BillingPortalUrlQuery = { __typename?: 'Query', billingPortalUrl: string };

export type CreateSubscriptionCheckoutMutationVariables = Exact<{
  planId: Scalars['ID']['input'];
}>;


export type CreateSubscriptionCheckoutMutation = { __typename?: 'Mutation', createSubscriptionCheckout: { __typename?: 'CheckoutPayload', sessionId: string, url: string } };

export type UpdateSubscriptionMutationVariables = Exact<{
  planId: Scalars['ID']['input'];
}>;


export type UpdateSubscriptionMutation = { __typename?: 'Mutation', updateSubscription: { __typename?: 'Subscription', id: string, userId: string, planId: string, status: string, currentPeriodStart: any, currentPeriodEnd: any, cancelAtPeriodEnd: boolean, plan: { __typename?: 'Plan', id: string, name: string, description?: string | null, price: number, currency: string, interval: string, features: Array<string> } } };

export type CancelSubscriptionMutationVariables = Exact<{ [key: string]: never; }>;


export type CancelSubscriptionMutation = { __typename?: 'Mutation', cancelSubscription: { __typename?: 'Subscription', id: string, userId: string, planId: string, status: string, currentPeriodStart: any, currentPeriodEnd: any, cancelAtPeriodEnd: boolean, plan: { __typename?: 'Plan', id: string, name: string } } };

export type UpdateNotificationPreferencesMutationVariables = Exact<{
  input: NotificationPreferencesInput;
}>;


export type UpdateNotificationPreferencesMutation = { __typename?: 'Mutation', updateNotificationPreferences: { __typename?: 'NotificationPreferences', emailEnabled: boolean, pushEnabled: boolean, inAppEnabled: boolean, channels: any } };

export type MarkNotificationReadMutationVariables = Exact<{
  notificationId: Scalars['ID']['input'];
}>;


export type MarkNotificationReadMutation = { __typename?: 'Mutation', markNotificationRead: boolean };

export type SendNotificationMutationVariables = Exact<{
  input: SendNotificationInput;
}>;


export type SendNotificationMutation = { __typename?: 'Mutation', sendNotification: boolean };

export type MyNotificationPreferencesQueryVariables = Exact<{ [key: string]: never; }>;


export type MyNotificationPreferencesQuery = { __typename?: 'Query', myNotificationPreferences: { __typename?: 'NotificationPreferences', emailEnabled: boolean, pushEnabled: boolean, inAppEnabled: boolean, channels: any } };

export type NotificationTokenQueryVariables = Exact<{ [key: string]: never; }>;


export type NotificationTokenQuery = { __typename?: 'Query', notificationToken: { __typename?: 'NotificationToken', token: string, socketUrl: string, expiresAt: any } };


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
export const MyAnalyticsDocument = gql`
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
 * __useMyAnalyticsQuery__
 *
 * To run a query within a React component, call `useMyAnalyticsQuery` and pass it any options that fit your needs.
 * When your component renders, `useMyAnalyticsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMyAnalyticsQuery({
 *   variables: {
 *      startDate: // value for 'startDate'
 *      endDate: // value for 'endDate'
 *   },
 * });
 */
export function useMyAnalyticsQuery(baseOptions?: Apollo.QueryHookOptions<MyAnalyticsQuery, MyAnalyticsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MyAnalyticsQuery, MyAnalyticsQueryVariables>(MyAnalyticsDocument, options);
      }
export function useMyAnalyticsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MyAnalyticsQuery, MyAnalyticsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MyAnalyticsQuery, MyAnalyticsQueryVariables>(MyAnalyticsDocument, options);
        }
export function useMyAnalyticsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MyAnalyticsQuery, MyAnalyticsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MyAnalyticsQuery, MyAnalyticsQueryVariables>(MyAnalyticsDocument, options);
        }
export type MyAnalyticsQueryHookResult = ReturnType<typeof useMyAnalyticsQuery>;
export type MyAnalyticsLazyQueryHookResult = ReturnType<typeof useMyAnalyticsLazyQuery>;
export type MyAnalyticsSuspenseQueryHookResult = ReturnType<typeof useMyAnalyticsSuspenseQuery>;
export type MyAnalyticsQueryResult = Apollo.QueryResult<MyAnalyticsQuery, MyAnalyticsQueryVariables>;
export const MySubscriptionDocument = gql`
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
 * __useMySubscriptionQuery__
 *
 * To run a query within a React component, call `useMySubscriptionQuery` and pass it any options that fit your needs.
 * When your component renders, `useMySubscriptionQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMySubscriptionQuery({
 *   variables: {
 *   },
 * });
 */
export function useMySubscriptionQuery(baseOptions?: Apollo.QueryHookOptions<MySubscriptionQuery, MySubscriptionQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MySubscriptionQuery, MySubscriptionQueryVariables>(MySubscriptionDocument, options);
      }
export function useMySubscriptionLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MySubscriptionQuery, MySubscriptionQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MySubscriptionQuery, MySubscriptionQueryVariables>(MySubscriptionDocument, options);
        }
export function useMySubscriptionSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MySubscriptionQuery, MySubscriptionQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MySubscriptionQuery, MySubscriptionQueryVariables>(MySubscriptionDocument, options);
        }
export type MySubscriptionQueryHookResult = ReturnType<typeof useMySubscriptionQuery>;
export type MySubscriptionLazyQueryHookResult = ReturnType<typeof useMySubscriptionLazyQuery>;
export type MySubscriptionSuspenseQueryHookResult = ReturnType<typeof useMySubscriptionSuspenseQuery>;
export type MySubscriptionQueryResult = Apollo.QueryResult<MySubscriptionQuery, MySubscriptionQueryVariables>;
export const UsersDocument = gql`
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

/**
 * __useUsersQuery__
 *
 * To run a query within a React component, call `useUsersQuery` and pass it any options that fit your needs.
 * When your component renders, `useUsersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useUsersQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *      offset: // value for 'offset'
 *   },
 * });
 */
export function useUsersQuery(baseOptions?: Apollo.QueryHookOptions<UsersQuery, UsersQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<UsersQuery, UsersQueryVariables>(UsersDocument, options);
      }
export function useUsersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<UsersQuery, UsersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<UsersQuery, UsersQueryVariables>(UsersDocument, options);
        }
export function useUsersSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<UsersQuery, UsersQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<UsersQuery, UsersQueryVariables>(UsersDocument, options);
        }
export type UsersQueryHookResult = ReturnType<typeof useUsersQuery>;
export type UsersLazyQueryHookResult = ReturnType<typeof useUsersLazyQuery>;
export type UsersSuspenseQueryHookResult = ReturnType<typeof useUsersSuspenseQuery>;
export type UsersQueryResult = Apollo.QueryResult<UsersQuery, UsersQueryVariables>;
export const UpdateProfileDocument = gql`
    mutation UpdateProfile($input: UpdateProfileInput!) {
  updateProfile(input: $input) {
    id
    email
    name
    updatedAt
  }
}
    `;
export type UpdateProfileMutationFn = Apollo.MutationFunction<UpdateProfileMutation, UpdateProfileMutationVariables>;

/**
 * __useUpdateProfileMutation__
 *
 * To run a mutation, you first call `useUpdateProfileMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateProfileMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateProfileMutation, { data, loading, error }] = useUpdateProfileMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateProfileMutation(baseOptions?: Apollo.MutationHookOptions<UpdateProfileMutation, UpdateProfileMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateProfileMutation, UpdateProfileMutationVariables>(UpdateProfileDocument, options);
      }
export type UpdateProfileMutationHookResult = ReturnType<typeof useUpdateProfileMutation>;
export type UpdateProfileMutationResult = Apollo.MutationResult<UpdateProfileMutation>;
export type UpdateProfileMutationOptions = Apollo.BaseMutationOptions<UpdateProfileMutation, UpdateProfileMutationVariables>;
export const AssignRoleDocument = gql`
    mutation AssignRole($userId: ID!, $roleId: ID!) {
  assignRole(userId: $userId, roleId: $roleId) {
    id
    email
    name
    roles {
      id
      name
      description
    }
    permissions
  }
}
    `;
export type AssignRoleMutationFn = Apollo.MutationFunction<AssignRoleMutation, AssignRoleMutationVariables>;

/**
 * __useAssignRoleMutation__
 *
 * To run a mutation, you first call `useAssignRoleMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAssignRoleMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [assignRoleMutation, { data, loading, error }] = useAssignRoleMutation({
 *   variables: {
 *      userId: // value for 'userId'
 *      roleId: // value for 'roleId'
 *   },
 * });
 */
export function useAssignRoleMutation(baseOptions?: Apollo.MutationHookOptions<AssignRoleMutation, AssignRoleMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AssignRoleMutation, AssignRoleMutationVariables>(AssignRoleDocument, options);
      }
export type AssignRoleMutationHookResult = ReturnType<typeof useAssignRoleMutation>;
export type AssignRoleMutationResult = Apollo.MutationResult<AssignRoleMutation>;
export type AssignRoleMutationOptions = Apollo.BaseMutationOptions<AssignRoleMutation, AssignRoleMutationVariables>;
export const RemoveRoleDocument = gql`
    mutation RemoveRole($userId: ID!, $roleId: ID!) {
  removeRole(userId: $userId, roleId: $roleId) {
    id
    email
    name
    roles {
      id
      name
      description
    }
    permissions
  }
}
    `;
export type RemoveRoleMutationFn = Apollo.MutationFunction<RemoveRoleMutation, RemoveRoleMutationVariables>;

/**
 * __useRemoveRoleMutation__
 *
 * To run a mutation, you first call `useRemoveRoleMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveRoleMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeRoleMutation, { data, loading, error }] = useRemoveRoleMutation({
 *   variables: {
 *      userId: // value for 'userId'
 *      roleId: // value for 'roleId'
 *   },
 * });
 */
export function useRemoveRoleMutation(baseOptions?: Apollo.MutationHookOptions<RemoveRoleMutation, RemoveRoleMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveRoleMutation, RemoveRoleMutationVariables>(RemoveRoleDocument, options);
      }
export type RemoveRoleMutationHookResult = ReturnType<typeof useRemoveRoleMutation>;
export type RemoveRoleMutationResult = Apollo.MutationResult<RemoveRoleMutation>;
export type RemoveRoleMutationOptions = Apollo.BaseMutationOptions<RemoveRoleMutation, RemoveRoleMutationVariables>;
export const CreateRoleDocument = gql`
    mutation CreateRole($input: CreateRoleInput!) {
  createRole(input: $input) {
    id
    name
    description
    permissions
    isSystem
    createdAt
  }
}
    `;
export type CreateRoleMutationFn = Apollo.MutationFunction<CreateRoleMutation, CreateRoleMutationVariables>;

/**
 * __useCreateRoleMutation__
 *
 * To run a mutation, you first call `useCreateRoleMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateRoleMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createRoleMutation, { data, loading, error }] = useCreateRoleMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateRoleMutation(baseOptions?: Apollo.MutationHookOptions<CreateRoleMutation, CreateRoleMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateRoleMutation, CreateRoleMutationVariables>(CreateRoleDocument, options);
      }
export type CreateRoleMutationHookResult = ReturnType<typeof useCreateRoleMutation>;
export type CreateRoleMutationResult = Apollo.MutationResult<CreateRoleMutation>;
export type CreateRoleMutationOptions = Apollo.BaseMutationOptions<CreateRoleMutation, CreateRoleMutationVariables>;
export const UsersListDocument = gql`
    query UsersList($limit: Int, $offset: Int) {
  users(limit: $limit, offset: $offset) {
    totalCount
    nodes {
      id
      email
      name
      createdAt
      updatedAt
      roles {
        id
        name
        description
      }
    }
  }
}
    `;

/**
 * __useUsersListQuery__
 *
 * To run a query within a React component, call `useUsersListQuery` and pass it any options that fit your needs.
 * When your component renders, `useUsersListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useUsersListQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *      offset: // value for 'offset'
 *   },
 * });
 */
export function useUsersListQuery(baseOptions?: Apollo.QueryHookOptions<UsersListQuery, UsersListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<UsersListQuery, UsersListQueryVariables>(UsersListDocument, options);
      }
export function useUsersListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<UsersListQuery, UsersListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<UsersListQuery, UsersListQueryVariables>(UsersListDocument, options);
        }
export function useUsersListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<UsersListQuery, UsersListQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<UsersListQuery, UsersListQueryVariables>(UsersListDocument, options);
        }
export type UsersListQueryHookResult = ReturnType<typeof useUsersListQuery>;
export type UsersListLazyQueryHookResult = ReturnType<typeof useUsersListLazyQuery>;
export type UsersListSuspenseQueryHookResult = ReturnType<typeof useUsersListSuspenseQuery>;
export type UsersListQueryResult = Apollo.QueryResult<UsersListQuery, UsersListQueryVariables>;
export const UserDocument = gql`
    query User($id: ID!) {
  user(id: $id) {
    id
    email
    name
    teamId
    createdAt
    updatedAt
    roles {
      id
      name
      description
      permissions
    }
    permissions
  }
}
    `;

/**
 * __useUserQuery__
 *
 * To run a query within a React component, call `useUserQuery` and pass it any options that fit your needs.
 * When your component renders, `useUserQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useUserQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useUserQuery(baseOptions: Apollo.QueryHookOptions<UserQuery, UserQueryVariables> & ({ variables: UserQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<UserQuery, UserQueryVariables>(UserDocument, options);
      }
export function useUserLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<UserQuery, UserQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<UserQuery, UserQueryVariables>(UserDocument, options);
        }
export function useUserSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<UserQuery, UserQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<UserQuery, UserQueryVariables>(UserDocument, options);
        }
export type UserQueryHookResult = ReturnType<typeof useUserQuery>;
export type UserLazyQueryHookResult = ReturnType<typeof useUserLazyQuery>;
export type UserSuspenseQueryHookResult = ReturnType<typeof useUserSuspenseQuery>;
export type UserQueryResult = Apollo.QueryResult<UserQuery, UserQueryVariables>;
export const RolesDocument = gql`
    query Roles {
  roles {
    id
    name
    description
    permissions
    isSystem
    createdAt
  }
}
    `;

/**
 * __useRolesQuery__
 *
 * To run a query within a React component, call `useRolesQuery` and pass it any options that fit your needs.
 * When your component renders, `useRolesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useRolesQuery({
 *   variables: {
 *   },
 * });
 */
export function useRolesQuery(baseOptions?: Apollo.QueryHookOptions<RolesQuery, RolesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<RolesQuery, RolesQueryVariables>(RolesDocument, options);
      }
export function useRolesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<RolesQuery, RolesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<RolesQuery, RolesQueryVariables>(RolesDocument, options);
        }
export function useRolesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<RolesQuery, RolesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<RolesQuery, RolesQueryVariables>(RolesDocument, options);
        }
export type RolesQueryHookResult = ReturnType<typeof useRolesQuery>;
export type RolesLazyQueryHookResult = ReturnType<typeof useRolesLazyQuery>;
export type RolesSuspenseQueryHookResult = ReturnType<typeof useRolesSuspenseQuery>;
export type RolesQueryResult = Apollo.QueryResult<RolesQuery, RolesQueryVariables>;
export const PlansDocument = gql`
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

/**
 * __usePlansQuery__
 *
 * To run a query within a React component, call `usePlansQuery` and pass it any options that fit your needs.
 * When your component renders, `usePlansQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = usePlansQuery({
 *   variables: {
 *   },
 * });
 */
export function usePlansQuery(baseOptions?: Apollo.QueryHookOptions<PlansQuery, PlansQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<PlansQuery, PlansQueryVariables>(PlansDocument, options);
      }
export function usePlansLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<PlansQuery, PlansQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<PlansQuery, PlansQueryVariables>(PlansDocument, options);
        }
export function usePlansSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<PlansQuery, PlansQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<PlansQuery, PlansQueryVariables>(PlansDocument, options);
        }
export type PlansQueryHookResult = ReturnType<typeof usePlansQuery>;
export type PlansLazyQueryHookResult = ReturnType<typeof usePlansLazyQuery>;
export type PlansSuspenseQueryHookResult = ReturnType<typeof usePlansSuspenseQuery>;
export type PlansQueryResult = Apollo.QueryResult<PlansQuery, PlansQueryVariables>;
export const BillingPortalUrlDocument = gql`
    query BillingPortalUrl {
  billingPortalUrl
}
    `;

/**
 * __useBillingPortalUrlQuery__
 *
 * To run a query within a React component, call `useBillingPortalUrlQuery` and pass it any options that fit your needs.
 * When your component renders, `useBillingPortalUrlQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useBillingPortalUrlQuery({
 *   variables: {
 *   },
 * });
 */
export function useBillingPortalUrlQuery(baseOptions?: Apollo.QueryHookOptions<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>(BillingPortalUrlDocument, options);
      }
export function useBillingPortalUrlLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>(BillingPortalUrlDocument, options);
        }
export function useBillingPortalUrlSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>(BillingPortalUrlDocument, options);
        }
export type BillingPortalUrlQueryHookResult = ReturnType<typeof useBillingPortalUrlQuery>;
export type BillingPortalUrlLazyQueryHookResult = ReturnType<typeof useBillingPortalUrlLazyQuery>;
export type BillingPortalUrlSuspenseQueryHookResult = ReturnType<typeof useBillingPortalUrlSuspenseQuery>;
export type BillingPortalUrlQueryResult = Apollo.QueryResult<BillingPortalUrlQuery, BillingPortalUrlQueryVariables>;
export const CreateSubscriptionCheckoutDocument = gql`
    mutation CreateSubscriptionCheckout($planId: ID!) {
  createSubscriptionCheckout(planId: $planId) {
    sessionId
    url
  }
}
    `;
export type CreateSubscriptionCheckoutMutationFn = Apollo.MutationFunction<CreateSubscriptionCheckoutMutation, CreateSubscriptionCheckoutMutationVariables>;

/**
 * __useCreateSubscriptionCheckoutMutation__
 *
 * To run a mutation, you first call `useCreateSubscriptionCheckoutMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateSubscriptionCheckoutMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createSubscriptionCheckoutMutation, { data, loading, error }] = useCreateSubscriptionCheckoutMutation({
 *   variables: {
 *      planId: // value for 'planId'
 *   },
 * });
 */
export function useCreateSubscriptionCheckoutMutation(baseOptions?: Apollo.MutationHookOptions<CreateSubscriptionCheckoutMutation, CreateSubscriptionCheckoutMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateSubscriptionCheckoutMutation, CreateSubscriptionCheckoutMutationVariables>(CreateSubscriptionCheckoutDocument, options);
      }
export type CreateSubscriptionCheckoutMutationHookResult = ReturnType<typeof useCreateSubscriptionCheckoutMutation>;
export type CreateSubscriptionCheckoutMutationResult = Apollo.MutationResult<CreateSubscriptionCheckoutMutation>;
export type CreateSubscriptionCheckoutMutationOptions = Apollo.BaseMutationOptions<CreateSubscriptionCheckoutMutation, CreateSubscriptionCheckoutMutationVariables>;
export const UpdateSubscriptionDocument = gql`
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
export type UpdateSubscriptionMutationFn = Apollo.MutationFunction<UpdateSubscriptionMutation, UpdateSubscriptionMutationVariables>;

/**
 * __useUpdateSubscriptionMutation__
 *
 * To run a mutation, you first call `useUpdateSubscriptionMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSubscriptionMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSubscriptionMutation, { data, loading, error }] = useUpdateSubscriptionMutation({
 *   variables: {
 *      planId: // value for 'planId'
 *   },
 * });
 */
export function useUpdateSubscriptionMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSubscriptionMutation, UpdateSubscriptionMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSubscriptionMutation, UpdateSubscriptionMutationVariables>(UpdateSubscriptionDocument, options);
      }
export type UpdateSubscriptionMutationHookResult = ReturnType<typeof useUpdateSubscriptionMutation>;
export type UpdateSubscriptionMutationResult = Apollo.MutationResult<UpdateSubscriptionMutation>;
export type UpdateSubscriptionMutationOptions = Apollo.BaseMutationOptions<UpdateSubscriptionMutation, UpdateSubscriptionMutationVariables>;
export const CancelSubscriptionDocument = gql`
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
export type CancelSubscriptionMutationFn = Apollo.MutationFunction<CancelSubscriptionMutation, CancelSubscriptionMutationVariables>;

/**
 * __useCancelSubscriptionMutation__
 *
 * To run a mutation, you first call `useCancelSubscriptionMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCancelSubscriptionMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [cancelSubscriptionMutation, { data, loading, error }] = useCancelSubscriptionMutation({
 *   variables: {
 *   },
 * });
 */
export function useCancelSubscriptionMutation(baseOptions?: Apollo.MutationHookOptions<CancelSubscriptionMutation, CancelSubscriptionMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CancelSubscriptionMutation, CancelSubscriptionMutationVariables>(CancelSubscriptionDocument, options);
      }
export type CancelSubscriptionMutationHookResult = ReturnType<typeof useCancelSubscriptionMutation>;
export type CancelSubscriptionMutationResult = Apollo.MutationResult<CancelSubscriptionMutation>;
export type CancelSubscriptionMutationOptions = Apollo.BaseMutationOptions<CancelSubscriptionMutation, CancelSubscriptionMutationVariables>;
export const UpdateNotificationPreferencesDocument = gql`
    mutation UpdateNotificationPreferences($input: NotificationPreferencesInput!) {
  updateNotificationPreferences(input: $input) {
    emailEnabled
    pushEnabled
    inAppEnabled
    channels
  }
}
    `;
export type UpdateNotificationPreferencesMutationFn = Apollo.MutationFunction<UpdateNotificationPreferencesMutation, UpdateNotificationPreferencesMutationVariables>;

/**
 * __useUpdateNotificationPreferencesMutation__
 *
 * To run a mutation, you first call `useUpdateNotificationPreferencesMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNotificationPreferencesMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNotificationPreferencesMutation, { data, loading, error }] = useUpdateNotificationPreferencesMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateNotificationPreferencesMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNotificationPreferencesMutation, UpdateNotificationPreferencesMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNotificationPreferencesMutation, UpdateNotificationPreferencesMutationVariables>(UpdateNotificationPreferencesDocument, options);
      }
export type UpdateNotificationPreferencesMutationHookResult = ReturnType<typeof useUpdateNotificationPreferencesMutation>;
export type UpdateNotificationPreferencesMutationResult = Apollo.MutationResult<UpdateNotificationPreferencesMutation>;
export type UpdateNotificationPreferencesMutationOptions = Apollo.BaseMutationOptions<UpdateNotificationPreferencesMutation, UpdateNotificationPreferencesMutationVariables>;
export const MarkNotificationReadDocument = gql`
    mutation MarkNotificationRead($notificationId: ID!) {
  markNotificationRead(notificationId: $notificationId)
}
    `;
export type MarkNotificationReadMutationFn = Apollo.MutationFunction<MarkNotificationReadMutation, MarkNotificationReadMutationVariables>;

/**
 * __useMarkNotificationReadMutation__
 *
 * To run a mutation, you first call `useMarkNotificationReadMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMarkNotificationReadMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [markNotificationReadMutation, { data, loading, error }] = useMarkNotificationReadMutation({
 *   variables: {
 *      notificationId: // value for 'notificationId'
 *   },
 * });
 */
export function useMarkNotificationReadMutation(baseOptions?: Apollo.MutationHookOptions<MarkNotificationReadMutation, MarkNotificationReadMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MarkNotificationReadMutation, MarkNotificationReadMutationVariables>(MarkNotificationReadDocument, options);
      }
export type MarkNotificationReadMutationHookResult = ReturnType<typeof useMarkNotificationReadMutation>;
export type MarkNotificationReadMutationResult = Apollo.MutationResult<MarkNotificationReadMutation>;
export type MarkNotificationReadMutationOptions = Apollo.BaseMutationOptions<MarkNotificationReadMutation, MarkNotificationReadMutationVariables>;
export const SendNotificationDocument = gql`
    mutation SendNotification($input: SendNotificationInput!) {
  sendNotification(input: $input)
}
    `;
export type SendNotificationMutationFn = Apollo.MutationFunction<SendNotificationMutation, SendNotificationMutationVariables>;

/**
 * __useSendNotificationMutation__
 *
 * To run a mutation, you first call `useSendNotificationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSendNotificationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sendNotificationMutation, { data, loading, error }] = useSendNotificationMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useSendNotificationMutation(baseOptions?: Apollo.MutationHookOptions<SendNotificationMutation, SendNotificationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SendNotificationMutation, SendNotificationMutationVariables>(SendNotificationDocument, options);
      }
export type SendNotificationMutationHookResult = ReturnType<typeof useSendNotificationMutation>;
export type SendNotificationMutationResult = Apollo.MutationResult<SendNotificationMutation>;
export type SendNotificationMutationOptions = Apollo.BaseMutationOptions<SendNotificationMutation, SendNotificationMutationVariables>;
export const MyNotificationPreferencesDocument = gql`
    query MyNotificationPreferences {
  myNotificationPreferences {
    emailEnabled
    pushEnabled
    inAppEnabled
    channels
  }
}
    `;

/**
 * __useMyNotificationPreferencesQuery__
 *
 * To run a query within a React component, call `useMyNotificationPreferencesQuery` and pass it any options that fit your needs.
 * When your component renders, `useMyNotificationPreferencesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMyNotificationPreferencesQuery({
 *   variables: {
 *   },
 * });
 */
export function useMyNotificationPreferencesQuery(baseOptions?: Apollo.QueryHookOptions<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>(MyNotificationPreferencesDocument, options);
      }
export function useMyNotificationPreferencesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>(MyNotificationPreferencesDocument, options);
        }
export function useMyNotificationPreferencesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>(MyNotificationPreferencesDocument, options);
        }
export type MyNotificationPreferencesQueryHookResult = ReturnType<typeof useMyNotificationPreferencesQuery>;
export type MyNotificationPreferencesLazyQueryHookResult = ReturnType<typeof useMyNotificationPreferencesLazyQuery>;
export type MyNotificationPreferencesSuspenseQueryHookResult = ReturnType<typeof useMyNotificationPreferencesSuspenseQuery>;
export type MyNotificationPreferencesQueryResult = Apollo.QueryResult<MyNotificationPreferencesQuery, MyNotificationPreferencesQueryVariables>;
export const NotificationTokenDocument = gql`
    query NotificationToken {
  notificationToken {
    token
    socketUrl
    expiresAt
  }
}
    `;

/**
 * __useNotificationTokenQuery__
 *
 * To run a query within a React component, call `useNotificationTokenQuery` and pass it any options that fit your needs.
 * When your component renders, `useNotificationTokenQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNotificationTokenQuery({
 *   variables: {
 *   },
 * });
 */
export function useNotificationTokenQuery(baseOptions?: Apollo.QueryHookOptions<NotificationTokenQuery, NotificationTokenQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NotificationTokenQuery, NotificationTokenQueryVariables>(NotificationTokenDocument, options);
      }
export function useNotificationTokenLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NotificationTokenQuery, NotificationTokenQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NotificationTokenQuery, NotificationTokenQueryVariables>(NotificationTokenDocument, options);
        }
export function useNotificationTokenSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NotificationTokenQuery, NotificationTokenQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NotificationTokenQuery, NotificationTokenQueryVariables>(NotificationTokenDocument, options);
        }
export type NotificationTokenQueryHookResult = ReturnType<typeof useNotificationTokenQuery>;
export type NotificationTokenLazyQueryHookResult = ReturnType<typeof useNotificationTokenLazyQuery>;
export type NotificationTokenSuspenseQueryHookResult = ReturnType<typeof useNotificationTokenSuspenseQuery>;
export type NotificationTokenQueryResult = Apollo.QueryResult<NotificationTokenQuery, NotificationTokenQueryVariables>;