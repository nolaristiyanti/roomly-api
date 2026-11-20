# Roomly API

REST API for **Room & Facility Booking Management**, built with Go, Gin, PostgreSQL, and Redis.

## Tech Stack

- **Language:** Go 1.27.2
- **Framework:** Gin
- **Database:** PostgreSQL 17
- **Cache & Session:** Redis 7
- **ORM:** GORM
- **Authentication:** JWT
- **API Documentation:** Swagger
- **Infrastructure:** Docker & Docker Compose

## Features

- User registration and login
- JWT authentication with Redis session management
- Role-based authorization (`admin` and `user`)
- User, room, and facility management
- Room booking and cancellation
- Booking validation and availability checking
- Database transaction and concurrency handling
- Structured logging and request ID
- Swagger API documentation

## Project Structure

```text
roomly-api/
├── common/
├── config/
├── database/
│   ├── migrations/
│   └── seeds/
├── docs/
├── domain/
│   ├── auth/
│   ├── booking/
│   ├── facility/
│   ├── room/
│   └── user/
├── middleware/
├── routes/
├── Dockerfile
├── compose.yaml
├── .env.example
├── go.mod
├── go.sum
└── main.go
```

## Getting Started

### Prerequisites

- Git
- Docker Desktop or Docker Engine
- Docker Compose

### 1. Clone Repository

```bash
git clone https://github.com/nolaristiyanti/roomly-api.git
cd roomly-api
```

### 2. Configure Environment

Copy the example environment file:

**Windows PowerShell**
```powershell
Copy-Item .env.example .env
```

**Linux / macOS**
```bash
cp .env.example .env
```

Update the `.env` file with your local configuration and credentials.

### 3. Create Docker Resources

Roomly currently uses an external Docker network and named volumes.

```bash
docker network create roomly-network

docker volume create roomly-postgres-data

docker volume create roomly-redis-data
```

### 4. Build and Run

```bash
docker compose up -d --build
```

Check container status:

```bash
docker compose ps
```

### 5. Database Migration and Seeding

After PostgreSQL is ready, execute the database migrations and seed scripts from the `database/` directory according to the project setup.

Ensure the database is initialized before using the API.

## Access the API

| Service           | URL                                       |
| API Health Check  | http://localhost:8080/health              |
| Swagger UI        | http://localhost:8080/swagger/index.html  |
| PostgreSQL        | localhost:5433                            |
| Redis             | localhost:6380                            |

## Docker Commands

```bash
# View running containers
docker compose ps

# View application logs
docker compose logs -f app

# Restart application
docker compose restart app

# Stop containers
docker compose down
```

**Note:** PostgreSQL and Redis data are stored in persistent Docker volumes. Avoid deleting these volumes unless you intend to remove the stored data.

## Booking Rules

- Bookings are available Monday to Friday.
- Operating hours: 08:00–17:00 (Asia/Jakarta).
- Booking duration uses full-hour intervals.
- Minimum duration: 1 hour.
- Maximum duration: 9 hours.
- Lunch time (12:00–13:00) remains bookable.
- Only confirmed bookings block room availability.

## API Documentation

Swagger UI is available at:

```text
http://localhost:8080/swagger/index.html
```

## License

This project is developed for learning and portfolio purposes.