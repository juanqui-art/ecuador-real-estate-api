# MCP Tools Usage Guide - Real Estate Platform

## 🧰 Complete MCP Stack Configured

### ✅ Installed Tools
1. **Context7** (@upstash/context7-mcp) - Project intelligence
2. **Sequential** (@modelcontextprotocol/server-sequential-thinking) - Methodology
3. **Magic** (@21st-dev/magic) - UI generation
4. **Puppeteer** (@kirkdeam/puppeteer-mcp-server) - E2E testing
5. **Filesystem** (@modelcontextprotocol/server-filesystem) - File operations
6. **PostgreSQL** (mcp-postgres-server) - Database optimization and queries
7. **OpenAPI** (universal-openapi-mcp) - API documentation and type generation

## 🎯 Role-Based Development Workflows

### 🏗️ **Architect Role** - Design Phase
**Tools**: Context7 + Sequential + Filesystem + PostgreSQL + OpenAPI

**Example Session**:
```
🎯 Task: "Design the dashboard architecture for Next.js integration"

Context7 will provide:
- Complete understanding of Go backend structure
- JWT authentication flow mapping
- Database relationships and constraints
- API endpoint inventory (56+ endpoints)

Sequential will plan:
- Component hierarchy design
- State management strategy
- API integration patterns
- Authentication flow implementation
- Responsive design approach

Filesystem will help:
- Create project structure
- Setup configuration files
- Generate boilerplate code

PostgreSQL will provide:
- Database schema insights
- Query optimization suggestions
- Performance analysis
- Index recommendations

OpenAPI will generate:
- Automatic API documentation
- TypeScript interfaces from Go structs
- Frontend API client code
- Type-safe integration patterns
```

**Expected Deliverables**:
- Dashboard architecture documentation
- Component hierarchy diagram
- API integration strategy
- File structure template

---

### 🎨 **Frontend Role** - Development Phase
**Tools**: Magic + Context7 + Puppeteer + OpenAPI

**Example Session**:
```
🎯 Task: "Build property management interface with search and filters"

Magic will generate:
- PropertyGrid component with shadcn/ui
- PropertyCard with image gallery
- SearchFilters with advanced options
- PropertyForm for creation/editing
- All with TypeScript, Tailwind CSS, accessibility

Context7 will ensure:
- Proper API integration with Go backend
- JWT authentication handling
- Role-based UI permissions
- Type-safe data handling

Puppeteer will create:
- Component interaction tests
- User workflow validations
- Visual regression tests

OpenAPI will provide:
- Auto-generated TypeScript interfaces
- Type-safe API client methods
- Real-time API documentation
- Go→TypeScript type mapping
```

**Expected Deliverables**:
- React components with full functionality
- API integration code
- TypeScript interfaces
- Component tests

---

### ⚙️ **Backend Role** - API Enhancement
**Tools**: Context7 + Sequential + Filesystem + PostgreSQL + OpenAPI

**Example Session**:
```
🎯 Task: "Add real-time notifications for property updates"

Context7 will understand:
- Current JWT authentication system
- Existing WebSocket infrastructure
- Database schema relationships
- User role permissions

Sequential will plan:
- WebSocket server implementation
- Event broadcasting system
- Authentication for real-time connections
- Database trigger setup

Filesystem will help:
- Create new handler files
- Update configuration
- Add migration files

PostgreSQL will provide:
- Database schema optimization
- Query performance analysis
- Real-time connection pooling
- Index recommendations for notifications

OpenAPI will generate:
- WebSocket endpoint documentation
- Real-time event type definitions
- Client SDK updates
- API versioning support
```

**Expected Deliverables**:
- WebSocket endpoints
- Real-time event system
- Database triggers
- Updated API documentation

---

### ✅ **QA Role** - Testing Phase
**Tools**: Puppeteer + Context7 + PostgreSQL

**Example Session**:
```
🎯 Task: "Test complete property management workflow"

Context7 will provide:
- Business logic understanding
- User role requirements
- API endpoint specifications
- Authentication requirements

Puppeteer will execute:
- End-to-end user workflows
- Role-based access testing
- Cross-browser validation
- Performance benchmarks

PostgreSQL will provide:
- Database integrity testing
- Query performance validation
- Connection pooling stress tests
- Data consistency verification
```

**Expected Deliverables**:
- Comprehensive test suites
- User workflow validation
- Performance reports
- Bug identification reports

---

### 🔍 **Analyzer Role** - Debugging & Optimization
**Tools**: All MCP tools available

**Example Session**:
```
🎯 Task: "Optimize property search performance"

Context7: Understand current architecture
Sequential: Plan optimization steps
Magic: Update UI components if needed
Puppeteer: Performance testing
Filesystem: Code modifications
PostgreSQL: Database optimization analysis
OpenAPI: API documentation updates
```

## 🚀 **FASE 2 Development Roadmap**

### Week 1: Foundation Setup
**Role**: Architect
**Tools**: Context7 + Sequential + Filesystem + PostgreSQL + OpenAPI

1. **Day 1-2**: Next.js project setup and architecture design
2. **Day 3-4**: Authentication integration with Go JWT
3. **Day 5**: API client and type generation setup

### Week 2: Core Components
**Role**: Frontend
**Tools**: Magic + Context7 + Puppeteer + OpenAPI

1. **Day 1-2**: Authentication UI (login, logout, token refresh)
2. **Day 3-4**: Property listing and search components
3. **Day 5**: Property creation/editing forms

### Week 3: Advanced Features
**Role**: Frontend + Backend
**Tools**: All MCP tools

1. **Day 1-2**: Image gallery and upload functionality
2. **Day 3-4**: User and agency management interfaces
3. **Day 5**: Dashboard analytics and statistics

### Week 4: Testing & Polish
**Role**: QA + Analyzer
**Tools**: Puppeteer + Context7 + PostgreSQL

1. **Day 1-2**: E2E testing implementation
2. **Day 3-4**: Performance optimization
3. **Day 5**: Final polish and deployment preparation

## 🎯 **Specific Use Cases**

### 1. **Generate Property Card Component**
```
Role: Frontend
Tools: Magic + Context7

Request: "Create a property card component for the dashboard grid"

Magic generates:
- TypeScript component with proper interfaces
- shadcn/ui Card, Badge, Button components
- Tailwind CSS responsive styling
- Image gallery integration
- Price formatting for Ecuador market
- Action buttons with role-based permissions

Context7 ensures:
- Proper integration with Go Property struct
- JWT authentication for actions
- Role-based button visibility
```

### 2. **Setup Authentication Flow**
```
Role: Architect + Frontend
Tools: Context7 + Sequential + Magic + OpenAPI

Request: "Implement JWT authentication flow in Next.js"

Sequential plans:
1. Setup NextAuth.js or custom JWT handling
2. Create login/logout components
3. Implement token refresh logic
4. Add protected route middleware
5. Create user context provider

Magic generates:
- Login form with validation
- User avatar dropdown
- Protected route components

Context7 ensures:
- Proper integration with Go JWT system
- Correct token validation approach

OpenAPI provides:
- Auto-generated TypeScript auth interfaces
- Type-safe authentication API methods
- Token refresh endpoint documentation
- Role-based permission types
```

### 3. **Test Property Management Workflow**
```
Role: QA
Tools: Puppeteer + Context7 + PostgreSQL

Request: "Test complete property CRUD workflow for agency user"

Context7 provides:
- Agency user permissions (can create, update, delete properties)
- Required JWT token format
- API endpoint specifications

Puppeteer executes:
1. Login as agency user
2. Navigate to properties section
3. Create new property with images
4. Edit property details
5. Delete property
6. Verify proper authentication throughout

PostgreSQL provides:
- Database transaction monitoring
- Query performance during testing
- Data consistency validation
- Connection pooling stress testing
```

## 🔧 **Configuration Details**

### Environment Variables Set
```bash
# Context7
PROJECT_ROOT=/Users/juanquizhpi/GolandProjects/realty-core
PROJECT_TYPE=fullstack
BACKEND_LANGUAGE=go
FRONTEND_FRAMEWORK=nextjs
AUTH_METHOD=jwt
DATABASE=postgresql

# Magic
UI_FRAMEWORK=nextjs
COMPONENT_LIBRARY=shadcn-ui
STYLING=tailwindcss
PROJECT_TYPE=dashboard

# Puppeteer
HEADLESS=true
VIEWPORT_WIDTH=1920
VIEWPORT_HEIGHT=1080
DEFAULT_TIMEOUT=30000

# PostgreSQL
DATABASE_URL=postgresql://juanquizhpi@localhost:5433/inmobiliaria_db?sslmode=disable
MAX_CONNECTIONS=10
QUERY_TIMEOUT=30000

# OpenAPI
PROJECT_ROOT=/Users/juanquizhpi/GolandProjects/realty-core
SPEC_PATH=/docs/api
AUTO_GENERATE=true
```

### File Permissions
- **Filesystem MCP**: Access to `/Users/juanquizhpi/GolandProjects/`
- **Context7**: Full project understanding
- **Sequential**: Workspace management
- **PostgreSQL**: Database connection and query access
- **OpenAPI**: API documentation and type generation

## 🎉 **Ready for FASE 2 Development**

The complete 7-tool MCP stack is now configured and ready to accelerate Next.js dashboard development. Each tool is optimized for the real estate platform's specific requirements and will provide intelligent assistance based on the existing Go backend architecture.

### 🔥 **New Capabilities Added**
- **PostgreSQL MCP**: Direct database insights, query optimization, and performance monitoring
- **OpenAPI MCP**: Automatic type generation, API documentation, and Go→TypeScript integration
- **Enhanced workflows**: Database-aware development and type-safe frontend integration

**Next Step**: Begin FASE 2 development with role-based MCP assistance!