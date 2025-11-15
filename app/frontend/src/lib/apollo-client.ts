import { ApolloClient, InMemoryCache, HttpLink, from, ApolloLink } from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { getToken, removeToken } from '@/contexts/AuthContext';

const httpLink = new HttpLink({
  uri: process.env.NEXT_PUBLIC_GRAPHQL_URL || 'http://localhost:4000/graphql',
  // Don't use credentials: 'include' with wildcard CORS
  // credentials: 'include',
});

// Auth middleware to add JWT token to requests
const authMiddleware = new ApolloLink((operation, forward) => {
  const token = getToken();
  
  operation.setContext(({ headers = {} }) => ({
    headers: {
      ...headers,
      ...(token ? { authorization: `Bearer ${token}` } : {}),
    },
  }));

  return forward(operation);
});

// Error handling middleware
const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path, extensions }) => {
      console.error(
        `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`,
        extensions
      );

      // Handle authentication errors
      if (extensions?.code === 'UNAUTHENTICATED') {
        // Clear token and redirect to login
        removeToken();
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
      }
    });
  }

  if (networkError) {
    console.error(`[Network error]: ${networkError}`);
  }
});

// Create Apollo Client
export const apolloClient = new ApolloClient({
  link: from([errorLink, authMiddleware, httpLink]),
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          me: {
            merge: true,
          },
        },
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
      errorPolicy: 'all',
    },
    query: {
      fetchPolicy: 'network-only',
      errorPolicy: 'all',
    },
    mutate: {
      errorPolicy: 'all',
    },
  },
});
