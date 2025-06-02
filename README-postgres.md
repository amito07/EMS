# PostgreSQL Docker Setup for EMS

This directory contains Docker configuration for running PostgreSQL locally for the EMS (Education Management System) project.

## Files

- `Dockerfile.postgres`: Custom PostgreSQL Docker image
- `docker-compose.yml`: Docker Compose configuration
- `init-scripts/01_init.sql`: Database initialization script

## Quick Start

### Option 1: Using Docker Compose (Recommended)

```powershell
# Start PostgreSQL container
docker-compose up -d

# Check container status
docker-compose ps

# View logs
docker-compose logs postgres

# Stop the container
docker-compose down
```

### Option 2: Using Docker directly

```powershell
# Build the custom PostgreSQL image
docker build -f Dockerfile.postgres -t ems-postgres .

# Run the container
docker run -d `
  --name ems_postgres `
  -e POSTGRES_DB=ems_db `
  -e POSTGRES_USER=ems_user `
  -e POSTGRES_PASSWORD=ems_password `
  -p 5432:5432 `
  -v ems_postgres_data:/var/lib/postgresql/data `
  ems-postgres

# Check if container is running
docker ps
```

## Database Connection Details

- **Host**: localhost
- **Port**: 5432
- **Database**: ems_db
- **Username**: ems_user
- **Password**: ems_password

## Connection String Examples

### Go (for your EMS application)

```go
connectionString := "host=localhost port=5432 user=ems_user password=ems_password dbname=ems_db sslmode=disable"
```

### psql command line

```powershell
psql -h localhost -p 5432 -U ems_user -d ems_db
```

## Management Commands

```powershell
# Connect to database using psql inside container
docker exec -it ems_postgres psql -U ems_user -d ems_db

# Backup database
docker exec ems_postgres pg_dump -U ems_user ems_db > backup.sql

# Restore database
docker exec -i ems_postgres psql -U ems_user -d ems_db < backup.sql

# Check database logs
docker logs ems_postgres

# Remove container and data (WARNING: This deletes all data)
docker-compose down -v
```

## Database Schema

The initialization script creates the following tables:

- `students`: Student information
- `teachers`: Teacher information
- `courses`: Course details
- `enrollments`: Student-course relationships

## Notes

- The database data is persisted in a Docker volume named `postgres_data`
- The container includes health checks to ensure PostgreSQL is ready
- Sample data is automatically inserted during first startup
- The container will restart automatically unless stopped manually
