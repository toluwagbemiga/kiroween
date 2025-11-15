/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
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


export const CallPromptDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CallPrompt"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"variables"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"JSON"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"callPrompt"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}},{"kind":"Argument","name":{"kind":"Name","value":"variables"},"value":{"kind":"Variable","name":{"kind":"Name","value":"variables"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"content"}},{"kind":"Field","name":{"kind":"Name","value":"model"}},{"kind":"Field","name":{"kind":"Name","value":"tokensUsed"}},{"kind":"Field","name":{"kind":"Name","value":"cost"}}]}}]}}]} as unknown as DocumentNode<CallPromptMutation, CallPromptMutationVariables>;
export const LoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Login"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"LoginInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"token"}},{"kind":"Field","name":{"kind":"Name","value":"refreshToken"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"expiresAt"}}]}}]}}]} as unknown as DocumentNode<LoginMutation, LoginMutationVariables>;
export const RegisterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Register"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"RegisterInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"register"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"token"}},{"kind":"Field","name":{"kind":"Name","value":"refreshToken"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"expiresAt"}}]}}]}}]} as unknown as DocumentNode<RegisterMutation, RegisterMutationVariables>;
export const ValidateTokenDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ValidateToken"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]} as unknown as DocumentNode<ValidateTokenQuery, ValidateTokenQueryVariables>;
export const LogoutDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"logout"}}]}}]} as unknown as DocumentNode<LogoutMutation, LogoutMutationVariables>;
export const MeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"enabledFeatures"}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]} as unknown as DocumentNode<MeQuery, MeQueryVariables>;
export const TrackEventDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"TrackEvent"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"TrackEventInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"trackEvent"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}]}]}}]} as unknown as DocumentNode<TrackEventMutation, TrackEventMutationVariables>;
export const IdentifyUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"IdentifyUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"properties"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"JSON"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"identifyUser"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"properties"},"value":{"kind":"Variable","name":{"kind":"Name","value":"properties"}}}]}]}}]} as unknown as DocumentNode<IdentifyUserMutation, IdentifyUserMutationVariables>;