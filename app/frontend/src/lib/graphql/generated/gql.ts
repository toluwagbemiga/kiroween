/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n  mutation CallPrompt($name: String!, $variables: JSON!) {\n    callPrompt(name: $name, variables: $variables) {\n      content\n      model\n      tokensUsed\n      cost\n    }\n  }\n": typeof types.CallPromptDocument,
    "\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n": typeof types.LoginDocument,
    "\n  mutation Register($input: RegisterInput!) {\n    register(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n": typeof types.RegisterDocument,
    "\n  query ValidateToken {\n    me {\n      id\n      email\n      name\n      roles {\n        id\n        name\n      }\n    }\n  }\n": typeof types.ValidateTokenDocument,
    "\n  mutation Logout {\n    logout\n  }\n": typeof types.LogoutDocument,
    "\n  query Me {\n    me {\n      id\n      email\n      name\n      enabledFeatures\n      roles {\n        id\n        name\n      }\n    }\n  }\n": typeof types.MeDocument,
    "\n  mutation TrackEvent($input: TrackEventInput!) {\n    trackEvent(input: $input)\n  }\n": typeof types.TrackEventDocument,
    "\n  mutation IdentifyUser($properties: JSON!) {\n    identifyUser(properties: $properties)\n  }\n": typeof types.IdentifyUserDocument,
    "\n  mutation CreateSubscriptionCheckout($planId: ID!) {\n    createSubscriptionCheckout(planId: $planId) {\n      sessionId\n      url\n    }\n  }\n": typeof types.CreateSubscriptionCheckoutDocument,
    "\n  mutation UpdateSubscription($planId: ID!) {\n    updateSubscription(planId: $planId) {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n      }\n    }\n  }\n": typeof types.UpdateSubscriptionDocument,
    "\n  mutation CancelSubscription {\n    cancelSubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n      }\n    }\n  }\n": typeof types.CancelSubscriptionDocument,
    "\n  mutation UpdateNotificationPreferences($input: NotificationPreferencesInput!) {\n    updateNotificationPreferences(input: $input) {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n": typeof types.UpdateNotificationPreferencesDocument,
    "\n  mutation MarkNotificationRead($notificationId: ID!) {\n    markNotificationRead(notificationId: $notificationId)\n  }\n": typeof types.MarkNotificationReadDocument,
    "\n  mutation SendNotification($input: SendNotificationInput!) {\n    sendNotification(input: $input)\n  }\n": typeof types.SendNotificationDocument,
    "\n  mutation UpdateProfile($input: UpdateProfileInput!) {\n    updateProfile(input: $input) {\n      id\n      email\n      name\n      updatedAt\n    }\n  }\n": typeof types.UpdateProfileDocument,
    "\n  mutation AssignRole($userId: ID!, $roleId: ID!) {\n    assignRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n": typeof types.AssignRoleDocument,
    "\n  mutation RemoveRole($userId: ID!, $roleId: ID!) {\n    removeRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n": typeof types.RemoveRoleDocument,
    "\n  mutation CreateRole($input: CreateRoleInput!) {\n    createRole(input: $input) {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n": typeof types.CreateRoleDocument,
    "\n  query MyAnalytics($startDate: Time, $endDate: Time) {\n    myAnalytics(startDate: $startDate, endDate: $endDate) {\n      totalEvents\n      uniqueUsers\n      eventsByType\n      topEvents {\n        eventName\n        count\n      }\n    }\n  }\n": typeof types.MyAnalyticsDocument,
    "\n  query Plans {\n    plans {\n      id\n      name\n      description\n      price\n      currency\n      interval\n      features\n      stripePriceId\n      isActive\n    }\n  }\n": typeof types.PlansDocument,
    "\n  query MySubscription {\n    mySubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      stripeSubscriptionId\n      createdAt\n      updatedAt\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n        stripePriceId\n        isActive\n      }\n    }\n  }\n": typeof types.MySubscriptionDocument,
    "\n  query BillingPortalUrl {\n    billingPortalUrl\n  }\n": typeof types.BillingPortalUrlDocument,
    "\n  query Users($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        roles {\n          id\n          name\n        }\n      }\n    }\n  }\n": typeof types.UsersDocument,
    "\n  query MyNotificationPreferences {\n    myNotificationPreferences {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n": typeof types.MyNotificationPreferencesDocument,
    "\n  query NotificationToken {\n    notificationToken {\n      token\n      socketUrl\n      expiresAt\n    }\n  }\n": typeof types.NotificationTokenDocument,
    "\n  query UsersList($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        updatedAt\n        roles {\n          id\n          name\n          description\n        }\n      }\n    }\n  }\n": typeof types.UsersListDocument,
    "\n  query User($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      teamId\n      createdAt\n      updatedAt\n      roles {\n        id\n        name\n        description\n        permissions\n      }\n      permissions\n    }\n  }\n": typeof types.UserDocument,
    "\n  query Roles {\n    roles {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n": typeof types.RolesDocument,
};
const documents: Documents = {
    "\n  mutation CallPrompt($name: String!, $variables: JSON!) {\n    callPrompt(name: $name, variables: $variables) {\n      content\n      model\n      tokensUsed\n      cost\n    }\n  }\n": types.CallPromptDocument,
    "\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n": types.LoginDocument,
    "\n  mutation Register($input: RegisterInput!) {\n    register(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n": types.RegisterDocument,
    "\n  query ValidateToken {\n    me {\n      id\n      email\n      name\n      roles {\n        id\n        name\n      }\n    }\n  }\n": types.ValidateTokenDocument,
    "\n  mutation Logout {\n    logout\n  }\n": types.LogoutDocument,
    "\n  query Me {\n    me {\n      id\n      email\n      name\n      enabledFeatures\n      roles {\n        id\n        name\n      }\n    }\n  }\n": types.MeDocument,
    "\n  mutation TrackEvent($input: TrackEventInput!) {\n    trackEvent(input: $input)\n  }\n": types.TrackEventDocument,
    "\n  mutation IdentifyUser($properties: JSON!) {\n    identifyUser(properties: $properties)\n  }\n": types.IdentifyUserDocument,
    "\n  mutation CreateSubscriptionCheckout($planId: ID!) {\n    createSubscriptionCheckout(planId: $planId) {\n      sessionId\n      url\n    }\n  }\n": types.CreateSubscriptionCheckoutDocument,
    "\n  mutation UpdateSubscription($planId: ID!) {\n    updateSubscription(planId: $planId) {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n      }\n    }\n  }\n": types.UpdateSubscriptionDocument,
    "\n  mutation CancelSubscription {\n    cancelSubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n      }\n    }\n  }\n": types.CancelSubscriptionDocument,
    "\n  mutation UpdateNotificationPreferences($input: NotificationPreferencesInput!) {\n    updateNotificationPreferences(input: $input) {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n": types.UpdateNotificationPreferencesDocument,
    "\n  mutation MarkNotificationRead($notificationId: ID!) {\n    markNotificationRead(notificationId: $notificationId)\n  }\n": types.MarkNotificationReadDocument,
    "\n  mutation SendNotification($input: SendNotificationInput!) {\n    sendNotification(input: $input)\n  }\n": types.SendNotificationDocument,
    "\n  mutation UpdateProfile($input: UpdateProfileInput!) {\n    updateProfile(input: $input) {\n      id\n      email\n      name\n      updatedAt\n    }\n  }\n": types.UpdateProfileDocument,
    "\n  mutation AssignRole($userId: ID!, $roleId: ID!) {\n    assignRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n": types.AssignRoleDocument,
    "\n  mutation RemoveRole($userId: ID!, $roleId: ID!) {\n    removeRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n": types.RemoveRoleDocument,
    "\n  mutation CreateRole($input: CreateRoleInput!) {\n    createRole(input: $input) {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n": types.CreateRoleDocument,
    "\n  query MyAnalytics($startDate: Time, $endDate: Time) {\n    myAnalytics(startDate: $startDate, endDate: $endDate) {\n      totalEvents\n      uniqueUsers\n      eventsByType\n      topEvents {\n        eventName\n        count\n      }\n    }\n  }\n": types.MyAnalyticsDocument,
    "\n  query Plans {\n    plans {\n      id\n      name\n      description\n      price\n      currency\n      interval\n      features\n      stripePriceId\n      isActive\n    }\n  }\n": types.PlansDocument,
    "\n  query MySubscription {\n    mySubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      stripeSubscriptionId\n      createdAt\n      updatedAt\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n        stripePriceId\n        isActive\n      }\n    }\n  }\n": types.MySubscriptionDocument,
    "\n  query BillingPortalUrl {\n    billingPortalUrl\n  }\n": types.BillingPortalUrlDocument,
    "\n  query Users($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        roles {\n          id\n          name\n        }\n      }\n    }\n  }\n": types.UsersDocument,
    "\n  query MyNotificationPreferences {\n    myNotificationPreferences {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n": types.MyNotificationPreferencesDocument,
    "\n  query NotificationToken {\n    notificationToken {\n      token\n      socketUrl\n      expiresAt\n    }\n  }\n": types.NotificationTokenDocument,
    "\n  query UsersList($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        updatedAt\n        roles {\n          id\n          name\n          description\n        }\n      }\n    }\n  }\n": types.UsersListDocument,
    "\n  query User($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      teamId\n      createdAt\n      updatedAt\n      roles {\n        id\n        name\n        description\n        permissions\n      }\n      permissions\n    }\n  }\n": types.UserDocument,
    "\n  query Roles {\n    roles {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n": types.RolesDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CallPrompt($name: String!, $variables: JSON!) {\n    callPrompt(name: $name, variables: $variables) {\n      content\n      model\n      tokensUsed\n      cost\n    }\n  }\n"): (typeof documents)["\n  mutation CallPrompt($name: String!, $variables: JSON!) {\n    callPrompt(name: $name, variables: $variables) {\n      content\n      model\n      tokensUsed\n      cost\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n"): (typeof documents)["\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation Register($input: RegisterInput!) {\n    register(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n"): (typeof documents)["\n  mutation Register($input: RegisterInput!) {\n    register(input: $input) {\n      token\n      refreshToken\n      user {\n        id\n        email\n        name\n        roles {\n          id\n          name\n        }\n      }\n      expiresAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query ValidateToken {\n    me {\n      id\n      email\n      name\n      roles {\n        id\n        name\n      }\n    }\n  }\n"): (typeof documents)["\n  query ValidateToken {\n    me {\n      id\n      email\n      name\n      roles {\n        id\n        name\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation Logout {\n    logout\n  }\n"): (typeof documents)["\n  mutation Logout {\n    logout\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query Me {\n    me {\n      id\n      email\n      name\n      enabledFeatures\n      roles {\n        id\n        name\n      }\n    }\n  }\n"): (typeof documents)["\n  query Me {\n    me {\n      id\n      email\n      name\n      enabledFeatures\n      roles {\n        id\n        name\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation TrackEvent($input: TrackEventInput!) {\n    trackEvent(input: $input)\n  }\n"): (typeof documents)["\n  mutation TrackEvent($input: TrackEventInput!) {\n    trackEvent(input: $input)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation IdentifyUser($properties: JSON!) {\n    identifyUser(properties: $properties)\n  }\n"): (typeof documents)["\n  mutation IdentifyUser($properties: JSON!) {\n    identifyUser(properties: $properties)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CreateSubscriptionCheckout($planId: ID!) {\n    createSubscriptionCheckout(planId: $planId) {\n      sessionId\n      url\n    }\n  }\n"): (typeof documents)["\n  mutation CreateSubscriptionCheckout($planId: ID!) {\n    createSubscriptionCheckout(planId: $planId) {\n      sessionId\n      url\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation UpdateSubscription($planId: ID!) {\n    updateSubscription(planId: $planId) {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n      }\n    }\n  }\n"): (typeof documents)["\n  mutation UpdateSubscription($planId: ID!) {\n    updateSubscription(planId: $planId) {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CancelSubscription {\n    cancelSubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n      }\n    }\n  }\n"): (typeof documents)["\n  mutation CancelSubscription {\n    cancelSubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      plan {\n        id\n        name\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation UpdateNotificationPreferences($input: NotificationPreferencesInput!) {\n    updateNotificationPreferences(input: $input) {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n"): (typeof documents)["\n  mutation UpdateNotificationPreferences($input: NotificationPreferencesInput!) {\n    updateNotificationPreferences(input: $input) {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation MarkNotificationRead($notificationId: ID!) {\n    markNotificationRead(notificationId: $notificationId)\n  }\n"): (typeof documents)["\n  mutation MarkNotificationRead($notificationId: ID!) {\n    markNotificationRead(notificationId: $notificationId)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation SendNotification($input: SendNotificationInput!) {\n    sendNotification(input: $input)\n  }\n"): (typeof documents)["\n  mutation SendNotification($input: SendNotificationInput!) {\n    sendNotification(input: $input)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation UpdateProfile($input: UpdateProfileInput!) {\n    updateProfile(input: $input) {\n      id\n      email\n      name\n      updatedAt\n    }\n  }\n"): (typeof documents)["\n  mutation UpdateProfile($input: UpdateProfileInput!) {\n    updateProfile(input: $input) {\n      id\n      email\n      name\n      updatedAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation AssignRole($userId: ID!, $roleId: ID!) {\n    assignRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n"): (typeof documents)["\n  mutation AssignRole($userId: ID!, $roleId: ID!) {\n    assignRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation RemoveRole($userId: ID!, $roleId: ID!) {\n    removeRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n"): (typeof documents)["\n  mutation RemoveRole($userId: ID!, $roleId: ID!) {\n    removeRole(userId: $userId, roleId: $roleId) {\n      id\n      email\n      name\n      roles {\n        id\n        name\n        description\n      }\n      permissions\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CreateRole($input: CreateRoleInput!) {\n    createRole(input: $input) {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n"): (typeof documents)["\n  mutation CreateRole($input: CreateRoleInput!) {\n    createRole(input: $input) {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query MyAnalytics($startDate: Time, $endDate: Time) {\n    myAnalytics(startDate: $startDate, endDate: $endDate) {\n      totalEvents\n      uniqueUsers\n      eventsByType\n      topEvents {\n        eventName\n        count\n      }\n    }\n  }\n"): (typeof documents)["\n  query MyAnalytics($startDate: Time, $endDate: Time) {\n    myAnalytics(startDate: $startDate, endDate: $endDate) {\n      totalEvents\n      uniqueUsers\n      eventsByType\n      topEvents {\n        eventName\n        count\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query Plans {\n    plans {\n      id\n      name\n      description\n      price\n      currency\n      interval\n      features\n      stripePriceId\n      isActive\n    }\n  }\n"): (typeof documents)["\n  query Plans {\n    plans {\n      id\n      name\n      description\n      price\n      currency\n      interval\n      features\n      stripePriceId\n      isActive\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query MySubscription {\n    mySubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      stripeSubscriptionId\n      createdAt\n      updatedAt\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n        stripePriceId\n        isActive\n      }\n    }\n  }\n"): (typeof documents)["\n  query MySubscription {\n    mySubscription {\n      id\n      userId\n      planId\n      status\n      currentPeriodStart\n      currentPeriodEnd\n      cancelAtPeriodEnd\n      stripeSubscriptionId\n      createdAt\n      updatedAt\n      plan {\n        id\n        name\n        description\n        price\n        currency\n        interval\n        features\n        stripePriceId\n        isActive\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query BillingPortalUrl {\n    billingPortalUrl\n  }\n"): (typeof documents)["\n  query BillingPortalUrl {\n    billingPortalUrl\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query Users($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        roles {\n          id\n          name\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query Users($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        roles {\n          id\n          name\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query MyNotificationPreferences {\n    myNotificationPreferences {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n"): (typeof documents)["\n  query MyNotificationPreferences {\n    myNotificationPreferences {\n      emailEnabled\n      pushEnabled\n      inAppEnabled\n      channels\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query NotificationToken {\n    notificationToken {\n      token\n      socketUrl\n      expiresAt\n    }\n  }\n"): (typeof documents)["\n  query NotificationToken {\n    notificationToken {\n      token\n      socketUrl\n      expiresAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query UsersList($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        updatedAt\n        roles {\n          id\n          name\n          description\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query UsersList($limit: Int, $offset: Int) {\n    users(limit: $limit, offset: $offset) {\n      totalCount\n      nodes {\n        id\n        email\n        name\n        createdAt\n        updatedAt\n        roles {\n          id\n          name\n          description\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query User($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      teamId\n      createdAt\n      updatedAt\n      roles {\n        id\n        name\n        description\n        permissions\n      }\n      permissions\n    }\n  }\n"): (typeof documents)["\n  query User($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      teamId\n      createdAt\n      updatedAt\n      roles {\n        id\n        name\n        description\n        permissions\n      }\n      permissions\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query Roles {\n    roles {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n"): (typeof documents)["\n  query Roles {\n    roles {\n      id\n      name\n      description\n      permissions\n      isSystem\n      createdAt\n    }\n  }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;