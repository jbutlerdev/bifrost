# Bifrost Project Context

## Project Overview

Bifrost is a high-performance AI gateway that unifies access to 15+ AI providers (OpenAI, Anthropic, AWS Bedrock, Google Vertex, Azure, Cerebras, Cohere, Mistral, Ollama, Groq, etc.) through a single OpenAI-compatible API. It provides automatic failover, load balancing, semantic caching, and enterprise-grade features.

Key features include:
- Unified Interface: Single OpenAI-compatible API for all providers
- Multi-Provider Support: Connect to 15+ AI providers
- Automatic Fallbacks: Seamless failover between providers and models
- Load Balancing: Intelligent request distribution across multiple API keys
- Semantic Caching: Intelligent response caching based on semantic similarity
- Model Context Protocol (MCP): Enable AI models to use external tools
- Multimodal Support: Support for text, images, audio, and streaming
- Custom Plugins: Extensible middleware architecture
- Governance: Usage tracking, rate limiting, and fine-grained access control
- Budget Management: Hierarchical cost control with virtual keys
- SSO Integration: Google and GitHub authentication support
- Observability: Native Prometheus metrics, distributed tracing, and logging

## Architecture

Bifrost uses a modular architecture:

```
bifrost/
├── npx/                 # NPX script for easy installation
├── core/                # Core functionality and shared components
│   ├── providers/       # Provider-specific implementations
│   ├── schemas/         # Interfaces and structs used throughout Bifrost
│   └── bifrost.go       # Main Bifrost implementation
├── framework/           # Framework components for data persistence
│   ├── configstore/     # Configuration storages
│   ├── logstore/        # Request logging storages
│   └── vectorstore/     # Vector storages
├── transports/          # HTTP gateway and other interface layers
│   └── bifrost-http/    # HTTP transport implementation
├── ui/                  # Web interface for HTTP gateway
├── plugins/             # Extensible plugin system
│   ├── governance/      # Budget management and access control
│   ├── logging/         # Request logging and analytics
│   ├── semanticcache/   # Intelligent response caching
│   └── telemetry/       # Monitoring and observability
├── docs/                # Documentation and guides
└── tests/               # Comprehensive test suites
```

## Core Components

### Core Module
The core module (`core/`) contains the main Bifrost implementation:
- `Bifrost` struct: Main orchestrator managing providers, plugins, and request processing
- Provider implementations for all supported AI services
- Request routing, queuing, and concurrency management
- Plugin pipeline system with PreHook/PostHook architecture

### Framework
The framework (`framework/`) provides infrastructure components:
- `configstore`: Persistent configuration management with SQLite/PostgreSQL support
- `logstore`: Request logging storage
- `vectorstore`: Vector storage for semantic caching

### Transports
Transports (`transports/`) expose Bifrost functionality:
- `bifrost-http`: High-performance HTTP server using FastHTTP
- OpenAI-compatible REST API endpoints
- Provider-specific endpoints for drop-in replacements
- WebSocket support for streaming
- Built-in web UI

### Plugins
Plugins (`plugins/`) extend Bifrost functionality:
- `governance`: Budget management, virtual keys, access control
- `semanticcache`: Semantic response caching using embeddings
- `logging`: Request/response logging and analytics
- `telemetry`: Monitoring and observability integration

## Development Setup

### Prerequisites
- Go 1.24+
- Node.js 18+ (for UI development)
- Docker (for containerized deployment)

### Quick Start
```bash
# Install and run locally via NPX
npx -y @maximhq/bifrost

# Or use Docker
docker run -p 8080:8080 maximhq/bifrost
```

### Local Development
```bash
# Start complete development environment (UI + API)
make dev

# Build the project
make build

# Run tests
make test
```

## Building and Running

### Make Commands
The project uses a comprehensive Makefile with these key commands:

- `make dev`: Start complete development environment (UI + API with proxy)
- `make build`: Build bifrost-http binary
- `make build-ui`: Build UI
- `make run`: Build and run bifrost-http (no hot reload)
- `make test`: Run tests for bifrost-http
- `make test-core`: Run core tests
- `make test-plugins`: Run plugin tests
- `make test-all`: Run all tests
- `make install-ui`: Install UI dependencies
- `make setup-workspace`: Set up Go workspace with all local modules

### Environment Variables
- `HOST`: Server host (default: localhost)
- `PORT`: Server port (default: 8080)
- `LOG_STYLE`: Logger output format: json|pretty (default: json)
- `LOG_LEVEL`: Logger level: debug|info|warn|error (default: info)
- `APP_DIR`: App data directory inside container (default: /app/data)

## Testing

Bifrost has comprehensive test coverage:
- Unit tests for core functionality
- Provider-specific integration tests
- Plugin tests
- End-to-end tests

Run tests with:
```bash
# Run all tests
make test-all

# Run specific test suites
make test-core
make test-plugins
make test
```

## Deployment Options

### 1. Gateway (HTTP API)
Best for language-agnostic integration, microservices, and production deployments:
```bash
# NPX - Get started in 30 seconds
npx -y @maximhq/bifrost

# Docker - Production ready
docker run -p 8080:8080 -v $(pwd)/data:/app/data maximhq/bifrost
```

### 2. Go SDK
Best for direct Go integration with maximum performance and control:
```bash
go get github.com/maximhq/bifrost/core
```

### 3. Drop-in Replacement
Best for migrating existing applications with zero code changes by replacing provider endpoints.

## Configuration

Bifrost supports multiple configuration methods:
- Web UI for visual configuration
- File-based configuration (config.json)
- Environment variables
- API-driven configuration
- Persistent storage via SQLite/PostgreSQL

## Plugin System

Bifrost has a powerful plugin system that allows extending functionality:
- PreHook/PostHook architecture
- TransportInterceptor for HTTP-level modifications
- Short-circuit capabilities
- Error handling and recovery mechanisms

Plugins can be loaded dynamically or statically and include:
- Semantic caching
- Governance and budget management
- Logging and analytics
- Telemetry and monitoring

## Performance

Bifrost is designed for high performance:
- Virtually zero overhead (<15 µs per request)
- Efficient queuing and request handling
- Memory pooling for reduced allocations
- Concurrent processing with goroutines
- Fast HTTP server using FastHTTP

## Contributing

See the contributing guide for:
- Setting up the development environment
- Code conventions and best practices
- How to submit pull requests
- Building and testing locally