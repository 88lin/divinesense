# DivineSense (神识)

**AI-Powered Personal Second Brain** — Automate tasks, filter information, amplify productivity through intelligent agents.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB.svg)](https://react.dev/)

> Forked from [usememos/memos](https://github.com/usememos/memos), enhanced with multi-agent AI system.

---

## Why DivineSense?

| **Efficiency** | **Knowledge** | **AI Agents** | **Privacy** |
|:-------------:|:-------------:|:-------------:|:-----------:|
| Automate tasks | Smart storage | Intent routing | Self-hosted |
| Save time | Semantic search | Multi-agent | Data privacy |

---

## Quick Start

### Docker (Basic Notes)

```bash
docker run -d --name divinesense -p 5230:5230 -v ~/.divinesense:/var/opt/divinesense hrygo/divinesense:stable
```

Access at http://localhost:5230

### Full AI Features (PostgreSQL Required)

```bash
# 1. Clone repository
git clone https://github.com/hrygo/divinesense.git && cd divinesense

# 2. Configure environment
cp .env.example .env
# Edit .env and add your API keys

# 3. Install dependencies
make deps-all

# 4. Start all services (PostgreSQL + Backend + Frontend)
make start
```

Access at http://localhost:25173

<details>
<summary><b>Service Management</b></summary>

```bash
make status   # Check service status
make logs     # View logs
make stop     # Stop services
make restart  # Restart services
```

</details>

---

## Deployment

### Docker Deployment (Simple Notes)

```bash
docker run -d --name divinesense -p 5230:5230 -v ~/.divinesense:/var/opt/divinesense hrygo/divinesense:stable
```

### Binary Deployment (Recommended for Geek Mode)

Binary deployment offers better performance and native Geek Mode support.

```bash
# One-click installation (default: Docker mode)
curl -fsSL https://raw.githubusercontent.com/hrygo/divinesense/main/deploy/aliyun/install.sh | sudo bash

# Binary mode (for Geek Mode)
curl -fsSL https://raw.githubusercontent.com/hrygo/divinesense/main/deploy/aliyun/install.sh | sudo bash -s -- --mode=binary
```

**Advantages:**
- ✅ Native Geek Mode (Claude Code CLI integration)
- ✅ Faster startup, lower overhead
- ✅ Easier updates

**Documentation:** [Binary Deployment Guide](docs/deployment/BINARY_DEPLOYMENT.md)

### Development Setup

```bash
# 1. Clone repository
git clone https://github.com/hrygo/divinesense.git && cd divinesense

# 2. Configure environment
cp .env.example .env
# Edit .env and add your API keys

# 3. Install dependencies
make deps-all

# 4. Start all services (PostgreSQL + Backend + Frontend)
make start
```

Access at http://localhost:25173

<details>
<summary><b>Service Management</b></summary>

```bash
make status   # Check service status
make logs     # View logs
make stop     # Stop services
make restart  # Restart services
```

</details>

<details>
<summary><b>Release Build</b></summary>

```bash
# Build release binaries (linux/amd64, linux/arm64)
make release-build VERSION=v1.0.0

# Package releases
make release-package VERSION=v1.0.0

# Full release workflow
make release-all VERSION=v1.0.0
```

</details>

---

### Note Taking
- Quick capture with Markdown support (KaTeX, Mermaid, GFM)
- Tag-based organization (`#tag`)
- Timeline view with filters
- File attachments (images, documents)
- Semantic search with hybrid BM25 + vector retrieval
- Memo relations and linking

### Schedule Management
- Calendar views (month/week/day/agenda)
- Natural language event creation
- Automatic conflict detection
- Drag-and-drop rescheduling
- Recurring events (daily/weekly/monthly/custom)
- Time zone support

### AI Agents

Three specialized "Parrot" agents with distinct personalities:

| Agent | Name | Purpose | Example |
|:-----:|:-----|:--------|:--------|
| **🦜** | **HuiHui** (灰灰) | Knowledge Retrieval | "What did I write about React?" |
| **🦜** | **JinGang** (金刚) | Schedule Management | "Schedule tomorrow's 3pm meeting" |
| **🦜** | **Amazing** (惊奇) | Comprehensive Assistant | "Summarize my week and upcoming tasks" |

**Smart Routing**:
- Rule-based matching (0ms) for common patterns
- History-aware routing (~10ms) for context
- LLM fallback (~400ms) for ambiguous inputs
- No manual agent selection needed

**Session Memory**:
- Conversation context persists across sessions
- 30-day retention with auto-cleanup
- Per-agent memory isolation

---

## Tech Stack

| Layer | Technology |
|:-----|:----------|
| **Backend** | Go 1.25+, Echo Framework, Connect RPC |
| **Frontend** | React 18, Vite 7, TypeScript, Tailwind CSS 4, Radix UI |
| **Database** | PostgreSQL 16+ (pgvector extension) |
| **AI Models** | DeepSeek V3, Qwen2.5-7B, bge-m3, bge-reranker-v2-m3 |

### Hybrid RAG Retrieval

```
Query → QueryRouter → BM25 + pgvector (HNSW) → Reranker → RRF Fusion
```

| Component | Technology | Purpose |
|:----------|:-----------|:--------|
| **Vector Search** | pgvector + HNSW index | Semantic similarity |
| **Full-Text** | PostgreSQL FTS + BM25 | Keyword matching |
| **Reranker** | BAAI/bge-reranker-v2-m3 | Result refinement |
| **Embedding** | BAAI/bge-m3 (1024d) | Text vectorization |
| **LLM** | DeepSeek V3 / Qwen2.5 | Response generation |

### Agent Architecture

```
ChatRouter (Intent Classification)
    ├── Rule-based (0ms) - keywords, patterns
    ├── History-aware (~10ms) - conversation context
    └── LLM fallback (~400ms) - semantic understanding

Routes to:
    ├── MemoParrot (灰灰) - memo_search tool
    ├── ScheduleParrotV2 (金刚) - schedule_add/query/update/find_free_time
    └── AmazingParrot (惊奇) - concurrent multi-tool orchestration
```

---

## Development

```bash
make start     # Start all services
make stop      # Stop all services
make status    # Check service status
make logs      # View logs
make test      # Run tests
make check-all # Run all checks (build, test, i18n)
```

**Documentation**:
- [Backend & Database](docs/dev-guides/BACKEND_DB.md) - API, DB schema, environment setup
- [Frontend Architecture](docs/dev-guides/FRONTEND.md) - Layouts, Tailwind pitfalls, components
- [System Architecture](docs/dev-guides/ARCHITECTURE.md) - Project structure, AI agents, data flow

---

## AI Database Schema (PostgreSQL)

| Table | Purpose |
|:-----|:--------|
| `memo_embedding` | Vector embeddings for semantic search |
| `conversation_context` | Session persistence for AI agents |
| `episodic_memory` | Long-term user memory and preferences |
| `user_preferences` | User communication settings |
| `agent_metrics` | Agent performance tracking (A/B testing) |

---

## License

[MIT](LICENSE) — Free to use, modify, and distribute.
