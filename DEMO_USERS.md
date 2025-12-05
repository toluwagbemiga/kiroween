# Demo Users & Test Accounts

## 🔐 Quick Login Credentials

### Super Admin Account
- **Email**: admin@haunted.dev
- **Password**: Demo123!
- **Role**: Super Admin
- **Permissions**: Full system access
- **Dashboard**: `/admin` - System administration panel
- **Use Case**: System administration, user management, role assignment

### Team Owner Accounts

#### Team 1 - Acme Corporation (Multi-Tenant SaaS)
- **Email**: owner@team1.com
- **Password**: Demo123!
- **Role**: Team Owner
- **Team**: Acme Corporation (Multi-tenant enabled)
- **Use Case**: Managing a SaaS with multiple customer tenants
- **Features**: Can enable multi-tenancy, invite members, manage billing

#### Team 2 - Startup Inc (Single Tenant)
- **Email**: owner@team2.com
- **Password**: Demo123!
- **Role**: Team Owner
- **Team**: Startup Inc (Single tenant mode)
- **Use Case**: Managing a single organization
- **Features**: Team management, member invitations

### Team Member Accounts

#### Member 1 (Acme Corporation)
- **Email**: member@team1.com
- **Password**: Demo123!
- **Role**: Team Member
- **Team**: Acme Corporation
- **Use Case**: Regular team member with limited permissions

#### Member 2 (General)
- **Email**: member@haunted.dev
- **Password**: Demo123!
- **Role**: Member
- **Use Case**: Standard access testing

### Viewer Accounts (Read-Only)

#### Viewer 1
- **Email**: viewer@haunted.dev
- **Password**: Demo123!
- **Role**: Viewer
- **Use Case**: Read-only access testing

#### Guest User
- **Email**: guest@haunted.dev
- **Password**: Demo123!
- **Role**: Viewer
- **Use Case**: Guest/limited access testing

---

## 🚀 Quick Start

### 1. Access the Application
```
Frontend: http://localhost:3000
GraphQL API: http://localhost:4000/graphql
Admin Dashboard: http://localhost:3000/admin
```

### 2. Login with Demo Accounts
Choose any account above based on what you want to test.

### 3. Explore Features
- **Dashboard**: `/dashboard` - Main user dashboard
- **Admin Panel**: `/admin` - Super admin only
- **Settings**: `/settings` - User and team settings
- **AI Assistant**: `/ai` - AI-powered features
- **Billing**: `/billing` - Subscription management

---

## 📦 Seeding Demo Data

### Using Node.js Script
```bash
cd demo
npm install pg
node seed-data.js
```

### Using Docker
```bash
docker-compose exec user-auth-service node /demo/seed-data.js
```

### Using SQL Directly
```bash
# Windows PowerShell
Get-Content demo/seed-users-simple.sql | docker exec -i kiroproject-postgres-1 psql -U haunted -d haunted

# Linux/Mac
cat demo/seed-users-simple.sql | docker exec -i kiroproject-postgres-1 psql -U haunted -d haunted
```

---

## 🏢 Team & Multi-Tenancy Features

### Multi-Tenancy (Acme Corporation)
- ✅ Enabled for SaaS use cases
- ✅ Supports multiple isolated customer environments
- ✅ Owner can manage tenant settings
- ✅ Up to 50 team members
- ✅ Advanced billing and analytics

### Single Tenancy (Startup Inc)
- ✅ Traditional single organization mode
- ✅ Simpler permission model
- ✅ Up to 10 team members
- ✅ Basic collaboration features

---

## 👥 Signup Flow

### As Team Owner (Create New Team)
1. Navigate to `/signup`
2. Click "Create Team (Owner)"
3. Fill in:
   - Full Name
   - Email address
   - Password
   - Confirm Password
4. Submit form
5. Automatically creates a new team
6. Redirected to dashboard
7. Can enable multi-tenancy in `/settings` → Team tab

### As Team Member (Join Existing Team)
1. Navigate to `/signup`
2. Click "Join Team (Member)"
3. Fill in:
   - Full Name
   - Email address
   - Password
   - Confirm Password
   - **Team ID** (get from team owner)
   - Invitation Token (optional)
4. Submit form
5. Status: **Pending** (awaiting owner approval)
6. Team owner approves in admin panel
7. Status changes to **Active**

---

## 🎛️ Admin Dashboard Features

Access at `/admin` with super admin account (admin@haunted.dev)

### Features:
- 📊 **System Stats**: Total users, roles, system status
- 👥 **User Management**: View all users, assign roles
- 🎭 **Role Management**: View and manage system roles
- 🔐 **Permission Overview**: See all permissions and assignments
- 📈 **Activity Monitoring**: Track system usage

### Admin Actions:
- Assign roles to any user
- View user details and team membership
- Monitor system health
- Manage global settings

---

## 🧪 Testing Scenarios

### Scenario 1: Multi-Tenant SaaS Setup
1. Login as `owner@team1.com`
2. Go to `/settings` → Team tab
3. Check "Enable Multi-Tenancy (SaaS Mode)"
4. Save settings
5. Invite team members
6. Manage customer tenants
7. Configure billing plans

### Scenario 2: Team Collaboration
1. Login as `owner@team2.com`
2. Go to `/settings` → Team tab
3. Copy your Team ID
4. Share Team ID with new member
5. New member signs up with Team ID
6. Approve member in admin panel
7. Assign appropriate role

### Scenario 3: System Administration
1. Login as `admin@haunted.dev`
2. Navigate to `/admin`
3. View all users and teams
4. Assign "admin" role to a user
5. Monitor system status
6. Review role permissions

### Scenario 4: Member Signup & Approval
1. Logout (if logged in)
2. Go to `/signup`
3. Select "Join Team (Member)"
4. Enter Team ID: (get from owner@team1.com)
5. Complete signup
6. Login as `owner@team1.com`
7. Go to team management
8. Approve pending member
9. Assign role to new member

---

## 🔑 Role Permissions Matrix

| Feature | Admin | Owner | Member | Viewer |
|---------|-------|-------|--------|--------|
| View Users | ✅ | ✅ | ✅ | ✅ |
| Create Users | ✅ | ✅ | ❌ | ❌ |
| Edit Users | ✅ | ✅ | ❌ | ❌ |
| Delete Users | ✅ | ❌ | ❌ | ❌ |
| Manage Roles | ✅ | ❌ | ❌ | ❌ |
| View Billing | ✅ | ✅ | ✅ | ❌ |
| Manage Billing | ✅ | ✅ | ❌ | ❌ |
| View Analytics | ✅ | ✅ | ✅ | ✅ |
| Team Settings | ✅ | ✅ | ❌ | ❌ |
| Multi-Tenancy | ✅ | ✅ | ❌ | ❌ |

---

## 🔧 Troubleshooting

### Can't Login?
```bash
# Check if services are running
docker-compose ps

# Check user-auth-service logs
docker logs kiroproject-user-auth-service-1

# Verify database
docker exec kiroproject-postgres-1 psql -U haunted -d haunted -c "SELECT email, name FROM users;"
```

### Password Not Working?
- All passwords are: `Demo123!` (case-sensitive)
- Includes exclamation mark at the end
- No spaces before or after

### Dashboard Pages Not Loading?
- The `(dashboard)` folder is a Next.js route group
- Access pages directly:
  - `/ai` - AI Assistant
  - `/settings` - Settings
  - `/billing` - Billing
- NOT `/dashboard/ai` or `/(dashboard)/ai`

### Team ID Not Found?
```bash
# Get all team IDs
docker exec kiroproject-postgres-1 psql -U haunted -d haunted -c "SELECT id, name, slug FROM teams;"
```

### Need to Reset Data?
```bash
# Re-run seed script
cd demo
node seed-data.js
```

---

## 📊 Service Endpoints

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | Next.js application |
| GraphQL API | http://localhost:4000/graphql | API Gateway |
| GraphQL Playground | http://localhost:4000/playground | Interactive API explorer |
| User Auth Service | http://localhost:3001 | gRPC service |
| Notifications | http://localhost:3002 | Socket.IO service |
| pgAdmin | http://localhost:5050 | Database UI |
| Unleash | http://localhost:4242 | Feature flags |

---

## 📝 Notes

- All demo users are created with the same password for easy testing
- Super admin has no team assignment (system-wide access)
- Team owners are automatically assigned to their teams
- Members must be approved by team owners
- Multi-tenancy can be toggled in team settings
- OSI-approved MIT License included in LICENSE file
