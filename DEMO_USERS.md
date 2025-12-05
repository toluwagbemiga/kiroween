# Demo Users - Login Credentials

All demo users have been successfully created in the database. You can use these credentials to test different user roles and permissions.

## 🔐 Login Credentials

### ADMIN (Full System Access)
**Full administrative access with all permissions**

- **Email:** `admin@haunted.dev`
- **Password:** `Demo123!`
- **Permissions:**
  - ✅ Users: Read, Write, Delete
  - ✅ Roles: Read, Write, Delete
  - ✅ Permissions: Read
  - ✅ Billing: Read, Write
  - ✅ Analytics: Read
  - ✅ Features: Read, Write

---

### MEMBERS (Standard Access)
**Standard team members with read access to most resources**

#### Member User
- **Email:** `member@haunted.dev`
- **Password:** `Demo123!`

#### John Doe
- **Email:** `john.doe@haunted.dev`
- **Password:** `Demo123!`

#### Jane Smith
- **Email:** `jane.smith@haunted.dev`
- **Password:** `Demo123!`

**Member Permissions:**
- ✅ Users: Read
- ✅ Billing: Read
- ✅ Analytics: Read
- ✅ Features: Read
- ❌ No write or delete permissions

---

### VIEWERS (Read-Only Access)
**Read-only access to view data without making changes**

#### Viewer User
- **Email:** `viewer@haunted.dev`
- **Password:** `Demo123!`

#### Guest User
- **Email:** `guest@haunted.dev`
- **Password:** `Demo123!`

**Viewer Permissions:**
- ✅ Users: Read
- ✅ Roles: Read
- ✅ Permissions: Read
- ✅ Analytics: Read
- ✅ Features: Read
- ❌ No write or delete permissions
- ❌ No billing access

---

## 🚀 Quick Start

1. **Access the application:**
   ```
   http://localhost:3000
   ```

2. **Login with any of the accounts above**

3. **Test different features based on role:**
   - **Admin:** Try creating/editing users, managing roles, accessing billing
   - **Member:** View users and analytics, but notice you can't edit
   - **Viewer:** Browse data but all write actions should be restricted

## 🔄 Re-seeding Data

If you need to reset or re-seed the demo users, run:

```bash
# Windows PowerShell
Get-Content demo/seed-users-simple.sql | docker exec -i kiroproject-postgres-1 psql -U haunted -d haunted

# Linux/Mac
cat demo/seed-users-simple.sql | docker exec -i kiroproject-postgres-1 psql -U haunted -d haunted
```

## 📊 Service Endpoints

- **Frontend:** http://localhost:3000
- **GraphQL API:** http://localhost:4000/graphql
- **GraphQL Playground:** http://localhost:4000/playground
- **Unleash (Feature Flags):** http://localhost:4242

## 🧪 Testing Permissions

### Admin Tests
1. Login as `admin@haunted.dev`
2. Navigate to Users page - should see all users
3. Try creating a new user - should succeed
4. Try editing roles - should succeed
5. Access billing page - should see all plans

### Member Tests
1. Login as `member@haunted.dev`
2. Navigate to Users page - should see users but no edit buttons
3. Try accessing analytics - should work
4. Billing page - read-only access

### Viewer Tests
1. Login as `viewer@haunted.dev`
2. All pages should be read-only
3. No edit/delete buttons should appear
4. Billing page should be restricted

## 🔧 Troubleshooting

### Can't login?
- Ensure all containers are running: `docker-compose ps`
- Check user-auth-service logs: `docker logs kiroproject-user-auth-service-1`
- Verify database connection: `docker exec kiroproject-postgres-1 psql -U haunted -d haunted -c "SELECT email, name FROM users;"`

### Password not working?
All passwords are: `Demo123!` (case-sensitive, includes exclamation mark)

### Need to reset a user?
Run the seed script again - it will update existing users with the correct password hash.
