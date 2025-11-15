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

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;