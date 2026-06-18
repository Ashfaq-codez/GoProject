# Go User Directory API

A containerized, RESTful Go API demonstrating clean architecture, PostgreSQL integration, and a Neo-Brutalist frontend. 

## Architecture & Tech Stack
* **Backend:** Go (GoFiber)
* **Database:** PostgreSQL (SQLC for type-safe queries)
* **Frontend:** HTML/JS + Tailwind CSS (served via Nginx)
* **Infrastructure:** Fully containerized with Docker & Docker Compose
* **Validation & Logging:** `go-playground/validator` and Uber Zap

## Features
* Explicit Clean Architecture (Router -> Handler -> Service -> Repository)
* Dynamic business logic (Exact age calculation handling leap years)
* Cross-Origin Resource Sharing (CORS) configured
* Backend pagination support (`limit` and `offset`)
* Automated database migrations via Docker entrypoint

## Running Locally

1. Ensure Docker and Docker Compose are installed.
2. Clone the repository:
   ```bash
   git clone [https://github.com/yourusername/goproject.git](https://github.com/yourusername/goproject.git)