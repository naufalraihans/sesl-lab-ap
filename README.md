# Laboratory Algorithm and Programming Management System (Web Lab-AP v3)

A comprehensive, state-of-the-art web platform developed to manage laboratory practicum sessions. This system facilitates secure examination environments, live coding capabilities, automated grading utilizing Large Language Models (LLM), and robust administrative dashboards for academic management.

## Table of Contents
1. [Architecture & Technology Stack](#architecture--technology-stack)
2. [Key Features](#key-features)
3. [Project Structure](#project-structure)
4. [Prerequisites](#prerequisites)
5. [Setup and Installation](#setup-and-installation)
6. [API Documentation (OpenAPI/Swagger)](#api-documentation-openapiswagger)
7. [Make Commands Reference](#make-commands-reference)

---

## Architecture & Technology Stack

The project adopts a Clean Architecture pattern in the backend to ensure a strict separation of concerns, dividing the codebase into Domain Entity, DTO, Repository, Usecase, and Delivery (HTTP Handlers) layers.

**Backend:**
- Language: Go (Golang)
- Database: PostgreSQL
- ORM: GORM (Implementing Soft Deletes)
- Authentication: JWT (JSON Web Tokens)
- File Storage: Supabase Storage
- API Documentation: Swagger (swaggo/swag)

**Frontend:**
- Framework: SvelteKit
- Language: TypeScript
- Styling: Tailwind CSS
- Editor Integrations: Monaco Editor (Live Coding), Edra Editor (WYSIWYG)

---

## Key Features

### 1. Robust Examination System
- Server-authoritative countdown timers to prevent client-side manipulation.
- Background "Auto-Save" mechanism to preserve student answers continuously.
- "Auto-Submit" functionality executed by background workers when timers expire or access is abruptly closed by administrators.
- Supports Pre-test (Essay), Post-test (Hybrid), and Practical Exams (Live Coding).

### 2. Live Coding & WYSIWYG Integration
- Embedded Monaco Editor for syntax-highlighted live coding assessments.
- Edra Rich Text Editor implementation for administrators to compose complex questions, featuring KaTeX support for mathematical equations and embedded images.

### 3. AI-Assisted Bulk Grading
- Automated assessment of essay and coding answers utilizing an LLM backend (Ollama).
- Processing is offloaded to background queues (workers) to prevent HTTP timeouts.
- Real-time polling interface on the frontend for grading progress tracking.

### 4. Advanced Data Management & Dashboards
- "Rekap Nilai" (Score Summary): A dynamic Pivot Table dashboard for calculating and aggregating final student grades.
- "Rekap Jawaban" (Global Answers Summary): A comprehensive data grid with filtering capabilities and bulk actions (Reset/Delete) to manage student submissions.
- Bulk Student Data Import via CSV with precise row-level error handling.

### 5. Security & Access Control
- Strict Role-Based Access Control (RBAC) middleware separating Admin (Laboratory Assistants) and Users (Students).
- Cryptographic PIN/Token Gate system requiring students to input dynamic tokens distributed inside the physical laboratory room.
- Soft Delete implementation on crucial tables to prevent catastrophic accidental data loss.

---

## Project Structure

The project is divided into two main components:

- `backend/`: Contains the Go server codebase, structured by clean architecture principles. Database migrations, OpenAPI definitions, and core logic reside here.
- `frontend/`: Contains the SvelteKit frontend codebase.
- `Makefile`: Root level automation script.

---

## Prerequisites

- Go 1.21 or higher
- Node.js v18 or higher
- PostgreSQL
- Ollama (Optional, required only if AI Grading is utilized locally)

---

## Setup and Installation

### 1. Environment Variables Configuration
Clone the provided `.env.example` file (if available) to `.env` in both the `backend/` and `frontend/` directories, and populate the required database credentials, Supabase keys, and JWT secrets.

### 2. Database Migration & Seeding
From the root directory, initialize your database schema and seed the initial administrative data:
```bash
make migrate-up
make seed
```

### 3. Running the Application
To run the backend server:
```bash
make run
```
The backend will run on port `8080` by default.

To run the frontend development server:
```bash
make fe-install
make fe-dev
```
The frontend will run on port `5173` by default.

---

## API Documentation (OpenAPI/Swagger)

The backend exposes an interactive OpenAPI (Swagger) interface for all registered endpoints, which serves as the primary contract between the backend and frontend.

To access the documentation:
1. Ensure the backend server is running (`make run`).
2. Navigate your browser to: `http://localhost:8080/swagger/index.html`

**Updating the Documentation:**
If you make changes to the Go HTTP handlers, you must regenerate the Swagger files before the UI will reflect your changes. Run the following command from the root directory:
```bash
make swag
```

---

## Make Commands Reference

The root `Makefile` provides several utilities to speed up development:

- `make run`: Starts the backend server.
- `make build`: Compiles the backend into a binary executable.
- `make migrate-up`: Executes all pending database migrations.
- `make migrate-down`: Rolls back the last applied database migration.
- `make seed`: Seeds the database with preliminary test data.
- `make swag`: Generates and updates OpenAPI documentation.
- `make tidy`: Executes `go mod tidy` in the backend directory.
- `make fe-install`: Installs all required NPM dependencies for the frontend.
- `make fe-dev`: Starts the SvelteKit development server.
- `make fe-build`: Compiles the frontend for production deployment.
