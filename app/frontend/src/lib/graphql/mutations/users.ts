import { gql } from '@apollo/client';

// Mutation to update user profile
export const UPDATE_USER_MUTATION = gql`
  mutation UpdateProfile($input: UpdateProfileInput!) {
    updateProfile(input: $input) {
      id
      email
      name
      updatedAt
    }
  }
`;

// Mutation to assign a role to a user
export const ASSIGN_ROLE_MUTATION = gql`
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

// Mutation to remove a role from a user
export const REMOVE_ROLE_MUTATION = gql`
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

// Mutation to create a new role
export const CREATE_ROLE_MUTATION = gql`
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
