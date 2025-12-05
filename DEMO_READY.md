# 🎉 Demo Environment Ready!

Your Haunted SaaS system has been successfully set up with demo users and is ready for testing!

## ✅ System Status

All services are running and ready:
- ✅ Frontend (Next.js) - http://localhost:3000
- ✅ GraphQL Gateway - http://localhost:4000
- ✅ User Auth Service - Running
- ✅ Billing Service - Running
- ✅ Analytics Service - Running
- ✅ Notifications Service - Running
- ✅ Feature Flags Service - Running
- ✅ LLM Gateway Service - Running
- ✅ PostgreSQL Database - Healthy
- ✅ Redis Cache - Healthy
- ✅ Unleash (Feature Flags) - http://localhost:4242

## 🔐 Demo User Accounts

### 👑 ADMIN - Full System Access
```
Email:    admin@haunted.dev
Password: Demo123!
```
**Can do:** Everything - create/edit/delete users, manage roles, configure billing, view analytics

### 👥 MEMBERS - Standard Access (3 accounts)
```
Email:    member@haunted.dev
Password: Demo123!

Email:    john.doe@haunted.dev
Password: Demo123!

Email:    jane.smith@haunted.dev
Password: Demo123!
```
**Can do:** View users, billing, analytics, and features (read-only)

### 👁️ VIEWERS - Read-Only Access (2 accounts)
```
Email:    viewer@haunted.dev
Password: Demo123!

Email:    guest@haunted.dev
Password: Demo123!
```
**Can do:** View users, roles, permissions, analytics, and features (read-only, no billing)

## 🚀 Getting Started

1. **Open your browser** and navigate to:
   ```
   http://localhost:3000
   ```

2. **Login** with any of the demo accounts above

3. **Explore the features:**
   - Dashboard - Overview of your account
   - Users - User management (permissions vary by role)
   - Analytics - View usage analytics
   - Billing - Subscription management
   - Notifications - Real-time notifications

## 🧪 Testing Different Roles

### Test Admin Capabilities
1. Login as `admin@haunted.dev`
2. Go to Users page → Click "Add User" (should work)
3. Try editing a user's role (should work)
4. Access Billing → Create/manage subscriptions (should work)

### Test Member Limitations
1. Login as `member@haunted.dev`
2. Go to Users page → Notice no "Add User" button
3. Try to edit users → Should be read-only
4. View Analytics → Should work fine

### Test Viewer Restrictions
1. Login as `viewer@haunted.dev`
2. Browse all pages → Everything is read-only
3. Try to access Billing → Should be restricted
4. No edit/delete buttons anywhere

## 📊 API Access

### GraphQL Playground
Access the interactive GraphQL playground:
```
http://localhost:4000/playground
```

### Example Query (after login)
```graphql
query {
  me {
    id
    email
    name
    roles {
      name
      permissions {
        name
        resource
        action
      }
    }
  }
}
```

## 🔧 Useful Commands

### View all users in database
```powershell
docker exec kiroproject-postgres-1 psql -U haunted -d haunted -c "SELECT email, name FROM users;"
```

### Check service logs
```powershell
# Frontend
docker logs kiroproject-frontend-1 --tail 50

# GraphQL Gateway
docker logs kiroproject-graphql-gateway-1 --tail 50

# User Auth Service
docker logs kiroproject-user-auth-service-1 --tail 50
```

### Restart all services
```powershell
docker-compose restart
```

### Re-seed demo users
```powershell
Get-Content demo/seed-users-simple.sql | docker exec -i kiroproject-postgres-1 psql -U haunted -d haunted
```

## 🎯 What to Test

### Authentication & Authorization
- ✅ Login with different roles
- ✅ Verify permission-based UI changes
- ✅ Test protected routes
- ✅ Check JWT token handling

### User Management
- ✅ View user list (all roles)
- ✅ Create new users (admin only)
- ✅ Edit user details (admin only)
- ✅ Assign roles (admin only)

### Billing
- ✅ View subscription plans
- ✅ Subscribe to a plan (requires Stripe test keys)
- ✅ View subscription status

### Analytics
- ✅ View analytics dashboard
- ✅ Track events
- ✅ View user activity

### Real-time Features
- ✅ Socket.IO notifications
- ✅ Live updates
- ✅ Connection status

### Feature Flags
- ✅ Toggle features via Unleash
- ✅ Gradual rollouts
- ✅ A/B testing

## 📝 Notes

- All demo users have the same password: `Demo123!` (case-sensitive)
- The system uses JWT tokens for authentication
- Sessions are stored in Redis
- All data is stored in PostgreSQL
- Feature flags are managed via Unleash (http://localhost:4242)

## 🐛 Troubleshooting

### Can't login?
1. Check if all containers are running: `docker-compose ps`
2. Verify user exists: `docker exec kiroproject-postgres-1 psql -U haunted -d haunted -c "SELECT email FROM users;"`
3. Check auth service logs: `docker logs kiroproject-user-auth-service-1`

### Frontend not loading?
1. Check frontend logs: `docker logs kiroproject-frontend-1`
2. Verify it's running: `curl http://localhost:3000`
3. Try restarting: `docker-compose restart frontend`

### GraphQL errors?
1. Check gateway logs: `docker logs kiroproject-graphql-gateway-1`
2. Verify backend services are running: `docker-compose ps`
3. Test direct service access via gRPC

## 🎊 You're All Set!

Your demo environment is fully configured and ready to use. Start by logging in as the admin user to explore all features, then try the other roles to see how permissions work.

**Happy testing! 🚀**
