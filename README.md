# 🧠 AI Hub

**AI Hub** is a microservice that serves as a unified gateway for connecting multiple company projects with LLM models (OpenAI, Claude, local models, etc.), providing centralized management of API keys, projects, tasks, prompts, logs, and security.

---

## 🚀 Key Features

- Centralized management of **API keys** from different LLM providers.
- Creation, execution, and history of **prompt-based tasks** (Prompt Pipelines).
- Ability to **test** individual prompts or complete execution flows.
- Drag & Drop **prompt constructor** with placeholders.
- **Version tracking** of prompts with rollback capability.
- Project-level API key association and management.
- Built-in support for **OpenAI**, **Gemini**, and **Llama (local)** models.

---

## 🧩 Architecture

AI Hub uses a microservice architecture allowing scalability for each component:

- **API Key Management**
- **Project & Task Management**
- **Prompt Builder & Runner**
- **Version History**
- **Security & Logging (in development)**

---

## 🔐 Supported Models

| Model     | Example Key   | Notes                         |
|-----------|----------------|-------------------------------|
| OpenAI    | `sk-...`        | Supports `gpt-3.5`, `gpt-4`   |
| Gemini    | `AIza...`       | Provided by Google            |
| Llama     | `sk-...`        | Self-hosted local model       |

---

## ⚙️ Environment Setup

This guide will help you fully set up and run AI Hub — including both frontend and backend.

---

### 🔧 Prerequisites

Make sure the following are installed on your system:

- **Node.js (v18 or later)**
- **npm** (comes with Node.js)
- **Git**
- **MongoDB Atlas** (or local MongoDB)

---

### 🟢 Check Node.js and npm

1. Open terminal and run:
   ```bash
   node -v
   npm -v
   ```

2. If not installed, download from:  
   [https://nodejs.org/](https://nodejs.org/)

---

### 🟤 Check Git

1. Run:
   ```bash
   git --version
   ```

2. If not installed, download from:  
   [https://git-scm.com/downloads](https://git-scm.com/downloads)

---

### 🍃 Prepare MongoDB

1. Register at [https://www.mongodb.com/cloud/atlas](https://www.mongodb.com/cloud/atlas)
2. Create a cluster and obtain a **connection string**
3. Or install MongoDB locally

---

## 🔁 Clone the Project

1. Clone the repository:
   ```bash
   git clone https://github.com/Dennis761/ai-hub.git
   ```

2. Navigate into the project folder:
   ```bash
   cd ai-hub
   ```

---

## 🛠️ Backend Setup

1. Go to the `ai-hub-backend` folder:
   ```bash
   cd ai-hub-backend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Fill in the existing `.env` file:
   ```
   PORT=8080
   MONGODB_URI=your_mongodb_uri
   OPENAI_KEY=your_openai_key
   GEMINI_KEY=your_gemini_key
   LLAMA_KEY=your_llama_key
   ```

4. Start the backend server:
   ```bash
   npm run dev
   ```

> Server will be running at: `http://localhost:8080`

---

## 🌐 Frontend Setup

1. Go to the `frontend` folder:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Create the `.env` file:
   ```
   VITE_API_URL=http://localhost:8080
   ```

4. Start the frontend:
   ```bash
   npm run dev
   ```

> App UI will be available at: `http://localhost:5173`

---

## ✅ Ready to Go

You can now:

- Manage API keys
- Create and run prompt-based task flows
- Interact with OpenAI, Gemini, and Llama models

