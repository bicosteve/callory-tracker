# 🥗 Callory Tracker

[![Go Build & Test](https://github.com/bicosteve/callory-tracker/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/bicosteve/callory-tracker/actions)

> A lightweight, highly performant, server-rendered web application for tracking what you eat and understanding your daily nutrition — built in Go.

---

## 🔗 Live Demo

Experience the live demo of Callory Tracker here:
👉 **[Callory Tracker Live Demo](https://youtu.be/7849vsrTUk4)**

---

## 📸 Demo Preview

| Dashboard / Running Daily Summary                                                                                          | Logging a New Meal Entry                                                                                            |
| -------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| ![Dashboard Preview](https://raw.githubusercontent.com/bicosteve/callory-tracker/main/ui/static/dashboard_placeholder.png) | ![Add Food Preview](https://raw.githubusercontent.com/bicosteve/callory-tracker/main/ui/static/add_placeholder.png) |

_(Note: Real screenshots can be committed to the `ui/static/` directory to display them live in your repository)._

---

## 🎯 The Problem It Solves

Most people have no clear idea of how much they actually eat in a day. Calories and macronutrients (protein, carbohydrates, and fat) silently add up across breakfast, lunch, dinner, and snacks — and by the time the effects show up, it's hard to know which habits to change.

Existing calorie-tracking apps are often:

- **Bloated** — packed with intrusive ads, heavy social feeds, and premium paywalls just to log a simple meal.
- **Privacy-unfriendly** — your deeply personal eating and health habits become marketing data for third parties.
- **Overly complex** — requiring slow barcode scanners and unreliable food databases for something that should be quick and manual.

**Callory Tracker** solves this with a deliberately simple approach: a fast, secure, self-hostable web app where you log each meal manually with its nutritional values, and instantly see your running daily totals. You own your data, the app stays completely out of your way, and the entire stack is small enough to run on a single inexpensive micro VM.

---

## 🚀 Features

- 🔐 **User Authentication** — Secure registration, login, and logout flow.
- 🍽 **Meal Management (CRUD)** — Create, read, update, and delete food entries easily.
- 📊 **Nutrition Analysis** — Automatic aggregation of total daily calories and macronutrients (carbs, protein, fats).
- 🛡 **Security-first Middleware** — CSRF protection, secure HTTP headers, panic recovery, and request logging.
- 🩺 **Health Endpoint** — A `/health` route for uptime checks, monitoring, and load balancer ping checks.
- 🧪 **Test Suite** — Fully test-covered with robust unit tests, helper validations, and mock-backed handler tests.

---

## 🧰 Tech Stack

| Layer      | Technology                                        |
| ---------- | ------------------------------------------------- |
| Language   | Go (Golang)                                       |
| Router     | [chi](https://github.com/go-chi/chi)              |
| Middleware | [alice](https://github.com/justinas/alice)        |
| CSRF       | [nosurf](https://github.com/justinas/nosurf)      |
| Sessions   | golangcollege/sessions                            |
| Database   | MySQL                                             |
| Templates  | Go `html/template` (embedded via `embed.FS`)      |
| Frontend   | Server-rendered HTML + CSS                        |
| Deployment | Docker, GitHub Actions CI/CD, Contabo VM / Heroku |

---

## 🛡️ Security Notes & Authentication Process

Callory Tracker utilizes industry-standard security patterns to guard user accounts and preserve data privacy:

1. **Password Hashing**: User credentials are protected using the **bcrypt** hashing algorithm with a cost factor of 12, ensuring slow, computationally-expensive defense against brute-force attacks.
2. **Session Management**: Session authentication is conducted through server-side sessions using cookie jars (via `golangcollege/sessions`). Cookies are configured with `HttpOnly` and `SameSite=Strict` (configured as `Secure: false` by default for local development over HTTP).
3. **CSRF Protection**: A global **CSRF (Cross-Site Request Forgery)** protection middleware (`nosurf`) wraps all state-changing endpoints. It injects a unique token in each page form and enforces validating incoming POST payloads.
4. **Secure HTTP Headers**: Custom headers are injected on every response:

- `X-Frame-Options: deny` (Blocks clickjacking)
- `X-XSS-Protection: 1; mode=block` (Mitigates cross-site scripting)
- `Content-Type: text/html; charset=utf-8` (Ensures explicit character encoding)

1. **Contextual Scoping**: Handlers retrieve the active user context injected by authentication middleware. Every query for food logging, editing, and summing is tightly scoped using the authenticated User ID, making sure users can **never** view or alter other users' data.

---

## 🏗 Architecture Overview

Callory Tracker follows a clean, idiomatic Go web-app layout that separates HTTP concerns from data access:

- **`cmd/web`** — the application entry point and HTTP layer: routing, handlers, middleware, template caching, and helpers.
- **`pkg/models`** — domain types (`User`, `Food`) and interfaces (`UserModelInterface`, `FoodModelInterface`) that define the data-access contract. This interface-driven design is what makes the handlers testable with mock implementations.
- **`pkg/models/mysql`** — the production MySQL implementation of those interfaces.
- **`pkg/models/mock`** — in-memory mock implementations used by the test suite.
- **`pkg/forms`** — reusable form parsing and validation (required fields, lengths, email format, password matching).
- **`pkg/helpers`, `pkg/configs`, `pkg/db`, `pkg/logger`** — supporting utilities for config loading, DB connection pooling, and structured logging.
- **`ui`** — HTML templates and CSS, embedded directly into the binary via `embed.FS` so the app ships as a single self-contained executable.

Every request passes through a standard middleware chain:

```
authenticate → recoverPanic → logRequest → secureHeaders → noSurf (CSRF)
```

Protected routes additionally pass through `requireAuthenticatedUser`, which redirects unauthenticated visitors to the login page.

---

## 🌐 Routes Overview

| Method | Path             | Auth Required | Description                                       |
| ------ | ---------------- | ------------- | ------------------------------------------------- |
| GET    | `/`              | No            | Home page / landing                               |
| GET    | `/health`        | No            | Health check for uptime monitors & load balancers |
| GET    | `/user/register` | No            | Registration form                                 |
| POST   | `/user/register` | No            | Register and create a new user account            |
| GET    | `/user/login`    | No            | Login form                                        |
| POST   | `/user/login`    | No            | Authenticate user and start session               |
| POST   | `/user/logout`   | No            | Terminate session and logout                      |
| GET    | `/user/me`       | Yes           | View profile details of currently logged-in user  |
| GET    | `/food/add`      | Yes           | Add new food item entry form                      |
| POST   | `/food/add`      | Yes           | Create a food entry                               |
| GET    | `/food/day`      | Yes           | View a specific food entry detail                 |
| GET    | `/food/get-edit` | Yes           | Get form to edit an existing food entry           |
| POST   | `/food/edit`     | Yes           | Update a food entry                               |
| POST   | `/food/delete`   | Yes           | Delete a food entry                               |
| POST   | `/food/total`    | Yes           | Get daily nutrition aggregate summary             |

---

## 📦 Project Structure

```bash
callory-tracker/
├── cmd/
│ └── web/ # HTTP layer
│ ├── main.go # Entry point: config, DB, sessions, server
│ ├── routes.go # Route + middleware wiring (chi + alice)
│ ├── handlers.go # Request handlers
│ ├── middleware.go # recoverPanic, logRequest, secureHeaders, auth, CSRF
│ ├── helper.go # serverError/clientError, template rendering
│ └── templates.go # Template cache + dynamic data
├── pkg/
│ ├── configs/ # Environment/config loading
│ ├── db/ # MySQL connection pool
│ ├── forms/ # Form parsing & validation
│ ├── helpers/ # Shared helper utilities
│ ├── logger/ # Structured logging
│ └── models/ # Domain types + interfaces
│ ├── mysql/ # Production MySQL implementation
│ └── mock/ # Mock implementations for tests
├── tables/ # SQL schema (Users.sql, Foods.sql)
├── ui/
│ ├── css/ # Stylesheets
│ ├── html/ # Templates (pages + partials)
│ └── efs.go # embed.FS for templates & static assets
├── Dockerfile
├── docker-compose.prod.yml
├── Makefile
├── Procfile
├── go.mod / go.sum
└── README.md
```

---

## 🛡️ CI/CD Status

Continuous Integration & Deployment is set up via **GitHub Actions** (`.github/workflows/ci-cd.yml`).

- **Test** (on push to `feat/*`, `fix/*`, `main` and PRs to `main`): provisions a clean Ubuntu runner, installs Go 1.20, downloads dependencies, runs `go vet ./...`, then runs the full suite with the race detector (`go test -race -vet=off ./...`).
- **Build & Push** (on push to `main` only, after tests pass): builds the Docker image with Buildx and pushes it to Docker Hub tagged with both the commit SHA and `latest`.
- **Deploy** (manual `workflow_dispatch` on `main`): copies `docker-compose.prod.yml` to the VM over SSH and pulls/restarts the stack.

---

## 🧪 Testing (How to Run Tests)

The project ships with a robust test suite covering form validation, utility helpers, custom HTTP middleware, and REST handler flows using mock database models. No external database connection is required to run the test suite.

Run tests using standard Go commands:

```bash
# Run all tests in the project
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests and evaluate test coverage
go test ./... -cover

# Run a specific test inside a package
go test -run=TestGetDay ./cmd/web -v
```

---

## 🚧 Known Limitations & Future Improvements

While functional, simple, and light, the following features are scope for future enhancement:

1. **Session Store Persistence**: Currently uses an in-memory session backing store. If the web server is restarted, all logged-in user sessions are cleared.

- _Improvement_: Transition sessions backend to **Redis** or a schema table in **MySQL**.

1. **Interactive Charting**: Daily macros are calculated and printed as static server-side HTML tables.

- _Improvement_: Add dynamic radial/bar charts utilizing **Chart.js** or light CSS progress meters.

1. **Nutritional Lookup Database**: Users must know and type the macros manually for every dish.

- _Improvement_: Integrate a free public food registry API (such as USDA FoodData or Open Food Facts).

1. **Automated Password Recovery**: User authentication lacks password reset flows (e.g. email reset token generation).

- _Improvement_: Integrate an email gateway service (such as SendGrid or mailgun) to authorize reset tokens.

---

## 🛠️ Local Development Getting Started

### 1. Prerequisites

- Go 1.20+ installed.
- A running MySQL database.

### 2. Configure Environment

Copy the example environment settings and customize:

```bash
cp .env.example .env
```

| Variable     | Description                                      | Example               |
| ------------ | ------------------------------------------------ | --------------------- |
| `DBUSER`     | MySQL username                                   | `root`                |
| `DBPASSWORD` | MySQL password                                   | `secret`              |
| `DBHOST`     | MySQL host                                       | `localhost`           |
| `DBPORT`     | MySQL port                                       | `3306`                |
| `DBNAME`     | Database name                                    | `callory`             |
| `SESSION`    | 32-byte key for session encryption               | `Nrqe6etTZ68Gymwx...` |
| `PORT`       | Port the app listens on (bare number, no colon)  | `4001`                |
| `DBSSLMODE`  | Enable TLS to the database (`true`/`false`)      | `false`               |
| `DBCACERT`   | Path to the CA certificate when `DBSSLMODE=true` | `./ca.pem`            |
| `ENV`        | Runtime environment (`prod` reads real env vars) | `prod`                |

### 3. Initialize Database Tables

```bash
mysql -u <user> -p <db_name> < tables/Users.sql
mysql -u <user> -p <db_name> < tables/Foods.sql
```

### 4. Build and Run

```bash
go mod tidy
go run ./cmd/web
```

Navigate to `http://localhost:4001` to test the application locally.

---

## 🚢 Production Deployment

The binary is fully compatible with containerized environments:

- **Docker**: Build a lightweight container image using the included multi-stage `Dockerfile` (built on Alpine Linux).
- **Docker Compose**: Production multi-container stacks are detailed in `docker-compose.prod.yml`.

```bash
# Build & Run Container
docker build -t callory-tracker .
docker run --env-file .env -p 4001:4001 callory-tracker
```

---

## 📄 License

This project is licensed under the repository license specifications. Feel free to clone, adapt, and self-host!
