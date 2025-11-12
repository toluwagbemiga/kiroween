# Feature Flags Integration Guide

## Overview

The feature flags service has been successfully integrated into the frontend application, creating a "Control Layer" that allows you to remotely control UI features without deploying new code. This enables:

- **A/B Testing**: Test different UI variations with different user segments
- **Gradual Rollouts**: Roll out new features to a percentage of users
- **Kill Switches**: Instantly disable problematic features
- **User Segmentation**: Show different features to different user groups
- **Development Flags**: Hide work-in-progress features from production users

## Architecture

### Backend (GraphQL Gateway)

**GraphQL Schema Updates:**
- Added `enabledFeatures: [String!]!` field to the `User` type
- The `me` query now fetches enabled feature flags from the feature-flags-service

**Implementation:**
- Location: `app/gateway/graphql-api-gateway/internal/resolvers/query.resolvers.go`
- The `Me` query calls `GetUserFeatures` gRPC endpoint
- Returns an array of enabled feature flag names (e.g., `["bento-grid", "new-chat-feature"]`)
- Gracefully handles feature-flags-service unavailability

### Feature Flags Service

**New gRPC Endpoint:**
- `GetUserFeatures(userId, teamId, properties)`: Returns all enabled features for a user
- Evaluates all feature flags in Unleash and returns only enabled ones
- Location: `app/services/feature-flags-service/internal/grpc_handlers.go`

### Frontend Integration

#### 1. Feature Flag Context (`src/contexts/FeatureFlagContext.tsx`)

Provides global access to feature flags:

```typescript
const { 
  enabledFeatures,    // Array of enabled feature names
  isFeatureEnabled,   // Function to check if a feature is enabled
  loading             // Loading state
} = useFeatureFlags();
```

**Features:**
- Fetches feature flags via GraphQL `me` query
- Caches flags in React Context
- Only fetches when user is authenticated
- Provides type-safe access to flags

#### 2. Feature Component (`src/components/Feature.tsx`)

Declarative component for feature-gated UI:

```typescript
<Feature name="bento-grid">
  <BentoGridLayout />
</Feature>
```

**Features:**
- Conditionally renders children based on feature flag
- Supports fallback UI when feature is disabled
- Handles loading states
- Zero runtime overhead when feature is disabled

#### 3. useFeature Hook

Programmatic access to feature flags:

```typescript
const isNewFeatureEnabled = useFeature('new-feature');
if (isNewFeatureEnabled) {
  // Do something
}
```

## Usage Examples

### Basic Feature Gating

```typescript
import { Feature } from '@/components/Feature';

function MyPage() {
  return (
    <div>
      <h1>My Page</h1>
      
      <Feature name="new-dashboard">
        <NewDashboard />
      </Feature>
    </div>
  );
}
```

### With Fallback UI

```typescript
<Feature 
  name="bento-grid"
  fallback={
    <Card>
      <CardContent>
        <p>The new layout is coming soon!</p>
      </CardContent>
    </Card>
  }
>
  <BentoGridLayout />
</Feature>
```

### Programmatic Feature Checks

```typescript
import { useFeature } from '@/components/Feature';

function MyComponent() {
  const hasNewFeature = useFeature('new-feature');
  const hasAdvancedMode = useFeature('advanced-mode');

  const handleClick = () => {
    if (hasNewFeature) {
      // Use new implementation
      newImplementation();
    } else {
      // Use old implementation
      oldImplementation();
    }
  };

  return (
    <button onClick={handleClick}>
      {hasAdvancedMode ? 'Advanced Action' : 'Basic Action'}
    </button>
  );
}
```

### Multiple Features

```typescript
<Feature name="feature-a">
  <FeatureA />
</Feature>

<Feature name="feature-b">
  <FeatureB />
</Feature>

<Feature name="feature-c" fallback={<OldFeatureC />}>
  <NewFeatureC />
</Feature>
```

### Nested Features

```typescript
<Feature name="new-dashboard">
  <Dashboard>
    <Feature name="advanced-charts">
      <AdvancedCharts />
    </Feature>
    
    <Feature name="real-time-updates">
      <RealTimeUpdates />
    </Feature>
  </Dashboard>
</Feature>
```

## Dashboard Example

The dashboard page demonstrates feature flag usage:

```typescript
<Feature 
  name="bento-grid"
  fallback={
    <Card>
      <CardContent className="p-6">
        <p className="text-gray-300">
          The new dashboard layout is not available yet. Check back soon!
        </p>
      </CardContent>
    </Card>
  }
>
  {/* Bento Grid Layout */}
  <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
    {/* Stats cards */}
  </div>
</Feature>
```

## Backend Configuration (Unleash)

The feature-flags-service uses Unleash for feature flag management. To configure:

### 1. Set Environment Variables

```bash
# Unleash Server URL
UNLEASH_URL=https://your-unleash-instance.com/api

# Unleash API Token
UNLEASH_API_TOKEN=your-api-token

# Environment (development, staging, production)
UNLEASH_ENVIRONMENT=production

# Application Name
UNLEASH_APP_NAME=haunted-saas
```

### 2. Create Feature Flags in Unleash

1. Log into your Unleash dashboard
2. Create a new feature flag (e.g., `bento-grid`)
3. Configure targeting rules:
   - **User IDs**: Target specific users
   - **User Properties**: Target by email domain, role, etc.
   - **Gradual Rollout**: Enable for X% of users
   - **Custom Strategies**: Complex targeting logic

### 3. Common Strategies

**Gradual Rollout:**
```
Strategy: Gradual rollout by userId
Percentage: 25%
```

**User Targeting:**
```
Strategy: UserIDs
User IDs: user-123, user-456
```

**Team Targeting:**
```
Strategy: Custom
Property: teamId
Operator: IN
Values: team-a, team-b
```

## Feature Flag Naming Conventions

Use kebab-case for feature flag names:

- ✅ `bento-grid`
- ✅ `new-chat-feature`
- ✅ `advanced-analytics`
- ✅ `dark-mode`
- ❌ `bentoGrid`
- ❌ `NEW_CHAT_FEATURE`
- ❌ `AdvancedAnalytics`

## Best Practices

### 1. Always Provide Fallbacks

```typescript
// Good
<Feature name="new-feature" fallback={<OldFeature />}>
  <NewFeature />
</Feature>

// Bad - users see nothing if feature is disabled
<Feature name="new-feature">
  <NewFeature />
</Feature>
```

### 2. Keep Feature Flags Temporary

Feature flags should be temporary. Once a feature is fully rolled out:
1. Remove the `<Feature>` wrapper
2. Delete the feature flag from Unleash
3. Clean up the old code path

### 3. Test Both States

Always test your application with the feature both enabled and disabled:

```typescript
// In development, you can override flags
localStorage.setItem('feature-override', JSON.stringify({
  'bento-grid': true,
  'new-feature': false
}));
```

### 4. Use Descriptive Names

```typescript
// Good
<Feature name="new-dashboard-layout">
<Feature name="advanced-search">
<Feature name="real-time-notifications">

// Bad
<Feature name="feature1">
<Feature name="test">
<Feature name="new">
```

### 5. Document Your Flags

Keep a list of active feature flags and their purpose:

```typescript
/**
 * Active Feature Flags:
 * 
 * - bento-grid: New dashboard layout with bento-style cards
 * - advanced-search: Enhanced search with filters and facets
 * - real-time-updates: WebSocket-based real-time data updates
 * - dark-mode: Dark theme support
 */
```

## Troubleshooting

### Feature flag not working

1. **Check if user is authenticated**: Feature flags only work for logged-in users
2. **Verify flag exists in Unleash**: Check your Unleash dashboard
3. **Check flag name**: Ensure the name matches exactly (case-sensitive)
4. **Check browser console**: Look for GraphQL errors
5. **Verify backend is running**: Ensure feature-flags-service is accessible

### Feature always disabled

1. **Check Unleash targeting rules**: Verify your user matches the targeting criteria
2. **Check environment**: Ensure `UNLEASH_ENVIRONMENT` matches your Unleash setup
3. **Check API token**: Verify `UNLEASH_API_TOKEN` is valid
4. **Check logs**: Look at feature-flags-service logs for errors

### Performance concerns

Feature flags are highly optimized:
- Flags are fetched once on login and cached
- No additional API calls per feature check
- Zero runtime overhead for disabled features
- Unleash SDK uses in-memory cache

## Advanced Usage

### Dynamic Feature Checks

```typescript
function DynamicFeatures() {
  const { enabledFeatures } = useFeatureFlags();
  
  return (
    <div>
      <h2>Your Enabled Features:</h2>
      <ul>
        {enabledFeatures.map(feature => (
          <li key={feature}>{feature}</li>
        ))}
      </ul>
    </div>
  );
}
```

### Feature-based Routing

```typescript
function AppRoutes() {
  const hasNewDashboard = useFeature('new-dashboard');
  
  return (
    <Routes>
      <Route 
        path="/dashboard" 
        element={hasNewDashboard ? <NewDashboard /> : <OldDashboard />} 
      />
    </Routes>
  );
}
```

### Analytics Integration

```typescript
function FeatureWithAnalytics() {
  const { trackEvent } = useAnalytics();
  const isEnabled = useFeature('new-feature');
  
  useEffect(() => {
    trackEvent('feature_flag_evaluated', {
      featureName: 'new-feature',
      enabled: isEnabled
    });
  }, [isEnabled]);
  
  return (
    <Feature name="new-feature">
      <NewFeature />
    </Feature>
  );
}
```

## Migration Guide

### Removing a Feature Flag

Once a feature is fully rolled out:

1. **Remove the Feature wrapper:**
   ```typescript
   // Before
   <Feature name="new-feature">
     <NewFeature />
   </Feature>
   
   // After
   <NewFeature />
   ```

2. **Delete old code paths:**
   ```typescript
   // Before
   <Feature name="new-feature" fallback={<OldFeature />}>
     <NewFeature />
   </Feature>
   
   // After
   <NewFeature />
   // Delete OldFeature component
   ```

3. **Remove from Unleash:** Archive the feature flag in Unleash dashboard

4. **Update documentation:** Remove from active flags list

## Next Steps

Consider adding feature flags for:
- New UI components or layouts
- Experimental features
- Beta features for specific users
- Performance optimizations
- Third-party integrations
- A/B test variations
- Seasonal features (e.g., holiday themes)
