-- CTF Database Initialization Script

-- Create NocoDB database (CTF database is created by POSTGRES_DB)
CREATE DATABASE nocodb_database;

-- Grant permissions to both databases
GRANT ALL PRIVILEGES ON DATABASE ctf_database TO ctf_user;
GRANT ALL PRIVILEGES ON DATABASE nocodb_database TO ctf_user;

-- Create extensions for CTF database
\c ctf_database;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create extensions for NocoDB database
\c nocodb_database;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Print success message
SELECT 'CTF and NocoDB databases initialized successfully!' AS message;