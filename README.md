# CTF Backend

A robust backend system for Capture The Flag (CTF) competitions built with Go and Gin framework.

## Features

- **User Management**: Registration, authentication, and user profiles
- **Challenge Management**: Create, update, and manage CTF challenges
- **Flag Submission**: Secure flag submission with rate limiting
- **Leaderboard**: Real-time scoring and rankings
- **Admin Panel**: Administrative controls for managing users and challenges
- **Rate Limiting**: Protection against brute force attacks
- **CORS Support**: Cross-origin resource sharing enabled
- **Database Integration**: PostgreSQL with GORM ORM

## Tech Stack

- **Backend**: Go 1.23+ with Gin framework
- **Database**: PostgreSQL (AWS RDS recommended for production)
- **ORM**: GORM
- **Authentication**: SHA-256 password hashing with custom tokens
- **Deployment**: AWS EC2, ECS, or Elastic Beanstalk

## Project Structure

```
CTF-backend/
├── api/
│   └── routes/           # HTTP route definitions
├── controllers/          # Request handlers
│   ├── admin_controller.go
│   ├── challenge_controller.go
│   └── user_controller.go
├── middleware/           # HTTP middleware
│   ├── admin_middleware.go
│   ├── auth_middleware.go
│   └── rate_limit_middleware.go
├── models/              # Database models
│   ├── challenge.go
│   ├── models.go
│   ├── submission.go
│   └── user.go
├── database/            # Database connection
│   └── database.go
├── main.go             # Application entry point
├── .env.example        # Environment variables template
└── AWS_DEPLOYMENT.md   # AWS deployment guide
```

## Local Development

### Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/TheLostLeo/CTF-backend.git
   cd CTF-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on the port specified in your `.env` file (default: 8080).

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `GET /api/v1/challenges` - List all challenges
- `GET /api/v1/challenges/:id` - Get specific challenge
- `GET /api/v1/leaderboard` - Get leaderboard
- `POST /api/v1/register` - User registration (rate limited)
- `POST /api/v1/login` - User login (rate limited)

### Protected Endpoints (Requires Authentication)

- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/challenges/:id/submit` - Submit flag (rate limited)

### Admin Endpoints (Requires Admin Access)

- `POST /api/v1/admin/challenges` - Create challenge
- `PUT /api/v1/admin/challenges/:id` - Update challenge
- `DELETE /api/v1/admin/challenges/:id` - Delete challenge
- `GET /api/v1/admin/users` - List all users
- `GET /api/v1/admin/dashboard` - Admin dashboard

## Authentication

The application uses a simple token-based authentication system. For production use, consider implementing JWT tokens.

**Headers required for protected endpoints:**
```
Authorization: Bearer <token>
```

## Rate Limiting

- **Authentication endpoints**: 10 requests per minute per IP
- **Flag submission**: 5 requests per minute per IP

## Database Models

### User
- ID, Username, Email, Password (hashed)
- Score, IsAdmin flag
- Created/Updated timestamps, soft delete support

### Challenge
- ID, Title, Description, Category
- Points, Flag, Hint, File URL
- Active status, Created/Updated timestamps

### Submission
- ID, User ID, Challenge ID
- Submitted Flag, Is Correct flag
- Created timestamp

## Deployment

For production deployment on AWS, see [AWS_DEPLOYMENT.md](AWS_DEPLOYMENT.md) for detailed instructions.

### Quick AWS Deployment Steps

1. Create AWS RDS PostgreSQL database
2. Launch EC2 instance or use Elastic Beanstalk
3. Configure environment variables
4. Build and deploy the application
5. Set up load balancer and SSL (recommended)

## Environment Variables

```bash
# Application
PORT=8080
GIN_MODE=release

# Database (AWS RDS)
DB_HOST=your-rds-endpoint.region.rds.amazonaws.com
DB_PORT=5432
DB_USER=your_db_username
DB_PASSWORD=your_db_password
DB_NAME=ctf_database
DB_SSLMODE=require

# Security
JWT_SECRET=your-super-secret-jwt-key
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Security Considerations

- Change default JWT secret in production
- Use HTTPS in production
- Regularly update dependencies
- Monitor for security vulnerabilities
- Use AWS Secrets Manager for sensitive configuration in production

## Support

For issues and questions, please create an issue in the GitHub repository.