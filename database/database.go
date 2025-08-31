package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establishes connection to PostgreSQL database with RDS KMS support
func ConnectDatabase() {
	var err error
	var dsn string

	// Check if using IAM database authentication
	useIAMAuth := os.Getenv("DB_USE_IAM_AUTH")
	if useIAMAuth == "true" {
		dsn, err = buildIAMAuthDSN()
		if err != nil {
			log.Fatal("Failed to build IAM auth connection string:", err)
		}
	} else {
		dsn = buildStandardDSN()
	}

	// Configure GORM with appropriate settings for RDS
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Set connection pool settings for RDS
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool for RDS with KMS
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Set connection pool parameters optimized for RDS
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	connMaxLifetime := getEnvAsInt("DB_CONN_MAX_LIFETIME", 300) // 5 minutes

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to AWS RDS PostgreSQL database successfully!")
	log.Printf("Connection pool: max_open=%d, max_idle=%d, max_lifetime=%ds",
		maxOpenConns, maxIdleConns, connMaxLifetime)
}

// buildStandardDSN creates a standard PostgreSQL connection string for RDS with KMS
func buildStandardDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Default to require SSL for RDS with KMS
	if sslmode == "" {
		sslmode = "require"
	}

	// Build connection string with additional RDS-specific parameters
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Add additional SSL parameters for RDS
	if sslmode != "disable" {
		dsn += " sslrootcert=rds-ca-2019-root.pem"
	}

	return dsn
}

// buildIAMAuthDSN creates a connection string using AWS IAM database authentication
func buildIAMAuthDSN() (string, error) {
	ctx := context.Background()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = "us-east-1" // Default region
	}

	// Build the endpoint for IAM auth
	endpoint := fmt.Sprintf("%s:%s", host, port)

	// Generate IAM auth token
	authToken, err := auth.BuildAuthToken(ctx, endpoint, region, user, cfg.Credentials)
	if err != nil {
		return "", fmt.Errorf("failed to generate IAM auth token: %w", err)
	}

	// Build connection string with IAM auth token
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, authToken, dbname,
	)

	return dsn, nil
}

// getEnvAsInt gets environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
