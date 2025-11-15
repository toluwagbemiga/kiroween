# Implementation Plan

- [x] 1. Fix build configuration and Tailwind CSS compilation





  - Create `postcss.config.js` file with Tailwind and Autoprefixer plugins
  - Update `next.config.js` to remove static export mode and enable SSR
  - Update frontend Dockerfile to use Next.js standalone build and run `next start`
  - Verify Tailwind CSS classes render correctly in the browser
  - _Requirements: 1.1, 1.2, 1.3, 8.1, 8.2, 8.3_

- [ ] 2. Set up GraphQL code generation and type safety
  - Install additional dependencies: `recharts`, `date-fns`, `react-hook-form`, `@tanstack/react-table`
  - Run GraphQL Code Generator to create TypeScript types and hooks from schema
  - Verify generated files in `src/lib/graphql/generated/` directory
  - Update imports in existing components to use generated types
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [ ] 3. Implement dashboard service integrations
- [ ] 3.1 Create GraphQL queries for dashboard data
  - Write `MY_ANALYTICS_QUERY` to fetch analytics summary
  - Write `MY_SUBSCRIPTION_QUERY` to fetch billing information
  - Write `USERS_QUERY` for user count statistics
  - Add queries to dashboard page component
  - _Requirements: 2.1, 2.2, 2.3_

- [ ] 3.2 Update dashboard page with real data
  - Replace placeholder stats with data from `useMyAnalyticsQuery` hook
  - Display subscription status from `useMySubscriptionQuery` hook
  - Implement loading states with skeleton loaders
  - Implement error handling with error boundaries
  - Remove feature flag fallback once data loads successfully
  - _Requirements: 2.1, 2.2, 2.4, 2.5, 10.1, 10.2_

- [ ] 3.3 Add dashboard charts and visualizations
  - Install and configure Recharts library
  - Create `AnalyticsChart` component for time-series data
  - Create `MetricCard` component for KPI display
  - Add charts to dashboard showing event trends
  - _Requirements: 2.1, 2.4_

- [ ] 4. Create users management page
- [ ] 4.1 Implement users list page
  - Create `/users/page.tsx` with DashboardLayout
  - Write `USERS_QUERY` with pagination support
  - Create `UserTable` component using TanStack Table
  - Implement search and filter functionality
  - Add loading and empty states
  - _Requirements: 3.1, 10.1_

- [ ] 4.2 Implement user creation and editing
  - Create `UserForm` component with React Hook Form
  - Write `CREATE_USER_MUTATION` and `UPDATE_USER_MUTATION`
  - Create modal dialog for user form
  - Implement form validation with Zod schema
  - Add role selection dropdown
  - Handle mutation errors and success states
  - _Requirements: 3.2, 3.3, 10.2, 10.3_

- [ ] 4.3 Implement user role and permission management
  - Create `RoleSelector` component for multi-select
  - Write `ASSIGN_ROLE_MUTATION` and `REMOVE_ROLE_MUTATION`
  - Create `PermissionMatrix` component for visual display
  - Implement permission checks before showing admin actions
  - Update Apollo cache after role changes
  - _Requirements: 3.3, 3.4, 3.5_

- [ ] 5. Create billing management page
- [ ] 5.1 Implement billing overview page
  - Create `/billing/page.tsx` with DashboardLayout
  - Write `PLANS_QUERY` to fetch available plans
  - Write `MY_SUBSCRIPTION_QUERY` for current subscription
  - Create `PlanCard` component to display plan details
  - Display current plan with features and pricing
  - _Requirements: 4.1, 10.1_

- [ ] 5.2 Implement subscription checkout flow
  - Write `CREATE_SUBSCRIPTION_CHECKOUT_MUTATION`
  - Create `CheckoutButton` component with Stripe integration
  - Handle redirect to Stripe Checkout
  - Implement success and cancel callback pages
  - Update subscription status after successful payment
  - _Requirements: 4.2, 4.4, 10.2_

- [ ] 5.3 Implement subscription management
  - Write `UPDATE_SUBSCRIPTION_MUTATION` and `CANCEL_SUBSCRIPTION_MUTATION`
  - Create upgrade/downgrade flow UI
  - Implement cancellation confirmation dialog
  - Display billing history with invoice list
  - Add download invoice functionality
  - _Requirements: 4.2, 4.3, 4.5_

- [ ] 6. Create analytics dashboard page
- [ ] 6.1 Implement analytics overview page
  - Create `/analytics/page.tsx` with DashboardLayout
  - Write `MY_ANALYTICS_QUERY` with date range parameters
  - Create `DateRangePicker` component for date selection
  - Display total events and unique users metrics
  - Create `EventTable` component for event log display
  - _Requirements: 2.1, 10.1_

- [ ] 6.2 Add analytics visualizations
  - Create `AnalyticsChart` component with Recharts
  - Display events by type in pie chart
  - Display event trends over time in line chart
  - Create `MetricCard` component for KPI display
  - Implement chart loading states
  - _Requirements: 2.1, 10.1_

- [ ] 6.3 Implement analytics event tracking
  - Write `TRACK_EVENT_MUTATION` for custom events
  - Enhance `analytics.ts` utility with batching
  - Auto-track page views on navigation
  - Track user interactions (button clicks, form submissions)
  - Track errors and performance metrics
  - _Requirements: 10.4_

- [ ] 7. Create notifications management page
- [ ] 7.1 Implement notifications list page
  - Create `/notifications/page.tsx` with DashboardLayout
  - Enhance Socket.IO integration to receive notifications
  - Create `NotificationList` component with scrollable feed
  - Create `NotificationItem` component for individual display
  - Display unread count in header
  - Implement mark as read functionality
  - _Requirements: 5.1, 5.2, 5.4_

- [ ] 7.2 Implement notification preferences
  - Write `MY_NOTIFICATION_PREFERENCES_QUERY`
  - Write `UPDATE_NOTIFICATION_PREFERENCES_MUTATION`
  - Create `PreferencesForm` component with channel toggles
  - Implement email, push, and in-app preference controls
  - Save preferences and update UI
  - _Requirements: 5.4, 10.2_

- [ ] 8. Enhance real-time Socket.IO integration
- [ ] 8.1 Improve Socket.IO connection management
  - Update `SocketProvider` with auto-reconnect logic
  - Implement JWT authentication for Socket.IO
  - Add connection state indicators in UI
  - Handle disconnect and reconnect events
  - Subscribe to user-specific channels
  - _Requirements: 5.1, 5.2, 5.3_

- [ ] 8.2 Integrate Socket.IO with Apollo Cache
  - Update Apollo cache when notifications arrive via Socket.IO
  - Refresh relevant queries on user update events
  - Update subscription status on billing events
  - Show toast notifications for important real-time events
  - _Requirements: 2.4, 5.2_

- [ ] 9. Enhance AI chat widget integration
- [ ] 9.1 Implement chat functionality
  - Review existing `ChatWidget` component
  - Write `CALL_PROMPT_MUTATION` for LLM requests
  - Write `AVAILABLE_PROMPTS_QUERY` for prompt templates
  - Implement message sending and receiving
  - Display chat history during session
  - _Requirements: 6.1, 6.2, 6.4_

- [ ] 9.2 Add chat UI enhancements
  - Implement typing indicators
  - Add message timestamps
  - Handle streaming responses (if supported)
  - Add error handling for failed LLM calls
  - Implement chat minimize/maximize functionality
  - _Requirements: 6.3, 6.5, 10.2_

- [ ] 10. Implement feature flag integration
- [ ] 10.1 Enhance feature flag context
  - Review existing `FeatureFlagContext`
  - Write `IS_FEATURE_ENABLED_QUERY` with caching
  - Write `FEATURE_VARIANT_QUERY` for A/B testing
  - Implement 5-minute cache TTL for flags
  - Add loading states while fetching flags
  - _Requirements: 7.1, 7.3, 10.1_

- [ ] 10.2 Apply feature flags to new pages
  - Wrap new dashboard features with `Feature` component
  - Add feature flags for advanced analytics
  - Add feature flags for billing features
  - Track feature usage in analytics
  - Re-evaluate flags on user context changes
  - _Requirements: 7.2, 7.4, 7.5_

- [ ] 11. Implement comprehensive error handling
- [ ] 11.1 Create error boundaries
  - Create global error boundary component
  - Wrap app in global error boundary
  - Create section-specific error boundaries for Dashboard, Users, Billing
  - Log errors to analytics service
  - Display user-friendly error messages
  - _Requirements: 10.2, 10.4_

- [ ] 11.2 Implement loading states
  - Create skeleton loader components for cards, tables, charts
  - Add loading states to all data-fetching components
  - Implement progressive loading (critical data first)
  - Add loading indicators to buttons during mutations
  - Show partial UI while secondary data loads
  - _Requirements: 10.1, 10.5_

- [ ] 11.3 Add form validation and error display
  - Implement Zod schemas for all forms
  - Display field-level validation errors
  - Highlight invalid inputs with error styling
  - Provide correction guidance in error messages
  - Prevent form submission when validation fails
  - _Requirements: 10.2, 10.3_

- [ ] 12. Create reusable UI components
- [ ] 12.1 Build data table component
  - Create `DataTable` component with TanStack Table
  - Implement sortable columns
  - Add pagination controls
  - Implement row selection
  - Add loading and empty states
  - _Requirements: 3.1, 10.1_

- [ ] 12.2 Build form components
  - Create `FormField` component for labeled inputs
  - Create `FormError` component for error messages
  - Create `FormSection` component for grouped fields
  - Create `SubmitButton` component with loading state
  - Ensure WCAG 2.1 Level AA accessibility
  - _Requirements: 1.5, 3.2, 4.2_

- [ ] 12.3 Build chart components
  - Create `LineChart` component for time series
  - Create `BarChart` component for comparisons
  - Create `PieChart` component for distributions
  - Create `StatCard` component for single metrics
  - Add responsive sizing and tooltips
  - _Requirements: 2.1_

- [ ] 13. Update Docker configuration and deployment
- [ ] 13.1 Update frontend Dockerfile
  - Modify Dockerfile to use Next.js standalone build
  - Change CMD to run `node server.js` instead of static serve
  - Ensure environment variables are passed correctly
  - Update health check to use Next.js health endpoint
  - _Requirements: 8.3, 8.4, 8.5_

- [ ] 13.2 Update docker-compose configuration
  - Verify all service ports are correct in docker-compose.yml
  - Ensure frontend can reach GraphQL gateway at port 4000
  - Ensure frontend can reach Socket.IO at port 3002
  - Update environment variables for frontend service
  - Test full stack startup with `docker-compose up`
  - _Requirements: 8.4, 8.5_

- [ ] 14. Testing and verification
- [ ] 14.1 Manual testing checklist
  - Verify all pages render with Tailwind CSS styling
  - Test login and registration flows
  - Test dashboard data loading from all services
  - Test user management CRUD operations
  - Test billing subscription flows
  - Test analytics data display
  - Test notifications real-time delivery
  - Test AI chat widget functionality
  - Test feature flag conditional rendering
  - Verify responsive design on mobile and tablet
  - _Requirements: All requirements_

- [ ] 14.2 Verify error handling and edge cases
  - Test behavior when GraphQL queries fail
  - Test behavior when Socket.IO disconnects
  - Test form validation with invalid inputs
  - Test loading states during slow network
  - Test authentication token expiration
  - Verify error messages are user-friendly
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [ ] 14.3 Performance verification
  - Measure page load times (target <3 seconds)
  - Check bundle size and optimize if needed
  - Verify code splitting is working
  - Test with throttled network connection
  - Monitor memory usage during extended use
  - _Requirements: 8.3_

- [ ] 14.4 Accessibility verification
  - Test keyboard navigation on all pages
  - Test with screen reader (NVDA or JAWS)
  - Verify color contrast ratios meet WCAG 2.1 AA
  - Check focus indicators are visible
  - Verify form labels and ARIA attributes
  - _Requirements: 1.5_
