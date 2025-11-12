# Gateway Fix Plan

## Root Cause
The gateway code was written assuming proto methods and fields that don't exist in the actual service proto files.

## Issues Found

### 1. User Auth Service
**Missing RPC Methods:**
- `GetUser` - dataloader and query resolvers call this
- `ListUsers` - query resolver calls this
- `ListRoles` - query resolver calls this  
- `GetRole` - query resolver calls this
- `ChangePassword` - mutation resolver calls this
- `UpdateUser` - mutation resolver calls this
- `AssignRole` - mutation resolver calls this
- `RemoveRole` - mutation resolver calls this

**Missing Fields:**
- `User.TeamId` - doesn't exist in proto
- `User.Permissions` - doesn't exist as direct field (only in roles)
- `RegisterRequest.TeamId` - doesn't exist
- `RegisterResponse.Token` - doesn't exist (only has User)
- `RegisterResponse.RefreshToken` - doesn't exist
- `RegisterResponse.ExpiresAt` - doesn't exist
- `LoginResponse.Token` - field is `access_token` not `Token`
- `LogoutRequest.UserId` - doesn't exist
- `LogoutRequest.Token` - field is `session_token`

### 2. Billing Service
**Missing RPC Methods:**
- `GetUserSubscription` - query resolver calls this
- `CreatePortalSession` - query resolver calls this

**Missing Fields:**
- `GetSubscriptionRequest.SubscriptionId` - field is `team_id`
- `CreateCheckoutSessionRequest.UserId` - field is `team_id`
- `CancelSubscriptionRequest.UserId` - field is `team_id`
- `UpdateSubscriptionRequest.UserId` - field is `team_id`
- `CreateCheckoutSessionResponse.Url` - field is `checkout_url`
- `CreateCustomerPortalSessionResponse.Url` - field is `portal_url`
- `Subscription.UserId` - field is `team_id`
- `Subscription.CancelAtPeriodEnd` - doesn't exist (has `cancel_at`)

### 3. LLM Gateway Service
**Missing RPC Methods:**
- `CallLLM` - mutation resolver calls this

**Missing Fields:**
- `CallPromptRequest.UserId` - doesn't exist
- `CallPromptRequest.PromptName` - field is `prompt_path`
- `CallPromptRequest.VariablesJson` - field is `variables_json`
- `CallPromptResponse.Content` - field is `response_text`
- `CallPromptResponse.Model` - field is `model_used`
- `CallPromptResponse.TokensUsed` - is object `token_usage.total_tokens`
- `CallPromptResponse.Cost` - doesn't exist
- `CallPromptResponse.FinishReason` - doesn't exist
- `GetPromptMetadataRequest.PromptName` - field is `prompt_path`
- `GetPromptMetadataResponse.Metadata` - doesn't exist (fields are at top level)
- `GetUsageStatsRequest.UserId` - doesn't exist
- `GetUsageStatsResponse.CallsByModelJson` - doesn't exist
- `GetUsageStatsResponse.TotalCalls` - field is `total_requests`
- `GetUsageStatsResponse.TotalCost` - doesn't exist

### 4. Notifications Service
**Missing RPC Methods:**
- `SendNotification` - mutation resolver calls this
- `UpdatePreferences` - mutation resolver calls this
- `MarkAsRead` - mutation resolver calls this
- `GenerateConnectionToken` - query resolver calls this
- `GetPreferences` - query resolver calls this

### 5. Analytics Service
**Missing RPC Methods:**
- `GetAnalyticsSummary` - query resolver calls this

**Missing Fields:**
- `TrackEventRequest.PropertiesJson` - field is `properties` (map)
- `TrackEventRequest.Timestamp` - field is `timestamp` (int64)
- `IdentifyUserRequest.PropertiesJson` - field is `properties` (map)

## Solution Strategy

Since all 6 services build successfully, the issue is the gateway is out of sync. We have two options:

1. **Update all proto files** to add missing methods and fields (LARGE EFFORT)
2. **Fix gateway to work with existing protos** (FASTER)

I'll go with option 2 - fix the gateway to work with what exists.
