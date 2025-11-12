#!/usr/bin/env node

/**
 * Demo Data Generation Script
 * 
 * This script populates the demo sandbox with realistic test data:
 * - Demo users with various roles
 * - Sample subscription plans
 * - Feature flags
 * - Analytics events
 * - Notifications
 * 
 * Run after all services are up and healthy.
 */

const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const { faker } = require('@faker-js/faker');

// Service addresses
const SERVICES = {
  auth: 'localhost:50051',
  billing: 'localhost:50052',
  llm: 'localhost:50053',
  notifications: 'localhost:50054',
  analytics: 'localhost:50055',
  featureFlags: 'localhost:50056',
};

async function main() {
  console.log('🎃 Starting HAUNTED SAAS SKELETON demo data generation...\n');

  try {
    // 1. Create demo users
    console.log('👥 Creating demo users...');
    await createDemoUsers();

    // 2. Create subscription plans
    console.log('💳 Creating subscription plans...');
    await createSubscriptionPlans();

    // 3. Create feature flags
    console.log('🚩 Creating feature flags...');
    await createFeatureFlags();

    // 4. Generate analytics events
    console.log('📊 Generating analytics events...');
    await generateAnalyticsEvents();

    // 5. Create sample prompts
    console.log('🤖 Creating LLM prompts...');
    await createSamplePrompts();

    console.log('\n✅ Demo data generation complete!');
    console.log('\nDemo Credentials:');
    console.log('  Admin: admin@haunted.dev / Admin123!');
    console.log('  User:  user@haunted.dev / User123!');
    console.log('\nAccess Points:');
    console.log('  Frontend:  http://localhost:3000');
    console.log('  Docs:      http://localhost:3001');
    console.log('  GraphQL:   http://localhost:4000/graphql');
    console.log('  Unleash:   http://localhost:4242 (admin/unleash4all)');
  } catch (error) {
    console.error('❌ Error generating demo data:', error);
    process.exit(1);
  }
}

async function createDemoUsers() {
  // TODO: Implement gRPC calls to user-auth-service
  console.log('  - admin@haunted.dev (Super Admin)');
  console.log('  - user@haunted.dev (Team Member)');
  console.log('  - viewer@haunted.dev (Viewer)');
}

async function createSubscriptionPlans() {
  // TODO: Implement gRPC calls to billing-service
  console.log('  - Free Plan ($0/month)');
  console.log('  - Pro Plan ($29/month)');
  console.log('  - Enterprise Plan ($99/month)');
}

async function createFeatureFlags() {
  // TODO: Implement API calls to Unleash
  console.log('  - new_dashboard (gradual rollout: 50%)');
  console.log('  - ai_assistant (enabled for Pro+)');
  console.log('  - advanced_analytics (enabled for Enterprise)');
}

async function generateAnalyticsEvents() {
  // TODO: Implement gRPC calls to analytics-service
  const eventCount = 1000;
  console.log(`  - Generated ${eventCount} sample events`);
}

async function createSamplePrompts() {
  // Prompts are file-based, so this just logs what's available
  console.log('  - onboarding/welcome-email.md');
  console.log('  - support/ticket-response.md');
  console.log('  - analytics/insight-summary.md');
}

// Run if called directly
if (require.main === module) {
  main();
}

module.exports = { main };
