# Project Structure

## Directory Organization

```
/Users/mac/workplace/go1.24/my/aichat/
├── .spec-workflow/              # Specification workflow files
│   ├── approvals/               # Approval requests
│   ├── archive/                 # Archived specifications
│   ├── config.example.toml      # Configuration example
│   ├── specs/                   # Specification documents
│   ├── steering/                # Steering documents
│   ├── templates/               # Specification templates
│   └── user-templates/          # User-defined templates
├── backend/                     # Backend source code
│   ├── api/                     # API definitions (protobuf files)
│   ├── cmd/                     # Main applications
│   │   └── server/              # Server application
│   ├── configs/                 # Configuration files
│   ├── internal/                # Private application and library code
│   │   ├── biz/                 # Business logic
│   │   ├── data/                # Data access layer
│   │   ├── server/              # HTTP and gRPC server
│   │   └── service/             # Service layer
│   ├── pkg/                     # Shared libraries and utilities
│   │   ├── agent/               # Agent capabilities
│   │   ├── auth/                # Authentication utilities
│   │   ├── db/                  # Database utilities
│   │   ├── initdata/            # Initial data
│   │   ├── mcp/                 # MCP service integration
│   │   └── redis/               # Redis utilities
│   ├── go.mod                   # Go module definition
│   └── go.sum                   # Go module checksums
├── frontend/                    # Frontend source code
│   ├── public/                  # Static assets
│   ├── src/                     # Source code
│   │   ├── api/                 # API clients
│   │   ├── assets/              # Static assets
│   │   ├── components/          # Reusable components
│   │   ├── plugins/             # Plugin integrations
│   │   ├── router/              # Routing configuration
│   │   ├── stores/              # State management
│   │   ├── utils/               # Utility functions
│   │   ├── views/               # Page components
│   │   ├── App.vue              # Root component
│   │   └── main.ts              # Entry point
│   ├── package.json             # Node.js dependencies
│   └── vite.config.ts           # Vite configuration
├── doc/                         # Documentation
├── script/                      # Scripts for development and deployment
└── README.md                    # Project overview
```

## Naming Conventions

### Files
- **Components/Modules**: `PascalCase` (e.g., `UserManagementView.vue`)
- **Services/Handlers**: `camelCase` (e.g., `userService.ts`)
- **Utilities/Helpers**: `camelCase` (e.g., `dateUtils.ts`)
- **Tests**: `[filename].test.ts` (e.g., `userService.test.ts`)

### Code
- **Classes/Types**: `PascalCase` (e.g., `UserModel`)
- **Functions/Methods**: `camelCase` (e.g., `getUserById`)
- **Constants**: `UPPER_SNAKE_CASE` (e.g., `MAX_RETRY_COUNT`)
- **Variables**: `camelCase` (e.g., `userName`)

## Import Patterns

### Import Order
1. External dependencies (Vue, AntV, etc.)
2. Internal modules (components, services, utilities)
3. Relative imports (siblings, parents)
4. Style imports (CSS, SCSS)

### Module/Package Organization
- Absolute imports from project root when possible
- Relative imports within modules
- Clear separation between frontend and backend modules
- Shared utilities in respective `utils` directories

## Code Structure Patterns

### Module/Class Organization
```
1. Imports/includes/dependencies
2. Constants and configuration
3. Type/interface definitions
4. Main implementation
5. Helper/utility functions
6. Exports/public API
```

### Function/Method Organization
```
1. Input validation first
2. Core logic in the middle
3. Error handling throughout
4. Clear return points
```

### File Organization Principles
- One component/module per file
- Related functionality grouped together
- Public API clearly exported
- Implementation details encapsulated

## Code Organization Principles

1. **Single Responsibility**: Each file should have one clear purpose
2. **Modularity**: Code should be organized into reusable modules
3. **Testability**: Structure code to be easily testable
4. **Consistency**: Follow patterns established in the codebase

## Module Boundaries

- **Core vs Extensions**: Core functionality vs extensible plugins
- **Public API vs Internal**: What's exposed vs implementation details  
- **Backend vs Frontend**: Clear separation of concerns between client and server
- **Business Logic vs Data Access**: Separation between business rules and data operations
- **Dependencies direction**: Higher-level modules depend on lower-level modules, not vice versa

## Code Size Guidelines

- **File size**: Maximum 500 lines per file
- **Function/Method size**: Maximum 50 lines per function
- **Class/Module complexity**: Maximum 10 public methods per class
- **Nesting depth**: Maximum 4 levels of nesting

## Dashboard/Monitoring Structure

```
frontend/
└── src/                     # Self-contained frontend subsystem
    ├── components/          # Reusable UI components using AntV
    ├── views/               # Page components
    ├── router/              # Vue Router configuration
    ├── stores/              # State management (Pinia)
    ├── api/                 # API clients
    └── utils/               # Frontend utilities
```

### Separation of Concerns
- Frontend isolated from backend implementation details
- Own build process with Vite
- Minimal dependencies on backend-specific code
- Can be developed and tested independently

## Documentation Standards
- All public APIs must have documentation
- Complex logic should include inline comments
- README files for major modules
- Follow JSDoc/GoDoc documentation conventions