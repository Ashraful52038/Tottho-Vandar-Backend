# Tottho Vandar Backend
[![Go](https://img.shields.io/badge/Go-1.25.7-blue.svg?logo=go&logoColor=white)](https://golang.org)
[![Echo](https://img.shields.io/badge/Echo-v4.15.0-black.svg?logo=go&logoColor=white)](https://echo.labstack.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg?logo=postgresql&logoColor=white)](https://postgresql.org)
[![GORM](https://img.shields.io/badge/GORM-v1.31.1-orange.svg?logo=go&logoColor=white)](https://gorm.io)
[![JWT](https://img.shields.io/badge/JWT-v5-green.svg)](https://jwt.io)
[![Docker](https://img.shields.io/badge/Docker-ready-blue.svg?logo=docker)](https://docker.com)

**Tottho Vandar** is a modern social blogging platform backend built with Go. It provides a complete RESTful API for user management, posts, comments, likes, follows, personalized feeds, tags, authentication, and media uploads. Designed with clean architecture for scalability and maintainability.

## ✨ Features

- 🔐 **Authentication & Authorization**: JWT-based auth, email verification, password reset/change
- 👥 **User Management**: Profiles, follow/unfollow users, view followers/following, likes/comments/posts
- 📝 **Posts**: CRUD, rich HTML content, image uploads, search, pagination, my posts
- 💬 **Comments**: Nested replies, CRUD, post-linked
- ❤️ **Likes**: Toggle likes on posts/comments
- 🏷️ **Tags**: CRUD (admin), popular tags, post tagging (many-to-many)
- 📡 **Feeds**: Personalized (followed users/tags), public feed
- 🖼️ **Media Upload**: Image upload to `/uploads` (static serve)
- 📧 **Emails**: SMTP integration (dev: Mailpit), notifications prepared (RabbitMQ ready)
- 🛡️ **Middleware**: Auth, CORS, logging, recovery, custom validation
- 📊 **Pagination**: Supported on all list endpoints (`page`, `limit`)
- 🐳 **Docker**: Postgres + pgAdmin via `docker-compose`

## 🚧 Upcoming Features (Roadmap)

- 🏷️ **Tag Manager** – Admin panel to create, edit, delete tags; users can follow tags.
- 👥 **Follow System** – Follow users and tags to get a personalized feed.
- 📚 **Reading List / Saved Posts** – Save posts to read later, with a dedicated saved page.
- 📝 **Draft Posts** – Save posts as drafts and publish them later.
- 🕒 **Latest Posts** – Dedicated tab for newest posts first.
- 🌟 **Featured Writers** – Show most followed or admin-selected writers in sidebar.
- 🖼️ **Profile Image Upload** – Users can upload and crop their avatar.
- 🔔 **Email Notifications** – Receive notifications for replies and follows (planned).

## 🏗️ Architecture

```
cmd/api/
├── main.go (Echo server setup)

internal/
├── config/ (env-based config)
├── domain/ (entities: User, Post, Comment, Like, Tag, Follow, Feed, Notification)
├── usecase/ (business logic)
├── repository/ (interfaces)
│   └── impl/postgres/ (GORM impl, DB conn)
└── delivery/http/
    ├── handler/ (endpoint handlers)
    ├── middleware/ (auth)
    └── router/ (Echo routes)
```

## 🚀 Quick Start

### 1. Clone & Install
```bash
git clone https://github.com/Ashraful52038/tottho-vandar-backend.git
cd tottho-vandar-backend
go mod tidy
```

### 2. Environment (.env)
```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tottho_vandar
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-key-change-in-prod
JWT_EXPIRATION=24

# Email (SMTP/Mailpit dev)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
EMAIL_FROM=noreply@totthovandar.com

# Frontend CORS
FRONTEND_URL=http://localhost:3000
```

### 3. Docker (Recommended)
```bash
docker-compose up -d  # Starts Postgres:5432, pgAdmin:5050 (admin/admin@tottho.com)
```

### 4. Run
```bash
# Migrations handled in code (AutoMigrate)
go run cmd/api/main.go
```
Server ready at `http://localhost:8080`

**pgAdmin**: http://localhost:5050 (Server: postgres:5432, user/pass: postgres/postgres)

## 📋 API Endpoints

Base: `/api`

### Auth (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Create account |
| POST | `/auth/login` | JWT token |
| GET | `/auth/verify-email` | Email verification |
| POST | `/auth/forget-password` | Send reset email |
| POST | `/auth/reset-password` | Reset via token |
| POST | `/auth/change-password` | Protected: Change pw |

### Posts (Public + Protected)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/posts` | | All posts (paginated) |
| GET | `/posts/search` | | Search posts |
| GET | `/posts/:id` | | Single post |
| POST | `/posts` | ✅ | Create |
| PUT | `/posts/:id` | ✅ Owner | Update |
| DELETE | `/posts/:id` | ✅ Owner | Delete |
| GET | `/posts/my-posts` | ✅ | My posts |

### Users (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/user/me` | Profile |
| PUT | `/user/me` | Update profile |
| GET/POST/DELETE | `/users/:id/*` | Profile, posts, comments, likes, follow/unfollow |

### Comments, Likes, Tags, Feed, Upload similarly structured (full list in `/internal/delivery/http/router/router.go`)

## 🧪 Testing

```bash
# Run tests
go test ./...

# With coverage
go test -v -cover ./...
```

## 🔧 Makefile Commands
```makefile
# View Makefile for: build, test, lint, docker-build, etc.
```

## 🐳 Docker Compose Services
- `postgres`: Database (5432)
- `pgadmin`: Admin UI (5050)

## 📈 Performance & Scaling
- GORM connection pooling
- JWT stateless
- Pagination everywhere
- Static file serving (/uploads)
- RabbitMQ-ready for email queues

## 🔒 Security
- bcrypt passwords
- JWT expiry (24h)
- Input validation (validator.v10)
- CORS restricted
- SQL injection safe (GORM)

## 🤝 Contributing
1. Fork → Clone → Create feature branch
2. `go mod tidy`
3. Commit: `git commit -m 'feat: add X'`
4. Push & PR to `main`

## 📄 License
MIT – See [LICENSE](LICENSE) (create if needed)

## 👨‍💻 Author
[Ashraful Islam](https://github.com/Ashraful52038) – Fullstack Go Developer

---

*Built with ❤️ using Clean Architecture & Go Best Practices*

