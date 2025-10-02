# Technology Stack

## Project Type
Web application with backend API service and frontend dashboard

## Core Technologies

### Primary Language(s)
- **Backend Language**: Go 1.24 with Kratos framework
- **Frontend Language**: TypeScript with Vue3
- **Runtime/Compiler**: Go compiler, Node.js runtime
- **Language-specific tools**: Go modules, npm, Vite

### Key Dependencies/Libraries
**Backend Dependencies:**
- **Kratos**: A microservice framework for Go applications
- **GORM**: ORM library for database operations
- **PostgreSQL**: Primary relational database
- **Redis**: In-memory data structure store for caching
- **RocketMQ**: Distributed messaging and streaming platform

**Frontend Dependencies:**
- **Vue3**: Progressive JavaScript framework for building user interfaces
- **TypeScript**: Typed superset of JavaScript
- **AntV**: Data visualization solution
- **Tailwind CSS**: Utility-first CSS framework
- **Vite**: Next generation frontend tooling

### Application Architecture
**Backend Architecture:**
- Microservice-oriented design using Kratos framework
- Clean architecture with separation of concerns
- RESTful API design for frontend communication
- Event-driven architecture with RocketMQ for asynchronous processing

**Frontend Architecture:**
- Component-based architecture with Vue3
- State management with Pinia or Vuex
- Responsive design with Tailwind CSS
- Data visualization with AntV

### Data Storage
- **Primary storage**: PostgreSQL for structured data persistence
- **Caching**: Redis for session management and performance optimization
- **Data formats**: JSON for API communication, Protocol Buffers for internal service communication

### External Integrations
- **MCP Services**: Bridge to third-party interfaces
- **Protocols**: HTTP/REST for API communication, gRPC for internal service communication
- **Authentication**: JWT-based authentication with OAuth2 support

### Monitoring & Dashboard Technologies
- **Dashboard Framework**: Vue3 with AntV for data visualization
- **Real-time Communication**: WebSocket for live updates
- **Visualization Libraries**: AntV for charts and graphs
- **State Management**: Pinia for frontend state management

## Development Environment

### Build & Development Tools
- **Backend Build System**: Go modules with Makefile
- **Frontend Build System**: Vite with npm scripts
- **Package Management**: Go modules for backend, npm for frontend
- **Development workflow**: Hot reload for frontend, live reloading for backend

### Code Quality Tools
- **Static Analysis**: golint, go vet for Go; ESLint, TSLint for TypeScript
- **Formatting**: gofmt for Go; Prettier for TypeScript
- **Testing Framework**: Go testing package for backend; Jest, Cypress for frontend
- **Documentation**: Swagger for API documentation

### Version Control & Collaboration
- **VCS**: Git with GitHub
- **Branching Strategy**: Git Flow
- **Code Review Process**: Pull requests with mandatory reviews

### Dashboard Development
- **Live Reload**: Vite's hot module replacement
- **Port Management**: Configurable ports via environment variables
- **Multi-Instance Support**: Docker-based deployment for multiple instances

## Deployment & Distribution
- **Target Platform(s)**: Cloud-based deployment (Docker containers)
- **Distribution Method**: Containerized deployment with Docker
- **Installation Requirements**: Docker, Docker Compose
- **Update Mechanism**: CI/CD pipeline with GitHub Actions

## Technical Requirements & Constraints

### Performance Requirements
- API response time < 200ms for 95% of requests
- Support 1000 concurrent users
- Database query response time < 100ms for simple queries

### Compatibility Requirements  
- **Platform Support**: Linux, macOS, Windows (development)
- **Browser Support**: Latest versions of Chrome, Firefox, Safari, Edge
- **Go Version**: 1.24+
- **Node.js Version**: 16+

### Security & Compliance
- **Security Requirements**: HTTPS, JWT authentication, input validation
- **Compliance Standards**: GDPR compliance for data protection
- **Threat Model**: Protection against SQL injection, XSS, CSRF attacks

### Scalability & Reliability
- **Expected Load**: 10,000 daily active users
- **Availability Requirements**: 99.9% uptime
- **Growth Projections**: Support for 100,000 users with horizontal scaling

## Technical Decisions & Rationale

### Decision Log
1. **Go with Kratos Framework**: Chosen for its microservice capabilities, strong typing, and performance. Alternatives considered were Java Spring Boot and Python FastAPI.
2. **Vue3 with TypeScript**: Selected for its reactive programming model, component-based architecture, and type safety. Alternatives considered were React and Angular.
3. **PostgreSQL**: Chosen for its reliability, advanced features, and ACID compliance. Alternatives considered were MySQL and MongoDB.
4. **Redis**: Selected for its performance and versatility as a caching layer. Alternatives considered were Memcached.
5. **RocketMQ**: Chosen for its distributed messaging capabilities and scalability. Alternatives considered were Apache Kafka and RabbitMQ.
6. **AntV**: Selected for its powerful data visualization capabilities and integration with Vue. Alternatives considered were Chart.js and D3.js.

## Known Limitations
- Initial learning curve for developers unfamiliar with Kratos framework
- Potential complexity in managing distributed messaging with RocketMQ
- Need for careful database schema design to optimize PostgreSQL performance