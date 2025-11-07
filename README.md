# 🧠 AI Hub

**AI Hub** is a unified gateway microservice for centralized interaction with multiple LLM (Large Language Model) providers.

---

## 🎯 Purpose

AI Hub serves as a **single entry point** for all services interacting with large language models (LLMs) such as:

* OpenAI
* Anthropic (Claude)
* Google Gemini
* AWS Bedrock
* Local models via Ollama

It provides centralized management of:

* API keys and provider routing
* prompt creation and versioning
* project/task lifecycle
* logging, monitoring, and analytics
* cost tracking and access control

---

## 🏗️ Architecture

The backend follows **Clean Architecture** and **Domain-Driven Design (DDD)** principles, ensuring:

* clear separation of concerns,
* independence from frameworks,
* scalability and testability.

---

### 🧩 Folder Structure Overview

```text
backend/
 ├── core/                       # Pure domain and application logic
 │   ├── app/                    # Application layer — orchestrates use-cases
 │   ├── domain/                 # Domain entities, aggregates, and value objects
 │   └── ports/                  # Interfaces for repositories, cache, external APIs
 │
 └── infra/                      # Infrastructure implementations
     ├── adapters/               # Implementations of core ports (bridge to infra)
     ├── cache/redis/            # Redis caching logic and configuration
     ├── composition/containers/ # Dependency Injection and service composition
     ├── config/                 # Environment and runtime configuration
     ├── db/mongoose/            # MongoDB persistence layer
     │   ├── mappers/            # Map DB models to domain entities and back
     │   ├── models/             # Database schemas and ODM/driver models
     │   ├── repos/              # Repository implementations for domain ports
     │   └── uow/                # Unit of Work abstraction for atomic operations
     ├── http/                   # HTTP transport layer
     │   ├── gin/                # Gin controllers and route definitions
     │   └── httperrormapper/    # Maps domain/application errors to HTTP responses
     └── services/               # Integrations with external services
````

---

## 🛠 Tech Stack

| **Component**  | **Technology**                                           |
| -------------- | -------------------------------------------------------- |
| Backend        | Go (Golang, Gin)                                         |
| Frontend       | React                                                    |
| Architecture   | Clean Architecture (hexagonal, ports-and-adapters) + DDD |
| Database       | MongoDB                                                  |
| Cache          | Redis                                                    |
| LLM SDK        | `grokify/gollm`                                          |
| Authentication | JWT                                                      |
| Email Service  | SMTP (Gmail)                                             |

---

## ⚙️ Environment Variables

```properties
PORT=5000

MONGODB_URI=mongodb+srv://<user>:<password>@<cluster>.mongodb.net/Ai-Hub
JWT_SECRET=<jwt_secret>
JWT_EXPIRES_IN=1h

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
EMAIL_USER=<email_address>
EMAIL_PASS=<app_password>

REDIS_URL=redis://127.0.0.1:6379
REDIS_PROJECT_EDIT_TTL=86400
REDIS_PROJECT_CACHE_TTL=3600

TIMEOUT=10000
CRYPTO_ALGORITHM=aes-256-cbc
KEY_ENCRYPT_SECRET=<32_byte_secret_key>
IV_LENGTH=16
```

---

## 🚀 Running Locally

### 1️⃣ Clone the repository

```bash
git clone https://github.com/Dennis761/ai-hub.git
cd ai-hub
```

### 2️⃣ Create `.env` file

Use the example above and place it in the project root.

### 3️⃣ Start the backend

```bash
cd backend
go mod tidy
go run ./infra
```

Backend will run on:

```text
http://localhost:5000
```

### 4️⃣ Start the frontend

```bash
cd frontend
npm install
npm start
```

Frontend will run on:

```text
http://localhost:5173
```

---

## 🧩 Core Features

| **Feature**                       | **Description**                                             |
| --------------------------------- | ----------------------------------------------------------- |
| 🔑 Centralized API Key Management | AES-256 encryption and secure lifecycle control             |
| 🌐 Unified Multi-Provider Gateway | One API for OpenAI, Claude, Gemini, Bedrock, and Ollama     |
| 🧱 Prompt Versioning              | Full history tracking, rollback, and version diffs          |
| 📁 Project & Task Lifecycle       | Manage creation, execution, and completion of LLM workflows |
| ⚡ Redis Caching                   | Session-level caching for fast response times               |
| 🔐 JWT Authentication             | Secure, role-based access and session management            |
| 🧩 Clean Architecture + DDD Core  | Modular, testable, and extensible domain-centered design    |
| 📊 Monitoring & Cost Tracking     | Real-time usage and billing insights per project/task       |

```
