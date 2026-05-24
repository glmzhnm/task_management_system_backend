# Task Manager API

A production-ready REST API built with Go (Gin), PostgreSQL, JWT authentication, microservices, and a full frontend UI.

## Architecture

```
backend/
├── main.go                          # Entry point
├── controllers/
│   ├── auth_controller.go           # Register, Login, Profile, ChangePassword
│   └── task_controller.go           # Full CRUD + Stats + Search
├── models/
│   ├── user.go                      # User model (GORM)
│   └── task.go                      # Task model with UserID FK
├── middlewares/
│   ├── auth_middleware.go           # JWT validation + user_id extraction
│   ├── logger_middleware.go         # Colored request logging
│   └── rate_limiter.go             # 100 req/min per IP
├── database/
│   └── db.go                        # GORM connection + golang-migrate
├── routes/
│   └── routes.go                    # All routes with CORS
├── migrations/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── notifier/
│   ├── main.go                      # Microservice (Resty v2 consumer)
│   └── Dockerfile
├── tests/
│   ├── auth_test.go                 # 5 auth handler tests
│   └── task_test.go                 # 5 task handler tests
├── Task_Manager_API.postman_collection.json  # 20 endpoints
├── docker-compose.yml
├── Dockerfile
└── MakeFile
```

## Quick Start

### With Docker (recommended)
```bash
docker-compose up --build
```
- API: http://localhost:8080
- Frontend: http://localhost:8080
- Notifier: http://localhost:8081

### Local Development
```bash
cp .env.example .env        
go mod tidy
make run
```

##  Running Tests
```bash
make test
make test-coverage
```

##  Postman
Import `Task_Manager_API.postman_collection.json` into Postman.
The **Login** request automatically saves the token to `{{token}}` — all protected endpoints use it.

##  API Endpoints

### Auth (public)
| Method | Path                  | Description        |
|--------|-----------------------|--------------------|
| POST   | /auth/register        | Register new user  |
| POST   | /auth/login           | Login → JWT token  |
| GET    | /health               | Health check       |

### Auth (protected — Bearer token)
| Method | Path              | Description           |
|--------|-------------------|-----------------------|
| GET    | /profile          | Get current user      |
| PATCH  | /profile/password | Change password       |
| GET | /profile/stats | Get user profile statistics |
| DELETE | /profile | Delete account |
### Tasks (protected)
| Method | Path | Description |
| :--- | :--- | :--- |
| POST | `/tasks` | Create task |
| GET | `/tasks` | Get all (own) tasks |
| GET | `/tasks/search` | Search by title/desc |
| GET | `/tasks/stats` | Aggregated stats |
| GET | `/tasks/due` | Get tasks due soon |
| GET | `/tasks/status/:status` | Filter by status |
| GET | `/tasks/priority/:priority` | Filter by priority |
| GET | `/tasks/:id` | Get by ID |
| PUT | `/tasks/:id` | Full update |
| PATCH | `/tasks/:id/status` | Update status only |
| PATCH | `/tasks/:id/priority` | Update priority only |
| DELETE | `/tasks/:id` | Delete task |
| DELETE | `/tasks/completed` | Bulk delete completed |

##  Security
- Passwords hashed with bcrypt (cost 12)
- JWT signed with HS256
- All task endpoints scope to authenticated user (ownership enforced)
- Rate limiting: 100 req/min per IP
- CORS configured for cross-origin frontend support
