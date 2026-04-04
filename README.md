# Tottho Vandar – Backend

[![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8?logo=go)](https://go.dev/)
[![Echo](https://img.shields.io/badge/Echo-4.15.0-000000?logo=go)](https://echo.labstack.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-5.3.1-000000?logo=jsonwebtokens)](https://jwt.io/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.13-FF6600?logo=rabbitmq)](https://www.rabbitmq.com/)

**Tottho Vandar** is a full‑stack blog platform. This is the **backend** built with Go (Echo), GORM, PostgreSQL, JWT, and optional RabbitMQ for email queues.

🔗 **Frontend repository:** [github.com/Ashraful52038/Tottho-Vandar-frontend](https://github.com/Ashraful52038/Tottho-Vandar-frontend)

---

## ✨ Features (Implemented)

- 🔐 **Authentication** – Register, login, email verification, password reset, change password
- 📝 **Post Management** – CRUD operations, rich content (HTML), featured images, tags
- 🏷️ **Tag System** – Admin can create/edit/delete tags; tags are linked to posts (many‑to‑many)
- 💬 **Comments & Replies** – Nested comments, delete, mention support
- ❤️ **Likes** – Like/unlike posts and comments
- 👥 **User Profiles** – View user info, posts, comments, likes, followers/following
- 🔔 **Follow System** – Follow users and tags → personalised feed
- 📸 **Image Upload** – Store images locally (./uploads) and serve statically
- 📧 **Email Notifications** – Verification, password reset, reply notifications (SMTP + optional RabbitMQ queue)
- 📄 **Pagination** – All list endpoints support `page` & `limit`
- 🛡️ **JWT Authentication** – Stateless, 24h expiry, middleware protected routes
- 🐳 **Docker Support** – Ready‑to‑use `docker-compose.yml` for PostgreSQL and pgAdmin

---

## 🧭 Roadmap (Upcoming Features)

- 📚 **Reading List / Saved Posts** – Save posts to read later, with a dedicated saved page
- 📝 **Draft Posts** – Save posts as drafts and publish them later
- 🕒 **Latest Posts** – Dedicated tab for newest posts first
- 🌟 **Featured Writers** – Show most followed or admin‑selected writers in sidebar
- 🖼️ **Profile Image Upload** – Users can upload and crop their avatar
- 🔔 **Full Email Notification System** – Notifications for replies and follows (already planned)

---

## 🛠️ Tech Stack

- **Language:** Go 1.25.7
- **Web Framework:** Echo v4.15.0
- **ORM:** GORM v1.31.1
- **Database:** PostgreSQL 15
- **Authentication:** JWT (golang-jwt/jwt/v5)
- **Validation:** go-playground/validator/v10
- **Email:** SMTP (Mailpit for development) + optional RabbitMQ (streadway/amqp)
- **Password Hashing:** bcrypt (via golang.org/x/crypto)
- **Configuration:** godotenv (`.env` file)

---

## 🚀 Getting Started

### Prerequisites

- Go 1.25.7+
- Docker & Docker Compose (optional, for PostgreSQL)
- Make (optional, for convenience)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ashraful52038/Tottho-Vandar-Backend.git
   cd Tottho-Vandar-Backend