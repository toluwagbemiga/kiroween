import { gql } from '@apollo/client';

// Query to fetch users with pagination (detailed version for user management)
export const USERS_LIST_QUERY = gql`
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

// Query to fetch a single user by ID
export const USER_QUERY = gql`
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

// Query to fetch all available roles
export const ROLES_QUERY = gql`
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
