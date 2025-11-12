# HAUNTED SAAS SKELETON Documentation

Docusaurus site with auto-generated API reference and comprehensive guides.

## Features

- ✅ Docusaurus 3.0
- ✅ MDX support
- ✅ Auto-generated API reference from OpenAPI
- ✅ Diátaxis framework (tutorials, how-to, reference, explanation)
- ✅ Code examples
- ✅ Search functionality
- ✅ Dark mode

## Structure

```
docs/
  ├── tutorial/
  │   ├── getting-started.md
  │   ├── authentication.md
  │   └── first-feature.md
  ├── how-to/
  │   ├── manage-users.md
  │   ├── configure-billing.md
  │   └── deploy-production.md
  ├── reference/
  │   ├── api/                  # Auto-generated from OpenAPI
  │   ├── architecture.md
  │   └── configuration.md
  └── explanation/
      ├── design-decisions.md
      └── security-model.md
```

## Diátaxis Framework

- **Tutorials**: Learning-oriented, step-by-step guides
- **How-to Guides**: Problem-oriented, practical steps
- **Reference**: Information-oriented, technical descriptions
- **Explanation**: Understanding-oriented, background and context

## OpenAPI Integration

Generate OpenAPI specs from gRPC:
```bash
# Use grpc-gateway or protoc-gen-openapi
protoc --openapi_out=docs/static/openapi \
  proto/**/*.proto
```

Then use Docusaurus OpenAPI plugin to render interactive API docs.

## Implementation Steps

1. Initialize Docusaurus project
2. Configure theme and plugins
3. Write tutorial content
4. Write how-to guides
5. Generate OpenAPI specs from proto files
6. Configure OpenAPI plugin
7. Write explanation content
8. Add code examples
9. Configure search
10. Deploy to Vercel

## Environment Variables

```bash
# For build
ALGOLIA_APP_ID=...
ALGOLIA_API_KEY=...
```

## Deployment

Deployed automatically via GitHub Actions to Vercel on push to main.

## Next Steps

1. Set up Docusaurus
2. Write getting started tutorial
3. Generate OpenAPI specs
4. Configure API reference
5. Write how-to guides
6. Add search
7. Deploy


## OpenAPI Generation

The `scripts/generate-openapi.js` script automatically generates API documentation:

1. Connects to the GraphQL API Gateway
2. Runs introspection query
3. Converts GraphQL schema to OpenAPI 3.0
4. Generates `static/openapi.yaml`

### Usage

```bash
# Make sure GraphQL API Gateway is running first
npm run generate-api
```

## Deployment Options

### Vercel (Recommended)
- Build Command: `npm install && npm run build`
- Output Directory: `build`
- Automatic deployments on push

### GitHub Pages
```bash
npm run deploy
```

### Docker
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY . .
RUN npm install && npm run build
FROM nginx:alpine
COPY --from=0 /app/build /usr/share/nginx/html
```

## Maintenance

- Run `npm run generate-api` after schema changes
- Run `npm run build` to verify before deploying
- Clear cache with `npm run clear` if issues occur
