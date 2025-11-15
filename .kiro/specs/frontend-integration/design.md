# Frontend Integration Design Document

## Overview

This design addresses the critical issues preventing the Haunted SaaS Skeleton frontend from functioning properly: missing Tailwind CSS styling, incomplete service integrations, and build configuration problems. The solution involves fixing the build process to properly compile Tailwind CSS, implementing complete GraphQL integrations for all backend services, creating missing UI pages, and switching from static export to a proper Next.js server deployment to support dynamic data fetching and real-time features.

## Architecture

### Current State Analysis

**Problems Identified:**
1. **Static Export Mode**: Next.js is configured with `output: 'export'` which creates a static site that cannot:
   - Make server-side API calls during build
   - Support dynamic routes properly
   - Handle real-time Socket.IO connections effectively
   - Use Next.js API routes

2. **Missing PostCSS Config**: No `postcss.config.js` file exists, which may cause Tailwind CSS compilation issues

3. **Incomplete Service Integration**: Dashboard shows placeholder content because GraphQL queries aren't implemented

4. **Missing Pages**: No pages exist for users, billing, analytics, or notifications management

### Proposed Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Next.js Application                      │
│                    (Server-Side Rendering)                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Pages      │  │  Components  │  │   Contexts   │      │
│  │              │  │              │  │              │      │
│  │ - Dashboard  │  │ - UI System  │  │ - Auth       │      │
│  │ - Users      │  │ - Layout     │  │ - Features   │      │
│  │ - Billing    │  │ - Forms      │  │ - Socket.IO  │      │
│  │ - Analytics  │  │ - Charts     │  │ - Analytics  │      │
│  │ - Settings   │  │ - Tables     │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           GraphQL Client (Apollo)                     │   │
│  │  - Generated Types & Hooks                            │   │
│  │  - Auth Middleware                                    │   │
│  │  - Error Handling                                     │   │
│  │  - Cache Management                                   │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                               │
└───────────────────────────┬───────────────────────────────────┘
                            │
                            │ HTTP/GraphQL
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              GraphQL API Gateway (Port 4000)                 │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│ User Auth    │   │   Billing    │   │  Analytics   │
│ Service      │   │   Service    │   │  Service     │
│ :50051       │   │   :50052     │   │  :50055      │
└──────────────┘   └──────────────┘   └──────────────┘
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│ Notifications│   │ Feature Flags│   │ LLM Gateway  │
│ Service      │   │ Service      │   │ Service      │
│ :50054       │   │ :50056       │   │ :50053       │
└──────────────┘   └──────────────┘   └──────────────┘
```

## Components and Interfaces

### 1. Build Configuration

#### PostCSS Configuration
Create `postcss.config.js` to ensure Tailwind CSS is properly processed:

```javascript
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

#### Next.js Configuration Updates
Modify `next.config.js` to:
- Remove `output: 'export'` to enable server-side rendering
- Keep image optimization disabled for Docker compatibility
- Ensure environment variables are properly passed

#### Dockerfile Updates
Update the Dockerfile to:
- Use Next.js standalone build for production
- Run `next start` instead of static file serving
- Properly handle environment variables at runtime

### 2. GraphQL Integration Layer

#### Generated Types and Hooks
Use GraphQL Code Generator to create:
- TypeScript types for all schema types
- React hooks for queries and mutations
- Automatic type inference for responses

**Key Queries to Implement:**
- `useMyAnalyticsQuery` - Dashboard statistics
- `useUsersQuery` - User management list
- `useMySubscriptionQuery` - Billing information
- `useIsFeatureEnabledQuery` - Feature flag checks
- `useMyNotificationPreferencesQuery` - Notification settings

**Key Mutations to Implement:**
- `useLoginMutation` - User authentication
- `useRegisterMutation` - User registration
- `useTrackEventMutation` - Analytics tracking
- `useCreateSubscriptionCheckoutMutation` - Billing checkout
- `useSendNotificationMutation` - Send notifications

#### Apollo Client Configuration
Enhance the existing Apollo Client setup with:
- Optimistic UI updates for mutations
- Proper error boundaries
- Retry logic for failed requests
- Cache invalidation strategies

### 3. Page Implementations

#### Dashboard Page (`/dashboard`)
**Current State**: Shows placeholder with feature flag fallback

**Enhanced Implementation**:
- Fetch real analytics data using `useMyAnalyticsQuery`
- Display subscription status from `useMySubscriptionQuery`
- Show recent activity from notifications service
- Real-time updates via Socket.IO integration
- Charts using a lightweight charting library (Chart.js or Recharts)

**Data Flow**:
```
Dashboard Component
  ├─> useMyAnalyticsQuery() → Analytics Service
  ├─> useMySubscriptionQuery() → Billing Service
  ├─> useSocketIO() → Notifications Service (real-time)
  └─> useFeatureFlags() → Feature Flags Service
```

#### Users Page (`/users`) - NEW
**Purpose**: User management interface for administrators

**Features**:
- User list table with pagination
- Search and filter capabilities
- Role assignment interface
- User creation/edit forms
- Permission management

**Components**:
- `UserTable` - Data table with sorting
- `UserForm` - Create/edit user modal
- `RoleSelector` - Multi-select for roles
- `PermissionMatrix` - Visual permission display

#### Billing Page (`/billing`) - NEW
**Purpose**: Subscription and payment management

**Features**:
- Current plan display
- Plan comparison table
- Upgrade/downgrade flows
- Billing history
- Invoice downloads
- Stripe checkout integration

**Components**:
- `PlanCard` - Subscription plan display
- `BillingHistory` - Invoice table
- `CheckoutButton` - Stripe integration
- `PaymentMethodForm` - Payment details

#### Analytics Page (`/analytics`) - NEW
**Purpose**: Detailed analytics dashboard

**Features**:
- Event tracking visualization
- User behavior analytics
- Custom date range selection
- Export capabilities
- Real-time event stream

**Components**:
- `AnalyticsChart` - Time-series visualization
- `EventTable` - Event log display
- `DateRangePicker` - Date selection
- `MetricCard` - KPI display

#### Notifications Page (`/notifications`) - NEW
**Purpose**: Notification center and preferences

**Features**:
- Notification list with read/unread status
- Mark as read functionality
- Notification preferences
- Channel configuration (email, push, in-app)

**Components**:
- `NotificationList` - Scrollable notification feed
- `NotificationItem` - Individual notification display
- `PreferencesForm` - Settings interface

### 4. UI Component Enhancements

#### Existing Components to Enhance
All existing UI components (Button, Card, Input, etc.) are already styled with Tailwind CSS. These will work once the build process is fixed.

#### New Components Needed

**DataTable Component**
- Sortable columns
- Pagination
- Row selection
- Loading states
- Empty states

**Chart Components**
- LineChart - Time series data
- BarChart - Comparative data
- PieChart - Distribution data
- StatCard - Single metric display

**Form Components**
- FormField - Labeled input wrapper
- FormError - Error message display
- FormSection - Grouped form fields
- SubmitButton - Loading state button

### 5. Real-time Integration

#### Socket.IO Client Setup
Enhance the existing `SocketProvider` to:
- Auto-reconnect on disconnect
- Handle authentication with JWT
- Subscribe to user-specific channels
- Emit and receive typed events

**Event Types**:
```typescript
interface SocketEvents {
  // Incoming
  'notification:new': (notification: Notification) => void;
  'user:updated': (user: User) => void;
  'subscription:changed': (subscription: Subscription) => void;
  
  // Outgoing
  'subscribe': (channels: string[]) => void;
  'unsubscribe': (channels: string[]) => void;
}
```

#### Real-time Updates
Integrate Socket.IO events with Apollo Cache:
- Update cache when notifications arrive
- Refresh queries on relevant events
- Show toast notifications for important events

### 6. Feature Flag Integration

#### Enhanced Feature Component
The existing `Feature` component works well. Extend it with:
- Loading states while fetching flags
- Error boundaries for flag failures
- Analytics tracking for feature usage

#### Feature Flag Usage Patterns
```typescript
// Conditional rendering
<Feature name="new-dashboard" fallback={<OldDashboard />}>
  <NewDashboard />
</Feature>

// Programmatic checks
const { isEnabled } = useFeatureFlag('advanced-analytics');
if (isEnabled) {
  // Show advanced features
}
```

### 7. Analytics Integration

#### Event Tracking Strategy
Enhance the existing `analytics.ts` utility:
- Auto-track page views
- Track user interactions (clicks, form submissions)
- Track errors and performance metrics
- Batch events for efficiency

**Event Categories**:
- `page_view` - Navigation tracking
- `user_action` - Button clicks, form submissions
- `feature_usage` - Feature flag interactions
- `error` - Error tracking
- `performance` - Load times, API latency

## Data Models

### Frontend State Management

#### Auth State
```typescript
interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: Error | null;
}
```

#### Feature Flag State
```typescript
interface FeatureFlagState {
  flags: Record<string, boolean>;
  variants: Record<string, FeatureVariant>;
  isLoading: boolean;
  lastFetched: Date | null;
}
```

#### Notification State
```typescript
interface NotificationState {
  notifications: Notification[];
  unreadCount: number;
  isConnected: boolean;
  preferences: NotificationPreferences;
}
```

### GraphQL Response Types
All types are generated from the schema using GraphQL Code Generator. Key types include:

- `User` - User profile and permissions
- `Subscription` - Billing subscription details
- `Plan` - Available subscription plans
- `AnalyticsSummary` - Analytics data
- `Notification` - Notification object
- `FeatureVariant` - Feature flag state

## Error Handling

### Error Boundary Strategy

#### Global Error Boundary
Wrap the entire app in an error boundary that:
- Catches React rendering errors
- Displays user-friendly error page
- Logs errors to analytics service
- Provides recovery options

#### Component-Level Error Boundaries
Wrap major sections (Dashboard, Users, Billing) in error boundaries to:
- Isolate failures
- Allow partial page functionality
- Show section-specific error messages

### GraphQL Error Handling

#### Network Errors
- Display toast notification
- Provide retry button
- Fall back to cached data if available
- Log to analytics service

#### Authentication Errors
- Clear invalid tokens
- Redirect to login page
- Preserve intended destination
- Show appropriate message

#### Validation Errors
- Display field-level errors
- Highlight invalid inputs
- Provide correction guidance
- Prevent form submission

### Loading States

#### Skeleton Screens
Use skeleton loaders for:
- Dashboard cards
- User tables
- Analytics charts
- Notification lists

#### Progressive Loading
- Load critical data first (user profile, permissions)
- Lazy load secondary data (analytics, notifications)
- Show partial UI while loading
- Indicate loading progress

## Testing Strategy

### Unit Testing
**Components to Test**:
- UI components (Button, Card, Input, etc.)
- Form validation logic
- Utility functions (formatters, validators)
- Context providers

**Testing Library**: React Testing Library
**Coverage Target**: >80% for utility functions and hooks

### Integration Testing
**Scenarios to Test**:
- Login flow end-to-end
- Dashboard data loading
- Form submission with validation
- Real-time notification reception
- Feature flag conditional rendering

**Tools**: React Testing Library + MSW (Mock Service Worker) for API mocking

### E2E Testing (Optional)
**Critical Paths**:
- User registration and login
- Subscription checkout flow
- User management operations
- Analytics event tracking

**Tool**: Playwright or Cypress (if time permits)

### Manual Testing Checklist
- [ ] All pages render with proper styling
- [ ] Tailwind CSS classes apply correctly
- [ ] GraphQL queries return data
- [ ] Mutations update UI optimistically
- [ ] Socket.IO connects and receives events
- [ ] Feature flags control visibility
- [ ] Analytics events are tracked
- [ ] Error states display properly
- [ ] Loading states show correctly
- [ ] Responsive design works on mobile
- [ ] Accessibility (keyboard navigation, screen readers)

## Performance Considerations

### Code Splitting
- Lazy load pages with `next/dynamic`
- Split large components
- Defer non-critical JavaScript

### Caching Strategy
- Apollo Client cache for GraphQL data
- Feature flag cache (5-minute TTL)
- Static assets cached by browser
- Service worker for offline support (future enhancement)

### Bundle Optimization
- Tree-shake unused code
- Minimize dependencies
- Use production builds
- Enable compression in Docker

### Real-time Optimization
- Batch Socket.IO events
- Debounce frequent updates
- Use connection pooling
- Implement backpressure handling

## Security Considerations

### Authentication
- JWT tokens stored in memory (not localStorage)
- Refresh token rotation
- Automatic token expiration handling
- CSRF protection for mutations

### Authorization
- Permission checks before rendering UI
- Server-side validation (never trust client)
- Role-based component visibility
- Audit logging for sensitive actions

### Data Protection
- No sensitive data in client-side logs
- Sanitize user inputs
- Escape rendered content
- Use HTTPS in production

### CORS Configuration
- Whitelist specific origins
- Credentials handling
- Preflight request optimization

## Deployment Configuration

### Docker Build Process
```dockerfile
# Multi-stage build
FROM node:18-alpine AS deps
# Install dependencies

FROM node:18-alpine AS builder
# Build Next.js app

FROM node:18-alpine AS runner
# Run Next.js server
CMD ["node", "server.js"]
```

### Environment Variables
**Build-time** (baked into bundle):
- `NEXT_PUBLIC_GRAPHQL_URL`
- `NEXT_PUBLIC_SOCKETIO_URL`
- `NEXT_PUBLIC_ANALYTICS_ENABLED`
- `NEXT_PUBLIC_FEATURE_FLAGS_ENABLED`

**Runtime** (server-side only):
- None needed for frontend (all public vars)

### Health Checks
- HTTP endpoint at `/api/health`
- Returns 200 OK when app is ready
- Checks GraphQL connectivity
- Monitors Socket.IO connection

## Migration Path

### Phase 1: Fix Build Process
1. Create `postcss.config.js`
2. Update `next.config.js` (remove static export)
3. Update Dockerfile for Next.js server
4. Verify Tailwind CSS compilation

### Phase 2: Implement Core Integrations
1. Run GraphQL Code Generator
2. Implement dashboard queries
3. Connect Socket.IO for notifications
4. Test feature flag integration

### Phase 3: Build Missing Pages
1. Create Users page with table
2. Create Billing page with Stripe
3. Create Analytics page with charts
4. Create Notifications page

### Phase 4: Polish and Testing
1. Add loading states
2. Implement error boundaries
3. Add form validation
4. Test all flows manually
5. Fix any remaining issues

## Dependencies

### New Dependencies Needed
```json
{
  "recharts": "^2.10.0",  // For charts
  "date-fns": "^3.0.0",   // For date formatting
  "react-hook-form": "^7.49.0",  // For forms
  "@tanstack/react-table": "^8.11.0"  // For data tables
}
```

### Existing Dependencies (Verified)
- ✅ `@apollo/client` - GraphQL client
- ✅ `socket.io-client` - Real-time communication
- ✅ `@heroicons/react` - Icons
- ✅ `tailwindcss` - Styling
- ✅ `@graphql-codegen/cli` - Type generation

## Success Criteria

1. **Styling Works**: All pages render with Tailwind CSS styling, no barebones HTML
2. **Dashboard Shows Data**: Real analytics, billing, and user data displayed
3. **All Services Integrated**: Users, billing, analytics, notifications, LLM chat all functional
4. **Real-time Works**: Socket.IO connects and delivers notifications
5. **Feature Flags Work**: Conditional rendering based on flags
6. **No Console Errors**: Clean browser console, no GraphQL errors
7. **Responsive Design**: Works on desktop, tablet, and mobile
8. **Accessibility**: Keyboard navigation and screen reader support
9. **Performance**: Page loads in <3 seconds, smooth interactions
10. **Docker Deployment**: Builds and runs successfully in Docker

## Open Questions

1. **Chart Library**: Recharts vs Chart.js vs Victory? (Recommendation: Recharts for React integration)
2. **Table Library**: TanStack Table vs React Table vs custom? (Recommendation: TanStack Table for features)
3. **Form Library**: React Hook Form vs Formik vs custom? (Recommendation: React Hook Form for performance)
4. **Testing Priority**: Focus on unit tests or integration tests first? (Recommendation: Integration tests for critical paths)
5. **Stripe Integration**: Use Stripe Elements or Checkout? (Recommendation: Checkout for simplicity)
