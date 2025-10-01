# Social Media Backend

A social media backend API built with Go (Golang) that enables users to create accounts, follow other users, post content, and interact through likes, comments, and notifications.

## 📋 Project Overview

**Social Media Backend** is a RESTful API service for a social networking platform. It provides user authentication, profile management, content creation, feeds, and notifications — all designed with scalability, low latency (via caching), and reliability in mind.

### Technologies Used

* **Go (Golang)** - Main programming language
* **Gin** - Web framework
* **PostgreSQL** - Database (via pgx/v5)
* **Redis** - Caching, token blacklist
* **Docker** - Containerization and deployment

---

## ✅ Features

* **Authentication**

  * User Registration
  * Login
  * Logout (with token blacklist)
* **User Profile**

  * Edit profile (name, avatar, bio)
  * Follow/unfollow users
* **Posts**

  * Create post (text, image, or both)
  * Upload multiple images in one post
  * Like and comment on posts
* **Feed**

  * View posts from followed users (sorted by newest first)
* **Notifications**

  * Receive notifications for follows, likes, and comments


---

## 🚀 Installation

### Prerequisites

* Go 1.25
* PostgreSQL
* Redis
* Docker & Docker Compose (for containerized deployment)

### Environment Variables

Create a `.env` file in the root directory:

```env
# PostgreSQL config
DB_USER=your_user
DB_PASS=your_pass
DB_HOST=pg-db
DB_PORT=5432
DB_NAME=your_db

# Redis config
RDBHOST=rdb
RDBPORT=6379

# JWT
JWT_SECRET=a-string-secret-at-least-256-bits-long
JWT_ISSUER=your_issue

# Compose overrides
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_pass
POSTGRES_DB=your_db

DB_HOST_MAKE=localhost
DB_PORT_MAKE=5422
```

---

### Setup Instructions

#### Option 1: Docker Deployment (Recommended)

1. Clone the repository

```bash
git clone https://github.com/yourusername/social-media-backend.git
cd social-media-backend
```

2. Create your `prod.env` file with the environment variables above

3. Start with Docker Compose

```bash
docker compose pull
docker compose up -d
```

4. Run database migrations

```bash
make -f Makefile.prod migrate-up
make -f Makefile.prod insert-seed
```

The application will be available at:

* **Backend API**: `http://localhost:8080`
* **PostgreSQL**: `localhost:5422`
* **Redis**: `localhost:6369`

#### Option 2: Local Development

```bash
git clone https://github.com/yourusername/social-media-backend.git
cd social-media-backend
```

1. Create `.env` file
2. Install dependencies

```bash
go mod download
```

3. Run migrations

```bash
make -f Makefile migrate-up
make -f Makefile insert-seed
```

4. Start the application

```bash
go run main.go
```

Server will run on `http://localhost:8080`.

---

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint         | Description   | Auth Required |
| ------ | ---------------- | ------------- | ------------- |
| POST   | `/auth/register` | Register user | ❌             |
| POST   | `/auth/login`    | User login    | ❌             |
| DELETE | `/auth/logout`   | User logout   | ✅             |

### User Endpoints

| Method | Endpoint                 | Description   | Auth Required |
| ------ | ------------------------ | ------------- | ------------- |
| PATCH  | `/user`                  | Edit profile  | ✅             |
| POST   | `/user/:targetID/follow` | Follow a user | ✅             |

### Post Endpoints

| Method | Endpoint        | Description       | Auth Required |
| ------ | --------------- | ----------------- | ------------- |
| POST   | `/post`         | Create a new post | ✅             |
| POST   | `/post/like`    | Like a post       | ✅             |
| POST   | `/post/comment` | Comment on a post | ✅             |

### Feed Endpoints

| Method | Endpoint | Description                                  | Auth Required |
| ------ | -------- | -------------------------------------------- | ------------- |
| GET    | `/feed`  | Get posts from followed users (newest first) | ✅             |

### Static Files

* Images are served under `/api/v1/img/*`

---

## 🔐 Authentication

All protected endpoints use JWT authentication.
Include the token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

---

## 🐳 Docker Architecture

The application runs via Docker Compose with:

* **pg-db**: PostgreSQL database (Port: 5422)
* **rdb**: Redis cache (Port: 6369)
* **backend**: Go API service (Port: 8080)

All services share a dedicated Docker network with persistent data volumes.

---

## 📝 Version History

### Version 1.0.0 (Current)

* Initial release
* User registration & authentication
* Profile management
* Posting (with images)
* Feed & interactions

---

## 👥 Contributors

* Radif - [@radifan9](https://github.com/radifan9)


