# Connecting EC2 Instance to RDS with KMS Encryption

This guide explains how to securely connect an EC2 instance to AWS RDS PostgreSQL with KMS encryption for your CTF backend application.

## Prerequisites

- AWS Account with appropriate permissions
- EC2 instance running Amazon Linux 2/2023
- Basic understanding of AWS services

## Step 1: Create KMS Key for RDS Encryption

### Using AWS Console
1. Navigate to **AWS KMS Console**
2. Click **"Create key"**
3. Select **"Symmetric"** key type
4. Choose **"Encrypt and decrypt"** usage
5. Add key description: "CTF Database Encryption Key"
6. Configure key policy to allow RDS service access
7. **Note the Key ID** for RDS configuration

### Using AWS CLI
```bash
# Create KMS key
aws kms create-key \
    --description "CTF Database Encryption Key" \
    --key-usage ENCRYPT_DECRYPT \
    --key-spec SYMMETRIC_DEFAULT

# Create alias for easier reference
aws kms create-alias \
    --alias-name alias/ctf-database-key \
    --target-key-id <KEY_ID_FROM_ABOVE>
```

## Step 2: Create RDS PostgreSQL with KMS Encryption

### Using AWS Console
1. Go to **AWS RDS Console**
2. Click **"Create database"**
3. Choose **"PostgreSQL"**
4. Select PostgreSQL version 13.0 or higher
5. Choose instance class:
   - **t3.micro** for testing
   - **t3.small+** for production
6. **Enable Encryption**:
   - ✅ Check "Enable encryption"
   - Select your custom KMS key or use AWS managed key
7. Configure database settings:
   - **DB Instance Identifier**: `ctf-database`
   - **Master Username**: `ctf_user`
   - **Master Password**: Create a secure password
   - **Database Name**: `ctf_database`
8. **Network & Security**:
   - **VPC**: Select your VPC
   - **Subnet group**: Choose private subnets (recommended)
   - **Public access**: **No** (recommended for security)
   - **VPC security groups**: Create new security group
9. **Additional Configuration**:
   - ✅ Enable automated backups
   - Set backup retention period (7-30 days)
   - ✅ Enable Performance Insights (optional)

### Using AWS CLI
```bash
# Create RDS instance with KMS encryption
aws rds create-db-instance \
    --db-instance-identifier ctf-database \
    --db-instance-class db.t3.micro \
    --engine postgres \
    --engine-version 15.4 \
    --master-username ctf_user \
    --master-user-password YOUR_SECURE_PASSWORD \
    --allocated-storage 20 \
    --storage-type gp3 \
    --storage-encrypted \
    --kms-key-id alias/ctf-database-key \
    --db-name ctf_database \
    --vpc-security-group-ids sg-xxxxxxxxx \
    --db-subnet-group-name your-subnet-group \
    --backup-retention-period 7 \
    --no-publicly-accessible
```

## Step 3: Configure Security Groups

### Create Security Group for RDS
```bash
# Create security group for RDS
aws ec2 create-security-group \
    --group-name ctf-rds-sg \
    --description "Security group for CTF RDS database" \
    --vpc-id vpc-xxxxxxxxx

# Allow PostgreSQL access from EC2 security group
aws ec2 authorize-security-group-ingress \
    --group-id sg-rds-xxxxxxxxx \
    --protocol tcp \
    --port 5432 \
    --source-group sg-ec2-xxxxxxxxx
```

### Security Group Rules
- **Inbound**: Port 5432 (PostgreSQL) from EC2 security group only
- **Outbound**: Default (all traffic allowed)

## Step 4: Configure EC2 Instance

### Install Required Packages
```bash
# Update system
sudo yum update -y

# Install PostgreSQL client for testing
sudo yum install postgresql -y

# Install Go (if not already installed)
sudo amazon-linux-extras install golang1.21 -y
```

### Create IAM Role for EC2 (Optional - for IAM Database Authentication)
```bash
# Create IAM policy for RDS access
cat > rds-policy.json << EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "rds-db:connect"
            ],
            "Resource": [
                "arn:aws:rds-db:us-east-1:ACCOUNT-ID:dbuser:ctf-database/ctf_user"
            ]
        }
    ]
}
EOF

# Create and attach policy
aws iam create-policy \
    --policy-name CTFDatabaseAccess \
    --policy-document file://rds-policy.json

aws iam attach-role-policy \
    --role-name YourEC2Role \
    --policy-arn arn:aws:iam::ACCOUNT-ID:policy/CTFDatabaseAccess
```

## Step 5: Test Database Connection

### Test with PostgreSQL Client
```bash
# Test standard authentication
psql "host=ctf-database.xxxxxxxxx.us-east-1.rds.amazonaws.com port=5432 dbname=ctf_database user=ctf_user sslmode=require"
```

### Download SSL Certificate (Recommended)
```bash
# Download RDS CA certificate
wget https://s3.amazonaws.com/rds-downloads/rds-ca-2019-root.pem

# Test with SSL certificate verification
psql "host=ctf-database.xxxxxxxxx.us-east-1.rds.amazonaws.com port=5432 dbname=ctf_database user=ctf_user sslmode=verify-ca sslrootcert=rds-ca-2019-root.pem"
```

## Step 6: Configure Application Environment

### Environment Variables
```bash
# Create .env file for your application
cat > .env << EOF
# Application Configuration
PORT=8080
GIN_MODE=release

# AWS RDS Database Configuration with KMS Encryption
DB_HOST=ctf-database.xxxxxxxxx.us-east-1.rds.amazonaws.com
DB_PORT=5432
DB_USER=ctf_user
DB_PASSWORD=your_secure_password
DB_NAME=ctf_database
DB_SSLMODE=require

# Database Connection Pool Settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=300

# AWS Configuration
AWS_REGION=us-east-1

# JWT Secret
JWT_SECRET=$(openssl rand -base64 32)
EOF
```

## Step 7: Security Best Practices

### Network Security
- ✅ Use **private subnets** for RDS
- ✅ **No public access** to RDS
- ✅ Restrict security groups to application tier only
- ✅ Use **VPC endpoints** for AWS services (optional)

### Encryption
- ✅ **KMS encryption** at rest (configured in Step 2)
- ✅ **SSL/TLS encryption** in transit (`sslmode=require`)
- ✅ Use **customer-managed KMS keys** for better control

### Access Control
- ✅ Use **IAM database authentication** (optional)
- ✅ **Rotate passwords** regularly if using standard auth
- ✅ Apply **principle of least privilege**
- ✅ Enable **database activity streams** (optional)

### Monitoring
- ✅ Enable **Performance Insights**
- ✅ Set up **CloudWatch alarms**
- ✅ Enable **Enhanced Monitoring**
- ✅ Configure **automated backups**

## Step 8: Troubleshooting

### Common Connection Issues

#### 1. Connection Timeout
```bash
# Check security groups
aws ec2 describe-security-groups --group-ids sg-xxxxxxxxx

# Test network connectivity
telnet your-rds-endpoint.region.rds.amazonaws.com 5432
```

#### 2. SSL Certificate Issues
```bash
# Verify SSL configuration
openssl s_client -connect your-rds-endpoint.region.rds.amazonaws.com:5432 -starttls postgres
```

#### 3. Authentication Failures
```bash
# Check RDS logs
aws rds describe-db-log-files --db-instance-identifier ctf-database
aws rds download-db-log-file-portion --db-instance-identifier ctf-database --log-file-name error/postgresql.log.2024-08-31-18
```

### Monitoring Commands
```bash
# Check RDS status
aws rds describe-db-instances --db-instance-identifier ctf-database

# Monitor connections
aws cloudwatch get-metric-statistics \
    --namespace AWS/RDS \
    --metric-name DatabaseConnections \
    --dimensions Name=DBInstanceIdentifier,Value=ctf-database \
    --start-time 2024-08-31T00:00:00Z \
    --end-time 2024-08-31T23:59:59Z \
    --period 3600 \
    --statistics Average
```

## Important Notes

- **KMS Keys**: Customer-managed keys provide better audit trails and control
- **SSL/TLS**: Always use `sslmode=require` or higher for production
- **Backups**: Automated backups are also encrypted with the same KMS key
- **Cross-Region**: KMS keys are region-specific
- **Performance**: Connection pooling is crucial for RDS performance
- **Costs**: Monitor RDS usage and optimize instance sizing

This setup ensures your CTF backend application connects securely to AWS RDS with KMS encryption, providing enterprise-grade security for your database communications.