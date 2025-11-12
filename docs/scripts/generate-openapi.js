#!/usr/bin/env node

/**
 * OpenAPI Generator for Haunted SaaS
 * 
 * This script introspects the GraphQL schema and generates an OpenAPI 3.0 specification.
 * The generated openapi.yaml is consumed by Docusaurus for interactive API documentation.
 */

const fs = require('fs');
const path = require('path');
const axios = require('axios');

const GRAPHQL_ENDPOINT = process.env.GRAPHQL_URL || 'http://localhost:8080/graphql';
const OUTPUT_PATH = path.join(__dirname, '../static/openapi.yaml');

// GraphQL Introspection Query
const INTROSPECTION_QUERY = `
  query IntrospectionQuery {
    __schema {
      queryType { name }
      mutationType { name }
      types {
        kind
        name
        description
        fields {
          name
          description
          args {
            name
            description
            type {
              kind
              name
              ofType {
                kind
                name
              }
            }
          }
          type {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
              }
            }
          }
        }
        inputFields {
          name
          description
          type {
            kind
            name
            ofType {
              kind
              name
            }
          }
        }
      }
    }
  }
`;

async function introspectSchema() {
  console.log(`Introspecting GraphQL schema from ${GRAPHQL_ENDPOINT}...`);
  
  try {
    const response = await axios.post(GRAPHQL_ENDPOINT, {
      query: INTROSPECTION_QUERY,
    });

    if (response.data.errors) {
      throw new Error(`GraphQL errors: ${JSON.stringify(response.data.errors)}`);
    }

    return response.data.data.__schema;
  } catch (error) {
    console.error('Failed to introspect schema:', error.message);
    throw error;
  }
}

function convertGraphQLToOpenAPI(schema) {
  const openapi = {
    openapi: '3.0.0',
    info: {
      title: 'Haunted SaaS API',
      version: '1.0.0',
      description: 'Unified GraphQL API for Haunted SaaS platform',
      contact: {
        name: 'API Support',
        email: 'support@haunted-saas.com',
      },
    },
    servers: [
      {
        url: 'http://localhost:8080',
        description: 'Development server',
      },
      {
        url: 'https://api.haunted-saas.com',
        description: 'Production server',
      },
    ],
    paths: {
      '/graphql': {
        post: {
          summary: 'GraphQL Endpoint',
          description: 'Execute GraphQL queries and mutations',
          tags: ['GraphQL'],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: {
                  type: 'object',
                  required: ['query'],
                  properties: {
                    query: {
                      type: 'string',
                      description: 'GraphQL query or mutation',
                    },
                    variables: {
                      type: 'object',
                      description: 'Query variables',
                    },
                    operationName: {
                      type: 'string',
                      description: 'Operation name',
                    },
                  },
                },
              },
            },
          },
          responses: {
            '200': {
              description: 'Successful response',
              content: {
                'application/json': {
                  schema: {
                    type: 'object',
                    properties: {
                      data: {
                        type: 'object',
                        description: 'Query result data',
                      },
                      errors: {
                        type: 'array',
                        description: 'GraphQL errors',
                        items: {
                          type: 'object',
                        },
                      },
                    },
                  },
                },
              },
            },
          },
          security: [
            {
              bearerAuth: [],
            },
          ],
        },
      },
    },
    components: {
      securitySchemes: {
        bearerAuth: {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'JWT',
          description: 'JWT token from login/register mutation',
        },
      },
      schemas: {},
    },
    tags: [
      { name: 'Authentication', description: 'User authentication and authorization' },
      { name: 'Users', description: 'User management' },
      { name: 'Billing', description: 'Subscription and billing management' },
      { name: 'Analytics', description: 'Event tracking and analytics' },
      { name: 'Notifications', description: 'Real-time notifications' },
      { name: 'Feature Flags', description: 'Feature flag management' },
      { name: 'LLM', description: 'AI and LLM operations' },
      { name: 'GraphQL', description: 'GraphQL endpoint' },
    ],
  };

  // Add GraphQL types as schemas
  const types = schema.types.filter(
    (type) =>
      !type.name.startsWith('__') &&
      type.kind === 'OBJECT' &&
      type.name !== schema.queryType?.name &&
      type.name !== schema.mutationType?.name
  );

  types.forEach((type) => {
    const properties = {};
    
    if (type.fields) {
      type.fields.forEach((field) => {
        properties[field.name] = {
          type: getOpenAPIType(field.type),
          description: field.description || '',
        };
      });
    }

    openapi.components.schemas[type.name] = {
      type: 'object',
      description: type.description || '',
      properties,
    };
  });

  return openapi;
}

function getOpenAPIType(graphqlType) {
  if (!graphqlType) return 'string';
  
  const typeName = graphqlType.name || graphqlType.ofType?.name || 'String';
  
  const typeMap = {
    String: 'string',
    Int: 'integer',
    Float: 'number',
    Boolean: 'boolean',
    ID: 'string',
    JSON: 'object',
    Time: 'string',
  };

  return typeMap[typeName] || 'object';
}

function generateYAML(obj, indent = 0) {
  const spaces = '  '.repeat(indent);
  let yaml = '';

  for (const [key, value] of Object.entries(obj)) {
    if (value === null || value === undefined) continue;

    if (Array.isArray(value)) {
      yaml += `${spaces}${key}:\n`;
      value.forEach((item) => {
        if (typeof item === 'object') {
          yaml += `${spaces}- \n${generateYAML(item, indent + 1)}`;
        } else {
          yaml += `${spaces}- ${item}\n`;
        }
      });
    } else if (typeof value === 'object') {
      yaml += `${spaces}${key}:\n${generateYAML(value, indent + 1)}`;
    } else {
      const valueStr = typeof value === 'string' && value.includes('\n')
        ? `|\n${spaces}  ${value.split('\n').join(`\n${spaces}  `)}`
        : value;
      yaml += `${spaces}${key}: ${valueStr}\n`;
    }
  }

  return yaml;
}

async function main() {
  try {
    console.log('🚀 Starting OpenAPI generation...\n');

    // Introspect GraphQL schema
    const schema = await introspectSchema();
    console.log('✅ Schema introspected successfully\n');

    // Convert to OpenAPI
    const openapi = convertGraphQLToOpenAPI(schema);
    console.log('✅ Converted to OpenAPI 3.0 format\n');

    // Generate YAML
    const yaml = generateYAML(openapi);

    // Ensure output directory exists
    const outputDir = path.dirname(OUTPUT_PATH);
    if (!fs.existsSync(outputDir)) {
      fs.mkdirSync(outputDir, { recursive: true });
    }

    // Write to file
    fs.writeFileSync(OUTPUT_PATH, yaml, 'utf8');
    console.log(`✅ OpenAPI spec written to: ${OUTPUT_PATH}\n`);

    console.log('🎉 OpenAPI generation complete!');
    console.log('\nNext steps:');
    console.log('  1. Run: npm run build');
    console.log('  2. Run: npm start');
    console.log('  3. Visit: http://localhost:3000/docs/reference/api\n');
  } catch (error) {
    console.error('❌ Error generating OpenAPI spec:', error.message);
    process.exit(1);
  }
}

main();
