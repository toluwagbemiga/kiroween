# HAUNTED SAAS SKELETON Frontend

Next.js 14+ application with TypeScript, Design System, and GraphQL integration.

## Features

- ✅ Next.js 14 with App Router
- ✅ TypeScript for type safety
- ✅ GraphQL client (Apollo or urql)
- ✅ Socket.IO for real-time notifications
- ✅ Design System components
- ✅ Bento Grid dashboard layout
- ✅ Glassmorphism aesthetic
- ✅ WCAG 2.1 Level AA accessibility
- ✅ Responsive design

## Tech Stack

- **Framework**: Next.js 14+
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **GraphQL**: Apollo Client or urql
- **Real-time**: Socket.IO client
- **State**: React Context + hooks
- **Forms**: React Hook Form + Zod validation
- **Testing**: Jest + React Testing Library

## Project Structure

```
app/
  ├── (auth)/
  │   ├── login/
  │   └── register/
  ├── (dashboard)/
  │   ├── layout.tsx
  │   ├── page.tsx              # Bento Grid dashboard
  │   ├── settings/
  │   └── billing/
  └── layout.tsx
components/
  ├── ui/                       # Design System components
  │   ├── Button.tsx
  │   ├── Card.tsx
  │   ├── Input.tsx
  │   └── ...
  ├── dashboard/
  │   ├── BentoGrid.tsx
  │   └── StatsCard.tsx
  └── notifications/
      └── NotificationToast.tsx
lib/
  ├── graphql/
  │   ├── client.ts
  │   └── queries.ts
  ├── socket.ts
  └── auth.ts
```

## Design System

Reusable components following:
- Consistent spacing and typography
- Accessible color contrast
- Keyboard navigation support
- Screen reader compatibility
- Focus indicators

## Bento Grid Dashboard

Modern dashboard layout with:
- Responsive grid system
- Glassmorphism cards
- Real-time data updates
- Interactive charts
- Quick actions

## Environment Variables

```bash
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/graphql
NEXT_PUBLIC_SOCKETIO_URL=http://localhost:3002
NEXT_PUBLIC_ANALYTICS_ENABLED=true
NEXT_PUBLIC_FEATURE_FLAGS_ENABLED=true
```

## GraphQL Code Generation

The project uses GraphQL Code Generator to create TypeScript types from the GraphQL schema.

**Run codegen:**
```bash
cd app/frontend
npm run codegen
```

This reads the schema from `../gateway/graphql-api-gateway/schema.graphqls` without needing the gateway to be running.

**Generated files:**
- `src/lib/graphql/generated/graphql.ts` - TypeScript types
- `src/lib/graphql/generated/hooks.tsx` - React hooks for queries/mutations
- `src/lib/graphql/generated/gql.ts` - GraphQL tag function
- `src/lib/graphql/generated/index.ts` - Main exports

**Important:** All GraphQL mutations must use the `input` object format:
```typescript
// ✅ Correct
mutation Login($input: LoginInput!) {
  login(input: $input) { ... }
}

// ❌ Wrong
mutation Login($email: String!, $password: String!) {
  login(email: $email, password: $password) { ... }
}
```

## Implementation Steps

1. Initialize Next.js 14 with TypeScript
2. Set up Tailwind CSS
3. Create Design System components
4. Implement GraphQL client
5. Add Socket.IO integration
6. Build authentication pages
7. Create Bento Grid dashboard
8. Add billing pages
9. Implement settings
10. Add accessibility features
11. Write component tests

## Authentication

```typescript
// Login flow
const { data } = await login({
  variables: { email, password }
});

// Store JWT
localStorage.setItem('token', data.login.accessToken);

// Use in GraphQL client
const client = new ApolloClient({
  headers: {
    authorization: `Bearer ${token}`
  }
});
```

## Real-time Notifications

```typescript
import io from 'socket.io-client';

const socket = io(process.env.NEXT_PUBLIC_SOCKETIO_URL, {
  auth: { token: getToken() }
});

socket.on('notification', (data) => {
  showToast(data);
});
```

## Accessibility

- Semantic HTML
- ARIA labels and roles
- Keyboard navigation
- Focus management
- Color contrast (WCAG AA)
- Screen reader testing

## Next Steps

1. Set up Next.js project
2. Configure Tailwind CSS
3. Build Design System
4. Implement auth pages
5. Create dashboard layout
6. Add GraphQL integration
7. Implement Socket.IO
8. Test accessibility
