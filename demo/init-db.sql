-- Initialize database with demo data
-- This script is run automatically when the demo sandbox starts

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Note: Individual service migrations will create their own tables
-- This file is for shared setup and demo data population

-- Demo data will be inserted after services start and run their migrations
