# Analytics Integration Guide

## Overview

The analytics service has been successfully integrated into the frontend application. This integration allows automatic tracking of user events, page views, and custom analytics throughout the application.

## Architecture

### Backend (GraphQL Gateway)

**GraphQL Mutations Added:**
- `trackEvent(input: TrackEventInput!)`: Track custom events with properties
- `identifyUser(properties: JSON!)`: Update user properties for analytics

**Implementation:**
- Location: `app/gateway/graphql-api-gateway/internal/resolvers/mutation.resolvers.go`
- The mutations make gRPC calls to the analytics-service
- Automatically enriches events with user ID from authentication context
- Handles JSON serialization of event properties

### Frontend Integration

#### 1. Analytics Hook (`src/lib/analytics.ts`)

Custom React hook that provides analytics functionality:

```typescript
const { 
  trackEvent,        // Track custom events
  identifyUser,      // Update user properties
  trackPageView,     // Track page views
  trackClick,        // Track button clicks
  trackFormSubmit,   // Track form submissions
  trackError,        // Track errors
  isEnabled          // Check if analytics is enabled
} = useAnalytics();
```

**Features:**
- Automatic user context enrichment (userId, email, timestamp, URL)
- Only tracks when user is authenticated
- Graceful error handling (analytics failures don't break the app)
- Development mode logging
- Type-safe with TypeScript

#### 2. Automatic Page View Tracking

**Component:** `src/app/analytics-tracker.tsx`

- Automatically tracks page views on every route change
- Integrated into root layout
- Uses Next.js `usePathname()` and `useSearchParams()` hooks
- Only tracks for authenticated users

#### 3. Login Event Tracking

**Location:** `src/app/login/page.tsx`

Tracks `user_logged_in` event on successful authentication with properties:
- `method`: 'login' or 'register'
- Plus automatic enrichment (userId, email, timestamp, URL)

## Usage Examples

### Track Custom Events

```typescript
import { useAnalytics } from '@/lib/analytics';

function MyComponent() {
  const { trackEvent } = useAnalytics();

  const handlePurchase = () => {
    trackEvent('purchase_completed', {
      productId: '123',
      amount: 99.99,
      currency: 'USD'
    });
  };

  return <button onClick={handlePurchase}>Buy Now</button>;
}
```

### Track Button Clicks

```typescript
const { trackClick } = useAnalytics();

<button onClick={() => {
  trackClick('upgrade_button', { plan: 'pro' });
  handleUpgrade();
}}>
  Upgrade to Pro
</button>
```

### Track Form Submissions

```typescript
const { trackFormSubmit } = useAnalytics();

const handleSubmit = (data) => {
  trackFormSubmit('contact_form', {
    fields: Object.keys(data),
    source: 'homepage'
  });
  // ... submit logic
};
```

### Track Errors

```typescript
const { trackError } = useAnalytics();

try {
  await riskyOperation();
} catch (error) {
  trackError(error.message, {
    operation: 'riskyOperation',
    context: 'checkout'
  });
}
```

### Identify Users

```typescript
const { identifyUser } = useAnalytics();

// Update user properties for analytics
identifyUser({
  plan: 'pro',
  company: 'Acme Inc',
  role: 'admin'
});
```

## Event Properties

All events are automatically enriched with:
- `userId`: Current user's ID
- `userEmail`: Current user's email
- `timestamp`: ISO 8601 timestamp
- `url`: Full URL of the page
- `pathname`: URL pathname

## Tracked Events

### Automatic Events
- `page_viewed`: Tracked on every route change
  - Properties: `pageName`, `referrer`

### Manual Events
- `user_logged_in`: Tracked on successful login/registration
  - Properties: `method` ('login' or 'register')

### Custom Events
You can track any custom event using `trackEvent()`:
- `button_clicked`
- `form_submitted`
- `error_occurred`
- `purchase_completed`
- `feature_used`
- etc.

## Backend Configuration

The analytics service supports multiple providers:
- **Mixpanel**: Set `MIXPANEL_TOKEN` environment variable
- **Amplitude**: Set `AMPLITUDE_API_KEY` environment variable
- **Segment**: Set `SEGMENT_WRITE_KEY` environment variable
- **Multi-provider**: Set multiple tokens to send to all providers

See `app/services/analytics-service/README.md` for backend configuration details.

## Testing

### Development Mode
In development, all analytics events are logged to the console:
```
[Analytics] Event tracked: page_viewed { userId: '123', pathname: '/dashboard', ... }
```

### Production Mode
Events are sent to configured analytics providers without console logging.

## Best Practices

1. **Don't over-track**: Only track meaningful user interactions
2. **Use descriptive event names**: Use snake_case (e.g., `button_clicked`, not `btnClick`)
3. **Include relevant properties**: Add context that helps understand user behavior
4. **Handle errors gracefully**: Analytics failures shouldn't break your app (already handled by the hook)
5. **Respect privacy**: Don't track sensitive information (passwords, credit cards, etc.)
6. **Test in development**: Check console logs to verify events are tracked correctly

## Troubleshooting

### Events not appearing in analytics dashboard

1. Check if user is authenticated (analytics only works for logged-in users)
2. Verify backend analytics service is running
3. Check environment variables for analytics provider tokens
4. Look for errors in browser console
5. Check backend logs for gRPC errors

### TypeScript errors

Make sure to regenerate GraphQL types after schema changes:
```bash
cd app/frontend
npm run codegen
```

## Next Steps

Consider adding tracking for:
- Feature usage (which features are most popular)
- User onboarding flow completion
- Subscription upgrades/downgrades
- Search queries
- Export/download actions
- Settings changes
- Error rates by page/feature
