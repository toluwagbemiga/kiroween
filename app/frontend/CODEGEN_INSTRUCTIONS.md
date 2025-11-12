# GraphQL Code Generation Instructions

## Overview

The frontend uses GraphQL Code Generator to automatically create TypeScript types and React hooks from the GraphQL schema.

## When to Regenerate

Regenerate the GraphQL types whenever:
- The GraphQL schema changes (new queries, mutations, or types)
- You add new GraphQL operations in the frontend
- You update existing GraphQL operations

## How to Regenerate

### Prerequisites

Make sure the GraphQL API Gateway is running and accessible at the URL specified in `codegen.ts`:
```typescript
schema: process.env.NEXT_PUBLIC_GRAPHQL_URL || 'http://localhost:8080/graphql'
```

### Run Code Generation

```bash
cd app/frontend
npm run codegen
```

This will:
1. Fetch the schema from the GraphQL endpoint
2. Generate TypeScript types for all GraphQL types
3. Generate React hooks for all queries and mutations
4. Create type-safe hooks in `src/generated/graphql.ts`

### Generated Files

The code generator creates:
- `src/generated/graphql.ts`: All TypeScript types and React hooks

### Example Generated Hooks

After running codegen, you'll have access to hooks like:

```typescript
import { 
  useTrackEventMutation,
  useIdentifyUserMutation,
  useLoginMutation,
  useMeQuery,
  // ... and many more
} from '@/generated/graphql';
```

## Configuration

The code generation is configured in `codegen.ts`:

```typescript
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: process.env.NEXT_PUBLIC_GRAPHQL_URL || 'http://localhost:8080/graphql',
  documents: ['src/**/*.{ts,tsx}'],
  generates: {
    './src/generated/': {
      preset: 'client',
      plugins: [],
      presetConfig: {
        gqlTagName: 'gql',
      },
    },
  },
  ignoreNoDocuments: true,
};

export default config;
```

## Troubleshooting

### Error: Cannot connect to GraphQL endpoint

**Solution:** Make sure the GraphQL API Gateway is running:
```bash
cd app/gateway/graphql-api-gateway
go run cmd/main.go
```

### Error: No documents found

**Solution:** Make sure you have GraphQL operations defined in your `.ts` or `.tsx` files using the `gql` tag:
```typescript
import { gql } from '@apollo/client';

const MY_QUERY = gql`
  query MyQuery {
    me {
      id
      email
    }
  }
`;
```

### Generated types are outdated

**Solution:** Delete the generated folder and regenerate:
```bash
rm -rf src/generated
npm run codegen
```

## Best Practices

1. **Run codegen after schema changes**: Always regenerate after updating the GraphQL schema
2. **Commit generated files**: Include generated files in version control for consistency
3. **Use generated hooks**: Prefer generated hooks over manual `useQuery`/`useMutation` for type safety
4. **Name operations**: Always name your GraphQL operations for better generated hook names

### Good Example:
```typescript
const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(input: { email: $email, password: $password }) {
      token
      user { id email }
    }
  }
`;
// Generates: useLoginMutation()
```

### Bad Example:
```typescript
const MUTATION = gql`
  mutation($email: String!) {
    login(input: { email: $email }) {
      token
    }
  }
`;
// Generates: useMutation() - not descriptive
```
