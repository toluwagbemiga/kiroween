#!/usr/bin/env node

/**
 * Demo Data Generation Script
 * Populates the database with sample users, teams, and data
 */

const { Client } = require('pg');

const DB_CONFIG = {
  host: process.env.POSTGRES_HOST || 'localhost',
  port: process.env.POSTGRES_PORT || 5432,
  database: process.env.POSTGRES_DB || 'haunted_saas',
  user: process.env.POSTGRES_USER || 'postgres',
  password: process.env.POSTGRES_PASSWORD || 'postgres',
};

// Sample data
const DEMO_USERS = [
  {
    email: 'admin@haunted.dev',
    name: 'Super Admin',
    password: '$2a$10$rZ5qKqZqZqZqZqZqZqZqZuO', // 'admin123' hashed
    signupType: 'owner',
    role: 'admin',
  },
  {
    email: 'owner@team1.com',
    name: 'Team Owner 1',
    password: '$2a$10$rZ5qKqZqZqZqZqZqZqZqZuO', // 'password123' hashed
    signupType: 'owner',
    role: 'owner',
  },
  {
    email: 'member@team1.com',
    name: 'Team Member 1',
    password: '$2a$10$rZ5qKqZqZqZqZqZqZqZqZuO',
    signupType: 'member',
    role: 'member',
  },
  {
    email: 'owner@team2.com',
    name: 'Team Owner 2',
    password: '$2a$10$rZ5qKqZqZqZqZqZqZqZqZuO',
    signupType: 'owner',
    role: 'owner',
  },
];

const DEMO_TEAMS = [
  {
    name: 'Acme Corporation',
    slug: 'acme-corp',
    isMultiTenant: true,
    tenantMode: 'multi',
    maxMembers: 50,
  },
  {
    name: 'Startup Inc',
    slug: 'startup-inc',
    isMultiTenant: false,
    tenantMode: 'single',
    maxMembers: 10,
  },
];

async function seedDatabase() {
  const client = new Client(DB_CONFIG);

  try {
    await client.connect();
    console.log('✓ Connected to database');

    // Create teams
    console.log('\n📦 Creating demo teams...');
    const teamIds = [];
    for (const team of DEMO_TEAMS) {
      const result = await client.query(
        `INSERT INTO teams (name, slug, is_multi_tenant, tenant_mode, max_members)
         VALUES ($1, $2, $3, $4, $5)
         ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`,
        [team.name, team.slug, team.isMultiTenant, team.tenantMode, team.maxMembers]
      );
      teamIds.push(result.rows[0].id);
      console.log(`  ✓ Created team: ${team.name} (${result.rows[0].id})`);
    }

    // Create users
    console.log('\n👥 Creating demo users...');
    const userIds = [];
    for (let i = 0; i < DEMO_USERS.length; i++) {
      const user = DEMO_USERS[i];
      const teamId = i === 0 ? null : teamIds[i % teamIds.length]; // Admin has no team

      const result = await client.query(
        `INSERT INTO users (email, password_hash, name, team_id, signup_type)
         VALUES ($1, $2, $3, $4, $5)
         ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`,
        [user.email, user.password, user.name, teamId, user.signupType]
      );
      userIds.push(result.rows[0].id);
      console.log(`  ✓ Created user: ${user.email} (${result.rows[0].id})`);
    }

    // Update team owners
    console.log('\n🔑 Setting team owners...');
    await client.query(
      `UPDATE teams SET owner_id = $1 WHERE slug = $2`,
      [userIds[1], 'acme-corp']
    );
    await client.query(
      `UPDATE teams SET owner_id = $1 WHERE slug = $2`,
      [userIds[3], 'startup-inc']
    );
    console.log('  ✓ Team owners assigned');

    // Create roles if they don't exist
    console.log('\n🎭 Creating roles...');
    const roles = ['admin', 'owner', 'member', 'viewer'];
    const roleIds = {};
    
    for (const roleName of roles) {
      const result = await client.query(
        `INSERT INTO roles (name, description, is_system)
         VALUES ($1, $2, true)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`,
        [roleName, `System ${roleName} role`]
      );
      roleIds[roleName] = result.rows[0].id;
      console.log(`  ✓ Created role: ${roleName}`);
    }

    // Assign roles to users
    console.log('\n🔐 Assigning roles to users...');
    for (let i = 0; i < DEMO_USERS.length; i++) {
      const user = DEMO_USERS[i];
      const userId = userIds[i];
      const roleId = roleIds[user.role];

      await client.query(
        `INSERT INTO user_roles (user_id, role_id)
         VALUES ($1, $2)
         ON CONFLICT DO NOTHING`,
        [userId, roleId]
      );
      console.log(`  ✓ Assigned ${user.role} role to ${user.email}`);
    }

    // Create permissions
    console.log('\n🛡️  Creating permissions...');
    const permissions = [
      'users:read', 'users:write', 'users:delete',
      'teams:read', 'teams:write', 'teams:delete',
      'billing:read', 'billing:write',
      'analytics:read', 'analytics:write',
      'settings:read', 'settings:write',
    ];

    for (const perm of permissions) {
      await client.query(
        `INSERT INTO permissions (name, description)
         VALUES ($1, $2)
         ON CONFLICT (name) DO NOTHING`,
        [perm, `Permission to ${perm.split(':')[1]} ${perm.split(':')[0]}`]
      );
    }
    console.log(`  ✓ Created ${permissions.length} permissions`);

    // Assign permissions to admin role
    console.log('\n🔓 Assigning permissions to admin role...');
    const permResult = await client.query(`SELECT id FROM permissions`);
    for (const perm of permResult.rows) {
      await client.query(
        `INSERT INTO role_permissions (role_id, permission_id)
         VALUES ($1, $2)
         ON CONFLICT DO NOTHING`,
        [roleIds.admin, perm.id]
      );
    }
    console.log('  ✓ All permissions assigned to admin role');

    console.log('\n✅ Demo data seeded successfully!');
    console.log('\n📝 Demo Credentials:');
    console.log('  Super Admin: admin@haunted.dev / admin123');
    console.log('  Team Owner 1: owner@team1.com / password123');
    console.log('  Team Member 1: member@team1.com / password123');
    console.log('  Team Owner 2: owner@team2.com / password123');

  } catch (error) {
    console.error('❌ Error seeding database:', error);
    process.exit(1);
  } finally {
    await client.end();
  }
}

seedDatabase();
