import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  // Use local schema file if available, otherwise connect to running gateway
  schema: process.env.GRAPHQL_SCHEMA_PATH || '../gateway/graphql-api-gateway/schema.graphqls',
  documents: ['src/**/*.{ts,tsx}', 'src/**/*.graphql'],
  generates: {
    'src/lib/graphql/generated/': {
      preset: 'client',
      plugins: [],
      presetConfig: {
        gqlTagName: 'gql',
      },
    },
    'src/lib/graphql/generated/hooks.tsx': {
      plugins: [
        'typescript',
        'typescript-operations',
        'typescript-react-apollo',
      ],
      config: {
        withHooks: true,
        withHOC: false,
        withComponent: false,
      },
    },
  },
};

export default config;
